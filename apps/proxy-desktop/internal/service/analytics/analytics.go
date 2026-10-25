package analytics

import (
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/deps"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/pricing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/storemaps"
)

type Service struct {
	d *deps.Desktop
}

func New(d *deps.Desktop) *Service {
	return &Service{d: d}
}

func (s *Service) ModelTokenUsageOverTime(tag string, days int) ([]domain.ModelTokenTimePoint, int, error) {
	points, err := s.d.Requests.ModelTokenUsageOverTime(s.d.Ctx(), tag, days)
	if err != nil {
		return nil, 0, err
	}
	if days <= 0 {
		days = 14
	}
	return points, days, nil
}

type HourlyModelPoint struct {
	Point            domain.HourlyModelRequestCount
	TotalCostUSD     float64
	HasTotalCost     bool
	GenerationTPS    float64
	HasGenerationTPS bool
}

type HourlySeriesMeta struct {
	Days       int
	RangeStart string
	RangeEnd   string
}

func (s *Service) HourlyRequestCountByModel(tag, projectID string, days int, rangeStart, rangeEnd string) ([]HourlyModelPoint, HourlySeriesMeta, error) {
	points, err := s.d.Requests.HourlyRequestCountByModel(s.d.Ctx(), tag, projectID, days, rangeStart, rangeEnd)
	if err != nil {
		return nil, HourlySeriesMeta{}, err
	}
	routes, err := s.d.Routes.List(s.d.Ctx())
	if err != nil {
		return nil, HourlySeriesMeta{}, err
	}
	priceByPublic := make(map[string]domain.ModelRoute, len(routes))
	for _, route := range routes {
		priceByPublic[route.PublicModelName] = storemaps.ModelRoute(route)
	}
	meta := HourlySeriesMeta{
		Days:       days,
		RangeStart: strings.TrimSpace(rangeStart),
		RangeEnd:   strings.TrimSpace(rangeEnd),
	}
	if meta.RangeStart != "" && meta.RangeEnd != "" {
		meta.Days = 0
	} else if meta.Days <= 0 {
		meta.Days = 7
	}
	out := make([]HourlyModelPoint, 0, len(points))
	for _, p := range points {
		hp := HourlyModelPoint{Point: p}
		if route, ok := priceByPublic[p.Model]; ok {
			cost, has := pricing.ModelUsageTotalCostUSD(route, p.InputTokens, p.CachedInputTokens, p.OutputTokens)
			hp.TotalCostUSD = cost
			hp.HasTotalCost = has
		}
		if p.DecodeDurationMS > 0 && p.OutputTokens > 0 {
			hp.GenerationTPS = float64(p.OutputTokens) / (float64(p.DecodeDurationMS) / 1000.0)
			hp.HasGenerationTPS = true
		}
		out = append(out, hp)
	}
	return out, meta, nil
}

func (s *Service) ModelUsageStats(tag, projectID string) ([]domain.ModelUsageStats, error) {
	items, err := s.d.Requests.ModelUsageStats(s.d.Ctx(), tag, projectID)
	if err != nil {
		return nil, err
	}
	routes, err := s.d.Routes.List(s.d.Ctx())
	if err != nil {
		return nil, err
	}
	priceByPublic := make(map[string]domain.ModelRoute, len(routes))
	for _, route := range routes {
		priceByPublic[route.PublicModelName] = storemaps.ModelRoute(route)
	}
	for i := range items {
		if route, ok := priceByPublic[items[i].Name]; ok {
			cost, has := pricing.ModelUsageTotalCostUSD(route, items[i].InputTokens, items[i].CachedInputTokens, items[i].OutputTokens)
			items[i].TotalCostUSD = cost
			items[i].HasTotalCost = has
		}
	}
	return items, nil
}

func (s *Service) Dashboard() (domain.DashboardStats, error) {
	return s.d.Requests.Dashboard(s.d.Ctx())
}
