/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cobol.md
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
	cobolCopyRegex = regexp.MustCompile(`(?i)COPY\s+([a-zA-Z0-9_\-]+)`)
	cobolCallRegex = regexp.MustCompile(`(?i)CALL\s+["']([a-zA-Z0-9_\-]+)["']`)
)

// ParseCobolSource analisa o conteúdo de um arquivo de código fonte COBOL (*.cbl / *.cob).
func ParseCobolSource(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "cobol",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range cobolCopyRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			copybook := strings.ToUpper(match[1])
			if !seen[copybook] {
				seen[copybook] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     copybook,
					Version:  "",
					Indirect: false,
					Type:     "copybook",
				})
			}
		}
	}

	for _, match := range cobolCallRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			callTarget := strings.ToUpper(match[1])
			if !seen[callTarget] {
				seen[callTarget] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     callTarget,
					Version:  "",
					Indirect: false,
					Type:     "call",
				})
			}
		}
	}

	return manifest, nil
}
