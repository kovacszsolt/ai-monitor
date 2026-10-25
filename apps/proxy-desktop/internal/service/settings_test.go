package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/storage"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/database"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyruntime"
)

func TestUpdateProxyVariables_proxyRunning(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "proxy.db")
	desktopCfg := config.Default()
	desktopCfg.Storage.DatabasePath = dbPath
	db, err := database.Open(desktopCfg.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := storage.MigrateEmbedded(db); err != nil {
		t.Fatal(err)
	}
	bridge, err := config.LoadBridge(desktopCfg, db)
	if err != nil {
		t.Fatal(err)
	}
	rt := proxyruntime.New()
	if err := rt.Start(context.Background(), bridge); err != nil {
		t.Fatal(err)
	}
	defer rt.Stop()

	svc := NewDesktopService(desktopCfg, bridge, db, rt)
	res := svc.UpdateProxyVariables(UpdateProxyVariablesDTO{ServerPort: 9999})
	if res.OK {
		t.Fatal("expected error when proxy running")
	}
	if res.OK || res.Error.Message == "" {
		t.Fatalf("expected error message, got %+v", res)
	}
}
