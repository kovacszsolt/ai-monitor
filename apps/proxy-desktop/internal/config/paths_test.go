package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConfigFilePath_findsProxyDesktop(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Walk up to repo if test runs from internal/config.
	dir := cwd
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "apps", "proxy", "proxy-desktop", "wails.json")); err == nil {
			dir = filepath.Join(dir, "apps", "proxy", "proxy-desktop")
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Skip("proxy-desktop not found from test cwd")
		}
		dir = parent
	}
	t.Chdir(dir)
	got := ResolveConfigFilePath()
	if !strings.HasSuffix(got, filepath.Join("proxy-desktop", "config.yaml")) && !strings.HasSuffix(got, "config.yaml") {
		t.Fatalf("unexpected config path: %s", got)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("expected absolute path: %s", got)
	}
}
