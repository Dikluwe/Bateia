/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/haxe.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"encoding/json"

	"bateia/01_core/domain"
)

type haxelibJsonStruct struct {
	Name         string            `json:"name"`
	Dependencies map[string]string `json:"dependencies"`
}

// ParseHaxelibJson analisa o conteúdo de haxelib.json (Haxe).
func ParseHaxelibJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "haxe",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	var parsed haxelibJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for _, depName := range sortedManifestKeys(parsed.Dependencies) {
			ver := parsed.Dependencies[depName]
			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     depName,
				Version:  ver,
				Indirect: false,
				Type:     "direct",
			})
		}
	} else {
		manifest.ModuleName = ref.Repo
	}

	return manifest, nil
}
