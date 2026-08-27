/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_recovery.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type InconclusiveCause string

const (
	CauseUnresolvedRevision    InconclusiveCause = "UNRESOLVED_REVISION"
	CauseIncompleteCollection  InconclusiveCause = "INCOMPLETE_COLLECTION"
	CauseMissingImplementation InconclusiveCause = "MISSING_IMPLEMENTATION_EVIDENCE"
	CauseMissingTest           InconclusiveCause = "MISSING_TEST_EVIDENCE"
	CauseUnknownLicense        InconclusiveCause = "UNKNOWN_LICENSE"
	CauseUnknownOperation      InconclusiveCause = "UNKNOWN_OPERATION"
	CauseConflictingEvidence   InconclusiveCause = "CONFLICTING_EVIDENCE"
	CauseBudgetExhausted       InconclusiveCause = "BUDGET_EXHAUSTED"
	CauseProviderFailure       InconclusiveCause = "PROVIDER_FAILURE"
	CauseUnsupportedArtifact   InconclusiveCause = "UNSUPPORTED_ARTIFACT"
	CauseRuleCoverageGap       InconclusiveCause = "RULE_COVERAGE_GAP"
	CauseInsufficientCorpus    InconclusiveCause = "INSUFFICIENT_CORPUS"
)

type InconclusiveFinding struct {
	ID          string             `json:"id"`
	Cause       InconclusiveCause  `json:"cause"`
	Candidate   string             `json:"candidate,omitempty"`
	Capability  string             `json:"capability,omitempty"`
	RuleID      string             `json:"rule_id"`
	Evidence    []EvidenceSource   `json:"evidence,omitempty"`
	Failures    []ExecutionFailure `json:"failures,omitempty"`
	Missing     []string           `json:"missing"`
	Recoverable bool               `json:"recoverable"`
	Message     string             `json:"message"`
}
type InconclusiveDiagnosis struct {
	SchemaVersion string                `json:"schema_version"`
	RunID         string                `json:"run_id"`
	DecisionHash  string                `json:"decision_hash"`
	Findings      []InconclusiveFinding `json:"findings"`
	Recoverable   int                   `json:"recoverable"`
	Terminal      int                   `json:"terminal"`
}
type RecoveryPolicy struct {
	SchemaVersion        string `json:"schema_version"`
	MaxRounds            int    `json:"max_rounds"`
	MaxCandidates        int    `json:"max_candidates"`
	AdditionalDocuments  int    `json:"additional_documents"`
	AdditionalCodeFiles  int    `json:"additional_code_files"`
	AdditionalTests      int    `json:"additional_tests"`
	AdditionalBytes      int64  `json:"additional_bytes"`
	AdditionalTreeDepth  int    `json:"additional_tree_depth"`
	RetryProviderFailure bool   `json:"retry_provider_failure"`
}
type RecoveryTarget struct {
	FindingID  string           `json:"finding_id"`
	Candidate  CommitRef        `json:"candidate"`
	Capability string           `json:"capability,omitempty"`
	Path       string           `json:"path,omitempty"`
	Patterns   []string         `json:"patterns,omitempty"`
	Kind       EvidenceKind     `json:"kind"`
	RuleID     string           `json:"rule_id"`
	Expected   []EvidenceStatus `json:"expected"`
	Reason     string           `json:"reason"`
}
type EvidenceRecoveryPlan struct {
	SchemaVersion string                `json:"schema_version"`
	Round         int                   `json:"round"`
	DiagnosisHash string                `json:"diagnosis_hash"`
	Policy        RecoveryPolicy        `json:"policy"`
	PolicyHash    string                `json:"policy_hash"`
	Targets       []RecoveryTarget      `json:"targets"`
	Skipped       []InconclusiveFinding `json:"skipped"`
}
type RecoveryMetrics struct {
	Round                  int      `json:"round"`
	FindingsBefore         int      `json:"findings_before"`
	TargetsPlanned         int      `json:"targets_planned"`
	TargetsExecuted        int      `json:"targets_executed"`
	TargetsNew             int      `json:"targets_new"`
	DuplicateContent       int      `json:"duplicate_content"`
	AdditionalBytes        int64    `json:"additional_bytes"`
	ProviderCalls          int      `json:"provider_calls"`
	StatusTransitions      int      `json:"status_transitions"`
	ResolvedFindings       int      `json:"resolved_findings"`
	RemainingFindings      int      `json:"remaining_findings"`
	DecisionBefore         string   `json:"decision_before"`
	DecisionAfter          string   `json:"decision_after"`
	UsefulEvidenceRate     *float64 `json:"useful_evidence_rate"`
	InconclusiveResolution *float64 `json:"inconclusive_resolution_rate"`
}
type RecoveryRoundResult struct {
	Bundles    []EvidenceBundle      `json:"bundles"`
	Assessment AssessmentReport      `json:"assessment"`
	Decision   ReuseDecisionReport   `json:"decision"`
	Diagnosis  InconclusiveDiagnosis `json:"diagnosis"`
	Metrics    RecoveryMetrics       `json:"metrics"`
}
type RecoveryRoundIndex struct {
	Round          int    `json:"round"`
	DiagnosisHash  string `json:"diagnosis_hash"`
	PlanHash       string `json:"plan_hash"`
	BundlesHash    string `json:"bundles_hash"`
	AssessmentHash string `json:"assessment_hash"`
	DecisionHash   string `json:"decision_hash"`
	MetricsHash    string `json:"metrics_hash"`
}
type RecoveryIndex struct {
	SchemaVersion   string               `json:"schema_version"`
	OriginalRunHash string               `json:"original_run_hash"`
	Rounds          []RecoveryRoundIndex `json:"rounds"`
}
