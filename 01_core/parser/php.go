/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/php.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"encoding/json"
	"fmt"
	"time"

	"bateia/01_core/domain"
)

// ComposerJSON representa a estrutura bruta do composer.json para unmarshal.
type ComposerJSON struct {
	Name       string            `json:"name"`
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

// ParseComposerJson extrai dependências de um manifesto composer.json de PHP.
func ParseComposerJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	var compJson ComposerJSON

	if err := json.Unmarshal([]byte(content), &compJson); err != nil {
		return nil, fmt.Errorf("failed to parse composer.json: %w", err)
	}

	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   compJson.Name,
		Language:     "php",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	// Mapear dependências de produção (direct)
	for name, version := range compJson.Require {
		if name == "php" || len(name) > 4 && name[:4] == "ext-" {
			continue
		}
		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     name,
			Version:  version,
			Indirect: false,
			Type:     "direct",
		})
	}

	// Mapear dependências de dev (dev)
	for name, version := range compJson.RequireDev {
		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     name,
			Version:  version,
			Indirect: false,
			Type:     "dev",
		})
	}

	return manifest, nil
}
