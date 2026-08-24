package model

import (
	"time"
)

type AIModel struct {
	ID               uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	ModelID          string           `gorm:"size:255;not null;uniqueIndex:uk_model_provider,priority:1" json:"model_id"`
	ProviderID       string           `gorm:"size:100;not null;index:idx_models_provider;uniqueIndex:uk_model_provider,priority:2" json:"provider_id"`
	Name             string           `gorm:"size:255;not null" json:"name"`
	Family           string           `gorm:"size:100;default:''" json:"family"`
	Attachment       bool             `gorm:"not null;default:false" json:"attachment"`
	Reasoning        bool             `gorm:"not null;default:false" json:"reasoning"`
	ToolCall         bool             `gorm:"column:tool_call;not null;default:false" json:"tool_call"`
	StructuredOutput *bool            `json:"structured_output"`
	Temperature      bool             `gorm:"not null;default:true" json:"temperature"`
	Knowledge        string           `gorm:"size:50;default:''" json:"knowledge"`
	ReleaseDate      string           `gorm:"size:50;default:''" json:"release_date"`
	LastUpdated      string           `gorm:"size:50;default:''" json:"last_updated"`
	OpenWeights      bool             `gorm:"not null;default:false" json:"open_weights"`
	Status           string           `gorm:"size:20" json:"status"`
	Interleaved      *InterleavedType `gorm:"type:json" json:"interleaved"`
	ModalitiesInput  StringSlice      `gorm:"type:json" json:"modalities_input"`
	ModalitiesOutput StringSlice      `gorm:"type:json" json:"modalities_output"`
	CostInput        *float64         `gorm:"type:decimal(10,4)" json:"cost_input"`
	CostOutput       *float64         `gorm:"type:decimal(10,4)" json:"cost_output"`
	CostReasoning    *float64         `gorm:"type:decimal(10,4)" json:"cost_reasoning"`
	CostCacheRead    *float64         `gorm:"type:decimal(10,4)" json:"cost_cache_read"`
	CostCacheWrite   *float64         `gorm:"type:decimal(10,4)" json:"cost_cache_write"`
	CostInputAudio   *float64         `gorm:"type:decimal(10,4)" json:"cost_input_audio"`
	CostOutputAudio  *float64         `gorm:"type:decimal(10,4)" json:"cost_output_audio"`
	LimitContext     *int             `gorm:"index:idx_models_context" json:"limit_context"`
	LimitInput       *int             `json:"limit_input"`
	LimitOutput      *int             `json:"limit_output"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Provider         *Provider        `gorm:"foreignKey:ProviderID;references:ProviderID" json:"provider,omitempty"`
}

func (AIModel) TableName() string { return "models" }

type ModelFilter struct {
	Page              int
	PageSize          int
	Query             string
	Sort              string
	Order             string
	Providers         []string
	InputTypes        []string
	OutputTypes       []string
	Reasoning         *bool
	ToolCall          *bool
	Vision            *bool
	Attachment        *bool
	OpenWeights       *bool
	StructuredOutput  *bool
	Temperature       *bool
	FreeOnly          *bool
	Under1            *bool
	PriceMin          *float64
	PriceMax          *float64
	PriceOutputMin    *float64
	PriceOutputMax    *float64
}

type ModelListResult struct {
	Data  []AIModel
	Total int64
}

type CompareResult struct {
	Data  []AIModel
	Count int
}

type ProviderWithCount struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Npm        string      `json:"npm"`
	Env        StringSlice `json:"env"`
	DocURL     string      `json:"doc_url"`
	ApiURL     string      `json:"api_url"`
	ModelCount int         `json:"model_count"`
}

type StatsData struct {
	Stats               StatsKPIs          `json:"stats"`
	Tiers               map[string]int     `json:"tiers"`
	Capabilities        map[string]CapItem `json:"capabilities"`
	TopProviders        []TopProvider      `json:"topProviders"`
	Modalities          ModalitiesData     `json:"modalities"`
	WeightsDistribution WeightsDist        `json:"weightsDistribution"`
	ContextDistribution []ContextItem      `json:"contextDistribution"`
	Recent              []RecentModel     `json:"recent"`
}

type StatsKPIs struct {
	TotalProviders   int     `json:"totalProviders"`
	TotalModels      int     `json:"totalModels"`
	FreeModelsCount  int     `json:"freeModelsCount"`
	FreeModelsPct    int     `json:"freeModelsPct"`
	MedianInputPrice float64 `json:"medianInputPrice"`
	ReasoningCount   int     `json:"reasoningCount"`
	ReasoningPct     int     `json:"reasoningPct"`
	MedianContext    int     `json:"medianContext"`
}

type CapItem struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

type TopProvider struct {
	ID    string `gorm:"column:id" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Count int    `gorm:"column:count" json:"count"`
}

type ModalitiesData struct {
	Input  []ModalityCount `json:"input"`
	Output []ModalityCount `json:"output"`
}

type ModalityCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type WeightsDist struct {
	Open    int `json:"open"`
	Closed  int `json:"closed"`
	OpenPct int `json:"openPct"`
}

type ContextItem struct {
	LimitContext int `json:"limit_context"`
}

type SyncResult struct {
	SyncedAt  string `json:"synced_at"`
	Skipped   bool   `json:"skipped,omitempty"`
	Providers int    `json:"providers,omitempty"`
	Models    int    `json:"models,omitempty"`
}

type SyncStatus struct {
	SyncedAt *string `json:"synced_at"`
}
