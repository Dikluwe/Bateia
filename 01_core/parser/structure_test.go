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
