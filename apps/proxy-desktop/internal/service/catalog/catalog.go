package catalog

import (
	"fmt"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/repository"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/deps"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/pricing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/storeerr"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/storemaps"
)

type Service struct {
	d *deps.Desktop
}

func New(d *deps.Desktop) *Service {
	return &Service{d: d}
}

func (s *Service) ListProjects() ([]domain.Project, error) {
	return repository.ListProjectsWithStats(s.d.Ctx(), s.d.DB)
}

func (s *Service) GetProject(id string) (domain.Project, error) {
	p, err := s.d.Projects.GetByID(s.d.Ctx(), id)
	if err != nil {
		return domain.Project{}, storeerr.Map(err)
	}
	return storemaps.Project(p), nil
}

func (s *Service) CreateProject(name, description string) (domain.Project, error) {
	p, err := s.d.Projects.Create(s.d.Ctx(), name, description)
	if err != nil {
		return domain.Project{}, storeerr.Map(err)
	}
	return storemaps.Project(p), nil
}

func (s *Service) UpdateProject(id, name, description string) (domain.Project, error) {
	p, err := s.d.Projects.Update(s.d.Ctx(), id, name, description)
	if err != nil {
		return domain.Project{}, storeerr.Map(err)
	}
	return storemaps.Project(p), nil
}

func (s *Service) SetProjectEnabled(id string, enabled bool) error {
	return storeerr.Map(s.d.Projects.SetEnabled(s.d.Ctx(), id, enabled))
}

func (s *Service) ListProviders() ([]domain.Provider, error) {
	items, err := s.d.Providers.List(s.d.Ctx())
	if err != nil {
		return nil, err
	}
	out := make([]domain.Provider, 0, len(items))
	for _, p := range items {
		out = append(out, storemaps.Provider(p))
	}
	return out, nil
}

func (s *Service) CreateProvider(name, providerType, baseURL, credential string, connectionTimeoutMS, requestTimeoutMS int) (domain.Provider, error) {
	p, err := s.d.Providers.Add(s.d.Ctx(), name, providers.Type(providerType), baseURL, credential, connectionTimeoutMS, requestTimeoutMS)
	if err != nil {
		return domain.Provider{}, storeerr.Map(err)
	}
	return storemaps.Provider(p), nil
}

func (s *Service) UpdateProvider(id, baseURL, credential string, connectionTimeoutMS, requestTimeoutMS int) (domain.Provider, error) {
	p, err := s.d.Providers.Update(s.d.Ctx(), id, baseURL, credential, connectionTimeoutMS, requestTimeoutMS)
	if err != nil {
		return domain.Provider{}, storeerr.Map(err)
	}
	return storemaps.Provider(p), nil
}

func (s *Service) SetProviderEnabled(id string, enabled bool) error {
	return storeerr.Map(s.d.Providers.SetEnabled(s.d.Ctx(), id, enabled))
}

func (s *Service) ListModelRoutes() ([]domain.ModelRoute, map[string]string, error) {
	routes, err := s.d.Routes.List(s.d.Ctx())
	if err != nil {
		return nil, nil, err
	}
	providers, err := s.d.Providers.List(s.d.Ctx())
	if err != nil {
		return nil, nil, err
	}
	names := make(map[string]string, len(providers))
	for _, p := range providers {
		names[p.ID] = p.Name
	}
	out := make([]domain.ModelRoute, 0, len(routes))
	for _, r := range routes {
		out = append(out, storemaps.ModelRoute(r))
	}
	return out, names, nil
}

func (s *Service) CreateModelRoute(
	publicModelName, providerID, upstreamModelName, protocol, modality string,
	inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M string,
) (domain.ModelRoute, string, error) {
	parsedModality, err := parseModality(modality)
	if err != nil {
		return domain.ModelRoute{}, "", err
	}
	inputUSD, cachedUSD, outputUSD, err := pricing.ParseModelRouteTokenPricing(inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M)
	if err != nil {
		return domain.ModelRoute{}, "", err
	}
	r, err := s.d.Routes.Add(s.d.Ctx(), publicModelName, providerID, upstreamModelName, routing.Protocol(protocol), parsedModality)
	if err != nil {
		return domain.ModelRoute{}, "", storeerr.Map(err)
	}
	if inputUSD != nil || cachedUSD != nil || outputUSD != nil {
		r, err = s.d.Routes.Update(s.d.Ctx(), r.ID, upstreamModelName, routing.Protocol(protocol), parsedModality, inputUSD, cachedUSD, outputUSD)
		if err != nil {
			return domain.ModelRoute{}, "", storeerr.Map(err)
		}
	}
	return storemaps.ModelRoute(r), providerName(s.d, r.ProviderID), nil
}

func (s *Service) UpdateModelRoute(
	id, upstreamModelName, protocol, modality string,
	inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M string,
) (domain.ModelRoute, string, error) {
	parsedModality, err := parseModality(modality)
	if err != nil {
		return domain.ModelRoute{}, "", err
	}
	inputUSD, cachedUSD, outputUSD, err := pricing.ParseModelRouteTokenPricing(inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M)
	if err != nil {
		return domain.ModelRoute{}, "", err
	}
	r, err := s.d.Routes.Update(s.d.Ctx(), id, upstreamModelName, routing.Protocol(protocol), parsedModality, inputUSD, cachedUSD, outputUSD)
	if err != nil {
		return domain.ModelRoute{}, "", storeerr.Map(err)
	}
	return storemaps.ModelRoute(r), providerName(s.d, r.ProviderID), nil
}

func parseModality(v string) (routing.Modality, error) {
	if v == "" {
		return routing.ModalityTextText, nil
	}
	switch routing.Modality(v) {
	case routing.ModalityTextText, routing.ModalityTextImage, routing.ModalityTextImageText:
		return routing.Modality(v), nil
	default:
		return "", fmt.Errorf("invalid modality %q (expected text-text, text-image, or text_image-text)", v)
	}
}

func (s *Service) SetModelRouteEnabled(id string, enabled bool) error {
	return storeerr.Map(s.d.Routes.SetEnabled(s.d.Ctx(), id, enabled))
}

func (s *Service) ListAPIKeys(projectID string) ([]domain.APIKey, error) {
	items, err := s.d.Keys.ListByProject(s.d.Ctx(), projectID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.APIKey, 0, len(items))
	for _, k := range items {
		out = append(out, storemaps.APIKey(k))
	}
	return out, nil
}

func (s *Service) CreateAPIKey(projectID, name string, allowAllModels bool) (domain.APIKey, string, error) {
	k, plaintext, err := s.d.Keys.Create(s.d.Ctx(), projectID, name, allowAllModels, nil)
	if err != nil {
		return domain.APIKey{}, "", err
	}
	return storemaps.APIKey(k), plaintext, nil
}

func (s *Service) RevokeAPIKey(id string) error {
	return storeerr.Map(s.d.Keys.Revoke(s.d.Ctx(), id))
}

func (s *Service) DeleteAPIKey(id string) error {
	return storeerr.Map(s.d.Keys.Delete(s.d.Ctx(), id))
}

func (s *Service) SetAPIKeyAllowAllModels(id string, allowAllModels bool) error {
	return repository.SetAPIKeyAllowAllModels(s.d.Ctx(), s.d.DB, id, allowAllModels)
}

func providerName(d *deps.Desktop, providerID string) string {
	if p, err := d.Providers.Get(d.Ctx(), providerID); err == nil {
		return p.Name
	}
	return ""
}
