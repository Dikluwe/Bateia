/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/swift.md
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
	swiftAppNameRegex = regexp.MustCompile(`name:\s*["']([^"']+)["']`)
	swiftDepRegex     = regexp.MustCompile(`\.package\s*\(\s*url:\s*["']([^"']+)["']`)
)

// ParsePackageSwift analisa o conteúdo de um manifesto Package.swift do ecossistema Swift.
func ParsePackageSwift(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "swift",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := swiftAppNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	seen := make(map[string]bool)

	for _, match := range swiftDepRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			url := match[1]
			parts := strings.Split(url, "/")
			pkgName := parts[len(parts)-1]
			pkgName = strings.TrimSuffix(pkgName, ".git")

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
