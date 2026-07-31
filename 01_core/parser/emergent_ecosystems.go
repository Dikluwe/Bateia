/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/emergent_ecosystems.md
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
	nimbleRequireRegex = regexp.MustCompile(`(?i)requires\s+["']([^"'\s>=<]+)`)
	shardNameRegex     = regexp.MustCompile(`(?m)^\s*name:\s*([^\s]+)`)
	shardDepRegex      = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-]+):`)
	vmodNameRegex      = regexp.MustCompile(`(?m)^\s*name:\s*['"]([^'"]+)['"]`)
	vmodDepRegex       = regexp.MustCompile(`['"]([^'"]+)['"]`)
	rebarDepRegex      = regexp.MustCompile(`\{\s*([a-zA-Z0-9_\-]+)\s*,`)
)

// ParseNimble analisa o conteúdo de um arquivo *.nimble (Nim).
func ParseNimble(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "nim",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range nimbleRequireRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			depName := match[1]
			if strings.HasPrefix(depName, "nim") {
				continue
			}

			if !seen[depName] {
				seen[depName] = true
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

// ParseShardYaml analisa o conteúdo de um arquivo shard.yml (Crystal).
func ParseShardYaml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "crystal",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := shardNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
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

		if strings.HasPrefix(trimmed, "dependencies:") {
			inDeps = true
			continue
		} else if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && strings.HasSuffix(trimmed, ":") {
			inDeps = false
		}

		if inDeps {
			if match := shardDepRegex.FindStringSubmatch(trimmed); len(match) >= 2 {
				depName := match[1]
				if depName != "github" && depName != "git" && depName != "version" && depName != "branch" {
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     depName,
						Version:  "",
						Indirect: false,
						Type:     "direct",
					})
				}
			}
		}
	}

	return manifest, nil
}

// ParseVMod analisa o conteúdo de um arquivo v.mod (V / Vlang).
func ParseVMod(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "v",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := vmodNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "deps:") {
			inDeps = true
		}

		if inDeps {
			for _, match := range vmodDepRegex.FindAllStringSubmatch(trimmed, -1) {
				if len(match) >= 2 {
					depName := match[1]
					manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
						Name:     depName,
						Version:  "",
						Indirect: false,
						Type:     "direct",
					})
				}
			}

			if strings.Contains(trimmed, "]") {
				inDeps = false
			}
		}
	}

	return manifest, nil
}

// ParseRebarConfig analisa o conteúdo de um arquivo rebar.config (Erlang).
func ParseRebarConfig(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "erlang",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)
	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{deps,") {
			inDeps = true
		}

		if inDeps {
			for _, match := range rebarDepRegex.FindAllStringSubmatch(trimmed, -1) {
				if len(match) >= 2 {
					depName := match[1]
					if depName != "deps" && depName != "erl_opts" && !seen[depName] {
						seen[depName] = true
						manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
							Name:     depName,
							Version:  "",
							Indirect: false,
							Type:     "direct",
						})
					}
				}
			}

			if strings.Contains(trimmed, "]}") || strings.HasSuffix(trimmed, "].") {
				inDeps = false
			}
		}
	}

	return manifest, nil
}
