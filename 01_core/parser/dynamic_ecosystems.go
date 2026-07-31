/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/dynamic_ecosystems.md
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
	fsharpRefRegex   = regexp.MustCompile(`(?i)<PackageReference\s+Include=["']([^"']+)["'](?:\s+Version=["']([^"']+)["'])?`)
	perlRequireRegex = regexp.MustCompile(`(?i)requires\s+['"]([^'"]+)['"](?:\s*,\s*['"]([^'"]+)['"])?`)
	cobolCopyRegex    = regexp.MustCompile(`(?i)\b(?:COPY|INCLUDE)\s+([A-Z0-9_\-]+)`)
)

// ParseFsProj analisa o conteúdo de um projeto F# (.fsproj).
func ParseFsProj(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "fsharp",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	for _, match := range fsharpRefRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			depName := match[1]
			version := ""
			if len(match) >= 3 {
				version = match[2]
			}

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  version,
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

// ParseCpanfile analisa o conteúdo de um cpanfile ou Makefile.PL (Perl).
func ParseCpanfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "perl",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	for _, match := range perlRequireRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			depName := match[1]
			version := ""
			if len(match) >= 3 {
				version = match[2]
			}

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  version,
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

// ParseCobolSource analisa o conteúdo de fontes COBOL (.cbl / .cob).
func ParseCobolSource(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "cobol",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range cobolCopyRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			copyName := match[1]
			if !seen[copyName] {
				seen[copyName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     copyName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
