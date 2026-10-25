package providers

import (
	"context"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/proxy"
)

// UnhealthyProvider is a provider that failed a pre-start connectivity check.
type UnhealthyProvider struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

// ListUnhealthy runs a lightweight reachability check for each enabled provider.
func (s *Store) ListUnhealthy(ctx context.Context) ([]UnhealthyProvider, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	failures := proxy.ListUnhealthyProviders(ctx, list)
	out := make([]UnhealthyProvider, 0, len(failures))
	for _, f := range failures {
		out = append(out, UnhealthyProvider{
			ID:      f.ProviderID,
			Name:    f.ProviderName,
			Message: f.Message,
		})
	}
	return out, nil
}
