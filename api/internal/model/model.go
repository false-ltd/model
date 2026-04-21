package model

import (
	"database/sql"
	"time"
)

type AIModel struct {
	ID               uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	ModelID          string           `gorm:"uniqueIndex;size:255;not null" json:"model_id"`
	ProviderID       string           `gorm:"index:idx_models_provider;size:100;not null" json:"provider_id"`
	Name             string           `gorm:"size:255;not null" json:"name"`
	Family           string           `gorm:"index:idx_models_family;size:100;default:''" json:"family"`
	Attachment       bool             `gorm:"not null;default:false" json:"attachment"`
	Reasoning        bool             `gorm:"index:idx_models_reasoning;not null;default:false" json:"reasoning"`
	ToolCall         bool             `gorm:"column:tool_call;not null;default:false" json:"tool_call"`
	StructuredOutput *bool            `json:"structured_output"`
	Temperature      bool             `gorm:"not null;default:true" json:"temperature"`
	Knowledge        string           `gorm:"size:50;default:''" json:"knowledge"`
	ReleaseDate      string           `gorm:"size:50;default:''" json:"release_date"`
	LastUpdated      string           `gorm:"size:50;default:''" json:"last_updated"`
	OpenWeights      bool             `gorm:"index:idx_models_open_weights;not null;default:false" json:"open_weights"`
	Status           string           `gorm:"size:20" json:"status"`
	Interleaved      *InterleavedType `gorm:"type:json" json:"interleaved"`
	ModalitiesInput  StringSlice      `gorm:"type:json" json:"modalities_input"`
	ModalitiesOutput StringSlice      `gorm:"type:json" json:"modalities_output"`
	CostInput        *float64         `gorm:"type:decimal(10,4);index:idx_models_cost_input" json:"cost_input"`
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
	Provider         *Provider        `gorm:"foreignKey:ProviderID;references:ProviderID" json:"providers,omitempty"`
}

func (AIModel) TableName() string { return "models" }

type ModelFilter struct {
	Page        int
	PageSize    int
	Query       string
	Sort        string
	Order       string
	Providers   []string
	InputTypes  []string
	OutputTypes []string
	Reasoning       *bool
	ToolCall        *bool
	Vision          *bool
	Attachment      *bool
	OpenWeights     *bool
	StructuredOutput *bool
	Temperature     *bool
	FreeOnly        *bool
	Under1      *bool
	PriceMin       *float64
	PriceMax       *float64
	PriceOutputMin *float64
	PriceOutputMax *float64
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
	ModelCount int         `json:"modelCount"`
}

type StatsData struct {
	Stats              StatsKPIs          `json:"stats"`
	Tiers              map[string]int     `json:"tiers"`
	Capabilities       map[string]CapItem `json:"capabilities"`
	TopProviders       []TopProvider      `json:"topProviders"`
	Modalities         ModalitiesData     `json:"modalities"`
	WeightsDistribution WeightsDist       `json:"weightsDistribution"`
	ContextDistribution []ContextItem     `json:"contextDistribution"`
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
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
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

type ModelsDevResponse map[string]ModelsDevProvider

type ModelsDevProvider struct {
	ID     string                    `json:"-"`
	Name   string                    `json:"name"`
	Npm    string                    `json:"npm"`
	Env    []string                  `json:"env"`
	Doc    string                    `json:"doc"`
	Api    string                    `json:"api"`
	Models map[string]ModelsDevModel `json:"-"`
}

type ModelsDevModel struct {
	ID               string              `json:"-"`
	Name             string              `json:"name"`
	Family           string              `json:"family"`
	Attachment       bool                `json:"attachment"`
	Reasoning        bool                `json:"reasoning"`
	ToolCall         bool                `json:"tool_call"`
	StructuredOutput *bool               `json:"structured_output"`
	Temperature      *bool               `json:"temperature"`
	Knowledge        string              `json:"knowledge"`
	ReleaseDate      string              `json:"release_date"`
	LastUpdated      string              `json:"last_updated"`
	OpenWeights      bool                `json:"open_weights"`
	Interleaved      interface{}         `json:"interleaved"`
	Modalities       ModelsDevModalities `json:"modalities"`
	Cost             ModelsDevCost       `json:"cost"`
	Limit            ModelsDevLimit      `json:"limit"`
	Status           string              `json:"status"`
}

type ModelsDevModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type ModelsDevCost struct {
	Input       *float64 `json:"input"`
	Output      *float64 `json:"output"`
	CacheRead   *float64 `json:"cache_read"`
	CacheWrite  *float64 `json:"cache_write"`
	Reasoning   *float64 `json:"reasoning"`
	InputAudio  *float64 `json:"input_audio"`
	OutputAudio *float64 `json:"output_audio"`
}

type ModelsDevLimit struct {
	Context *int `json:"context"`
	Input   *int `json:"input"`
	Output  *int `json:"output"`
}

var _ sql.NullFloat64
var _ sql.NullInt64
var _ sql.NullString
var _ sql.NullBool
