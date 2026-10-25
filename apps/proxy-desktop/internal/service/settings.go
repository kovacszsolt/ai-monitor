package service

import (
	"fmt"
	"strings"

	proxysettings "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/settings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/database"
)

type SettingsDTO struct {
	ConfigFilePath        string `json:"config_file_path"`
	DatabasePath          string `json:"database_path"`
	EffectiveDatabasePath string `json:"effective_database_path"`
	AutoResolve           bool   `json:"auto_resolve"`
	ProxyConfigPath       string `json:"proxy_config_path"`
	BridgeMessage         string `json:"bridge_message,omitempty"`
	ConfigBridgeOK        bool   `json:"config_bridge_ok"`
}

type ProxyVariablesDTO struct {
	ServerHost                string `json:"server_host"`
	ServerPort                int    `json:"server_port"`
	ServerMaxRequestBytes     int64  `json:"server_max_request_bytes"`
	LoggingLevel              string `json:"logging_level"`
	LoggingFormat             string `json:"logging_format"`
	CaptureMode               string `json:"capture_mode"`
	CaptureStoreContent       bool   `json:"capture_store_content"`
	CaptureMaxPromptBytes     int64  `json:"capture_max_prompt_bytes"`
	CaptureMaxResponseBytes   int64  `json:"capture_max_response_bytes"`
}

type UpdateProxyVariablesDTO struct {
	ServerHost                string `json:"server_host"`
	ServerPort                int    `json:"server_port"`
	ServerMaxRequestBytes     int64  `json:"server_max_request_bytes"`
	LoggingLevel              string `json:"logging_level"`
	LoggingFormat             string `json:"logging_format"`
	CaptureStoreContent       bool   `json:"capture_store_content"`
	CaptureMaxPromptBytes     int64  `json:"capture_max_prompt_bytes"`
	CaptureMaxResponseBytes   int64  `json:"capture_max_response_bytes"`
}

func (s *DesktopService) GetSettings() Result[SettingsDTO] {
	filePath, snap, err := config.LoadFileSnapshot()
	if err != nil {
		return fail[SettingsDTO](err)
	}
	stored := strings.TrimSpace(snap.Storage.DatabasePath)
	dto := SettingsDTO{
		ConfigFilePath:        filePath,
		DatabasePath:          stored,
		EffectiveDatabasePath: s.deps.Cfg.Storage.DatabasePath,
		AutoResolve:           stored == "",
		ProxyConfigPath:       s.deps.Bridge.ProxyPath,
		BridgeMessage:         s.deps.Bridge.BridgeMessage,
		ConfigBridgeOK:        s.deps.Bridge.BridgeMessage == "",
	}
	return ok(dto)
}

func (s *DesktopService) UpdateDatabasePath(databasePath string) Result[ConnectionInfoDTO] {
	if s.deps.Runtime != nil && s.deps.Runtime.Status().Running {
		return fail[ConnectionInfoDTO](fmt.Errorf("stop the proxy before changing the database path"))
	}
	filePath := config.ConfigFilePath()
	if err := config.WriteDatabasePath(filePath, databasePath); err != nil {
		return fail[ConnectionInfoDTO](err)
	}
	cfg, err := config.Load()
	if err != nil {
		return fail[ConnectionInfoDTO](err)
	}
	db, cfg, err := database.OpenWithMigrations(cfg)
	if err != nil {
		return fail[ConnectionInfoDTO](err)
	}
	bridge, err := config.LoadBridge(cfg, db)
	if err != nil {
		_ = db.Close()
		return fail[ConnectionInfoDTO](err)
	}
	if err := s.reconnect(cfg, bridge, db); err != nil {
		_ = db.Close()
		return fail[ConnectionInfoDTO](err)
	}
	if s.onBridgeReload != nil {
		s.onBridgeReload(bridge)
	}
	return ok(s.GetConnectionInfo())
}

func (s *DesktopService) GetProxyVariables() Result[ProxyVariablesDTO] {
	bootstrap, _, err := config.ProxyYAMLBootstrap(s.deps.Cfg)
	if err != nil {
		return fail[ProxyVariablesDTO](err)
	}
	v, err := proxysettings.ReadProxyVariables(s.ctx(), s.deps.DB, bootstrap)
	if err != nil {
		return fail[ProxyVariablesDTO](err)
	}
	dto := ProxyVariablesDTO{
		ServerHost:              v.ServerHost,
		ServerPort:              v.ServerPort,
		ServerMaxRequestBytes:   v.ServerMaxRequestBytes,
		LoggingLevel:            v.LoggingLevel,
		LoggingFormat:           v.LoggingFormat,
		CaptureMode:             v.CaptureMode,
		CaptureStoreContent:     v.CaptureStoreContent,
		CaptureMaxPromptBytes:   v.CaptureMaxPromptBytes,
		CaptureMaxResponseBytes: v.CaptureMaxResponseBytes,
	}
	return ok(dto)
}

func (s *DesktopService) UpdateProxyVariables(in UpdateProxyVariablesDTO) Result[ProxyVariablesDTO] {
	if s.deps.Runtime != nil && s.deps.Runtime.Status().Running {
		return fail[ProxyVariablesDTO](fmt.Errorf("stop the proxy before changing proxy settings"))
	}
	bootstrap, _, err := config.ProxyYAMLBootstrap(s.deps.Cfg)
	if err != nil {
		return fail[ProxyVariablesDTO](err)
	}
	_, err = proxysettings.SaveProxyVariables(s.ctx(), s.deps.DB, bootstrap, proxysettings.ProxyVariables{
		ServerHost:              in.ServerHost,
		ServerPort:              in.ServerPort,
		ServerMaxRequestBytes:   in.ServerMaxRequestBytes,
		LoggingLevel:            in.LoggingLevel,
		LoggingFormat:           in.LoggingFormat,
		CaptureMode:             "metadata",
		CaptureStoreContent:     in.CaptureStoreContent,
		CaptureMaxPromptBytes:   in.CaptureMaxPromptBytes,
		CaptureMaxResponseBytes: in.CaptureMaxResponseBytes,
	})
	if err != nil {
		return fail[ProxyVariablesDTO](err)
	}
	bridge, err := config.LoadBridge(s.deps.Cfg, s.deps.DB)
	if err != nil {
		return fail[ProxyVariablesDTO](err)
	}
	s.deps.Bridge = bridge
	s.deps.RebindDB(s.deps.DB)
	if s.onBridgeReload != nil {
		s.onBridgeReload(bridge)
	}
	return s.GetProxyVariables()
}
