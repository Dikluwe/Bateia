/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/regional_discovery_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

import (
	"encoding/json"
	"testing"
)

func TestM21RegionalDiscoveryOwnershipContract(t *testing.T) {
	provenance := RegionalQueryProvenance{
		Locale:             "ja-JP",
		CatalogVersion:     "regional-recipes-v1",
		Concept:            "recipe-manager",
		RuleID:             "recipe-manager",
		SourceQuery:        "gerenciador de receitas",
		DerivationStrategy: "controlled-catalog-exact-v1",
		ProposalID:         "regional-ja-JP-recipe-manager",
	}
	encoded, err := json.Marshal(SearchExpansion{ID: provenance.ProposalID, Query: "レシピ管理", Regional: &provenance})
	if err != nil {
		t.Fatal(err)
	}
	var root map[string]any
	if err := json.Unmarshal(encoded, &root); err != nil {
		t.Fatal(err)
	}
	if _, exists := root["search_locale"]; exists {
		t.Fatalf("SearchExpansion exposed a non-contract top-level locale: %s", encoded)
	}
	regional, ok := root["regional"].(map[string]any)
	if !ok || regional["locale"] != "ja-JP" || regional["proposal_id"] != provenance.ProposalID || regional["source_query"] != provenance.SourceQuery {
		t.Fatalf("SearchExpansion.Regional lost sealed provenance: %s", encoded)
	}
}
