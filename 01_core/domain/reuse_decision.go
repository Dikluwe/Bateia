/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/reuse_decision.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

type ReuseDecision string

const (
	DecisionFork      ReuseDecision = "FORK"
	DecisionExtend    ReuseDecision = "EXTEND"
	DecisionCompose   ReuseDecision = "COMPOSE"
	DecisionIntegrate ReuseDecision = "INTEGRATE"
	DecisionReference ReuseDecision = "REFERENCE"
	DecisionBuild     ReuseDecision = "BUILD"
	DecisionUnknown   ReuseDecision = "INCONCLUSIVE"
)

type IntegrationEffort string

const (
	EffortLow     IntegrationEffort = "LOW"
	EffortMedium  IntegrationEffort = "MEDIUM"
	EffortHigh    IntegrationEffort = "HIGH"
	EffortUnknown IntegrationEffort = "INCONCLUSIVE"
)

type ReuseDecisionPolicy struct {
	SchemaVersion       string   `json:"schema_version"`
	AllowedLicenses     []string `json:"allowed_licenses"`
	RequireLocal        bool     `json:"require_local"`
	AllowArchived       bool     `json:"allow_archived"`
	MaxComponents       int      `json:"max_components"`
	ExtendMinCoverage   float64  `json:"extend_min_coverage"`
	IntegrationMinValue float64  `json:"integration_min_value"`
}

type ReuseDecisionRequest struct {
	Idea        IdeaSpec              `json:"idea"`
	Assessments []CandidateAssessment `json:"assessments"`
	Policy      ReuseDecisionPolicy   `json:"policy"`
	Composition *CompositionInput     `json:"composition,omitempty"`
}

type CandidateCapabilityStatus struct {
	Candidate string           `json:"candidate"`
	Status    EvidenceStatus   `json:"status"`
	Evidence  []EvidenceSource `json:"evidence,omitempty"`
}

type CapabilityComparison struct {
	Capability string                      `json:"capability"`
	Candidates []CandidateCapabilityStatus `json:"candidates"`
}

type DecisionReason struct {
	RuleID    string           `json:"rule_id"`
	Effect    string           `json:"effect"`
	Candidate string           `json:"candidate,omitempty"`
	Evidence  []EvidenceSource `json:"evidence,omitempty"`
	Message   string           `json:"message"`
}

type DecisionAlternative struct {
	Decision  ReuseDecision `json:"decision"`
	Candidate string        `json:"candidate,omitempty"`
	Coverage  float64       `json:"coverage"`
	Message   string        `json:"message"`
}

type EffortFactor struct {
	Factor   string `json:"factor"`
	Effect   string `json:"effect"`
	Evidence string `json:"evidence"`
}

type EffortAssessment struct {
	Method      string            `json:"method,omitempty"`
	Limitations string            `json:"limitations,omitempty"`
	Level       IntegrationEffort `json:"level"`
	Factors     []EffortFactor    `json:"factors"`
}

type ReuseDecisionReport struct {
	Investigation     *InvestigationPlan       `json:"investigation_plan,omitempty"`
	Composition       *CompositionReport       `json:"composition,omitempty"`
	TechnicalEvidence *TechnicalEvidenceReview `json:"technical_evidence,omitempty"`
	SchemaVersion     string                   `json:"schema_version"`
	Policy            ReuseDecisionPolicy      `json:"policy"`
	PolicyHash        string                   `json:"policy_hash"`
	Decision          ReuseDecision            `json:"decision"`
	Primary           *ResolvedRepository      `json:"primary,omitempty"`
	Components        []ReusableComponent      `json:"components,omitempty"`
	Matrix            []CapabilityComparison   `json:"matrix"`
	Gaps              []string                 `json:"gaps"`
	Blockers          []DecisionReason         `json:"blockers,omitempty"`
	Reasons           []DecisionReason         `json:"reasons"`
	Alternatives      []DecisionAlternative    `json:"alternatives"`
	Effort            EffortAssessment         `json:"effort"`
}
