/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/objc.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"regexp"
	"time"

	"bateia/01_core/domain"
)

var (
	podfileRegex = regexp.MustCompile(`(?i)pod\s+['"]([^'"]+)['"]`)
)

// ParsePodfile analisa o conteúdo de um Podfile (Objective-C).
func ParsePodfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "objective-c",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range podfileRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			podName := match[1]
			if !seen[podName] {
				seen[podName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     podName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
