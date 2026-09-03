/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_run_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestCanonicalRecommendationRequestHashNormalizesSemanticSets(t *testing.T) {
	request := domain.RecommendationRequest{
		Idea:   domain.IdeaSpec{Problem: "inventory", Users: []string{"shops", "homes"}, Capabilities: []string{"forecast", "alerts"}, Constraints: []string{"private"}, NonGoals: []string{"ads"}, PreferredStack: []string{"Go"}},
		Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT", "Apache-2.0"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .3},
		Target: domain.SoftwareAlternativeTarget{Name: "pantry"},
		Search: domain.RecommendationSearch{Strategy: "literal", Locales: []string{"pt-BR", "en-US"}, MaxPerQuery: 2, MaxCandidates: 4},
	}
	reordered := request
	reordered.Idea.Users = []string{"homes", "shops"}
	reordered.Policy.AllowedLicenses = []string{"Apache-2.0", "MIT"}
	reordered.Search.Locales = []string{"en-US", "pt-BR"}

	canonical, hash, err := CanonicalRecommendationRequestHash(request)
	if err != nil {
		t.Fatal(err)
	}
	other, otherHash, err := CanonicalRecommendationRequestHash(reordered)
	if err != nil {
		t.Fatal(err)
	}
	if hash != otherHash || !reflect.DeepEqual(canonical, other) {
		t.Fatalf("semantic set order changed identity: %q != %q", hash, otherHash)
	}
}
