/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/uuid.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

import "testing"

func TestGenerateRepoUUID(t *testing.T) {
	uuid1 := GenerateRepoUUID("ast-grep/ast-grep", "a1b2c3d4")
	uuid2 := GenerateRepoUUID("ast-grep/ast-grep", "a1b2c3d4")
	uuid3 := GenerateRepoUUID("ast-grep/ast-grep", "e5f6g7h8")

	if uuid1 != uuid2 {
		t.Fatalf("Esperava UUIDs determinísticos idênticos para o mesmo repo@commit")
	}

	if uuid1 == uuid3 {
		t.Fatalf("Esperava UUIDs diferentes para commits distintos")
	}

	if len(uuid1) != 36 {
		t.Errorf("Esperava tamanho de UUID de 36 caracteres, obteve %d", len(uuid1))
	}
}
