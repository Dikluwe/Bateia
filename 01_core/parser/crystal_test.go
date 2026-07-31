/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/crystal.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseShardYaml_Basic(t *testing.T) {
	content := `name: kemal_app
version: 0.1.0

dependencies:
  kemal:
    github: kemalcr/kemal
  db:
    github: crystal-lang/crystal-db
`

	ref := domain.CommitRef{Owner: "crystal-community", Repo: "kemal_app"}
	manifest, err := ParseShardYaml(ref, content)
	if err != nil {
		t.Fatalf("ParseShardYaml falhou: %v", err)
	}

	if manifest.ModuleName != "kemal_app" {
		t.Errorf("Esperado ModuleName 'kemal_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "crystal" {
		t.Errorf("Esperado Language 'crystal', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (kemal, db), obtido %d", len(manifest.Dependencies))
	}
}
