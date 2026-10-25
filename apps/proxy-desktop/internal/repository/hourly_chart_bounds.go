package repository

import (
	"fmt"
	"strings"
	"time"
)

const maxHourlyChartSpanDays = 90

type hourlyChartBounds struct {
	Since           string
	Until           string
	Days            int
	RangeStart      string
	RangeEnd        string
	UsedCustomRange bool
}

func resolveHourlyChartBounds(days int, rangeStart, rangeEnd string) (hourlyChartBounds, error) {
	rs := strings.TrimSpace(rangeStart)
	re := strings.TrimSpace(rangeEnd)
	if rs != "" || re != "" {
		if rs == "" || re == "" {
			return hourlyChartBounds{}, fmt.Errorf("custom range requires both start and end dates (YYYY-MM-DD)")
		}
		startDay, err := time.ParseInLocation("2006-01-02", rs, time.UTC)
		if err != nil {
			return hourlyChartBounds{}, fmt.Errorf("invalid range start date")
		}
		endDay, err := time.ParseInLocation("2006-01-02", re, time.UTC)
		if err != nil {
			return hourlyChartBounds{}, fmt.Errorf("invalid range end date")
		}
		if endDay.Before(startDay) {
			return hourlyChartBounds{}, fmt.Errorf("range end must be on or after start")
		}
		spanDays := int(endDay.Sub(startDay).Hours()/24) + 1
		if spanDays > maxHourlyChartSpanDays {
			return hourlyChartBounds{}, fmt.Errorf("custom range cannot exceed %d days", maxHourlyChartSpanDays)
		}
		until := endDay.AddDate(0, 0, 1)
		return hourlyChartBounds{
			Since:           startDay.Format(time.RFC3339Nano),
			Until:           until.Format(time.RFC3339Nano),
			RangeStart:      rs,
			RangeEnd:        re,
			UsedCustomRange: true,
		}, nil
	}

	if days <= 0 {
		days = 7
	}
	if days > maxHourlyChartSpanDays {
		days = maxHourlyChartSpanDays
	}
	now := time.Now().UTC()
	return hourlyChartBounds{
		Since: now.AddDate(0, 0, -days).Format(time.RFC3339Nano),
		Until: now.Format(time.RFC3339Nano),
		Days:  days,
	}, nil
}
