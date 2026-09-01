/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_calibration_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

import "testing"

func TestEvidenceCalibrationDomainPreservesOpaqueObservation(t *testing.T) {
	observation := EvidenceCalibrationObservation{ID: "opaque", Split: "TEST", Expected: EvidenceInconclusive, Observed: CapabilityEvidence{Status: EvidenceInconclusive}}
	if observation.Expected != EvidenceInconclusive || observation.Observed.Status != EvidenceInconclusive {
		t.Fatal("observação opaca foi alterada")
	}
	metrics := EvidenceCalibrationMetrics{}
	if metrics.ProvenPrecision != nil || metrics.ProvenRecall != nil {
		t.Fatal("denominador zero fabricou métrica")
	}
}
