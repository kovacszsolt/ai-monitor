package airequests

import (
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/deps"
)

type Service struct {
	d *deps.Desktop
}

func New(d *deps.Desktop) *Service {
	return &Service{d: d}
}

type ListResult struct {
	Items    []domain.AIRequest
	Total    int64
	Page     int
	PageSize int
}

func (s *Service) List(page, pageSize int, projectID, state, tag, model, sortBy string, sortDesc bool) (ListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	items, total, err := s.d.Requests.List(s.d.Ctx(), repository.AIRequestListFilter{
		ProjectID:      projectID,
		State:          state,
		Tag:            tag,
		ModelRouteID: model,
		SortBy:         sortBy,
		SortDesc:       sortDesc,
	}, pageSize, offset)
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) ListTagFilterOptions() (enabled bool, tags []string, hasUntagged bool, err error) {
	return s.d.Requests.ListTagFilterOptions(s.d.Ctx())
}

func (s *Service) BulkSetRequestTags(ids []string, tag string) (int64, error) {
	return s.d.Requests.BulkSetRequestTags(s.d.Ctx(), ids, tag)
}

func (s *Service) Delete(id string) error {
	return s.d.Requests.Delete(s.d.Ctx(), id)
}

func (s *Service) BulkDelete(ids []string) (int64, error) {
	return s.d.Requests.BulkDelete(s.d.Ctx(), ids)
}

func (s *Service) ListModelRouteFilterOptions() ([]repository.ModelRouteFilterOption, error) {
	return s.d.Requests.ListModelRouteFilterOptions(s.d.Ctx())
}

func (s *Service) GetRequestContent(id string) (domain.AIRequestContent, error) {
	return s.d.Requests.GetRequestContent(s.d.Ctx(), id)
}
