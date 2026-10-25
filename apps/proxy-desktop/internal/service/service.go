package service

import (
	"context"
	"database/sql"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/buildinfo"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/database"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyruntime"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/airequests"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/analytics"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/catalog"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/deps"
)

type DesktopService struct {
	deps       *deps.Desktop
	catalog    *catalog.Service
	analytics  *analytics.Service
	airequests *airequests.Service

	onBridgeReload func(config.Bridge)
	wailsCtx       context.Context
}

func (s *DesktopService) SetBridgeReloadHandler(fn func(config.Bridge)) {
	s.onBridgeReload = fn
}

func (s *DesktopService) reconnect(cfg config.Config, bridge config.Bridge, db *sql.DB) error {
	if s.deps.DB != nil {
		s.deps.DB.Close()
	}
	s.deps.Cfg = cfg
	s.deps.Bridge = bridge
	s.deps.RebindDB(db)
	return nil
}

func NewDesktopService(cfg config.Config, bridge config.Bridge, db *sql.DB, runtime *proxyruntime.Runtime) *DesktopService {
	s := &DesktopService{
		onBridgeReload: nil,
	}
	s.deps = deps.NewDesktop(cfg, bridge, db, runtime, s.ctx)
	s.catalog = catalog.New(s.deps)
	s.analytics = analytics.New(s.deps)
	s.airequests = airequests.New(s.deps)
	return s
}

func (s *DesktopService) ctx() context.Context {
	if s.wailsCtx != nil {
		return s.wailsCtx
	}
	return context.Background()
}

func (s *DesktopService) GetAbout() AboutDTO {
	conn := s.GetConnectionInfo()
	about := AboutDTO{
		AppName:          buildinfo.AppName,
		Version:          buildinfo.Version,
		BuildTime:        conn.BuildTime,
		Author:           "Kovacs Zsolt",
		License:          buildinfo.License,
		DatabasePath:     conn.DatabasePath,
		ProxyConfigPath:  conn.ProxyConfigPath,
		SchemaReady:      conn.SchemaReady,
		MigrationVersion: conn.MigrationVersion,
		MigrationDirty:   conn.MigrationDirty,
		ConfigBridgeOK:   conn.ConfigBridgeOK,
		BridgeMessage:    conn.BridgeMessage,
	}
	settings := s.GetSettings()
	if settings.OK {
		about.ConfigFilePath = settings.Data.ConfigFilePath
		if about.DatabasePath == "" {
			about.DatabasePath = settings.Data.EffectiveDatabasePath
		}
		if about.ProxyConfigPath == "" {
			about.ProxyConfigPath = settings.Data.ProxyConfigPath
		}
	}
	return about
}

func (s *DesktopService) GetConnectionInfo() ConnectionInfoDTO {
	info := ConnectionInfoDTO{
		DatabasePath:    s.deps.Cfg.Storage.DatabasePath,
		Connected:       false,
		SchemaReady:     false,
		BuildTime:       buildinfo.BuildTime,
		ProxyConfigPath: s.deps.Bridge.ProxyPath,
		ConfigBridgeOK:  s.deps.Bridge.BridgeMessage == "",
		BridgeMessage:   s.deps.Bridge.BridgeMessage,
	}
	if s.deps.Runtime != nil {
		info.ProxyRunning = s.deps.Runtime.Status().Running
	}
	if err := database.Ping(s.ctx(), s.deps.DB); err != nil {
		info.Message = err.Error()
		return info
	}
	info.Connected = true
	ver, dirty, _ := database.MigrationVersion(s.ctx(), s.deps.DB)
	info.MigrationVersion = ver
	info.MigrationDirty = dirty
	if err := database.SchemaReady(s.ctx(), s.deps.DB); err != nil {
		info.Message = err.Error()
		return info
	}
	info.SchemaReady = true
	return info
}

func (s *DesktopService) CloseDB() {
	if s.deps.DB != nil {
		s.deps.DB.Close()
	}
}
