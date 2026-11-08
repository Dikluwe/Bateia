/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/julia.md
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
	juliaNameRegex = regexp.MustCompile(`(?m)^\s*name\s*=\s*["']([^"']+)["']`)
	juliaDepRegex  = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-]+)\s*=\s*["']([^"']+)["']`)
)

// ParseJuliaProjectToml analisa o conteúdo de um Project.toml do ecossistema Julia.
func ParseJuliaProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "julia",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := juliaNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "[deps]") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "[") {
			inDeps = false
		}

		if inDeps {
			if match := juliaDepRegex.FindStringSubmatch(trimmed); len(match) >= 3 {
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
