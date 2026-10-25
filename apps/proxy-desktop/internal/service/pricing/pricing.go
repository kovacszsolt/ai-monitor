package pricing

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

func ParseModelRouteTokenPricing(inputUSD, cachedUSD, outputUSD string) (input, cached, output *float64, err error) {
	input, err = parseOptionalUSDPer1M(inputUSD)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("input token price: %w", err)
	}
	cached, err = parseOptionalUSDPer1M(cachedUSD)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cached input token price: %w", err)
	}
	output, err = parseOptionalUSDPer1M(outputUSD)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("output token price: %w", err)
	}
	return input, cached, output, nil
}

const tokensPerMillion = 1_000_000.0

// ModelUsageTotalCostUSD estimates spend from aggregated token sums and per-1M USD prices on the model route.
func ModelUsageTotalCostUSD(route domain.ModelRoute, inputTokens, cachedInputTokens, outputTokens int64) (float64, bool) {
	hasPrice := route.InputTokenUSDPer1M != nil || route.CachedInputTokenUSDPer1M != nil || route.OutputTokenUSDPer1M != nil
	if !hasPrice {
		return 0, false
	}
	uncached := inputTokens - cachedInputTokens
	if uncached < 0 {
		uncached = 0
	}
	var cost float64
	if route.InputTokenUSDPer1M != nil {
		cost += float64(uncached) * *route.InputTokenUSDPer1M / tokensPerMillion
	}
	cachedUSDPer1M := route.CachedInputTokenUSDPer1M
	if cachedUSDPer1M == nil {
		cachedUSDPer1M = route.InputTokenUSDPer1M
	}
	if cachedUSDPer1M != nil {
		cost += float64(cachedInputTokens) * *cachedUSDPer1M / tokensPerMillion
	}
	if route.OutputTokenUSDPer1M != nil {
		cost += float64(outputTokens) * *route.OutputTokenUSDPer1M / tokensPerMillion
	}
	return cost, true
}

func parseOptionalUSDPer1M(s string) (*float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q", s)
	}
	if v < 0 {
		return nil, errors.New("must be non-negative")
	}
	return &v, nil
}
