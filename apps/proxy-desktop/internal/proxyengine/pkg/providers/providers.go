package providers

import (
	"context"
	"database/sql"

	iprov "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

type Provider = iprov.Provider
type Type = iprov.Type

var ErrNotFound = iprov.ErrNotFound

const (
	TypeOpenAI           = iprov.TypeOpenAI
	TypeAnthropic        = iprov.TypeAnthropic
	TypeOpenAICompatible = iprov.TypeOpenAICompatible
)

type Store struct{ inner *iprov.Store }

func NewStore(db *sql.DB) *Store { return &Store{inner: iprov.NewStore(db)} }

func (s *Store) Add(ctx context.Context, name string, ptype Type, baseURL, credential string, connMS, reqMS int) (Provider, error) {
	return s.inner.Add(ctx, name, ptype, baseURL, credential, connMS, reqMS)
}

func (s *Store) List(ctx context.Context) ([]Provider, error) { return s.inner.List(ctx) }

func (s *Store) Get(ctx context.Context, id string) (Provider, error) { return s.inner.Get(ctx, id) }

func (s *Store) Update(ctx context.Context, id string, baseURL, credential string, connMS, reqMS int) (Provider, error) {
	return s.inner.Update(ctx, id, baseURL, credential, connMS, reqMS)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	return s.inner.SetEnabled(ctx, id, enabled)
}
