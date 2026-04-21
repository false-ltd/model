package repository

import (
	"time"

	"github.com/false-ltd/model/api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProviderRepo struct {
	db *gorm.DB
}

func NewProviderRepo(db *gorm.DB) *ProviderRepo {
	return &ProviderRepo{db: db}
}

func (r *ProviderRepo) FindAll() ([]model.Provider, error) {
	var providers []model.Provider
	err := r.db.Order("name ASC").Find(&providers).Error
	return providers, err
}

func (r *ProviderRepo) UpsertAll(providers []model.Provider) error {
	if len(providers) == 0 {
		return nil
	}
	const batchSize = 500
	for i := 0; i < len(providers); i += batchSize {
		end := i + batchSize
		if end > len(providers) {
			end = len(providers)
		}
		if err := r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "provider_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "npm", "env", "doc_url", "api_url", "updated_at"}),
		}).Create(providers[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProviderRepo) UpdateSyncedAt(t time.Time) error {
	return r.db.Model(&model.Provider{}).Where("1 = 1").Update("synced_at", t).Error
}

func (r *ProviderRepo) GetLatestSyncedAt() (*time.Time, error) {
	var provider model.Provider
	err := r.db.Where("synced_at IS NOT NULL").Order("synced_at DESC").Limit(1).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return provider.SyncedAt, nil
}

func (r *ProviderRepo) CountModelByProvider() (map[string]int, error) {
	type result struct {
		ProviderID string
		Count      int
	}
	var results []result
	err := r.db.Model(&model.AIModel{}).
		Select("provider_id, count(*) as count").
		Group("provider_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]int, len(results))
	for _, r := range results {
		m[r.ProviderID] = r.Count
	}
	return m, nil
}

func (r *ProviderRepo) DeleteNotInProviderIDs(keepIDs []string) error {
	if len(keepIDs) == 0 {
		return nil
	}
	return r.db.Where("provider_id NOT IN ?", keepIDs).Delete(&model.Provider{}).Error
}
