/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/scripting.md
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
	juliaNameRegex       = regexp.MustCompile(`(?m)^\s*name\s*=\s*["']([^"']+)["']`)
	juliaDepLineRegex    = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-]+)\s*=\s*["']([^"']+)["']`)
	rPackageRegex        = regexp.MustCompile(`(?i)Package:\s*([a-zA-Z0-9_\-\.]+)`)
	luaPackageRegex      = regexp.MustCompile(`(?i)package\s*=\s*["']([^"']+)["']`)
	luaDepRegex          = regexp.MustCompile(`["']([a-zA-Z0-9_\-\./]+)\s*(?:[~=<>]*)?["']`)
	duneNameRegex        = regexp.MustCompile(`\(name\s+([a-zA-Z0-9_\-]+)\)`)
	duneDependsRegex     = regexp.MustCompile(`\(depends\s+([^\)]+)\)`)
)

// ParseJuliaProjectToml analisa o conteúdo de um Project.toml de Julia.
func ParseJuliaProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "julia",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			currentSection = strings.ToLower(trimmed)
			continue
		}

		if manifest.ModuleName == "" && (currentSection == "" || currentSection == "[deps]") {
			if matches := juliaNameRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				manifest.ModuleName = matches[1]
			}
		}

		if currentSection == "[deps]" {
			if matches := juliaDepLineRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     matches[1],
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}

// ParseRDescription analisa o conteúdo de um arquivo DESCRIPTION de R.
func ParseRDescription(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "r",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := rPackageRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "Imports:") || strings.HasPrefix(trimmed, "Depends:") {
			currentSection = "direct"
			trimmed = strings.TrimPrefix(strings.TrimPrefix(trimmed, "Imports:"), "Depends:")
		} else if strings.HasPrefix(trimmed, "Suggests:") {
			currentSection = "dev"
			trimmed = strings.TrimPrefix(trimmed, "Suggests:")
		} else if strings.Contains(trimmed, ":") {
			currentSection = ""
		}

		if currentSection != "" {
			pkgs := strings.Split(trimmed, ",")
			for _, pkg := range pkgs {
				cleanPkg := strings.TrimSpace(pkg)
				if idx := strings.Index(cleanPkg, "("); idx != -1 {
					cleanPkg = strings.TrimSpace(cleanPkg[:idx])
				}
				if cleanPkg != "" && cleanPkg != "R" {
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     cleanPkg,
						Version:  "",
						Indirect: false,
						Type:     currentSection,
					})
				}
			}
		}
	}

	return manifest, nil
}

// ParseLuaRockspec analisa o conteúdo de um manifesto .rockspec de Lua.
func ParseLuaRockspec(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "lua",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := luaPackageRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := luaDepRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) >= 2 {
			pkg := strings.TrimSpace(match[1])
			if pkg != "" && !strings.HasPrefix(pkg, "lua") && !seen[pkg] {
				seen[pkg] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     pkg,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseOCamlDune analisa o conteúdo de um manifesto dune-project de OCaml.
func ParseOCamlDune(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "ocaml",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := duneNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	if matches := duneDependsRegex.FindStringSubmatch(content); len(matches) >= 2 {
		depsRaw := matches[1]
		tokens := strings.Fields(depsRaw)
		for _, tok := range tokens {
			clean := strings.Trim(tok, "()")
			if clean != "" && !strings.HasPrefix(clean, ":") {
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     clean,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
