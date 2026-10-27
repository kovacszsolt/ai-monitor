package repository

import (
	"testing"
)

func TestNormalizeChartBucketUnit(t *testing.T) {
	tests := []struct {
		in   string
		want ChartBucketUnit
	}{
		{"", ChartBucketHour},
		{"hour", ChartBucketHour},
		{"minute", ChartBucketMinute},
		{"15m", ChartBucketQuarter},
		{"quarter", ChartBucketQuarter},
		{"day", ChartBucketDay},
		{"bogus", ChartBucketHour},
	}
	for _, tc := range tests {
		if got := NormalizeChartBucketUnit(tc.in); got != tc.want {
			t.Fatalf("NormalizeChartBucketUnit(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidateChartBucketForBounds(t *testing.T) {
	bounds7d := hourlyChartBounds{Days: 7}
	if err := validateChartBucketForBounds(ChartBucketHour, bounds7d); err != nil {
		t.Fatalf("hour/7d: %v", err)
	}
	if err := validateChartBucketForBounds(ChartBucketMinute, bounds7d); err == nil {
		t.Fatal("expected minute/7d to fail")
	}
	bounds1d := hourlyChartBounds{Days: 1}
	if err := validateChartBucketForBounds(ChartBucketMinute, bounds1d); err != nil {
		t.Fatalf("minute/1d: %v", err)
	}
}
