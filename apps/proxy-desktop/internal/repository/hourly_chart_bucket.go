package repository

import (
	"fmt"
	"strings"
	"time"
)

// ChartBucketUnit is the time granularity for usage-over-time chart aggregates.
type ChartBucketUnit string

const (
	ChartBucketMinute  ChartBucketUnit = "minute"
	ChartBucketQuarter ChartBucketUnit = "quarter"
	ChartBucketHour    ChartBucketUnit = "hour"
	ChartBucketDay     ChartBucketUnit = "day"
)

func NormalizeChartBucketUnit(raw string) ChartBucketUnit {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "minute", "min":
		return ChartBucketMinute
	case "quarter", "15m", "15min":
		return ChartBucketQuarter
	case "day":
		return ChartBucketDay
	case "hour", "":
		return ChartBucketHour
	default:
		return ChartBucketHour
	}
}

func chartBucketStartSQL(unit ChartBucketUnit) string {
	switch unit {
	case ChartBucketDay:
		return `substr(r.started_at, 1, 10) || 'T00:00:00Z'`
	case ChartBucketMinute:
		return `substr(r.started_at, 1, 16) || ':00Z'`
	case ChartBucketQuarter:
		return `substr(r.started_at, 1, 14) || printf('%02d', (CAST(substr(r.started_at, 15, 2) AS INTEGER) / 15) * 15) || ':00Z'`
	default:
		return `substr(r.started_at, 1, 13) || ':00:00Z'`
	}
}

func chartBucketSpanDays(bounds hourlyChartBounds) int {
	if bounds.UsedCustomRange {
		startDay, err := parseUTCDate(bounds.RangeStart)
		if err != nil {
			return maxHourlyChartSpanDays
		}
		endDay, err := parseUTCDate(bounds.RangeEnd)
		if err != nil {
			return maxHourlyChartSpanDays
		}
		return int(endDay.Sub(startDay).Hours()/24) + 1
	}
	if bounds.Days <= 0 {
		return 7
	}
	return bounds.Days
}

func maxSpanDaysForChartBucket(unit ChartBucketUnit) int {
	switch unit {
	case ChartBucketMinute:
		return 2
	case ChartBucketQuarter:
		return 14
	case ChartBucketDay:
		return maxHourlyChartSpanDays
	default:
		return maxHourlyChartSpanDays
	}
}

func validateChartBucketForBounds(unit ChartBucketUnit, bounds hourlyChartBounds) error {
	span := chartBucketSpanDays(bounds)
	maxSpan := maxSpanDaysForChartBucket(unit)
	if span > maxSpan {
		return fmt.Errorf(
			"%s buckets support at most %d day(s) in the selected range (%d days); use a shorter range or coarser buckets",
			unit,
			maxSpan,
			span,
		)
	}
	return nil
}

func parseUTCDate(ymd string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", strings.TrimSpace(ymd), time.UTC)
}
