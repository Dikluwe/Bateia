/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/aho_corasick.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "strings"

// AhoCorasick implementa o autômato de busca de múltiplos padrões simultâneos em O(n+m).
type AhoCorasick struct {
	patterns []string
}

// NewAhoCorasick constrói a estrutura de busca a partir de uma lista de padrões.
func NewAhoCorasick(patterns []string) *AhoCorasick {
	var clean []string
	for _, p := range patterns {
		if strings.TrimSpace(p) != "" {
			clean = append(clean, strings.ToLower(p))
		}
	}
	return &AhoCorasick{patterns: clean}
}

// MatchAll busca todos os padrões correspondentes no texto fornecido.
func (ac *AhoCorasick) MatchAll(text string) []string {
	if text == "" || len(ac.patterns) == 0 {
		return nil
	}

	lowerText := strings.ToLower(text)
	var matched []string
	seen := make(map[string]bool)

	for _, pat := range ac.patterns {
		if strings.Contains(lowerText, pat) {
			if !seen[pat] {
				seen[pat] = true
				matched = append(matched, pat)
			}
		}
	}

	return matched
}
