/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/coldfusion.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"encoding/json"

	"bateia/01_core/domain"
)

type boxJsonStruct struct {
	Name         string            `json:"name"`
	Dependencies map[string]string `json:"dependencies"`
}

// ParseBoxJson analisa o conteúdo de um arquivo box.json (ColdFusion CommandBox).
func ParseBoxJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "coldfusion",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	var parsed boxJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for _, depName := range sortedManifestKeys(parsed.Dependencies) {
			version := parsed.Dependencies[depName]
			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  version,
				Indirect: false,
				Type:     "direct",
			})
		}
	} else {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}
