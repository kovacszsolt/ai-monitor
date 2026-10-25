package metrics

import "math"

// UsageInput holds nullable token counts from provider-normalized usage.
type UsageInput struct {
	Input       *int64
	CachedInput *int64
	Output      *int64
}

// TimingInput holds nullable timing measurements in milliseconds.
type TimingInput struct {
	TotalDurationMS   *int64
	TTFTMS            *int64
	PrefillDurationMS *int64
	DecodeDurationMS  *int64
}

// Result is the computed LLM token throughput metrics snapshot.
type Result struct {
	UncachedInputTokens *int64
	PrefillDurationMS   *int64
	DecodeDurationMS    *int64

	PromptTPS           *float64
	GenerationTPS       *float64
	EffectiveTPS        *float64
	CacheHitRate        *float64

	GenerationSpeedEstimated bool
	PromptSpeedEstimated     bool

	Completed bool
}

// Compute derives metrics from usage and timing. Missing data stays nil.
func Compute(usage UsageInput, timing TimingInput, completed bool) Result {
	out := Result{Completed: completed}

	if timing.PrefillDurationMS != nil && *timing.PrefillDurationMS > 0 {
		out.PrefillDurationMS = timing.PrefillDurationMS
	}
	if timing.DecodeDurationMS != nil && *timing.DecodeDurationMS > 0 {
		out.DecodeDurationMS = timing.DecodeDurationMS
	}

	uncached, cacheRate := deriveUncachedAndCacheRate(usage.Input, usage.CachedInput)
	out.UncachedInputTokens = uncached
	out.CacheHitRate = cacheRate

	if completed {
		out.EffectiveTPS = effectiveTPS(usage.Input, usage.Output, timing.TotalDurationMS)
		out.GenerationTPS = generationTPS(usage.Output, out.DecodeDurationMS)
		out.PromptTPS = promptTPS(out.UncachedInputTokens, out.PrefillDurationMS)
	}

	return out
}

// ApplyEstimationFlags marks whether TPS values were estimated (caller sets after Compute).
func (r *Result) ApplyEstimationFlags(generationEstimated, promptEstimated bool) {
	r.GenerationSpeedEstimated = generationEstimated && r.GenerationTPS != nil
	r.PromptSpeedEstimated = promptEstimated && r.PromptTPS != nil
}

func deriveUncachedAndCacheRate(input, cached *int64) (*int64, *float64) {
	if input == nil || *input < 0 {
		return nil, nil
	}
	if cached == nil {
		return nil, nil
	}
	if *cached < 0 || *cached > *input {
		return nil, nil
	}
	uncached := *input - *cached
	if *input == 0 {
		return ptrInt64(0), nil
	}
	rate := float64(*cached) / float64(*input) * 100.0
	if !finite(rate) || rate < 0 || rate > 100 {
		return ptrInt64(uncached), nil
	}
	return ptrInt64(uncached), &rate
}

func effectiveTPS(input, output *int64, totalMS *int64) *float64 {
	if totalMS == nil || *totalMS <= 0 {
		return nil
	}
	if input == nil || output == nil {
		return nil
	}
	if *input < 0 || *output < 0 {
		return nil
	}
	tokens := float64(*input + *output)
	sec := float64(*totalMS) / 1000.0
	return divTPS(tokens, sec)
}

func generationTPS(output *int64, decodeMS *int64) *float64 {
	if output == nil || *output <= 0 {
		return nil
	}
	if decodeMS == nil || *decodeMS <= 0 {
		return nil
	}
	sec := float64(*decodeMS) / 1000.0
	return divTPS(float64(*output), sec)
}

func promptTPS(uncached *int64, prefillMS *int64) *float64 {
	if uncached == nil || *uncached <= 0 {
		return nil
	}
	if prefillMS == nil || *prefillMS <= 0 {
		return nil
	}
	sec := float64(*prefillMS) / 1000.0
	return divTPS(float64(*uncached), sec)
}

func divTPS(tokens, seconds float64) *float64 {
	if seconds <= 0 || tokens < 0 {
		return nil
	}
	v := tokens / seconds
	if !finite(v) {
		return nil
	}
	return &v
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func ptrInt64(n int64) *int64 {
	return &n
}

// DecodeDurationFromStream returns decode window ms when first and last output token times are known.
func DecodeDurationFromStream(firstMS, lastMS *int64) *int64 {
	if firstMS == nil || lastMS == nil {
		return nil
	}
	if *lastMS <= *firstMS {
		return nil
	}
	d := *lastMS - *firstMS
	return &d
}
