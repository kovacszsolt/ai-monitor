package metrics

import (
	"math"
	"testing"
)

func TestUncachedAndCacheHitRate(t *testing.T) {
	u := UsageInput{Input: ptr(10000), CachedInput: ptr(8000)}
	r := Compute(u, TimingInput{}, true)
	if r.UncachedInputTokens == nil || *r.UncachedInputTokens != 2000 {
		t.Fatalf("uncached: %v", r.UncachedInputTokens)
	}
	if r.CacheHitRate == nil || math.Abs(*r.CacheHitRate-80) > 0.001 {
		t.Fatalf("cache rate: %v", r.CacheHitRate)
	}
}

func TestGenerationTPS(t *testing.T) {
	decode := int64(30000)
	u := UsageInput{Output: ptr(1500)}
	timing := TimingInput{DecodeDurationMS: &decode}
	r := Compute(u, timing, true)
	if r.GenerationTPS == nil || math.Abs(*r.GenerationTPS-50) > 0.001 {
		t.Fatalf("generation tps: %v", r.GenerationTPS)
	}
}

func TestPromptTPS(t *testing.T) {
	prefill := int64(10000)
	u := UsageInput{Input: ptr(10000), CachedInput: ptr(8000)}
	timing := TimingInput{PrefillDurationMS: &prefill}
	r := Compute(u, timing, true)
	if r.PromptTPS == nil || math.Abs(*r.PromptTPS-200) > 0.001 {
		t.Fatalf("prompt tps: %v", r.PromptTPS)
	}
}

func TestEffectiveTPS(t *testing.T) {
	total := int64(40000)
	u := UsageInput{Input: ptr(10000), CachedInput: ptr(8000), Output: ptr(1500)}
	timing := TimingInput{TotalDurationMS: &total}
	r := Compute(u, timing, true)
	if r.EffectiveTPS == nil || math.Abs(*r.EffectiveTPS-287.5) > 0.001 {
		t.Fatalf("effective tps: %v", r.EffectiveTPS)
	}
}

func TestEdgeCases(t *testing.T) {
	cases := []struct {
		name string
		u    UsageInput
		t    TimingInput
		ok   func(Result) bool
	}{
		{
			name: "input zero cache rate null",
			u:    UsageInput{Input: ptr(0), CachedInput: ptr(0)},
			ok: func(r Result) bool { return r.CacheHitRate == nil },
		},
		{
			name: "cached unknown",
			u:    UsageInput{Input: ptr(100), CachedInput: nil},
			ok: func(r Result) bool { return r.UncachedInputTokens == nil && r.CacheHitRate == nil },
		},
		{
			name: "cached exceeds input",
			u:    UsageInput{Input: ptr(10), CachedInput: ptr(20)},
			ok: func(r Result) bool { return r.UncachedInputTokens == nil && r.CacheHitRate == nil },
		},
		{
			name: "decode zero",
			u:    UsageInput{Output: ptr(10)},
			t:    TimingInput{DecodeDurationMS: ptr(0)},
			ok: func(r Result) bool { return r.GenerationTPS == nil },
		},
		{
			name: "not completed no tps",
			u:    UsageInput{Input: ptr(10), Output: ptr(5)},
			t:    TimingInput{TotalDurationMS: ptr(1000)},
			ok: func(r Result) bool {
				return r.EffectiveTPS == nil && r.GenerationTPS == nil
			},
		},
		{
			name: "missing output",
			u:    UsageInput{Input: ptr(10)},
			t:    TimingInput{TotalDurationMS: ptr(1000), DecodeDurationMS: ptr(500)},
			ok: func(r Result) bool { return r.GenerationTPS == nil },
		},
	}
	for _, tc := range cases {
		r := Compute(tc.u, tc.t, true)
		if tc.name == "not completed no tps" {
			r = Compute(tc.u, tc.t, false)
		}
		if !tc.ok(r) {
			t.Fatalf("%s: unexpected %+v", tc.name, r)
		}
	}
}

func TestDecodeDurationFromStream(t *testing.T) {
	first := int64(100)
	last := int64(400)
	d := DecodeDurationFromStream(&first, &last)
	if d == nil || *d != 300 {
		t.Fatalf("decode duration %v", d)
	}
	if DecodeDurationFromStream(&last, &first) != nil {
		t.Fatal("expected nil for invalid order")
	}
}

func ptr(n int64) *int64 { return &n }
