/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/html_selector.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractHTMLTags extrai os conteúdos das tags especificadas em um texto HTML/XML sem dependências externas em L1.
func ExtractHTMLTags(htmlContent string, tagName string) []string {
	if htmlContent == "" || tagName == "" {
		return nil
	}

	pattern := fmt.Sprintf(`(?is)<%s\b[^>]*>(.*?)</%s>`, regexp.QuoteMeta(tagName), regexp.QuoteMeta(tagName))
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(htmlContent, -1)
	var results []string

	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, strings.TrimSpace(match[1]))
		}
	}

	return results
}
