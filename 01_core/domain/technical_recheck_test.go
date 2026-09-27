/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_recheck_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestTechnicalRecheckContract(t *testing.T) {
	req := domain.TechnicalSourceRecheckRequest{
		Candidates: []domain.RepositoryCandidate{{Owner: "o", Repo: "r"}},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.TechnicalSourceRecheckRequest
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored.Candidates) != 1 || restored.Candidates[0].Repo != "r" {
		t.Fatalf("recheck request mismatch: %+v", restored)
	}
}
