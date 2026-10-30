package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

func TestListUpstreamModels_ok(t *testing.T) {
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
		w.Write([]byte(`{"object":"list","data":[{"id":"gpt-4o","object":"model"},{"id":"gpt-4o-mini","object":"model"}]}`))
	}))
	defer srv.Close()

	p := providers.Provider{
		ProviderType:        providers.TypeOpenAI,
		BaseURL:             srv.URL,
		Credential:          "secret",
		ConnectionTimeoutMS: 2000,
		RequestTimeoutMS:    5000,
	}
	models, err := ListUpstreamModels(context.Background(), p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 2 || models[0].ID != "gpt-4o" || models[1].ID != "gpt-4o-mini" {
		t.Fatalf("models=%#v", models)
	}
}

func TestListUpstreamModels_authFailure(t *testing.T) {
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
	_, err := ListUpstreamModels(context.Background(), p)
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("expected auth error, got %v", err)
	}
}

func TestParseUpstreamModelsList_dedupes(t *testing.T) {
	body := []byte(`{"data":[{"id":"a"},{"id":"a"},{"id":"b"}]}`)
	models, err := parseUpstreamModelsList(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 2 {
		t.Fatalf("len=%d", len(models))
	}
}
