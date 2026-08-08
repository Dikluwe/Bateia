/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/linter_engine.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestLinterEngine_RunConcurrent(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "test-engine",
		Language:   "go",
	}

	tree := []domain.DirectoryNode{
		{Path: "main.go", Type: "file"},
	}

	engine := NewLinterEngine(2)
	violations := engine.RunConcurrent(manifest, tree)

	if len(violations) < 2 {
		t.Fatalf("Esperava pelo menos 2 violações (Missing Tests e DevOps), obteve %d", len(violations))
	}
}
