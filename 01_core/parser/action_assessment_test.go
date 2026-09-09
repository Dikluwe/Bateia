/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/action_assessment_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"bateia/01_core/usecase"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func action61Fixture(name, body string) (domain.ImplementationSearchReport, []domain.ExtractedFile) {
	src := "package neutral\nfunc " + name + "() {\n" + body + "\n}\n"
	return domain.ImplementationSearchReport{Status: "complete", Hits: []domain.ImplementationHit{{Path: "neutral.go", Symbol: name, Kind: "definition", Retrieval: "lexical", Line: 2, EndLine: strings.Count(src, "\n"), SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(src))), MatchedTerms: []string{"cycle"}}}}, []domain.ExtractedFile{{Path: "neutral.go", Content: src, Size: int64(len(src))}}
}
func action61One(t *testing.T, query string, r domain.ImplementationSearchReport, f []domain.ExtractedFile) domain.RequirementActionAssessment {
	t.Helper()
	got := parser.AssessImplementationActions(query, r, f)
	if len(got) != 1 {
		t.Fatalf("one interpretation expected: %+v", got)
	}
	return got[0]
}
func TestAction61QueryBoundary(t *testing.T) {
	r, f := action61Fixture("StartCycle", "")
	for _, tc := range []struct{ q, action string }{
		{"iniciar ciclo", "start"}, {"começar ciclo", "start"}, {"start cycle", "start"}, {"subscribe cycle", "start"},
		{"parar ciclo", "stop"}, {"interromper ciclo", "stop"}, {"cancelar ciclo", "stop"}, {"cessar ciclo", "stop"}, {"stop cycle", "stop"}, {"cancel cycle", "stop"}, {"unsubscribe cycle", "stop"},
		{"restart cycle", ""}, {"starting cycle", ""}, {"iniciarmos ciclo", ""}, {"cycle monitor", ""}, {"start stop cycle", ""}, {"iniciar parar ciclo", ""}, {"não iniciar ciclo", ""}, {"iniciar ciclo sem bloqueio", ""}, {"never start cycle", ""}, {"start `Cycle`", ""}, {"`StartCycle`", ""},
	} {
		t.Run(tc.q, func(t *testing.T) {
			a := action61One(t, tc.q, r, f)
			if a.Query != tc.q || a.RequiredAction != tc.action || a.Reason == "" {
				t.Fatalf("query contract: %+v", a)
			}
			if tc.action == "" {
				if a.Status != "not_evaluated" || len(a.Candidates) != 0 {
					t.Fatalf("unsupported query classified: %+v", a)
				}
			} else if a.Status != "hypothesis" || len(a.Candidates) != 1 {
				t.Fatalf("hypothesis missing: %+v", a)
			}
		})
	}
	a := parser.AssessImplementationActions("stop cycle\nstart cycle\nstop cycle", r, f)
	b := parser.AssessImplementationActions("start cycle\nstop cycle", r, f)
	if len(a) != 2 || !reflect.DeepEqual(a, b) {
		t.Fatalf("canonical interpretations lost: %+v / %+v", a, b)
	}
}
func TestAction61NamesCallsAndLexicalBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name, body, relation string
		count                int
	}{
		{"StartCycle", "", "aligned_hint", 1}, {"StopCycle", "", "opposed_hint", 1},
		{"StartCycle", "svc.StopCycle()", "conflicting_hints", 2}, {"InspectCycle", "svc.SubscribeCycle()", "aligned_hint", 1},
		{"InspectCycle", "UnsubscribeCycle()", "opposed_hint", 1}, {"InspectCycle", "", "unknown", 0},
		{"RestartCycle", "restartCycle()", "unknown", 0}, {"NotStartCycle", "NeverStopCycle(); WithoutWatchCycle(); NoCancelCycle()", "unknown", 0},
		{"InspectCycle", "// StartCycle()\n_ = \"StopCycle()\"\n_ = StartCycle\n_ = func(){ StartCycle(); StopCycle() }", "unknown", 0},
		{"HTTPStartCycle", "", "aligned_hint", 1}, {"watch_cycle", "", "aligned_hint", 1},
		{"InspectCycle", "if false { StartCycle() }", "aligned_hint", 1},
	} {
		t.Run(tc.name+tc.body, func(t *testing.T) {
			r, f := action61Fixture(tc.name, tc.body)
			a := action61One(t, "start cycle", r, f)
			if len(a.Candidates) != 1 {
				t.Fatal(a)
			}
			c := a.Candidates[0]
			if c.Relation != tc.relation || len(c.Evidence) != tc.count || c.Reason == "" {
				t.Fatalf("syntactic boundary: %+v", c)
			}
			if c.Source.Symbol != tc.name || c.Source.SHA256 != r.Hits[0].SHA256 || c.EndLine != r.Hits[0].EndLine {
				t.Fatalf("identity lost: %+v", c)
			}
		})
	}
	for _, term := range []string{"add", "start", "subscribe", "register", "watch", "remove", "stop", "unsubscribe", "unregister", "cancel", "close", "unwatch"} {
		r, f := action61Fixture("InspectCycle", term+"()")
		c := action61One(t, "start cycle", r, f).Candidates[0]
		if len(c.Evidence) != 1 || c.Evidence[0].Term != term || c.Evidence[0].Snippet != term || c.Evidence[0].Kind != "call_name" || c.Evidence[0].Line != 3 || c.Evidence[0].EndLine != 3 {
			t.Fatalf("token evidence %s: %+v", term, c)
		}
	}
}
func TestAction61RejectsUnboundBytes(t *testing.T) {
	cases := map[string]func(*domain.ImplementationSearchReport, *[]domain.ExtractedFile){
		"hash": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			r.Hits[0].SHA256 = strings.Repeat("f", 64)
		},
		"missing":   func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { *f = nil },
		"duplicate": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { *f = append(*f, (*f)[0]) },
		"size":      func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { (*f)[0].Size++ },
		"line":      func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].Line++ },
		"end":       func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].EndLine++ },
		"symbol":    func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].Symbol = "Other" },
		"documentation": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			r.Hits[0].Retrieval = "documentation"
		},
		"language": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			r.Hits[0].Path = "neutral.rs"
			(*f)[0].Path = "neutral.rs"
		},
		"invalid_ast": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			s := "package neutral\nfunc StartCycle( {"
			(*f)[0].Content = s
			(*f)[0].Size = int64(len(s))
			r.Hits[0].SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
		},
		"bodyless": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			s := "package neutral\nfunc StartCycle()\n"
			(*f)[0].Content = s
			(*f)[0].Size = int64(len(s))
			r.Hits[0].EndLine = 2
			r.Hits[0].SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
		},
	}
	for name, mut := range cases {
		t.Run(name, func(t *testing.T) {
			r, f := action61Fixture("StartCycle", "")
			mut(&r, &f)
			a := action61One(t, "start cycle", r, f)
			if len(a.Candidates) != 1 || a.Candidates[0].Relation != "unknown" || len(a.Candidates[0].Evidence) != 0 || a.Candidates[0].Reason == "" {
				t.Fatalf("untrusted candidate classified: %+v", a)
			}
		})
	}
}
func TestAction61BudgetIsNotPartialClassification(t *testing.T) {
	for _, n := range []int{32, 33} {
		body := strings.Repeat("StartCycle()\n", n-1) + "StopCycle()"
		r, f := action61Fixture("InspectCycle", body)
		c := action61One(t, "start cycle", r, f).Candidates[0]
		if n == 32 {
			if c.Relation != "conflicting_hints" || len(c.Evidence) != 32 {
				t.Fatalf("budget boundary lost: %+v", c)
			}
		} else if c.Relation != "unknown" || len(c.Evidence) > 32 || c.Reason == "" {
			t.Fatalf("truncated evidence classified: %+v", c)
		}
	}
}
func TestAction61PhysicalLinesDeterminismAndPurity(t *testing.T) {
	src := "package neutral\n//line imaginary.go:900\nfunc StartCycle() { StopCycle() }\n"
	r, f := action61Fixture("StartCycle", "")
	f[0].Content = src
	f[0].Size = int64(len(src))
	r.Hits[0].SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(src)))
	r.Hits[0].Line = 3
	r.Hits[0].EndLine = 3
	r.Hits = append(r.Hits, r.Hits[0])
	r.Hits[0].Path = "z.go"
	original := append([]domain.ImplementationHit(nil), r.Hits...)
	a := action61One(t, "start cycle", r, f)
	if len(a.Candidates) != 2 || a.Candidates[0].Source.Path != "neutral.go" || a.Candidates[0].Relation != "conflicting_hints" {
		t.Fatalf("physical identity/order: %+v", a)
	}
	for _, e := range a.Candidates[0].Evidence {
		if e.Line != 3 || e.EndLine != 3 {
			t.Fatalf("virtual line used: %+v", e)
		}
	}
	if !reflect.DeepEqual(r.Hits, original) {
		t.Fatal("hits reordered in place")
	}
	r.Hits[0], r.Hits[1] = r.Hits[1], r.Hits[0]
	if !reflect.DeepEqual(a, action61One(t, "start cycle", r, f)) {
		t.Fatal("input order affected projection")
	}
}
func TestAction61ValidationReplacesForgedProjection(t *testing.T) {
	r, f := action61Fixture("StartCycle", "")
	r.ActionAssessments = []domain.RequirementActionAssessment{{Query: "forged", Status: "verified", Candidates: []domain.ActionCandidate{{Relation: "satisfied"}}}}
	got := usecase.ValidateImplementationRequirements("start cycle", r, f)
	if len(got.ActionAssessments) != 1 || got.ActionAssessments[0].Status != "hypothesis" || len(got.ActionAssessments[0].Candidates) != 1 || got.ActionAssessments[0].Candidates[0].Relation != "aligned_hint" {
		t.Fatalf("fresh projection missing: %+v", got.ActionAssessments)
	}
	if r.ActionAssessments[0].Status != "verified" || got.Status != r.Status || !reflect.DeepEqual(got.Hits, r.Hits) {
		t.Fatal("unrelated state changed")
	}
	stale := usecase.ValidateImplementationRequirements("start cycle", got, nil)
	for _, a := range stale.ActionAssessments {
		for _, c := range a.Candidates {
			if c.Relation != "unknown" || len(c.Evidence) != 0 {
				t.Fatalf("stale action survived: %+v", c)
			}
		}
	}
}
