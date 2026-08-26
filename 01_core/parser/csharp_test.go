/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/csharp_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseCsProj_Basic(t *testing.T) {
	content := `<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="Microsoft.EntityFrameworkCore.SqlServer" Version="8.0.0" />
    <PackageReference Include="Serilog.AspNetCore" Version="8.0.0" />
    <PackageReference Include="xunit" Version="2.6.2" />
    <PackageReference Include="Microsoft.NET.Test.Sdk" Version="17.8.0" />
  </ItemGroup>
</Project>`

	ref := domain.CommitRef{Owner: "dotnet", Repo: "sample-api"}
	manifest, err := ParseCsProj(ref, content)
	if err != nil {
		t.Fatalf("ParseCsProj falhou: %v", err)
	}

	if manifest.Language != "csharp" {
		t.Errorf("Esperado Language 'csharp', obtido '%s'", manifest.Language)
	}

	directCount := 0
	devCount := 0
	for _, dep := range manifest.Dependencies {
		if dep.Type == "direct" {
			directCount++
		} else if dep.Type == "dev" {
			devCount++
		}
	}

	if directCount != 2 {
		t.Errorf("Esperado 2 dependências direct (EF Core, Serilog), obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (xunit, Test.Sdk), obtido %d", devCount)
	}
}
