/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/perl_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseCpanfile_Basic(t *testing.T) {
	content := `requires 'perl', '5.008001';
requires 'Dancer2', '0.300000';
requires 'DBI', '1.640';
`

	ref := domain.CommitRef{Owner: "perl-community", Repo: "DancerApp"}
	manifest, err := ParseCpanfile(ref, content)
	if err != nil {
		t.Fatalf("ParseCpanfile falhou: %v", err)
	}

	if manifest.Language != "perl" {
		t.Errorf("Esperado Language 'perl', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (Dancer2, DBI), obtido %d", len(manifest.Dependencies))
	}
}
