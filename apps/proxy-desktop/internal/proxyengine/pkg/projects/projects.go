package projects

import (
	"context"
	"database/sql"

	iproj "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/projects"
)

type Project = iproj.Project

var (
	ErrNotFound  = iproj.ErrNotFound
	ErrDuplicate = iproj.ErrDuplicate
)

type Store struct{ inner *iproj.Store }

func NewStore(db *sql.DB) *Store { return &Store{inner: iproj.NewStore(db)} }

func (s *Store) Create(ctx context.Context, name, description string) (Project, error) {
	return s.inner.Create(ctx, name, description)
}

func (s *Store) List(ctx context.Context) ([]Project, error) { return s.inner.List(ctx) }

func (s *Store) GetByID(ctx context.Context, id string) (Project, error) { return s.inner.GetByID(ctx, id) }

func (s *Store) Update(ctx context.Context, id, name, description string) (Project, error) {
	return s.inner.Update(ctx, id, name, description)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	return s.inner.SetEnabled(ctx, id, enabled)
}
