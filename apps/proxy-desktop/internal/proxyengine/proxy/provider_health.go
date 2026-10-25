package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

// ProviderHealthFailure describes a provider that did not pass a connectivity check.
type ProviderHealthFailure struct {
	ProviderID   string
	ProviderName string
	Message      string
}

// CheckProviderReachability issues a lightweight GET /v1/models against the provider base URL.
func CheckProviderReachability(ctx context.Context, p providers.Provider) error {
	if strings.TrimSpace(p.BaseURL) == "" {
		return fmt.Errorf("base URL is not configured")
	}
	switch p.ProviderType {
	case providers.TypeOpenAI, providers.TypeOpenAICompatible:
		if strings.TrimSpace(p.Credential) == "" && p.ProviderType == providers.TypeOpenAI {
			return fmt.Errorf("API credential is not configured")
		}
	case providers.TypeAnthropic:
		if strings.TrimSpace(p.Credential) == "" {
			return fmt.Errorf("API credential is not configured")
		}
	default:
		return fmt.Errorf("unsupported provider type %q", p.ProviderType)
	}

	connMS := p.ConnectionTimeoutMS
	if connMS <= 0 {
		connMS = 10_000
	}
	reqMS := p.RequestTimeoutMS
	if reqMS <= 0 {
		reqMS = 30_000
	}
	if reqMS > 30_000 {
		reqMS = 30_000
	}

	client := NewHTTPClient(connMS, reqMS)
	url := joinUpstreamURL(p.BaseURL, "/v1/models")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if err := applyProviderAuth(p, req); err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication failed (HTTP %d)", resp.StatusCode)
	}
	return fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
}

// ListUnhealthyProviders checks each enabled provider and returns those that fail the reachability probe.
func ListUnhealthyProviders(ctx context.Context, list []providers.Provider) []ProviderHealthFailure {
	var (
		mu  sync.Mutex
		out []ProviderHealthFailure
		wg  sync.WaitGroup
	)
	for _, p := range list {
		if !p.Enabled {
			continue
		}
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			if err := CheckProviderReachability(ctx, p); err != nil {
				mu.Lock()
				out = append(out, ProviderHealthFailure{
					ProviderID:   p.ID,
					ProviderName: p.Name,
					Message:      err.Error(),
				})
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()
	return out
}
