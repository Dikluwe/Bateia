/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/racket.md
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
	racketDepRegex = regexp.MustCompile(`["']([a-zA-Z0-9_\-]+)["']`)
)

// ParseInfoRkt analisa o conteúdo de info.rkt (Racket).
func ParseInfoRkt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "racket",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "deps") {
			for _, match := range racketDepRegex.FindAllStringSubmatch(trimmed, -1) {
				if len(match) >= 2 {
					depName := match[1]
					if depName != "base" && depName != "deps" && !seen[depName] {
						seen[depName] = true
						manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
							Name:     depName,
							Version:  "",
							Indirect: false,
							Type:     "direct",
						})
					}
				}
			}
		}
	}

	return manifest, nil
}
