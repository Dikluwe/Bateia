/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/haxe_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseHaxelibJson_Basic(t *testing.T) {
	content := `{
  "name": "heaps_game",
  "version": "1.0.0",
  "dependencies": {
    "heaps": "",
    "hlsdl": ""
  }
}`

	ref := domain.CommitRef{Owner: "haxe-community", Repo: "heaps_game"}
	manifest, err := ParseHaxelibJson(ref, content)
	if err != nil {
		t.Fatalf("ParseHaxelibJson falhou: %v", err)
	}

	if manifest.Language != "haxe" {
		t.Errorf("Esperado Language 'haxe', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (heaps, hlsdl), obtido %d", len(manifest.Dependencies))
	}
}
