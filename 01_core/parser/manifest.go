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

	if strings.HasSuffix(base, ".csproj") {
		return "csharp"
	}

	switch base {
	case "go.mod":
		return "gomod"
	case "package.json":
		return "packagejson"
	case "pyproject.toml":
		return "pyproject"
	case "requirements.txt":
		return "requirements"
	case "cargo.toml":
		return "cargo"
	case "pom.xml":
		return "pomxml"
	case "build.gradle", "build.gradle.kts":
		return "gradle"
	case "composer.json":
		return "composer"
	case "cmakelists.txt":
		return "cmake"
	case "vcpkg.json":
		return "vcpkg"
	case "conanfile.txt":
		return "conan"
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
	case "cargo":
		return ParseCargoToml(ref, content)
	case "pomxml":
		return ParsePomXml(ref, content)
	case "gradle":
		return ParseBuildGradle(ref, content)
	case "composer":
		return ParseComposerJson(ref, content)
	case "cmake":
		return ParseCMakeListsTxt(ref, content)
	case "vcpkg":
		return ParseVcpkgJson(ref, content)
	case "conan":
		return ParseConanfileTxt(ref, content)
	case "csharp":
		return ParseCsProj(ref, content)
	default:
		return nil, fmt.Errorf("unsupported manifest type: %s (file: %s)", manifestType, filename)
	}
}
