/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/market_discovery_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarketDiscoveryProductKeepsUnknownDimensionsAsEmptyArrays(t *testing.T) {
	product := ProductCandidate{
		Identity: ProductIdentity{ConnectorID: "apple-itunes-search-software", ExternalID: "1"},
		Name:     "Example", ProductKind: "application", ListingLocales: []string{}, ObservedMarkets: []string{"JP"},
		PublisherCountries: []string{}, RepositoryLinks: []ProductRepositoryLink{}, Origins: []MarketCandidateOrigin{}, Witnesses: []IdentityWitness{},
	}
	encoded, err := json.Marshal(product)
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{`"listing_locales":[]`, `"publisher_countries":[]`, `"observed_markets":["JP"]`} {
		if !strings.Contains(string(encoded), fragment) {
			t.Fatalf("missing %s in %s", fragment, encoded)
		}
	}
}
