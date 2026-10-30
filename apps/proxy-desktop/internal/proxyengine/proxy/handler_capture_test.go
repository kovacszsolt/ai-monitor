package proxy

import (
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

func TestExtractPromptCaptureSkipsOversizedBody(t *testing.T) {
	h := &Handler{
		svc: Services{
			Capture: config.CaptureConfig{
				StoreContent:   true,
				MaxPromptBytes: 100,
			},
		},
	}
	body := make([]byte, 200)
	for i := range body {
		body[i] = 'a'
	}
	prompt, errs := h.extractPromptCapture(routing.ProtocolOpenAIChat, body, "req-1")
	if prompt != nil || len(errs) != 0 {
		t.Fatalf("expected skip, got prompt=%v errs=%v", prompt, errs)
	}
}

func TestExtractPromptCaptureStoresSmallBody(t *testing.T) {
	h := &Handler{
		svc: Services{
			Capture: config.CaptureConfig{
				StoreContent:   true,
				MaxPromptBytes: 4096,
			},
		},
	}
	body := []byte(`{"model":"m","messages":[{"role":"user","content":"hi"}]}`)
	prompt, errs := h.extractPromptCapture(routing.ProtocolOpenAIChat, body, "req-1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errs: %v", errs)
	}
	if prompt == nil || prompt.LastUserText != "hi" {
		t.Fatalf("expected prompt capture, got %v", prompt)
	}
}
