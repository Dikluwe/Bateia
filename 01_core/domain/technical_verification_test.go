/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_verification_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestTechnicalVerificationContract(t *testing.T) {
	v := domain.TechnicalEvidenceVerification{
		SourceAuthentication: "local-tree-sha256",
		IntegrityVerified:    true,
		BehaviorVerified:     false,
		FindingID:            "find-01",
		Candidate:            "owner/repo",
		Predicate:            domain.TechnicalLicenseID,
		Value:                "MIT",
		Status:               domain.TechnicalVerified,
		Message:              "Licença verificada contra texto canônico",
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.TechnicalEvidenceVerification
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Status != domain.TechnicalVerified || !restored.IntegrityVerified || restored.Value != "MIT" {
		t.Fatalf("verification mismatch: %+v", restored)
	}
}
