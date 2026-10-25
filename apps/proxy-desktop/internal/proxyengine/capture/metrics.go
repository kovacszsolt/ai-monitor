package capture

import (
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/metrics"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/usage"
)

// ApplyTokenMetrics fills computed metric fields on the record from usage and timing inputs.
func ApplyTokenMetrics(r *Record, u usage.Usage, timing metrics.TimingInput, completed bool, generationEstimated, promptEstimated bool) {
	res := metrics.Compute(usage.ToMetricsInput(u), timing, completed)
	res.ApplyEstimationFlags(generationEstimated, promptEstimated)

	r.UncachedInputTokens = res.UncachedInputTokens
	if res.PrefillDurationMS != nil {
		r.PrefillDurationMS = res.PrefillDurationMS
	}
	if res.DecodeDurationMS != nil {
		r.DecodeDurationMS = res.DecodeDurationMS
	}
	r.PromptTPS = res.PromptTPS
	r.GenerationTPS = res.GenerationTPS
	r.EffectiveTPS = res.EffectiveTPS
	r.CacheHitRate = res.CacheHitRate
	r.GenerationSpeedEstimated = res.GenerationSpeedEstimated
	r.PromptSpeedEstimated = res.PromptSpeedEstimated
}
