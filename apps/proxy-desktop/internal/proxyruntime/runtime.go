package proxyruntime

import (
	"context"
	"errors"
	"fmt"
	"sync"

	proxyapp "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/app"
	proxyconfig "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/config"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
)

// Runtime runs ai-proxy serve in-process (separate DB connection from the desktop UI).
type Runtime struct {
	mu     sync.Mutex
	cfg    proxyconfig.Config
	app    *proxyapp.App
	runCtx context.Context
	cancel context.CancelFunc
	err    error
}

func New() *Runtime {
	return &Runtime{}
}

// Start loads the config bridge, opens the proxy app, and serves HTTP until Stop or parent ctx cancel.
func (r *Runtime) Start(parent context.Context, bridge config.Bridge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.app != nil {
		return fmt.Errorf("proxy already running")
	}
	cfg := bridge.ProxyMerged
	proxyapp.SetupLogger(cfg.Logging)
	a, err := proxyapp.Open(cfg, proxyapp.DefaultMigrationsDir())
	if err != nil {
		return err
	}
	runCtx, cancel := context.WithCancel(parent)
	r.cfg = cfg
	r.app = a
	r.runCtx = runCtx
	r.cancel = cancel
	r.err = nil
	go func() {
		err := a.Run(runCtx)
		r.mu.Lock()
		defer r.mu.Unlock()
		if err != nil && !errors.Is(err, context.Canceled) {
			r.err = err
		}
		r.app = nil
		r.cancel = nil
	}()
	return nil
}

// Stop shuts down a running proxy instance.
func (r *Runtime) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel != nil {
		r.cancel()
	}
}

// Status reports whether the proxy HTTP server is running.
func (r *Runtime) Status() Status {
	r.mu.Lock()
	defer r.mu.Unlock()
	st := Status{
		Running: r.app != nil,
	}
	if r.app != nil {
		st.ListenAddr = r.cfg.ListenAddr()
	}
	if r.err != nil {
		st.Error = r.err.Error()
	}
	return st
}

// WriterStats returns capture writer queue metrics when the proxy is running.
func (r *Runtime) WriterStats() (inserted, dropped int64, queueLen int, ok bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.app == nil {
		return 0, 0, 0, false
	}
	ins, dropped, qlen := r.app.Writer.Stats()
	return ins, dropped, qlen, true
}

type Status struct {
	Running    bool
	ListenAddr string
	Error      string
}
