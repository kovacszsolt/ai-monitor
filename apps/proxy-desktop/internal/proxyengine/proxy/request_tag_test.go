package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestTagFromHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	if requestTagFromHeader(r) != "" {
		t.Fatal("expected empty without header")
	}
	r.Header.Set("X-AI-Proxy-Tag", "  my-tag  ")
	if requestTagFromHeader(r) != "my-tag" {
		t.Fatalf("got %q", requestTagFromHeader(r))
	}
	r.Header.Set("X-AI-Proxy-Tag", "   ")
	if requestTagFromHeader(r) != "" {
		t.Fatal("whitespace only should be empty")
	}
	long := strings.Repeat("a", maxRequestTagLen+10)
	r.Header.Set("X-AI-Proxy-Tag", long)
	got := requestTagFromHeader(r)
	if len(got) != maxRequestTagLen {
		t.Fatalf("expected len %d, got %d", maxRequestTagLen, len(got))
	}
}
