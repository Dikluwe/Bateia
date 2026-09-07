/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cyclomatic_lexical_test.md
 * @layer L1
 * @updated 2026-09-07
 */

package parser

import "testing"

// These regressions share authorship with the implementation.
func TestGoLexicalEstimateExcludesCommentsAndLiterals(t *testing.T) {
	cases := []string{
		"// if for case && ||\n",
		"/* if for case && || */",
		`var text = "if for case && ||"`,
		"var text = `if for case && ||`",
		`var letter = 'i'`,
		`var text = "escaped \" if &&"`,
	}
	for _, code := range cases {
		if got := EstimateGoLexicalComplexity(code); got != 1 {
			t.Errorf("literal/comment %q: got %d, want 1", code, got)
		}
	}
}

func TestGoLexicalEstimateCountsControlsAndOperators(t *testing.T) {
	cases := []struct {
		code string
		want int
	}{
		{"", 1},
		{"if a&&b||c {}", 4},
		{"for _, value := range items { _ = value }", 2},
		{"switch x { case 1: f(); case 2: g(); default: h() }", 3},
		{"if a {} else if b {} else {}", 3},
		{"// if fake {}\nif real { println(\"for &&\") }", 2},
	}
	for _, tc := range cases {
		if got := EstimateGoLexicalComplexity(tc.code); got != tc.want {
			t.Errorf("%q: got %d, want %d", tc.code, got, tc.want)
		}
		if got := CalculateCyclomaticComplexity(tc.code); got != tc.want {
			t.Errorf("legacy %q: got %d, want %d", tc.code, got, tc.want)
		}
	}
}
