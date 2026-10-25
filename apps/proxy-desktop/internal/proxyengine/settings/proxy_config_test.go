package settings

import (
	"context"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

func TestSeedAndLoadProxyConfig(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	base := config.Default()
	base.Server.Port = 9090
	base.Logging.Level = "debug"
	if err := SeedFromConfig(ctx, db, base); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadProxyConfig(ctx, db, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Server.Port != 9090 {
		t.Fatalf("port %d", loaded.Server.Port)
	}
	if loaded.Logging.Level != "debug" {
		t.Fatalf("level %s", loaded.Logging.Level)
	}
	// Seed does not overwrite.
	base.Server.Port = 7777
	if err := SeedFromConfig(ctx, db, base); err != nil {
		t.Fatal(err)
	}
	loaded2, err := LoadProxyConfig(ctx, db, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if loaded2.Server.Port != 9090 {
		t.Fatalf("expected seed skip, got port %d", loaded2.Server.Port)
	}
}

func TestSaveProxyVariables(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	_, err := SaveProxyVariables(ctx, db, config.Default(), ProxyVariables{
		ServerHost:              "0.0.0.0",
		ServerPort:              18080,
		ServerMaxRequestBytes:   5 * 1024 * 1024,
		LoggingLevel:            "warn",
		LoggingFormat:           "json",
		CaptureMode:             "metadata",
		CaptureStoreContent:     true,
		CaptureMaxPromptBytes:   1000,
		CaptureMaxResponseBytes: 2000,
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := ReadProxyVariables(ctx, db, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if v.ServerPort != 18080 || !v.CaptureStoreContent {
		t.Fatalf("unexpected %+v", v)
	}
}
