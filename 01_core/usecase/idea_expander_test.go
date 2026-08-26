/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_expander_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func completeIdeaFixture() domain.IdeaSpec {
	return domain.IdeaSpec{
		Problem: "  mostrar   consequência estrutural ",
		Users:   []string{"agentes", "mantenedores"},
		Capabilities: []string{
			"encontrar dependentes transitivos",
			"visualizar ciclos entre módulos",
		},
		Constraints: []string{"execução local"},
		NonGoals:    []string{"runtime tracing"}, PreferredStack: []string{"Rust"},
	}
}

func TestExpandIdeaProducesTraceableFacets(t *testing.T) {
	plan, err := ExpandIdea(completeIdeaFixture())
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Expansions) != 4 {
		t.Fatalf("esperava problema, duas capacidades e restrição; recebeu %+v", plan.Expansions)
	}
	wantFacets := []domain.ExpansionFacet{domain.FacetProblem, domain.FacetCapability, domain.FacetCapability, domain.FacetConstraint}
	for i, want := range wantFacets {
		if plan.Expansions[i].Facet != want || plan.Expansions[i].ID == "" || len(plan.Expansions[i].Covers) == 0 {
			t.Fatalf("expansão %d não rastreável: %+v", i, plan.Expansions[i])
		}
	}
	if plan.Idea.Problem != "mostrar consequência estrutural" {
		t.Fatalf("problema não normalizado: %q", plan.Idea.Problem)
	}
}

func TestExpandIdeaIsDeterministicAndDoesNotMutateInput(t *testing.T) {
	input := completeIdeaFixture()
	original := completeIdeaFixture()
	first, err := ExpandIdea(input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ExpandIdea(input)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("saída não determinística:\n%s\n%s", firstJSON, secondJSON)
	}
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("entrada modificada: antes=%+v depois=%+v", original, input)
	}
}

func TestExpandIdeaDeduplicatesCaseAndWhitespace(t *testing.T) {
	plan, err := ExpandIdea(domain.IdeaSpec{Capabilities: []string{"Impact Analysis", " impact   analysis ", "Callers"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Expansions) != 2 {
		t.Fatalf("duplicata não removida: %+v", plan.Expansions)
	}
}

func TestExpandIdeaRejectsEmptyOrStackOnlyIdea(t *testing.T) {
	for _, spec := range []domain.IdeaSpec{{}, {PreferredStack: []string{"Go"}}, {NonGoals: []string{"runtime"}}} {
		if _, err := ExpandIdea(spec); err == nil {
			t.Fatalf("ideia inválida aceita: %+v", spec)
		}
	}
}

func TestExpandIdeaNeverPromotesNonGoals(t *testing.T) {
	plan, err := ExpandIdea(domain.IdeaSpec{Problem: "dependency graph", NonGoals: []string{"runtime tracing"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, expansion := range plan.Expansions {
		if strings.Contains(strings.ToLower(expansion.Query), "runtime") {
			t.Fatalf("non-goal promovido: %+v", expansion)
		}
	}
}

func TestExpandIdeaMergesCoverageWhenQueriesDeduplicate(t *testing.T) {
	plan, err := ExpandIdea(domain.IdeaSpec{Problem: "impact analysis", Capabilities: []string{" Impact Analysis "}})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Expansions) != 1 || len(plan.Expansions[0].Covers) != 2 {
		t.Fatalf("cobertura perdida na deduplicação: %+v", plan.Expansions)
	}
}

func TestExpansionPlanKeepsTargetSeparateFromExternalQuery(t *testing.T) {
	target := domain.IdeaSpec{Problem: "ferramenta de impacto estrutural para agentes"}
	externalQuery := "code intelligence AST dependency graph"
	plan, err := ExpandIdea(target)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Idea.Problem != target.Problem || plan.Expansions[0].Query == externalQuery {
		t.Fatalf("alvo substituído por query externa: %+v", plan)
	}
}
