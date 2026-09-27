/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation_intake.md
 * @layer L1
 * @updated 2026-09-24
 */

package domain

import "encoding/json"

const RecommendationIntakeSchemaVersion = "1"

type InputOrigin string

const (
	InputExplicit            InputOrigin = "EXPLICIT"
	InputVersionedDefault    InputOrigin = "VERSIONED_DEFAULT"
	InputDerived             InputOrigin = "DERIVED"
	InputReplayed            InputOrigin = "REPLAYED"
	InputEnvironmentObserved InputOrigin = "ENVIRONMENT_OBSERVED"
	InputUnresolved          InputOrigin = "UNRESOLVED"
)

type InputAuthority string

const (
	AuthorityUnresolved       InputAuthority = "UNRESOLVED"
	AuthorityUser             InputAuthority = "USER"
	AuthorityProductPolicy    InputAuthority = "PRODUCT_POLICY"
	AuthorityVerifiedArtifact InputAuthority = "VERIFIED_ARTIFACT"
	AuthorityRuntime          InputAuthority = "RUNTIME"
	AuthorityAIProposal       InputAuthority = "AI_PROPOSAL"
)

type IntakeStatus string

const (
	IntakeReady      IntakeStatus = "READY"
	IntakeNeedsInput IntakeStatus = "NEEDS_INPUT"
	IntakeInvalid    IntakeStatus = "INVALID"
)

type NeedPhase string

const (
	NeedInitial   NeedPhase = "INITIAL"
	NeedExecution NeedPhase = "EXECUTION"
)

type ResolverKind string

const (
	ResolverCompute              ResolverKind = "COMPUTE"
	ResolverCollect              ResolverKind = "COLLECT"
	ResolverAskUser              ResolverKind = "ASK_USER"
	ResolverRequestAuthorization ResolverKind = "REQUEST_AUTHORIZATION"
	ResolverRequireRiskOwner     ResolverKind = "REQUIRE_RISK_OWNER"
	ResolverTerminal             ResolverKind = "TERMINAL"
)

type AnswerKind string

const (
	AnswerText       AnswerKind = "TEXT"
	AnswerStringList AnswerKind = "STRING_LIST"
	AnswerBoolean    AnswerKind = "BOOLEAN"
)

// Pointer fields preserve omission separately from explicit false, zero,
// empty string and empty list.
type RecommendationIntakeDraft struct {
	Idea      *IdeaIntakeInput      `json:"idea,omitempty"`
	Policy    *PolicyIntakeInput    `json:"policy,omitempty"`
	Target    *TargetIntakeInput    `json:"target,omitempty"`
	Search    *SearchIntakeInput    `json:"search,omitempty"`
	Execution *ExecutionIntakeInput `json:"execution,omitempty"`
}

type IdeaIntakeInput struct {
	Problem             *string   `json:"problem,omitempty"`
	Users               *[]string `json:"users,omitempty"`
	Capabilities        *[]string `json:"capabilities,omitempty"`
	Qualities           *[]string `json:"qualities,omitempty"`
	OperatingConditions *[]string `json:"operating_conditions,omitempty"`
	Constraints         *[]string `json:"constraints,omitempty"`
	NonGoals            *[]string `json:"non_goals,omitempty"`
	PreferredStack      *[]string `json:"preferred_stack,omitempty"`
}

type PolicyIntakeInput struct {
	SchemaVersion       *string   `json:"schema_version,omitempty"`
	AllowedLicenses     *[]string `json:"allowed_licenses,omitempty"`
	RequireLocal        *bool     `json:"require_local,omitempty"`
	AllowArchived       *bool     `json:"allow_archived,omitempty"`
	MaxComponents       *int      `json:"max_components,omitempty"`
	ExtendMinCoverage   *float64  `json:"extend_min_coverage,omitempty"`
	IntegrationMinValue *float64  `json:"integration_min_value,omitempty"`
}

type TargetIntakeInput struct {
	Name *string `json:"name,omitempty"`
}

type SearchIntakeInput struct {
	Strategy      *string   `json:"strategy,omitempty"`
	Language      *string   `json:"language,omitempty"`
	Locales       *[]string `json:"locales,omitempty"`
	MaxPerQuery   *int      `json:"max_per_query,omitempty"`
	MaxCandidates *int      `json:"max_candidates,omitempty"`
}

type ExecutionIntakeInput struct {
	Mode                    *string   `json:"mode,omitempty"`
	Provider                *string   `json:"provider,omitempty"`
	LocalPaths              *[]string `json:"local_paths,omitempty"`
	ReplayDirectory         *string   `json:"replay_directory,omitempty"`
	OutputDirectory         *string   `json:"output_directory,omitempty"`
	Resume                  *bool     `json:"resume,omitempty"`
	TimeoutSeconds          *int      `json:"timeout_seconds,omitempty"`
	AllowCandidateExecution *bool     `json:"allow_candidate_execution,omitempty"`
	EvidenceRulesRef        *string   `json:"evidence_rules_ref,omitempty"`
	ScreeningModelRef       *string   `json:"screening_model_ref,omitempty"`
	IdempotencyKey          *string   `json:"idempotency_key,omitempty"`
}

type RecommendationExecutionProfile struct {
	Mode                    string   `json:"mode"`
	Provider                string   `json:"provider"`
	LocalPaths              []string `json:"local_paths"`
	ReplayDirectory         string   `json:"replay_directory,omitempty"`
	OutputDirectory         string   `json:"output_directory,omitempty"`
	Resume                  bool     `json:"resume"`
	TimeoutSeconds          int      `json:"timeout_seconds"`
	AllowCandidateExecution bool     `json:"allow_candidate_execution"`
	EvidenceRulesRef        string   `json:"evidence_rules_ref,omitempty"`
	ScreeningModelRef       string   `json:"screening_model_ref,omitempty"`
	IdempotencyKey          string   `json:"idempotency_key,omitempty"`
}

type InputProvenance struct {
	Path      string         `json:"path"`
	Origin    InputOrigin    `json:"origin"`
	Authority InputAuthority `json:"authority"`
	SourceRef string         `json:"source_ref,omitempty"`
	Digest    string         `json:"digest,omitempty"`
}

type InputIssue struct {
	Path     string `json:"path"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Blocking bool   `json:"blocking"`
}

type DecisionEffect struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

type InformationNeed struct {
	ID              string           `json:"id"`
	Phase           NeedPhase        `json:"phase"`
	Subject         string           `json:"subject"`
	TargetPath      string           `json:"target_path,omitempty"`
	ClaimID         string           `json:"claim_id,omitempty"`
	ReasonCode      string           `json:"reason_code"`
	Resolver        ResolverKind     `json:"resolver"`
	AnswerKind      AnswerKind       `json:"answer_kind,omitempty"`
	AllowedValues   []string         `json:"allowed_values,omitempty"`
	Unit            string           `json:"unit,omitempty"`
	ValidationRule  string           `json:"validation_rule,omitempty"`
	DecisionEffects []DecisionEffect `json:"decision_effects"`
	Blocking        bool             `json:"blocking"`
	Priority        int              `json:"priority"`
}

type InformationAnswer struct {
	NeedID     string          `json:"need_id"`
	AnswerKind AnswerKind      `json:"answer_kind"`
	Value      json.RawMessage `json:"value"`
	Authority  InputAuthority  `json:"authority"`
	SourceRef  string          `json:"source_ref,omitempty"`
	Confirmed  bool            `json:"confirmed"`
	Supersedes string          `json:"supersedes,omitempty"`
}

type RecommendationIntake struct {
	SchemaVersion string                         `json:"schema_version"`
	Request       RecommendationRequest          `json:"request"`
	Execution     RecommendationExecutionProfile `json:"execution"`
	Provenance    []InputProvenance              `json:"provenance"`
	Issues        []InputIssue                   `json:"issues"`
	Questions     []InformationNeed              `json:"questions"`
	Answers       []InformationAnswer            `json:"answers"`
	Status        IntakeStatus                   `json:"status"`
	InputHash     string                         `json:"input_hash"`
}

type InformationNeedReport struct {
	SchemaVersion string            `json:"schema_version"`
	IntakeHash    string            `json:"intake_hash"`
	Needs         []InformationNeed `json:"needs"`
}
