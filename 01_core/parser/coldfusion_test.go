/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/coldfusion.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseBoxJson_Basic(t *testing.T) {
	content := `{
  "name": "coldbox_app",
  "version": "1.0.0",
  "dependencies": {
    "coldbox": "^6.8.0",
    "cborm": "^5.0.0"
  }
}`

	ref := domain.CommitRef{Owner: "ortus", Repo: "coldbox_app"}
	manifest, err := ParseBoxJson(ref, content)
	if err != nil {
		t.Fatalf("ParseBoxJson falhou: %v", err)
	}

	if manifest.Language != "coldfusion" {
		t.Errorf("Esperado Language 'coldfusion', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (coldbox, cborm), obtido %d", len(manifest.Dependencies))
	}
}
