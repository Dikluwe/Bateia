/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/structure.md
 * @layer L1
 * @updated 2026-07-25
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestEvaluateStructure(t *testing.T) {
	nodes := []domain.DirectoryNode{
		{Path: "cmd", Type: "dir"},
		{Path: "cmd/app/main.go", Type: "file"},
		{Path: "internal", Type: "dir"},
		{Path: "internal/domain/user.go", Type: "file"},
	}

	// 1. Testar aprovação quando todos os mustHave estão presentes
	score := EvaluateStructure(nodes, []string{"cmd", "internal"})
	if !score.Approved {
		t.Errorf("Deveria ter aprovado a estrutura, mas reprovou com motivo: %s", score.Reason)
	}

	// 2. Testar reprovação quando um item de mustHave está ausente
	scoreMissing := EvaluateStructure(nodes, []string{"cmd", "pkg"})
	if scoreMissing.Approved {
		t.Errorf("Deveria ter reprovado pela ausência da pasta 'pkg'")
	}
}

func TestMatchGlobRecursive(t *testing.T) {
	tests := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"src/**/*.rs", "src/parser/ast/rust.rs", true},
		{"src/**/*.rs", "src/main.rs", true},
		{"src/**/*.rs", "cmd/main.go", false},
		{"pkg/**/*_test.go", "pkg/domain/usecase/funnel_test.go", true},
		{"tests/**/*.py", "tests/unit/auth/test_login.py", true},
		{"**/__tests__/**/*.ts", "src/components/__tests__/button.spec.ts", true},
		{"**/*.go", "main.go", true},
		{"cmd/*", "cmd/app/main.go", false},
		{"cmd/**", "cmd/app/main.go", true},
	}

	for _, tt := range tests {
		result := MatchGlobRecursive(tt.pattern, tt.path)
		if result != tt.expected {
			t.Errorf("MatchGlobRecursive(%q, %q) = %v; esperava %v", tt.pattern, tt.path, result, tt.expected)
		}
	}
}
