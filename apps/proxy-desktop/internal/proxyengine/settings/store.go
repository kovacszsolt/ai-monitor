package settings

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

const KeyKeyPepper = "key_pepper"

var ErrNotFound = errors.New("setting not found")
var ErrEmptyValue = errors.New("setting value is empty")

func Get(ctx context.Context, db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrEmptyValue
	}
	return value, nil
}

func Set(ctx context.Context, db *sql.DB, key, value string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := db.ExecContext(ctx, `
		INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value, now)
	return err
}

// migrateLegacyKeyPepper copies proxy_key_pepper (old schema) into settings then drops the legacy table.
func migrateLegacyKeyPepper(ctx context.Context, db *sql.DB) error {
	var legacy string
	err := db.QueryRowContext(ctx, `
		SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'proxy_key_pepper'`).Scan(&legacy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	var pepper, updated string
	err = db.QueryRowContext(ctx, `SELECT pepper, updated_at FROM proxy_key_pepper WHERE id = 1`).Scan(&pepper, &updated)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if strings.TrimSpace(pepper) != "" {
		if err := Set(ctx, db, KeyKeyPepper, strings.TrimSpace(pepper)); err != nil {
			return err
		}
	}
	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS proxy_key_pepper`)
	return err
}
