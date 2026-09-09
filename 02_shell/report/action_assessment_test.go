/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/action_assessment_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package report_test

import (
	"bateia/01_core/domain"
	"bateia/02_shell/report"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"testing"
)

func TestAction61MarkdownBoundsIntegrationAndJSON(t *testing.T) {
	a := domain.RequirementActionAssessment{Query: "start cycle", RequiredAction: "start", Status: "hypothesis", Reason: "ACTION_LIMIT_SYNTAX_ONLY"}
	for i := 0; i < 14; i++ {
		c := domain.ActionCandidate{Source: domain.ImplementationReference{Path: fmt.Sprintf("action_file_%02d.go", i), Symbol: fmt.Sprintf("ActionSymbol%02d", i), SHA256: fmt.Sprintf("action_hash_%02d", i), Line: 17}, EndLine: 19, Relation: "aligned_hint", Reason: "ACTION_CANDIDATE_LIMIT"}
		for j := 0; j < 6; j++ {
			c.Evidence = append(c.Evidence, domain.ActionEvidence{Kind: "call_name", Action: "start", Term: "start", Line: 18, EndLine: 18, Snippet: fmt.Sprintf("ActionEvidence%02dX%02d", i, j)})
		}
		a.Candidates = append(a.Candidates, c)
	}
	var sb strings.Builder
	report.WriteActionAssessments(&sb, []domain.RequirementActionAssessment{a})
	direct := sb.String()
	r := domain.ImplementationSearchReport{ActionAssessments: []domain.RequirementActionAssessment{a}}
	md := report.GenerateMarkdownReport(&domain.SearchReport{Results: []domain.ExtractedRepository{{ImplementationSearch: &r}}})
	for _, raw := range []string{direct, md} {
		out := html.UnescapeString(raw)
		for _, token := range []string{"start", "hypothesis", "aligned_hint", "ACTION_LIMIT_SYNTAX_ONLY", "ACTION_CANDIDATE_LIMIT", "17", "18", "JSON"} {
			if !strings.Contains(out, token) {
				t.Errorf("missing %s in %s", token, out)
			}
		}
		for i := 0; i < 14; i++ {
			for _, token := range []string{fmt.Sprintf("ActionSymbol%02d", i), fmt.Sprintf("action_file_%02d.go", i), fmt.Sprintf("action_hash_%02d", i)} {
				if strings.Contains(out, token) != (i < 12) {
					t.Errorf("candidate limit for %s", token)
				}
			}
			for j := 0; j < 6; j++ {
				token := fmt.Sprintf("ActionEvidence%02dX%02d", i, j)
				if strings.Contains(out, token) != (i < 12 && j < 4) {
					t.Errorf("evidence limit for %s", token)
				}
			}
		}
		if !strings.Contains(out, "2") {
			t.Error("missing omitted count")
		}
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.ImplementationSearchReport
	if err = json.Unmarshal(b, &restored); err != nil {
		t.Fatal(err)
	}
	if len(restored.ActionAssessments) != 1 || len(restored.ActionAssessments[0].Candidates) != 14 {
		t.Fatal("JSON candidate truncation")
	}
	for _, c := range restored.ActionAssessments[0].Candidates {
		if len(c.Evidence) != 6 {
			t.Fatal("JSON evidence truncation")
		}
	}
}
func TestAction61MarkdownEscapesUntrustedTextAndShowsNotEvaluated(t *testing.T) {
	a := domain.RequirementActionAssessment{Query: "<script>alert(1)</script>\n# forged heading", Status: "not_evaluated", Reason: "ACTION_UNSUPPORTED_REASON"}
	var sb strings.Builder
	report.WriteActionAssessments(&sb, []domain.RequirementActionAssessment{a})
	out := sb.String()
	visible := html.UnescapeString(out)
	if !strings.Contains(visible, "not_evaluated") || !strings.Contains(visible, "ACTION_UNSUPPORTED_REASON") {
		t.Fatalf("empty assessment hidden: %s", out)
	}
	if strings.Contains(out, "<script>") || strings.Contains(out, "\n# forged heading") {
		t.Fatalf("untrusted query rendered as markup: %s", out)
	}
}

func TestAction61MarkdownStablePriorityBeforeLimitWithoutMutation(t *testing.T) {
	relations := []string{"unknown", "future_relation", "opposed_hint", "conflicting_hints", "unknown", "opposed_hint", "conflicting_hints", "unknown", "opposed_hint", "unknown", "conflicting_hints", "opposed_hint", "aligned_hint", "aligned_hint", "unknown", "aligned_hint"}
	a := domain.RequirementActionAssessment{Query: "start cycle", RequiredAction: "start", Status: "hypothesis", Reason: "syntax only"}
	for i, relation := range relations {
		a.Candidates = append(a.Candidates, domain.ActionCandidate{Source: domain.ImplementationReference{Path: fmt.Sprintf("priority%02d.go", i), Symbol: fmt.Sprintf("PriorityCandidate%02dToken", i), SHA256: fmt.Sprintf("digest%02d", i), Line: 2}, EndLine: 3, Relation: relation, Reason: "syntax only", Evidence: []domain.ActionEvidence{{Kind: "call_name", Action: "start", Term: "start", Line: 2, EndLine: 2, Snippet: fmt.Sprintf("PriorityEvidence%02d", i)}}})
	}
	assessments := []domain.RequirementActionAssessment{a}
	r := domain.ImplementationSearchReport{ActionAssessments: assessments}
	before, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	report.WriteActionAssessments(&sb, assessments)
	afterDirect, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(afterDirect) {
		t.Fatal("direct presenter mutated report or shared candidate slice")
	}
	integrated := report.GenerateMarkdownReport(&domain.SearchReport{Results: []domain.ExtractedRepository{{ImplementationSearch: &r}}})
	afterIntegrated, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(afterIntegrated) {
		t.Fatal("integrated presenter mutated report or JSON order")
	}
	want := []int{12, 13, 15, 3, 6, 10, 2, 5, 8, 11, 0, 1}
	for _, raw := range []string{sb.String(), integrated} {
		out := html.UnescapeString(raw)
		previous := -1
		for _, i := range want {
			token := fmt.Sprintf("PriorityCandidate%02dToken", i)
			pos := strings.Index(out, token)
			if pos < 0 || pos <= previous {
				t.Fatalf("priority/stability before cut violated at %s: %s", token, out)
			}
			previous = pos
		}
		for _, i := range []int{4, 7, 9, 14} {
			if strings.Contains(out, fmt.Sprintf("PriorityCandidate%02dToken", i)) {
				t.Fatalf("candidate %d should be beyond presentation limit", i)
			}
		}
	}
}
