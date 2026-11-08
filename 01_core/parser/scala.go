/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/scala.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	sbtNameRegex = regexp.MustCompile(`(?m)^name\s*:=\s*["']([^"']+)["']`)
	sbtDepRegex  = regexp.MustCompile(`["']([^"']+)["']\s*%%?\s*["']([^"']+)["']\s*%\s*["']([^"']+)["']`)
)

// ParseBuildSbt analisa o conteúdo de um manifesto build.sbt do ecossistema Scala.
func ParseBuildSbt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "scala",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	if match := sbtNameRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = match[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	seen := make(map[string]bool)

	for _, match := range sbtDepRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 4 {
			depName := match[2]
			ver := match[3]
			if !seen[depName] {
				seen[depName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  ver,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
