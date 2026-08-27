/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_queries_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func functionalFixture() domain.FunctionalLexicon {
	return domain.FunctionalLexicon{SchemaVersion: "1", Version: "t", Frames: map[domain.DiscoveryConceptID]domain.FunctionalFrame{"repository-recommendation": {ID: "f", ConceptID: "repository-recommendation", Inputs: []string{"natural language project idea"}, Actions: []string{"find and rank"}, Objects: []string{"GitHub repositories"}, Outputs: []string{"repository shortlist"}, RuleID: "r"}}, Forbidden: []string{"secret competitor"}}
}
func TestFunctionalExpansionDescribesJobAndPreservesOrigin(t *testing.T) {
	o := domain.OntologyExpansionResult{Matches: []domain.ConceptMatch{{ConceptID: "repository-recommendation", SourceID: "s"}}}
	templates := []domain.FunctionalQueryTemplate{{ID: "T1", Pattern: "{action} {object} from {input}"}, {ID: "T2", Pattern: "tool that {action} {object} and returns {output}"}}
	got, err := ExpandFunctionalQueries(o, functionalFixture(), templates)
	if err != nil || len(got.Queries) != 2 {
		t.Fatalf("%+v %v", got, err)
	}
	if got.Queries[0].Text != "find and rank GitHub repositories from natural language project idea" || got.Queries[0].SourceIDs[0] != "s" {
		t.Fatalf("bad query: %+v", got.Queries[0])
	}
}
func TestFunctionalAuditRejectsExpectedOutcomeLeak(t *testing.T) {
	l := functionalFixture()
	e := domain.FunctionalExpansion{Queries: []domain.FunctionalQuery{{ID: "x", Text: "find secret competitor", Actions: []string{"find"}, Objects: []string{"tool"}, Inputs: []string{"idea"}}}}
	if AuditFunctionalQueries(e, l)[0].Verdict != "REJECTED" {
		t.Fatal("leak accepted")
	}
}
