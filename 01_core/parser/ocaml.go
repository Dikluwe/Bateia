/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ocaml.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	ocamlNameRegex = regexp.MustCompile(`\(name\s+([^\s\)]+)\)`)
	ocamlDepRegex  = regexp.MustCompile(`\(depends\s+([^\)]+)\)`)
	ocamlItemRegex = regexp.MustCompile(`([a-zA-Z0-9_\-]+)`)
)

// ParseOCamlDune analisa o conteúdo de um dune-project ou arquivo *.opam (OCaml).
func ParseOCamlDune(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "ocaml",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := ocamlNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	if match := ocamlDepRegex.FindStringSubmatch(content); len(match) >= 2 {
		depsRaw := match[1]
		items := ocamlItemRegex.FindAllString(depsRaw, -1)
		seen := make(map[string]bool)
		for _, item := range items {
			if item != "dune" && !seen[item] {
				seen[item] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     item,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
