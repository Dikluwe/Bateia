/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/regex_validator_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestIsSafeRegex(t *testing.T) {
	safePatterns := []string{
		`^[a-zA-Z0-9_-]+$`,
		`func\s+([a-zA-Z0-9_]+)`,
		`import\s+"[^"]+"`,
	}

	for _, p := range safePatterns {
		if !IsSafeRegex(p) {
			t.Errorf("IsSafeRegex(%q) = false; esperava true", p)
		}
	}

	unsafePatterns := []string{
		`(a+)+`,
		`(.*)*`,
		`[invalid regex (`,
	}

	for _, p := range unsafePatterns {
		if IsSafeRegex(p) {
			t.Errorf("IsSafeRegex(%q) = true; esperava false", p)
		}
	}
}
