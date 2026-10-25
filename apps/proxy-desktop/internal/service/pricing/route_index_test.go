package pricing

import (
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

func TestModelRoutePriceIndex_LookupByID(t *testing.T) {
	routes := []domain.ModelRoute{{
		ID:              "route-1",
		PublicModelName: "gpt-5.4",
	}}
	idx := NewModelRoutePriceIndex(routes)

	if _, ok := idx.LookupByID("route-1"); !ok {
		t.Fatal("expected route by id")
	}
	if _, ok := idx.LookupByID("gpt-5.4"); ok {
		t.Fatal("must not resolve by public model name")
	}
}
