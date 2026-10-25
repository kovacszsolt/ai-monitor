package proxy

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/auth"
)

const (
	ErrorSourceProxy    = "proxy"
	ErrorSourceProvider = "provider"
)

func logProxyAuthFailure(r *http.Request, reason string, err error) {
	slog.Warn("proxy client authentication failed",
		"error_source", ErrorSourceProxy,
		"auth_stage", "proxy_api_key",
		"reason", reason,
		"path", r.URL.Path,
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
		"err", errString(err),
	)
}

func authFailureReason(err error) string {
	switch {
	case errors.Is(err, auth.ErrInvalidKey):
		return "invalid_key_format"
	case errors.Is(err, auth.ErrKeyNotFound):
		return "key_not_found"
	case errors.Is(err, auth.ErrKeyDisabled):
		return "key_disabled"
	case errors.Is(err, auth.ErrKeyExpired):
		return "key_expired"
	case errors.Is(err, auth.ErrProjectDisabled):
		return "project_disabled"
	default:
		return "authentication_failed"
	}
}

func logProxyRequestError(source string, stage string, attrs ...any) {
	args := []any{
		"error_source", source,
		"stage", stage,
	}
	args = append(args, attrs...)
	if source == ErrorSourceProxy {
		slog.Warn("proxy request error", args...)
	} else {
		slog.Warn("upstream provider response", args...)
	}
}

func logProxyRequestInfo(attrs ...any) {
	slog.Info("proxy request", attrs...)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func upstreamHost(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Host
}

// previewJSONError extracts a short, non-secret message snippet from an upstream body.
func previewJSONError(body []byte, max int) string {
	if len(body) == 0 {
		return ""
	}
	if max <= 0 {
		max = 512
	}
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		s := strings.TrimSpace(string(body))
		if len(s) > max {
			return s[:max] + "…"
		}
		return s
	}
	if e, ok := v["error"].(map[string]interface{}); ok {
		if msg, ok := e["message"].(string); ok {
			return truncate(msg, max)
		}
	}
	if msg, ok := v["message"].(string); ok {
		return truncate(msg, max)
	}
	b, _ := json.Marshal(v)
	return truncate(string(b), max)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func classifyUpstreamError(status int) string {
	switch status {
	case http.StatusUnauthorized:
		return "upstream_auth_rejected"
	case http.StatusForbidden:
		return "upstream_forbidden"
	case http.StatusTooManyRequests:
		return "upstream_rate_limited"
	case http.StatusBadRequest:
		return "upstream_bad_request"
	default:
		if status >= 500 {
			return "upstream_server_error"
		}
		return "upstream_client_error"
	}
}
