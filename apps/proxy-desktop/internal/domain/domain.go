package domain

import "time"

type Project struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RequestCount   int64
	LastRequestAt *time.Time
}

type ProviderType string

const (
	ProviderOpenAI           ProviderType = "openai"
	ProviderAnthropic        ProviderType = "anthropic"
	ProviderOpenAICompatible ProviderType = "openai-compatible"
	ProviderBifrost          ProviderType = "bifrost"
)

type Provider struct {
	ID                  string
	Name                string
	ProviderType        ProviderType
	BaseURL             string
	Credential          string
	Enabled             bool
	ConnectionTimeoutMS int
	RequestTimeoutMS    int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Protocol string

const (
	ProtocolOpenAIChat        Protocol = "openai_chat"
	ProtocolOpenAIResponses   Protocol = "openai_responses"
	ProtocolAnthropicMessages Protocol = "anthropic_messages"
)

type ModelRoute struct {
	ID                       string
	PublicModelName          string
	ProviderID               string
	UpstreamModelName        string
	Protocol                 Protocol
	Enabled                  bool
	InputTokenUSDPer1M       *float64
	CachedInputTokenUSDPer1M *float64
	OutputTokenUSDPer1M      *float64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type APIKey struct {
	ID             string
	ProjectID      string
	Name           string
	Secret         string
	AllowAllModels bool
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      *time.Time
	LastUsedAt     *time.Time
}

type ChatPinnedPrompt struct {
	ID        string
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
