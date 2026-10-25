package pricing

import (
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

// ModelRoutePriceIndex resolves route pricing by stable model_route_id.
type ModelRoutePriceIndex struct {
	byID map[string]domain.ModelRoute
}

func NewModelRoutePriceIndex(routes []domain.ModelRoute) ModelRoutePriceIndex {
	idx := ModelRoutePriceIndex{byID: make(map[string]domain.ModelRoute, len(routes))}
	for _, route := range routes {
		id := strings.TrimSpace(route.ID)
		if id == "" {
			continue
		}
		idx.byID[id] = route
	}
	return idx
}

func (idx ModelRoutePriceIndex) LookupByID(id string) (domain.ModelRoute, bool) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ModelRoute{}, false
	}
	route, ok := idx.byID[id]
	return route, ok
}
