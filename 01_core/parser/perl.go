/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/perl.md
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
	perlRequiresRegex = regexp.MustCompile(`(?m)^\s*requires\s+["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`)
)

// ParseCpanfile analisa o conteúdo de um manifesto cpanfile ou Makefile.PL (Perl).
func ParseCpanfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "perl",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range perlRequiresRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			moduleName := match[1]
			ver := ""
			if len(match) >= 3 {
				ver = match[2]
			}

			if moduleName != "perl" && !seen[moduleName] {
				seen[moduleName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     moduleName,
					Version:  ver,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
