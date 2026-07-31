/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/r.md
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
	rPackageRegex = regexp.MustCompile(`(?m)^Package:\s*([^\r\n]+)`)
	rImportsRegex = regexp.MustCompile(`(?m)^Imports:\s*([^\r\n]+(?:\r?\n\s+[^\r\n]+)*)`)
)

// ParseRDescription analisa o conteúdo de um manifesto DESCRIPTION do ecossistema R.
func ParseRDescription(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "r",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if match := rPackageRegex.FindStringSubmatch(content); len(match) >= 2 {
		manifest.ModuleName = strings.TrimSpace(match[1])
	} else {
		manifest.ModuleName = ref.Repo
	}

	if match := rImportsRegex.FindStringSubmatch(content); len(match) >= 2 {
		rawImports := match[1]
		rawImports = strings.ReplaceAll(rawImports, "\n", " ")
		rawImports = strings.ReplaceAll(rawImports, "\r", " ")
		items := strings.Split(rawImports, ",")
		for _, item := range items {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				parts := strings.Fields(trimmed)
				depName := parts[0]
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
