/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_interpretation_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

func TestQuery0048OperatorsAndExplicitSymbols(t *testing.T) {
	cases := []struct {
		query, mode                      string
		symbols, markers, wanted, absent []string
	}{
		{"cache sem lock", "natural", nil, []string{"sem"}, []string{"cache", "lock"}, []string{"sem"}},
		{"cache NÃO lock", "natural", nil, []string{"não"}, []string{"cache", "lock"}, []string{"não"}},
		{"cache nao nunca not without never lock", "natural", nil, []string{"nao", "nunca", "not", "without", "never"}, []string{"cache", "lock"}, []string{"nao", "nunca", "not", "without", "never"}},
		{"sem", "symbol", []string{"sem"}, nil, []string{"sem"}, nil},
		{"not", "symbol", []string{"not"}, nil, []string{"not"}, nil},
		{"without", "symbol", []string{"without"}, nil, []string{"without"}, nil},
		{"`sem`", "symbol", []string{"sem"}, nil, []string{"sem"}, nil},
		{"cache `sem` lock", "mixed", []string{"sem"}, nil, []string{"sem", "cache", "lock"}, nil},
		{"cache CacheWithoutLock sem lock", "mixed", []string{"CacheWithoutLock"}, []string{"sem"}, []string{"cache", "lock", "without"}, []string{"sem"}},
		{"cache cache_without_lock sem lock", "mixed", []string{"cache_without_lock"}, []string{"sem"}, []string{"cache", "lock", "without"}, []string{"sem"}},
		{"cache `sem lock", "natural", nil, []string{"sem"}, []string{"cache", "lock"}, []string{"sem"}},
		{"cache `sem lock`", "natural", nil, []string{"sem"}, []string{"cache", "lock"}, []string{"sem"}},
		{"sem not without", "natural", nil, []string{"sem", "not", "without"}, nil, []string{"sem", "not", "without"}},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			got := parser.InterpretImplementationQuery("  " + tc.query + "  ")
			if len(got) != 1 {
				t.Fatalf("interpretations=%v", got)
			}
			q := got[0]
			if q.Original != tc.query || q.Mode != tc.mode {
				t.Errorf("original/mode=%q/%q", q.Original, q.Mode)
			}
			if !query0048Equal(q.Symbols, tc.symbols) {
				t.Errorf("symbols=%v want=%v", q.Symbols, tc.symbols)
			}
			markers := []string{}
			for _, c := range q.Constraints {
				markers = append(markers, strings.ToLower(c.Marker))
				if c.Text != tc.query || c.Status != "not_verified" {
					t.Errorf("constraint=%+v", c)
				}
			}
			sort.Strings(markers)
			wantMarkers := append([]string(nil), tc.markers...)
			sort.Strings(wantMarkers)
			if !query0048Equal(markers, wantMarkers) {
				t.Errorf("markers=%v want=%v", markers, wantMarkers)
			}
			for _, term := range tc.wanted {
				if !query0048Contains(q.Terms, term) {
					t.Errorf("missing term %q in %v", term, q.Terms)
				}
			}
			for _, term := range tc.absent {
				if query0048Contains(q.Terms, term) {
					t.Errorf("operator %q leaked into terms %v", term, q.Terms)
				}
			}
		})
	}
}

func TestQuery0048CanonicalLinesAndJSON(t *testing.T) {
	a := parser.InterpretImplementationQuery("cache sem lock\n\n `sem` \ncache sem lock\n")
	b := parser.InterpretImplementationQuery("`sem`\ncache sem lock")
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("line reorder/dedup changed interpretation: %v / %v", a, b)
	}
	if len(a) != 2 || a[0].Original != "`sem`" {
		t.Fatalf("not sorted: %v", a)
	}
	symbols := parser.InterpretImplementationQuery("cache ZebraLock AlphaLock ZebraLock")[0].Symbols
	if !query0048Equal(symbols, []string{"AlphaLock", "ZebraLock"}) {
		t.Fatalf("symbols not sorted/deduped: %v", symbols)
	}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"original":`, `"mode":`, `"terms":`, `"symbols":`, `"constraints":`, `"marker":`, `"text":`, `"status":"not_verified"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("JSON missing %s: %s", key, data)
		}
	}
}

func TestQuery0048RankingAndRequirementProjection(t *testing.T) {
	files := []domain.ExtractedFile{{Path: "fixture.go", Content: "package fixture\nfunc sem() {}\nfunc CacheLock() {}\nfunc Other() {}\n"}}
	natural := parser.RankGoImplementations("cache sem lock", files, 12)
	positive := parser.RankGoImplementations("cache lock", files, 12)
	if !reflect.DeepEqual(natural.Hits, positive.Hits) {
		t.Errorf("natural operator changed ranking evidence: %+v versus %+v", natural.Hits, positive.Hits)
	}
	if len(natural.QueryInterpretation) != 1 || len(natural.QueryInterpretation[0].Constraints) != 1 {
		t.Errorf("report lost interpretation: %+v", natural.QueryInterpretation)
	}
	for _, h := range natural.Hits {
		for _, w := range h.TermWeights {
			if w.Term == "sem" {
				t.Error("natural operator received weight")
			}
		}
		if query0048Contains(h.MatchedTerms, "sem") {
			t.Error("natural operator matched")
		}
	}
	for _, query := range []string{"sem", "`sem`"} {
		report := parser.RankGoImplementations(query, files, 12)
		found := false
		for _, h := range report.Hits {
			if h.Symbol == "sem" && query0048Contains(h.MatchedTerms, "sem") {
				found = true
			}
		}
		if !found {
			t.Errorf("explicit symbol unavailable for %q: %+v", query, report.Hits)
		}
	}
	requirements := parser.DescribeImplementationRequirements("cache sem lock", natural)
	if len(requirements) != 1 {
		t.Fatalf("requirements=%v", requirements)
	}
	r := requirements[0]
	if r.Status != "evidence_located" || len(r.Sources) == 0 {
		t.Errorf("positive evidence not located: %+v", r)
	}
	if len(r.Constraints) != 1 || r.Constraints[0].Status != "not_verified" || r.Constraints[0].Text != "cache sem lock" {
		t.Errorf("constraint promoted or lost: %+v", r)
	}
	for _, query := range []string{"sem not without", "cache sem lock"} {
		noSource := parser.DescribeImplementationRequirements(query, domain.ImplementationSearchReport{})
		if len(noSource) != 1 || noSource[0].Status != "not_verified" || len(noSource[0].Constraints) == 0 {
			t.Errorf("empty old report incorrectly projected: %+v", noSource)
		}
	}
	old := natural
	old.QueryInterpretation = nil
	oldReq := parser.DescribeImplementationRequirements("cache sem lock", old)
	if !reflect.DeepEqual(oldReq, requirements) {
		t.Errorf("old report projection lost query semantics: %+v", oldReq)
	}
}

func query0048Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func query0048Contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func TestQuery0048ExplicitShortAndStopwordSymbols(t *testing.T) {
	files := []domain.ExtractedFile{{Path: "explicit.go", Content: "package fixture\nfunc s() {}\nfunc use() {}\n"}}
	for _, name := range []string{"s", "use"} {
		for _, query := range []string{name, "`" + name + "`"} {
			t.Run(query, func(t *testing.T) {
				interpreted := parser.InterpretImplementationQuery(query)
				if len(interpreted) != 1 {
					t.Fatalf("interpretation=%+v", interpreted)
				}
				q := interpreted[0]
				if q.Mode != "symbol" || !query0048Equal(q.Symbols, []string{name}) || !query0048Contains(q.Terms, name) {
					t.Errorf("explicit identifier token lost: %+v", q)
				}
				report := parser.RankGoImplementations(query, files, 12)
				exact := false
				for _, hit := range report.Hits {
					if hit.Path == "explicit.go" && hit.Symbol == name {
						exact = true
					}
				}
				if !exact {
					t.Errorf("explicit identifier has no exact hit: %+v", report.Hits)
				}
				requirements := parser.DescribeImplementationRequirements(query, report)
				if len(requirements) != 1 {
					t.Fatalf("requirements=%+v", requirements)
				}
				requirement := requirements[0]
				found := false
				for _, source := range requirement.Sources {
					if source.Path == "explicit.go" && source.Symbol == name {
						found = true
					}
				}
				if requirement.Status != "evidence_located" || !found || !query0048Contains(requirement.MatchedTerms, name) {
					t.Errorf("explicit identifier has no projected exact evidence: %+v", requirement)
				}
			})
		}
	}
}
