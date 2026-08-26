/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

import (
	"testing"
)

func TestEvidenceTypes(t *testing.T) {
	if EvidenceProven != "PROVEN" {
		t.Errorf("expected PROVEN, got %s", EvidenceProven)
	}
}
