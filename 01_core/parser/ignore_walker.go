/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ignore_walker.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "strings"

// IgnoreWalker filtra nós de diretório usando padrões ignorados no estilo .gitignore.
type IgnoreWalker struct {
	patterns []string
}

// NewIgnoreWalker cria um IgnoreWalker com regras padrão (vendor, node_modules, target, .git).
func NewIgnoreWalker(customPatterns []string) *IgnoreWalker {
	defaults := []string{
		"vendor", "node_modules", "target", ".git", "build", "dist", ".cache",
	}
	all := append(defaults, customPatterns...)
	return &IgnoreWalker{patterns: all}
}

// ShouldIgnore retorna true se o caminho corresponder a qualquer regra ignorada.
func (w *IgnoreWalker) ShouldIgnore(path string) bool {
	if path == "" {
		return false
	}
	lower := strings.ToLower(path)
	for _, pat := range w.patterns {
		if strings.Contains(lower, strings.ToLower(pat)) {
			return true
		}
	}
	return false
}
