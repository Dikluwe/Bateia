/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/c.md
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
	makeTargetRegex = regexp.MustCompile(`(?m)^(?:TARGET|PROG)\s*[?:]*=\s*(\S+)`)
	makeLibsRegex   = regexp.MustCompile(`-l([a-zA-Z0-9_\-]+)`)
	mesonProjectRx  = regexp.MustCompile(`project\s*\(\s*'([^']+)'`)
	mesonDepRx      = regexp.MustCompile(`dependency\s*\(\s*'([^']+)'`)
)

// ParseMakefile analisa o conteúdo de um Makefile para projetos C puros.
func ParseMakefile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "c",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	// Extrair nome do projeto de TARGET ou PROG
	if m := makeTargetRegex.FindStringSubmatch(content); len(m) >= 2 {
		manifest.ModuleName = m[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	// Extrair dependências de -l<lib> em qualquer linha (LDFLAGS, LDLIBS, regras)
	matches := makeLibsRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 2 {
			lib := strings.ToLower(match[1])
			if !seen[lib] {
				seen[lib] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     lib,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseMesonBuild analisa o conteúdo de um meson.build para projetos C puros.
func ParseMesonBuild(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "c",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	// Extrair nome do projeto de project('name', ...)
	if m := mesonProjectRx.FindStringSubmatch(content); len(m) >= 2 {
		manifest.ModuleName = m[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	// Extrair dependências de dependency('lib')
	depMatches := mesonDepRx.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	for _, match := range depMatches {
		if len(match) >= 2 {
			dep := strings.ToLower(match[1])
			if !seen[dep] {
				seen[dep] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     dep,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
