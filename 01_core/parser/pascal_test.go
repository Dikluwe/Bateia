/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/pascal_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePascalDelphi_Basic(t *testing.T) {
	content := `<Project xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <ItemGroup>
    <ItemName>vcl.dcp</ItemName>
    <ItemName>rtl.dcp</ItemName>
  </ItemGroup>
</Project>`

	ref := domain.CommitRef{Owner: "embarcadero", Repo: "DelphiApp"}
	manifest, err := ParsePascalDelphi(ref, content)
	if err != nil {
		t.Fatalf("ParsePascalDelphi falhou: %v", err)
	}

	if manifest.Language != "pascal" {
		t.Errorf("Esperado Language 'pascal', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vcl.dcp, rtl.dcp), obtido %d", len(manifest.Dependencies))
	}
}
