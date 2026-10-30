package config

import (
	"path/filepath"
	"testing"
)

func TestResolveDefaultDatabasePath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := ResolveDefaultDatabasePath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".config", "ai-proxy-manager", "proxy.db")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolveDefaultDatabasePath_XDGConfigHome(t *testing.T) {
	base := filepath.Join(t.TempDir(), "xdg-config")
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("HOME", t.TempDir())

	got, err := ResolveDefaultDatabasePath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(base, "ai-proxy-manager", "proxy.db")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
