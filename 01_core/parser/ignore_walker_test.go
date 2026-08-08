/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ignore_walker.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestIgnoreWalker_ShouldIgnore(t *testing.T) {
	walker := NewIgnoreWalker([]string{"tmp"})

	tests := []struct {
		path     string
		expected bool
	}{
		{"vendor/github.com/foo", true},
		{"node_modules/express", true},
		{"target/debug/app", true},
		{"src/main.go", false},
		{"tmp/cache.json", true},
	}

	for _, tt := range tests {
		res := walker.ShouldIgnore(tt.path)
		if res != tt.expected {
			t.Errorf("ShouldIgnore(%q) = %v; esperava %v", tt.path, res, tt.expected)
		}
	}
}
