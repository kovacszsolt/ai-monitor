package routing

import (
	"context"
	"database/sql"

	iroute "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

type Route = iroute.Route
type Protocol = iroute.Protocol
type Modality = iroute.Modality

var ErrNotFound = iroute.ErrNotFound

const (
	ModalityTextText      = iroute.ModalityTextText
	ModalityTextImage     = iroute.ModalityTextImage
	ModalityTextImageText = iroute.ModalityTextImageText

	ProtocolOpenAIChat        = iroute.ProtocolOpenAIChat
	ProtocolOpenAIResponses   = iroute.ProtocolOpenAIResponses
	ProtocolOpenAIImages      = iroute.ProtocolOpenAIImages
	ProtocolAnthropicMessages = iroute.ProtocolAnthropicMessages
)

type Store struct{ inner *iroute.Store }

func NewStore(db *sql.DB) *Store { return &Store{inner: iroute.NewStore(db)} }

func (s *Store) Add(ctx context.Context, public, providerID, upstream string, protocol Protocol, modality Modality) (Route, error) {
	return s.inner.Add(ctx, public, providerID, upstream, protocol, modality)
}

func (s *Store) List(ctx context.Context) ([]Route, error) { return s.inner.List(ctx) }

func (s *Store) Get(ctx context.Context, id string) (Route, error) {
	return s.inner.Get(ctx, id)
}

func (s *Store) Update(ctx context.Context, id, upstream string, protocol Protocol, modality Modality, inputUSD, cachedUSD, outputUSD *float64) (Route, error) {
	return s.inner.Update(ctx, id, upstream, protocol, modality, inputUSD, cachedUSD, outputUSD)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	return s.inner.SetEnabled(ctx, id, enabled)
}
