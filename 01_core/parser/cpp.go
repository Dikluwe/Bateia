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
	cmakeProjectBody   = regexp.MustCompile(`(?i)project\s*\(([^)]+)\)`)
	cmakeFindPkgRegex  = regexp.MustCompile(`(?i)find_package\s*\(\s*([a-zA-Z0-9_\-]+)`)
	conanRequiresRegex = regexp.MustCompile(`(?m)^\s*([a-zA-Z0-9_\-\.]+)/([^@\s]+)`)
	cmakeCXXRegex      = regexp.MustCompile(`\bCXX\b`)
	cmakeCLangRegex    = regexp.MustCompile(`\bC\b`)
	conanLibcxxRegex   = regexp.MustCompile(`(?m)^\s*compiler\.libcxx\s*=`)
)

type VcpkgJSON struct {
	Name         string        `json:"name"`
	Dependencies []interface{} `json:"dependencies"`
}

// detectCMakeLanguage inspeciona o corpo da declaração project() e retorna
// "c" se apenas C estiver declarado (sem CXX), ou "cpp" caso contrário.
// Quando o project() omite a lista de linguagens, o padrão é "cpp".
func detectCMakeLanguage(content string) string {
	m := cmakeProjectBody.FindStringSubmatch(content)
	if len(m) < 2 {
		return "cpp"
	}
	body := m[1]
	if cmakeCXXRegex.MatchString(body) {
		return "cpp"
	}
	if cmakeCLangRegex.MatchString(body) {
		return "c"
	}
	// project(name) sem linguagem declarada → padrão histórico C++
	return "cpp"
}

// detectConanLanguage retorna "c" se o conanfile.txt não declara
// compiler.libcxx (indicador de runtime C++), e "cpp" caso contrário.
func detectConanLanguage(content string) string {
	if conanLibcxxRegex.MatchString(content) {
		return "cpp"
	}
	return "c"
}

// ParseCMakeListsTxt analisa o conteúdo de um CMakeLists.txt.
// Detecta automaticamente se o projeto é C ou C++ pela declaração project().
func ParseCMakeListsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	lang := detectCMakeLanguage(content)

	manifest := &domain.Manifest{
		Ref:          ref,
		Language:     lang,
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
// Exclusivo de C++ — sempre produz Language: "cpp".
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
// Detecta automaticamente C vs C++ pela presença de compiler.libcxx nas settings.
func ParseConanfileTxt(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	lang := detectConanLanguage(content)

	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     lang,
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
