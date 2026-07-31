/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/chapel.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMasonToml_Basic(t *testing.T) {
	content := `[brick]
name = "my_chapel_app"
version = "0.1.0"

[dependencies]
chpl_net = "*"
`

	ref := domain.CommitRef{Owner: "chapel-lang", Repo: "my_chapel_app"}
	manifest, err := ParseMasonToml(ref, content)
	if err != nil {
		t.Fatalf("ParseMasonToml falhou: %v", err)
	}

	if manifest.ModuleName != "my_chapel_app" {
		t.Errorf("Esperado ModuleName 'my_chapel_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "chapel" {
		t.Errorf("Esperado Language 'chapel', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 1 {
		t.Fatalf("Esperado 1 dependência (chpl_net), obtido %d", len(manifest.Dependencies))
	}
}
