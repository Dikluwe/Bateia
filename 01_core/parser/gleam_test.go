/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/gleam.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseGleamToml_Basic(t *testing.T) {
	content := `name = "my_gleam_app"
version = "1.0.0"

[dependencies]
lustre = ">= 1.0.0"
gleam_http = ">= 3.0.0"
`

	ref := domain.CommitRef{Owner: "gleam-lang", Repo: "my_gleam_app"}
	manifest, err := ParseGleamToml(ref, content)
	if err != nil {
		t.Fatalf("ParseGleamToml falhou: %v", err)
	}

	if manifest.ModuleName != "my_gleam_app" {
		t.Errorf("Esperado ModuleName 'my_gleam_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "gleam" {
		t.Errorf("Esperado Language 'gleam', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (lustre, gleam_http), obtido %d", len(manifest.Dependencies))
	}
}
