package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/config"
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

// ErrSyncInProgress is returned when a sync is already running; the caller
// should report it as a conflict rather than an internal error.
var ErrSyncInProgress = errors.New("sync already in progress")

type SyncService struct {
	providerRepo *repository.ProviderRepo
	modelRepo    *repository.ModelRepo
	db           *gorm.DB
	cfg          *config.SyncConfig
	cache        *cache.Cache
	httpClient   *http.Client
	mu           sync.Mutex
}

func NewSyncService(providerRepo *repository.ProviderRepo, modelRepo *repository.ModelRepo, db *gorm.DB, cfg *config.SyncConfig, c *cache.Cache) *SyncService {
	return &SyncService{
		providerRepo: providerRepo,
		modelRepo:    modelRepo,
		db:           db,
		cfg:          cfg,
		cache:        c,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *SyncService) GetStatus() (*model.SyncStatus, error) {
	ts, err := s.providerRepo.GetLatestSyncedAt()
	if err != nil {
		return nil, err
	}
	if ts == nil {
		return &model.SyncStatus{SyncedAt: nil}, nil
	}
	str := ts.Format(time.RFC3339)
	return &model.SyncStatus{SyncedAt: &str}, nil
}

func (s *SyncService) Trigger() (*model.SyncResult, error) {
	// Serialize triggers: cron and manual POST /sync can race, and without
	// the lock both would pass the cooldown check and double-write.
	if !s.mu.TryLock() {
		return nil, ErrSyncInProgress
	}
	defer s.mu.Unlock()

	latest, err := s.providerRepo.GetLatestSyncedAt()
	if err != nil {
		return nil, err
	}
	if latest != nil {
		elapsed := time.Since(*latest)
		if elapsed < time.Duration(s.cfg.CooldownMinutes)*time.Minute {
			str := latest.Format(time.RFC3339)
			return &model.SyncResult{SyncedAt: str, Skipped: true}, nil
		}
	}

	resp, err := s.httpClient.Get(s.cfg.ModelsDevURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models.dev: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("models.dev returned status %d", resp.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse models.dev JSON: %w", err)
	}
	if len(raw) == 0 {
		return nil, errors.New("models.dev returned no providers; refusing to wipe local data")
	}

	var providers []model.Provider
	var models []model.AIModel
	var keepProviderIDs []string
	modelsByProvider := map[string][]string{}

	for providerID, rawProvider := range raw {
		var tmp map[string]json.RawMessage
		if err := json.Unmarshal(rawProvider, &tmp); err != nil {
			return nil, fmt.Errorf("failed to parse provider %q: %w", providerID, err)
		}

		p := struct {
			Name string
			Npm  string
			Doc  string
			Api  string
			Env  []string
		}{
			Name: extractString(tmp, "name"),
			Npm:  extractString(tmp, "npm"),
			Doc:  extractString(tmp, "doc"),
			Api:  extractString(tmp, "api"),
		}
		if raw, ok := tmp["env"]; ok {
			json.Unmarshal(raw, &p.Env)
		}

		env := model.StringSlice{}
		if p.Env != nil {
			env = model.StringSlice(p.Env)
		}

		providers = append(providers, model.Provider{
			ProviderID: providerID,
			Name:       p.Name,
			Npm:        p.Npm,
			Env:        env,
			DocURL:     p.Doc,
			ApiURL:     p.Api,
		})
		keepProviderIDs = append(keepProviderIDs, providerID)

		if rawModels, ok := tmp["models"]; ok {
			var modelsMap map[string]json.RawMessage
			if err := json.Unmarshal(rawModels, &modelsMap); err != nil {
				return nil, fmt.Errorf("failed to parse models of provider %q: %w", providerID, err)
			}

			for modelID, rawModel := range modelsMap {
				m := parseModel(rawModel, providerID, modelID)
				models = append(models, m)
				modelsByProvider[providerID] = append(modelsByProvider[providerID], modelID)
			}
		}
	}

	if len(models) == 0 {
		return nil, errors.New("models.dev returned no models; refusing to wipe local data")
	}

	now := time.Now()

	// All writes in one transaction: either the dataset is fully replaced
	// or the previous state stays intact and readable.
	err = s.db.Transaction(func(tx *gorm.DB) error {
		ptx := s.providerRepo.WithTx(tx)
		mtx := s.modelRepo.WithTx(tx)

		if err := ptx.UpsertAll(providers); err != nil {
			return fmt.Errorf("failed to upsert providers: %w", err)
		}
		if err := mtx.UpsertAll(models); err != nil {
			return fmt.Errorf("failed to upsert models: %w", err)
		}
		for providerID, keepModelIDs := range modelsByProvider {
			if err := mtx.DeleteStaleForProvider(providerID, keepModelIDs); err != nil {
				return fmt.Errorf("failed to delete removed models of provider %q: %w", providerID, err)
			}
		}
		if err := mtx.DeleteByProviderIDsNotIn(keepProviderIDs); err != nil {
			return fmt.Errorf("failed to delete models of removed providers: %w", err)
		}
		if err := ptx.DeleteNotInProviderIDs(keepProviderIDs); err != nil {
			return fmt.Errorf("failed to delete removed providers: %w", err)
		}
		if err := ptx.UpdateSyncedAt(now); err != nil {
			return fmt.Errorf("failed to record sync time: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.cache.Invalidate()

	return &model.SyncResult{
		SyncedAt:  now.Format(time.RFC3339),
		Providers: len(providers),
		Models:    len(models),
	}, nil
}

func parseModel(raw json.RawMessage, providerID, modelID string) model.AIModel {
	var tmp map[string]json.RawMessage
	json.Unmarshal(raw, &tmp)

	m := model.AIModel{
		ModelID:     modelID,
		ProviderID:  providerID,
		Name:        extractString(tmp, "name"),
		Family:      extractString(tmp, "family"),
		Attachment:  extractBool(tmp, "attachment"),
		Reasoning:   extractBool(tmp, "reasoning"),
		ToolCall:    extractBool(tmp, "tool_call"),
		Temperature: func() bool {
			if v, ok := tmp["temperature"]; ok {
				var b *bool
				json.Unmarshal(v, &b)
				if b != nil {
					return *b
				}
			}
			return true
		}(),
		OpenWeights: extractBool(tmp, "open_weights"),
		Status:      extractStringDefault(tmp, "status", ""),
		Knowledge:   extractStringDefault(tmp, "knowledge", ""),
		ReleaseDate: extractStringDefault(tmp, "release_date", ""),
		LastUpdated: extractStringDefault(tmp, "last_updated", ""),
	}

	if v, ok := tmp["structured_output"]; ok {
		var b *bool
		json.Unmarshal(v, &b)
		m.StructuredOutput = b
	}

	// Upstream encodes interleaved as `false` when unsupported; treat
	// null/false as absent and keep everything else (bool/string/object).
	if v, ok := tmp["interleaved"]; ok {
		trimmed := bytes.TrimSpace(v)
		if len(trimmed) > 0 && string(trimmed) != "null" && string(trimmed) != "false" {
			var val model.InterleavedType
			if err := json.Unmarshal(v, &val.Data); err == nil && val.Data != nil {
				m.Interleaved = &val
			}
		}
	}

	if v, ok := tmp["modalities"]; ok {
		var mod struct {
			Input  model.StringSlice `json:"input"`
			Output model.StringSlice `json:"output"`
		}
		json.Unmarshal(v, &mod)
		m.ModalitiesInput = mod.Input
		m.ModalitiesOutput = mod.Output
	}

	if v, ok := tmp["cost"]; ok {
		var cost struct {
			Input       *float64 `json:"input"`
			Output      *float64 `json:"output"`
			Reasoning   *float64 `json:"reasoning"`
			CacheRead   *float64 `json:"cache_read"`
			CacheWrite  *float64 `json:"cache_write"`
			InputAudio  *float64 `json:"input_audio"`
			OutputAudio *float64 `json:"output_audio"`
		}
		json.Unmarshal(v, &cost)
		m.CostInput = cost.Input
		m.CostOutput = cost.Output
		m.CostReasoning = cost.Reasoning
		m.CostCacheRead = cost.CacheRead
		m.CostCacheWrite = cost.CacheWrite
		m.CostInputAudio = cost.InputAudio
		m.CostOutputAudio = cost.OutputAudio
	}

	if v, ok := tmp["limit"]; ok {
		var limit struct {
			Context *int `json:"context"`
			Input   *int `json:"input"`
			Output  *int `json:"output"`
		}
		json.Unmarshal(v, &limit)
		m.LimitContext = limit.Context
		m.LimitInput = limit.Input
		m.LimitOutput = limit.Output
	}

	return m
}

func extractString(tmp map[string]json.RawMessage, key string) string {
	if v, ok := tmp[key]; ok {
		var s string
		json.Unmarshal(v, &s)
		return s
	}
	return ""
}

func extractStringDefault(tmp map[string]json.RawMessage, key, def string) string {
	s := extractString(tmp, key)
	if s == "" {
		return def
	}
	return s
}

func extractBool(tmp map[string]json.RawMessage, key string) bool {
	if v, ok := tmp[key]; ok {
		var b bool
		json.Unmarshal(v, &b)
		return b
	}
	return false
}
