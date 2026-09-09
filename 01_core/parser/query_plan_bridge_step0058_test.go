/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_plan_bridge_step0058_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"reflect"
	"strings"
	"testing"
)

func bridge58Query(t *testing.T, original string) domain.ImplementationQuery {
	t.Helper()
	got := parser.InterpretImplementationQuery(original)
	if len(got) != 1 || got[0].Original != original {
		t.Fatalf("single original line lost: %+v", got)
	}
	return got[0]
}

func bridge58Receipts(q domain.ImplementationQuery, prefix string) map[string]domain.ImplementationQueryExpansion {
	out := map[string]domain.ImplementationQueryExpansion{}
	for _, e := range q.Expansions {
		if strings.HasPrefix(e.RuleID, prefix) {
			out[e.RuleID] = e
		}
	}
	return out
}

func bridge58Has(xs []string, term string) bool {
	for _, x := range xs {
		if x == term {
			return true
		}
	}
	return false
}

func TestBridge0058PlanProvenance(t *testing.T) {
	for _, original := range []string{
		"Preciso monitorar arquivos mantendo ordem peculiar?",
		"interpolar cores",
		"serializar arquivos",
		"comprimir dados",
		"união bits",
		"comprimir compressão dados data",
	} {
		t.Run(original, func(t *testing.T) {
			plan := parser.PlanDiscoveryQuery(original)
			if plan.Status != "hypothesis" {
				t.Fatalf("stable prerequisite: %+v", plan)
			}
			q := bridge58Query(t, original)
			if q.Mode != "natural" || len(q.Symbols) != 0 || len(q.Constraints) != 0 {
				t.Fatalf("changed query interpretation: %+v", q)
			}
			want := map[string]domain.ImplementationQueryExpansion{}
			for _, rule := range plan.RuleIDs {
				var terms []string
				if strings.HasPrefix(rule, "action.") {
					terms = plan.ActionTerms
				}
				if strings.HasPrefix(rule, "object.") {
					terms = plan.ObjectTerms
				}
				if terms == nil {
					continue
				}
				id := "0058-v1:" + plan.Version + ":" + rule
				want[id] = domain.ImplementationQueryExpansion{Source: original, Terms: terms, RuleID: id, Kind: "translation", Status: "hypothesis"}
			}
			got := bridge58Receipts(q, "0058-")
			count := 0
			for _, e := range q.Expansions {
				if strings.HasPrefix(e.RuleID, "0058-") {
					count++
				}
			}
			if len(want) != 2 || count != 2 || !reflect.DeepEqual(got, want) {
				t.Fatalf("exact action/object provenance got=%+v want=%+v count=%d", got, want, count)
			}
			for _, group := range [][]string{plan.ActionTerms, plan.ObjectTerms} {
				for _, term := range group {
					if !bridge58Has(q.Terms, term) {
						t.Errorf("missing %q in %v", term, q.Terms)
					}
				}
			}
			seen := map[string]bool{}
			for _, term := range q.Terms {
				if seen[term] {
					t.Errorf("duplicate term %q", term)
				}
				seen[term] = true
			}
		})
	}
}

func TestBridge0058Boundaries(t *testing.T) {
	for _, original := range []string{
		"comprimir serializar dados", "comprimir dados arquivos", "comprimir", "dados",
		"comprimir dados repo:org/x", "comprimir dados && x", `comprimir "dados"`,
		"comprimir dados OR x", "comprimir dados (x)", "descomprimir dados",
	} {
		if p := parser.PlanDiscoveryQuery(original); p.Status != "literal_only" {
			t.Fatalf("prerequisite %q %+v", original, p)
		}
		if q := bridge58Query(t, original); len(bridge58Receipts(q, "0058-")) != 0 {
			t.Errorf("literal plan expanded: %+v", q)
		}
	}
	for _, marker := range []string{"sem", "não", "nao", "nunca", "not", "without", "never"} {
		original := "monitorar arquivos " + marker + " atrasar"
		q := bridge58Query(t, original)
		if len(bridge58Receipts(q, "0058-")) != 0 || len(q.Constraints) != 1 {
			t.Fatalf("negation expanded or lost: %+v", q)
		}
		want := domain.ImplementationQueryConstraint{Marker: marker, Text: original, Status: "not_verified"}
		if !reflect.DeepEqual(q.Constraints[0], want) {
			t.Errorf("constraint promoted: %+v", q.Constraints)
		}
	}
	for _, tc := range []struct{ text, mode, symbol string }{
		{"monitorar arquivos WatchDirectory", "mixed", "WatchDirectory"},
		{"monitorar arquivos watch_directory", "mixed", "watch_directory"},
		{"monitorar arquivos `watch`", "mixed", "watch"},
		{"WatchDirectory", "symbol", "WatchDirectory"},
		{"`watch`", "symbol", "watch"},
	} {
		q := bridge58Query(t, tc.text)
		if q.Mode != tc.mode || !reflect.DeepEqual(q.Symbols, []string{tc.symbol}) || len(bridge58Receipts(q, "0058-")) != 0 {
			t.Errorf("mode/symbol protection: %+v", q)
		}
	}
	for _, q := range parser.InterpretImplementationQuery("monitorar\narquivos") {
		if len(bridge58Receipts(q, "0058-")) != 0 {
			t.Errorf("cross-line expansion: %+v", q)
		}
	}
}

func TestBridge0058Preserves0052(t *testing.T) {
	for _, pair := range [][2]string{{"serializar", "dados"}, {"comprimir", "dados"}} {
		want := map[string]domain.ImplementationQueryExpansion{}
		for _, word := range pair {
			for id, e := range bridge58Receipts(bridge58Query(t, word+" conteúdo"), "0052-") {
				want[id] = e
			}
		}
		q := bridge58Query(t, strings.Join(pair[:], " "))
		if len(want) != 3 || !reflect.DeepEqual(bridge58Receipts(q, "0052-"), want) {
			t.Fatalf("0052 receipts changed: got=%+v want=%+v", q.Expansions, want)
		}
	}
}

func TestBridge0058DeterminismAndOwnership(t *testing.T) {
	a := "monitorar arquivos"
	b := "serializar dados"
	want := parser.InterpretImplementationQuery(a + "\n" + b)
	for i := 0; i < 12; i++ {
		got := parser.InterpretImplementationQuery("  " + b + " \n" + a + "\n" + a + "\n")
		if !reflect.DeepEqual(got, want) {
			t.Fatal("line permutation/dedup/repetition changed output")
		}
	}
	mutated := parser.InterpretImplementationQuery(a + "\n" + b)
	for i := range mutated {
		if len(mutated[i].Terms) > 0 {
			mutated[i].Terms[0] = "tampered"
		}
		for j := range mutated[i].Expansions {
			if len(mutated[i].Expansions[j].Terms) > 0 {
				mutated[i].Expansions[j].Terms[0] = "tampered"
			}
		}
	}
	if got := parser.InterpretImplementationQuery(a + "\n" + b); !reflect.DeepEqual(got, want) {
		t.Fatal("caller mutation leaked into later query")
	}
}

func TestBridge0058IndependentRankingFixture(t *testing.T) {
	files := []domain.ExtractedFile{{Path: "neutral/operations.go", Content: "package neutral\nfunc WatchDirectory(directory string) string { return directory }\nfunc BlendColour(colour int) int { return colour }\nfunc MarshalDirectory(directory string) []byte { return []byte(directory) }\nfunc CountRows(rows []int) int { return len(rows) }\n"}}
	before := append([]domain.ExtractedFile(nil), files...)
	for _, tc := range []struct{ query, symbol, object string }{
		{"monitorar arquivos", "WatchDirectory", "directory"},
		{"interpolar cores", "BlendColour", "colour"},
		{"serializar arquivos", "MarshalDirectory", "directory"},
	} {
		t.Run(tc.query, func(t *testing.T) {
			plan := parser.PlanDiscoveryQuery(tc.query)
			if plan.Status != "hypothesis" {
				t.Fatalf("missing stable plan: %+v", plan)
			}
			terms := append(append([]string(nil), plan.ActionTerms...), plan.ObjectTerms...)
			for _, query := range []string{strings.Join(terms, " "), tc.query} {
				r := parser.RankGoImplementations(query, files, 100)
				found := false
				for _, h := range r.Hits {
					if h.Path != files[0].Path || h.Symbol != tc.symbol {
						continue
					}
					found = true
					if !bridge58Has(h.MatchedTerms, tc.object) {
						t.Errorf("object absent from target evidence: %+v", h)
					}
					action := false
					for _, term := range plan.ActionTerms {
						if bridge58Has(h.MatchedTerms, term) {
							action = true
						}
					}
					if !action {
						t.Errorf("action absent from target evidence: %+v", h)
					}
				}
				if !found {
					t.Errorf("query=%q missing target %s: %+v", query, tc.symbol, r.Hits)
				}
				if query == tc.query && !reflect.DeepEqual(r.QueryInterpretation, parser.InterpretImplementationQuery(query)) {
					t.Error("ranking lost query provenance")
				}
			}
		})
	}
	if !reflect.DeepEqual(files, before) {
		t.Fatal("ranking mutated caller fixture")
	}
}
