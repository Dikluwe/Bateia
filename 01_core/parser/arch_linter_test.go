/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter_test.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"strings"
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

func TestEvaluateArchViolations_PolyglotRust(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "mado-linter",
		Language:   "rust",
		Dependencies: []domain.Dependency{
			{Name: "sqlx", Version: "0.7.0"},
		},
	}

	tree := []domain.DirectoryNode{
		{Path: "src/main.rs", Type: "file"},
		{Path: "tests/cli_test.rs", Type: "file"},
	}

	violations := EvaluateArchViolations(manifest, tree)
	if len(violations) != 2 {
		t.Fatalf("Esperava 2 violações em Rust, obteve %d", len(violations))
	}

	for _, v := range violations {
		if v.Rule == "MISSING_INFRASTRUCTURE_ISOLATION" {
			if !strings.Contains(v.Suggestion, "Cargo") && !strings.Contains(v.Suggestion, "trait") && !strings.Contains(v.Suggestion, "src/infra") {
				t.Errorf("Dica do linteiro deveria ser adaptada para Rust, obteve: %s", v.Suggestion)
			}
		}
	}
}

func TestCalculateHealthScore(t *testing.T) {
	violations := []domain.ArchViolation{
		{Severity: "HIGH"},
		{Severity: "MEDIUM"},
	}

	score := CalculateHealthScore(violations)
	if score != 70 {
		t.Fatalf("Esperava HealthScore 70 (100 - 20 - 10), obteve %d", score)
	}
}
