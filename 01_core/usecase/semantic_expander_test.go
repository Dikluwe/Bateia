/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/semantic_expander_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

type semanticExpanderFunc func(domain.ExpansionPlan) ([]domain.SearchExpansion, error)

func (function semanticExpanderFunc) Expand(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
	return function(plan)
}

func proposal(source domain.SearchExpansion, query string) domain.SearchExpansion {
	return domain.SearchExpansion{
		Facet: source.Facet, Query: query, Rationale: "test rule", Covers: cloneStrings(source.Covers),
		Derivation: &domain.ExpansionDerivation{Strategy: "lexical", RuleIDs: []string{"rule-1"}, SourceIDs: []string{source.ID}},
	}
}

func TestEnrichExpansionPlanPreservesLiteralAndInput(t *testing.T) {
	literal, err := ExpandIdea(domain.IdeaSpec{Problem: "impacto estrutural"})
	if err != nil {
		t.Fatal(err)
	}
	original, _ := json.Marshal(literal)
	result, err := EnrichExpansionPlan(*literal, semanticExpanderFunc(func(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
		return []domain.SearchExpansion{proposal(plan.Expansions[0], "change blast radius")}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(literal)
	if string(original) != string(after) {
		t.Fatalf("plano de entrada foi modificado")
	}
	if !reflect.DeepEqual(result.Plan.Expansions[0], literal.Expansions[0]) {
		t.Fatalf("expansão literal mudou")
	}
	if len(result.Plan.Expansions) != 2 || result.Plan.Expansions[1].Derivation == nil {
		t.Fatalf("enriquecimento ausente: %+v", result)
	}
}

func TestEnrichExpansionPlanDeduplicatesAndMergesTraceability(t *testing.T) {
	literal, _ := ExpandIdea(domain.IdeaSpec{Problem: "impact analysis"})
	result, err := EnrichExpansionPlan(*literal, semanticExpanderFunc(func(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
		first := proposal(plan.Expansions[0], " Impact   Analysis ")
		first.Covers = []string{"capabilities[0]"}
		return []domain.SearchExpansion{first}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Plan.Expansions) != 1 || result.DuplicatesRemoved != 1 {
		t.Fatalf("duplicata não removida: %+v", result)
	}
	got := result.Plan.Expansions[0]
	if len(got.Covers) != 2 || got.Derivation == nil || len(got.Derivation.SourceIDs) != 1 {
		t.Fatalf("rastreabilidade não mesclada: %+v", got)
	}
}

func TestEnrichExpansionPlanRejectsMissingProvenance(t *testing.T) {
	literal, _ := ExpandIdea(domain.IdeaSpec{Problem: "impact analysis"})
	_, err := EnrichExpansionPlan(*literal, semanticExpanderFunc(func(domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
		return []domain.SearchExpansion{{Query: "blast radius"}}, nil
	}))
	if err == nil {
		t.Fatal("expansão sem proveniência foi aceita")
	}
}

func TestEnrichExpansionPlanUsesNonGoalsAsBarrier(t *testing.T) {
	literal, _ := ExpandIdea(domain.IdeaSpec{Problem: "impact analysis", NonGoals: []string{"runtime monitoring"}})
	result, err := EnrichExpansionPlan(*literal, semanticExpanderFunc(func(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
		return []domain.SearchExpansion{proposal(plan.Expansions[0], "runtime monitoring platform")}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Plan.Expansions) != 1 {
		t.Fatalf("non-goal promovido: %+v", result.Plan.Expansions)
	}
}

func TestEvaluateExpansionStrategyIsDeterministicAndHandlesEmptyLists(t *testing.T) {
	corpus := domain.EvaluationCorpus{Cases: []domain.EvaluationCase{{ID: "one", Idea: domain.IdeaSpec{Problem: "impact analysis"}}}}
	first, err := EvaluateExpansionStrategy(corpus, "literal", nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := EvaluateExpansionStrategy(corpus, "literal", nil)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("avaliação não determinística")
	}
	if first.Metrics.Coverage != 0 || first.Metrics.Contamination != 0 || first.Metrics.Duplication != 0 {
		t.Fatalf("semântica de lista vazia inválida: %+v", first.Metrics)
	}
}

func TestEvaluateExpansionStrategyMeasuresCoverageContaminationAndNovelty(t *testing.T) {
	corpus := domain.EvaluationCorpus{Cases: []domain.EvaluationCase{{
		ID: "one", Idea: domain.IdeaSpec{Problem: "impact analysis"},
		ExpectedConcepts: []string{"blast radius"}, ForbiddenConcepts: []string{"runtime monitoring"},
	}}}
	report, err := EvaluateExpansionStrategy(corpus, "lexical", semanticExpanderFunc(func(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
		return []domain.SearchExpansion{proposal(plan.Expansions[0], "blast radius")}, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if report.Metrics.Coverage != 1 || report.Metrics.Contamination != 0 || report.Metrics.Novelty != 0.5 {
		t.Fatalf("métricas incorretas: %+v", report.Metrics)
	}
}
