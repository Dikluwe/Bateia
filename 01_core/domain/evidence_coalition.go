/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_coalition.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain

const (
	EvidenceCoalitionInputSchemaVersion  = "evidence-coalition-input-v1"
	EvidenceCoalitionReportSchemaVersion = "evidence-coalition-report-v1"
)

type EvidenceCoalitionRole string

const (
	EvidenceCoalitionRoleImplementation  EvidenceCoalitionRole = "IMPLEMENTATION"
	EvidenceCoalitionRoleIndependentTest EvidenceCoalitionRole = "INDEPENDENT_TEST"
	EvidenceCoalitionRoleExecution       EvidenceCoalitionRole = "EXECUTION"
	EvidenceCoalitionRoleDocumentation   EvidenceCoalitionRole = "DOCUMENTATION"
)

type EvidenceCoalitionPolarity string

const (
	EvidenceCoalitionPolaritySupport  EvidenceCoalitionPolarity = "SUPPORT"
	EvidenceCoalitionPolarityConflict EvidenceCoalitionPolarity = "CONFLICT"
)

type EvidenceOriginStatus string

const (
	EvidenceOriginRecorded EvidenceOriginStatus = "RECORDED"
	EvidenceOriginSame     EvidenceOriginStatus = "SAME_ORIGIN"
	EvidenceOriginDistinct EvidenceOriginStatus = "DISTINCT_ORIGIN"
	EvidenceOriginUnknown  EvidenceOriginStatus = "UNKNOWN_ORIGIN"
)

type EvidenceCoalitionOutcome string

const (
	EvidenceCoalitionSatisfied                EvidenceCoalitionOutcome = "SATISFIED"
	EvidenceCoalitionConflict                 EvidenceCoalitionOutcome = "CONFLICT"
	EvidenceCoalitionMissingRoles             EvidenceCoalitionOutcome = "MISSING_ROLES"
	EvidenceCoalitionUnknownIndependence      EvidenceCoalitionOutcome = "UNKNOWN_INDEPENDENCE"
	EvidenceCoalitionInsufficientIndependence EvidenceCoalitionOutcome = "INSUFFICIENT_INDEPENDENCE"
	EvidenceCoalitionInvalidInput             EvidenceCoalitionOutcome = "INVALID_INPUT"
)

const (
	EvidenceCoalitionDecisionImpactNone = "NONE"
	EvidenceCoalitionReportComplete     = "COMPLETE"
	EvidenceCoalitionReportUnknown      = "UNKNOWN"
)

type EvidenceCoalitionInput struct {
	SchemaVersion string                   `json:"schema_version"`
	Rules         []EvidenceCoalitionRule  `json:"rules"`
	Claims        []EvidenceCoalitionClaim `json:"claims"`
}

type EvidenceCoalitionRule struct {
	ID                     string                  `json:"id"`
	RequiredRoles          []EvidenceCoalitionRole `json:"required_roles"`
	MinimumDistinctOrigins int                     `json:"minimum_distinct_origins"`
}

type EvidenceCoalitionClaim struct {
	ID             string                     `json:"id"`
	RuleID         string                     `json:"rule_id"`
	BaselineStatus EvidenceStatus             `json:"baseline_status"`
	Witnesses      []EvidenceCoalitionWitness `json:"witnesses"`
}

type EvidenceCoalitionWitness struct {
	ID           string                    `json:"id"`
	Role         EvidenceCoalitionRole     `json:"role"`
	Polarity     EvidenceCoalitionPolarity `json:"polarity"`
	OriginStatus EvidenceOriginStatus      `json:"origin_status"`
	OriginID     string                    `json:"origin_id"`
	Source       EvidenceSource            `json:"source"`
}

type EvidenceCoalitionAssessment struct {
	ClaimID               string                   `json:"claim_id"`
	RuleID                string                   `json:"rule_id"`
	BaselineStatus        EvidenceStatus           `json:"baseline_status"`
	DecisionImpact        string                   `json:"decision_impact"`
	Outcome               EvidenceCoalitionOutcome `json:"outcome"`
	Independence          EvidenceOriginStatus     `json:"independence"`
	SupportingWitnessIDs  []string                 `json:"supporting_witness_ids"`
	ConflictingWitnessIDs []string                 `json:"conflicting_witness_ids"`
	MissingRoles          []EvidenceCoalitionRole  `json:"missing_roles"`
	DistinctOrigins       []string                 `json:"distinct_origins"`
	Reason                string                   `json:"reason"`
}

type EvidenceCoalitionMetrics struct {
	Claims                   int  `json:"claims"`
	Witnesses                int  `json:"witnesses"`
	Satisfied                int  `json:"satisfied"`
	Conflicts                int  `json:"conflicts"`
	MissingRoles             int  `json:"missing_roles"`
	UnknownIndependence      int  `json:"unknown_independence"`
	InsufficientIndependence int  `json:"insufficient_independence"`
	Invalid                  int  `json:"invalid"`
	Operations               int  `json:"operations"`
	LimitExceeded            bool `json:"limit_exceeded"`
	VerdictChanges           int  `json:"verdict_changes"`
}

type EvidenceCoalitionReport struct {
	SchemaVersion string                        `json:"schema_version"`
	Status        string                        `json:"status"`
	Reason        string                        `json:"reason"`
	Assessments   []EvidenceCoalitionAssessment `json:"assessments"`
	Metrics       EvidenceCoalitionMetrics      `json:"metrics"`
}
