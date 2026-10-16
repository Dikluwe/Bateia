/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_analysis_test.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func analysisFixture() (domain.ImprovementSynthesisRequest, improvementAnalysisBasis) {
	a, b := "mech_"+strings.Repeat("a", 64), "mech_"+strings.Repeat("b", 64)
	request := domain.ImprovementSynthesisRequest{
		SourceArtifacts: []domain.SourceArtifactRef{{RunID: domain.RunID("run_" + strings.Repeat("1", 32)), RunKind: domain.RunKindRecommendation, Artifact: domain.ArtifactMechanismCatalog, SHA256: strings.Repeat("c", 64)}},
		Policy:          domain.ImprovementSynthesisPolicy{MaxHypotheses: 2, MaxMechanismsPerHypothesis: 2, RequireCounterEvidenceSearch: true},
		Proposals:       []domain.ImprovementProposal{{ProposalID: "proposal-1", Title: "Separate analysis from execution", TargetProperty: "recoverability", ExpectedEffect: "preserve conclusions across retries", MechanismRefs: []string{a, b}, Assumptions: []string{"explicit boundaries"}, Derivation: domain.DerivationExternal}},
	}
	basis := improvementAnalysisBasis{mechanisms: map[string]domain.ArchitecturalMechanism{
		a: {ID: a, Name: "State isolation", Purpose: "recoverability", Strategy: "explicit states", Status: domain.MechanismProven, Novelty: domain.NoveltyConventional},
		b: {ID: b, Name: "Evidence boundary", Purpose: "traceability", Strategy: "immutable evidence", Status: domain.MechanismPartial, Novelty: domain.NoveltyUnknown},
	}}
	return request, basis
}

func TestImprovementAnalysisPreservesLegacyProjection(t *testing.T) {
	request, basis := analysisFixture()
	before, _ := json.Marshal(request)
	result := synthesizeImprovement(request, basis)
	c, h, g, p, v := result.comparisons, result.hypotheses, result.gaps, result.plans, result.provenance
	projection := struct {
		Comparisons []domain.ArchitectureComparison
		Hypotheses  []domain.ImprovementHypothesis
		Gaps        []domain.HypothesisGap
		Plans       []domain.ValidationPlanEntry
		Provenance  []domain.SynthesisProvenance
	}{c, h, g, p, v}
	encoded, err := json.Marshal(projection)
	if err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(encoded))
	const want = "5b2134a4a8c8017031a38d8fbe662a26658ee1e50aac4197c0df841fbe729e24"
	if digest != want {
		t.Fatalf("legacy projection digest: %s", digest)
	}
	after, _ := json.Marshal(request)
	if string(before) != string(after) {
		t.Fatal("analysis mutated its request")
	}
	if len(h) != 1 || h[0].Status != domain.HypothesisProposed || h[0].AdoptionDisposition != "NOT_AN_ADOPTION_DECISION" || h[0].TransferMode != domain.TransferReferenceOnly {
		t.Fatalf("unsupported promotion: %+v", h)
	}
	if len(h[0].IntegrationUnknowns) != 5 || len(p) != 1 {
		t.Fatal("composition unknowns or validation plan lost")
	}
}

func TestImprovementAnalysisDoesNotPromoteDegradedEvidence(t *testing.T) {
	request, basis := analysisFixture()
	basis.degraded = true
	result := synthesizeImprovement(request, basis)
	hypotheses, gaps := result.hypotheses, result.gaps
	if len(hypotheses) != 1 || hypotheses[0].Status != domain.HypothesisInconclusive {
		t.Fatalf("degraded evidence promoted: %+v", hypotheses)
	}
	found := false
	for _, gap := range gaps {
		if gap.Code == "SOURCE_RUN_NOT_SUCCEEDED" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing source degradation gap")
	}
}

func TestImprovementAnalysisRuleDerivationIgnoresMapInsertionOrder(t *testing.T) {
	request, basis := analysisFixture()
	request.Proposals = nil
	other := improvementAnalysisBasis{mechanisms: map[string]domain.ArchitecturalMechanism{}}
	keys := []string{"mech_" + strings.Repeat("b", 64), "mech_" + strings.Repeat("a", 64)}
	for _, key := range keys {
		other.mechanisms[key] = basis.mechanisms[key]
	}
	a := synthesizeImprovement(request, basis)
	b := synthesizeImprovement(request, other)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("map order changed analysis")
	}
}

func TestImprovementAnalysisResultStatusDoesNotRepresentAdoption(t *testing.T) {
	request, basis := analysisFixture()
	result := synthesizeImprovement(request, basis)
	if result.status != domain.ImprovementPartial {
		t.Fatalf("unrecognized license and incomplete mechanism must remain partial: %s", result.status)
	}
	request.Proposals = nil
	result = synthesizeImprovement(request, improvementAnalysisBasis{})
	if result.status != domain.ImprovementInconclusive || len(result.hypotheses) != 0 {
		t.Fatal("no hypotheses must remain inconclusive")
	}
}
