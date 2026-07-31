/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/coldfusion.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"encoding/json"
	"time"

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
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	var parsed boxJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for depName, version := range parsed.Dependencies {
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
