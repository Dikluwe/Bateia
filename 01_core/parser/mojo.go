/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/mojo.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"strings"
	"time"

	"bateia/01_core/domain"
)

// ParseMojoProjectToml analisa o conteúdo de mojoproject.toml ou pixi.toml (Mojo).
func ParseMojoProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "mojo",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	manifest.ModuleName = ref.Repo
	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name =") {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				manifest.ModuleName = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}

		if strings.HasPrefix(trimmed, "[dependencies]") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "[") {
			inDeps = false
		}

		if inDeps {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				depName := strings.TrimSpace(parts[0])
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
