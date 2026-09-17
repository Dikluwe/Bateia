/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0073_implementation_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestStep0073PreservesLegacyCompoundTestSymbols(t *testing.T) {
	rule := findRule("blast radius", DefaultEvidenceRules())
	bundle := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Reports blast radius for a code change."},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "conformance_test.go", Content: "func TestNoSpilloverWhenTheChangeStaysInScope(t *testing.T) {}"},
	)
	got := assessCapability("blast radius", bundle, rule)
	if got.Status != domain.EvidencePartial {
		t.Fatalf("compound test symbol lost: %+v", got)
	}
}
