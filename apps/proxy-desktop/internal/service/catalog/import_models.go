package catalog

import (
	"fmt"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/proxy"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/service/storeerr"
)

// UpstreamModelItem is one upstream model id with import status for a provider.
type UpstreamModelItem struct {
	ID              string
	AlreadyImported bool
}

// ImportModelRoutesResult summarizes a bulk import operation.
type ImportModelRoutesResult struct {
	Created int
	Skipped []string
}

func (s *Service) ListProviderUpstreamModels(providerID string) ([]UpstreamModelItem, error) {
	p, err := s.d.Providers.Get(s.d.Ctx(), providerID)
	if err != nil {
		return nil, storeerr.Map(err)
	}

	upstream, err := proxy.ListUpstreamModels(s.d.Ctx(), p)
	if err != nil {
		return nil, err
	}

	routes, err := s.d.Routes.List(s.d.Ctx())
	if err != nil {
		return nil, err
	}
	importedUpstream := importedUpstreamSet(routes, providerID)

	out := make([]UpstreamModelItem, 0, len(upstream))
	for _, m := range upstream {
		_, already := importedUpstream[m.ID]
		out = append(out, UpstreamModelItem{
			ID:              m.ID,
			AlreadyImported: already,
		})
	}
	return out, nil
}

func (s *Service) ImportModelRoutes(providerID string, upstreamIDs []string, publicNamePrefix string) (ImportModelRoutesResult, error) {
	p, err := s.d.Providers.Get(s.d.Ctx(), providerID)
	if err != nil {
		return ImportModelRoutesResult{}, storeerr.Map(err)
	}

	routes, err := s.d.Routes.List(s.d.Ctx())
	if err != nil {
		return ImportModelRoutesResult{}, err
	}

	importedUpstream := importedUpstreamSet(routes, providerID)
	publicNames := make(map[string]struct{}, len(routes))
	for _, r := range routes {
		publicNames[r.PublicModelName] = struct{}{}
	}

	protocol := defaultImportProtocol(p.ProviderType)
	modality := routing.ModalityTextText

	var result ImportModelRoutesResult
	for _, rawID := range upstreamIDs {
		upstreamID := strings.TrimSpace(rawID)
		if upstreamID == "" {
			continue
		}
		if _, ok := importedUpstream[upstreamID]; ok {
			result.Skipped = append(result.Skipped, skipReason(upstreamID, "already imported for this provider"))
			continue
		}
		publicName := publicNamePrefix + upstreamID
		if strings.TrimSpace(publicName) == "" {
			result.Skipped = append(result.Skipped, skipReason(upstreamID, "public model name would be empty"))
			continue
		}
		if _, ok := publicNames[publicName]; ok {
			result.Skipped = append(result.Skipped, skipReason(upstreamID, "public name already in use"))
			continue
		}
		if _, err := s.d.Routes.Add(s.d.Ctx(), publicName, providerID, upstreamID, protocol, modality); err != nil {
			result.Skipped = append(result.Skipped, skipReason(upstreamID, storeerr.Map(err).Error()))
			continue
		}
		result.Created++
		publicNames[publicName] = struct{}{}
		importedUpstream[upstreamID] = struct{}{}
	}
	return result, nil
}

func importedUpstreamSet(routes []routing.Route, providerID string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, r := range routes {
		if r.ProviderID != providerID {
			continue
		}
		out[r.UpstreamModelName] = struct{}{}
	}
	return out
}

func defaultImportProtocol(ptype providers.Type) routing.Protocol {
	if ptype == providers.TypeAnthropic {
		return routing.ProtocolAnthropicMessages
	}
	return routing.ProtocolOpenAIChat
}

func skipReason(upstreamID, reason string) string {
	return fmt.Sprintf("%s: %s", upstreamID, reason)
}
