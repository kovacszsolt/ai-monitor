package service

func (s *DesktopService) ListProjects() Result[[]ProjectDTO] {
	items, err := s.catalog.ListProjects()
	if err != nil {
		return fail[[]ProjectDTO](err)
	}
	out := make([]ProjectDTO, 0, len(items))
	for _, p := range items {
		out = append(out, mapProject(p))
	}
	return ok(out)
}

func (s *DesktopService) GetProject(id string) Result[ProjectDTO] {
	p, err := s.catalog.GetProject(id)
	if err != nil {
		return fail[ProjectDTO](err)
	}
	return ok(mapProject(p))
}

func (s *DesktopService) CreateProject(name, description string) Result[ProjectDTO] {
	p, err := s.catalog.CreateProject(name, description)
	if err != nil {
		return fail[ProjectDTO](err)
	}
	return ok(mapProject(p))
}

func (s *DesktopService) UpdateProject(id, name, description string) Result[ProjectDTO] {
	p, err := s.catalog.UpdateProject(id, name, description)
	if err != nil {
		return fail[ProjectDTO](err)
	}
	return ok(mapProject(p))
}

func (s *DesktopService) SetProjectEnabled(id string, enabled bool) Result[bool] {
	if err := s.catalog.SetProjectEnabled(id, enabled); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) ListProviders() Result[[]ProviderDTO] {
	items, err := s.catalog.ListProviders()
	if err != nil {
		return fail[[]ProviderDTO](err)
	}
	out := make([]ProviderDTO, 0, len(items))
	for _, p := range items {
		out = append(out, mapProvider(p))
	}
	return ok(out)
}

func (s *DesktopService) CreateProvider(name, providerType, baseURL, credential string, connectionTimeoutMS, requestTimeoutMS int) Result[ProviderDTO] {
	p, err := s.catalog.CreateProvider(name, providerType, baseURL, credential, connectionTimeoutMS, requestTimeoutMS)
	if err != nil {
		return fail[ProviderDTO](err)
	}
	return ok(mapProvider(p))
}

func (s *DesktopService) UpdateProvider(id, baseURL, credential string, connectionTimeoutMS, requestTimeoutMS int) Result[ProviderDTO] {
	p, err := s.catalog.UpdateProvider(id, baseURL, credential, connectionTimeoutMS, requestTimeoutMS)
	if err != nil {
		return fail[ProviderDTO](err)
	}
	return ok(mapProvider(p))
}

func (s *DesktopService) SetProviderEnabled(id string, enabled bool) Result[bool] {
	if err := s.catalog.SetProviderEnabled(id, enabled); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) ListModelRoutes() Result[[]ModelRouteDTO] {
	routes, names, err := s.catalog.ListModelRoutes()
	if err != nil {
		return fail[[]ModelRouteDTO](err)
	}
	out := make([]ModelRouteDTO, 0, len(routes))
	for _, r := range routes {
		out = append(out, mapModelRoute(r, names[r.ProviderID]))
	}
	return ok(out)
}

func (s *DesktopService) CreateModelRoute(
	publicModelName, providerID, upstreamModelName, protocol, modality string,
	inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M string,
) Result[ModelRouteDTO] {
	r, pname, err := s.catalog.CreateModelRoute(publicModelName, providerID, upstreamModelName, protocol, modality, inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M)
	if err != nil {
		return fail[ModelRouteDTO](err)
	}
	return ok(mapModelRoute(r, pname))
}

func (s *DesktopService) UpdateModelRoute(
	id, upstreamModelName, protocol, modality string,
	inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M string,
) Result[ModelRouteDTO] {
	r, pname, err := s.catalog.UpdateModelRoute(id, upstreamModelName, protocol, modality, inputTokenUSDPer1M, cachedInputTokenUSDPer1M, outputTokenUSDPer1M)
	if err != nil {
		return fail[ModelRouteDTO](err)
	}
	return ok(mapModelRoute(r, pname))
}

func (s *DesktopService) SetModelRouteEnabled(id string, enabled bool) Result[bool] {
	if err := s.catalog.SetModelRouteEnabled(id, enabled); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) ListProviderUpstreamModels(providerID string) Result[[]UpstreamModelDTO] {
	items, err := s.catalog.ListProviderUpstreamModels(providerID)
	if err != nil {
		return fail[[]UpstreamModelDTO](err)
	}
	out := make([]UpstreamModelDTO, 0, len(items))
	for _, m := range items {
		out = append(out, UpstreamModelDTO{
			ID:              m.ID,
			AlreadyImported: m.AlreadyImported,
		})
	}
	return ok(out)
}

func (s *DesktopService) ImportModelRoutes(providerID string, upstreamIDs []string, publicNamePrefix string) Result[ImportModelRoutesResultDTO] {
	res, err := s.catalog.ImportModelRoutes(providerID, upstreamIDs, publicNamePrefix)
	if err != nil {
		return fail[ImportModelRoutesResultDTO](err)
	}
	return ok(ImportModelRoutesResultDTO{
		Created: res.Created,
		Skipped: res.Skipped,
	})
}

func (s *DesktopService) ListAPIKeys(projectID string) Result[[]APIKeyDTO] {
	items, err := s.catalog.ListAPIKeys(projectID)
	if err != nil {
		return fail[[]APIKeyDTO](err)
	}
	out := make([]APIKeyDTO, 0, len(items))
	for _, k := range items {
		out = append(out, mapAPIKey(k))
	}
	return ok(out)
}

func (s *DesktopService) CreateAPIKey(projectID, name string, allowAllModels bool) Result[CreateAPIKeyResultDTO] {
	k, plaintext, err := s.catalog.CreateAPIKey(projectID, name, allowAllModels)
	if err != nil {
		return fail[CreateAPIKeyResultDTO](err)
	}
	keyDTO := mapAPIKey(k)
	if keyDTO.APIKey == "" {
		keyDTO.APIKey = plaintext
	}
	return ok(CreateAPIKeyResultDTO{Key: keyDTO, Plaintext: plaintext})
}

func (s *DesktopService) RevokeAPIKey(id string) Result[bool] {
	if err := s.catalog.RevokeAPIKey(id); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) DeleteAPIKey(id string) Result[bool] {
	if err := s.catalog.DeleteAPIKey(id); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) SetAPIKeyAllowAllModels(id string, allowAllModels bool) Result[bool] {
	if err := s.catalog.SetAPIKeyAllowAllModels(id, allowAllModels); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}
