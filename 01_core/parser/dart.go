/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/dart.md
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
	pubspecNameRegex = regexp.MustCompile(`(?m)^\s*name:\s*([^\r\n]+)`)
	pubspecDepRegex  = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-]+):\s*([^\r\n]*)`)
)

// ParsePubspecYaml analisa o conteúdo de um manifesto pubspec.yaml do ecossistema Dart/Flutter.
func ParsePubspecYaml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "dart",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := pubspecNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = strings.TrimSpace(match[1])
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "dependencies:") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "dev_dependencies:") || (strings.HasSuffix(trimmed, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")) {
			inDeps = false
		}

		if inDeps {
			if match := pubspecDepRegex.FindStringSubmatch(trimmed); len(match) >= 2 {
				depName := match[1]
				if depName != "flutter" && depName != "sdk" {
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     depName,
						Version:  strings.TrimSpace(match[2]),
						Indirect: false,
						Type:     "direct",
					})
				}
			}
		}
	}

	return manifest, nil
}
