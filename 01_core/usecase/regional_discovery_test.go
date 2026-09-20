/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/regional_discovery_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"bateia/01_core/domain"
	"context"
	"reflect"
	"testing"
)

func TestM21RegionalDiscoveryOwnershipContract(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"ja-JP", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	provider := &m21RegionalProvider{}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "organizar receitas culinárias"}, Expansions: regional.Proposals}
	snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 7, 1, provider)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(snapshot.SearchLocales, []string{"pt-BR", "ja-JP"}) {
		t.Fatalf("DiscoverySnapshot lost canonical regional ownership: %v", snapshot.SearchLocales)
	}
}
