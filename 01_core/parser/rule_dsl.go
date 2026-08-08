/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/rule_dsl.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "strings"

// DSLRule representa uma regra declarativa de validação arquitetural ou sintática.
type DSLRule struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Pattern  string `json:"pattern"`
}

// EvaluateDSL executa as regras declarativas sobre o texto fornecido e retorna os nomes das regras violadas.
func EvaluateDSL(rules []DSLRule, text string) []string {
	if text == "" || len(rules) == 0 {
		return nil
	}

	lowerText := strings.ToLower(text)
	var violated []string

	for _, rule := range rules {
		if rule.Pattern == "" {
			continue
		}
		if strings.Contains(lowerText, strings.ToLower(rule.Pattern)) {
			violated = append(violated, rule.Name)
		}
	}

	return violated
}
