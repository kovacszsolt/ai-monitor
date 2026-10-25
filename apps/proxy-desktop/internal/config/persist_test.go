package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWriteDatabasePath(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	dbPath := filepath.Join(dir, "custom.db")
	if err := WriteDatabasePath(cfgPath, dbPath); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), dbPath) {
		t.Fatalf("expected path in file: %s", data)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Storage.DatabasePath != dbPath {
		t.Fatalf("got %q want %q", cfg.Storage.DatabasePath, dbPath)
	}
}
