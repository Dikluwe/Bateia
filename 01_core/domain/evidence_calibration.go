/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_calibration.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

type EvidenceCalibrationObservation struct {
	ID        string              `json:"id"`
	Split     string              `json:"split"`
	Expected  EvidenceStatus      `json:"expected"`
	Observed  CapabilityEvidence  `json:"observed"`
	Evidence  []CollectedEvidence `json:"evidence,omitempty"`
	Diagnosis EvidenceMissCause   `json:"diagnosis,omitempty"`
}

type EvidenceCalibrationMetrics struct {
	Cases             int      `json:"cases"`
	ExpectedProven    int      `json:"expected_proven"`
	ObservedProven    int      `json:"observed_proven"`
	TrueProven        int      `json:"true_proven"`
	ProvenPrecision   *float64 `json:"proven_precision"`
	ProvenRecall      *float64 `json:"proven_recall"`
	OpaquePromotions  int      `json:"opaque_promotions"`
	InvalidSources    int      `json:"invalid_sources"`
	UndiagnosedMisses int      `json:"undiagnosed_misses"`
}

type EvidenceCalibrationReport struct {
	SchemaVersion string                                `json:"schema_version"`
	Verdict       string                                `json:"verdict"`
	BySplit       map[string]EvidenceCalibrationMetrics `json:"by_split"`
	Aggregate     EvidenceCalibrationMetrics            `json:"aggregate"`
}
