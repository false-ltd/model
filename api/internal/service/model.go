package service

import (
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

type ModelService struct {
	modelRepo *repository.ModelRepo
}

func NewModelService(modelRepo *repository.ModelRepo) *ModelService {
	return &ModelService{modelRepo: modelRepo}
}

func (s *ModelService) List(f *model.ModelFilter) (*model.ModelListResult, error) {
	return s.modelRepo.FindWithFilter(f)
}

func (s *ModelService) GetByID(id uint) (*model.AIModel, error) {
	return s.modelRepo.FindByID(id)
}

func (s *ModelService) Compare(ids []uint) (*model.CompareResult, error) {
	if len(ids) == 0 {
		return &model.CompareResult{Data: []model.AIModel{}, Count: 0}, nil
	}

	models, err := s.modelRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	orderMap := make(map[uint]int, len(ids))
	for i, id := range ids {
		orderMap[id] = i
	}

	modelMap := make(map[uint]model.AIModel, len(models))
	for _, m := range models {
		modelMap[m.ID] = m
	}

	sorted := make([]model.AIModel, 0, len(models))
	for _, id := range ids {
		if m, ok := modelMap[id]; ok {
			sorted = append(sorted, m)
		}
	}

	return &model.CompareResult{Data: sorted, Count: len(sorted)}, nil
}
