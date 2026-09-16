/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_discovery.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain

import "time"

type DiscoveryRequest struct {
	RequirementIDs     []string             `json:"requirement_ids,omitempty"`
	Derivation         *ExpansionDerivation `json:"derivation,omitempty"`
	ExpansionID        string               `json:"expansion_id"`
	Query              string               `json:"query"`
	Discriminators     []string             `json:"discriminators,omitempty"`
	RepositoryScope    string               `json:"repository_scope,omitempty"`
	Language           string               `json:"language,omitempty"`
	SearchLocale       string               `json:"search_locale,omitempty"`
	RegionalProposalID string               `json:"regional_proposal_id,omitempty"`
	MaxResults         int                  `json:"max_results"`
	AttemptTimeout     time.Duration        `json:"attempt_timeout,omitempty"`
	MaxAttempts        int                  `json:"max_attempts,omitempty"`
}

type ProviderAttempt struct {
	Coverage             *DiscoveryCoverage             `json:"coverage,omitempty"`
	CandidateDiagnostics []DiscoveryCandidateDiagnostic `json:"candidate_diagnostics,omitempty"`
	QueryPlan            *DiscoveryQueryPlan            `json:"query_plan,omitempty"`
	RawCandidates        int                            `json:"raw_candidates,omitempty"`
	RejectedCandidates   int                            `json:"rejected_candidates,omitempty"`
	Query                string                         `json:"query"`
	AttemptIndex         int                            `json:"attempt_index,omitempty"`
	StartedAt            time.Time                      `json:"started_at,omitempty"`
	Duration             time.Duration                  `json:"duration,omitempty"`
	Status               string                         `json:"status,omitempty"`
	ResultCount          int                            `json:"result_count"`
	ErrorClass           string                         `json:"error_class,omitempty"`
	Error                string                         `json:"error,omitempty"`
	SearchLocale         string                         `json:"search_locale,omitempty"`
	RegionalProposalID   string                         `json:"regional_proposal_id,omitempty"`
}

type ProviderDiscovery struct {
	Provider   string                   `json:"provider"`
	Attempts   []ProviderAttempt        `json:"attempts"`
	Candidates []RepositoryCandidate    `json:"candidates"`
	Matches    []ProviderCandidateMatch `json:"matches,omitempty"`
}

type DiscoveredCandidate struct {
	Repository   RepositoryCandidate `json:"repository"`
	ExpansionIDs []string            `json:"expansion_ids"`
	Queries      []string            `json:"queries"`
	Origins      []CandidateOrigin   `json:"origins,omitempty"`
}

type ExpansionDiscovery struct {
	Expansion SearchExpansion   `json:"expansion"`
	Provider  ProviderDiscovery `json:"provider"`
}

type DiscoverySnapshot struct {
	SchemaVersion string                `json:"schema_version"`
	Idea          IdeaSpec              `json:"idea"`
	Strategy      string                `json:"strategy"`
	Language      string                `json:"language,omitempty"`
	SearchLocales []string              `json:"search_locales,omitempty"`
	PerExpansion  []ExpansionDiscovery  `json:"per_expansion"`
	Candidates    []DiscoveredCandidate `json:"candidates"`
	Failures      []ExecutionFailure    `json:"failures,omitempty"`
	Receipt       *DiscoveryRunReceipt  `json:"receipt,omitempty"`
}

type DiscoveryRunStatus string

const (
	DiscoveryRunComplete DiscoveryRunStatus = "COMPLETE"
	DiscoveryRunPartial  DiscoveryRunStatus = "PARTIAL"
	DiscoveryRunFailed   DiscoveryRunStatus = "FAILED"
)

type DiscoveryCompletionCause string

const (
	CompletionFinished       DiscoveryCompletionCause = "FINISHED"
	CompletionTotalTimeout   DiscoveryCompletionCause = "TOTAL_TIMEOUT"
	CompletionCancelled      DiscoveryCompletionCause = "CANCELLED"
	CompletionProviderFailed DiscoveryCompletionCause = "PROVIDER_FAILED"
	CompletionInvalidInput   DiscoveryCompletionCause = "INVALID_INPUT"
)

type DiscoveryExecutionPolicy struct {
	SchemaVersion  string        `json:"schema_version"`
	TotalTimeout   time.Duration `json:"total_timeout"`
	AttemptTimeout time.Duration `json:"attempt_timeout"`
	Concurrency    int           `json:"concurrency"`
	MaxAttempts    int           `json:"max_attempts"`
}
type ProviderAttemptReceipt struct {
	ExpansionID        string        `json:"expansion_id"`
	AttemptIndex       int           `json:"attempt_index"`
	PlannedQuery       string        `json:"planned_query"`
	EffectiveQuery     string        `json:"effective_query"`
	StartedAt          time.Time     `json:"started_at,omitempty"`
	Duration           time.Duration `json:"duration,omitempty"`
	Status             string        `json:"status"`
	ResultCount        int           `json:"result_count"`
	ErrorClass         string        `json:"error_class,omitempty"`
	SafeMessage        string        `json:"safe_message,omitempty"`
	SearchLocale       string        `json:"search_locale,omitempty"`
	RegionalProposalID string        `json:"regional_proposal_id,omitempty"`
}
type DiscoveryRunReceipt struct {
	SchemaVersion   string                   `json:"schema_version"`
	Status          DiscoveryRunStatus       `json:"status"`
	CompletionCause DiscoveryCompletionCause `json:"completion_cause"`
	Policy          DiscoveryExecutionPolicy `json:"policy"`
	Planned         int                      `json:"planned"`
	Started         int                      `json:"started"`
	Completed       int                      `json:"completed"`
	TimedOut        int                      `json:"timed_out"`
	Cancelled       int                      `json:"cancelled"`
	Failed          int                      `json:"failed"`
	Attempts        []ProviderAttemptReceipt `json:"attempts"`
	Warnings        []string                 `json:"warnings,omitempty"`
}
type DiscoveryProgressEvent struct {
	ExpansionID  string `json:"expansion_id"`
	AttemptIndex int    `json:"attempt_index"`
	Status       string `json:"status"`
	Completed    int    `json:"completed"`
	Planned      int    `json:"planned"`
}
type DiscoveryProgressReporter interface{ ReportDiscoveryProgress(DiscoveryProgressEvent) }

type RelevanceLabel string

const (
	RelevanceRelevant   RelevanceLabel = "relevant"
	RelevanceIrrelevant RelevanceLabel = "irrelevant"
	RelevanceUncertain  RelevanceLabel = "uncertain"
)

type CandidateJudgment struct {
	Repository string         `json:"repository"`
	Label      RelevanceLabel `json:"label"`
	Reason     string         `json:"reason"`
}

type DiscoveryJudgments struct {
	SchemaVersion string              `json:"schema_version"`
	Judgments     []CandidateJudgment `json:"judgments"`
}

type DiscoveryMetrics struct {
	JudgedPrecision  *float64 `json:"judged_precision"`
	UncertainRate    float64  `json:"uncertain_rate"`
	UnjudgedRate     float64  `json:"unjudged_rate"`
	ExpansionHitRate float64  `json:"expansion_hit_rate"`
	UniqueYield      float64  `json:"unique_yield"`
	OverlapRate      float64  `json:"overlap_rate"`
	FailureRate      float64  `json:"failure_rate"`
	RelevantCount    int      `json:"relevant_count"`
	IrrelevantCount  int      `json:"irrelevant_count"`
	UncertainCount   int      `json:"uncertain_count"`
	UnjudgedCount    int      `json:"unjudged_count"`
}

type ExpansionDiscoveryEvaluation struct {
	ExpansionID    string `json:"expansion_id"`
	CandidateCount int    `json:"candidate_count"`
	RelevantCount  int    `json:"relevant_count"`
	Hit            bool   `json:"hit"`
}

type DiscoveryEvaluation struct {
	SchemaVersion string                         `json:"schema_version"`
	Strategy      string                         `json:"strategy"`
	Metrics       DiscoveryMetrics               `json:"metrics"`
	PerExpansion  []ExpansionDiscoveryEvaluation `json:"per_expansion"`
}
