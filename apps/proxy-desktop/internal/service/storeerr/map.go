package storeerr

import (
	"errors"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/auth"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/projects"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
)

// Map translates proxyengine store errors into repository sentinel errors for API responses.
func Map(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, projects.ErrNotFound):
		return repository.ErrProjectNotFound
	case errors.Is(err, projects.ErrDuplicate):
		return repository.ErrProjectDuplicate
	case errors.Is(err, providers.ErrNotFound):
		return repository.ErrProviderNotFound
	case errors.Is(err, routing.ErrNotFound):
		return repository.ErrModelRouteNotFound
	case errors.Is(err, auth.ErrKeyNotFound):
		return repository.ErrAPIKeyNotFound
	default:
		return err
	}
}
