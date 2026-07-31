/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/clojure.md
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
	clojureNameRegex = regexp.MustCompile(`\(defproject\s+([^\s]+)`)
	clojureDepRegex  = regexp.MustCompile(`\[([a-zA-Z0-9_\-\.\/]+)\s+["']([^"']+)["']\]`)
)

// ParseProjectClj analisa o conteúdo de um project.clj ou deps.edn do ecossistema Clojure.
func ParseProjectClj(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "clojure",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := clojureNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	seen := make(map[string]bool)

	for _, match := range clojureDepRegex.FindAllStringSubmatch(content, -1) {
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
