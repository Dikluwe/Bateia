/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/odin_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseOlsJson_Basic(t *testing.T) {
	content := `{
  "collections": [
    { "name": "core", "path": "core" },
    { "name": "vendor", "path": "vendor" }
  ]
}`

	ref := domain.CommitRef{Owner: "odin-lang", Repo: "OdinGame"}
	manifest, err := ParseOlsJson(ref, content)
	if err != nil {
		t.Fatalf("ParseOlsJson falhou: %v", err)
	}

	if manifest.Language != "odin" {
		t.Errorf("Esperado Language 'odin', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 coleções (core, vendor), obtido %d", len(manifest.Dependencies))
	}
}
