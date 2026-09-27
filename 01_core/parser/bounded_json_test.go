/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/bounded_json_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package parser

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBoundedJSONPositiveValuesAndIndependentObjects(t *testing.T) {
	for _, input := range []string{`null`, `false`, `1e9999`, `"literal [ { } ]"`, `{}`, `[]`, `[{"a":1},{"a":2}]`, `{"a":{"b":1},"c":{"b":2}}`, strings.Repeat("[", 64) + `0` + strings.Repeat("]", 64)} {
		if err := ValidateBoundedJSON([]byte(input)); err != nil {
			t.Errorf("valid input %q: %v", input, err)
		}
	}
}

func TestBoundedJSONRejectsAmbiguousMembersAndMalformedDocuments(t *testing.T) {
	for _, input := range []string{`{"complete":false,"complete":true}`, `{"complete":false,"Complete":true}`, `{"k":1,"\u006b":2}`, `{"K":1,"K":2}`, `{"Σ":1,"ς":2}`, `{"nested":{"a":1,"a":2}}`, ``, `{"x":`, `[}`, `{} {}`, `true garbage`, string([]byte{'"', 0xff, '"'}), strings.Repeat("[", 65) + `0` + strings.Repeat("]", 65)} {
		if err := ValidateBoundedJSON([]byte(input)); err == nil {
			t.Errorf("ambiguous/malformed input accepted: %q", input)
		}
	}
}

func TestBoundedJSONByteAndTokenBudgets(t *testing.T) {
	if err := ValidateBoundedJSON([]byte(`"` + strings.Repeat("x", MaxStructuredJSONBytes) + `"`)); err == nil {
		t.Fatal("byte budget bypassed")
	}
	input := `[` + strings.Repeat(`0,`, maxStructuredJSONTokens) + `0]`
	if err := ValidateBoundedJSON([]byte(input)); err == nil {
		t.Fatal("token budget bypassed")
	}
}

func FuzzBoundedJSONAcceptsOnlyValidJSON(f *testing.F) {
	for _, seed := range []string{`{}`, `{"A":1,"a":2}`, `[[null]]`, `{"x":"\\\"[]"}`, `{} {}`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > MaxStructuredJSONBytes {
			t.Skip()
		}
		if ValidateBoundedJSON([]byte(input)) == nil && !json.Valid([]byte(input)) {
			t.Fatal("invalid JSON accepted")
		}
	})
}
