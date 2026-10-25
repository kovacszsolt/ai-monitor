package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const KeyPrefix = "aip_"

var (
	ErrKeyNotFound     = errors.New("api key not found")
	ErrKeyDisabled     = errors.New("api key disabled")
	ErrKeyExpired      = errors.New("api key expired")
	ErrInvalidKey      = errors.New("invalid api key")
	ErrProjectDisabled = errors.New("project disabled")
)

type APIKey struct {
	ID             string
	ProjectID      string
	Name           string
	Secret         string // plaintext aip_… stored in proxy_api_keys.api_key
	AllowAllModels bool
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      *time.Time
	LastUsedAt     *time.Time
}

type AuthenticatedKey struct {
	APIKey
	ProjectEnabled bool
}

type KeyStore struct {
	db *sql.DB
}

func NewKeyStore(db *sql.DB) *KeyStore {
	return &KeyStore{db: db}
}

func (s *KeyStore) Create(ctx context.Context, projectID, name string, allowAll bool, expiresAt *time.Time) (APIKey, string, error) {
	secret, err := generateSecret()
	if err != nil {
		return APIKey{}, "", err
	}
	plaintext := KeyPrefix + secret
	now := time.Now().UTC()
	k := APIKey{
		ID:             uuid.NewString(),
		ProjectID:      projectID,
		Name:           name,
		Secret:         plaintext,
		AllowAllModels: allowAll,
		Enabled:        true,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      expiresAt,
	}
	var exp *string
	if expiresAt != nil {
		v := expiresAt.UTC().Format(time.RFC3339Nano)
		exp = &v
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO proxy_api_keys (id, project_id, name, api_key, allow_all_models, enabled, created_at, updated_at, expires_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?, ?)`,
		k.ID, k.ProjectID, k.Name, plaintext, boolInt(allowAll),
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), exp)
	if err != nil {
		return APIKey{}, "", err
	}
	return k, plaintext, nil
}

func (s *KeyStore) Authenticate(ctx context.Context, bearer string) (AuthenticatedKey, error) {
	token := strings.TrimSpace(bearer)
	if token == "" || !strings.HasPrefix(token, KeyPrefix) {
		return AuthenticatedKey{}, ErrInvalidKey
	}
	row := s.db.QueryRowContext(ctx, `
		SELECT k.id, k.project_id, k.name, k.api_key, k.allow_all_models, k.enabled, k.created_at, k.updated_at, k.expires_at, k.last_used_at, p.enabled
		FROM proxy_api_keys k
		JOIN projects p ON p.id = k.project_id
		WHERE k.api_key = ?`, token)

	var k APIKey
	var allowAll, enabled, projectEnabled int
	var created, updated string
	var expires, lastUsed sql.NullString
	if err := row.Scan(&k.ID, &k.ProjectID, &k.Name, &k.Secret, &allowAll, &enabled, &created, &updated, &expires, &lastUsed, &projectEnabled); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthenticatedKey{}, ErrKeyNotFound
		}
		return AuthenticatedKey{}, err
	}
	k.AllowAllModels = allowAll == 1
	k.Enabled = enabled == 1
	k.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	k.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	if expires.Valid {
		t, _ := time.Parse(time.RFC3339Nano, expires.String)
		k.ExpiresAt = &t
	}
	if lastUsed.Valid {
		t, _ := time.Parse(time.RFC3339Nano, lastUsed.String)
		k.LastUsedAt = &t
	}
	if !k.Enabled {
		return AuthenticatedKey{}, ErrKeyDisabled
	}
	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return AuthenticatedKey{}, ErrKeyExpired
	}
	if projectEnabled != 1 {
		return AuthenticatedKey{}, ErrProjectDisabled
	}
	if subtle.ConstantTimeCompare([]byte(k.Secret), []byte(token)) != 1 {
		return AuthenticatedKey{}, ErrKeyNotFound
	}
	now := time.Now().UTC()
	_, _ = s.db.ExecContext(ctx, `UPDATE proxy_api_keys SET last_used_at = ? WHERE id = ?`,
		now.Format(time.RFC3339Nano), k.ID)
	return AuthenticatedKey{APIKey: k, ProjectEnabled: true}, nil
}

func (s *KeyStore) ListByProject(ctx context.Context, projectID string) ([]APIKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, project_id, name, api_key, allow_all_models, enabled, created_at, updated_at, expires_at, last_used_at
		FROM proxy_api_keys WHERE project_id = ? ORDER BY name`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanKeys(rows)
}

func (s *KeyStore) Get(ctx context.Context, id string) (APIKey, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, project_id, name, api_key, allow_all_models, enabled, created_at, updated_at, expires_at, last_used_at
		FROM proxy_api_keys WHERE id = ?`, id)
	return scanKeyRow(row)
}

func (s *KeyStore) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM proxy_api_keys WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrKeyNotFound
	}
	return nil
}

func (s *KeyStore) Revoke(ctx context.Context, id string) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE proxy_api_keys SET enabled = 0, updated_at = ? WHERE id = ?`,
		now.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrKeyNotFound
	}
	return nil
}

func (s *KeyStore) AllowedModelRouteIDs(ctx context.Context, keyID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT model_route_id FROM api_key_allowed_models WHERE api_key_id = ?`, keyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *KeyStore) SetAllowedModels(ctx context.Context, keyID string, routeIDs []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM api_key_allowed_models WHERE api_key_id = ?`, keyID); err != nil {
		return err
	}
	for _, rid := range routeIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO api_key_allowed_models (api_key_id, model_route_id) VALUES (?, ?)`, keyID, rid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanKeys(rows *sql.Rows) ([]APIKey, error) {
	var out []APIKey
	for rows.Next() {
		k, err := scanKeyFromRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

func scanKeyRow(row *sql.Row) (APIKey, error) {
	var k APIKey
	var allowAll, enabled int
	var created, updated string
	var expires, lastUsed sql.NullString
	if err := row.Scan(&k.ID, &k.ProjectID, &k.Name, &k.Secret, &allowAll, &enabled, &created, &updated, &expires, &lastUsed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return APIKey{}, ErrKeyNotFound
		}
		return APIKey{}, err
	}
	k.AllowAllModels = allowAll == 1
	k.Enabled = enabled == 1
	k.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	k.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	if expires.Valid {
		t, _ := time.Parse(time.RFC3339Nano, expires.String)
		k.ExpiresAt = &t
	}
	if lastUsed.Valid {
		t, _ := time.Parse(time.RFC3339Nano, lastUsed.String)
		k.LastUsedAt = &t
	}
	return k, nil
}

func scanKeyFromRows(rows *sql.Rows) (APIKey, error) {
	var k APIKey
	var allowAll, enabled int
	var created, updated string
	var expires, lastUsed sql.NullString
	if err := rows.Scan(&k.ID, &k.ProjectID, &k.Name, &k.Secret, &allowAll, &enabled, &created, &updated, &expires, &lastUsed); err != nil {
		return APIKey{}, err
	}
	k.AllowAllModels = allowAll == 1
	k.Enabled = enabled == 1
	k.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	k.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	if expires.Valid {
		t, _ := time.Parse(time.RFC3339Nano, expires.String)
		k.ExpiresAt = &t
	}
	if lastUsed.Valid {
		t, _ := time.Parse(time.RFC3339Nano, lastUsed.String)
		k.LastUsedAt = &t
	}
	return k, nil
}
