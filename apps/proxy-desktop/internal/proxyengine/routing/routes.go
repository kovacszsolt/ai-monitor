package routing

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("model route not found")
var ErrDisabled = errors.New("model route disabled")

type Protocol string

const (
	ProtocolOpenAIChat        Protocol = "openai_chat"
	ProtocolOpenAIResponses   Protocol = "openai_responses"
	ProtocolOpenAIImages      Protocol = "openai_images"
	ProtocolAnthropicMessages Protocol = "anthropic_messages"
)

type Modality string

const (
	ModalityTextText      Modality = "text-text"
	ModalityTextImage     Modality = "text-image"
	ModalityTextImageText Modality = "text_image-text"
)

type Route struct {
	ID                       string
	PublicModelName          string
	ProviderID               string
	UpstreamModelName        string
	Protocol                 Protocol
	Modality                 Modality
	Enabled                  bool
	InputTokenUSDPer1M       *float64
	CachedInputTokenUSDPer1M *float64
	OutputTokenUSDPer1M      *float64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Add(ctx context.Context, public, providerID, upstream string, protocol Protocol, modality Modality) (Route, error) {
	now := time.Now().UTC()
	if modality == "" {
		modality = ModalityTextText
	}
	r := Route{
		ID:                uuid.NewString(),
		PublicModelName:   public,
		ProviderID:        providerID,
		UpstreamModelName: upstream,
		Protocol:          protocol,
		Modality:          modality,
		Enabled:           true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO model_routes (
			id, public_model_name, provider_id, upstream_model_name, protocol, modality, enabled,
			input_token_usd_per_1m, cached_input_token_usd_per_1m, output_token_usd_per_1m,
			created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?, ?)`,
		r.ID, r.PublicModelName, r.ProviderID, r.UpstreamModelName, string(r.Protocol), string(r.Modality),
		nullFloatArg(r.InputTokenUSDPer1M), nullFloatArg(r.CachedInputTokenUSDPer1M), nullFloatArg(r.OutputTokenUSDPer1M),
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return Route{}, err
	}
	return r, nil
}

func (s *Store) GetByPublicModel(ctx context.Context, public string) (Route, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, public_model_name, provider_id, upstream_model_name, protocol, modality, enabled,
			input_token_usd_per_1m, cached_input_token_usd_per_1m, output_token_usd_per_1m,
			created_at, updated_at
		FROM model_routes WHERE public_model_name = ?`, public)
	return scanRouteRow(row)
}

func (s *Store) Get(ctx context.Context, id string) (Route, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, public_model_name, provider_id, upstream_model_name, protocol, modality, enabled,
			input_token_usd_per_1m, cached_input_token_usd_per_1m, output_token_usd_per_1m,
			created_at, updated_at
		FROM model_routes WHERE id = ?`, id)
	return scanRouteRow(row)
}

func (s *Store) List(ctx context.Context) ([]Route, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, public_model_name, provider_id, upstream_model_name, protocol, modality, enabled,
			input_token_usd_per_1m, cached_input_token_usd_per_1m, output_token_usd_per_1m,
			created_at, updated_at
		FROM model_routes ORDER BY public_model_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Route
	for rows.Next() {
		r, err := scanRoute(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) Update(ctx context.Context, id, upstream string, protocol Protocol, modality Modality, inputUSDPer1M, cachedUSDPer1M, outputUSDPer1M *float64) (Route, error) {
	now := time.Now().UTC()
	if modality == "" {
		modality = ModalityTextText
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE model_routes SET
			upstream_model_name = ?, protocol = ?, modality = ?,
			input_token_usd_per_1m = ?, cached_input_token_usd_per_1m = ?, output_token_usd_per_1m = ?,
			updated_at = ?
		WHERE id = ?`,
		upstream, string(protocol), string(modality),
		nullFloatArg(inputUSDPer1M), nullFloatArg(cachedUSDPer1M), nullFloatArg(outputUSDPer1M),
		now.Format(time.RFC3339Nano), id)
	if err != nil {
		return Route{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return Route{}, ErrNotFound
	}
	return s.Get(ctx, id)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE model_routes SET enabled = ?, updated_at = ? WHERE id = ?`,
		boolInt(enabled), now.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanRoute(rows *sql.Rows) (Route, error) {
	var r Route
	var protocol, modality string
	var enabled int
	var inputUSD, cachedUSD, outputUSD sql.NullFloat64
	var created, updated string
	if err := rows.Scan(
		&r.ID, &r.PublicModelName, &r.ProviderID, &r.UpstreamModelName, &protocol, &modality, &enabled,
		&inputUSD, &cachedUSD, &outputUSD, &created, &updated,
	); err != nil {
		return Route{}, err
	}
	r.Protocol = Protocol(protocol)
	r.Modality = normalizeModality(modality)
	r.Enabled = enabled == 1
	r.InputTokenUSDPer1M = nullFloatPtr(inputUSD)
	r.CachedInputTokenUSDPer1M = nullFloatPtr(cachedUSD)
	r.OutputTokenUSDPer1M = nullFloatPtr(outputUSD)
	r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	r.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return r, nil
}

func scanRouteRow(row *sql.Row) (Route, error) {
	var r Route
	var protocol, modality string
	var enabled int
	var inputUSD, cachedUSD, outputUSD sql.NullFloat64
	var created, updated string
	if err := row.Scan(
		&r.ID, &r.PublicModelName, &r.ProviderID, &r.UpstreamModelName, &protocol, &modality, &enabled,
		&inputUSD, &cachedUSD, &outputUSD, &created, &updated,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Route{}, ErrNotFound
		}
		return Route{}, err
	}
	r.Protocol = Protocol(protocol)
	r.Modality = normalizeModality(modality)
	r.Enabled = enabled == 1
	r.InputTokenUSDPer1M = nullFloatPtr(inputUSD)
	r.CachedInputTokenUSDPer1M = nullFloatPtr(cachedUSD)
	r.OutputTokenUSDPer1M = nullFloatPtr(outputUSD)
	r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	r.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return r, nil
}

func normalizeModality(v string) Modality {
	switch Modality(v) {
	case ModalityTextText, ModalityTextImage, ModalityTextImageText:
		return Modality(v)
	default:
		return ModalityTextText
	}
}

func nullFloatArg(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullFloatPtr(n sql.NullFloat64) *float64 {
	if !n.Valid {
		return nil
	}
	v := n.Float64
	return &v
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
