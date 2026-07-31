/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/erlang.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
)

var (
	erlangDepRegex = regexp.MustCompile(`\{([a-zA-Z0-9_\-]+),\s*["']([^"']*)["']`)
)

// ParseRebarConfig analisa o conteúdo de um manifesto rebar.config do ecossistema Erlang.
func ParseRebarConfig(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "erlang",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	inDeps := false
	seen := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{deps,") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "]}") && inDeps {
			inDeps = false
		}

		if inDeps {
			for _, match := range erlangDepRegex.FindAllStringSubmatch(trimmed, -1) {
				if len(match) >= 2 {
					depName := match[1]
					if !seen[depName] {
						seen[depName] = true
						manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
							Name:     depName,
							Version:  match[2],
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
