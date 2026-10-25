package providers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("provider not found")

type Type string

const (
	TypeOpenAI           Type = "openai"
	TypeAnthropic        Type = "anthropic"
	TypeOpenAICompatible Type = "openai-compatible"
)

type Provider struct {
	ID                  string
	Name                string
	ProviderType        Type
	BaseURL             string
	Credential          string
	Enabled             bool
	ConnectionTimeoutMS int
	RequestTimeoutMS    int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Add(ctx context.Context, name string, ptype Type, baseURL, credential string, connMS, reqMS int) (Provider, error) {
	now := time.Now().UTC()
	p := Provider{
		ID:                  uuid.NewString(),
		Name:                name,
		ProviderType:        ptype,
		BaseURL:             strings.TrimRight(baseURL, "/"),
		Credential:          strings.TrimSpace(credential),
		Enabled:             true,
		ConnectionTimeoutMS: connMS,
		RequestTimeoutMS:    reqMS,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if p.ConnectionTimeoutMS <= 0 {
		p.ConnectionTimeoutMS = 10000
	}
	if p.RequestTimeoutMS <= 0 {
		p.RequestTimeoutMS = 300000
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO providers (id, name, provider_type, base_url, credential, enabled, connection_timeout_ms, request_timeout_ms, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?, ?, ?)`,
		p.ID, p.Name, string(p.ProviderType), p.BaseURL, p.Credential,
		p.ConnectionTimeoutMS, p.RequestTimeoutMS,
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return Provider{}, err
	}
	return p, nil
}

func (s *Store) List(ctx context.Context) ([]Provider, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, provider_type, base_url, credential, enabled, connection_timeout_ms, request_timeout_ms, created_at, updated_at
		FROM providers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Provider
	for rows.Next() {
		p, err := scanProvider(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) Get(ctx context.Context, id string) (Provider, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, provider_type, base_url, credential, enabled, connection_timeout_ms, request_timeout_ms, created_at, updated_at
		FROM providers WHERE id = ?`, id)
	return scanProviderRow(row)
}

func (s *Store) GetByName(ctx context.Context, name string) (Provider, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, provider_type, base_url, credential, enabled, connection_timeout_ms, request_timeout_ms, created_at, updated_at
		FROM providers WHERE name = ?`, name)
	return scanProviderRow(row)
}

func (s *Store) Update(ctx context.Context, id string, baseURL, credential string, connMS, reqMS int) (Provider, error) {
	now := time.Now().UTC()
	baseURL = strings.TrimRight(baseURL, "/")
	credential = strings.TrimSpace(credential)
	var res sql.Result
	var err error
	if credential == "" {
		res, err = s.db.ExecContext(ctx, `
			UPDATE providers SET base_url = ?, connection_timeout_ms = ?, request_timeout_ms = ?, updated_at = ?
			WHERE id = ?`,
			baseURL, connMS, reqMS, now.Format(time.RFC3339Nano), id)
	} else {
		res, err = s.db.ExecContext(ctx, `
			UPDATE providers SET base_url = ?, credential = ?, connection_timeout_ms = ?, request_timeout_ms = ?, updated_at = ?
			WHERE id = ?`,
			baseURL, credential, connMS, reqMS, now.Format(time.RFC3339Nano), id)
	}
	if err != nil {
		return Provider{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Provider{}, ErrNotFound
	}
	return s.Get(ctx, id)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE providers SET enabled = ?, updated_at = ? WHERE id = ?`,
		boolInt(enabled), now.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (p Provider) APIKey() (string, error) {
	return strings.TrimSpace(p.Credential), nil
}

func (p Provider) RequireAPIKey() error {
	if strings.TrimSpace(p.Credential) == "" {
		return fmt.Errorf("provider %q has no API credential configured", p.Name)
	}
	return nil
}

func scanProvider(rows *sql.Rows) (Provider, error) {
	var p Provider
	var ptype string
	var enabled int
	var created, updated string
	if err := rows.Scan(&p.ID, &p.Name, &ptype, &p.BaseURL, &p.Credential, &enabled,
		&p.ConnectionTimeoutMS, &p.RequestTimeoutMS, &created, &updated); err != nil {
		return Provider{}, err
	}
	p.ProviderType = Type(ptype)
	p.Enabled = enabled == 1
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, nil
}

func scanProviderRow(row *sql.Row) (Provider, error) {
	var p Provider
	var ptype string
	var enabled int
	var created, updated string
	if err := row.Scan(&p.ID, &p.Name, &ptype, &p.BaseURL, &p.Credential, &enabled,
		&p.ConnectionTimeoutMS, &p.RequestTimeoutMS, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Provider{}, ErrNotFound
		}
		return Provider{}, err
	}
	p.ProviderType = Type(ptype)
	p.Enabled = enabled == 1
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, nil
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
