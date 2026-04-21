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
	if f.Reasoning != nil && *f.Reasoning {
		qb = qb.Where("models.reasoning = true")
	}
	if f.ToolCall != nil && *f.ToolCall {
		qb = qb.Where("models.tool_call = true")
	}
	if f.Vision != nil && *f.Vision {
		qb = qb.Where("JSON_CONTAINS(models.modalities_input, '\"image\"')")
	}
	if f.Attachment != nil && *f.Attachment {
		qb = qb.Where("models.attachment = true")
	}
	if f.OpenWeights != nil && *f.OpenWeights {
		qb = qb.Where("models.open_weights = true")
	}
	if f.StructuredOutput != nil && *f.StructuredOutput {
		qb = qb.Where("models.structured_output = true")
	}
	if f.Temperature != nil && *f.Temperature {
		qb = qb.Where("models.temperature = true")
	}
	if f.FreeOnly != nil && *f.FreeOnly {
		qb = qb.Where("models.cost_input = 0")
	}
	if f.Under1 != nil && *f.Under1 {
		qb = qb.Where("models.cost_input > 0 AND models.cost_input < 1")
	}
	if f.PriceMin != nil || f.PriceMax != nil {
		qb = qb.Where(
			"models.cost_input IS NULL OR (models.cost_input >= ? AND models.cost_input <= ?)",
			floatPtrToVal(f.PriceMin, 0),
			floatPtrToVal(f.PriceMax, 100),
		)
	}
	if f.PriceOutputMin != nil || f.PriceOutputMax != nil {
		qb = qb.Where(
			"models.cost_output IS NULL OR (models.cost_output >= ? AND models.cost_output <= ?)",
			floatPtrToVal(f.PriceOutputMin, 0),
			floatPtrToVal(f.PriceOutputMax, 100),
		)
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

func (r *ModelRepo) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&model.AIModel{}).Count(&count).Error
	return count, err
}

func (r *ModelRepo) CountWhere(query string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.Model(&model.AIModel{}).Where(query, args...).Count(&count).Error
	return count, err
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

func (r *ModelRepo) FindProviderModelCounts() ([]struct {
	ProviderID   string
	ProviderName string
}, error) {
	var results []struct {
		ProviderID   string
		ProviderName string
	}
	err := r.db.Model(&model.AIModel{}).
		Select("models.provider_id, providers.name as provider_name").
		Joins("JOIN providers ON providers.provider_id = models.provider_id").
		Find(&results).Error
	return results, err
}

func (r *ModelRepo) DeleteNotInModelIDs(keepIDs []string) error {
	if len(keepIDs) == 0 {
		return nil
	}
	return r.db.Where("model_id NOT IN ?", keepIDs).Delete(&model.AIModel{}).Error
}

func (r *ModelRepo) DeleteNotInProviderIDs(keepIDs []string) error {
	if len(keepIDs) == 0 {
		return nil
	}
	return r.db.Where("provider_id NOT IN ?", keepIDs).Delete(&model.Provider{}).Error
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

func floatPtrToVal(p *float64, fallback float64) float64 {
	if p != nil {
		return *p
	}
	return fallback
}
