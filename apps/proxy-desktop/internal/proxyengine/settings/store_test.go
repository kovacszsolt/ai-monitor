package settings

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/storage"
)

func TestSettingsGetSet(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	if err := Set(ctx, db, "key_pepper", "pepper-a"); err != nil {
		t.Fatal(err)
	}
	v, err := Get(ctx, db, "key_pepper")
	if err != nil || v != "pepper-a" {
		t.Fatalf("got %q err %v", v, err)
	}
	if err := Set(ctx, db, "key_pepper", "pepper-b"); err != nil {
		t.Fatal(err)
	}
	v, err = Get(ctx, db, "key_pepper")
	if err != nil || v != "pepper-b" {
		t.Fatalf("got %q err %v", v, err)
	}
}

func TestMigrateLegacyKeyPepper(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	_, err := db.Exec(`CREATE TABLE proxy_key_pepper (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		pepper TEXT NOT NULL,
		updated_at TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO proxy_key_pepper (id, pepper, updated_at) VALUES (1, 'legacy-pepper', '2020-01-01T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}
	if err := migrateLegacyKeyPepper(ctx, db); err != nil {
		t.Fatal(err)
	}
	v, err := Get(ctx, db, KeyKeyPepper)
	if err != nil || v != "legacy-pepper" {
		t.Fatalf("got %q err %v", v, err)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dir := t.TempDir()
	cfg := config.StorageConfig{
		DatabasePath:  filepath.Join(dir, "test.db"),
		JournalMode:   "WAL",
		BusyTimeoutMS: 5000,
		Synchronous:   "NORMAL",
	}
	db, err := storage.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.MigrateEmbedded(db); err != nil {
		t.Fatal(err)
	}
	return db
}
