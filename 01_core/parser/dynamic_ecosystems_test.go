/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/dynamic_ecosystems.md
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
    <PackageReference Include="FSharp.Core" Version="7.0.0" />
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
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}

func TestParseCpanfile_Basic(t *testing.T) {
	content := `requires 'Dancer2', '0.400000';
requires 'DBI', '1.643';
`

	ref := domain.CommitRef{Owner: "perl-community", Repo: "WebApp"}
	manifest, err := ParseCpanfile(ref, content)
	if err != nil {
		t.Fatalf("ParseCpanfile falhou: %v", err)
	}

	if manifest.Language != "perl" {
		t.Errorf("Esperado Language 'perl', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}

func TestParseCobolSource_Basic(t *testing.T) {
	content := `       IDENTIFICATION DIVISION.
       PROGRAM-ID. HELLO-WORLD.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       COPY CUSTREC.
       COPY WORKREC.
       PROCEDURE DIVISION.
           DISPLAY "HELLO WORLD".
           STOP RUN.`

	ref := domain.CommitRef{Owner: "mainframe", Repo: "COBOLApp"}
	manifest, err := ParseCobolSource(ref, content)
	if err != nil {
		t.Fatalf("ParseCobolSource falhou: %v", err)
	}

	if manifest.Language != "cobol" {
		t.Errorf("Esperado Language 'cobol', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 COPYbooks (CUSTREC, WORKREC), obtido %d", len(manifest.Dependencies))
	}
}
