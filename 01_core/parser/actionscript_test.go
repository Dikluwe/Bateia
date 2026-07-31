/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/actionscript.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

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
