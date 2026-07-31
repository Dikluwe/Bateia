/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/modern_ecosystems.md
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
	pubspecNameRegex     = regexp.MustCompile(`(?m)^\s*name:\s*([a-zA-Z0-9_]+)`)
	pubspecDepRegex      = regexp.MustCompile(`(?m)^\s{2}([a-zA-Z0-9_]+):\s*(?:\^|>=|~>|v)?([0-9\.\-]+|\{)?`)
	sbtNameRegex         = regexp.MustCompile(`name\s*:=\s*["']([^"']+)["']`)
	sbtDepRegex          = regexp.MustCompile(`["']([^"']+)["']\s*%%\s*["']([^"']+)["']\s*%\s*["']([^"']+)["']`)
	cljNameRegex         = regexp.MustCompile(`defproject\s+([^\s]+)`)
	cljDepRegex          = regexp.MustCompile(`\[([a-zA-Z0-9_\-\./]+)\s+["']([^"']+)["']\]`)
	zigNameRegex         = regexp.MustCompile(`\.name\s*=\s*["']([^"']+)["']`)
	zigDepRegex          = regexp.MustCompile(`\.([a-zA-Z0-9_\-]+)\s*=\s*\.\{`)
)

// ParsePubspecYaml analisa o conteúdo de um pubspec.yaml (Dart / Flutter).
func ParsePubspecYaml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "dart",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := pubspecNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		if strings.TrimSpace(trimmed) == "" || strings.HasPrefix(strings.TrimSpace(trimmed), "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "dependencies:") {
			currentSection = "direct"
			continue
		} else if strings.HasPrefix(trimmed, "dev_dependencies:") {
			currentSection = "dev"
			continue
		} else if !strings.HasPrefix(trimmed, " ") && !strings.HasPrefix(trimmed, "\t") {
			currentSection = ""
		}

		if currentSection != "" {
			if matches := pubspecDepRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				depName := matches[1]
				version := ""
				if len(matches) >= 3 && matches[2] != "{" {
					version = matches[2]
				}

				if depName != "flutter" && depName != "sdk" {
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     depName,
						Version:  version,
						Indirect: false,
						Type:     currentSection,
					})
				}
			}
		}
	}

	return manifest, nil
}

// ParseBuildSbt analisa o conteúdo de um build.sbt (Scala).
func ParseBuildSbt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "scala",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := sbtNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := sbtDepRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			groupID := match[1]
			artifactID := match[2]
			version := match[3]

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     groupID + "::" + artifactID,
				Version:  version,
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

// ParseProjectClj analisa o conteúdo de um project.clj ou deps.edn (Clojure).
func ParseProjectClj(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "clojure",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := cljNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := cljDepRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			depName := match[1]
			version := match[2]

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

// ParseBuildZigZon analisa o conteúdo de um build.zig.zon ou build.zig (Zig).
func ParseBuildZigZon(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "zig",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := zigNameRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	matches := zigDepRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			depName := match[1]
			if depName != "dependencies" && depName != "paths" {
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
