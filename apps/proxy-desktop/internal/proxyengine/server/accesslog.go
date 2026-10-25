package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func accessLogMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			level := slog.LevelInfo
			if status >= 500 {
				level = slog.LevelError
			} else if status >= 400 {
				level = slog.LevelWarn
			}

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", middleware.GetReqID(r.Context()),
				"bytes", ww.BytesWritten(),
			}
			if src := ww.Header().Get("X-AI-Proxy-Error-Source"); src != "" {
				attrs = append(attrs, "error_source", src)
			} else if status == http.StatusUnauthorized || status == http.StatusForbidden {
				// Response may be from upstream without our header — still log for triage.
				if !strings.HasPrefix(r.URL.Path, "/health/") {
					attrs = append(attrs, "error_source_hint", "check_json_body_or_prior_logs")
				}
			}

			slog.Log(r.Context(), level, "http request completed", attrs...)
		})
	}
}
