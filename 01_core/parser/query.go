/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"strings"
)

var defaultStopWords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true, "is": true, "are": true,
	"in": true, "on": true, "at": true, "to": true, "for": true, "with": true, "by": true,
	"of": true, "from": true, "about": true, "simple": true, "open": true, "source": true,
	"free": true, "best": true, "great": true, "tool": true, "app": true, "application": true,
	"de": true, "do": true, "da": true, "para": true, "com": true, "um": true, "uma": true,
}

// NormalizeQuery sanitiza a entrada do usuário removendo pontuações, stop words e espaços duplicados.
func NormalizeQuery(raw string) string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '/' || r == '@' || r == '.' {
			return r
		}
		return ' '
	}, raw)

	words := strings.Fields(cleaned)
	filtered := make([]string, 0, len(words))

	for _, w := range words {
		lower := strings.ToLower(w)
		if !defaultStopWords[lower] || len(words) <= 2 {
			filtered = append(filtered, w)
		}
	}

	if len(filtered) == 0 {
		return strings.Join(words, " ")
	}

	return strings.Join(filtered, " ")
}

// ExpandKeywords expande conceitos conhecidos de arquitetura para termos técnicos específicos do ecossistema.
func ExpandKeywords(query string, lang string) []string {
	lower := strings.ToLower(query)
	keywords := []string{query}

	// Mapeamentos conceituais agnósticos e por linguagem
	if strings.Contains(lower, "hexagonal") || strings.Contains(lower, "clean") {
		keywords = append(keywords, "domain usecase ports adapters")
	}
	if strings.Contains(lower, "time") || strings.Contains(lower, "tracking") {
		keywords = append(keywords, "timer timesheet tracking")
	}
	if strings.Contains(lower, "agent") || strings.Contains(lower, "harness") {
		keywords = append(keywords, "agent llm prompt tools")
	}

	return keywords
}
