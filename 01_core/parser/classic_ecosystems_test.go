/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/classic_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseFortranFpm_Basic(t *testing.T) {
	content := `
name = "lapack_wrapper"
version = "0.1.0"

[dependencies]
stdlib = "*"
minpack = { git = "https://github.com/fortran-lang/minpack" }
`

	ref := domain.CommitRef{Owner: "fortran-lang", Repo: "lapack_wrapper"}
	manifest, err := ParseFortranFpm(ref, content)
	if err != nil {
		t.Fatalf("ParseFortranFpm falhou: %v", err)
	}

	if manifest.ModuleName != "lapack_wrapper" {
		t.Errorf("Esperado ModuleName 'lapack_wrapper', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "fortran" {
		t.Errorf("Esperado Language 'fortran', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (stdlib, minpack), obtido %d", len(manifest.Dependencies))
	}
}

func TestParsePascalDelphi_Basic(t *testing.T) {
	content := `<Project xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <ItemGroup>
    <DCCReference Include="MainForm.pas" />
    <DCCReference Include="DatabaseModule.pas" />
  </ItemGroup>
</Project>`

	ref := domain.CommitRef{Owner: "delphi", Repo: "MyDelphiApp"}
	manifest, err := ParsePascalDelphi(ref, content)
	if err != nil {
		t.Fatalf("ParsePascalDelphi falhou: %v", err)
	}

	if manifest.Language != "pascal" {
		t.Errorf("Esperado Language 'pascal', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências DCCReference, obtido %d", len(manifest.Dependencies))
	}
}
