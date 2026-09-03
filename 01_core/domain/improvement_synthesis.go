/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/improvement_synthesis.md
 * @layer L1
 * @updated 2026-09-02
 */

package domain

import "context"

type ImprovementHypothesisStatus string

const (
	HypothesisProposed     ImprovementHypothesisStatus = "PROPOSED"
	HypothesisSupported    ImprovementHypothesisStatus = "SUPPORTED"
	HypothesisConflicting  ImprovementHypothesisStatus = "CONFLICTING"
	HypothesisInconclusive ImprovementHypothesisStatus = "INCONCLUSIVE"
)

type ImprovementSynthesisStatus string

const (
	ImprovementComplete     ImprovementSynthesisStatus = "COMPLETE"
	ImprovementPartial      ImprovementSynthesisStatus = "PARTIAL"
	ImprovementInconclusive ImprovementSynthesisStatus = "INCONCLUSIVE"
)

type TransferMode string

const (
	TransferComponentReuse TransferMode = "COMPONENT_REUSE"
	TransferPatternAdapt   TransferMode = "PATTERN_ADAPTATION"
	TransferReferenceOnly  TransferMode = "REFERENCE_ONLY"
	TransferUnknown        TransferMode = "UNKNOWN"
)

type HypothesisDerivation string

const (
	DerivationRuleDerived HypothesisDerivation = "RULE_DERIVED"
	DerivationExternal    HypothesisDerivation = "EXTERNAL_AGENT_PROPOSAL"
)

type TargetEvidenceState string

const (
	TargetPresent TargetEvidenceState = "PRESENT"
	TargetAbsent  TargetEvidenceState = "ABSENT"
	TargetPartial TargetEvidenceState = "PARTIAL"
	TargetUnknown TargetEvidenceState = "UNKNOWN"
)

type ImprovementLicenseAssessment string

const (
	ImprovementLicenseCompatible   ImprovementLicenseAssessment = "COMPATIBLE"
	ImprovementLicenseIncompatible ImprovementLicenseAssessment = "INCOMPATIBLE"
	ImprovementLicenseUnrecognized ImprovementLicenseAssessment = "UNRECOGNIZED"
)

type CounterEvidenceSearchStatus string

const (
	CounterEvidenceFound           CounterEvidenceSearchStatus = "FOUND"
	CounterEvidenceNoneInScope     CounterEvidenceSearchStatus = "NONE_FOUND_IN_SCOPE"
	CounterEvidenceBudgetExhausted CounterEvidenceSearchStatus = "BUDGET_EXHAUSTED"
)

type SourceArtifactRef struct {
	RunID    RunID           `json:"run_id"`
	RunKind  RunKind         `json:"run_kind"`
	Artifact RunArtifactName `json:"artifact"`
	SHA256   string          `json:"sha256"`
}

type ImprovementProposal struct {
	ProposalID             string               `json:"proposal_id"`
	Title                  string               `json:"title"`
	TargetProperty         string               `json:"target_property"`
	ExpectedEffect         string               `json:"expected_effect"`
	MechanismRefs          []string             `json:"mechanism_refs"`
	ProductObservationRefs []string             `json:"product_observation_refs"`
	Assumptions            []string             `json:"assumptions"`
	Derivation             HypothesisDerivation `json:"derivation"`
	Rationale              string               `json:"rationale"`
}

type ValidationStep struct {
	Question        string `json:"question"`
	Observable      string `json:"observable"`
	SafeProcedure   string `json:"safe_procedure"`
	PassCriterion   string `json:"pass_criterion"`
	RejectCriterion string `json:"reject_criterion"`
	RelativeCost    string `json:"relative_cost"`
	MissingEvidence string `json:"missing_evidence"`
}

type CounterEvidenceSearch struct {
	Status     CounterEvidenceSearchStatus `json:"status"`
	ScopeRefs  []string                    `json:"scope_refs"`
	ResultRefs []string                    `json:"result_refs"`
}

type ImprovementHypothesis struct {
	ID                     string                       `json:"id"`
	Title                  string                       `json:"title"`
	TargetProperty         string                       `json:"target_property"`
	TargetState            TargetEvidenceState          `json:"target_state"`
	TargetGapRefs          []string                     `json:"target_gap_refs"`
	ExpectedEffect         string                       `json:"expected_effect"`
	MechanismRefs          []string                     `json:"mechanism_refs"`
	ProductObservationRefs []string                     `json:"product_observation_refs"`
	TransferMode           TransferMode                 `json:"transfer_mode"`
	License                ImprovementLicenseAssessment `json:"license"`
	Assumptions            []string                     `json:"assumptions"`
	Preconditions          []string                     `json:"preconditions"`
	IntegrationUnknowns    []string                     `json:"integration_unknowns"`
	Risks                  []string                     `json:"risks"`
	Tradeoffs              []string                     `json:"tradeoffs"`
	EvidenceRefs           []string                     `json:"evidence_refs"`
	CounterEvidenceRefs    []string                     `json:"counter_evidence_refs"`
	CounterEvidenceSearch  CounterEvidenceSearch        `json:"counterevidence_search"`
	ValidationPlan         []ValidationStep             `json:"validation_plan"`
	Derivation             HypothesisDerivation         `json:"derivation"`
	Status                 ImprovementHypothesisStatus  `json:"status"`
	AdoptionDisposition    string                       `json:"adoption_disposition"`
}

type ImprovementSynthesisPolicy struct {
	MaxHypotheses                int  `json:"max_hypotheses"`
	MaxMechanismsPerHypothesis   int  `json:"max_mechanisms_per_hypothesis"`
	RequireCounterEvidenceSearch bool `json:"require_counterevidence_search"`
}

type ImprovementSynthesisRequest struct {
	Idea            IdeaSpec                   `json:"idea"`
	SourceArtifacts []SourceArtifactRef        `json:"source_artifacts"`
	Proposals       []ImprovementProposal      `json:"proposals"`
	Policy          ImprovementSynthesisPolicy `json:"policy"`
}

type ImprovementSynthesisStartReceipt struct {
	SchemaVersion string                  `json:"schema_version"`
	RunID         RunID                   `json:"run_id"`
	RunKind       RunKind                 `json:"run_kind"`
	Lifecycle     RecommendationLifecycle `json:"lifecycle"`
	Version       uint64                  `json:"version"`
}

type ImprovementSynthesisStarter interface {
	StartImprovementSynthesis(context.Context, string, ImprovementSynthesisRequest) (ImprovementSynthesisStartReceipt, error)
}

type ImprovementSynthesisExecutor interface {
	ExecuteImprovementSynthesis(context.Context, RunID) error
}

type ArchitectureComparison struct {
	TargetProperty      string              `json:"target_property"`
	TargetState         TargetEvidenceState `json:"target_state"`
	TargetRef           string              `json:"target_ref"`
	MechanismID         string              `json:"mechanism_id"`
	SourceRunID         RunID               `json:"source_run_id"`
	Novelty             MechanismNovelty    `json:"novelty"`
	IntegrationUnknowns []string            `json:"integration_unknowns"`
}

type HypothesisGap struct {
	HypothesisID string `json:"hypothesis_id"`
	Code         string `json:"code"`
	SourceRef    string `json:"source_ref"`
}

type ValidationPlanEntry struct {
	HypothesisID string         `json:"hypothesis_id"`
	StepIndex    int            `json:"step_index"`
	Step         ValidationStep `json:"step"`
}

type SynthesisProvenance struct {
	SourceRunID  RunID                `json:"source_run_id,omitempty"`
	Artifact     RunArtifactName      `json:"artifact,omitempty"`
	SHA256       string               `json:"sha256,omitempty"`
	HypothesisID string               `json:"hypothesis_id,omitempty"`
	Derivation   HypothesisDerivation `json:"derivation,omitempty"`
	RuleIDs      []string             `json:"rule_ids"`
	ScopeRefs    []string             `json:"scope_refs"`
}
