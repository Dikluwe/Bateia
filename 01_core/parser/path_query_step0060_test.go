/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/path_query_step0060_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"reflect"
	"sort"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

func path60Query(t *testing.T, input string) domain.ImplementationQuery {
	t.Helper()
	qs := parser.InterpretImplementationQuery(input)
	if len(qs) != 1 {
		t.Fatalf("interpretations: %+v", qs)
	}
	if qs[0].Original != input {
		t.Fatalf("original changed: %+v", qs[0])
	}
	return qs[0]
}

func path60Has(terms []string, term string) bool {
	for _, v := range terms {
		if v == term {
			return true
		}
	}
	return false
}

func TestPath0060TranslationReceipts(t *testing.T) {
	for _, tc := range []struct {
		input   string
		sources []string
	}{
		{"caminho conteúdo", []string{"caminho"}},
		{"caminhos conteúdo", []string{"caminhos"}},
		{"Caminho conteúdo", []string{"caminho"}},
		{"Caminhos conteúdo", []string{"caminhos"}},
		{"caminhos caminho caminhos caminho conteúdo", []string{"caminho", "caminhos"}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			q := path60Query(t, tc.input)
			if q.Mode != "natural" || len(q.Symbols) != 0 || len(q.Constraints) != 0 {
				t.Fatalf("query identity: %+v", q)
			}
			if len(q.Expansions) != len(tc.sources) {
				t.Fatalf("receipts: %+v", q.Expansions)
			}
			sources := []string{}
			for _, e := range q.Expansions {
				if e.RuleID != "0060-v1:path:translation" || e.Kind != "translation" || e.Status != "hypothesis" || !reflect.DeepEqual(e.Terms, []string{"path"}) {
					t.Fatalf("invalid receipt: %+v", e)
				}
				sources = append(sources, e.Source)
			}
			sort.Strings(sources)
			if !reflect.DeepEqual(sources, tc.sources) {
				t.Fatalf("sources %v want %v", sources, tc.sources)
			}
			if !path60Has(q.Terms, "path") {
				t.Fatalf("translation absent: %v", q.Terms)
			}
			for _, s := range tc.sources {
				if !path60Has(q.Terms, s) {
					t.Errorf("original term %s lost: %v", s, q.Terms)
				}
			}
			if !sort.StringsAreSorted(q.Terms) {
				t.Errorf("unsorted terms: %v", q.Terms)
			}
			for i := 1; i < len(q.Terms); i++ {
				if q.Terms[i] == q.Terms[i-1] {
					t.Errorf("duplicate term: %v", q.Terms)
				}
			}
			for i := 0; i < 4; i++ {
				if next := path60Query(t, tc.input); !reflect.DeepEqual(q, next) {
					t.Fatal("non-deterministic query")
				}
			}
		})
	}
	a := parser.InterpretImplementationQuery("caminho conteúdo\ncaminhos conteúdo\ncaminho conteúdo")
	b := parser.InterpretImplementationQuery("caminhos conteúdo\ncaminho conteúdo")
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("line canonicalization: %+v / %+v", a, b)
	}
}

func TestPath0060ExcludedInputs(t *testing.T) {
	for _, input := range []string{
		"caminho", "caminhos", "`caminho`", "`caminhos`", "`caminho` conteúdo", "`outro` caminho conteúdo",
		"OutroSimbolo caminho conteúdo", "outro_simbolo caminhos conteúdo", "CaminhoLocal conteúdo", "caminho_local conteúdo",
		"descaminho conteúdo", "caminhox conteúdo", "encaminhar conteúdo", "path conteúdo", "paths conteúdo", "parar conteúdo",
	} {
		t.Run(input, func(t *testing.T) {
			q := path60Query(t, input)
			for _, e := range q.Expansions {
				if e.RuleID == "0060-v1:path:translation" || e.Source == "caminho" || e.Source == "caminhos" {
					t.Errorf("excluded query expanded: %+v", q)
				}
			}
			if input != "path conteúdo" && path60Has(q.Terms, "path") {
				t.Errorf("unlicensed translation: %+v", q)
			}
			if path60Has(q.Terms, "remove") {
				t.Errorf("invented parar translation: %+v", q)
			}
		})
	}
	for _, marker := range []string{"sem", "não", "nao", "nunca", "not", "without", "never"} {
		input := "caminho conteúdo " + marker + " mudança"
		q := path60Query(t, input)
		if len(q.Expansions) != 0 || path60Has(q.Terms, "path") {
			t.Errorf("negation expanded: %+v", q)
		}
		if len(q.Constraints) != 1 || q.Constraints[0].Marker != marker || q.Constraints[0].Text != input || q.Constraints[0].Status != "not_verified" {
			t.Errorf("constraint changed: %+v", q)
		}
	}
}

func TestPath0060PreservesLegacyAndDoesNotCascade(t *testing.T) {
	base := path60Query(t, "ordenar conteúdo")
	q := path60Query(t, "ordenar caminho conteúdo")
	for _, term := range base.Terms {
		if !path60Has(q.Terms, term) {
			t.Errorf("legacy term lost: %s", term)
		}
	}
	legacy := []domain.ImplementationQueryExpansion{}
	for _, e := range q.Expansions {
		if e.RuleID != "0060-v1:path:translation" {
			legacy = append(legacy, e)
		}
	}
	if !reflect.DeepEqual(legacy, base.Expansions) {
		t.Fatalf("legacy expansion changed: %+v / %+v", legacy, base.Expansions)
	}
	q = path60Query(t, "caminho path conteúdo")
	if len(q.Expansions) != 1 || q.Expansions[0].Source != "caminho" {
		t.Fatalf("derived expansion: %+v", q)
	}
	count := 0
	for _, term := range q.Terms {
		if term == "path" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("global translation not deduplicated: %+v", q)
	}
}

func TestPath0060NeutralBodyRetrievalRemainsLexical(t *testing.T) {
	src := "package sample\nfunc Process(value string) string {\n path := value\n return path\n}\n"
	files := []domain.ExtractedFile{{Path: "neutral/a.go", Content: src, Size: int64(len(src))}, {Path: "neutral/b.go", Content: "package sample\nfunc Other() int { return 7 }\n"}}
	for _, query := range []string{"caminho conteúdo", "caminhos conteúdo"} {
		report := parser.RankGoImplementations(query, files, 8)
		found := false
		for _, hit := range report.Hits {
			if hit.Path == "neutral/a.go" && hit.Symbol == "Process" {
				found = true
				if hit.Retrieval != "lexical" || !path60Has(hit.MatchedTerms, "path") {
					t.Errorf("translation promoted or absent: %+v", hit)
				}
			}
		}
		if !found {
			t.Fatalf("body-only path not recovered for %q: %+v", query, report)
		}
		if !reflect.DeepEqual(report.QueryInterpretation, parser.InterpretImplementationQuery(query)) {
			t.Error("report lost provenance")
		}
		reqs := parser.DescribeImplementationRequirements(query, report)
		if len(reqs) != 1 || reqs[0].Status != "evidence_located" {
			t.Fatalf("lexical requirement status: %+v", reqs)
		}
		for _, behavior := range parser.InspectImplementationBehavior(report, files) {
			if behavior.Status != "observed" && behavior.Status != "not_verified" {
				t.Errorf("behavior verified by translation: %+v", behavior)
			}
		}
	}
	for _, query := range []string{"`caminho`", "`caminhos`"} {
		report := parser.RankGoImplementations(query, files, 8)
		if len(report.Hits) != 0 {
			t.Errorf("code glossary acquired translation: %q %+v", query, report.Hits)
		}
	}
}
