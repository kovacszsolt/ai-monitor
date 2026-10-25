package service

import (
	"errors"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
)

var errInvalidExportFormat = errors.New("export format must be jsonl or csv")

var (
	errProjectNotFound    = repository.ErrProjectNotFound
	errProjectDuplicate   = repository.ErrProjectDuplicate
	errProviderNotFound   = repository.ErrProviderNotFound
	errModelRouteNotFound = repository.ErrModelRouteNotFound
	errAPIKeyNotFound     = repository.ErrAPIKeyNotFound
	errAIRequestNotFound         = repository.ErrAIRequestNotFound
	errRequestTagColumnMissing   = repository.ErrRequestTagColumnMissing
	errNoRequestIDs              = repository.ErrNoRequestIDs
	errChatPinnedPromptNotFound  = repository.ErrChatPinnedPromptNotFound
)

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ToErrorDTO(err error) ErrorDTO {
	if err == nil {
		return ErrorDTO{Code: "unknown", Message: "unknown error"}
	}
	msg := err.Error()
	switch {
	case errors.Is(err, errProjectNotFound), errors.Is(err, errProviderNotFound),
		errors.Is(err, errModelRouteNotFound), errors.Is(err, errAPIKeyNotFound),
		errors.Is(err, errAIRequestNotFound), errors.Is(err, errChatPinnedPromptNotFound):
		return ErrorDTO{Code: "not_found", Message: msg}
	case errors.Is(err, errRequestTagColumnMissing):
		return ErrorDTO{Code: "schema", Message: msg}
	case errors.Is(err, errNoRequestIDs):
		return ErrorDTO{Code: "invalid", Message: msg}
	case errors.Is(err, errProjectDuplicate):
		return ErrorDTO{Code: "duplicate", Message: msg}
	default:
		lower := strings.ToLower(msg)
		if strings.Contains(lower, "already exists") || strings.Contains(lower, "unique") {
			return ErrorDTO{Code: "duplicate", Message: msg}
		}
		return ErrorDTO{Code: "internal", Message: msg}
	}
}

type Result[T any] struct {
	OK    bool     `json:"ok"`
	Data  T        `json:"data"`
	Error ErrorDTO `json:"error,omitempty"`
}

func ok[T any](data T) Result[T] {
	return Result[T]{OK: true, Data: data}
}

func fail[T any](err error) Result[T] {
	return Result[T]{OK: false, Error: ToErrorDTO(err)}
}
