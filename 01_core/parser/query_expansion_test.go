/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_expansion_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func query52(t *testing.T, s string) domain.ImplementationQuery {
	t.Helper()
	q := parser.InterpretImplementationQuery(s)
	if len(q) != 1 {
		t.Fatalf("single line: %#v", q)
	}
	if q[0].Original != s {
		t.Fatalf("original changed: %q", q[0].Original)
	}
	return q[0]
}
func has52(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
func TestExpansion0052CatalogAndProvenance(t *testing.T) {
	catalog := []struct{ triggers, translation, mechanism string }{
		{"descomprimir descompressão descompressao", "decompress decompression", "decode decoded decompresses"},
		{"comprimir compressão compressao", "compress compression", "encode encoded compresses"},
		{"serializar serialização serializacao", "serialize serialization", "marshal"},
		{"desserializar desserialização desserializacao", "deserialize deserialization", "unmarshal"},
		{"ordenar ordenação ordenacao", "sort sorting", ""}, {"comparar comparação comparacao", "compare comparison", ""},
		{"bloco blocos", "block blocks", ""}, {"corrompido corrompidos corrompida corrompidas", "corrupt corrupted", ""},
		{"dados", "data", ""}, {"destino", "destination", "dst"}, {"reutilizar", "reuse", ""}, {"rejeitar", "reject", ""},
	}
	ids := map[string]string{}
	for _, c := range catalog {
		family := strings.Fields(c.triggers)[0]
		for _, source := range strings.Fields(c.triggers) {
			for _, capital := range []bool{false, true} {
				word := source
				if capital {
					r := []rune(word)
					r[0] = []rune(strings.ToUpper(string(r[0])))[0]
					word = string(r)
				}
				q := query52(t, word+" conteúdo")
				expected := map[string][]string{"translation": strings.Fields(c.translation)}
				if c.mechanism != "" {
					expected["mechanism"] = strings.Fields(c.mechanism)
				}
				if len(q.Expansions) != len(expected) {
					t.Fatalf("%s receipts %#v", word, q.Expansions)
				}
				seen := map[string]bool{}
				for _, e := range q.Expansions {
					if e.Source != source || e.Status != "hypothesis" || !strings.HasPrefix(e.RuleID, "0052-v1") {
						t.Fatalf("invalid provenance: %#v", e)
					}
					want, ok := expected[e.Kind]
					sort.Strings(want)
					if !ok || seen[e.Kind] || !reflect.DeepEqual(e.Terms, want) {
						t.Fatalf("%s receipt %#v expected %#v", word, e, expected)
					}
					seen[e.Kind] = true
					key := family + "/" + e.Kind
					if old := ids[key]; old != "" && old != e.RuleID {
						t.Fatalf("unstable family rule %s", key)
					}
					ids[key] = e.RuleID
					for _, term := range want {
						if !has52(q.Terms, term) {
							t.Fatalf("%s missing %s", word, term)
						}
					}
				}
			}
		}
	}
	used := map[string]string{}
	for key, id := range ids {
		if old := used[id]; old != "" && old != key {
			t.Fatalf("same ID for distinct families/types: %s and %s", old, key)
		}
		used[id] = key
	}
}
func TestExpansion0052BoundariesAndDeterminism(t *testing.T) {
	for _, s := range []string{"descomprimir", "`descomprimir` conteúdo", "DescomprimirBloco conteúdo", "descomprimir_bloco conteúdo", "predescomprimir conteúdo", "descomprimirx conteúdo", "decode conteúdo", "decompression conteúdo", "inventário conteúdo"} {
		if q := query52(t, s); len(q.Expansions) != 0 {
			t.Errorf("must not expand %q: %#v", s, q.Expansions)
		}
	}
	for _, marker := range []string{"sem", "não", "nao", "nunca", "not", "without", "never"} {
		if q := query52(t, "descomprimir conteúdo "+marker+" alocar"); len(q.Expansions) != 0 {
			t.Errorf("negation %s expanded", marker)
		}
	}
	q := query52(t, "descomprimir descomprimir conteúdo")
	if len(q.Expansions) != 2 {
		t.Fatalf("duplicate/cascade receipts: %#v", q.Expansions)
	}
	for i := 0; i < 12; i++ {
		if next := query52(t, q.Original); !reflect.DeepEqual(q, next) {
			t.Fatal("nondeterminism")
		}
	}
	q = query52(t, "`descomprimir` ordenar conteúdo")
	if !has52(q.Symbols, "descomprimir") || len(q.Expansions) != 1 || q.Expansions[0].Source != "ordenar" {
		t.Fatalf("mixed symbol protection: %#v", q)
	}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err = json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	receipts, ok := raw["expansions"].([]any)
	if !ok || len(receipts) != 1 {
		t.Fatalf("json receipts %s", b)
	}
	for _, key := range []string{"source", "terms", "rule_id", "kind", "status"} {
		if _, ok := receipts[0].(map[string]any)[key]; !ok {
			t.Fatalf("missing json %s: %s", key, b)
		}
	}
}
