/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/raku_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMeta6Json_Basic(t *testing.T) {
	content := `{
  "name": "Cro::HTTP",
  "version": "0.8.7",
  "depends": [
    "Cro::Core",
    "IO::Socket::Async::SSL"
  ]
}`

	ref := domain.CommitRef{Owner: "cro", Repo: "Cro-HTTP"}
	manifest, err := ParseMeta6Json(ref, content)
	if err != nil {
		t.Fatalf("ParseMeta6Json falhou: %v", err)
	}

	if manifest.ModuleName != "Cro::HTTP" {
		t.Errorf("Esperado ModuleName 'Cro::HTTP', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "raku" {
		t.Errorf("Esperado Language 'raku', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (Cro::Core, IO::Socket::Async::SSL), obtido %d", len(manifest.Dependencies))
	}
}
