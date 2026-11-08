/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_possibilities.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

// Operational coverage is an opportunity to discover, not a measure of insight quality.
type InsightOpportunityCoverage struct {
	AttemptedPerspectives []string `json:"attempted_perspectives"`
	CompletedPerspectives []string `json:"completed_perspectives"`
	AttemptedSources      []string `json:"attempted_sources"`
}

// A lead witnesses a reference in a document; the target has not been read.
type InsightReferenceLead struct {
	ID                   string                  `json:"id"`
	SourceCandidateID    string                  `json:"source_candidate_id"`
	SourceIdentity       ProductIdentity         `json:"source_identity"`
	Source               CommitRef               `json:"source"`
	Reference            LocatedInsightReference `json:"reference"`
	Target               ProductIdentity         `json:"target"`
	Status               string                  `json:"status"`
	MatchedQuestionTerms []string                `json:"matched_question_terms"`
	SelectionReasons     []string                `json:"selection_reasons"`
}
