/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/mojo.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMojoProjectToml_Basic(t *testing.T) {
	content := `name = "max_app"
version = "0.1.0"

[dependencies]
max = "*"
`

	ref := domain.CommitRef{Owner: "modular", Repo: "max_app"}
	manifest, err := ParseMojoProjectToml(ref, content)
	if err != nil {
		t.Fatalf("ParseMojoProjectToml falhou: %v", err)
	}

	if manifest.ModuleName != "max_app" {
		t.Errorf("Esperado ModuleName 'max_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "mojo" {
		t.Errorf("Esperado Language 'mojo', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 1 {
		t.Fatalf("Esperado 1 dependência (max), obtido %d", len(manifest.Dependencies))
	}
}
