/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/v_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseVMod_Basic(t *testing.T) {
	content := `Module {
	name: 'vweb_app'
	description: 'Simple web app in V'
	version: '0.1.0'
	dependencies: [
		'vweb',
		'sqlite'
	]
}`

	ref := domain.CommitRef{Owner: "vlang", Repo: "vweb_app"}
	manifest, err := ParseVMod(ref, content)
	if err != nil {
		t.Fatalf("ParseVMod falhou: %v", err)
	}

	if manifest.ModuleName != "vweb_app" {
		t.Errorf("Esperado ModuleName 'vweb_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "v" {
		t.Errorf("Esperado Language 'v', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vweb, sqlite), obtido %d", len(manifest.Dependencies))
	}
}
