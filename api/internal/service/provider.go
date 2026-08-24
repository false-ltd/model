package service

import (
	"time"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

const providersCacheKey = "providers:list"
const providersCacheTTL = 5 * time.Minute

type ProviderService struct {
	providerRepo *repository.ProviderRepo
	cache        *cache.Cache
}

func NewProviderService(providerRepo *repository.ProviderRepo, c *cache.Cache) *ProviderService {
	return &ProviderService{providerRepo: providerRepo, cache: c}
}

// List returns providers with model counts, cached until invalidation on
// data sync.
func (s *ProviderService) List() ([]model.ProviderWithCount, int, error) {
	v, err := s.cache.GetOrCompute(providersCacheKey, providersCacheTTL, func() (any, error) {
		providers, err := s.providerRepo.FindAll()
		if err != nil {
			return nil, err
		}

		counts, err := s.providerRepo.CountModelByProvider()
		if err != nil {
			return nil, err
		}

		result := make([]model.ProviderWithCount, 0, len(providers))
		for _, p := range providers {
			result = append(result, model.ProviderWithCount{
				ID:         p.ProviderID,
				Name:       p.Name,
				Npm:        p.Npm,
				Env:        p.Env,
				DocURL:     p.DocURL,
				ApiURL:     p.ApiURL,
				ModelCount: counts[p.ProviderID],
			})
		}
		return result, nil
	})
	if err != nil {
		return nil, 0, err
	}
	result := v.([]model.ProviderWithCount)
	return result, len(result), nil
}
