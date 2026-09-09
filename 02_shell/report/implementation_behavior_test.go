/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/implementation_behavior_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package report_test

import (
	"bateia/01_core/domain"
	"bateia/02_shell/report"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBehavior57MarkdownBoundsAndCompleteJSON(t *testing.T) {
	r := domain.ImplementationSearchReport{Status: "complete"}
	for i := 0; i < 12; i++ {
		b := domain.ImplementationBehavior{Source: domain.ImplementationReference{Path: "a.go", Symbol: fmt.Sprintf("BehaviorFunction%02d", i), Line: i + 1, SHA256: "hash"}, EndLine: i + 1, Status: "observed", Reason: "syntax", Limitations: []string{"BEHAVIOR_LIMITATION_EXECUTION", "BEHAVIOR_LIMITATION_ALIAS", "BEHAVIOR_LIMITATION_PRESERVATION", "BEHAVIOR_LIMITATION_INTENT"}}
		for j := 0; j < 11; j++ {
			b.Facts = append(b.Facts, domain.BehaviorFact{Kind: "assignment", Line: j + 1, EndLine: j + 1, Target: "x", Snippet: fmt.Sprintf("fact_token_%02d_%02d", i, j)})
		}
		r.Behavior = append(r.Behavior, b)
	}
	s := domain.SearchReport{Results: []domain.ExtractedRepository{{Ref: domain.CommitRef{Owner: "fixture", Repo: "behavior", Commit: "fixed"}, ImplementationSearch: &r}}}
	md := report.GenerateMarkdownReport(&s)
	for i := 0; i < 12; i++ {
		name := fmt.Sprintf("BehaviorFunction%02d", i)
		if strings.Contains(md, name) != (i < 10) {
			t.Errorf("function bound violated for %s", name)
		}
		for j := 0; j < 11; j++ {
			token := fmt.Sprintf("fact_token_%02d_%02d", i, j)
			if strings.Contains(md, token) != (i < 10 && j < 8) {
				t.Errorf("fact bound violated for %s", token)
			}
		}
	}
	for _, marker := range r.Behavior[0].Limitations {
		if !strings.Contains(md, marker) {
			t.Errorf("missing limitation %s", marker)
		}
	}
	lower := strings.ToLower(md)
	if !(strings.Contains(lower, "omit") || strings.Contains(lower, "omiss")) || !strings.Contains(md, "2") || !strings.Contains(md, "3") {
		t.Errorf("missing omission counts: %s", md)
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.SearchReport
	if err = json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	bs := restored.Results[0].ImplementationSearch.Behavior
	if len(bs) != 12 {
		t.Fatalf("JSON truncated functions: %d", len(bs))
	}
	for _, b := range bs {
		if len(b.Facts) != 11 {
			t.Fatalf("JSON truncated facts: %+v", b)
		}
	}
}
func TestBehavior57MarkdownIncludesEmptyAndNotEvaluated(t *testing.T) {
	r := domain.ImplementationSearchReport{Behavior: []domain.ImplementationBehavior{{Source: domain.ImplementationReference{Path: "a.go", Symbol: "EmptyFunction", Line: 2}, Status: "observed", EndLine: 2, Limitations: []string{"EMPTY_OBSERVATION_LIMIT"}}, {Source: domain.ImplementationReference{Path: "missing.go", Symbol: "UnknownFunction", Line: 2}, Status: "not_evaluated", Reason: "MISSING_SOURCE_REASON"}}}
	md := report.GenerateMarkdownReport(&domain.SearchReport{Results: []domain.ExtractedRepository{{ImplementationSearch: &r}}})
	for _, token := range []string{"EmptyFunction", "UnknownFunction", "EMPTY_OBSERVATION_LIMIT", "MISSING_SOURCE_REASON"} {
		if !strings.Contains(md, token) {
			t.Errorf("missing %s in %s", token, md)
		}
	}
}
