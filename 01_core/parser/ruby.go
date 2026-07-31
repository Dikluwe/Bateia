/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ruby.md
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
	rubyGemRegex = regexp.MustCompile(`(?m)^\s*gem\s+["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`)
)

// ParseGemfile analisa o conteúdo de um manifesto Gemfile do ecossistema Ruby.
func ParseGemfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "ruby",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range rubyGemRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			gemName := match[1]
			ver := ""
			if len(match) >= 3 {
				ver = match[2]
			}

			if !seen[gemName] {
				seen[gemName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     gemName,
					Version:  ver,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
