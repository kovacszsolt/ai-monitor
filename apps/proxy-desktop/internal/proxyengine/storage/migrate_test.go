package storage

import (
	"path/filepath"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/migrations"
)

func TestMigrateEmbedded(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cfg := config.StorageConfig{
		DatabasePath:  dbPath,
		JournalMode:   "WAL",
		BusyTimeoutMS: 5000,
		Synchronous:   "NORMAL",
	}
	db, err := Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := MigrateFS(db, migrations.FS, "."); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='projects'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("projects table missing after migrate: %v n=%d", err, n)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM providers`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("expected 4 default providers on fresh migrate, got %d", n)
	}
	var gs4URL string
	if err := db.QueryRow(`SELECT base_url FROM providers WHERE name = 'GS4'`).Scan(&gs4URL); err != nil {
		t.Fatal(err)
	}
	if gs4URL != "http://127.0.0.1:8003/v1" {
		t.Fatalf("unexpected GS4 base_url: %q", gs4URL)
	}
}
