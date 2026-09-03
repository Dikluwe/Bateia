/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_synthesis_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestImprovementRequestRequiresFrozenSourceAndCounterevidencePolicy(t *testing.T) {
	request := domain.ImprovementSynthesisRequest{Idea: domain.IdeaSpec{Problem: "p", Capabilities: []string{"c"}}, Policy: domain.ImprovementSynthesisPolicy{MaxHypotheses: 1, MaxMechanismsPerHypothesis: 1}}
	if _, err := normalizeImprovementSynthesisRequest(request); err == nil {
		t.Fatal("request without frozen source and counterevidence policy was accepted")
	}
}
