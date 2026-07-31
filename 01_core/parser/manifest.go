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
	if strings.HasSuffix(base, ".rockspec") {
		return "rockspec"
	}
	if strings.HasSuffix(base, ".opam") {
		return "ocaml"
	}
	if strings.HasSuffix(base, ".cabal") {
		return "haskell"
	}
	if strings.HasSuffix(base, ".dproj") || strings.HasSuffix(base, ".lpi") {
		return "pascal"
	}
	if strings.HasSuffix(base, ".vbp") {
		return "basic"
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
	case "project.toml":
		return "julia"
	case "description":
		return "rdescription"
	case "dune-project":
		return "ocaml"
	case "mix.exs":
		return "mix"
	case "gemfile":
		return "gemfile"
	case "package.swift":
		return "swiftpackage"
	case "stack.yaml":
		return "haskell"
	case "pubspec.yaml":
		return "pubspec"
	case "build.sbt":
		return "sbt"
	case "project.clj", "deps.edn":
		return "clojure"
	case "build.zig.zon", "build.zig":
		return "zig"
	case "fpm.toml":
		return "fortran"
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
	case "julia":
		return ParseJuliaProjectToml(ref, content)
	case "rdescription":
		return ParseRDescription(ref, content)
	case "rockspec":
		return ParseLuaRockspec(ref, content)
	case "ocaml":
		return ParseOCamlDune(ref, content)
	case "mix":
		return ParseMixExs(ref, content)
	case "gemfile":
		return ParseGemfile(ref, content)
	case "swiftpackage":
		return ParsePackageSwift(ref, content)
	case "haskell":
		return ParseHaskellCabal(ref, content)
	case "pubspec":
		return ParsePubspecYaml(ref, content)
	case "sbt":
		return ParseBuildSbt(ref, content)
	case "clojure":
		return ParseProjectClj(ref, content)
	case "zig":
		return ParseBuildZigZon(ref, content)
	case "pascal":
		return ParsePascalDelphi(ref, content)
	case "fortran":
		return ParseFortranFpm(ref, content)
	case "basic":
		return ParseBasicVbp(ref, content)
	default:
		return nil, fmt.Errorf("unsupported manifest type: %s (file: %s)", manifestType, filename)
	}
}
