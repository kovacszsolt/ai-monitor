package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/settings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
)

// OpenWithMigrations opens the DB, runs embedded migrations, and seeds proxy settings from bootstrap YAML.
func OpenWithMigrations(cfg config.Config) (*sql.DB, config.Config, error) {
	db, err := Open(cfg.Storage)
	if err != nil {
		return nil, config.Config{}, err
	}
	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, config.Config{}, err
	}
	ctx := context.Background()
	bootstrap, _, err := config.ProxyYAMLBootstrap(cfg)
	if err != nil {
		db.Close()
		return nil, config.Config{}, err
	}
	if err := settings.SeedFromConfig(ctx, db, bootstrap); err != nil {
		db.Close()
		return nil, config.Config{}, fmt.Errorf("seed proxy settings: %w", err)
	}
	if _, err := settings.LoadProxyConfig(ctx, db, bootstrap); err != nil {
		db.Close()
		return nil, config.Config{}, fmt.Errorf("load proxy settings: %w", err)
	}
	return db, cfg, nil
}
