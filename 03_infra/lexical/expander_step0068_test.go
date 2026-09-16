/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/lexical/expander_step0068_test.md
 * @layer L3
 * @updated 2026-09-11
 */

package lexical

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const step0068Portuguese = "predicados geometricos de orientacao e in-circle robustos a casos quase degenerados"

func step0068SourcePlan(query string) domain.ExpansionPlan {
	return domain.ExpansionPlan{
		Idea: domain.IdeaSpec{Capabilities: []string{query}},
		Expansions: []domain.SearchExpansion{{
			ID: "exp-capability-source", Facet: domain.FacetCapability, Query: query,
			Covers: []string{"capabilities[0]"},
		}},
	}
}

func TestStep0068DefaultCatalogCompilesExactTechnicalQueries(t *testing.T) {
	expander, err := NewDefault(3)
	if err != nil {
		t.Fatal(err)
	}
	input := step0068SourcePlan(step0068Portuguese)
	before, _ := json.Marshal(input)
	got, err := expander.Expand(input)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(input)
	if string(before) != string(after) {
		t.Fatal("expander modificou o plano literal")
	}

	want := []struct {
		rule, query string
		discrims    []string
	}{
		{"robust-predicates.translation.en", "robust geometric orientation and in-circle predicates for nearly degenerate cases", []string{"orientation", "in-circle"}},
		{"robust-predicates.symbol-pair", "orient2d incircle", []string{"orient2d", "incircle"}},
		{"robust-predicates.reference-pair", "Shewchuk exact predicates", []string{"Shewchuk", "predicates"}},
	}
	if len(got) != len(want) {
		t.Fatalf("propostas=%d, quero 3: %+v", len(got), got)
	}
	for index, expected := range want {
		proposal := got[index]
		if proposal.Query != expected.query || !reflect.DeepEqual(proposal.Discriminators, expected.discrims) {
			t.Fatalf("proposta[%d]=%+v, esperava query=%q discriminadores=%v", index, proposal, expected.query, expected.discrims)
		}
		if proposal.Facet != domain.FacetCapability || !reflect.DeepEqual(proposal.Covers, []string{"capabilities[0]"}) {
			t.Fatalf("proposta[%d] perdeu facet/covers: %+v", index, proposal)
		}
		if proposal.Derivation == nil || proposal.Derivation.Strategy != "technical-lexicon-v1" || !reflect.DeepEqual(proposal.Derivation.RuleIDs, []string{expected.rule}) || !reflect.DeepEqual(proposal.Derivation.SourceIDs, []string{"exp-capability-source"}) {
			t.Fatalf("proposta[%d] perdeu proveniência técnica: %+v", index, proposal.Derivation)
		}
	}
	for _, proposal := range got {
		if proposal.Query == "predicados geometricos" || proposal.Query == "geometric predicates" {
			t.Fatalf("fallback genérico proibido foi produzido: %q", proposal.Query)
		}
	}
}

func TestStep0068EnrichmentPreservesOriginalPhraseAndStableIDs(t *testing.T) {
	expander, err := NewDefault(3)
	if err != nil {
		t.Fatal(err)
	}
	literal := step0068SourcePlan(step0068Portuguese)
	result, err := usecase.EnrichExpansionPlan(literal, expander)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Plan.Expansions) != 4 || result.Plan.Expansions[0].Query != step0068Portuguese || result.Plan.Idea.Capabilities[0] != step0068Portuguese {
		t.Fatalf("frase original não permaneceu intacta: %+v", result.Plan)
	}
	seenIDs := map[string]bool{result.Plan.Expansions[0].ID: true}
	for _, proposal := range result.Plan.Expansions[1:] {
		if strings.TrimSpace(proposal.ID) == "" || seenIDs[proposal.ID] {
			t.Fatalf("ID técnico vazio ou duplicado: %+v", proposal)
		}
		seenIDs[proposal.ID] = true
	}
}

func TestStep0068DuplicateSourcesDeduplicateQueriesAndPreserveProvenance(t *testing.T) {
	expander, err := NewDefault(3)
	if err != nil {
		t.Fatal(err)
	}
	plan := step0068SourcePlan(step0068Portuguese)
	duplicate := plan.Expansions[0]
	duplicate.ID = "exp-capability-duplicate"
	duplicate.Covers = []string{"capabilities[1]"}
	plan.Expansions = append(plan.Expansions, duplicate)

	result, err := usecase.EnrichExpansionPlan(plan, expander)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Plan.Expansions) != 5 || result.ProposedCount != 6 || result.DuplicatesRemoved != 3 {
		t.Fatalf("duplicação de origem não foi deduplicada: %+v", result)
	}
	for _, proposal := range result.Plan.Expansions[2:] {
		if proposal.Derivation == nil || !reflect.DeepEqual(proposal.Derivation.SourceIDs, []string{"exp-capability-source", "exp-capability-duplicate"}) {
			t.Fatalf("origens duplicadas não foram preservadas: %+v", proposal)
		}
	}
}

func TestStep0068CatalogBoundaryAndTwentyNormalizedInputs(t *testing.T) {
	expander, err := NewDefault(3)
	if err != nil {
		t.Fatal(err)
	}
	for _, opaque := range []string{
		"predicados geometricos",
		"robust geometric predicates",
		"predicados geometricos de distancia robustos",
		"predicados geométricos de orientacao e in-circle robustos a casos quase degenerados",
	} {
		got, err := expander.Expand(step0068SourcePlan(opaque))
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("entrada parcial/opaca %q ativou catálogo: %+v", opaque, got)
		}
	}
	embedded, err := expander.Expand(step0068SourcePlan("biblioteca C++ para " + step0068Portuguese + " com licença permissiva"))
	if err != nil {
		t.Fatal(err)
	}
	if len(embedded) != 3 {
		t.Fatalf("frase integral contida não ativou as três regras: %+v", embedded)
	}

	var canonical []domain.SearchExpansion
	for index := 0; index < 20; index++ {
		query := step0068Portuguese
		if index%2 == 0 {
			query = strings.ToUpper(query)
		}
		if index%3 == 0 {
			query = "  " + strings.ReplaceAll(query, " ", "   ") + "  "
		}
		got, err := expander.Expand(step0068SourcePlan(query))
		if err != nil {
			t.Fatal(err)
		}
		if index == 0 {
			canonical = got
			continue
		}
		if !reflect.DeepEqual(got, canonical) {
			t.Fatalf("permutação normalizada %d mudou propostas:\n%+v\n%+v", index, canonical, got)
		}
	}
}

func TestStep0068TechnicalVocabularyValidationAndLimits(t *testing.T) {
	long65 := strings.Repeat("x", 65)
	long257 := strings.Repeat("q", 257)
	invalid := []string{
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["query"],"strategy":"technical-lexicon-v1","discriminators":["one"]}]}`,
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["query"],"strategy":"technical-lexicon-v1","discriminators":["one","ONE"]}]}`,
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["query"],"strategy":"technical-lexicon-v1","discriminators":["one","two","three"]}]}`,
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["query"],"strategy":"unknown","discriminators":["one","two"]}]}`,
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["` + long257 + `"],"strategy":"technical-lexicon-v1","discriminators":["one","two"]}]}`,
		`{"rules":[{"id":"t","when_any":["source"],"expand_to":["query"],"strategy":"technical-lexicon-v1","discriminators":["` + long65 + `","two"]}]}`,
	}
	for _, raw := range invalid {
		if _, err := NewFromJSON([]byte(raw), 3); err == nil {
			t.Fatalf("vocabulário técnico inválido aceito: %s", raw)
		}
	}

	expander, err := NewDefault(2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := expander.Expand(step0068SourcePlan(step0068Portuguese))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Query == got[1].Query {
		t.Fatalf("limite explícito por origem não foi preservado: %+v", got)
	}

	expander, err = NewDefault(3)
	if err != nil {
		t.Fatal(err)
	}
	over := domain.ExpansionPlan{Expansions: make([]domain.SearchExpansion, 62)}
	for index := range over.Expansions {
		over.Expansions[index] = domain.SearchExpansion{ID: "source-" + strings.Repeat("x", index%5+1), Query: "opaque"}
	}
	over.Expansions[0] = step0068SourcePlan(step0068Portuguese).Expansions[0]
	if proposals, err := expander.Expand(over); err == nil {
		t.Fatalf("limite total 64 não produziu erro/Unknown: base=%d propostas=%d", len(over.Expansions), len(proposals))
	}
}

func TestStep0068LegacyRuleIsByteCompatible(t *testing.T) {
	raw := []byte(`{"rules":[{"id":"legacy","when_any":["impact"],"expand_to":["code impact analysis","change blast radius"]}]}`)
	expander, err := NewFromJSON(raw, 2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := expander.Expand(sourcePlan("source impact analysis"))
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"id":"","facet":"problem","query":"code impact analysis","rationale":"ampliar vocabulário por regra lexical legacy","covers":["problem"],"derivation":{"strategy":"lexical","rule_ids":["legacy"],"source_ids":["source-1"]}},{"id":"","facet":"problem","query":"change blast radius","rationale":"ampliar vocabulário por regra lexical legacy","covers":["problem"],"derivation":{"strategy":"lexical","rule_ids":["legacy"],"source_ids":["source-1"]}}]`
	if string(encoded) != want {
		t.Fatalf("semântica serializada da regra legada mudou:\ngot  %s\nwant %s", encoded, want)
	}
}
