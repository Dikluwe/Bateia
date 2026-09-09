/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/implementation_search_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
)

func TestImplementationSearchExactGenericAndLexical(t *testing.T) {
	source := "package sample\n\nfunc GatherItems[T any](items []T) []T { return items }\n"
	files := []domain.ExtractedFile{{Path: "neutral/a.go", Content: source}, {Path: "GatherItems.go", Content: "package sample\n// GatherItems is not implemented here\nconst note = `GatherItems`\nfunc Other() {}\n"}, {Path: "README.md", Content: "GatherItems exists"}}
	report := parser.RankGoImplementations("GatherItems", files, 8)
	if len(report.Hits) == 0 {
		t.Fatalf("missing exact definition: %+v", report)
	}
	hit := report.Hits[0]
	if hit.Path != "neutral/a.go" || hit.Symbol != "GatherItems" || hit.Line != 3 || hit.EndLine != 3 || hit.SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(source))) {
		t.Fatalf("wrong exact definition: %+v", hit)
	}
	for _, h := range report.Hits {
		if h.Symbol == "GatherItems" && h.Path != "neutral/a.go" {
			t.Fatalf("mention became definition: %+v", h)
		}
	}
	for _, tc := range []struct{ query, symbol, source string }{
		{"recipe ingredients", "RecipeIngredients", "package sample\nfunc RecipeIngredients() []string { return []string{\"flour\"} }\n"},
		{"enqueue queue items", "EnqueueItems", "package sample\nfunc EnqueueItems(queue chan int, item int) { queue <- item }\n"},
	} {
		t.Run(tc.symbol, func(t *testing.T) {
			r := parser.RankGoImplementations(tc.query, []domain.ExtractedFile{{Path: "misc/z.go", Content: tc.source}}, 4)
			if len(r.Hits) == 0 || r.Hits[0].Symbol != tc.symbol {
				t.Fatalf("lexical retrieval missing: %+v", r)
			}
		})
	}
}

func TestImplementationSearchOrderDedupAndLimit(t *testing.T) {
	a := domain.ExtractedFile{Path: "a.go", Content: "package sample\nfunc QueueOne() {}\n"}
	b := domain.ExtractedFile{Path: "b.go", Content: "package sample\nfunc QueueTwo() {}\n"}
	first := parser.RankGoImplementations("queue", []domain.ExtractedFile{a, b}, 8)
	second := parser.RankGoImplementations("queue", []domain.ExtractedFile{b, a, a}, 8)
	if len(first.Hits) != 2 || !reflect.DeepEqual(first.Hits, second.Hits) {
		t.Fatalf("unstable hits: %+v / %+v", first, second)
	}
	if !reflect.DeepEqual(first.Examined, second.Examined) || first.BytesRead != second.BytesRead || first.FilesRead != second.FilesRead {
		t.Fatalf("duplicate/order changed accounting: %+v / %+v", first, second)
	}
	limited := parser.RankGoImplementations("queue", []domain.ExtractedFile{a, b}, 1)
	if len(limited.Hits) != 1 {
		t.Fatalf("result cap not enforced: %+v", limited)
	}
}

func TestImplementationSearchUnknownAndScopedAbsence(t *testing.T) {
	for _, f := range []domain.ExtractedFile{{Path: "bad.go", Content: "package sample\nfunc Hidden( {"}, {Path: "a.py", Content: "def Hidden(): pass"}} {
		r := parser.RankGoImplementations("Hidden", []domain.ExtractedFile{f}, 4)
		if len(r.Hits) != 0 || r.Reason == "" || r.Status == "complete" {
			t.Fatalf("opaque input represented as complete: %+v", r)
		}
	}
	r := parser.RankGoImplementations("MissingSymbol", []domain.ExtractedFile{{Path: "a.go", Content: "package sample\nfunc Other() {}"}}, 4)
	if len(r.Hits) != 0 || r.Reason == "" {
		t.Fatalf("absence needs explicit scoped reason: %+v", r)
	}
}
