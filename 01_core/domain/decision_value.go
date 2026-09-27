/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/decision_value.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

const DecisionValueSchemaVersion = "1"

type DecisionImpactClass string

const (
	DecisionPivotal                   DecisionImpactClass = "DECISION_PIVOTAL"
	DecisionNeutral                   DecisionImpactClass = "DECISION_NEUTRAL"
	DecisionCounterfactualUnsupported DecisionImpactClass = "COUNTERFACTUAL_UNSUPPORTED"
)

type InvestigationMethod string

type InvestigationSecurityRisk string

const (
	InvestigationBoundedRead   InvestigationSecurityRisk = "BOUNDED_UNTRUSTED_READ"
	InvestigationHostExecution InvestigationSecurityRisk = "CANDIDATE_CODE_EXECUTION"
)

const (
	InvestigationInspectArtifact      InvestigationMethod = "INSPECT_ARTIFACT"
	InvestigationRetryProvider        InvestigationMethod = "RETRY_PROVIDER"
	InvestigationRunCompatibilityTest InvestigationMethod = "RUN_COMPATIBILITY_TEST"
)

type InvestigationCost struct {
	ProviderCalls      int  `json:"provider_calls"`
	Documents          int  `json:"documents"`
	CodeFiles          int  `json:"code_files"`
	Tests              int  `json:"tests"`
	CandidateExecution bool `json:"candidate_execution"`
}

type CounterfactualOutcome struct {
	OutcomeCode string         `json:"outcome_code"`
	Status      EvidenceStatus `json:"status,omitempty"`
	Decision    ReuseDecision  `json:"decision"`
}

type InvestigationValueOption struct {
	SecurityRisk          InvestigationSecurityRisk `json:"security_risk"`
	SecurityRationale     string                    `json:"security_rationale"`
	ActionID              string                    `json:"action_id"`
	FindingID             string                    `json:"finding_id"`
	Target                RecoveryTarget            `json:"target"`
	Method                InvestigationMethod       `json:"method"`
	Impact                DecisionImpactClass       `json:"impact"`
	EpistemicBasis        string                    `json:"epistemic_basis"`
	Outcomes              []CounterfactualOutcome   `json:"outcomes"`
	PossibleDecisions     []ReuseDecision           `json:"possible_decisions"`
	ResolvedOutcomeCount  int                       `json:"resolved_outcome_count"`
	Cost                  InvestigationCost         `json:"cost"`
	CostParetoDominated   bool                      `json:"cost_pareto_dominated"`
	RequiresAuthorization bool                      `json:"requires_authorization"`
	Executable            bool                      `json:"executable"`
	Reason                string                    `json:"reason"`
}

type DecisionValueRequest struct {
	Idea                    IdeaSpec              `json:"idea"`
	Assessment              AssessmentReport      `json:"assessment"`
	DecisionPolicy          ReuseDecisionPolicy   `json:"decision_policy"`
	Diagnosis               InconclusiveDiagnosis `json:"diagnosis"`
	RecoveryPolicy          RecoveryPolicy        `json:"recovery_policy"`
	Round                   int                   `json:"round"`
	PreviousTargets         []RecoveryTarget      `json:"previous_targets"`
	AllowCandidateExecution bool                  `json:"allow_candidate_execution"`
}

type DecisionValuePlan struct {
	SchemaVersion       string                     `json:"schema_version"`
	Method              string                     `json:"method"`
	BaselineDecision    ReuseDecision              `json:"baseline_decision"`
	Options             []InvestigationValueOption `json:"options"`
	RecommendedActionID string                     `json:"recommended_action_id,omitempty"`
	ExecutablePlan      EvidenceRecoveryPlan       `json:"executable_plan"`
	Questions           []InformationNeed          `json:"questions"`
	StopReason          string                     `json:"stop_reason"`
}
