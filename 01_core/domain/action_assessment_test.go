/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/action_assessment_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestAction61JSONContract(t *testing.T) {
	a := domain.RequirementActionAssessment{Query: "iniciar ciclo", RequiredAction: "start", Status: "hypothesis", Reason: "syntax only", Candidates: []domain.ActionCandidate{{Source: domain.ImplementationReference{Path: "neutral.go", Symbol: "StartCycle", Line: 2, SHA256: "digest"}, EndLine: 4, Relation: "aligned_hint", Reason: "no execution proof", Evidence: []domain.ActionEvidence{{Kind: "declaration_name", Action: "start", Term: "start", Line: 2, EndLine: 2, Snippet: "StartCycle"}}}}}
	r := domain.ImplementationSearchReport{ActionAssessments: []domain.RequirementActionAssessment{a}}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"action_assessments", "query", "required_action", "status", "reason", "candidates", "source", "end_line", "relation", "evidence", "kind", "action", "term", "line", "snippet"} {
		if !strings.Contains(string(b), `"`+k+`":`) {
			t.Errorf("missing JSON field %s: %s", k, b)
		}
	}
	var back domain.ImplementationSearchReport
	if err = json.Unmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.ActionAssessments, back.ActionAssessments) {
		t.Fatal("JSON loses action evidence")
	}
	b, err = json.Marshal(domain.ImplementationSearchReport{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "action_assessments") {
		t.Fatal("empty projection must be omitted")
	}
}
