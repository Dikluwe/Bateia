/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/requirement_retrieval_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func requirement0046Hit() domain.ImplementationHit {
	return domain.ImplementationHit{Retrieval: "lexical", Path: "queue.go", Symbol: "Queue", Line: 3, SHA256: strings.Repeat("a", 64), MatchedTerms: []string{"queue", "token"}}
}

func TestRequirement0046MultilineIntersection(t *testing.T) {
	hit := requirement0046Hit()
	got := parser.DescribeImplementationRequirements("\n queue save \n token \nqueue save\n\n", domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{hit}})
	if len(got) != 2 {
		t.Fatalf("requirements = %#v", got)
	}
	if got[0].Requirement != "queue save" || got[1].Requirement != "token" {
		t.Fatalf("dedup/order: %#v", got)
	}
	for _, r := range got {
		if r.Status != "evidence_located" || r.Reason == "" || len(r.Sources) != 1 {
			t.Fatalf("lexical evidence: %#v", r)
		}
	}
	if !reflect.DeepEqual(got[0].MatchedTerms, []string{"queue"}) || !reflect.DeepEqual(got[0].MissingTerms, []string{"save"}) {
		t.Fatalf("intersection must exclude token: %#v", got[0])
	}
	if !reflect.DeepEqual(got[1].MatchedTerms, []string{"token"}) || len(got[1].MissingTerms) != 0 {
		t.Fatalf("second intersection: %#v", got[1])
	}
	want := domain.ImplementationReference{Path: hit.Path, Symbol: hit.Symbol, Line: hit.Line, SHA256: hit.SHA256}
	if got[0].Sources[0] != want {
		t.Fatalf("identity: %#v", got[0].Sources)
	}
}

func TestRequirement0046RejectUnusableSources(t *testing.T) {
	cases := map[string]func(*domain.ImplementationHit){
		"direct call":        func(h *domain.ImplementationHit) { h.Retrieval = "direct_call" },
		"unspecified origin": func(h *domain.ImplementationHit) { h.Retrieval = "" },
		"missing path":       func(h *domain.ImplementationHit) { h.Path = "" },
		"missing symbol":     func(h *domain.ImplementationHit) { h.Symbol = "" },
		"zero line":          func(h *domain.ImplementationHit) { h.Line = 0 },
		"negative line":      func(h *domain.ImplementationHit) { h.Line = -1 },
		"missing hash":       func(h *domain.ImplementationHit) { h.SHA256 = "" },
		"short hash":         func(h *domain.ImplementationHit) { h.SHA256 = "abc" },
		"nonhex hash":        func(h *domain.ImplementationHit) { h.SHA256 = strings.Repeat("z", 64) },
		"unrelated terms":    func(h *domain.ImplementationHit) { h.MatchedTerms = []string{"auth"} },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			h := requirement0046Hit()
			change(&h)
			got := parser.DescribeImplementationRequirements("queue", domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{h}})
			if len(got) != 1 || got[0].Status != "not_verified" || len(got[0].Sources) != 0 || len(got[0].MatchedTerms) != 0 || !reflect.DeepEqual(got[0].MissingTerms, []string{"queue"}) {
				t.Fatalf("unusable source became evidence: %#v", got)
			}
		})
	}
}

func TestRequirement0046DeterminismDedupAndPurity(t *testing.T) {
	a := requirement0046Hit()
	b := a
	b.Path = "a.go"
	b.Symbol = "Other"
	b.Line = 8
	report := domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{a, b, a}}
	before, _ := json.Marshal(report)
	got := parser.DescribeImplementationRequirements("queue", report)
	after, _ := json.Marshal(report)
	if string(before) != string(after) {
		t.Fatal("projection mutated report")
	}
	if len(got) != 1 || len(got[0].Sources) != 2 || got[0].Sources[0].Path != "a.go" || got[0].Sources[1].Path != "queue.go" {
		t.Fatalf("sources must dedup and sort: %#v", got)
	}
	report.Hits = []domain.ImplementationHit{b, a, b}
	again := parser.DescribeImplementationRequirements("queue\nqueue", report)
	if !reflect.DeepEqual(got, again) {
		t.Fatalf("permutation changed result: %#v vs %#v", got, again)
	}
}

func TestRequirement0046NoEvidenceIsNotVerified(t *testing.T) {
	for _, report := range []domain.ImplementationSearchReport{{}, {Status: "unsupported"}, {Status: "partial", Limited: true}} {
		got := parser.DescribeImplementationRequirements("queue", report)
		if len(got) != 1 || got[0].Status != "not_verified" || len(got[0].Sources) != 0 || got[0].Reason == "" {
			t.Fatalf("no support promoted: %#v", got)
		}
	}
	h := requirement0046Hit()
	got := parser.DescribeImplementationRequirements("a the how", domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{h}})
	if len(got) != 1 || got[0].Status != "not_verified" || len(got[0].Sources) != 0 || len(got[0].MatchedTerms) != 0 || len(got[0].MissingTerms) != 0 {
		t.Fatalf("stopwords became evidence: %#v", got)
	}
	if got := parser.DescribeImplementationRequirements(" \n\t\n", domain.ImplementationSearchReport{}); len(got) != 0 {
		t.Fatalf("blank query: %#v", got)
	}
}
