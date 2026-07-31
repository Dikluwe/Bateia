/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/web_ecosystems.md
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
	mixAppNameRegex    = regexp.MustCompile(`app:\s*:([a-zA-Z0-9_]+)`)
	mixDepTupleRegex   = regexp.MustCompile(`\{:([a-zA-Z0-9_]+),\s*["']([^"']+)["']`)
	gemNameRegex       = regexp.MustCompile(`(?m)^\s*gem\s+["']([^"']+)["'](?:\s*,\s*["']([^"']+)["'])?`)
	swiftPkgNameRegex  = regexp.MustCompile(`name:\s*["']([^"']+)["']`)
	swiftDepRegex      = regexp.MustCompile(`\.package\s*\(\s*url:\s*["']([^"']+)["']`)
	cabalNameRegex     = regexp.MustCompile(`(?i)^name:\s*([a-zA-Z0-9_\-]+)`)
	cabalBuildDepRegex = regexp.MustCompile(`build-depends:`)
)

// ParseMixExs analisa o conteúdo de um arquivo mix.exs de Elixir.
func ParseMixExs(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "elixir",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := mixAppNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := mixDepTupleRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) >= 3 {
			depName := match[1]
			version := match[2]
			if !seen[depName] {
				seen[depName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  version,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseGemfile analisa o conteúdo de um Gemfile de Ruby.
func ParseGemfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "ruby",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	inDevGroup := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "group ") && (strings.Contains(trimmed, ":development") || strings.Contains(trimmed, ":test")) {
			inDevGroup = true
		} else if trimmed == "end" {
			inDevGroup = false
		}

		if matches := gemNameRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
			gemName := matches[1]
			version := ""
			if len(matches) >= 3 {
				version = matches[2]
			}

			depType := "direct"
			if inDevGroup || strings.Contains(gemName, "rspec") || strings.Contains(gemName, "rubocop") || strings.Contains(gemName, "capybara") {
				depType = "dev"
			}

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     gemName,
				Version:  version,
				Indirect: false,
				Type:     depType,
			})
		}
	}

	return manifest, nil
}

// ParsePackageSwift analisa o conteúdo de um Package.swift de Swift.
func ParsePackageSwift(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "swift",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := swiftPkgNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := swiftDepRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			url := match[1]
			parts := strings.Split(url, "/")
			pkgName := parts[len(parts)-1]
			pkgName = strings.TrimSuffix(pkgName, ".git")

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     pkgName,
				Version:  "",
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

// ParseHaskellCabal analisa o conteúdo de um arquivo *.cabal de Haskell.
func ParseHaskellCabal(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "haskell",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	inBuildDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if manifest.ModuleName == "" {
			if matches := cabalNameRegex.FindStringSubmatch(line); len(matches) >= 2 {
				manifest.ModuleName = matches[1]
			}
		}

		if cabalBuildDepRegex.MatchString(trimmed) {
			inBuildDeps = true
			trimmed = strings.TrimPrefix(trimmed, "build-depends:")
		} else if inBuildDeps && strings.Contains(line, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inBuildDeps = false
		}

		if inBuildDeps {
			pkgs := strings.Split(trimmed, ",")
			for _, p := range pkgs {
				cleanPkg := strings.TrimSpace(p)
				if idx := strings.IndexAny(cleanPkg, "<>= "); idx != -1 {
					cleanPkg = strings.TrimSpace(cleanPkg[:idx])
				}
				if cleanPkg != "" && cleanPkg != "base" {
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     cleanPkg,
						Version:  "",
						Indirect: false,
						Type:     "direct",
					})
				}
			}
		}
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}
