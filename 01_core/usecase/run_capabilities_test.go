/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/run_capabilities_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"testing"
)

type capabilityExtractor struct {
	inputs  []domain.MechanismExtractionInput
	failure error
}

func (e *capabilityExtractor) ExtractMechanisms(_ context.Context, in domain.MechanismExtractionInput) (domain.MechanismExtractionResult, error) {
	e.inputs = append(e.inputs, in)
	if e.failure != nil {
		return domain.MechanismExtractionResult{}, e.failure
	}
	id := "z"
	if len(e.inputs) > 1 {
		id = "a"
	}
	return domain.MechanismExtractionResult{Mechanisms: []domain.ArchitecturalMechanism{{ID: id}}, Gaps: []domain.MechanismGap{{Code: id}}}, nil
}
func TestRunCapabilityExtractionOwnsConfigurationAndOrdering(t *testing.T) {
	catalog, digest := ActiveMechanismRuleCatalog()
	extractor := &capabilityExtractor{}
	baseline := &domain.ConventionalBaseline{BaselineID: "declared"}
	var capability mechanismRunCapability
	if !capability.configure(extractor, catalog, digest, baseline) {
		t.Fatal("valid catalog rejected")
	}
	catalog[0] ^= 1
	if capability.configure(extractor, catalog, digest, nil) {
		t.Fatal("corrupt catalog accepted")
	}
	mechanisms, gaps, err := capability.extract(context.Background(), domain.IdeaRunResult{Bundles: []domain.EvidenceBundle{{}, {}}})
	if err != nil || len(mechanisms) != 2 || mechanisms[0].ID != "a" || len(gaps) != 2 || gaps[0].Code != "a" {
		t.Fatalf("aggregation: %v %v %v", mechanisms, gaps, err)
	}
	if len(extractor.inputs) != 2 || extractor.inputs[0].ConventionalBaseline != baseline || extractor.inputs[0].RuleCatalogJSON[0] == catalog[0] {
		t.Fatal("configuration was not retained")
	}
	extractor.failure = errors.New("extractor failed")
	mechanisms, gaps, err = capability.extract(context.Background(), domain.IdeaRunResult{Bundles: []domain.EvidenceBundle{{}}})
	if !errors.Is(err, extractor.failure) || mechanisms != nil || gaps != nil {
		t.Fatal("partial result escaped failure")
	}
}
func TestSynthesisAdmissionConfigurationIsIndependent(t *testing.T) {
	service := NewRecommendationRunService(nil, nil, nil, nil)
	service.ConfigureImprovementSynthesis("invalid")
	if service.synthesis.enabled() {
		t.Fatal("invalid catalog enabled synthesis")
	}
	service.ConfigureImprovementSynthesis(activeMechanismRulesSHA256)
	if !service.synthesis.enabled() || service.mechanisms.enabled() {
		t.Fatal("synthesis changed extraction configuration")
	}
}

func TestConfiguredMechanismFactoryRejectsBeforeExtraction(t *testing.T) {
	catalog, digest := ActiveMechanismRuleCatalog()
	if extractor, err := NewConfiguredMechanismExtractor(catalog, digest); err != nil || extractor == nil {
		t.Fatal("valid configuration rejected")
	}
	catalog[0] ^= 1
	if extractor, err := NewConfiguredMechanismExtractor(catalog, digest); err != domain.ErrIntegrity || extractor != nil {
		t.Fatal("corrupt configuration accepted")
	}
}
