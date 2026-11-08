/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/elixir.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	elixirAppNameRegex = regexp.MustCompile(`app:\s*:([a-zA-Z0-9_]+)`)
	elixirDepRegex     = regexp.MustCompile(`\{:([a-zA-Z0-9_]+),\s*["']([^"']+)["']`)
)

// ParseMixExs analisa o conteúdo de um manifesto mix.exs do ecossistema Elixir.
func ParseMixExs(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "elixir",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := elixirAppNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	seen := make(map[string]bool)

	for _, match := range elixirDepRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 3 {
			depName := match[1]
			ver := match[2]
			if !seen[depName] {
				seen[depName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  ver,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
