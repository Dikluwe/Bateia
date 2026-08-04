/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-08-04
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

func TestEvaluateArchViolations_MissingTestsAndDevOps(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "sample-service",
	}

	tree := []domain.DirectoryNode{
		{Path: "main.go", Type: "file"},
		{Path: "service.go", Type: "file"},
	}

	violations := EvaluateArchViolations(manifest, tree)
	if len(violations) != 2 {
		t.Fatalf("Esperava 2 violações (Missing Tests & DevOps), obteve %d", len(violations))
	}

	ruleMap := make(map[string]bool)
	for _, v := range violations {
		ruleMap[v.Rule] = true
	}

	if !ruleMap["MISSING_TEST_SUITE"] {
		t.Errorf("Esperava violação MISSING_TEST_SUITE")
	}
	if !ruleMap["MISSING_DEVOPS_AUTOMATION"] {
		t.Errorf("Esperava violação MISSING_DEVOPS_AUTOMATION")
	}
}
