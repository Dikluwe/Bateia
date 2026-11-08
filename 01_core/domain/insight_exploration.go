/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/insight_exploration.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type InsightExplorationRequest struct {
	AlgorithmIdeas     bool             `json:"algorithm_ideas,omitempty"`
	ReviewAlternatives bool             `json:"review_alternatives,omitempty"`
	SchemaVersion      string           `json:"schema_version"`
	Question           string           `json:"question"`
	FamiliarApproaches []string         `json:"familiar_approaches"`
	Bundles            []EvidenceBundle `json:"bundles"`
	MaxPassages        int              `json:"max_passages"`
}

type InsightReadingCard struct {
	ID                       string                  `json:"id"`
	Source                   ImmutableCandidateRef   `json:"source"`
	Kind                     EvidenceKind            `json:"kind"`
	Reference                LocatedInsightReference `json:"reference"`
	Origin                   string                  `json:"origin"`
	MatchedQuestionTerms     []string                `json:"matched_question_terms"`
	DistinctExcerptWords     int                     `json:"distinct_excerpt_words"`
	FamiliarSubstringMatches []string                `json:"familiar_substring_matches"`
	SelectionReasons         []string                `json:"selection_reasons"`
}

type InsightExplorationTask struct {
	Perspective string   `json:"perspective"`
	Question    string   `json:"question"`
	CardRefs    []string `json:"card_refs"`
	Gaps        []string `json:"gaps"`
}

type InsightResearchDirection struct {
	Perspective string `json:"perspective"`
	Query       string `json:"query"`
	Status      string `json:"status"`
}

type InsightExplorationDiagnostics struct {
	SuppliedSources    int  `json:"supplied_sources"`
	SelectedSources    int  `json:"selected_sources"`
	EligibleFiles      int  `json:"eligible_files"`
	IgnoredItems       int  `json:"ignored_items"`
	PrefixLimitedFiles int  `json:"prefix_limited_files"`
	DuplicatePassages  int  `json:"duplicate_passages"`
	SelectionLimited   bool `json:"selection_limited"`
	IncompleteBundles  int  `json:"incomplete_bundles"`
}

type InsightExplorationReport struct {
	AlgorithmAgenda    []string                      `json:"algorithm_agenda,omitempty"`
	SchemaVersion      string                        `json:"schema_version"`
	ID                 string                        `json:"id"`
	InputSHA256        string                        `json:"input_sha256"`
	Status             string                        `json:"status"`
	Question           string                        `json:"question"`
	FamiliarApproaches []string                      `json:"familiar_approaches"`
	QuestionTerms      []string                      `json:"question_terms"`
	Cards              []InsightReadingCard          `json:"cards"`
	Tasks              []InsightExplorationTask      `json:"tasks"`
	ResearchDirections []InsightResearchDirection    `json:"research_directions"`
	AIInstructions     []string                      `json:"ai_instructions"`
	Diagnostics        InsightExplorationDiagnostics `json:"diagnostics"`
	Gaps               []string                      `json:"gaps"`
	NoveltyState       string                        `json:"novelty_state"`
	EffectState        string                        `json:"effect_state"`
	Limitations        []string                      `json:"limitations"`
}

type InsightCollectionReport struct {
	SchemaVersion          string                   `json:"schema_version"`
	Bundles                []EvidenceBundle         `json:"bundles"`
	UnexploredRepositories []string                 `json:"unexplored_repositories"`
	Exploration            InsightExplorationReport `json:"exploration"`
}
