/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_recheck.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

type TechnicalSourceRecheckRequest struct {
	Batch      TechnicalEvidenceBatch `json:"batch"`
	Candidates []RepositoryCandidate  `json:"candidates"`
}

type TechnicalSourceRecheckReport struct {
	SchemaVersion string                  `json:"schema_version"`
	Basis         string                  `json:"verification_basis"`
	Review        TechnicalEvidenceReview `json:"review"`
	Bundles       []EvidenceBundle        `json:"bundles"`
	Failures      []ExecutionFailure      `json:"failures,omitempty"`
}
