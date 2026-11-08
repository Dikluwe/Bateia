/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/zig.md
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
	zigNameRegex = regexp.MustCompile(`(?m)^\s*\.name\s*=\s*["']([^"']+)["']`)
	zigDepRegex  = regexp.MustCompile(`(?m)^\s*\.([a-zA-Z0-9_\-]+)\s*=\s*\.\{`)
)

// ParseBuildZigZon analisa o conteúdo de um manifesto build.zig.zon do ecossistema Zig.
func ParseBuildZigZon(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "zig",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := zigNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ".dependencies") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "}") && inDeps {
			inDeps = false
		}

		if inDeps {
			if match := zigDepRegex.FindStringSubmatch(trimmed); len(match) >= 2 {
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

	return manifest, nil
}
