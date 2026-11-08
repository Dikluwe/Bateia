/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/basic.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"
	"strings"

	"bateia/01_core/domain"
)

var (
	basicRefRegex = regexp.MustCompile(`(?i)Reference=\*\\G([^\r\n#]+)`)
)

// ParseBasicVbp analisa o conteúdo de um manifesto .vbp (BASIC / Visual Basic).
func ParseBasicVbp(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "basic",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range basicRefRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			refPath := match[1]
			parts := strings.Split(refPath, "#")
			refName := parts[len(parts)-1]

			if !seen[refName] {
				seen[refName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     refName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
