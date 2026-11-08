/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/v.md
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
	vNameRegex = regexp.MustCompile(`(?m)^\s*name\s*:\s*["']([^"']+)["']`)
	vDepRegex  = regexp.MustCompile(`["']([a-zA-Z0-9_\-]+)["']`)
)

// ParseVMod analisa o conteúdo de um arquivo v.mod do ecossistema V (Vlang).
func ParseVMod(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "v",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := vNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "dependencies:") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "]") && inDeps {
			inDeps = false
		}

		if inDeps {
			for _, match := range vDepRegex.FindAllStringSubmatch(trimmed, -1) {
				if len(match) >= 2 {
					depName := match[1]
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

	return manifest, nil
}
