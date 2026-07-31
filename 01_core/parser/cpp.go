/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cpp.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
)

var (
	cmakeProjectRegex  = regexp.MustCompile(`(?i)project\s*\(\s*([a-zA-Z0-9_\-]+)`)
	cmakeFindPkgRegex  = regexp.MustCompile(`(?i)find_package\s*\(\s*([a-zA-Z0-9_\-]+)`)
	conanRequiresRegex = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-\.]+)/([^@\s]+)`)
)

type VcpkgJSON struct {
	Name         string   `json:"name"`
	Dependencies []interface{} `json:"dependencies"`
}

// ParseCMakeListsTxt analisa o conteúdo de um CMakeLists.txt.
func ParseCMakeListsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     "cpp",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if matches := cmakeProjectRegex.FindStringSubmatch(content); len(matches) >= 2 {
		manifest.ModuleName = matches[1]
	} else {
		manifest.ModuleName = ref.Repo
	}

	pkgMatches := cmakeFindPkgRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)

	for _, match := range pkgMatches {
		if len(match) >= 2 {
			pkgName := strings.ToLower(match[1])
			if !seen[pkgName] {
				seen[pkgName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     pkgName,
					Version:  "",
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}

// ParseVcpkgJson analisa o conteúdo de um vcpkg.json.
func ParseVcpkgJson(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	var vcpkg VcpkgJSON
	if err := json.Unmarshal([]byte(content), &vcpkg); err != nil {
		return nil, fmt.Errorf("failed to parse vcpkg.json: %w", err)
	}

	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   vcpkg.Name,
		Language:     "cpp",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	if manifest.ModuleName == "" {
		manifest.ModuleName = ref.Repo
	}

	for _, depRaw := range vcpkg.Dependencies {
		depName := ""
		switch v := depRaw.(type) {
		case string:
			depName = v
		case map[string]interface{}:
			if nameVal, ok := v["name"].(string); ok {
				depName = nameVal
			}
		}

		if depName != "" {
			manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
				Name:     strings.ToLower(depName),
				Version:  "",
				Indirect: false,
				Type:     "direct",
			})
		}
	}

	return manifest, nil
}

// ParseConanfileTxt analisa o conteúdo de um conanfile.txt.
func ParseConanfileTxt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "cpp",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	lines := strings.Split(content, "\n")
	inRequires := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inRequires = (strings.ToLower(trimmed) == "[requires]")
			continue
		}

		if inRequires {
			if matches := conanRequiresRegex.FindStringSubmatch(trimmed); len(matches) >= 3 {
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     strings.ToLower(matches[1]),
					Version:  matches[2],
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
