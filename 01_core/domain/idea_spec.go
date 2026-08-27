/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_spec.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

type IdeaSpec struct {
	Problem        string   `json:"problem"`
	Users          []string `json:"users,omitempty"`
	Capabilities   []string `json:"capabilities"`
	Constraints    []string `json:"constraints,omitempty"`
	NonGoals       []string `json:"non_goals,omitempty"`
	PreferredStack []string `json:"preferred_stack,omitempty"`
}

type ExpansionFacet string

const (
	FacetProblem      ExpansionFacet = "problem"
	FacetCapability   ExpansionFacet = "capability"
	FacetArchitecture ExpansionFacet = "architecture"
	FacetConstraint   ExpansionFacet = "constraint"
)

type SearchExpansion struct {
	ID              string               `json:"id"`
	Facet           ExpansionFacet       `json:"facet"`
	Query           string               `json:"query"`
	Discriminators  []string             `json:"discriminators,omitempty"`
	RepositoryScope string               `json:"repository_scope,omitempty"`
	Rationale       string               `json:"rationale"`
	Covers          []string             `json:"covers"`
	Derivation      *ExpansionDerivation `json:"derivation,omitempty"`
}

type ExpansionDerivation struct {
	Strategy  string   `json:"strategy"`
	RuleIDs   []string `json:"rule_ids"`
	SourceIDs []string `json:"source_ids"`
}

type ExpansionPlan struct {
	Idea       IdeaSpec          `json:"idea"`
	Expansions []SearchExpansion `json:"expansions"`
}

type EvaluationCase struct {
	ID                string   `json:"id"`
	Idea              IdeaSpec `json:"idea"`
	ExpectedConcepts  []string `json:"expected_concepts"`
	ForbiddenConcepts []string `json:"forbidden_concepts"`
}

type EvaluationCorpus struct {
	Cases []EvaluationCase `json:"cases"`
}

type ExpansionMetrics struct {
	Coverage      float64 `json:"coverage"`
	Contamination float64 `json:"contamination"`
	Novelty       float64 `json:"novelty"`
	Duplication   float64 `json:"duplication"`
}

type CaseEvaluation struct {
	ID                     string           `json:"id"`
	Metrics                ExpansionMetrics `json:"metrics"`
	FoundExpected          []string         `json:"found_expected"`
	FoundForbidden         []string         `json:"found_forbidden"`
	ExpansionCount         int              `json:"expansion_count"`
	NewExpansionCount      int              `json:"new_expansion_count"`
	ProposedExpansionCount int              `json:"proposed_expansion_count"`
}

type EvaluationReport struct {
	Strategy string           `json:"strategy"`
	Metrics  ExpansionMetrics `json:"metrics"`
	Cases    []CaseEvaluation `json:"cases"`
}
