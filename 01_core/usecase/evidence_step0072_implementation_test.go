/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0072_implementation_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestStep0072TargetedCoveragePreservesPositiveEvidenceDespiteGeneralSampleLimit(t *testing.T) {
	const capability = "contiguous span view"
	rule := domain.CapabilityRule{
		ID: "cpp-span", Basis: domain.EvidenceRuleDeclared,
		CapabilityTerms: []string{capability},
		StrongEvidence:  []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"span"}}},
	}
	ref := domain.CommitRef{Owner: "local", Repo: "libstdcxx", Commit: "immutable-13"}
	content := "template<class T> class span {};"
	base := domain.EvidenceBundle{
		Repository: ref, Complete: false, TreeComplete: true,
		Items: []domain.CollectedEvidence{{Kind: domain.EvidenceFile, Ref: ref, Path: "span", Content: content, ContentHash: contentHash(content)}},
	}

	t.Run("targeted examined keeps positive signal", func(t *testing.T) {
		bundle := base
		bundle.RequirementCoverage = []domain.RequirementCollectionCoverage{{
			RequirementKind: "capability", Requirement: capability, RuleID: rule.ID,
			Status: domain.RequirementTargetedExamined, SearchComplete: true,
			CandidatePaths: []string{"span"}, ExaminedPaths: []string{"span"}, SkippedPaths: []string{},
		}}
		got := assessCapability(capability, bundle, &rule)
		if got.Status != domain.EvidencePartial || got.ReasonCode == "COLLECTION_INCOMPLETE" || len(got.Sources) == 0 || got.Confidence == 0 {
			t.Fatalf("cobertura dirigida completa apagou evidência positiva: %+v", got)
		}
	})

	t.Run("legacy bundle remains globally inconclusive", func(t *testing.T) {
		got := assessCapability(capability, base, &rule)
		if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "COLLECTION_INCOMPLETE" {
			t.Fatalf("bundle legado ignorou incompletude global: %+v", got)
		}
	})

	t.Run("incomplete targeted coverage remains inconclusive", func(t *testing.T) {
		bundle := base
		bundle.RequirementCoverage = []domain.RequirementCollectionCoverage{{
			RequirementKind: "capability", Requirement: capability, RuleID: rule.ID,
			Status: domain.RequirementCollectionIncomplete, SearchComplete: true,
			CandidatePaths: []string{"span", "experimental/span"}, ExaminedPaths: []string{"span"}, SkippedPaths: []string{"experimental/span"},
		}}
		got := assessCapability(capability, bundle, &rule)
		if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "COLLECTION_INCOMPLETE" {
			t.Fatalf("cobertura dirigida incompleta promoveu sinal positivo: %+v", got)
		}
	})
}
