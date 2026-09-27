/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/technical_evidence.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

// TechnicalEvidence is an external, claim-level observation. It is deliberately
// separate from EvidenceSource: importing a finding must not make it admissible
// proof of a repository capability without the normal evidence checks.
type TechnicalEvidence struct {
	ID          string                     `json:"id"`
	Technique   TechnicalEvidenceTechnique `json:"technique"`
	Candidate   string                     `json:"candidate"`
	Claim       string                     `json:"claim"`
	Predicate   TechnicalEvidencePredicate `json:"predicate,omitempty"`
	Value       string                     `json:"value,omitempty"`
	Result      TechnicalEvidenceResult    `json:"result"`
	Source      TechnicalEvidenceSource    `json:"source"`
	Method      string                     `json:"method"`
	ObservedAt  string                     `json:"observed_at"`
	ValidUntil  string                     `json:"valid_until,omitempty"`
	Limitations []string                   `json:"limitations,omitempty"`
}

type TechnicalEvidenceTechnique string

const (
	TechniqueCodeSearch           TechnicalEvidenceTechnique = "CODE_SEARCH"
	TechniqueRepositoryMining     TechnicalEvidenceTechnique = "REPOSITORY_MINING"
	TechniqueAPIRecommendation    TechnicalEvidenceTechnique = "API_RECOMMENDATION"
	TechniqueSoftwareArchaeology  TechnicalEvidenceTechnique = "SOFTWARE_ARCHAEOLOGY"
	TechniqueArchitectureRecovery TechnicalEvidenceTechnique = "ARCHITECTURE_RECOVERY"
	TechniqueImpactAnalysis       TechnicalEvidenceTechnique = "IMPACT_ANALYSIS"
	TechniqueSBOM                 TechnicalEvidenceTechnique = "SBOM_SCA"
	TechniqueLicense              TechnicalEvidenceTechnique = "LICENSE_COMPATIBILITY"
	TechniqueProjectHealth        TechnicalEvidenceTechnique = "PROJECT_HEALTH"
	TechniqueBuildProvenance      TechnicalEvidenceTechnique = "BUILD_PROVENANCE"
)

type TechnicalEvidenceResult string

const (
	TechnicalSupports     TechnicalEvidenceResult = "SUPPORTS"
	TechnicalRefutes      TechnicalEvidenceResult = "REFUTES"
	TechnicalInconclusive TechnicalEvidenceResult = "INCONCLUSIVE"
)

type TechnicalEvidenceSource struct {
	Ref      string `json:"ref"`
	Revision string `json:"revision,omitempty"`
	Digest   string `json:"digest,omitempty"`
}

type TechnicalEvidenceBatch struct {
	SchemaVersion string              `json:"schema_version"`
	Findings      []TechnicalEvidence `json:"findings"`
}

type TechnicalClaimStatus string

const (
	TechnicalClaimSupported    TechnicalClaimStatus = "SUPPORTED"
	TechnicalClaimRefuted      TechnicalClaimStatus = "REFUTED"
	TechnicalClaimConflicting  TechnicalClaimStatus = "CONFLICTING"
	TechnicalClaimInconclusive TechnicalClaimStatus = "INCONCLUSIVE"
	TechnicalClaimStale        TechnicalClaimStatus = "STALE"
)

type TechnicalClaimReview struct {
	Candidate   string               `json:"candidate"`
	Claim       string               `json:"claim"`
	Status      TechnicalClaimStatus `json:"status"`
	FindingIDs  []string             `json:"finding_ids"`
	StaleIDs    []string             `json:"stale_ids,omitempty"`
	Limitations []string             `json:"limitations,omitempty"`
}

type TechnicalEvidenceReview struct {
	ContentPolicy  string                          `json:"content_policy"`
	SchemaVersion  string                          `json:"schema_version"`
	AsOf           string                          `json:"as_of"`
	DecisionImpact string                          `json:"decision_impact"`
	Claims         []TechnicalClaimReview          `json:"claims"`
	Verifications  []TechnicalEvidenceVerification `json:"verifications,omitempty"`
	Warnings       []string                        `json:"warnings,omitempty"`
}
