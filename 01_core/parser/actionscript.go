/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/actionscript.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	flexClassRegex = regexp.MustCompile(`(?i)<path-element>([^<]+)</path-element>`)
)

// ParseFlexConfig analisa o conteúdo de um arquivo flex-config.xml (ActionScript / Flex).
func ParseFlexConfig(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "actionscript",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range flexClassRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			elemPath := match[1]
			if !seen[elemPath] {
				seen[elemPath] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     elemPath,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
