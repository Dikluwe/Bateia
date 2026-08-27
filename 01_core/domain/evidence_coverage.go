/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_coverage.md
 * @layer L1
 * @updated 2026-08-26
 */

package domain

type EvidenceMissCause string

const (
	MissTargetNotPlanned         EvidenceMissCause = "TARGET_NOT_PLANNED"
	MissTargetNotCollected       EvidenceMissCause = "TARGET_NOT_COLLECTED"
	MissBudgetExhausted          EvidenceMissCause = "BUDGET_EXHAUSTED"
	MissLanguageUnsupported      EvidenceMissCause = "LANGUAGE_UNSUPPORTED"
	MissPattern                  EvidenceMissCause = "PATTERN_MISS"
	MissSymbolLinkage            EvidenceMissCause = "SYMBOL_LINKAGE_MISS"
	MissStrengthUnderestimated   EvidenceMissCause = "STRENGTH_UNDERESTIMATED"
	MissInsufficientIndependence EvidenceMissCause = "INSUFFICIENT_INDEPENDENCE"
	MissNegativeConflict         EvidenceMissCause = "NEGATIVE_CONFLICT"
	MissGroundTruthAmbiguous     EvidenceMissCause = "GROUND_TRUTH_AMBIGUOUS"
)

type EvidenceMissDiagnosis struct {
	Repository     string            `json:"repository"`
	Capability     string            `json:"capability"`
	Aliases        []string          `json:"aliases,omitempty"`
	Expected       EvidenceStatus    `json:"expected"`
	Observed       EvidenceStatus    `json:"observed"`
	PrimaryCause   EvidenceMissCause `json:"primary_cause"`
	SupportingData []EvidenceSource  `json:"supporting_data,omitempty"`
	MissingTargets []string          `json:"missing_targets,omitempty"`
	RuleID         string            `json:"rule_id,omitempty"`
	Explanation    string            `json:"explanation"`
}

type EvidenceCoverageMetrics struct {
	ExpectedProvenCount int      `json:"expected_proven_count"`
	EmittedProvenCount  int      `json:"emitted_proven_count"`
	ProvenRecall        *float64 `json:"proven_recall"`
	EvidenceMetrics
	MissCountByCause map[EvidenceMissCause]int `json:"miss_count_by_cause"`
}

type EvidenceCoverageReport struct {
	SchemaVersion string                    `json:"schema_version"`
	RulesVersion  string                    `json:"rules_version"`
	RulesHash     string                    `json:"rules_hash"`
	Diagnoses     []EvidenceMissDiagnosis   `json:"diagnoses"`
	CauseCounts   map[EvidenceMissCause]int `json:"cause_counts"`
	Metrics       EvidenceCoverageMetrics   `json:"metrics"`
}
