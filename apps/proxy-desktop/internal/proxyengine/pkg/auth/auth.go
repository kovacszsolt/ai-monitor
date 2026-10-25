package auth

import (
	"context"
	"database/sql"
	"time"

	iauth "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/auth"
)

type APIKey = iauth.APIKey

var ErrKeyNotFound = iauth.ErrKeyNotFound

type KeyStore struct{ inner *iauth.KeyStore }

func NewKeyStore(db *sql.DB) *KeyStore {
	return &KeyStore{inner: iauth.NewKeyStore(db)}
}

func (s *KeyStore) Create(ctx context.Context, projectID, name string, allowAll bool, expiresAt *time.Time) (APIKey, string, error) {
	return s.inner.Create(ctx, projectID, name, allowAll, expiresAt)
}

func (s *KeyStore) ListByProject(ctx context.Context, projectID string) ([]APIKey, error) {
	return s.inner.ListByProject(ctx, projectID)
}

func (s *KeyStore) Get(ctx context.Context, id string) (APIKey, error) {
	return s.inner.Get(ctx, id)
}

func (s *KeyStore) Revoke(ctx context.Context, id string) error { return s.inner.Revoke(ctx, id) }

func (s *KeyStore) Delete(ctx context.Context, id string) error { return s.inner.Delete(ctx, id) }
