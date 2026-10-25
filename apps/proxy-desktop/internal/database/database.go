package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
)

// Open opens an existing proxy SQLite database. No migrations are run.
func Open(cfg config.StorageConfig) (*sql.DB, error) {
	if cfg.DatabasePath == "" {
		return nil, fmt.Errorf("database path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		return nil, fmt.Errorf("ensure db directory: %w", err)
	}
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(%d)",
		cfg.DatabasePath, cfg.BusyTimeoutMS)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	if cfg.JournalMode != "" {
		if _, err := db.Exec(fmt.Sprintf("PRAGMA journal_mode=%s", cfg.JournalMode)); err != nil {
			db.Close()
			return nil, err
		}
	}
	if cfg.Synchronous != "" {
		if _, err := db.Exec(fmt.Sprintf("PRAGMA synchronous=%s", cfg.Synchronous)); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}

func Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

// SchemaReady checks that core proxy tables exist (embedded migrations).
func SchemaReady(ctx context.Context, db *sql.DB) error {
	const q = `
		SELECT COUNT(*) FROM sqlite_master
		WHERE type = 'table' AND name IN ('projects','providers','model_routes','proxy_api_keys','ai_requests')`
	var n int
	if err := db.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return fmt.Errorf("schema check: %w", err)
	}
	if n < 5 {
		return fmt.Errorf("database schema incomplete: run migrations on this file first (Settings → Database maintenance)")
	}
	return nil
}
