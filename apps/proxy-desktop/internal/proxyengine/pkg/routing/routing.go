package routing

import (
	"context"
	"database/sql"

	iroute "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

type Route = iroute.Route
type Protocol = iroute.Protocol

var ErrNotFound = iroute.ErrNotFound

type Store struct{ inner *iroute.Store }

func NewStore(db *sql.DB) *Store { return &Store{inner: iroute.NewStore(db)} }

func (s *Store) Add(ctx context.Context, public, providerID, upstream string, protocol Protocol) (Route, error) {
	return s.inner.Add(ctx, public, providerID, upstream, protocol)
}

func (s *Store) List(ctx context.Context) ([]Route, error) { return s.inner.List(ctx) }

func (s *Store) Get(ctx context.Context, id string) (Route, error) {
	return s.inner.Get(ctx, id)
}

func (s *Store) Update(ctx context.Context, id, upstream string, protocol Protocol, inputUSD, cachedUSD, outputUSD *float64) (Route, error) {
	return s.inner.Update(ctx, id, upstream, protocol, inputUSD, cachedUSD, outputUSD)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	return s.inner.SetEnabled(ctx, id, enabled)
}
