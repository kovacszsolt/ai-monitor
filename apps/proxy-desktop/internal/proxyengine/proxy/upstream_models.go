package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
)

// UpstreamModel is one model id from an OpenAI-style GET /v1/models list.
type UpstreamModel struct {
	ID string
}

// ListUpstreamModels fetches model ids from the provider's GET /v1/models endpoint.
func ListUpstreamModels(ctx context.Context, p providers.Provider) ([]UpstreamModel, error) {
	if err := validateProviderForModelsList(p); err != nil {
		return nil, err
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
		return nil, err
	}
	if err := applyProviderAuth(p, req); err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("read upstream models: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("authentication failed (HTTP %d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
	}

	return parseUpstreamModelsList(body)
}

func validateProviderForModelsList(p providers.Provider) error {
	if strings.TrimSpace(p.BaseURL) == "" {
		return fmt.Errorf("base URL is not configured")
	}
	switch p.ProviderType {
	case providers.TypeOpenAI, providers.TypeOpenAICompatible:
		if strings.TrimSpace(p.Credential) == "" && p.ProviderType == providers.TypeOpenAI {
			return fmt.Errorf("API credential is not configured")
		}
	case providers.TypeAnthropic, providers.TypeBifrost:
		if strings.TrimSpace(p.Credential) == "" {
			return fmt.Errorf("API credential is not configured")
		}
	default:
		return fmt.Errorf("unsupported provider type %q", p.ProviderType)
	}
	return nil
}

func parseUpstreamModelsList(body []byte) ([]UpstreamModel, error) {
	var v struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("invalid models list json: %w", err)
	}
	out := make([]UpstreamModel, 0, len(v.Data))
	seen := make(map[string]struct{}, len(v.Data))
	for _, item := range v.Data {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, UpstreamModel{ID: id})
	}
	return out, nil
}
