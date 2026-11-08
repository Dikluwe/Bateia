/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/crystal.md
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
	shardNameRegex = regexp.MustCompile(`(?m)^\s*name:\s*([^\r\n]+)`)
	shardDepRegex  = regexp.MustCompile(`^([a-zA-Z0-9_\-]+):`)
)

// ParseShardYaml analisa o conteúdo de um arquivo shard.yml do ecossistema Crystal.
func ParseShardYaml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "crystal",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := shardNameRegex.FindStringSubmatch(content); len(match) >= 2 {
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
		} else if strings.HasPrefix(trimmed, "development_dependencies:") || (strings.HasSuffix(trimmed, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")) {
			inDeps = false
		}

		if inDeps && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
			if match := shardDepRegex.FindStringSubmatch(trimmed); len(match) >= 2 {
				depName := match[1]
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
