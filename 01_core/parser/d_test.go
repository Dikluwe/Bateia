/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/d.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseDubJson_Basic(t *testing.T) {
	content := `{
  "name": "vibe_app",
  "dependencies": {
    "vibe-d": "~>0.9.5",
    "mir-algorithm": ">=3.10.0"
  }
}`

	ref := domain.CommitRef{Owner: "dlang", Repo: "vibe_app"}
	manifest, err := ParseDubJson(ref, content)
	if err != nil {
		t.Fatalf("ParseDubJson falhou: %v", err)
	}

	if manifest.ModuleName != "vibe_app" {
		t.Errorf("Esperado ModuleName 'vibe_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "d" {
		t.Errorf("Esperado Language 'd', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vibe-d, mir-algorithm), obtido %d", len(manifest.Dependencies))
	}
}
