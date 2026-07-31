/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/csharp.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"regexp"
	"strings"
	"time"

	"bateia/01_core/domain"
)

var (
	csPkgRefRegex = regexp.MustCompile(`(?i)<PackageReference\s+Include=["']([^"']+)["'](?:\s+Version=["']([^"']+)["'])?`)
)

// ParseCsProj analisa o conteúdo de um arquivo de projeto C# .csproj.
func ParseCsProj(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "csharp",
		Dependencies: make([]domain.Dependency, 0),
		ConsultadoEm: time.Now().UTC(),
		Content:      []byte(content),
	}

	matches := csPkgRefRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		pkgName := strings.TrimSpace(match[1])
		version := ""
		if len(match) >= 3 {
			version = strings.TrimSpace(match[2])
		}

		key := strings.ToLower(pkgName)
		if seen[key] {
			continue
		}
		seen[key] = true

		depType := "direct"
		if isTestPackage(key) {
			depType = "dev"
		}

		manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
			Name:     pkgName,
			Version:  version,
			Indirect: false,
			Type:     depType,
		})
	}

	return manifest, nil
}

func isTestPackage(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "xunit") ||
		strings.Contains(lower, "nunit") ||
		strings.Contains(lower, "mstest") ||
		strings.Contains(lower, "test.sdk") ||
		strings.Contains(lower, "moq") ||
		strings.Contains(lower, "nsubstitute") ||
		strings.Contains(lower, "fluentassertions")
}
