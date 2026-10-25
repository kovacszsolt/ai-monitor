package capture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

func TestShouldFileDumpErrorsOnly(t *testing.T) {
	cfg := config.FileDumpConfig{Enabled: true, ErrorsOnly: true}
	if !ShouldFileDump(cfg, 200, "parse_error", StateCompleted, "") {
		t.Fatal("expected parse_error to dump")
	}
	if ShouldFileDump(cfg, 200, "ok", StateCompleted, "") {
		t.Fatal("expected ok 200 to skip")
	}
	if !ShouldFileDump(cfg, 502, "n/a", StateFailed, "upstream") {
		t.Fatal("expected http error to dump")
	}
}

func TestWriteHTTPExchangeDump(t *testing.T) {
	dir := t.TempDir()
	cfg := config.FileDumpConfig{
		Enabled:      true,
		Directory:    dir,
		MaxBodyBytes: 1024,
	}
	path, err := WriteHTTPExchangeDump(cfg, HTTPExchangeDump{
		RequestID: "req/1",
		StartedAt: "20260101T120000.000000000",
		Request:   BodyPart{Bytes: 3, Body: "hey"},
		Response:  BodyPart{Bytes: 5, Body: "world"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(path) != dir {
		t.Fatalf("path dir: %s", path)
	}
}
