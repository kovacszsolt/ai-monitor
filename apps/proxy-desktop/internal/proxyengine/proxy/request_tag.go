package proxy

import (
	"log/slog"
	"net/http"
	"strings"
)

const maxRequestTagLen = 256

func requestTagFromHeader(r *http.Request) string {
	s := strings.TrimSpace(r.Header.Get("X-AI-Proxy-Tag"))
	if s == "" {
		return ""
	}
	if len(s) > maxRequestTagLen {
		slog.Debug("request tag truncated", "max_len", maxRequestTagLen)
		return s[:maxRequestTagLen]
	}
	return s
}
