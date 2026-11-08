/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/d.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"encoding/json"

	"bateia/01_core/domain"
)

type dubJsonStruct struct {
	Name         string                 `json:"name"`
	Dependencies map[string]interface{} `json:"dependencies"`
}

// ParseDubJson analisa o conteúdo de um arquivo dub.json (D language).
func ParseDubJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "d",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	var parsed dubJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for _, depName := range sortedManifestKeys(parsed.Dependencies) {
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
