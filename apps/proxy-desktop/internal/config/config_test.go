package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProxyDBFrom(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data", "proxy.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dbPath, []byte("sqlite"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "apps", "proxy", "desktop")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	found, ok := findProxyDBFrom(sub)
	if !ok {
		t.Fatal("expected to find proxy.db")
	}
	if found != dbPath {
		t.Fatalf("got %q want %q", found, dbPath)
	}
}
