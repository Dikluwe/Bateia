/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/capability_compilation_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestCapabilityCompilationJSONContract(t *testing.T) {
	original := domain.CapabilityCompilation{
		SchemaVersion: "1", IntakeHash: "hash", IntakeStatus: domain.IntakeReady,
		Goal: "buscar arquivos", Users: []string{"analista"},
		Requirements: []domain.CompiledRequirement{{
			ID: "req-1", Kind: "capability", Text: "buscar trechos",
			SourcePath: "idea.capabilities[0]", Authority: domain.AuthorityUser, Priority: "MUST",
		}},
		Questions: []domain.CompiledEvidenceQuestion{{
			ID: "eq-1", RequirementID: "req-1", Channel: "tests",
			Question: "que teste?", Support: "teste", Refute: "falha",
			Limit: "não é prova", Status: "PENDING",
		}},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.CapabilityCompilation
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Requirements[0].SourcePath != "idea.capabilities[0]" || restored.Questions[0].Status != "PENDING" {
		t.Fatalf("lost traceability or state: %+v", restored)
	}
}
