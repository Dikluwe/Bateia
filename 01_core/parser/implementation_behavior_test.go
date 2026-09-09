/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/implementation_behavior_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"bateia/01_core/usecase"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func behavior57Fixture(src string, start, end int) (domain.ImplementationSearchReport, []domain.ExtractedFile) {
	return domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{{Path: "sample.go", Symbol: "Change", Kind: "definition", Line: start, EndLine: end, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(src))), Retrieval: "lexical", MatchedTerms: []string{"change"}}}}, []domain.ExtractedFile{{Path: "sample.go", Content: src, Size: int64(len(src))}}
}
func behavior57Observed(t *testing.T, got []domain.ImplementationBehavior) domain.ImplementationBehavior {
	t.Helper()
	if len(got) != 1 || got[0].Status != "observed" {
		t.Fatalf("expected observed function: %+v", got)
	}
	b := got[0]
	limits := strings.ToLower(strings.Join(b.Limitations, " "))
	for _, alternatives := range [][]string{{"flux", "execu", "flow"}, {"alias"}, {"chamad", "call"}, {"indiret", "indirect"}, {"preserv"}, {"inten", "intent"}} {
		found := false
		for _, term := range alternatives {
			found = found || strings.Contains(limits, term)
		}
		if !found {
			t.Errorf("missing limitation concerning %v: %+v", alternatives, b.Limitations)
		}
	}
	return b
}
func TestBehavior57ExactSyntaxAndClosureExclusion(t *testing.T) {
	src := "package sample\nfunc Change(input []int) []int {\n // input is immutable; this comment is not evidence\n alias := input\n alias[0], input[1] = 2, 3\n input[0]++\n closure := func() { input[2] = 9; input[2]--; return }\n opaque(input)\n if false { input[3] = 8 }\n _ = closure\n return append([]int(nil), input...)\n}\n"
	r, f := behavior57Fixture(src, 2, 12)
	before, _ := json.Marshal(struct {
		R domain.ImplementationSearchReport
		F []domain.ExtractedFile
	}{r, f})
	b := behavior57Observed(t, parser.InspectImplementationBehavior(r, f))
	want := []domain.BehaviorFact{
		{Kind: "assignment", Line: 4, EndLine: 4, Target: "alias", Snippet: "alias := input"},
		{Kind: "assignment", Line: 5, EndLine: 5, Target: "alias[0]", Snippet: "alias[0], input[1] = 2, 3"},
		{Kind: "assignment", Line: 5, EndLine: 5, Target: "input[1]", Snippet: "alias[0], input[1] = 2, 3"},
		{Kind: "increment_decrement", Line: 6, EndLine: 6, Target: "input[0]", Snippet: "input[0]++"},
		{Kind: "assignment", Line: 7, EndLine: 7, Target: "closure", Snippet: "closure := func() { input[2] = 9; input[2]--; return }"},
		{Kind: "assignment", Line: 9, EndLine: 9, Target: "input[3]", Snippet: "input[3] = 8"},
		{Kind: "assignment", Line: 10, EndLine: 10, Target: "_", Snippet: "_ = closure"},
		{Kind: "return", Line: 11, EndLine: 11, Target: "", Snippet: "return append([]int(nil), input...)"},
	}
	if !reflect.DeepEqual(b.Facts, want) {
		t.Fatalf("facts must describe syntax only:\ngot %#v\nwant %#v", b.Facts, want)
	}
	h := r.Hits[0]
	if b.Source != (domain.ImplementationReference{Path: h.Path, Symbol: h.Symbol, Line: h.Line, SHA256: h.SHA256}) || b.EndLine != 12 {
		t.Fatalf("identity lost: %+v", b)
	}
	after, _ := json.Marshal(struct {
		R domain.ImplementationSearchReport
		F []domain.ExtractedFile
	}{r, f})
	if string(before) != string(after) {
		t.Fatal("inspection mutated inputs")
	}
}
func TestBehavior57ReturnAndOpaqueCallAreNotVerification(t *testing.T) {
	for _, body := range []string{"opaque(input)", "return input", "return []int{1}", "// input[0] = 10\n"} {
		src := "package sample\nfunc Change(input []int) []int {\n" + body + "\n}\n"
		end := strings.Count(src, "\n")
		r, f := behavior57Fixture(src, 2, end)
		b := behavior57Observed(t, parser.InspectImplementationBehavior(r, f))
		want := 0
		if strings.HasPrefix(body, "return") {
			want = 1
		}
		if len(b.Facts) != want {
			t.Fatalf("calls/comments must not become facts: %+v", b)
		}
	}
}
func TestBehavior57RejectUnboundIdentity(t *testing.T) {
	src := "package sample\nfunc Change(input []int) { input[0] = 1 }\n"
	cases := map[string]func(*domain.ImplementationSearchReport, *[]domain.ExtractedFile){
		"hash": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			r.Hits[0].SHA256 = strings.Repeat("0", 64)
		},
		"missing":   func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { *f = nil },
		"duplicate": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { *f = append(*f, (*f)[0]) },
		"size":      func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { (*f)[0].Size++ },
		"start":     func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].Line = 1 },
		"end":       func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].EndLine = 3 },
		"symbol":    func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].Symbol = "Other" },
		"kind":      func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) { r.Hits[0].Kind = "call" },
		"cpp": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			r.Hits[0].Path = "sample.cpp"
			(*f)[0].Path = "sample.cpp"
		},
		"invalid": func(r *domain.ImplementationSearchReport, f *[]domain.ExtractedFile) {
			s := "package sample\nfunc Change( {"
			(*f)[0].Content = s
			(*f)[0].Size = int64(len(s))
			r.Hits[0].SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
		},
	}
	for name, mut := range cases {
		t.Run(name, func(t *testing.T) {
			r, f := behavior57Fixture(src, 2, 2)
			mut(&r, &f)
			got := parser.InspectImplementationBehavior(r, f)
			if len(got) != 1 || got[0].Status != "not_evaluated" || got[0].Reason == "" || len(got[0].Facts) != 0 {
				t.Fatalf("unbound identity observed: %+v", got)
			}
		})
	}
}
func TestBehavior57MethodIdentityAndReceiverWrite(t *testing.T) {
	src := "package sample\ntype A struct { n int }; type B struct { n int }\nfunc (a *A) Change() { a.n-- }\nfunc (b *B) Change() { b.n = 7 }\n"
	files := []domain.ExtractedFile{{Path: "sample.go", Content: src, Size: int64(len(src))}}
	r, _ := behavior57Fixture(src, 3, 3)
	o := behavior57Observed(t, parser.InspectImplementationBehavior(r, files))
	if len(o.Facts) != 1 || o.Facts[0].Target != "a.n" || o.Facts[0].Kind != "increment_decrement" {
		t.Fatalf("receiver write lost: %+v", o)
	}
	r.Hits[0].EndLine = 4
	o2 := parser.InspectImplementationBehavior(r, files)
	if len(o2) != 1 || o2[0].Status != "not_evaluated" || len(o2[0].Facts) != 0 {
		t.Fatalf("homonym crossed receiver/range: %+v", o2)
	}
	r.Hits[0].Line = 4
	o = behavior57Observed(t, parser.InspectImplementationBehavior(r, files))
	if len(o.Facts) != 1 || o.Facts[0].Target != "b.n" {
		t.Fatalf("wrong receiver: %+v", o)
	}
}
func TestBehavior57ValidationReplacesStaleProjection(t *testing.T) {
	src := "package sample\nfunc Change(input []int) []int { return input }\n"
	r, f := behavior57Fixture(src, 2, 2)
	r.Behavior = []domain.ImplementationBehavior{{Status: "verified", Facts: []domain.BehaviorFact{{Kind: "stale"}}}}
	before, _ := json.Marshal(r)
	got := usecase.ValidateImplementationRequirements("change", r, f)
	b := behavior57Observed(t, got.Behavior)
	if len(b.Facts) != 1 || b.Facts[0].Kind != "return" {
		t.Fatalf("stale projection persisted: %+v", b)
	}
	after, _ := json.Marshal(r)
	if string(before) != string(after) {
		t.Fatal("validator mutated input")
	}
	got = usecase.ValidateImplementationRequirements("change", got, nil)
	for _, b := range got.Behavior {
		if b.Status == "observed" || len(b.Facts) > 0 {
			t.Fatalf("stale facts survived missing bytes: %+v", b)
		}
	}
}
