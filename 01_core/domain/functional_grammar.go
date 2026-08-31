/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/functional_grammar.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type FunctionalGrammar struct {
	SchemaVersion string                  `json:"schema_version"`
	Version       string                  `json:"version"`
	Rules         []FunctionalGrammarRule `json:"rules"`
	Morphology    []MorphologyRule        `json:"morphology"`
	Hash          string                  `json:"hash"`
}

type FunctionalGrammarRule struct {
	ID            string   `json:"id"`
	Intent        string   `json:"intent"`
	Pattern       string   `json:"pattern"`
	RequiredSlots []string `json:"required_slots"`
	OptionalSlots []string `json:"optional_slots,omitempty"`
	Priority      int      `json:"priority"`
	MaxVariants   int      `json:"max_variants"`
	SubjectKind   string   `json:"subject_kind"`
}

type MorphologyRule struct {
	ID          string `json:"id"`
	SubjectKind string `json:"subject_kind"`
	Source      string `json:"source"`
	Realized    string `json:"realized"`
}

type GrammarFinding struct {
	QueryID         string   `json:"query_id"`
	GrammarVersion  string   `json:"grammar_version"`
	RuleID          string   `json:"rule_id"`
	Verdict         string   `json:"verdict"`
	UnresolvedVerbs []string `json:"unresolved_verbs,omitempty"`
	MissingSlots    []string `json:"missing_slots,omitempty"`
	ForbiddenHits   []string `json:"forbidden_hits,omitempty"`
	Reasons         []string `json:"reasons"`
}

type RealizedFunctionalQuery struct {
	FunctionalQuery
	GrammarVersion string   `json:"grammar_version"`
	GrammarHash    string   `json:"grammar_hash"`
	GrammarRuleID  string   `json:"grammar_rule_id"`
	MorphologyIDs  []string `json:"morphology_ids,omitempty"`
}

type FunctionalGrammarExpansion struct {
	SchemaVersion string                    `json:"schema_version"`
	GrammarHash   string                    `json:"grammar_hash"`
	Queries       []RealizedFunctionalQuery `json:"queries"`
	Unknown       []GrammarFinding          `json:"unknown,omitempty"`
}
