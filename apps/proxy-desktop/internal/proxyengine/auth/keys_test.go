package auth

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/projects"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/storage"
)

func testDB(t *testing.T) *sql.DB {
	dir := t.TempDir()
	cfg := config.StorageConfig{DatabasePath: filepath.Join(dir, "test.db"), JournalMode: "WAL", BusyTimeoutMS: 5000, Synchronous: "DELETE"}
	db, err := storage.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.RunMigrations(db, ""); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAPIKeyStoredPlaintext(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ks := NewKeyStore(db)
	ps := projects.NewStore(db)
	p, err := ps.Create(context.Background(), "p1", "")
	if err != nil {
		t.Fatal(err)
	}
	_, plaintext, err := ks.Create(context.Background(), p.ID, "dev", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if plaintext == "" {
		t.Fatal("expected plaintext")
	}
	row := db.QueryRow(`SELECT api_key FROM proxy_api_keys LIMIT 1`)
	var stored string
	if err := row.Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if stored != plaintext {
		t.Fatalf("expected plaintext in db, got %q", stored)
	}
	auth, err := ks.Authenticate(context.Background(), plaintext)
	if err != nil {
		t.Fatal(err)
	}
	if auth.ProjectID != p.ID {
		t.Fatalf("wrong project")
	}
}
