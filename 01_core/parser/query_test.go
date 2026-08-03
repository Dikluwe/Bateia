/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"testing"
)

func TestNormalizeQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "a simple open-source time tracking tool",
			expected: "open-source time tracking",
		},
		{
			input:    "Clean Architecture in Go!",
			expected: "Clean Architecture Go",
		},
		{
			input:    "solidtime",
			expected: "solidtime",
		},
	}

	for _, tt := range tests {
		got := NormalizeQuery(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeQuery(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}

func TestExpandKeywords(t *testing.T) {
	keywords := ExpandKeywords("time tracking", "go")
	if len(keywords) < 2 {
		t.Errorf("Esperava pelo menos 2 keywords expandidas, obteve %d", len(keywords))
	}
}
