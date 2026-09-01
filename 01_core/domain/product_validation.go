/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/product_validation.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

type DiscoveryArmCompletion struct {
	Arm       string `json:"arm"`
	Planned   int    `json:"planned"`
	Completed int    `json:"completed"`
}

type ProductValidationInput struct {
	SchemaVersion string                   `json:"schema_version"`
	Domain        string                   `json:"domain"`
	CaseIDs       []string                 `json:"case_ids"`
	DiscoveryArms []DiscoveryArmCompletion `json:"discovery_arms"`
	SourceReceipt string                   `json:"source_receipt"`
	Oracle        string                   `json:"oracle"`
	Evidence      string                   `json:"evidence_calibration"`
}

type ProductArmMetric struct {
	Arm            string  `json:"arm"`
	Planned        int     `json:"planned"`
	Completed      int     `json:"completed"`
	CompletionRate float64 `json:"completion_rate"`
	Passed         bool    `json:"passed"`
}

type ProductValidationReport struct {
	SchemaVersion string             `json:"schema_version"`
	Domain        string             `json:"domain"`
	Verdict       string             `json:"verdict"`
	FailedPhase   string             `json:"failed_phase,omitempty"`
	Discovery     []ProductArmMetric `json:"discovery"`
	Evidence      string             `json:"evidence_status"`
	Decision      string             `json:"decision_status"`
	Utility       string             `json:"utility_status"`
	Diagnostics   []string           `json:"diagnostics"`
}
