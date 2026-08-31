/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_grammar_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func grammarFixture() domain.FunctionalGrammar {
	return domain.FunctionalGrammar{SchemaVersion: "2", Version: "v2-test", Rules: []domain.FunctionalGrammarRule{
		{ID: "imp", Intent: "adoption-decision", Pattern: "{action} {object} from {input} and produce {output}", RequiredSlots: []string{"action", "object", "input", "output"}, Priority: 1, SubjectKind: "imperative"},
		{ID: "third", Intent: "adoption-decision", Pattern: "tool that {action} {object} from {input} and returns {output}", RequiredSlots: []string{"action", "object", "input", "output"}, Priority: 2, SubjectKind: "singular_3rd"},
	}, Morphology: []domain.MorphologyRule{{ID: "m-evaluate", SubjectKind: "singular_3rd", Source: "evaluate", Realized: "evaluates"}}}
}
func expansionFixture() domain.FunctionalExpansion {
	return domain.FunctionalExpansion{Frames: []domain.FunctionalFrame{{ID: "frame", ConceptID: "adoption-decision", SourceIDs: []string{"source"}, Inputs: []string{"GitHub repository"}, Actions: []string{"evaluate"}, Objects: []string{"open-source tool"}, Outputs: []string{"adoption verdict"}}}}
}

func TestRealizeFunctionalGrammarUsesIntentRulesAndMorphology(t *testing.T) {
	got, err := RealizeFunctionalGrammar(expansionFixture(), grammarFixture())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Queries) != 2 {
		t.Fatalf("queries=%d", len(got.Queries))
	}
	texts := map[string]bool{}
	for _, q := range got.Queries {
		texts[q.Text] = true
		if q.GrammarRuleID == "" || q.GrammarHash == "" {
			t.Fatal("missing lineage")
		}
	}
	if !texts["evaluate open-source tool from GitHub repository and produce adoption verdict"] {
		t.Fatal("imperative missing")
	}
	if !texts["tool that evaluates open-source tool from GitHub repository and returns adoption verdict"] {
		t.Fatal("third person missing")
	}
}
func TestUnknownVerbFailsClosed(t *testing.T) {
	e := expansionFixture()
	e.Frames[0].Actions = []string{"assess"}
	got, err := RealizeFunctionalGrammar(e, grammarFixture())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Queries) != 1 || len(got.Unknown) == 0 {
		t.Fatalf("queries=%d unknown=%d", len(got.Queries), len(got.Unknown))
	}
}
func TestMissingSlotFailsClosed(t *testing.T) {
	e := expansionFixture()
	e.Frames[0].Outputs = nil
	got, _ := RealizeFunctionalGrammar(e, grammarFixture())
	if len(got.Queries) != 0 || len(got.Unknown) < 1 {
		t.Fatal("missing slot advanced")
	}
}
func TestReorderingDoesNotChangeOutput(t *testing.T) {
	g := grammarFixture()
	a, _ := RealizeFunctionalGrammar(expansionFixture(), g)
	g.Rules[0], g.Rules[1] = g.Rules[1], g.Rules[0]
	b, _ := RealizeFunctionalGrammar(expansionFixture(), g)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("reordering changed output")
	}
}
func TestAuditRejectsForbiddenAndHashMutation(t *testing.T) {
	g := grammarFixture()
	got, _ := RealizeFunctionalGrammar(expansionFixture(), g)
	got.Queries[0].Text += " gitmvp"
	got.Queries[0].GrammarHash = "tampered"
	fs := AuditFunctionalGrammar(got, g, []string{"gitmvp"})
	if fs[0].Verdict != "REJECTED" || len(fs[0].ForbiddenHits) == 0 {
		t.Fatal("mutation survived")
	}
}
