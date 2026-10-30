package proxy

import (
	"net/http"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

func TestApplyProviderAuth_Bifrost(t *testing.T) {
	t.Run("virtual key via x-bf-vk", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://example/v1/chat/completions", nil)
		p := providers.Provider{
			Name:         "bifrost",
			ProviderType: providers.TypeBifrost,
			Credential:   "vk-legacy-no-prefix",
		}
		if err := applyProviderAuth(p, req); err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("x-bf-vk"); got != "vk-legacy-no-prefix" {
			t.Fatalf("x-bf-vk: got %q", got)
		}
		if req.Header.Get("Authorization") != "" {
			t.Fatalf("expected no Authorization for legacy key, got %q", req.Header.Get("Authorization"))
		}
	})

	t.Run("sk-bf- prefix also sets bearer", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://example/v1/chat/completions", nil)
		p := providers.Provider{
			Name:         "bifrost",
			ProviderType: providers.TypeBifrost,
			Credential:   "sk-bf-test-key",
		}
		if err := applyProviderAuth(p, req); err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("x-bf-vk"); got != "sk-bf-test-key" {
			t.Fatalf("x-bf-vk: got %q", got)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer sk-bf-test-key" {
			t.Fatalf("Authorization: got %q", got)
		}
	})

	t.Run("openai-compatible sends x-bf-vk for Bifrost-style gateways", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://example/v1/chat/completions", nil)
		p := providers.Provider{
			Name:         "spark",
			ProviderType: providers.TypeOpenAICompatible,
			Credential:   "vk-test",
		}
		if err := applyProviderAuth(p, req); err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("x-bf-vk"); got != "vk-test" {
			t.Fatalf("x-bf-vk: got %q", got)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer vk-test" {
			t.Fatalf("Authorization: got %q", got)
		}
	})

	t.Run("missing credential", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://example/v1/chat/completions", nil)
		p := providers.Provider{
			Name:         "bifrost",
			ProviderType: providers.TypeBifrost,
		}
		if err := applyProviderAuth(p, req); err == nil {
			t.Fatal("expected error")
		}
	})
}
