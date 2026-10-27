package proxy

import (
	"net/http"
	"testing"
)

func TestHTTPHeaderMap(t *testing.T) {
	h := http.Header{}
	h.Add("Content-Type", "application/json")
	h.Add("X-Custom", "a")
	h.Add("X-Custom", "b")

	got := httpHeaderMap(h)
	if got["Content-Type"] != "application/json" {
		t.Fatalf("Content-Type: %q", got["Content-Type"])
	}
	if got["X-Custom"] != "a, b" {
		t.Fatalf("X-Custom: %q", got["X-Custom"])
	}
}
