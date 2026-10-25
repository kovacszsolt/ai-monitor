package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/auth"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/capture"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/logging"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/projects"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/proxy"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/server"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/settings"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/storage"
)

type App struct {
	Config    config.Config
	DB        *sql.DB
	Writer    *capture.Writer
	Projects  *projects.Store
	Keys      *auth.KeyStore
	Providers *providers.Store
	Routes    *routing.Store
	Server    *server.Server
}

func Open(cfg config.Config, migrationsDir string) (*App, error) {
	db, err := storage.Open(cfg.Storage)
	if err != nil {
		return nil, err
	}
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	if err := storage.RunMigrations(db, migrationsDir); err != nil {
		db.Close()
		return nil, err
	}
	ctx := context.Background()
	if err := settings.SeedFromConfig(ctx, db, cfg); err != nil {
		db.Close()
		return nil, fmt.Errorf("seed proxy settings: %w", err)
	}
	cfg, err = settings.LoadProxyConfig(ctx, db, cfg)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("load proxy settings: %w", err)
	}
	writer := capture.NewWriter(db, cfg.Persistence)
	writer.Start()
	a := &App{
		Config:    cfg,
		DB:        db,
		Writer:    writer,
		Projects:  projects.NewStore(db),
		Keys:      auth.NewKeyStore(db),
		Providers: providers.NewStore(db),
		Routes:    routing.NewStore(db),
	}
	h := proxy.NewHandler(proxy.Services{
		Keys:      a.Keys,
		Providers: a.Providers,
		Routes:    a.Routes,
		Writer:    writer,
		MaxBody:   cfg.Server.MaxRequestBytes,
		Capture:   cfg.Capture,
	})
	a.Server = server.New(cfg.ListenAddr(), h)
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.Server.Start()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = a.Server.Shutdown(shutdownCtx)
		_ = a.Writer.Shutdown(shutdownCtx)
		a.DB.Close()
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func DefaultMigrationsDir() string {
	if v := os.Getenv("AI_PROXY_MIGRATIONS"); v != "" {
		return v
	}
	// When running from apps/proxy/proxy-desktop project root (cwd).
	if _, err := os.Stat("migrations"); err == nil {
		return "migrations"
	}
	exe, err := os.Executable()
	if err != nil {
		return "migrations"
	}
	exeDir := filepath.Dir(exe)
		// Prefer migrations next to the binary or one level up (legacy layouts).
	for _, candidate := range []string{
		filepath.Join(exeDir, "migrations"),
		filepath.Join(exeDir, "..", "migrations"),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return filepath.Join(exeDir, "migrations")
}

func SetupLogger(cfg config.LoggingConfig) {
	logging.Setup(cfg)
}
