/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_discovery.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

type DiscoveryRequest struct {
	ExpansionID string `json:"expansion_id"`
	Query       string `json:"query"`
	Language    string `json:"language,omitempty"`
	MaxResults  int    `json:"max_results"`
}

type ProviderAttempt struct {
	Query       string `json:"query"`
	ResultCount int    `json:"result_count"`
	Error       string `json:"error,omitempty"`
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
	PerExpansion  []ExpansionDiscovery  `json:"per_expansion"`
	Candidates    []DiscoveredCandidate `json:"candidates"`
	Failures      []ExecutionFailure    `json:"failures,omitempty"`
}

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
