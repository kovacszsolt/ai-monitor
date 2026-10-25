package domain

import "time"

type AIRequest struct {
	ID              string
	ProjectID       string
	ProjectName     string
	ProviderID      string
	ProviderName    string
	RequestedModel  string
	ActualModel     string
	Endpoint        string
	Protocol        string
	StartedAt       time.Time
	CompletedAt     *time.Time
	DurationMS      *int64
	TTFTMS          *int64
	StatusCode      *int
	RequestState    string
	Streaming       bool
	InputTokens        *int64
	OutputTokens       *int64
	CachedInputTokens        *int64
	UncachedInputTokens      *int64
	PromptTPS                *float64
	GenerationTPS            *float64
	EffectiveTPS             *float64
	CacheHitRate             *float64
	GenerationSpeedEstimated bool
	PromptSpeedEstimated     bool
	UsageStatus              string
	RequestBytes    int64
	ResponseBytes   int64
	ErrorType       string
	Tag             string
}

type RequestStateCount struct {
	State string
	Count int64
}

type NamedCount struct {
	Name  string
	Count int64
}

type TokenThroughput struct {
	ModelRouteID      string
	Name              string
	InputTokens       int64
	OutputTokens      int64
	CachedInputTokens int64
	DurationMS        int64
	DecodeDurationMS  int64
	EffectiveTPS      float64
	GenerationTPS     float64
	HasGenerationTPS  bool
}

// UsageGroupStats aggregates request counts and token sums per label (proxy tag or model).
type UsageGroupStats struct {
	ModelRouteID      string
	Name              string
	RequestCount      int64
	InputTokens       int64
	OutputTokens      int64
	CachedInputTokens int64
}

// ModelUsageStats is usage per model plus aggregated throughput (when duration is available).
type ModelUsageStats struct {
	UsageGroupStats
	EffectiveTPS     float64
	GenerationTPS    float64
	HasGenerationTPS bool
	TotalCostUSD     float64
	HasTotalCost     bool
}

// ModelTokenTimePoint is one hourly bucket for a model (tag-filtered queries).
type ModelTokenTimePoint struct {
	BucketStart  string
	ModelRouteID string
	Model        string
	InputTokens  int64
	OutputTokens int64
}

// HourlyModelRequestCount is one UTC hour bucket for a public model (chart aggregates).
type HourlyModelRequestCount struct {
	BucketStart       string
	ModelRouteID      string
	Model             string
	RequestCount      int64
	InputTokens       int64
	CachedInputTokens int64
	OutputTokens      int64
	DecodeDurationMS  int64
}

type DashboardStats struct {
	TotalRequests      int64
	RequestsLast24h    int64
	SumInputTokens     int64
	SumOutputTokens    int64
	SumCachedInTokens  int64
	AvgDurationMS      float64
	StreamingCount     int64
	ByState            []RequestStateCount
	TopModels          []NamedCount
	TopProjects        []NamedCount
	UsageByTag []UsageGroupStats
}
