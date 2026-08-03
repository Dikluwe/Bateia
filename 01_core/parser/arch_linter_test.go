/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestEvaluateArchViolations_HighCoupling(t *testing.T) {
	deps := make([]domain.Dependency, 40)
	manifest := &domain.Manifest{
		ModuleName:   "heavy-repo",
		Dependencies: deps,
	}

	violations := EvaluateArchViolations(manifest, nil)
	if len(violations) != 1 {
		t.Fatalf("Esperava 1 violação por acoplamento, obteve %d", len(violations))
	}
	if violations[0].Rule != "HIGH_DEPENDENCY_COUPLING" {
		t.Errorf("Esperava regra HIGH_DEPENDENCY_COUPLING, obteve %s", violations[0].Rule)
	}
}
