/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cobol_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseCobolSource_Basic(t *testing.T) {
	content := `       IDENTIFICATION DIVISION.
       PROGRAM-ID. HELLO.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       COPY CUSTREC.
       PROCEDURE DIVISION.
           CALL 'DBLOGIC' USING CUST-ID.
           STOP RUN.`

	ref := domain.CommitRef{Owner: "mainframe", Repo: "banking_cobol"}
	manifest, err := ParseCobolSource(ref, content)
	if err != nil {
		t.Fatalf("ParseCobolSource falhou: %v", err)
	}

	if manifest.Language != "cobol" {
		t.Errorf("Esperado Language 'cobol', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (CUSTREC, DBLOGIC), obtido %d", len(manifest.Dependencies))
	}
}
