package service

import (
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

type ClientErrorLogEntryDTO struct {
	ID        string `json:"id"`
	Level     string `json:"level"`
	Source    string `json:"source"`
	Message   string `json:"message"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

type ClientErrorLogListDTO struct {
	Items    []ClientErrorLogEntryDTO `json:"items"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}

func (s *DesktopService) LogClientError(level, source, message, detail string) Result[ClientErrorLogEntryDTO] {
	entry, err := s.deps.ClientErrors.Insert(s.ctx(), level, source, message, detail)
	if err != nil {
		return fail[ClientErrorLogEntryDTO](err)
	}
	return ok(mapClientErrorLogEntry(entry))
}

func (s *DesktopService) ListClientErrors(page, pageSize int) Result[ClientErrorLogListDTO] {
	res, err := s.deps.ClientErrors.List(s.ctx(), page, pageSize)
	if err != nil {
		return fail[ClientErrorLogListDTO](err)
	}
	items := make([]ClientErrorLogEntryDTO, 0, len(res.Items))
	for _, e := range res.Items {
		items = append(items, mapClientErrorLogEntry(e))
	}
	return ok(ClientErrorLogListDTO{
		Items:    items,
		Total:    res.Total,
		Page:     res.Page,
		PageSize: res.PageSize,
	})
}

func (s *DesktopService) ClearClientErrors() Result[bool] {
	if err := s.deps.ClientErrors.Clear(s.ctx()); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func mapClientErrorLogEntry(e domain.ClientErrorLogEntry) ClientErrorLogEntryDTO {
	return ClientErrorLogEntryDTO{
		ID:        e.ID,
		Level:     e.Level,
		Source:    e.Source,
		Message:   e.Message,
		Detail:    e.Detail,
		CreatedAt: formatTime(e.CreatedAt),
	}
}
