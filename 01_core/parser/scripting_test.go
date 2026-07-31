/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/scripting.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseJuliaProjectToml_Basic(t *testing.T) {
	content := `
name = "DataAnalysis"
uuid = "1234-5678"

[deps]
DataFrames = "a93c63e6-5b32-5aa1-634d-178039b1d722"
CSV = "336ed68f-a31e-5942-800f-90e4a6fa8087"
`

	ref := domain.CommitRef{Owner: "julia", Repo: "DataAnalysis.jl"}
	manifest, err := ParseJuliaProjectToml(ref, content)
	if err != nil {
		t.Fatalf("ParseJuliaProjectToml falhou: %v", err)
	}

	if manifest.ModuleName != "DataAnalysis" {
		t.Errorf("Esperado ModuleName 'DataAnalysis', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "julia" {
		t.Errorf("Esperado Language 'julia', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}

func TestParseRDescription_Basic(t *testing.T) {
	content := `
Package: ggplot2
Type: Package
Imports:
    cli,
    glue,
    rlang (>= 1.1.0)
Suggests:
    covr,
    testthat (>= 3.1.5)
`

	ref := domain.CommitRef{Owner: "tidyverse", Repo: "ggplot2"}
	manifest, err := ParseRDescription(ref, content)
	if err != nil {
		t.Fatalf("ParseRDescription falhou: %v", err)
	}

	if manifest.ModuleName != "ggplot2" {
		t.Errorf("Esperado ModuleName 'ggplot2', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "r" {
		t.Errorf("Esperado Language 'r', obtido '%s'", manifest.Language)
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

	if directCount != 3 {
		t.Errorf("Esperado 3 dependências direct, obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (covr, testthat), obtido %d", devCount)
	}
}
