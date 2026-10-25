package service

import (
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
)

func (s *DesktopService) GetModelTokenUsageOverTime(tag string, days int) Result[ModelTokenTimeSeriesDTO] {
	points, days, err := s.analytics.ModelTokenUsageOverTime(tag, days)
	if err != nil {
		return fail[ModelTokenTimeSeriesDTO](err)
	}
	out := ModelTokenTimeSeriesDTO{
		BucketUnit: "hour",
		Days:       days,
		Points:     make([]ModelTokenTimePointDTO, 0, len(points)),
	}
	for _, p := range points {
		out.Points = append(out.Points, ModelTokenTimePointDTO{
			BucketStart:  p.BucketStart,
			Model:        p.Model,
			InputTokens:  p.InputTokens,
			OutputTokens: p.OutputTokens,
		})
	}
	return ok(out)
}

func (s *DesktopService) GetHourlyRequestCountByModel(tag, projectID string, days int, rangeStart, rangeEnd string) Result[HourlyModelRequestSeriesDTO] {
	points, meta, err := s.analytics.HourlyRequestCountByModel(tag, projectID, days, rangeStart, rangeEnd)
	if err != nil {
		return fail[HourlyModelRequestSeriesDTO](err)
	}
	out := HourlyModelRequestSeriesDTO{
		BucketUnit: "hour",
		Days:       meta.Days,
		RangeStart: meta.RangeStart,
		RangeEnd:   meta.RangeEnd,
		Points:     make([]HourlyModelRequestCountDTO, 0, len(points)),
	}
	for _, hp := range points {
		p := hp.Point
		dto := HourlyModelRequestCountDTO{
			BucketStart:       p.BucketStart,
			Model:             p.Model,
			RequestCount:      p.RequestCount,
			InputTokens:       p.InputTokens,
			CachedInputTokens: p.CachedInputTokens,
			OutputTokens:      p.OutputTokens,
			TotalCostUSD:      hp.TotalCostUSD,
			HasTotalCost:      hp.HasTotalCost,
			GenerationTPS:     hp.GenerationTPS,
			HasGenerationTPS:  hp.HasGenerationTPS,
		}
		out.Points = append(out.Points, dto)
	}
	return ok(out)
}

func (s *DesktopService) GetModelUsageStats(tag, projectID string) Result[[]ModelUsageStatsDTO] {
	items, err := s.analytics.ModelUsageStats(tag, projectID)
	if err != nil {
		return fail[[]ModelUsageStatsDTO](err)
	}
	out := make([]ModelUsageStatsDTO, 0, len(items))
	for _, m := range items {
		out = append(out, mapModelUsageStats(m))
	}
	return ok(out)
}

func (s *DesktopService) GetDashboardStats() Result[DashboardStatsDTO] {
	stats, err := s.analytics.Dashboard()
	if err != nil {
		return fail[DashboardStatsDTO](err)
	}
	return ok(mapDashboard(stats))
}

func (s *DesktopService) ListAIRequestTagOptions() Result[AIRequestTagOptionsDTO] {
	enabled, tags, hasUntagged, err := s.airequests.ListTagFilterOptions()
	if err != nil {
		return fail[AIRequestTagOptionsDTO](err)
	}
	return ok(AIRequestTagOptionsDTO{
		Enabled:     enabled,
		Tags:        tags,
		HasUntagged: hasUntagged,
		NoTagLabel:  repository.NoTagLabel,
	})
}

func (s *DesktopService) BulkSetAIRequestTag(ids []string, tag string) Result[BulkSetAIRequestTagResultDTO] {
	n, err := s.airequests.BulkSetRequestTags(ids, tag)
	if err != nil {
		return fail[BulkSetAIRequestTagResultDTO](err)
	}
	return ok(BulkSetAIRequestTagResultDTO{Updated: n})
}

func (s *DesktopService) DeleteAIRequest(id string) Result[bool] {
	if err := s.airequests.Delete(id); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) BulkDeleteAIRequests(ids []string) Result[BulkDeleteAIRequestResultDTO] {
	n, err := s.airequests.BulkDelete(ids)
	if err != nil {
		return fail[BulkDeleteAIRequestResultDTO](err)
	}
	return ok(BulkDeleteAIRequestResultDTO{Deleted: n})
}

func (s *DesktopService) ListAIRequestModelOptions() Result[[]string] {
	models, err := s.airequests.ListRequestedModelFilterOptions()
	if err != nil {
		return fail[[]string](err)
	}
	return ok(models)
}

func (s *DesktopService) GetAIRequestContent(id string) Result[AIRequestContentDTO] {
	content, err := s.airequests.GetRequestContent(id)
	if err != nil {
		return fail[AIRequestContentDTO](err)
	}
	return ok(mapAIRequestContent(content))
}

func (s *DesktopService) ListAIRequests(page, pageSize int, projectID, state, tag, model, sortBy string, sortDesc bool) Result[AIRequestListDTO] {
	res, err := s.airequests.List(page, pageSize, projectID, state, tag, model, sortBy, sortDesc)
	if err != nil {
		return fail[AIRequestListDTO](err)
	}
	dtos := make([]AIRequestDTO, 0, len(res.Items))
	for _, item := range res.Items {
		dtos = append(dtos, mapAIRequest(item))
	}
	return ok(AIRequestListDTO{
		Items:    dtos,
		Total:    res.Total,
		Page:     res.Page,
		PageSize: res.PageSize,
	})
}
