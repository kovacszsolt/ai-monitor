package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

func TestCheckProviderReachability_ok(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	p := providers.Provider{
		ID:                  "id",
		Name:                "test",
		ProviderType:        providers.TypeOpenAI,
		BaseURL:             srv.URL,
		Credential:          "secret",
		Enabled:             true,
		ConnectionTimeoutMS: 2000,
		RequestTimeoutMS:    5000,
	}
	if err := CheckProviderReachability(context.Background(), p); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestCheckProviderReachability_authFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer srv.Close()

	p := providers.Provider{
		ProviderType:        providers.TypeOpenAI,
		BaseURL:             srv.URL,
		Credential:          "bad",
		ConnectionTimeoutMS: 2000,
		RequestTimeoutMS:    5000,
	}
	err := CheckProviderReachability(context.Background(), p)
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestListUnhealthyProviders_skipsDisabled(t *testing.T) {
	failures := ListUnhealthyProviders(context.Background(), []providers.Provider{
		{
			Name:                "disabled",
			ProviderType:        providers.TypeOpenAI,
			BaseURL:             "http://127.0.0.1:1",
			Credential:          "x",
			Enabled:             false,
			ConnectionTimeoutMS: 100,
			RequestTimeoutMS:    100,
		},
	})
	if len(failures) != 0 {
		t.Fatalf("expected no checks for disabled provider, got %d", len(failures))
	}
}
