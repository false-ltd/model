package service

import (
	"math"
	"sort"
	"time"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/model"
	"github.com/false-ltd/model/api/internal/repository"
)

const statsCacheKey = "stats"
const statsCacheTTL = 5 * time.Minute

type StatsService struct {
	modelRepo    *repository.ModelRepo
	providerRepo *repository.ProviderRepo
	cache        *cache.Cache
}

func NewStatsService(modelRepo *repository.ModelRepo, providerRepo *repository.ProviderRepo, c *cache.Cache) *StatsService {
	return &StatsService{modelRepo: modelRepo, providerRepo: providerRepo, cache: c}
}

// Get returns aggregate stats, cached until invalidation on data sync.
func (s *StatsService) Get() (*model.StatsData, error) {
	v, err := s.cache.GetOrCompute(statsCacheKey, statsCacheTTL, s.compute)
	if err != nil {
		return nil, err
	}
	return v.(*model.StatsData), nil
}

func (s *StatsService) compute() (any, error) {
	agg, err := s.modelRepo.AggregateCounts()
	if err != nil {
		return nil, err
	}

	providers, err := s.providerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	prices, err := s.modelRepo.FindColumnValues("cost_input")
	if err != nil {
		return nil, err
	}
	medianPrice := medianFloat(prices)

	contexts, err := s.modelRepo.FindIntColumnValues("limit_context")
	if err != nil {
		return nil, err
	}
	medianCtx := medianInt(contexts)

	allPrices, err := s.modelRepo.FindAllPrices()
	if err != nil {
		return nil, err
	}
	tiers := calcTiers(allPrices)

	capabilities := map[string]model.CapItem{
		"reasoning":   {Count: int(agg.Reasoning), Total: int(agg.Total)},
		"tool_call":   {Count: int(agg.ToolCall), Total: int(agg.Total)},
		"vision":      {Count: int(agg.Vision), Total: int(agg.Total)},
		"attachment":  {Count: int(agg.Attachment), Total: int(agg.Total)},
		"temperature": {Count: int(agg.Temperature), Total: int(agg.Total)},
	}

	topProviders, err := s.modelRepo.FindTopProviders(5)
	if err != nil {
		return nil, err
	}

	modData, err := s.modelRepo.FindAllModalityData()
	if err != nil {
		return nil, err
	}
	modalities := calcModalities(modData)

	open := int(agg.OpenWeights)
	total := int(agg.Total)
	weights := model.WeightsDist{
		Open:    open,
		Closed:  total - open,
		OpenPct: pct(open, total),
	}

	recent, err := s.modelRepo.FindRecentModels(12)
	if err != nil {
		return nil, err
	}

	ctxLimits, err := s.modelRepo.FindAllContextLimits()
	if err != nil {
		return nil, err
	}
	ctxDist := make([]model.ContextItem, 0, len(ctxLimits))
	for _, v := range ctxLimits {
		ctxDist = append(ctxDist, model.ContextItem{LimitContext: v})
	}

	return &model.StatsData{
		Stats: model.StatsKPIs{
			TotalProviders:   len(providers),
			TotalModels:      total,
			FreeModelsCount:  int(agg.Free),
			FreeModelsPct:    pct(int(agg.Free), total),
			MedianInputPrice: medianPrice,
			ReasoningCount:   int(agg.Reasoning),
			ReasoningPct:     pct(int(agg.Reasoning), total),
			MedianContext:    medianCtx,
		},
		Tiers:               tiers,
		Capabilities:        capabilities,
		TopProviders:        topProviders,
		Modalities:          modalities,
		WeightsDistribution: weights,
		ContextDistribution: ctxDist,
		Recent:              recent,
	}, nil
}

// medianFloat sorts defensively; the repo already returns ordered values.
func medianFloat(arr []float64) float64 {
	sort.Float64s(arr)
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
	sort.Ints(arr)
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
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Type < result[j].Type
	})
	return result
}
