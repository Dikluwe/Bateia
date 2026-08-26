/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ocaml_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseOCamlDune_Basic(t *testing.T) {
	content := `(lang dune 3.0)
(name my_ocaml_app)
(depends dune alcotest dream)
`

	ref := domain.CommitRef{Owner: "ocaml-community", Repo: "my_ocaml_app"}
	manifest, err := ParseOCamlDune(ref, content)
	if err != nil {
		t.Fatalf("ParseOCamlDune falhou: %v", err)
	}

	if manifest.ModuleName != "my_ocaml_app" {
		t.Errorf("Esperado ModuleName 'my_ocaml_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "ocaml" {
		t.Errorf("Esperado Language 'ocaml', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (alcotest, dream), obtido %d", len(manifest.Dependencies))
	}
}
