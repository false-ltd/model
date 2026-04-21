package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/false-ltd/model/api/internal/config"
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

type SyncService struct {
	providerRepo *repository.ProviderRepo
	modelRepo    *repository.ModelRepo
	cfg          *config.SyncConfig
	httpClient   *http.Client
}

func NewSyncService(providerRepo *repository.ProviderRepo, modelRepo *repository.ModelRepo, cfg *config.SyncConfig) *SyncService {
	return &SyncService{
		providerRepo: providerRepo,
		modelRepo:    modelRepo,
		cfg:          cfg,
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

	var providers []model.Provider
	var models []model.AIModel
	var keepProviderIDs []string
	var keepModelIDs []string

	for providerID, rawProvider := range raw {
		var tmp map[string]json.RawMessage
		json.Unmarshal(rawProvider, &tmp)

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
			json.Unmarshal(rawModels, &modelsMap)

			for modelID, rawModel := range modelsMap {
				m := parseModel(rawModel, providerID, modelID)
				models = append(models, m)
				keepModelIDs = append(keepModelIDs, modelID)
			}
		}
	}

	if err := s.providerRepo.UpsertAll(providers); err != nil {
		return nil, fmt.Errorf("failed to upsert providers: %w", err)
	}

	if err := s.modelRepo.UpsertAll(models); err != nil {
		return nil, fmt.Errorf("failed to upsert models: %w", err)
	}

	s.modelRepo.DeleteNotInModelIDs(keepModelIDs)
	s.providerRepo.DeleteNotInProviderIDs(keepProviderIDs)

	now := time.Now()
	s.providerRepo.UpdateSyncedAt(now)

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

	if v, ok := tmp["interleaved"]; ok {
		var val model.InterleavedType
		json.Unmarshal(v, &val.Data)
		if val.Data != nil && val.Data != false {
			m.Interleaved = &val
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
