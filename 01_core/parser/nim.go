/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/nim.md
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
	nimPackageRegex = regexp.MustCompile(`(?m)^\s*packageName\s*=\s*["']([^"']+)["']`)
	nimRequiresRegex = regexp.MustCompile(`requires\s+["']([^"']+)["']`)
)

// ParseNimble analisa o conteúdo de um arquivo *.nimble do ecossistema Nim.
func ParseNimble(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "nim",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := nimPackageRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	seen := make(map[string]bool)

	for _, match := range nimRequiresRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			reqStr := match[1]
			parts := strings.Fields(reqStr)
			if len(parts) >= 1 {
				depName := parts[0]
				if depName != "nim" && !seen[depName] {
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

	return manifest, nil
}
