/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/julia.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseJuliaProjectToml_Basic(t *testing.T) {
	content := `name = "DataAnalysisApp"
uuid = "12345678-1234-1234-1234-1234567890ab"
version = "0.1.0"

[deps]
DataFrames = "a93c6304-b5b0-51b6-ad7b-0e0084b60433"
Plots = "91a5bc11-a23a-5474-8092-5641e46f6d8d"
`

	ref := domain.CommitRef{Owner: "julia-data", Repo: "DataAnalysisApp"}
	manifest, err := ParseJuliaProjectToml(ref, content)
	if err != nil {
		t.Fatalf("ParseJuliaProjectToml falhou: %v", err)
	}

	if manifest.ModuleName != "DataAnalysisApp" {
		t.Errorf("Esperado ModuleName 'DataAnalysisApp', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "julia" {
		t.Errorf("Esperado Language 'julia', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (DataFrames, Plots), obtido %d", len(manifest.Dependencies))
	}
}
