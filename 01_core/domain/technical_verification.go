/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_verification.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

// Only dedicated verifiers may interpret these predicates. An arbitrary
// natural-language claim never becomes decision evidence.
type TechnicalEvidencePredicate string

const (
	TechnicalLicenseID     TechnicalEvidencePredicate = "LICENSE_ID"
	TechnicalSBOMComponent TechnicalEvidencePredicate = "SBOM_COMPONENT"
)

type TechnicalVerificationStatus string

const (
	TechnicalVerified     TechnicalVerificationStatus = "VERIFIED"
	TechnicalUnverified   TechnicalVerificationStatus = "UNVERIFIED"
	TechnicalUnsupported  TechnicalVerificationStatus = "UNSUPPORTED"
	TechnicalStale        TechnicalVerificationStatus = "STALE"
	TechnicalContradicted TechnicalVerificationStatus = "CONTRADICTED"
)

type TechnicalEvidenceVerification struct {
	SourceAuthentication string                      `json:"source_authentication"`
	IntegrityVerified    bool                        `json:"integrity_verified"`
	InterpretationBasis  string                      `json:"interpretation_basis,omitempty"`
	BehaviorVerified     bool                        `json:"behavior_verified"`
	FindingID            string                      `json:"finding_id"`
	Candidate            string                      `json:"candidate"`
	Predicate            TechnicalEvidencePredicate  `json:"predicate,omitempty"`
	Value                string                      `json:"value,omitempty"`
	Status               TechnicalVerificationStatus `json:"status"`
	Source               *EvidenceSource             `json:"source,omitempty"`
	Message              string                      `json:"message"`
}
