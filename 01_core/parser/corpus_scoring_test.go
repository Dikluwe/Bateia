/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/corpus_scoring_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func corpus0047Files() []domain.ExtractedFile {
	files := []domain.ExtractedFile{{Path: "a.go", Content: "package sample\nfunc Alpha(){ quartz(); parcel() }"}}
	for i := 0; i < 8; i++ {
		files = append(files, domain.ExtractedFile{Path: fmt.Sprintf("d%d.go", i), Content: fmt.Sprintf("package sample\nfunc Other%d(){ parcel() }", i)})
	}
	return files
}

func TestCorpus0047RareContribution(t *testing.T) {
	r := parser.RankGoImplementations("quartz parcel", corpus0047Files(), 24)
	var rare, common *domain.ImplementationTermWeight
	for _, h := range r.Hits {
		if h.Path != "a.go" {
			continue
		}
		for _, w := range h.TermWeights {
			w := w
			if w.Term == "quartz" {
				rare = &w
			}
			if w.Term == "parcel" {
				common = &w
			}
		}
	}
	if rare == nil || common == nil {
		t.Fatalf("missing both contributions: %+v", r.Hits)
	}
	if rare.Field != common.Field || rare.DocumentFrequency >= common.DocumentFrequency || rare.Weight <= common.Weight {
		t.Fatalf("rare term must discriminate within same field: rare=%+v common=%+v", rare, common)
	}
}

func TestCorpus0047LengthNormalization(t *testing.T) {
	var body strings.Builder
	body.WriteString("package sample\nfunc Alpha(){ quartz();")
	for i := 0; i < 200; i++ {
		fmt.Fprintf(&body, "unrelated%d();", i)
	}
	body.WriteString("}")
	files := []domain.ExtractedFile{{Path: "a-long.go", Content: body.String()}, {Path: "z-short.go", Content: "package sample\nfunc Beta(){ quartz() }"}}
	r := parser.RankGoImplementations("quartz", files, 24)
	if len(r.Hits) < 2 || r.Hits[0].Path != "z-short.go" {
		t.Fatalf("concise relevant body should precede bloated body: %+v", r.Hits)
	}
}

func TestCorpus0047FiniteDeterministicDuplicate(t *testing.T) {
	files := corpus0047Files()
	r := parser.RankGoImplementations("quartz parcel", files, 24)
	count := 0
	for _, h := range r.Hits {
		for _, w := range h.TermWeights {
			count++
			if w.Term == "" || w.Field == "" || w.DocumentFrequency < 1 || w.CorpusDocuments < w.DocumentFrequency || math.IsNaN(w.Weight) || math.IsInf(w.Weight, 0) {
				t.Fatalf("invalid contribution: %+v", w)
			}
		}
	}
	if count == 0 {
		t.Fatal("no contribution exposed")
	}
	reordered := append([]domain.ExtractedFile{}, files...)
	for i, j := 0, len(reordered)-1; i < j; i, j = i+1, j-1 {
		reordered[i], reordered[j] = reordered[j], reordered[i]
	}
	reordered = append(reordered, files[0], files[1])
	got := parser.RankGoImplementations("quartz parcel", reordered, 24)
	if !reflect.DeepEqual(r.Hits, got.Hits) {
		t.Fatalf("order or duplicates changed hits/statistics:\n%+v\n%+v", r.Hits, got.Hits)
	}
}

func TestCorpus0047ExactAndASTOnly(t *testing.T) {
	files := []domain.ExtractedFile{{Path: "exact.go", Content: "package sample\nfunc Quartz(){}"}, {Path: "body.go", Content: "package sample\nfunc Alpha(){ Quartz() }"}, {Path: "decoy.go", Content: "package sample\n// quartz\nfunc Message() string { return \"quartz\" }"}}
	r := parser.RankGoImplementations("Quartz", files, 24)
	if len(r.Hits) == 0 || r.Hits[0].Path != "exact.go" {
		t.Fatalf("exact priority lost: %+v", r.Hits)
	}
	for _, h := range r.Hits {
		if h.Path == "decoy.go" {
			t.Fatalf("comment/literal ranked: %+v", h)
		}
	}
	absent := parser.RankGoImplementations("zephyr", []domain.ExtractedFile{{Path: "notes.go", Content: "package sample\n// zephyr\nfunc Message() string { return \"zephyr\" }"}}, 24)
	if len(absent.Hits) != 0 {
		t.Fatalf("literal-only term produced hits: %+v", absent.Hits)
	}
}
