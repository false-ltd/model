package repository

import (
	"fmt"

	"github.com/false-ltd/model/api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModelRepo struct {
	db *gorm.DB
}

func NewModelRepo(db *gorm.DB) *ModelRepo {
	return &ModelRepo{db: db}
}

// WithTx returns a shallow copy of the repo bound to a transaction.
func (r *ModelRepo) WithTx(tx *gorm.DB) *ModelRepo {
	return &ModelRepo{db: tx}
}

func (r *ModelRepo) FindByID(id uint) (*model.AIModel, error) {
	var m model.AIModel
	err := r.db.Preload("Provider").First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ModelRepo) FindByIDs(ids []uint) ([]model.AIModel, error) {
	var models []model.AIModel
	err := r.db.Preload("Provider").Where("id IN ?", ids).Find(&models).Error
	return models, err
}

func (r *ModelRepo) FindWithFilter(f *model.ModelFilter) (*model.ModelListResult, error) {
	qb := r.db.Model(&model.AIModel{}).Preload("Provider")

	if f.Query != "" {
		like := "%" + f.Query + "%"
		qb = qb.Where("models.name LIKE ? OR models.family LIKE ? OR models.model_id LIKE ?", like, like, like)
	}
	if len(f.Providers) > 0 {
		qb = qb.Where("models.provider_id IN ?", f.Providers)
	}
	if len(f.InputTypes) > 0 {
		for _, t := range f.InputTypes {
			qb = qb.Where("JSON_CONTAINS(models.modalities_input, ?)", fmt.Sprintf(`"%s"`, t))
		}
	}
	if len(f.OutputTypes) > 0 {
		for _, t := range f.OutputTypes {
			qb = qb.Where("JSON_CONTAINS(models.modalities_output, ?)", fmt.Sprintf(`"%s"`, t))
		}
	}
	if f.Reasoning != nil {
		qb = qb.Where("models.reasoning = ?", *f.Reasoning)
	}
	if f.ToolCall != nil {
		qb = qb.Where("models.tool_call = ?", *f.ToolCall)
	}
	if f.Vision != nil {
		if *f.Vision {
			qb = qb.Where("JSON_CONTAINS(models.modalities_input, '\"image\"')")
		} else {
			qb = qb.Where("models.modalities_input IS NULL OR NOT JSON_CONTAINS(models.modalities_input, '\"image\"')")
		}
	}
	if f.Attachment != nil {
		qb = qb.Where("models.attachment = ?", *f.Attachment)
	}
	if f.OpenWeights != nil {
		qb = qb.Where("models.open_weights = ?", *f.OpenWeights)
	}
	if f.StructuredOutput != nil {
		if *f.StructuredOutput {
			qb = qb.Where("models.structured_output = true")
		} else {
			qb = qb.Where("models.structured_output IS NULL OR models.structured_output = false")
		}
	}
	if f.Temperature != nil {
		qb = qb.Where("models.temperature = ?", *f.Temperature)
	}
	if f.FreeOnly != nil && *f.FreeOnly {
		qb = qb.Where("models.cost_input = 0 AND models.cost_output = 0")
	}
	if f.Under1 != nil && *f.Under1 {
		qb = qb.Where("models.cost_input > 0 AND models.cost_input < 1")
	}
	// Price ranges exclude unpriced models: a model without a known price
	// does not belong to any price range.
	if f.PriceMin != nil {
		qb = qb.Where("models.cost_input IS NOT NULL AND models.cost_input >= ?", *f.PriceMin)
	}
	if f.PriceMax != nil {
		qb = qb.Where("models.cost_input IS NOT NULL AND models.cost_input <= ?", *f.PriceMax)
	}
	if f.PriceOutputMin != nil {
		qb = qb.Where("models.cost_output IS NOT NULL AND models.cost_output >= ?", *f.PriceOutputMin)
	}
	if f.PriceOutputMax != nil {
		qb = qb.Where("models.cost_output IS NOT NULL AND models.cost_output <= ?", *f.PriceOutputMax)
	}

	var total int64
	if err := qb.Count(&total).Error; err != nil {
		return nil, err
	}

	sortCol := mapSortColumn(f.Sort)
	orderDir := "ASC"
	if f.Order == "desc" {
		orderDir = "DESC"
	}
	// Nullable columns sort NULLs first in MySQL ASC (last in DESC) — push
	// them to the end so price/limit sorting leads with real values.
	if _, nullable := nullableSortColumns[sortCol]; nullable && orderDir == "ASC" {
		qb = qb.Order(sortCol + " IS NULL")
	}
	qb = qb.Order(sortCol + " " + orderDir)

	offset := (f.Page - 1) * f.PageSize
	var data []model.AIModel
	if err := qb.Offset(offset).Limit(f.PageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	return &model.ModelListResult{Data: data, Total: total}, nil
}

func (r *ModelRepo) UpsertAll(models []model.AIModel) error {
	if len(models) == 0 {
		return nil
	}
	const batchSize = 500
	for i := 0; i < len(models); i += batchSize {
		end := i + batchSize
		if end > len(models) {
			end = len(models)
		}
		if err := r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "model_id"}, {Name: "provider_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"provider_id", "name", "family", "attachment", "reasoning",
				"tool_call", "structured_output", "temperature", "knowledge",
				"release_date", "last_updated", "open_weights", "status",
				"interleaved", "modalities_input", "modalities_output",
				"cost_input", "cost_output", "cost_reasoning",
				"cost_cache_read", "cost_cache_write",
				"cost_input_audio", "cost_output_audio",
				"limit_context", "limit_input", "limit_output", "updated_at",
			}),
		}).Create(models[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

// CountAggregates is the single-query aggregate over the whole models table
// used by the stats endpoint. Boolean SUM() expressions return NULL on an
// empty table, hence COALESCE.
type CountAggregates struct {
	Total        int64 `gorm:"column:total"`
	Free         int64 `gorm:"column:free"`
	Reasoning    int64 `gorm:"column:reasoning_count"`
	ToolCall     int64 `gorm:"column:tool_call_count"`
	Vision       int64 `gorm:"column:vision_count"`
	Attachment   int64 `gorm:"column:attachment_count"`
	Temperature  int64 `gorm:"column:temperature_count"`
	OpenWeights  int64 `gorm:"column:open_weights_count"`
}

func (r *ModelRepo) AggregateCounts() (*CountAggregates, error) {
	var agg CountAggregates
	err := r.db.Model(&model.AIModel{}).Select(`
		COUNT(*) AS total,
		COALESCE(SUM(cost_input = 0 AND cost_output = 0), 0) AS free,
		COALESCE(SUM(reasoning), 0) AS reasoning_count,
		COALESCE(SUM(tool_call), 0) AS tool_call_count,
		COALESCE(SUM(JSON_CONTAINS(modalities_input, '"image"')), 0) AS vision_count,
		COALESCE(SUM(attachment), 0) AS attachment_count,
		COALESCE(SUM(temperature), 0) AS temperature_count,
		COALESCE(SUM(open_weights), 0) AS open_weights_count
	`).Scan(&agg).Error
	return &agg, err
}

func (r *ModelRepo) FindColumnValues(column string) ([]float64, error) {
	var values []float64
	err := r.db.Model(&model.AIModel{}).
		Where(column+" IS NOT NULL AND "+column+" > 0").
		Order(column + " ASC").
		Pluck(column, &values).Error
	return values, err
}

func (r *ModelRepo) FindIntColumnValues(column string) ([]int, error) {
	var values []int
	err := r.db.Model(&model.AIModel{}).
		Where(column+" IS NOT NULL").
		Order(column + " ASC").
		Pluck(column, &values).Error
	return values, err
}

func (r *ModelRepo) FindAllPrices() ([]struct {
	CostInput  *float64
	CostOutput *float64
}, error) {
	var results []struct {
		CostInput  *float64
		CostOutput *float64
	}
	err := r.db.Model(&model.AIModel{}).
		Select("cost_input, cost_output").
		Find(&results).Error
	return results, err
}

func (r *ModelRepo) FindAllModalityData() ([]struct {
	ModalitiesInput  model.StringSlice
	ModalitiesOutput model.StringSlice
}, error) {
	var results []struct {
		ModalitiesInput  model.StringSlice
		ModalitiesOutput model.StringSlice
	}
	err := r.db.Model(&model.AIModel{}).
		Select("modalities_input, modalities_output").
		Find(&results).Error
	return results, err
}

func (r *ModelRepo) FindAllContextLimits() ([]int, error) {
	var values []int
	err := r.db.Model(&model.AIModel{}).
		Where("limit_context > 0").
		Pluck("limit_context", &values).Error
	return values, err
}

// FindTopProviders returns the providers with the most models, computed
// entirely in SQL.
func (r *ModelRepo) FindTopProviders(limit int) ([]model.TopProvider, error) {
	var results []model.TopProvider
	err := r.db.Model(&model.AIModel{}).
		Select("models.provider_id AS id, providers.name AS name, COUNT(*) AS count").
		Joins("JOIN providers ON providers.provider_id = models.provider_id").
		Group("models.provider_id, providers.name").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

// FindRecentModels returns the newest released models for the homepage
// ticker. release_date is a VARCHAR ISO date, so lexicographic order works.
func (r *ModelRepo) FindRecentModels(limit int) ([]model.RecentModel, error) {
	var out []model.RecentModel
	err := r.db.Model(&model.AIModel{}).
		Select("id, name, release_date").
		Where("release_date IS NOT NULL AND release_date <> ''").
		Order("release_date DESC, id DESC").
		Limit(limit).
		Find(&out).Error
	return out, err
}

// FindAtlasPoints streams the minimal columns needed by the atlas point
// cloud; capability count is computed in SQL to avoid shipping raw flags.
func (r *ModelRepo) FindAtlasPoints() ([]model.AtlasPoint, error) {
	var points []model.AtlasPoint
	err := r.db.Model(&model.AIModel{}).
		Select(`
			id,
			name AS n,
			model_id AS m,
			provider_id AS p,
			cost_input AS ci, cost_output AS co, limit_context AS ctx, open_weights AS ow,
			(reasoning + tool_call + attachment + open_weights + (temperature) +
			 COALESCE(structured_output, 0) +
			 COALESCE(JSON_CONTAINS(modalities_input, '"image"'), 0)) AS k
		`).
		Order("id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (r *ModelRepo) FindAllIDs() ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.AIModel{}).Pluck("id", &ids).Error
	return ids, err
}

// DeleteStaleForProvider removes models of one provider that are no longer
// in the upstream catalog. Cleanup must be scoped per (provider, model)
// pair: a global model_id NOT IN would keep stale rows whose model_id is
// still offered by a different provider.
func (r *ModelRepo) DeleteStaleForProvider(providerID string, keepModelIDs []string) error {
	if len(keepModelIDs) == 0 {
		return nil
	}
	return r.db.Where("provider_id = ? AND model_id NOT IN ?", providerID, keepModelIDs).Delete(&model.AIModel{}).Error
}

// DeleteByProviderIDsNotIn removes models belonging to providers that are
// gone from the upstream catalog (no FK, so this must precede the provider
// delete).
func (r *ModelRepo) DeleteByProviderIDsNotIn(keepProviderIDs []string) error {
	if len(keepProviderIDs) == 0 {
		return nil
	}
	return r.db.Where("provider_id NOT IN ?", keepProviderIDs).Delete(&model.AIModel{}).Error
}

func mapSortColumn(s string) string {
	columns := map[string]string{
		"name":              "models.name",
		"model_id":          "models.model_id",
		"provider_id":       "models.provider_id",
		"family":            "models.family",
		"cost_input":        "models.cost_input",
		"cost_output":       "models.cost_output",
		"cost_reasoning":    "models.cost_reasoning",
		"cost_cache_read":   "models.cost_cache_read",
		"cost_cache_write":  "models.cost_cache_write",
		"cost_input_audio":  "models.cost_input_audio",
		"cost_output_audio": "models.cost_output_audio",
		"limit_context":     "models.limit_context",
		"limit_input":       "models.limit_input",
		"limit_output":      "models.limit_output",
		"release_date":      "models.release_date",
		"last_updated":      "models.last_updated",
		"tool_call":         "models.tool_call",
		"reasoning":         "models.reasoning",
		"open_weights":      "models.open_weights",
		"structured_output": "models.structured_output",
		"attachment":        "models.attachment",
		"temperature":       "models.temperature",
	}
	if col, ok := columns[s]; ok {
		return col
	}
	return "models.name"
}

var nullableSortColumns = map[string]struct{}{
	"models.cost_input":        {},
	"models.cost_output":       {},
	"models.cost_reasoning":    {},
	"models.cost_cache_read":   {},
	"models.cost_cache_write":  {},
	"models.cost_input_audio":  {},
	"models.cost_output_audio": {},
	"models.limit_context":     {},
	"models.limit_input":       {},
	"models.limit_output":      {},
	"models.release_date":      {},
	"models.last_updated":      {},
}
