package repository

import (
	"testing"
	"time"
)

func TestResolveHourlyChartBoundsPresetDays(t *testing.T) {
	b, err := resolveHourlyChartBounds(7, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if b.Days != 7 || b.UsedCustomRange {
		t.Fatalf("unexpected bounds: %+v", b)
	}
	since, err := time.Parse(time.RFC3339Nano, b.Since)
	if err != nil {
		t.Fatal(err)
	}
	if since.After(time.Now().UTC()) {
		t.Fatal("since in future")
	}
}

func TestResolveHourlyChartBoundsCustom(t *testing.T) {
	b, err := resolveHourlyChartBounds(0, "2026-03-01", "2026-03-03")
	if err != nil {
		t.Fatal(err)
	}
	if !b.UsedCustomRange || b.RangeStart != "2026-03-01" || b.RangeEnd != "2026-03-03" {
		t.Fatalf("unexpected bounds: %+v", b)
	}
	until, err := time.Parse(time.RFC3339Nano, b.Until)
	if err != nil {
		t.Fatal(err)
	}
	wantUntil := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	if !until.Equal(wantUntil) {
		t.Fatalf("until = %v want %v", until, wantUntil)
	}
}

func TestResolveHourlyChartBoundsCustomInvalid(t *testing.T) {
	_, err := resolveHourlyChartBounds(0, "2026-03-05", "2026-03-01")
	if err == nil {
		t.Fatal("expected error")
	}
}
