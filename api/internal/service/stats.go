package service

import (
	"math"

	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

type StatsService struct {
	modelRepo    *repository.ModelRepo
	providerRepo *repository.ProviderRepo
}

func NewStatsService(modelRepo *repository.ModelRepo, providerRepo *repository.ProviderRepo) *StatsService {
	return &StatsService{modelRepo: modelRepo, providerRepo: providerRepo}
}

func (s *StatsService) Get() (*model.StatsData, error) {
	totalProviders, _ := s.providerRepo.FindAll()
	totalModels, _ := s.modelRepo.CountAll()
	freeCount, _ := s.modelRepo.CountWhere("cost_input = 0 AND cost_output = 0")
	reasoningCount, _ := s.modelRepo.CountWhere("reasoning = true")
	toolCallCount, _ := s.modelRepo.CountWhere("tool_call = true")
	visionCount, _ := s.modelRepo.CountWhere("JSON_CONTAINS(modalities_input, '\"image\"')")
	attachmentCount, _ := s.modelRepo.CountWhere("attachment = true")
	temperatureCount, _ := s.modelRepo.CountWhere("temperature = true")
	openCount, _ := s.modelRepo.CountWhere("open_weights = true")

	total := int(totalModels)

	prices, _ := s.modelRepo.FindColumnValues("cost_input")
	medianPrice := medianFloat(prices)

	contexts, _ := s.modelRepo.FindIntColumnValues("limit_context")
	medianCtx := medianInt(contexts)

	allPrices, _ := s.modelRepo.FindAllPrices()
	tiers := calcTiers(allPrices)

	capabilities := map[string]model.CapItem{
		"reasoning":   {Count: int(reasoningCount), Total: total},
		"tool_call":   {Count: int(toolCallCount), Total: total},
		"vision":      {Count: int(visionCount), Total: total},
		"attachment":  {Count: int(attachmentCount), Total: total},
		"temperature": {Count: int(temperatureCount), Total: total},
	}

	provCounts, _ := s.modelRepo.FindProviderModelCounts()
	topProviders := calcTopProviders(provCounts)

	modData, _ := s.modelRepo.FindAllModalityData()
	modalities := calcModalities(modData)

	open := int(openCount)
	weights := model.WeightsDist{
		Open:    open,
		Closed:  total - open,
		OpenPct: pct(open, total),
	}

	ctxLimits, _ := s.modelRepo.FindAllContextLimits()
	ctxDist := make([]model.ContextItem, 0, len(ctxLimits))
	for _, v := range ctxLimits {
		ctxDist = append(ctxDist, model.ContextItem{LimitContext: v})
	}

	return &model.StatsData{
		Stats: model.StatsKPIs{
			TotalProviders:   len(totalProviders),
			TotalModels:      total,
			FreeModelsCount:  int(freeCount),
			FreeModelsPct:    pct(int(freeCount), total),
			MedianInputPrice: medianPrice,
			ReasoningCount:   int(reasoningCount),
			ReasoningPct:     pct(int(reasoningCount), total),
			MedianContext:    medianCtx,
		},
		Tiers:               tiers,
		Capabilities:        capabilities,
		TopProviders:        topProviders,
		Modalities:          modalities,
		WeightsDistribution: weights,
		ContextDistribution: ctxDist,
	}, nil
}

func medianFloat(arr []float64) float64 {
	n := len(arr)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return arr[n/2]
	}
	return (arr[n/2-1] + arr[n/2]) / 2
}

func medianInt(arr []int) int {
	n := len(arr)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return arr[n/2]
	}
	return (arr[n/2-1] + arr[n/2]) / 2
}

func pct(part, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(part) / float64(total) * 100))
}

func calcTiers(prices []struct {
	CostInput  *float64
	CostOutput *float64
}) map[string]int {
	tiers := map[string]int{"free": 0, "budget": 0, "standard": 0, "premium": 0, "enterprise": 0}
	for _, p := range prices {
		in := floatPtr(p.CostInput)
		out := floatPtr(p.CostOutput)
		maxPrice := math.Max(in, out)
		switch {
		case maxPrice == 0:
			tiers["free"]++
		case maxPrice < 1:
			tiers["budget"]++
		case maxPrice < 5:
			tiers["standard"]++
		case maxPrice < 15:
			tiers["premium"]++
		default:
			tiers["enterprise"]++
		}
	}
	return tiers
}

func floatPtr(p *float64) float64 {
	if p != nil {
		return *p
	}
	return 0
}

func calcTopProviders(data []struct {
	ProviderID   string
	ProviderName string
}) []model.TopProvider {
	counts := make(map[string]*model.TopProvider)
	for _, d := range data {
		if _, ok := counts[d.ProviderID]; !ok {
			counts[d.ProviderID] = &model.TopProvider{ID: d.ProviderID, Name: d.ProviderName}
		}
		counts[d.ProviderID].Count++
	}
	result := make([]model.TopProvider, 0, len(counts))
	for _, v := range counts {
		result = append(result, *v)
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Count > result[i].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func calcModalities(data []struct {
	ModalitiesInput  model.StringSlice
	ModalitiesOutput model.StringSlice
}) model.ModalitiesData {
	inputCounts := map[string]int{}
	outputCounts := map[string]int{}
	for _, d := range data {
		for _, m := range d.ModalitiesInput {
			inputCounts[m]++
		}
		for _, m := range d.ModalitiesOutput {
			outputCounts[m]++
		}
	}
	return model.ModalitiesData{
		Input:  toModSlice(inputCounts),
		Output: toModSlice(outputCounts),
	}
}

func toModSlice(m map[string]int) []model.ModalityCount {
	result := make([]model.ModalityCount, 0, len(m))
	for t, c := range m {
		result = append(result, model.ModalityCount{Type: t, Count: c})
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Count > result[i].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
