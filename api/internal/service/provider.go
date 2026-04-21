package service

import (
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

type ProviderService struct {
	providerRepo *repository.ProviderRepo
}

func NewProviderService(providerRepo *repository.ProviderRepo) *ProviderService {
	return &ProviderService{providerRepo: providerRepo}
}

func (s *ProviderService) List() ([]model.ProviderWithCount, int, error) {
	providers, err := s.providerRepo.FindAll()
	if err != nil {
		return nil, 0, err
	}

	counts, err := s.providerRepo.CountModelByProvider()
	if err != nil {
		return nil, 0, err
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

	return result, len(result), nil
}
