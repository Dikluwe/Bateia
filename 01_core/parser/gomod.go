/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/gomod.md
 * @layer L1
 * @updated 2026-07-25
 */

package parser

import (
	"strings"
	"time"

	"bateia/01_core/domain"
)

// ParseGoMod faz o parsing puro do conteúdo de um arquivo go.mod sem fazer I/O.
func ParseGoMod(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
	}

	lines := strings.Split(content, "\n")
	inRequireBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if strings.HasPrefix(trimmed, "module ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				manifest.ModuleName = parts[1]
			}
			continue
		}

		if strings.HasPrefix(trimmed, "go ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				manifest.GoVersion = parts[1]
			}
			continue
		}

		if trimmed == "require (" {
			inRequireBlock = true
			continue
		}

		if inRequireBlock && trimmed == ")" {
			inRequireBlock = false
			continue
		}

		if strings.HasPrefix(trimmed, "require ") && !inRequireBlock {
			dep := parseDependencyLine(strings.TrimPrefix(trimmed, "require "))
			if dep != nil {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
			continue
		}

		if inRequireBlock {
			dep := parseDependencyLine(trimmed)
			if dep != nil {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}
	}

	return manifest, nil
}

func parseDependencyLine(line string) *domain.Dependency {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	dep := &domain.Dependency{
		Name:     parts[0],
		Version:  parts[1],
		Indirect: false,
	}

	if len(parts) >= 3 && strings.Contains(line, "// indirect") {
		dep.Indirect = true
	}

	return dep
}
