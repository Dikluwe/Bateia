/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/python.md
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
	poetryNameRegex    = regexp.MustCompile(`(?m)^\s*name\s*=\s*["']([^"']+)["']`)
	depLineRegex       = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)\s*([<>=!~^].*)?$`)
	tomlPoetryDepRegex = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-\.]+)\s*=\s*(.+)$`)
)

// ParsePyProjectToml analisa o conteúdo de um arquivo pyproject.toml.
func ParsePyProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "python",
		Dependencies: make([]domain.Dependency, 0),
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

		// Extrair ModuleName de [tool.poetry] ou [project]
		if manifest.ModuleName == "" && (currentSection == "[tool.poetry]" || currentSection == "[project]") {
			if matches := poetryNameRegex.FindStringSubmatch(trimmed); len(matches) >= 2 {
				manifest.ModuleName = matches[1]
			}
		}

		// Extrair dependências de produção
		if currentSection == "[tool.poetry.dependencies]" || currentSection == "[project.dependencies]" {
			dep := parseTomlDepLine(trimmed, "direct")
			if dep != nil && dep.Name != "python" {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}

		// Extrair dependências de dev
		if strings.Contains(currentSection, "dev") || currentSection == "[tool.poetry.group.dev.dependencies]" {
			dep := parseTomlDepLine(trimmed, "dev")
			if dep != nil && dep.Name != "python" {
				manifest.Dependencies = append(manifest.Dependencies, *dep)
			}
		}
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}

// ParseRequirementsTxt analisa o conteúdo de um arquivo requirements.txt.
func ParseRequirementsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "python",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "-") {
			continue
		}

		// Remover comentários inline (ex: `fastapi # web framework`)
		if idx := strings.Index(trimmed, "#"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}

		matches := depLineRegex.FindStringSubmatch(trimmed)
		if len(matches) >= 2 {
			version := ""
			if len(matches) >= 3 {
				version = strings.TrimSpace(matches[2])
			}

			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     strings.ToLower(matches[1]),
				Version:  version,
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

func parseTomlDepLine(line, depType string) *domain.Dependency {
	matches := tomlPoetryDepRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}

	name := strings.ToLower(strings.TrimSpace(matches[1]))
	rawVersion := strings.Trim(strings.TrimSpace(matches[2]), `"'`)

	return &domain.Dependency{
		Name:     name,
		Version:  rawVersion,
		Indirect: false,
		Type:     depType,
	}
}
