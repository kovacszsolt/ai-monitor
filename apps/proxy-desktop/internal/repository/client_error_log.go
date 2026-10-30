package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

const (
	maxClientErrorMessageLen = 8 * 1024
	maxClientErrorDetailLen  = 64 * 1024
	clientErrorLogMaxRows    = 500
)

type ClientErrorLog struct {
	db *sql.DB
}

func NewClientErrorLog(db *sql.DB) *ClientErrorLog {
	return &ClientErrorLog{db: db}
}

func (r *ClientErrorLog) Insert(ctx context.Context, level, source, message, detail string) (domain.ClientErrorLogEntry, error) {
	level = normalizeClientErrorLevel(level)
	source = trimClientErrorField(source, 256)
	message = trimClientErrorField(message, maxClientErrorMessageLen)
	detail = trimClientErrorField(detail, maxClientErrorDetailLen)
	if message == "" {
		message = "(empty message)"
	}
	now := time.Now().UTC()
	entry := domain.ClientErrorLogEntry{
		ID:        uuid.NewString(),
		Level:     level,
		Source:    source,
		Message:   message,
		Detail:    detail,
		CreatedAt: now,
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO client_error_log (id, level, source, message, detail, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.Level, entry.Source, entry.Message, entry.Detail,
		now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.ClientErrorLogEntry{}, err
	}
	if err := r.trimToMaxRows(ctx); err != nil {
		return entry, err
	}
	return entry, nil
}

type ClientErrorLogListResult struct {
	Items    []domain.ClientErrorLogEntry
	Total    int64
	Page     int
	PageSize int
}

func (r *ClientErrorLog) List(ctx context.Context, page, pageSize int) (ClientErrorLogListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM client_error_log`).Scan(&total); err != nil {
		return ClientErrorLogListResult{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, level, source, message, detail, created_at
		FROM client_error_log
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return ClientErrorLogListResult{}, err
	}
	defer rows.Close()

	items := make([]domain.ClientErrorLogEntry, 0, pageSize)
	for rows.Next() {
		e, err := scanClientErrorLogEntry(rows)
		if err != nil {
			return ClientErrorLogListResult{}, err
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return ClientErrorLogListResult{}, err
	}
	return ClientErrorLogListResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (r *ClientErrorLog) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM client_error_log`)
	return err
}

func (r *ClientErrorLog) trimToMaxRows(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM client_error_log
		WHERE id NOT IN (
			SELECT id FROM client_error_log
			ORDER BY created_at DESC
			LIMIT ?
		)`, clientErrorLogMaxRows)
	return err
}

func scanClientErrorLogEntry(scanner interface {
	Scan(dest ...any) error
}) (domain.ClientErrorLogEntry, error) {
	var (
		e         domain.ClientErrorLogEntry
		createdAt string
	)
	if err := scanner.Scan(&e.ID, &e.Level, &e.Source, &e.Message, &e.Detail, &createdAt); err != nil {
		return domain.ClientErrorLogEntry{}, err
	}
	t, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		t, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return domain.ClientErrorLogEntry{}, err
		}
	}
	e.CreatedAt = t.UTC()
	return e, nil
}

func normalizeClientErrorLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "warn", "warning":
		return "warn"
	case "info":
		return "info"
	default:
		return "error"
	}
}

func trimClientErrorField(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "…"
}
