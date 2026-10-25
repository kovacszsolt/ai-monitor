package storage

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "modernc.org/sqlite"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/migrations"
)

func Open(cfg config.StorageConfig) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
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
	if _, err := db.Exec(fmt.Sprintf("PRAGMA journal_mode=%s", cfg.JournalMode)); err != nil {
		db.Close()
		return nil, err
	}
	if cfg.Synchronous != "" {
		if _, err := db.Exec(fmt.Sprintf("PRAGMA synchronous=%s", cfg.Synchronous)); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}

func Migrate(db *sql.DB, migrationsDir string) error {
	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	sourceURL := fmt.Sprintf("file://%s", abs)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite3", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// MigrateEmbedded runs migrations from the embedded SQL bundle (see migrations/embed.go).
func MigrateEmbedded(db *sql.DB) error {
	return MigrateFS(db, migrations.FS, ".")
}

// MigrateFS runs migrations from an fs.FS source (e.g. embed or test fixture).
func MigrateFS(db *sql.DB, migrationFS fs.FS, root string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	source, err := iofs.New(migrationFS, root)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// RunMigrations applies migrations from the embedded SQL bundle shipped with the module.
// migrationsDir is ignored (kept for callers that pass DefaultMigrationsDir()); running a
// second pass from disk caused version skew when the DB was already at a newer embedded version.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	return MigrateEmbedded(db)
}

func Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}
