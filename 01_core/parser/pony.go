/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/pony.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"encoding/json"

	"bateia/01_core/domain"
)

type corralJsonStruct struct {
	Deps []struct {
		Locator string `json:"locator"`
	} `json:"deps"`
}

// ParseCorralJson analisa o conteúdo de corral.json (Pony).
func ParseCorralJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "pony",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	var parsed corralJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		for _, dep := range parsed.Deps {
			if dep.Locator != "" {
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     dep.Locator,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
