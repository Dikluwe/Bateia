/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/fsharp.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"regexp"

	"bateia/01_core/domain"
)

var (
	fsharpPackageReferenceRegex = regexp.MustCompile(`(?i)<PackageReference\s+Include=["']([^"']+)["'](?:\s+Version=["']([^"']+)["'])?`)
)

// ParseFsProj analisa o conteúdo de um arquivo .fsproj (F#).
func ParseFsProj(ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   ref.Repo,
		Language:     "fsharp",
		Dependencies: make([]domain.Dependency, 0),
		Content:      []byte(content),
	}

	seen := make(map[string]bool)

	for _, match := range fsharpPackageReferenceRegex.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			pkgName := match[1]
			ver := ""
			if len(match) >= 3 {
				ver = match[2]
			}

			if !seen[pkgName] {
				seen[pkgName] = true
				manifest.Dependencies = append(manifest.Dependencies, domain.Dependency{
					Name:     pkgName,
					Version:  ver,
					Indirect: false,
					Type:     "direct",
				})
			}
		}
	}

	return manifest, nil
}
