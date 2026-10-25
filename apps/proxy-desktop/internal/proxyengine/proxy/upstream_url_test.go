package proxy

import "testing"

func TestJoinUpstreamURL(t *testing.T) {
	tests := []struct {
		base, endpoint, want string
	}{
		{"http://127.0.0.1:8001", "/v1/chat/completions", "http://127.0.0.1:8001/v1/chat/completions"},
		{"http://127.0.0.1:8001/v1", "/v1/chat/completions", "http://127.0.0.1:8001/v1/chat/completions"},
		{"http://127.0.0.1:8001/v1/", "/v1/chat/completions", "http://127.0.0.1:8001/v1/chat/completions"},
		{"https://api.openai.com", "/v1/chat/completions", "https://api.openai.com/v1/chat/completions"},
	}
	for _, tc := range tests {
		got := joinUpstreamURL(tc.base, tc.endpoint)
		if got != tc.want {
			t.Fatalf("joinUpstreamURL(%q, %q) = %q, want %q", tc.base, tc.endpoint, got, tc.want)
		}
	}
}
