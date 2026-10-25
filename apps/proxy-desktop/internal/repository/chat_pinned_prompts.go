package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

var ErrChatPinnedPromptNotFound = errors.New("pinned prompt not found")

const maxChatPinnedBodyLen = 32 * 1024

type ChatPinnedPrompts struct {
	db *sql.DB
}

func NewChatPinnedPrompts(db *sql.DB) *ChatPinnedPrompts {
	return &ChatPinnedPrompts{db: db}
}

func (r *ChatPinnedPrompts) List(ctx context.Context) ([]domain.ChatPinnedPrompt, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, body, created_at, updated_at
		FROM chat_pinned_prompts
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ChatPinnedPrompt
	for rows.Next() {
		p, err := scanChatPinnedPrompt(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *ChatPinnedPrompts) Create(ctx context.Context, title, body string) (domain.ChatPinnedPrompt, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return domain.ChatPinnedPrompt{}, errors.New("prompt body is required")
	}
	if len(body) > maxChatPinnedBodyLen {
		return domain.ChatPinnedPrompt{}, errors.New("prompt body is too long")
	}
	title = strings.TrimSpace(title)
	now := time.Now().UTC()
	p := domain.ChatPinnedPrompt{
		ID:        uuid.NewString(),
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chat_pinned_prompts (id, title, body, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		p.ID, p.Title, p.Body,
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.ChatPinnedPrompt{}, err
	}
	return p, nil
}

func (r *ChatPinnedPrompts) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM chat_pinned_prompts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrChatPinnedPromptNotFound
	}
	return nil
}

func scanChatPinnedPrompt(rows *sql.Rows) (domain.ChatPinnedPrompt, error) {
	var p domain.ChatPinnedPrompt
	var created, updated string
	if err := rows.Scan(&p.ID, &p.Title, &p.Body, &created, &updated); err != nil {
		return domain.ChatPinnedPrompt{}, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, nil
}
