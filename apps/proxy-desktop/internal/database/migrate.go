package database

import (
	"database/sql"
	"fmt"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/storage"
)

// RunMigrations applies embedded proxy SQLite schema migrations (migrations/).
func RunMigrations(db *sql.DB) error {
	if err := storage.MigrateEmbedded(db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
