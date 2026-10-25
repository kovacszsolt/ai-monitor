package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/proxy"
)

type Server struct {
	httpServer *http.Server
	handler    *proxy.Handler
}

func New(addr string, h *proxy.Handler) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(accessLogMiddleware())
	r.Use(corsMiddleware)
	r.Use(h.AuthMiddleware)
	h.Register(r)
	return &Server{
		handler: h,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	slog.Info("listening", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
