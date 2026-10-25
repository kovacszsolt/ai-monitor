package service

import (
	"strings"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

type AboutDTO struct {
	AppName          string `json:"app_name"`
	Version          string `json:"version"`
	BuildTime        string `json:"build_time"`
	Author           string `json:"author"`
	License          string `json:"license"`
	ConfigFilePath   string `json:"config_file_path"`
	DatabasePath     string `json:"database_path"`
	ProxyConfigPath  string `json:"proxy_config_path"`
	SchemaReady      bool   `json:"schema_ready"`
	MigrationVersion uint   `json:"migration_version"`
	MigrationDirty   bool   `json:"migration_dirty"`
	ConfigBridgeOK   bool   `json:"config_bridge_ok"`
	BridgeMessage    string `json:"bridge_message,omitempty"`
}

type ConnectionInfoDTO struct {
	DatabasePath    string `json:"database_path"`
	Connected       bool   `json:"connected"`
	SchemaReady     bool   `json:"schema_ready"`
	Message         string `json:"message,omitempty"`
	BuildTime       string `json:"build_time"`
	ProxyConfigPath string `json:"proxy_config_path,omitempty"`
	ConfigBridgeOK  bool   `json:"config_bridge_ok"`
	BridgeMessage   string `json:"bridge_message,omitempty"`
	MigrationVersion uint  `json:"migration_version"`
	MigrationDirty   bool  `json:"migration_dirty"`
	ProxyRunning     bool  `json:"proxy_running"`
}

type ProxyConfigSummaryDTO struct {
	ListenAddr      string `json:"listen_addr"`
	DatabasePath    string `json:"database_path"`
	CaptureMode     string `json:"capture_mode"`
	StoreContent    bool   `json:"store_content"`
	ProxyConfigPath string `json:"proxy_config_path"`
	BridgeOK        bool   `json:"bridge_ok"`
	BridgeMessage   string `json:"bridge_message,omitempty"`
}

type StorageStatsDTO struct {
	RequestRows  int64 `json:"request_rows"`
	Inserted     int64 `json:"inserted"`
	Dropped      int64 `json:"dropped"`
	QueueLen     int   `json:"queue_len"`
	ProxyRunning bool  `json:"proxy_running"`
}

type ExportPayloadDTO struct {
	Filename      string `json:"filename"`
	ContentBase64 string `json:"content_base64"`
	ByteCount     int    `json:"byte_count"`
}

type UnhealthyProviderDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ProviderHealthCheckDTO struct {
	UnhealthyProviders []UnhealthyProviderDTO `json:"unhealthy_providers"`
}

type StartProxyResultDTO struct {
	Started            bool                   `json:"started"`
	UnhealthyProviders []UnhealthyProviderDTO `json:"unhealthy_providers"`
}

type ProxyStatusDTO struct {
	Running    bool   `json:"running"`
	ListenAddr string `json:"listen_addr,omitempty"`
	Error      string `json:"error,omitempty"`
	ConfigPath string `json:"config_path,omitempty"`
}

type ProjectDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	RequestCount    int64  `json:"request_count"`
	LastRequestAt   string `json:"last_request_at,omitempty"`
}

type ProviderDTO struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	ProviderType        string `json:"provider_type"`
	BaseURL             string `json:"base_url"`
	HasCredential       bool   `json:"has_credential"`
	Enabled             bool   `json:"enabled"`
	ConnectionTimeoutMS int    `json:"connection_timeout_ms"`
	RequestTimeoutMS    int    `json:"request_timeout_ms"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type ModelRouteDTO struct {
	ID                       string   `json:"id"`
	PublicModelName          string   `json:"public_model_name"`
	ProviderID               string   `json:"provider_id"`
	ProviderName             string   `json:"provider_name,omitempty"`
	UpstreamModelName        string   `json:"upstream_model_name"`
	Protocol                 string   `json:"protocol"`
	Enabled                  bool     `json:"enabled"`
	InputTokenUSDPer1M       *float64 `json:"input_token_usd_per_1m,omitempty"`
	CachedInputTokenUSDPer1M *float64 `json:"cached_input_token_usd_per_1m,omitempty"`
	OutputTokenUSDPer1M      *float64 `json:"output_token_usd_per_1m,omitempty"`
	CreatedAt                string   `json:"created_at"`
	UpdatedAt                string   `json:"updated_at"`
}

type APIKeyDTO struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	Name           string `json:"name"`
	APIKey         string `json:"api_key"`
	AllowAllModels bool   `json:"allow_all_models"`
	Enabled        bool   `json:"enabled"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ExpiresAt      string `json:"expires_at,omitempty"`
	LastUsedAt     string `json:"last_used_at,omitempty"`
}

type CreateAPIKeyResultDTO struct {
	Key       APIKeyDTO `json:"key"`
	Plaintext string    `json:"plaintext"`
}

type AIRequestDTO struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	ProjectName    string `json:"project_name"`
	ProviderID     string `json:"provider_id"`
	ProviderName   string `json:"provider_name"`
	RequestedModel string `json:"requested_model"`
	ActualModel    string `json:"actual_model,omitempty"`
	Endpoint       string `json:"endpoint"`
	Protocol       string `json:"protocol"`
	StartedAt      string `json:"started_at"`
	CompletedAt    string `json:"completed_at,omitempty"`
	DurationMS     int64  `json:"duration_ms"`
	HasDuration    bool   `json:"has_duration"`
	TTFTMS         int64  `json:"ttft_ms"`
	HasTTFT        bool   `json:"has_ttft"`
	StatusCode     int    `json:"status_code"`
	HasStatusCode  bool   `json:"has_status_code"`
	RequestState   string `json:"request_state"`
	Streaming      bool   `json:"streaming"`
	InputTokens    int64  `json:"input_tokens"`
	HasInputTokens bool   `json:"has_input_tokens"`
	OutputTokens          int64  `json:"output_tokens"`
	HasOutputTokens       bool   `json:"has_output_tokens"`
	CachedInputTokens        int64   `json:"cached_input_tokens"`
	HasCachedInputTokens     bool    `json:"has_cached_input_tokens"`
	UncachedInputTokens      int64   `json:"uncached_input_tokens"`
	HasUncachedInputTokens   bool    `json:"has_uncached_input_tokens"`
	PromptTPS                float64 `json:"prompt_tps"`
	HasPromptTPS             bool    `json:"has_prompt_tps"`
	PromptSpeedEstimated     bool    `json:"prompt_speed_estimated"`
	GenerationTPS            float64 `json:"generation_tps"`
	HasGenerationTPS         bool    `json:"has_generation_tps"`
	GenerationSpeedEstimated bool    `json:"generation_speed_estimated"`
	EffectiveTPS             float64 `json:"effective_tps"`
	HasEffectiveTPS          bool    `json:"has_effective_tps"`
	CacheHitRate             float64 `json:"cache_hit_rate"`
	HasCacheHitRate          bool    `json:"has_cache_hit_rate"`
	UsageStatus              string  `json:"usage_status"`
	RequestBytes   int64  `json:"request_bytes"`
	ResponseBytes  int64  `json:"response_bytes"`
	ErrorType      string `json:"error_type,omitempty"`
	Tag            string `json:"tag,omitempty"`
}

type AIRequestListDTO struct {
	Items    []AIRequestDTO `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type AIRequestTagOptionsDTO struct {
	Enabled     bool     `json:"enabled"`
	Tags        []string `json:"tags"`
	HasUntagged bool     `json:"has_untagged"`
	NoTagLabel  string   `json:"no_tag_label"`
}

type BulkSetAIRequestTagResultDTO struct {
	Updated int64 `json:"updated"`
}

type BulkDeleteAIRequestResultDTO struct {
	Deleted int64 `json:"deleted"`
}

type AIRequestCaptureErrorDTO struct {
	Stage        string `json:"stage"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type AIRequestContentDTO struct {
	RequestID      string `json:"request_id"`
	RequestState   string `json:"request_state,omitempty"`
	ProxyError     string `json:"proxy_error,omitempty"`
	CaptureEnabled bool   `json:"capture_enabled"`
	HasContent     bool   `json:"has_content"`

	LastUserText    string `json:"last_user_text,omitempty"`
	MessagesJSON    string `json:"messages_json,omitempty"`
	PromptTruncated bool   `json:"prompt_truncated"`

	AssistantText      string `json:"assistant_text,omitempty"`
	ResponseTruncated  bool   `json:"response_truncated"`
	ToolCallsJSON      string `json:"tool_calls_json,omitempty"`
	ToolCallsTruncated bool   `json:"tool_calls_truncated"`
	RawBody            string `json:"raw_body,omitempty"`
	RawBodyTruncated   bool   `json:"raw_body_truncated"`

	CaptureErrors []AIRequestCaptureErrorDTO `json:"capture_errors,omitempty"`
}

type StateCountDTO struct {
	State string `json:"state"`
	Count int64  `json:"count"`
}

type NamedCountDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TokenThroughputDTO struct {
	Name              string  `json:"name"`
	InputTokens       int64   `json:"input_tokens"`
	OutputTokens      int64   `json:"output_tokens"`
	CachedInputTokens int64   `json:"cached_input_tokens"`
	DurationMS        int64   `json:"duration_ms"`
	DecodeDurationMS  int64   `json:"decode_duration_ms"`
	EffectiveTPS      float64 `json:"effective_tps"`
	GenerationTPS     float64 `json:"generation_tps"`
	HasGenerationTPS  bool    `json:"has_generation_tps"`
}

type UsageGroupStatsDTO struct {
	Name              string `json:"name"`
	RequestCount      int64  `json:"request_count"`
	InputTokens       int64  `json:"input_tokens"`
	OutputTokens      int64  `json:"output_tokens"`
	CachedInputTokens int64  `json:"cached_input_tokens"`
}

type ModelUsageStatsDTO struct {
	Name              string  `json:"name"`
	RequestCount      int64   `json:"request_count"`
	InputTokens       int64   `json:"input_tokens"`
	OutputTokens      int64   `json:"output_tokens"`
	CachedInputTokens int64   `json:"cached_input_tokens"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	HasTotalCost      bool    `json:"has_total_cost"`
	EffectiveTPS      float64 `json:"effective_tps"`
	GenerationTPS     float64 `json:"generation_tps"`
	HasGenerationTPS  bool    `json:"has_generation_tps"`
}

type ModelTokenTimePointDTO struct {
	BucketStart  string `json:"bucket_start"`
	Model        string `json:"model"`
	InputTokens  int64  `json:"input_tokens"`
	OutputTokens int64  `json:"output_tokens"`
}

type ModelTokenTimeSeriesDTO struct {
	BucketUnit string                   `json:"bucket_unit"`
	Days       int                      `json:"days"`
	Points     []ModelTokenTimePointDTO `json:"points"`
}

type HourlyModelRequestCountDTO struct {
	BucketStart       string  `json:"bucket_start"`
	Model             string  `json:"model"`
	RequestCount      int64   `json:"request_count"`
	InputTokens       int64   `json:"input_tokens"`
	CachedInputTokens int64   `json:"cached_input_tokens"`
	OutputTokens      int64   `json:"output_tokens"`
	TotalCostUSD      float64 `json:"total_cost_usd"`
	HasTotalCost      bool    `json:"has_total_cost"`
	GenerationTPS     float64 `json:"generation_tps"`
	HasGenerationTPS  bool    `json:"has_generation_tps"`
}

type HourlyModelRequestSeriesDTO struct {
	BucketUnit string                       `json:"bucket_unit"`
	Days       int                          `json:"days"`
	RangeStart string                       `json:"range_start"`
	RangeEnd   string                       `json:"range_end"`
	Points     []HourlyModelRequestCountDTO `json:"points"`
}

type DashboardStatsDTO struct {
	TotalRequests      int64                `json:"total_requests"`
	RequestsLast24h    int64                `json:"requests_last_24h"`
	SumInputTokens     int64                `json:"sum_input_tokens"`
	SumOutputTokens    int64                `json:"sum_output_tokens"`
	SumCachedInTokens  int64                `json:"sum_cached_input_tokens"`
	AvgDurationMS      float64              `json:"avg_duration_ms"`
	StreamingCount     int64                `json:"streaming_count"`
	ByState            []StateCountDTO      `json:"by_state"`
	TopModels          []NamedCountDTO      `json:"top_models"`
	TopProjects        []NamedCountDTO      `json:"top_projects"`
	UsageByTag []UsageGroupStatsDTO `json:"usage_by_tag"`
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func mapProject(p domain.Project) ProjectDTO {
	dto := ProjectDTO{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Enabled:      p.Enabled,
		CreatedAt:    formatTime(p.CreatedAt),
		UpdatedAt:    formatTime(p.UpdatedAt),
		RequestCount: p.RequestCount,
	}
	if p.LastRequestAt != nil {
		dto.LastRequestAt = formatTime(*p.LastRequestAt)
	}
	return dto
}

func mapProvider(p domain.Provider) ProviderDTO {
	return ProviderDTO{
		ID:                  p.ID,
		Name:                p.Name,
		ProviderType:        string(p.ProviderType),
		BaseURL:             p.BaseURL,
		HasCredential:       strings.TrimSpace(p.Credential) != "",
		Enabled:             p.Enabled,
		ConnectionTimeoutMS: p.ConnectionTimeoutMS,
		RequestTimeoutMS:    p.RequestTimeoutMS,
		CreatedAt:           formatTime(p.CreatedAt),
		UpdatedAt:           formatTime(p.UpdatedAt),
	}
}

func mapModelRoute(r domain.ModelRoute, providerName string) ModelRouteDTO {
	return ModelRouteDTO{
		ID:                       r.ID,
		PublicModelName:          r.PublicModelName,
		ProviderID:               r.ProviderID,
		ProviderName:             providerName,
		UpstreamModelName:        r.UpstreamModelName,
		Protocol:                 string(r.Protocol),
		Enabled:                  r.Enabled,
		InputTokenUSDPer1M:       r.InputTokenUSDPer1M,
		CachedInputTokenUSDPer1M: r.CachedInputTokenUSDPer1M,
		OutputTokenUSDPer1M:      r.OutputTokenUSDPer1M,
		CreatedAt:                formatTime(r.CreatedAt),
		UpdatedAt:                formatTime(r.UpdatedAt),
	}
}

func mapAIRequestContent(c domain.AIRequestContent) AIRequestContentDTO {
	dto := AIRequestContentDTO{
		RequestID:          c.RequestID,
		RequestState:       c.RequestState,
		ProxyError:         c.ProxyError,
		CaptureEnabled:     c.CaptureEnabled,
		HasContent:         c.HasContent,
		LastUserText:       c.LastUserText,
		MessagesJSON:       c.MessagesJSON,
		PromptTruncated:    c.PromptTruncated,
		AssistantText:      c.AssistantText,
		ResponseTruncated:  c.ResponseTruncated,
		ToolCallsJSON:      c.ToolCallsJSON,
		ToolCallsTruncated: c.ToolCallsTruncated,
		RawBody:            c.RawBody,
		RawBodyTruncated:   c.RawBodyTruncated,
	}
	for _, e := range c.CaptureErrors {
		dto.CaptureErrors = append(dto.CaptureErrors, AIRequestCaptureErrorDTO{
			Stage:        e.Stage,
			ErrorCode:    e.ErrorCode,
			ErrorMessage: e.ErrorMessage,
		})
	}
	return dto
}

func mapAIRequest(r domain.AIRequest) AIRequestDTO {
	dto := AIRequestDTO{
		ID:             r.ID,
		ProjectID:      r.ProjectID,
		ProjectName:    r.ProjectName,
		ProviderID:     r.ProviderID,
		ProviderName:   r.ProviderName,
		RequestedModel: r.RequestedModel,
		ActualModel:    r.ActualModel,
		Endpoint:       r.Endpoint,
		Protocol:       r.Protocol,
		StartedAt:      formatTime(r.StartedAt),
		RequestState:   r.RequestState,
		Streaming:      r.Streaming,
		UsageStatus:    r.UsageStatus,
		RequestBytes:   r.RequestBytes,
		ResponseBytes:  r.ResponseBytes,
		ErrorType:      r.ErrorType,
		Tag:            r.Tag,
	}
	if r.CompletedAt != nil {
		dto.CompletedAt = formatTime(*r.CompletedAt)
	}
	if r.DurationMS != nil {
		dto.DurationMS = *r.DurationMS
		dto.HasDuration = true
	}
	if r.TTFTMS != nil {
		dto.TTFTMS = *r.TTFTMS
		dto.HasTTFT = true
	}
	if r.StatusCode != nil {
		dto.StatusCode = *r.StatusCode
		dto.HasStatusCode = true
	}
	if r.InputTokens != nil {
		dto.InputTokens = *r.InputTokens
		dto.HasInputTokens = true
	}
	if r.OutputTokens != nil {
		dto.OutputTokens = *r.OutputTokens
		dto.HasOutputTokens = true
	}
	if r.CachedInputTokens != nil {
		dto.CachedInputTokens = *r.CachedInputTokens
		dto.HasCachedInputTokens = true
	}
	if r.UncachedInputTokens != nil {
		dto.UncachedInputTokens = *r.UncachedInputTokens
		dto.HasUncachedInputTokens = true
	}
	if r.PromptTPS != nil {
		dto.PromptTPS = *r.PromptTPS
		dto.HasPromptTPS = true
	}
	if r.GenerationTPS != nil {
		dto.GenerationTPS = *r.GenerationTPS
		dto.HasGenerationTPS = true
	}
	if r.EffectiveTPS != nil {
		dto.EffectiveTPS = *r.EffectiveTPS
		dto.HasEffectiveTPS = true
	}
	if r.CacheHitRate != nil {
		dto.CacheHitRate = *r.CacheHitRate
		dto.HasCacheHitRate = true
	}
	dto.GenerationSpeedEstimated = r.GenerationSpeedEstimated
	dto.PromptSpeedEstimated = r.PromptSpeedEstimated
	return dto
}

func mapTokenThroughput(t domain.TokenThroughput) TokenThroughputDTO {
	return TokenThroughputDTO{
		Name:              t.Name,
		InputTokens:       t.InputTokens,
		OutputTokens:      t.OutputTokens,
		CachedInputTokens: t.CachedInputTokens,
		DurationMS:        t.DurationMS,
		DecodeDurationMS:  t.DecodeDurationMS,
		EffectiveTPS:      t.EffectiveTPS,
		GenerationTPS:     t.GenerationTPS,
		HasGenerationTPS:  t.HasGenerationTPS,
	}
}

func mapUsageGroupStats(g domain.UsageGroupStats) UsageGroupStatsDTO {
	return UsageGroupStatsDTO{
		Name:              g.Name,
		RequestCount:      g.RequestCount,
		InputTokens:       g.InputTokens,
		OutputTokens:      g.OutputTokens,
		CachedInputTokens: g.CachedInputTokens,
	}
}

func mapModelUsageStats(m domain.ModelUsageStats) ModelUsageStatsDTO {
	return ModelUsageStatsDTO{
		Name:              m.Name,
		RequestCount:      m.RequestCount,
		InputTokens:       m.InputTokens,
		OutputTokens:      m.OutputTokens,
		CachedInputTokens: m.CachedInputTokens,
		TotalCostUSD:      m.TotalCostUSD,
		HasTotalCost:      m.HasTotalCost,
		EffectiveTPS:      m.EffectiveTPS,
		GenerationTPS:     m.GenerationTPS,
		HasGenerationTPS:  m.HasGenerationTPS,
	}
}

func mapDashboard(d domain.DashboardStats) DashboardStatsDTO {
	out := DashboardStatsDTO{
		TotalRequests:     d.TotalRequests,
		RequestsLast24h:   d.RequestsLast24h,
		SumInputTokens:    d.SumInputTokens,
		SumOutputTokens:   d.SumOutputTokens,
		SumCachedInTokens: d.SumCachedInTokens,
		AvgDurationMS:     d.AvgDurationMS,
		StreamingCount:    d.StreamingCount,
	}
	for _, s := range d.ByState {
		out.ByState = append(out.ByState, StateCountDTO{State: s.State, Count: s.Count})
	}
	for _, m := range d.TopModels {
		out.TopModels = append(out.TopModels, NamedCountDTO{Name: m.Name, Count: m.Count})
	}
	for _, p := range d.TopProjects {
		out.TopProjects = append(out.TopProjects, NamedCountDTO{Name: p.Name, Count: p.Count})
	}
	for _, g := range d.UsageByTag {
		out.UsageByTag = append(out.UsageByTag, mapUsageGroupStats(g))
	}
	return out
}

func mapAPIKey(k domain.APIKey) APIKeyDTO {
	return APIKeyDTO{
		ID:             k.ID,
		ProjectID:      k.ProjectID,
		Name:           k.Name,
		APIKey:         k.Secret,
		AllowAllModels: k.AllowAllModels,
		Enabled:        k.Enabled,
		CreatedAt:      formatTime(k.CreatedAt),
		UpdatedAt:      formatTime(k.UpdatedAt),
		ExpiresAt:      formatTimePtr(k.ExpiresAt),
		LastUsedAt:     formatTimePtr(k.LastUsedAt),
	}
}
