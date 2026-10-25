package pricing

import (
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

func TestParseModelRouteTokenPricing(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		cached  string
		out     string
		wantErr bool
	}{
		{name: "all empty", in: "", cached: "", out: ""},
		{name: "values", in: "3", cached: "0.75", out: "15"},
		{name: "negative", in: "-1", wantErr: true},
		{name: "invalid", in: "x", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, cached, out, err := ParseModelRouteTokenPricing(tt.in, tt.cached, tt.out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.in == "" && in != nil {
				t.Fatalf("input should be nil")
			}
			if tt.in == "3" && (in == nil || *in != 3) {
				t.Fatalf("input = %v", in)
			}
			if tt.cached == "0.75" && (cached == nil || *cached != 0.75) {
				t.Fatalf("cached = %v", cached)
			}
			if tt.out == "15" && (out == nil || *out != 15) {
				t.Fatalf("output = %v", out)
			}
		})
	}
}

func TestModelUsageTotalCostUSD(t *testing.T) {
	in, cached, out := 3.0, 0.75, 15.0
	route := domain.ModelRoute{
		InputTokenUSDPer1M:       &in,
		CachedInputTokenUSDPer1M: &cached,
		OutputTokenUSDPer1M:      &out,
	}
	cost, ok := ModelUsageTotalCostUSD(route, 1_200_000, 200_000, 50_000)
	if !ok {
		t.Fatal("expected has cost")
	}
	want := 1.0*3 + 0.2*0.75 + 0.05*15
	if cost != want {
		t.Fatalf("cost = %v want %v", cost, want)
	}
}
