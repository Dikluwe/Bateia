/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_evidence_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestTechnicalEvidenceContract(t *testing.T) {
	te := domain.TechnicalEvidence{
		ID:        "te-01",
		Technique: domain.TechniqueLicense,
		Candidate: "owner/repo",
		Claim:     "Licença MIT declarada",
		Predicate: domain.TechnicalLicenseID,
		Value:     "MIT",
		Result:    domain.TechnicalSupports,
		Source: domain.TechnicalEvidenceSource{
			Ref:      "README.md",
			Revision: "HEAD",
			Digest:   "sha256:abc",
		},
		Method:     "inspeção de texto canônico",
		ObservedAt: "2026-09-26T00:00:00Z",
	}
	batch := domain.TechnicalEvidenceBatch{
		SchemaVersion: "1",
		Findings:      []domain.TechnicalEvidence{te},
	}
	data, err := json.Marshal(batch)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.TechnicalEvidenceBatch
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored.Findings) != 1 || restored.Findings[0].Predicate != domain.TechnicalLicenseID {
		t.Fatalf("restored technical evidence mismatch: %+v", restored)
	}
}
