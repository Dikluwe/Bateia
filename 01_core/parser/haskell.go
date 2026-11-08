/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/haskell.md
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
	haskellNameRegex = regexp.MustCompile(`(?m)^name:\s*([^\r\n]+)`)
	haskellDepRegex  = regexp.MustCompile(`([a-zA-Z0-9_\-]+)\s*(?:>=|==|<=|>|<)?`)
)

// ParseHaskellCabal analisa o conteúdo de um arquivo *.cabal ou stack.yaml do ecossistema Haskell.
func ParseHaskellCabal(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "haskell",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := haskellNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = strings.TrimSpace(match[1])
	} else {
		manifest.ModuleName = ref.Repo
	}

	lines := strings.Split(content, "\n")
	inDeps := false
	seen := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "build-depends:") {
			inDeps = true
			trimmed = strings.TrimPrefix(trimmed, "build-depends:")
		} else if strings.HasPrefix(trimmed, "executable") || strings.HasPrefix(trimmed, "library") || strings.HasPrefix(trimmed, "test-suite") {
			inDeps = false
		}

		if inDeps && trimmed != "" {
			items := strings.Split(trimmed, ",")
			for _, item := range items {
				itemTrimmed := strings.TrimSpace(item)
				if itemTrimmed != "" {
					matches := haskellDepRegex.FindStringSubmatch(itemTrimmed)
					if len(matches) >= 2 {
						depName := matches[1]
						if depName != "base" && !seen[depName] {
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
			}
		}
	}

	return manifest, nil
}
