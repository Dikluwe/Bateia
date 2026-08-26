/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/fortran_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseFortranFpm_Basic(t *testing.T) {
	content := `name = "fast_matrix"
version = "1.0.0"

[dependencies]
lapack = "*"
stdlib = "*"
`

	ref := domain.CommitRef{Owner: "fortran-lang", Repo: "fast_matrix"}
	manifest, err := ParseFortranFpm(ref, content)
	if err != nil {
		t.Fatalf("ParseFortranFpm falhou: %v", err)
	}

	if manifest.ModuleName != "fast_matrix" {
		t.Errorf("Esperado ModuleName 'fast_matrix', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "fortran" {
		t.Errorf("Esperado Language 'fortran', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (lapack, stdlib), obtido %d", len(manifest.Dependencies))
	}
}
