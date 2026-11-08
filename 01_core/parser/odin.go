/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/odin.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	odinColRegex = regexp.MustCompile(`["']([a-zA-Z0-9_\-]+)["']\s*:\s*["']([^"']+)["']`)
)

// ParseOlsJson analisa o conteúdo de ols.json ou odin.json (Odin).
func ParseOlsJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "odin",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range odinColRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 3 {
			colName := match[1]
			if !seen[colName] {
				seen[colName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     colName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
