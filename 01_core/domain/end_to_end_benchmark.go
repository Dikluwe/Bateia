/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/end_to_end_benchmark.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type BenchmarkCaseID string
type BenchmarkSplit string

const (
	SplitContract   BenchmarkSplit = "CONTRACT"
	SplitTrain      BenchmarkSplit = "TRAIN"
	SplitValidation BenchmarkSplit = "VALIDATION"
	SplitTest       BenchmarkSplit = "TEST"
)

type DossierUtilityLabel string

const (
	UtilityUseful      DossierUtilityLabel = "USEFUL"
	UtilityNotUseful   DossierUtilityLabel = "NOT_USEFUL"
	UtilityUnjudgeable DossierUtilityLabel = "UNJUDGEABLE"
)

type CorrespondenceVerdict string

const (
	CorrespondencePreserved CorrespondenceVerdict = "PRESERVED"
	CorrespondenceViolated  CorrespondenceVerdict = "VIOLATED"
	CorrespondenceUnknown   CorrespondenceVerdict = "UNKNOWN"
)

type ReuseDecisionJudgment struct {
	Expected        ReuseDecision   `json:"expected"`
	Acceptable      []ReuseDecision `json:"acceptable,omitempty"`
	Primary         string          `json:"primary,omitempty"`
	Components      []string        `json:"components,omitempty"`
	RequiredRules   []string        `json:"required_rules"`
	RequiredSources []string        `json:"required_sources"`
	Forbidden       []ReuseDecision `json:"forbidden"`
	Rationale       string          `json:"rationale"`
}
type RecoveryJudgment struct {
	RecoverableCauses  []InconclusiveCause `json:"recoverable_causes"`
	TerminalCauses     []InconclusiveCause `json:"terminal_causes"`
	ExpectedNewKinds   []EvidenceKind      `json:"expected_new_kinds"`
	ForbiddenTargets   []string            `json:"forbidden_targets"`
	MaxRounds          int                 `json:"max_rounds"`
	MaxAdditionalBytes int64               `json:"max_additional_bytes"`
	ExpectedFinal      ReuseDecision       `json:"expected_final"`
}
type DossierUtilityJudgment struct {
	Label              DossierUtilityLabel `json:"label"`
	DecisionActionable bool                `json:"decision_actionable"`
	EvidenceSufficient bool                `json:"evidence_sufficient"`
	BlockersVisible    bool                `json:"blockers_visible"`
	RepositoriesOpened int                 `json:"repositories_opened"`
	Rationale          string              `json:"rationale"`
}
type BenchmarkProvenance struct {
	CapturedAt       string            `json:"captured_at"`
	SourceMode       IdeaRunMode       `json:"source_mode"`
	Provider         string            `json:"provider,omitempty"`
	OriginalRunID    string            `json:"original_run_id"`
	OriginalRunHash  string            `json:"original_run_hash"`
	ArtifactHashes   map[string]string `json:"artifact_hashes"`
	CandidateCommits map[string]string `json:"candidate_commits"`
	CollectorVersion string            `json:"collector_version"`
	Notes            []string          `json:"notes,omitempty"`
}
type EndToEndBenchmarkCase struct {
	SchemaVersion    string                 `json:"schema_version"`
	ID               BenchmarkCaseID        `json:"id"`
	Idea             IdeaSpec               `json:"idea"`
	Discovery        DiscoverySnapshot      `json:"discovery"`
	ScreeningModel   ScreeningModel         `json:"screening_model"`
	Bundles          []EvidenceBundle       `json:"bundles"`
	EvidenceRules    CapabilityRuleCatalog  `json:"evidence_rules"`
	ReusePolicy      ReuseDecisionPolicy    `json:"reuse_policy"`
	RecoveryPolicy   RecoveryPolicy         `json:"recovery_policy"`
	ExpectedEvidence []EvidenceJudgment     `json:"expected_evidence"`
	ExpectedDecision ReuseDecisionJudgment  `json:"expected_decision"`
	ExpectedRecovery RecoveryJudgment       `json:"expected_recovery"`
	UtilityJudgment  DossierUtilityJudgment `json:"utility_judgment"`
	Provenance       BenchmarkProvenance    `json:"provenance"`
	Split            BenchmarkSplit         `json:"split"`
	Group            string                 `json:"group"`
}
type RoleCapabilityReceipt struct {
	Role         string            `json:"role"`
	Executor     string            `json:"executor"`
	Readable     []string          `json:"readable"`
	Writable     []string          `json:"writable"`
	InputHashes  map[string]string `json:"input_hashes"`
	OutputHashes map[string]string `json:"output_hashes"`
	Limitation   string            `json:"limitation,omitempty"`
}
type BenchmarkManifest struct {
	ProtocolVersion   string                     `json:"protocol_version"`
	PromptPath        string                     `json:"prompt_path"`
	PromptHash        string                     `json:"prompt_hash"`
	BaselineCommit    string                     `json:"baseline_commit"`
	ContractPath      string                     `json:"contract_path"`
	ContractHash      string                     `json:"contract_hash"`
	UnknownPolicyHash string                     `json:"unknown_policy_hash"`
	CaseHashes        map[BenchmarkCaseID]string `json:"case_hashes"`
	SplitGroups       map[BenchmarkCaseID]string `json:"split_groups"`
	RoleCapabilities  []RoleCapabilityReceipt    `json:"role_capabilities"`
}
type BenchmarkCorpus struct {
	SchemaVersion string                  `json:"schema_version"`
	Cases         []EndToEndBenchmarkCase `json:"cases"`
}
type CorrespondenceResult struct {
	CaseID    BenchmarkCaseID       `json:"case_id"`
	Verdict   CorrespondenceVerdict `json:"verdict"`
	Witnesses []string              `json:"witnesses"`
}
type BenchmarkValidationReport struct {
	SchemaVersion string                 `json:"schema_version"`
	ManifestHash  string                 `json:"manifest_hash"`
	Results       []CorrespondenceResult `json:"results"`
	Preserved     int                    `json:"preserved"`
	Violated      int                    `json:"violated"`
	Unknown       int                    `json:"unknown"`
	Valid         bool                   `json:"valid"`
}
type BenchmarkMutation struct {
	ID   string                `json:"id"`
	Case EndToEndBenchmarkCase `json:"case"`
}
type BenchmarkMutationReport struct {
	SchemaVersion string   `json:"schema_version"`
	Total         int      `json:"total"`
	Rejected      int      `json:"rejected"`
	MutationScore *float64 `json:"mutation_score"`
	Survivors     []string `json:"survivors"`
	Sealed        bool     `json:"sealed"`
}
type BenchmarkCaseResult struct {
	CaseID              BenchmarkCaseID       `json:"case_id"`
	Correspondence      CorrespondenceVerdict `json:"correspondence"`
	EvidenceMatches     int                   `json:"evidence_matches"`
	EvidenceDivergences int                   `json:"evidence_divergences"`
	DecisionObserved    ReuseDecision         `json:"decision_observed"`
	DecisionAccepted    bool                  `json:"decision_accepted"`
	ForbiddenDecision   bool                  `json:"forbidden_decision"`
	RecoveryExecuted    bool                  `json:"recovery_executed"`
	RecoveryFinal       ReuseDecision         `json:"recovery_final,omitempty"`
	Utility             DossierUtilityLabel   `json:"utility"`
	Witnesses           []string              `json:"witnesses"`
}
type EndToEndBenchmarkMetrics struct {
	CorrespondencePreservedRate float64  `json:"correspondence_preserved_rate"`
	CorrespondenceUnknownRate   float64  `json:"correspondence_unknown_rate"`
	EvidencePrecision           *float64 `json:"evidence_precision"`
	EvidenceRecall              *float64 `json:"evidence_recall"`
	DecisionAccuracy            *float64 `json:"decision_accuracy"`
	ForbiddenDecisionRate       float64  `json:"forbidden_decision_rate"`
	ConclusiveCoverage          float64  `json:"conclusive_coverage"`
	CorrectInconclusiveRate     *float64 `json:"correct_inconclusive_rate"`
	RecoveryResolutionRate      *float64 `json:"recovery_resolution_rate"`
	RecoveryUsefulEvidenceRate  *float64 `json:"recovery_useful_evidence_rate"`
	UtilityRate                 *float64 `json:"utility_rate"`
}
type EndToEndBenchmarkReport struct {
	SchemaVersion string                   `json:"schema_version"`
	Split         BenchmarkSplit           `json:"split"`
	Cases         []BenchmarkCaseResult    `json:"cases"`
	Metrics       EndToEndBenchmarkMetrics `json:"metrics"`
	Passed        bool                     `json:"passed"`
}
