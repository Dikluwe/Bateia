/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/regex_validator.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"regexp"
	"strings"
)

// IsSafeRegex analisa se uma expressão regular possui padrões propensos a retrocesso catastrófico (ReDoS).
func IsSafeRegex(pattern string) bool {
	if pattern == "" {
		return true
	}

	// Padrões perigosos típicos de ReDoS: (a+)+ ou (a*)* ou (.*)*
	dangerousSubpatterns := []string{
		"(a+)+", "(a*)*", "(.*)*", "(.+)+", "([a-z]+)+", "([a-z]*)*",
	}

	lowerPat := strings.ToLower(pattern)
	for _, dangerous := range dangerousSubpatterns {
		if strings.Contains(lowerPat, dangerous) {
			return false
		}
	}

	// Tentar compilar com a engine nativa do Go (RE2 é garantida O(n), mas verificamos integridade)
	_, err := regexp.Compile(pattern)
	return err == nil
}
