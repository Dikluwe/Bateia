/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/extended_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseDubJson_Basic(t *testing.T) {
	content := `{
  "name": "vibe_app",
  "description": "A simple vibe.d application.",
  "dependencies": {
    "vibe-d": "~>0.9.5",
    "mir-algorithm": ">=3.10.0"
  }
}`

	ref := domain.CommitRef{Owner: "dlang", Repo: "vibe_app"}
	manifest, err := ParseDubJson(ref, content)
	if err != nil {
		t.Fatalf("ParseDubJson falhou: %v", err)
	}

	if manifest.ModuleName != "vibe_app" {
		t.Errorf("Esperado ModuleName 'vibe_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "d" {
		t.Errorf("Esperado Language 'd', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vibe-d, mir-algorithm), obtido %d", len(manifest.Dependencies))
	}
}

func TestParseGleamToml_Basic(t *testing.T) {
	content := `
name = "my_gleam_app"
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

func TestParseBoxJson_Basic(t *testing.T) {
	content := `{
  "name": "coldbox_app",
  "version": "1.0.0",
  "dependencies": {
    "coldbox": "^6.8.0",
    "cborm": "^5.0.0"
  }
}`

	ref := domain.CommitRef{Owner: "ortus", Repo: "coldbox_app"}
	manifest, err := ParseBoxJson(ref, content)
	if err != nil {
		t.Fatalf("ParseBoxJson falhou: %v", err)
	}

	if manifest.Language != "coldfusion" {
		t.Errorf("Esperado Language 'coldfusion', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (coldbox, cborm), obtido %d", len(manifest.Dependencies))
	}
}

func TestParseFlexConfig_Basic(t *testing.T) {
	content := `<flex-config>
  <compiler>
    <library-path>
      <path-element>libs/robotlegs.swc</path-element>
      <path-element>libs/flexunit.swc</path-element>
    </library-path>
  </compiler>
</flex-config>`

	ref := domain.CommitRef{Owner: "adobe", Repo: "FlexApp"}
	manifest, err := ParseFlexConfig(ref, content)
	if err != nil {
		t.Fatalf("ParseFlexConfig falhou: %v", err)
	}

	if manifest.Language != "actionscript" {
		t.Errorf("Esperado Language 'actionscript', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (robotlegs.swc, flexunit.swc), obtido %d", len(manifest.Dependencies))
	}
}
