/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/recommendation.md
 * @layer L1
 * @updated 2026-09-01
 */

package domain

type RecommendationSearch struct {
	Strategy      string   `json:"strategy"`
	Language      string   `json:"language,omitempty"`
	Locales       []string `json:"locales"`
	MaxPerQuery   int      `json:"max_per_query"`
	MaxCandidates int      `json:"max_candidates"`
}

type RecommendationRequest struct {
	Idea   IdeaSpec                  `json:"idea"`
	Policy ReuseDecisionPolicy       `json:"policy"`
	Target SoftwareAlternativeTarget `json:"target"`
	Search RecommendationSearch      `json:"search"`
}

type RecommendationSummary struct {
	ExpansionCount int `json:"expansion_count"`
	CandidateCount int `json:"candidate_count"`
	AssessedCount  int `json:"assessed_count"`
}

type RecommendationResult struct {
	SchemaVersion string                    `json:"schema_version"`
	Status        IdeaRunStatus             `json:"status"`
	Summary       RecommendationSummary     `json:"summary"`
	Decision      ReuseDecisionReport       `json:"decision"`
	Alternatives  SoftwareAlternativeReport `json:"alternatives"`
	Failures      []ExecutionFailure        `json:"failures"`
	Warnings      []string                  `json:"warnings"`
}
