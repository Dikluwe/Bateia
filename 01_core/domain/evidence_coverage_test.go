/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_coverage_test.md
 * @layer L1
 * @updated 2026-08-26
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEvidenceMissCauseSchemaIsStable(t *testing.T) {
	causes := []EvidenceMissCause{
		MissTargetNotPlanned, MissTargetNotCollected, MissBudgetExhausted,
		MissLanguageUnsupported, MissPattern, MissSymbolLinkage,
		MissStrengthUnderestimated, MissInsufficientIndependence,
		MissNegativeConflict, MissGroundTruthAmbiguous,
	}
	for _, cause := range causes {
		if cause == "" {
			t.Fatal("causa vazia")
		}
	}
	data, err := json.Marshal(EvidenceCoverageReport{SchemaVersion: "1", CauseCounts: map[EvidenceMissCause]int{MissPattern: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"PATTERN_MISS":1`) {
		t.Fatalf("schema inesperado: %s", data)
	}
}
