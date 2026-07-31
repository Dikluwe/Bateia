/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/package_json.md
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

// PackageJSON representa a estrutura bruta do package.json para unmarshal.
type PackageJSON struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

// ParsePackageJson extrai dependências de um manifesto package.json.
func ParsePackageJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	var pkgJson PackageJSON

	if err := json.Unmarshal([]byte(content), &pkgJson); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   pkgJson.Name,
		Language:     "typescript",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	// Mapear dependências de produção (direct)
	for name, version := range pkgJson.Dependencies {
		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     name,
			Version:  version,
			Indirect: false,
			Type:     "direct",
		})
	}

	// Mapear dependências de desenvolvimento (dev)
	for name, version := range pkgJson.DevDependencies {
		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     name,
			Version:  version,
			Indirect: false,
			Type:     "dev",
		})
	}

	// Mapear dependências de runtime peer (peer)
	for name, version := range pkgJson.PeerDependencies {
		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     name,
			Version:  version,
			Indirect: false,
			Type:     "peer",
		})
	}

	return manifest, nil
}
