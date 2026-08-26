/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ada_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseAlireToml_Basic(t *testing.T) {
	content := `name = "my_ada_app"
version = "1.0.0"

[dependencies]
gnatcoll = "*"
sdlz = "^2.0.0"
`

	ref := domain.CommitRef{Owner: "ada-lang", Repo: "my_ada_app"}
	manifest, err := ParseAlireToml(ref, content)
	if err != nil {
		t.Fatalf("ParseAlireToml falhou: %v", err)
	}

	if manifest.ModuleName != "my_ada_app" {
		t.Errorf("Esperado ModuleName 'my_ada_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "ada" {
		t.Errorf("Esperado Language 'ada', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (gnatcoll, sdlz), obtido %d", len(manifest.Dependencies))
	}
}
