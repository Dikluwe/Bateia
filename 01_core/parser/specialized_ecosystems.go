/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/specialized_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
)

var (
	podfileRegex   = regexp.MustCompile(`(?i)pod\s+['"]([^'"]+)['"]`)
	racketDepRegex = regexp.MustCompile(`["']([a-zA-Z0-9_\-]+)["']`)
	odinColRegex   = regexp.MustCompile(`["']([a-zA-Z0-9_\-]+)["']\s*:\s*["']([^"']+)["']`)
)

type haxelibJsonStruct struct {
	Name         string            `json:"name"`
	Dependencies map[string]string `json:"dependencies"`
}

type meta6JsonStruct struct {
	Name    string   `json:"name"`
	Depends []string `json:"depends"`
}

type corralJsonStruct struct {
	Deps []struct {
		Locator string `json:"locator"`
	} `json:"deps"`
}

// ParseAlireToml analisa o conteúdo de um arquivo alire.toml (Ada).
func ParseAlireToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "ada",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	manifest.ModuleName = ref.Repo
	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name =") {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				manifest.ModuleName = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}

		if strings.HasPrefix(trimmed, "[dependencies]") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "[") {
			inDeps = false
		}

		if inDeps {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				depName := strings.TrimSpace(parts[0])
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseMojoProjectToml analisa o conteúdo de mojoproject.toml ou pixi.toml (Mojo).
func ParseMojoProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "mojo",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	manifest.ModuleName = ref.Repo
	lines := strings.Split(content, "\n")
	inDeps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name =") {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				manifest.ModuleName = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}

		if strings.HasPrefix(trimmed, "[dependencies]") {
			inDeps = true
			continue
		} else if strings.HasPrefix(trimmed, "[") {
			inDeps = false
		}

		if inDeps {
			parts := strings.Split(trimmed, "=")
			if len(parts) >= 2 {
				depName := strings.TrimSpace(parts[0])
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParsePodfile analisa o conteúdo de um Podfile (Objective-C).
func ParsePodfile(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "objective-c",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range podfileRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			podName := match[1]
			if !seen[podName] {
				seen[podName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     podName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseHaxelibJson analisa o conteúdo de haxelib.json (Haxe).
func ParseHaxelibJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "haxe",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	var parsed haxelibJsonStruct
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Name != "" {
			manifest.ModuleName = parsed.Name
		} else {
			manifest.ModuleName = ref.Repo
		}

		for depName, ver := range parsed.Dependencies {
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

// ParseMeta6Json analisa o conteúdo de META6.json (Raku).
func ParseMeta6Json(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "raku",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
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

// ParseInfoRkt analisa o conteúdo de info.rkt (Racket).
func ParseInfoRkt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "racket",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range racketDepRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			depName := match[1]
			if depName != "base" && !seen[depName] {
				seen[depName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     depName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseMasonToml analisa o conteúdo de mason.toml (Chapel).
func ParseMasonToml(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	return ParseAlireToml(ref, content)
}

// ParseCorralJson analisa o conteúdo de corral.json (Pony).
func ParseCorralJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "pony",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
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

// ParseOlsJson analisa o conteúdo de ols.json ou odin.json (Odin).
func ParseOlsJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "odin",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range odinColRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 3 {
			colName := match[1]
			if !seen[colName] {
				seen[colName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     colName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
