/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/pascal.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	pascalPackageRegex = regexp.MustCompile(`(?i)<ItemName>([^<]+)</ItemName>`)
)

// ParsePascalDelphi analisa o conteúdo de um arquivo .dproj ou .lpi (Pascal/Delphi).
func ParsePascalDelphi(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "pascal",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range pascalPackageRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			pkgName := match[1]
			if !seen[pkgName] {
				seen[pkgName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     pkgName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
