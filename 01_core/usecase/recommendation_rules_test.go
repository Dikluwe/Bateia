/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_rules_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestRecommendationServiceWithEvidenceRules(t *testing.T) {
	rules := DefaultEvidenceRules()
	svc, err := NewRecommendationServiceWithEvidenceRules(nil, nil, nil, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}

	invalid := domain.CapabilityRuleCatalog{
		SchemaVersion: "invalid",
	}
	_, err = NewRecommendationServiceWithEvidenceRules(nil, nil, nil, invalid)
	if err == nil {
		t.Fatalf("expected error for invalid catalog")
	}
}
