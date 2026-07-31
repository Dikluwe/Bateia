/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/manifest.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"bateia/01_core/domain"
)

// DetectManifestType identifica o tipo de manifesto pelo nome do arquivo.
func DetectManifestType(filename string) string {
	base := filepath.Base(strings.ToLower(filename))

	switch base {
	case "go.mod":
		return "gomod"
	case "package.json":
		return "packagejson"
	case "pyproject.toml":
		return "pyproject"
	case "requirements.txt":
		return "requirements"
	default:
		return ""
	}
}

// ParseManifest roteia a execução para o parser correto baseado no arquivo.
func ParseManifest(filename string, ref domain.CommitRef, content string) (*domain.Manifest, error) {
	manifestType := DetectManifestType(filename)

	switch manifestType {
	case "gomod":
		return ParseGoMod(ref, content)
	case "packagejson":
		return ParsePackageJson(ref, content)
	case "pyproject":
		return ParsePyProjectToml(ref, content)
	case "requirements":
		return ParseRequirementsTxt(ref, content)
	default:
		return nil, fmt.Errorf("unsupported manifest type: %s (file: %s)", manifestType, filename)
	}
}
