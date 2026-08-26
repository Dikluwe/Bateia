/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/r_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseRDescription_Basic(t *testing.T) {
	content := `Package: DataVizApp
Title: Data Visualization Tool
Version: 1.0.0
Imports: ggplot2,
    dplyr (>= 1.0.0),
    shiny
`

	ref := domain.CommitRef{Owner: "r-stats", Repo: "DataVizApp"}
	manifest, err := ParseRDescription(ref, content)
	if err != nil {
		t.Fatalf("ParseRDescription falhou: %v", err)
	}

	if manifest.ModuleName != "DataVizApp" {
		t.Errorf("Esperado ModuleName 'DataVizApp', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "r" {
		t.Errorf("Esperado Language 'r', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 3 {
		t.Fatalf("Esperado 3 dependências (ggplot2, dplyr, shiny), obtido %d", len(manifest.Dependencies))
	}
}
