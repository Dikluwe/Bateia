/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/nil_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"regexp"
	"strings"

	"bateia/01_core/domain"
)

var (
	// Detecta desreferenciação direta como ptr.Field ou *ptr sem verificação de nil prévia
	uncheckedPointerRegex = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\s*\(`)
)

// EvaluateNilPointerRisk analisa sintaticamente o código de um arquivo em busca de possíveis desreferenciações diretas de ponteiro.
func EvaluateNilPointerRisk(filename string, content string) []domain.ArchViolation {
	var violations []domain.ArchViolation

	if content == "" {
		return violations
	}

	lines := strings.Split(content, "\n")
	hasNilCheck := false

	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "if ") && (strings.Contains(trimmed, "== nil") || strings.Contains(trimmed, "!= nil") || strings.Contains(trimmed, ".is_some()")) {
			hasNilCheck = true
		}

		if uncheckedPointerRegex.MatchString(trimmed) && !hasNilCheck && idx > 5 {
			matches := uncheckedPointerRegex.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				varName := matches[1]
				if varName != "fmt" && varName != "log" && varName != "os" && varName != "strings" && varName != "filepath" && varName != "json" {
					violations = append(violations, domain.ArchViolation{
						Rule:        "NIL_POINTER_RISK",
						Severity:    "MEDIUM",
						Message:     "Risco de pânico por ponteiro nulo (Nil Pointer Dereference) detectado sem verificação prévia de nil.",
						Target:      filename,
						Line:        idx + 1,
						CodeSnippet: trimmed,
						Suggestion:  "Adicione uma verificação prévia de nil (ex: if ptr != nil) antes de acessar seus campos.",
					})
					// Limita a 1 violação por arquivo para evitar spam
					break
				}
			}
		}
	}

	return violations
}
