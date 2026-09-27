/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_sbom_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestTechnicalSBOMDeclarationVerification(t *testing.T) {
	validSBOM := `{"bomFormat": "CycloneDX", "specVersion": "1.6", "components": [{"bom-ref": "comp-1", "name": "target-lib", "version": "1.0.0"}]}`

	status, msg := verifyTechnicalSBOMDeclaration(validSBOM, "target-lib@1.0.0")
	if status != domain.TechnicalVerified {
		t.Fatalf("expected TechnicalVerified, got %v (%s)", status, msg)
	}

	statusMissing, msgMissing := verifyTechnicalSBOMDeclaration(validSBOM, "absent-lib@1.0.0")
	if statusMissing != domain.TechnicalUnverified {
		t.Fatalf("expected TechnicalUnverified, got %v (%s)", statusMissing, msgMissing)
	}

	invalidJSON := `{"bomFormat": "CycloneDX", invalid`
	statusInv, msgInv := verifyTechnicalSBOMDeclaration(invalidJSON, "target-lib@1.0.0")
	if statusInv != domain.TechnicalUnverified {
		t.Fatalf("expected TechnicalUnverified, got %v (%s)", statusInv, msgInv)
	}
}
