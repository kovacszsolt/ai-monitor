package service

import (
	"context"
	"database/sql"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyruntime"
)

// ProxyRuntimeService exposes ai-proxy serve lifecycle to the Wails UI.
type ProxyRuntimeService struct {
	bridge    config.Bridge
	runtime   *proxyruntime.Runtime
	parent    context.Context
	providers *providers.Store
}

func NewProxyRuntimeService(parent context.Context, bridge config.Bridge, runtime *proxyruntime.Runtime, db *sql.DB) *ProxyRuntimeService {
	return &ProxyRuntimeService{
		bridge:    bridge,
		runtime:   runtime,
		parent:    parent,
		providers: providers.NewStore(db),
	}
}

func (s *ProxyRuntimeService) SetContext(ctx context.Context) {
	s.parent = ctx
}

func (s *ProxyRuntimeService) UpdateBridge(bridge config.Bridge) {
	s.bridge = bridge
}

func (s *ProxyRuntimeService) listUnhealthyProviders(ctx context.Context) ([]UnhealthyProviderDTO, error) {
	if s.providers == nil {
		return nil, nil
	}
	failures, err := s.providers.ListUnhealthy(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]UnhealthyProviderDTO, 0, len(failures))
	for _, f := range failures {
		out = append(out, UnhealthyProviderDTO{
			ID:      f.ID,
			Name:    f.Name,
			Message: f.Message,
		})
	}
	return out, nil
}

func (s *ProxyRuntimeService) CheckProviders() Result[ProviderHealthCheckDTO] {
	ctx := s.parent
	if ctx == nil {
		ctx = context.Background()
	}
	unhealthy, err := s.listUnhealthyProviders(ctx)
	if err != nil {
		return fail[ProviderHealthCheckDTO](err)
	}
	return ok(ProviderHealthCheckDTO{UnhealthyProviders: unhealthy})
}

func (s *ProxyRuntimeService) StartProxy() Result[StartProxyResultDTO] {
	ctx := s.parent
	if ctx == nil {
		ctx = context.Background()
	}

	unhealthy, err := s.listUnhealthyProviders(ctx)
	if err != nil {
		return fail[StartProxyResultDTO](err)
	}

	if err := s.runtime.Start(s.parent, s.bridge); err != nil {
		return fail[StartProxyResultDTO](err)
	}
	return ok(StartProxyResultDTO{
		Started:            true,
		UnhealthyProviders: unhealthy,
	})
}

func (s *ProxyRuntimeService) StopProxy() Result[bool] {
	s.runtime.Stop()
	return ok(true)
}

func (s *ProxyRuntimeService) GetProxyStatus() Result[ProxyStatusDTO] {
	st := s.runtime.Status()
	return ok(ProxyStatusDTO{
		Running:    st.Running,
		ListenAddr: st.ListenAddr,
		Error:      st.Error,
		ConfigPath: s.bridge.ProxyPath,
	})
}
