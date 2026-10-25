package config

import (
	"os"
	"path/filepath"
	"testing"

	proxyconfig "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"
)

func TestValidateBridgeAlignment_relativeProxyPath(t *testing.T) {
	desktop := Config{
		Storage: StorageConfig{DatabasePath: "/repo/data/proxy.db"},
	}
	proxy := proxyconfig.Default()
	proxy.Storage.DatabasePath = "../../data/proxy.db"
	proxyConfig := "/repo/apps/proxy-desktop/proxy.bootstrap.yaml"
	if msg := validateBridgeAlignment(desktop, proxy, proxyConfig); msg != "" {
		t.Fatalf("expected relative proxy path to resolve: %s", msg)
	}
}

func TestValidateBridgeAlignment_ok(t *testing.T) {
	desktop := Config{
		Storage: StorageConfig{DatabasePath: "/tmp/proxy.db"},
	}
	proxy := proxyconfig.Default()
	proxy.Storage.DatabasePath = "/tmp/proxy.db"
	if msg := validateBridgeAlignment(desktop, proxy, "/tmp/config.yaml"); msg != "" {
		t.Fatalf("unexpected bridge message: %s", msg)
	}
}

func TestNormalizePath(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "data", "proxy.db")
	if normalizePath(p) == "" {
		t.Fatal("expected non-empty")
	}
}

func TestResolveProxyConfigPath_fromEnv(t *testing.T) {
	f := filepath.Join(t.TempDir(), "proxy.yaml")
	if err := os.WriteFile(f, []byte("server:\n  port: 9090\nstorage:\n  database_path: /x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AI_PROXY_CONFIG", f)
	got, err := ResolveProxyConfigPath()
	if err != nil || got != f {
		t.Fatalf("got %q err %v", got, err)
	}
}
