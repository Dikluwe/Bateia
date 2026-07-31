/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/classic_ecosystems.md
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
	delphiRefRegex     = regexp.MustCompile(`(?i)<DCCReference\s+Include=["']([^"']+)["']`)
	lazarusPkgRegex    = regexp.MustCompile(`(?i)<ItemName\s+Value=["']([^"']+)["']`)
	fortranNameRegex   = regexp.MustCompile(`(?m)^\s*name\s*=\s*["']([^"']+)["']`)
	fortranDepRegex    = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-]+)\s*=\s*`)
	basicRefRegex      = regexp.MustCompile(`(?i)(?:Reference|Object)\s*=\s*([^\r\n]+)`)
)

// ParsePascalDelphi analisa o conteúdo de um projeto Pascal/Delphi (.dproj / .lpi).
func ParsePascalDelphi(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "pascal",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	// Inspecionar referências DCC (Delphi)
	for _, match := range delphiRefRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			refName := match[1]
			if !seen[refName] {
				seen[refName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     refName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	// Inspecionar pacotes Lazarus (.lpi)
	for _, match := range lazarusPkgRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			pkgName := match[1]
			if !seen[pkgName] {
				seen[pkgName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     pkgName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseFortranFpm analisa o conteúdo de um fpm.toml (Fortran Package Manager).
func ParseFortranFpm(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "fortran",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := fortranNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "[dependencies]") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "[") {
			inDeps = false
		}

		if inDeps {
			if matches := fortranDepRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				depName := matches[1]
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

// ParseBasicVbp analisa o conteúdo de um arquivo de projeto .vbp (Visual Basic).
func ParseBasicVbp(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "basic",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	for _, match := range basicRefRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			refString := match[1]
			parts := strings.Split(refString, "#")
			depName := parts[len(parts)-1]
			depName = strings.TrimSpace(depName)

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  "",
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}
