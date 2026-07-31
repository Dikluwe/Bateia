/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/chapel.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"bateia/01_core/domain"
)

// ParseMasonToml analisa o conteúdo de mason.toml (Chapel).
func ParseMasonToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest, err := ParseAlireToml(ref, content)
	if err == nil && manifest != nil {
		manifest.Language = "chapel"
	}
	return manifest, err
}
