package storemaps

import (
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/auth"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/projects"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
)

func Project(p projects.Project) domain.Project {
	return domain.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Enabled:     p.Enabled,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func Provider(p providers.Provider) domain.Provider {
	return domain.Provider{
		ID:                  p.ID,
		Name:                p.Name,
		ProviderType:        domain.ProviderType(p.ProviderType),
		BaseURL:             p.BaseURL,
		Credential:          p.Credential,
		Enabled:             p.Enabled,
		ConnectionTimeoutMS: p.ConnectionTimeoutMS,
		RequestTimeoutMS:    p.RequestTimeoutMS,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

func ModelRoutes(routes []routing.Route) []domain.ModelRoute {
	out := make([]domain.ModelRoute, 0, len(routes))
	for _, r := range routes {
		out = append(out, ModelRoute(r))
	}
	return out
}

func ModelRoute(r routing.Route) domain.ModelRoute {
	return domain.ModelRoute{
		ID:                       r.ID,
		PublicModelName:          r.PublicModelName,
		ProviderID:               r.ProviderID,
		UpstreamModelName:        r.UpstreamModelName,
		Protocol:                 domain.Protocol(r.Protocol),
		Modality:                 domain.Modality(r.Modality),
		Enabled:                  r.Enabled,
		InputTokenUSDPer1M:       r.InputTokenUSDPer1M,
		CachedInputTokenUSDPer1M: r.CachedInputTokenUSDPer1M,
		OutputTokenUSDPer1M:      r.OutputTokenUSDPer1M,
		CreatedAt:                r.CreatedAt,
		UpdatedAt:                r.UpdatedAt,
	}
}

func APIKey(k auth.APIKey) domain.APIKey {
	return domain.APIKey{
		ID:             k.ID,
		ProjectID:      k.ProjectID,
		Name:           k.Name,
		Secret:         k.Secret,
		AllowAllModels: k.AllowAllModels,
		Enabled:        k.Enabled,
		CreatedAt:      k.CreatedAt,
		UpdatedAt:      k.UpdatedAt,
		ExpiresAt:      k.ExpiresAt,
		LastUsedAt:     k.LastUsedAt,
	}
}
