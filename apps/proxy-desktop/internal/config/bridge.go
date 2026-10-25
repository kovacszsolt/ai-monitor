package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	proxyconfig "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"
	proxysettings "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/settings"
)

var ErrBridgeMismatch = errors.New("desktop and proxy config disagree on database path")

// Bridge holds desktop UI config and optional bootstrap YAML for serve/capture/server settings.
type Bridge struct {
	Desktop       Config
	ProxyPath     string
	ProxyFile     proxyconfig.Config
	ProxyMerged   proxyconfig.Config
	BridgeMessage string
}

// ProxyYAMLBootstrap loads proxy YAML (or defaults) with desktop storage for settings seed/load.
func ProxyYAMLBootstrap(desktop Config) (proxyconfig.Config, string, error) {
	proxyPath, err := ResolveProxyConfigPath()
	if err != nil {
		return proxyconfig.Config{}, "", err
	}
	var proxyFile proxyconfig.Config
	if proxyPath != "" {
		proxyFile, err = proxyconfig.Load(proxyPath)
		if err != nil {
			return proxyconfig.Config{}, "", fmt.Errorf("load proxy config: %w", err)
		}
	} else {
		proxyFile = proxyconfig.Default()
		proxyconfig.ApplyEnvOverrides(&proxyFile)
	}
	merged := proxyFile
	merged.Storage = proxyStorageFromDesktop(desktop.Storage)
	return merged, proxyPath, nil
}

// LoadBridge loads the proxy YAML and overlays runtime settings from the database.
func LoadBridge(desktop Config, db *sql.DB) (Bridge, error) {
	merged, proxyPath, err := ProxyYAMLBootstrap(desktop)
	if err != nil {
		return Bridge{}, err
	}
	if db != nil {
		ctx := context.Background()
		if err := proxysettings.SeedFromConfig(ctx, db, merged); err != nil {
			return Bridge{}, fmt.Errorf("seed proxy settings: %w", err)
		}
		merged, err = proxysettings.LoadProxyConfig(ctx, db, merged)
		if err != nil {
			return Bridge{}, fmt.Errorf("load proxy settings: %w", err)
		}
	}
	if err := merged.Validate(); err != nil {
		return Bridge{}, fmt.Errorf("proxy config (merged): %w", err)
	}
	// Reload raw file cfg for alignment checks (without DB overlay).
	var proxyFileForAlign proxyconfig.Config
	if proxyPath != "" {
		proxyFileForAlign, err = proxyconfig.Load(proxyPath)
		if err != nil {
			return Bridge{}, fmt.Errorf("load proxy config: %w", err)
		}
	} else {
		proxyFileForAlign = proxyconfig.Default()
		proxyconfig.ApplyEnvOverrides(&proxyFileForAlign)
	}
	b := Bridge{
		Desktop:     desktop,
		ProxyPath:   proxyPath,
		ProxyFile:   proxyFileForAlign,
		ProxyMerged: merged,
	}
	if msg := validateBridgeAlignment(desktop, proxyFileForAlign, proxyPath); msg != "" {
		b.BridgeMessage = msg
	}
	return b, nil
}

func validateBridgeAlignment(desktop Config, proxyFile proxyconfig.Config, proxyConfigPath string) string {
	var issues []string
	if proxyFile.Storage.DatabasePath != "" {
		proxyDB := resolvePathRelativeTo(proxyFile.Storage.DatabasePath, proxyConfigDir(proxyConfigPath))
		if normalizePath(desktop.Storage.DatabasePath) != normalizePath(proxyDB) {
			issues = append(issues, fmt.Sprintf(
				"database_path mismatch (desktop=%s proxy=%s)",
				desktop.Storage.DatabasePath, proxyFile.Storage.DatabasePath,
			))
		}
	}
	if len(issues) == 0 {
		return ""
	}
	return strings.Join(issues, "; ")
}

func proxyStorageFromDesktop(s StorageConfig) proxyconfig.StorageConfig {
	return proxyconfig.StorageConfig{
		DatabasePath:  s.DatabasePath,
		JournalMode:   s.JournalMode,
		BusyTimeoutMS: s.BusyTimeoutMS,
		Synchronous:   s.Synchronous,
	}
}

func proxyConfigDir(proxyConfigPath string) string {
	if strings.TrimSpace(proxyConfigPath) == "" {
		return ""
	}
	return filepath.Dir(proxyConfigPath)
}

func resolvePathRelativeTo(p, baseDir string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if filepath.IsAbs(p) || baseDir == "" {
		return p
	}
	return filepath.Clean(filepath.Join(baseDir, p))
}

func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return filepath.Clean(p)
	}
	return abs
}

// ResolveProxyConfigPath returns the YAML path for proxy bootstrap settings, or "" if none found.
func ResolveProxyConfigPath() (string, error) {
	if v := strings.TrimSpace(os.Getenv("AI_PROXY_CONFIG")); v != "" {
		return v, nil
	}
	candidates := []string{
		"proxy.bootstrap.yaml",
		"proxy.bootstrap.example.yaml",
		"config.yaml",
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "proxy.bootstrap.yaml"),
			filepath.Join(cwd, "apps/proxy-desktop/proxy.bootstrap.yaml"),
		)
		for dir := cwd; ; dir = filepath.Dir(dir) {
			candidates = append(candidates, filepath.Join(dir, "apps/proxy-desktop/proxy.bootstrap.yaml"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	seen := make(map[string]struct{})
	for _, c := range candidates {
		c = filepath.Clean(c)
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", nil
}
