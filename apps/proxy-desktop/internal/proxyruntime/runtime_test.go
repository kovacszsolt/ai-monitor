package proxyruntime

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	proxyconfig "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
)

func TestRuntimeStartStop(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "proxy.db")
	desktop := config.Config{
		Storage: config.StorageConfig{
			DatabasePath:  dbPath,
			JournalMode:   "WAL",
			BusyTimeoutMS: 5000,
			Synchronous:   "NORMAL",
		},
	}
	merged := proxyconfig.Default()
	merged.Storage.DatabasePath = dbPath
	merged.Server.Port = 18080
	bridge := config.Bridge{Desktop: desktop, ProxyMerged: merged}

	rt := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rt.Start(ctx, bridge); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if rt.Status().Running {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !rt.Status().Running {
		t.Fatal("expected proxy running")
	}
	rt.Stop()
	time.Sleep(200 * time.Millisecond)
	if rt.Status().Running {
		t.Fatal("expected proxy stopped")
	}
}
