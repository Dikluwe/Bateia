/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/fsharp_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseFsProj_Basic(t *testing.T) {
	content := `<Project Sdk="Microsoft.NET.Sdk">
  <ItemGroup>
    <PackageReference Include="Giraffe" Version="6.0.0" />
    <PackageReference Include="Expecto" Version="9.0.0" />
  </ItemGroup>
</Project>`

	ref := domain.CommitRef{Owner: "giraffe-fsharp", Repo: "GiraffeApp"}
	manifest, err := ParseFsProj(ref, content)
	if err != nil {
		t.Fatalf("ParseFsProj falhou: %v", err)
	}

	if manifest.Language != "fsharp" {
		t.Errorf("Esperado Language 'fsharp', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 PackageReferences (Giraffe, Expecto), obtido %d", len(manifest.Dependencies))
	}
}
