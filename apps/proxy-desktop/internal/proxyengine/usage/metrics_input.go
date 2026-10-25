package usage

import "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/metrics"

// ToMetricsInput converts captured usage to metrics.UsageInput.
func ToMetricsInput(u Usage) metrics.UsageInput {
	return metrics.UsageInput{
		Input:       u.InputTokens,
		CachedInput: u.CachedInputTokens,
		Output:      u.OutputTokens,
	}
}
