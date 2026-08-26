/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/lexical/expander_test.md
 * @layer L3
 * @updated 2026-08-24
 */

package lexical

import (
	"testing"

	"bateia/01_core/domain"
)

func sourcePlan(query string) domain.ExpansionPlan {
	return domain.ExpansionPlan{Expansions: []domain.SearchExpansion{{
		ID: "source-1", Facet: domain.FacetProblem, Query: query, Covers: []string{"problem"},
	}}}
}

func TestNewFromJSONRejectsInvalidVocabulary(t *testing.T) {
	invalid := [][]byte{[]byte(`{`), []byte(`{"rules":[]}`), []byte(`{"rules":[{"id":"x","when_any":[],"expand_to":["y"]}]}`)}
	for _, data := range invalid {
		if _, err := NewFromJSON(data, 2); err == nil {
			t.Fatalf("vocabulário inválido aceito: %s", data)
		}
	}
}

func TestExpanderMatchesWholePhraseAndCarriesProvenance(t *testing.T) {
	expander, err := NewFromJSON([]byte(`{"rules":[{"id":"impact","when_any":["impact"],"expand_to":["blast radius"]}]}`), 2)
	if err != nil {
		t.Fatal(err)
	}
	none, err := expander.Expand(sourcePlan("impactful change"))
	if err != nil {
		t.Fatal(err)
	}
	if len(none) != 0 {
		t.Fatalf("substring parcial ativou regra: %+v", none)
	}
	got, err := expander.Expand(sourcePlan("source impact analysis"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Derivation == nil || got[0].Derivation.RuleIDs[0] != "impact" {
		t.Fatalf("proveniência ausente: %+v", got)
	}
}

func TestExpanderRespectsLimitPerSource(t *testing.T) {
	expander, err := NewFromJSON([]byte(`{"rules":[{"id":"impact","when_any":["impact"],"expand_to":["one","two","three"]}]}`), 2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := expander.Expand(sourcePlan("impact analysis"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("limite ignorado: %+v", got)
	}
}
