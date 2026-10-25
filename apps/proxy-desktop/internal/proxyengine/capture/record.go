package capture

import "time"

type RequestState string

const (
	StateStarted     RequestState = "started"
	StateCompleted   RequestState = "completed"
	StateFailed      RequestState = "failed"
	StateCancelled   RequestState = "cancelled"
	StateInterrupted RequestState = "interrupted"
)

// Progress is the capture lifecycle (in-flight vs closed).
type Progress string

const (
	ProgressStarted  Progress = "started"
	ProgressFinished Progress = "finished"
)

// PersistPhase controls how the writer persists a record.
type PersistPhase int

const (
	// PhaseOneShot is a single full INSERT (failures before upstream send).
	PhaseOneShot PersistPhase = 0
	PhaseBegin   PersistPhase = 1
	PhaseFinish  PersistPhase = 2
)

type Record struct {
	ID                string
	ProjectID         string
	ProxyAPIKeyID     string
	ProviderID        string
	ModelRouteID      string
	RequestedModel    string
	ActualModel       string
	Endpoint          string
	Protocol          string
	StartedAt         time.Time
	CompletedAt       *time.Time
	DurationMS        *int64
	UpstreamLatencyMS *int64
	TTFBMS            *int64
	TTFTMS            *int64
	StatusCode        *int
	RequestState      RequestState
	Progress          Progress
	Phase             PersistPhase
	Streaming         bool
	InputTokens       *int64
	OutputTokens      *int64
	CachedInputTokens *int64
	CacheWriteTokens  *int64
	ReasoningTokens   *int64
	UncachedInputTokens      *int64
	PrefillDurationMS        *int64
	DecodeDurationMS         *int64
	PromptTPS                *float64
	GenerationTPS            *float64
	EffectiveTPS             *float64
	CacheHitRate             *float64
	GenerationSpeedEstimated bool
	PromptSpeedEstimated     bool
	UsageStatus       string
	RawUsageJSON      string
	RequestBytes      int64
	ResponseBytes     int64
	ErrorType         string
	TimeoutFlag       bool
	CancelledFlag     bool
	Tag               string

	Prompt   *PromptCapture
	Response *ResponseCapture
	Errors   []CaptureError
}

// CaptureError is persisted in ai_request_errors when content capture fails.
type CaptureError struct {
	Stage   string // prompt | response
	Code    string
	Message string
}

// PromptCapture is persisted in ai_request_prompts (not part of ai_requests row).
type PromptCapture struct {
	LastUserText    string
	MessagesJSON    string
	PromptTruncated bool
}

// ResponseCapture is persisted in ai_request_responses.
type ResponseCapture struct {
	AssistantText      string
	ResponseTruncated  bool
	ToolCallsJSON      string
	ToolCallsTruncated bool
	RawBody            string
	RawBodyTruncated   bool
}
