/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cyclomatic.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"regexp"
	"strings"
)

var (
	// Palavras-chave de bifurcação sintática para cálculo de v(G) de McCabe
	decisionRegex = regexp.MustCompile(`\b(if|else if|for|range|while|switch|case|catch|&&|\|\|)\b`)
)

// CalculateCyclomaticComplexity computa a complexidade ciclomática de McCabe v(G) = E - N + 2P (estimada por nós de decisão + 1).
func CalculateCyclomaticComplexity(code string) int {
	if strings.TrimSpace(code) == "" {
		return 1
	}

	matches := decisionRegex.FindAllString(code, -1)
	return len(matches) + 1
}
