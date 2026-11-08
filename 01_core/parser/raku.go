/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/raku.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"encoding/json"

	"bateia/01_core/domain"
)

type meta6JsonStruct struct {
	Name    string   `json:"name"`
	Depends []string `json:"depends"`
}

// ParseMeta6Json analisa o conteúdo de META6.json (Raku).
func ParseMeta6Json(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "raku",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	var parsed meta6JsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for _, depName := range parsed.Depends {
			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  "",
				Indirect: false,
				Type:     "direct",
			})
		}
	} else {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}
