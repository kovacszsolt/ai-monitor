package service

import (
	"strings"
	"unicode/utf8"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

const maxChatPinnedTitleLen = 80

type ChatPinnedPromptDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *DesktopService) ListChatPinnedPrompts() Result[[]ChatPinnedPromptDTO] {
	items, err := s.deps.ChatPins.List(s.ctx())
	if err != nil {
		return fail[[]ChatPinnedPromptDTO](err)
	}
	out := make([]ChatPinnedPromptDTO, 0, len(items))
	for _, p := range items {
		out = append(out, mapChatPinnedPrompt(p))
	}
	return ok(out)
}

func (s *DesktopService) CreateChatPinnedPrompt(body, title string) Result[ChatPinnedPromptDTO] {
	title = strings.TrimSpace(title)
	if title == "" {
		title = derivePinnedPromptTitle(body)
	}
	p, err := s.deps.ChatPins.Create(s.ctx(), title, body)
	if err != nil {
		return fail[ChatPinnedPromptDTO](err)
	}
	return ok(mapChatPinnedPrompt(p))
}

func (s *DesktopService) DeleteChatPinnedPrompt(id string) Result[bool] {
	if err := s.deps.ChatPins.Delete(s.ctx(), id); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func mapChatPinnedPrompt(p domain.ChatPinnedPrompt) ChatPinnedPromptDTO {
	return ChatPinnedPromptDTO{
		ID:        p.ID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: formatTime(p.CreatedAt),
		UpdatedAt: formatTime(p.UpdatedAt),
	}
}

func derivePinnedPromptTitle(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return "Pinned question"
	}
	line := body
	if idx := strings.IndexAny(body, "\r\n"); idx >= 0 {
		line = body[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		line = body
	}
	if utf8.RuneCountInString(line) <= maxChatPinnedTitleLen {
		return line
	}
	runes := []rune(line)
	return string(runes[:maxChatPinnedTitleLen]) + "…"
}
