package service

import (
	"time"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

const atlasCacheKey = "atlas"
const atlasCacheTTL = 5 * time.Minute

type AtlasService struct {
	modelRepo    *repository.ModelRepo
	providerRepo *repository.ProviderRepo
	cache        *cache.Cache
}

func NewAtlasService(modelRepo *repository.ModelRepo, providerRepo *repository.ProviderRepo, c *cache.Cache) *AtlasService {
	return &AtlasService{modelRepo: modelRepo, providerRepo: providerRepo, cache: c}
}

// Get returns the atlas point cloud, cached until invalidation on data sync.
func (s *AtlasService) Get() (*model.AtlasData, error) {
	v, err := s.cache.GetOrCompute(atlasCacheKey, atlasCacheTTL, s.compute)
	if err != nil {
		return nil, err
	}
	return v.(*model.AtlasData), nil
}

func (s *AtlasService) compute() (any, error) {
	providers, err := s.providerRepo.FindAll()
	if err != nil {
		return nil, err
	}
	points, err := s.modelRepo.FindAtlasPoints()
	if err != nil {
		return nil, err
	}

	out := model.AtlasData{
		Providers: make([]model.AtlasProvider, 0, len(providers)),
		Points:    points,
	}
	for _, p := range providers {
		out.Providers = append(out.Providers, model.AtlasProvider{ID: p.ProviderID, Name: p.Name})
	}
	return &out, nil
}
