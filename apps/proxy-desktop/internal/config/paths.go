package config

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveConfigFilePath returns an absolute path for desktop config.yaml suitable for read/write.
func ResolveConfigFilePath() string {
	if v := strings.TrimSpace(os.Getenv("AI_PROXY_DESKTOP_CONFIG")); v != "" {
		if abs, err := filepath.Abs(v); err == nil {
			return abs
		}
		return v
	}
	if cwd, err := os.Getwd(); err == nil {
		if p, ok := findDesktopConfigFrom(cwd); ok {
			return p
		}
	}
	if exe, err := os.Executable(); err == nil {
		if p, ok := findDesktopConfigFrom(filepath.Dir(exe)); ok {
			return p
		}
	}
	// Last resort: cwd-relative (may fail under Wails if cwd is read-only).
	if abs, err := filepath.Abs("config.yaml"); err == nil {
		return abs
	}
	return "config.yaml"
}

func findDesktopConfigFrom(start string) (string, bool) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for i := 0; i < 16; i++ {
		if _, err := os.Stat(filepath.Join(dir, "wails.json")); err == nil {
			cfg := filepath.Join(dir, "config.yaml")
			return cfg, true
		}
		// Monorepo path when running from repo root or other packages.
		nested := filepath.Join(dir, "apps", "proxy", "proxy-desktop", "config.yaml")
		if _, err := os.Stat(nested); err == nil {
			return nested, true
		}
		nestedDir := filepath.Join(dir, "apps", "proxy", "proxy-desktop")
		if _, err := os.Stat(filepath.Join(nestedDir, "wails.json")); err == nil {
			return filepath.Join(nestedDir, "config.yaml"), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Legacy: config.yaml in cwd when wails.json not found (e.g. tests).
	if _, err := os.Stat(filepath.Join(start, "config.yaml")); err == nil {
		if abs, err := filepath.Abs(filepath.Join(start, "config.yaml")); err == nil {
			return abs, true
		}
	}
	return "", false
}
