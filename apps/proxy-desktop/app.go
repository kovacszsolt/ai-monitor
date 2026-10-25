package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/database"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyruntime"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service"
)

type App struct {
	ctx         context.Context
	service     *service.DesktopService
	proxy       *service.ProxyRuntimeService
	runtime     *proxyruntime.Runtime
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	db, cfg, err := database.OpenWithMigrations(cfg)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	bridge, err := config.LoadBridge(cfg, db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("config bridge: %w", err)
	}
	rt := proxyruntime.New()
	svc := service.NewDesktopService(cfg, bridge, db, rt)
	proxySvc := service.NewProxyRuntimeService(context.Background(), bridge, rt, db)
	svc.SetBridgeReloadHandler(proxySvc.UpdateBridge)
	return &App{
		service: svc,
		proxy:   proxySvc,
		runtime: rt,
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.proxy != nil {
		a.proxy.SetContext(ctx)
	}
	if a.service != nil {
		a.service.SetWailsContext(ctx)
		if res := a.service.RunDatabaseMigrations(); !res.OK {
			slog.Warn("database migration on startup", "error", res.Error.Message)
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.runtime != nil {
		a.runtime.Stop()
	}
	if a.service != nil {
		a.service.CloseDB()
	}
}
