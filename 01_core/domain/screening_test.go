/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/screening_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUncalibratedDecisionSerializesNullProbability(t *testing.T) {
	data, err := json.Marshal(ScreeningDecision{Repository: "org/repo", CalibrationStatus: "uncalibrated"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"probability":null`) {
		t.Fatalf("ausência de probabilidade não explícita: %s", data)
	}
}
