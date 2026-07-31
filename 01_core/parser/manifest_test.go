/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/manifest.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestDetectManifestType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"go.mod", "gomod"},
		{"/path/to/go.mod", "gomod"},
		{"package.json", "packagejson"},
		{"/root/package.json", "packagejson"},
		{"pyproject.toml", "pyproject"},
		{"requirements.txt", "requirements"},
		{"Cargo.toml", "cargo"},
		{"pom.xml", "pomxml"},
		{"build.gradle", "gradle"},
		{"build.gradle.kts", "gradle"},
		{"composer.json", "composer"},
		{"CMakeLists.txt", "cmake"},
		{"vcpkg.json", "vcpkg"},
		{"conanfile.txt", "conan"},
		{"MyProject.csproj", "csharp"},
		{"Project.toml", "julia"},
		{"DESCRIPTION", "rdescription"},
		{"mylib-1.0.rockspec", "rockspec"},
		{"dune-project", "ocaml"},
		{"mix.exs", "mix"},
		{"Gemfile", "gemfile"},
		{"Package.swift", "swiftpackage"},
		{"myproject.cabal", "haskell"},
		{"pubspec.yaml", "pubspec"},
		{"build.sbt", "sbt"},
		{"project.clj", "clojure"},
		{"build.zig.zon", "zig"},
		{"MyProject.dproj", "pascal"},
		{"fpm.toml", "fortran"},
		{"Project1.vbp", "basic"},
		{"App.fsproj", "fsharp"},
		{"cpanfile", "perl"},
		{"MAIN.cbl", "cobol"},
		{"unknown.txt", ""},
	}

	for _, tt := range tests {
		got := DetectManifestType(tt.filename)
		if got != tt.expected {
			t.Errorf("DetectManifestType(%s) = %s; esperado %s", tt.filename, got, tt.expected)
		}
	}
}

func TestParseManifest_GoMod(t *testing.T) {
	content := "module github.com/foo/bar\ngo 1.22"
	ref := domain.CommitRef{Owner: "foo", Repo: "bar"}

	manifest, err := ParseManifest("go.mod", ref, content)
	if err != nil {
		t.Fatalf("ParseManifest go.mod falhou: %v", err)
	}

	if manifest.Language != "go" {
		t.Errorf("Esperado Language 'go', obtido '%s'", manifest.Language)
	}
	if manifest.ModuleName != "github.com/foo/bar" {
		t.Errorf("Esperado ModuleName 'github.com/foo/bar', obtido '%s'", manifest.ModuleName)
	}
}

func TestParseManifest_PackageJson(t *testing.T) {
	content := `{"name": "my-node-app", "version": "1.0.0"}`
	ref := domain.CommitRef{Owner: "foo", Repo: "my-node-app"}

	manifest, err := ParseManifest("package.json", ref, content)
	if err != nil {
		t.Fatalf("ParseManifest package.json falhou: %v", err)
	}

	if manifest.Language != "typescript" {
		t.Errorf("Esperado Language 'typescript', obtido '%s'", manifest.Language)
	}
	if manifest.ModuleName != "my-node-app" {
		t.Errorf("Esperado ModuleName 'my-node-app', obtido '%s'", manifest.ModuleName)
	}
}

func TestParseManifest_Unsupported(t *testing.T) {
	ref := domain.CommitRef{Owner: "foo", Repo: "bar"}
	_, err := ParseManifest("unsupported.file", ref, "content")
	if err == nil {
		t.Error("Deveria retornar erro para arquivo não suportado")
	}
}
