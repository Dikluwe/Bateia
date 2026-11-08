/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/manifest_time_test.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func manifestTimeCorpus() map[string]string {
	names := []string{"go.mod", "package.json", "pyproject.toml", "requirements.txt", "Cargo.toml", "pom.xml", "build.gradle", "composer.json", "CMakeLists.txt", "vcpkg.json", "conanfile.txt", "app.csproj", "Project.toml", "DESCRIPTION", "app.rockspec", "dune-project", "mix.exs", "Gemfile", "Package.swift", "app.cabal", "pubspec.yaml", "build.sbt", "project.clj", "build.zig.zon", "app.dproj", "fpm.toml", "app.vbp", "app.fsproj", "cpanfile", "app.cbl", "app.nimble", "shard.yml", "v.mod", "rebar.config", "dub.json", "gleam.toml", "box.json", "flex-config.xml", "alire.toml", "mojoproject.toml", "Podfile", "haxelib.json", "META6.json", "info.rkt", "mason.toml", "corral.json", "ols.json", "Makefile", "meson.build"}
	corpus := map[string]string{}
	for _, name := range names {
		content := ""
		if strings.HasSuffix(strings.ToLower(name), ".json") {
			content = "{}"
		}
		corpus[name] = content
	}
	corpus["go.mod"] = "module example.org/app\ngo 1.25\nrequire example.org/dep v1.2.3\n"
	corpus["package.json"] = `{"name":"app","dependencies":{"dep":"1.2.3"}}`
	corpus["Cargo.toml"] = "[package]\nname = \"app\"\n[dependencies]\nserde = \"1\"\n"
	corpus["requirements.txt"] = "requests==2.31.0\n"
	corpus["pom.xml"] = "<project><groupId>org.example</groupId><artifactId>app</artifactId></project>"
	corpus["app.csproj"] = "<Project><PropertyGroup><TargetFramework>net8.0</TargetFramework></PropertyGroup></Project>"
	return corpus
}

func TestManifestParsingCharacterization(t *testing.T) {
	projection := map[string]*domain.Manifest{}
	ref := domain.CommitRef{Owner: "example", Repo: "app", Commit: "abc123"}
	for name, content := range manifestTimeCorpus() {
		manifest, err := ParseManifest(name, ref, content)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		manifest.ConsultadoEm = time.Time{}
		projection[name] = manifest
	}
	encoded, err := json.Marshal(projection)
	if err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(encoded))
	const expected = "68f0ed92e2eab8424eabd5b1293756c074156855e8a5c398fdf3bcc9bd65ebcc"
	if digest != expected {
		t.Fatalf("manifest projection digest = %s, want %s", digest, expected)
	}
}

func TestManifestParsingHasExplicitCollectionTime(t *testing.T) {
	ref := domain.CommitRef{Owner: "example", Repo: "app", Commit: "abc123"}
	collectedAt := time.Date(2026, 10, 28, 9, 30, 12, 123456789, time.FixedZone("BRT", -3*60*60))
	for name, content := range manifestTimeCorpus() {
		t.Run(name, func(t *testing.T) {
			plain, err := ParseManifest(name, ref, content)
			if err != nil {
				t.Fatal(err)
			}
			if !plain.ConsultadoEm.IsZero() {
				t.Fatalf("parser invented collection time: %v", plain.ConsultadoEm)
			}
			for _, instant := range []time.Time{collectedAt, time.Time{}} {
				first, err := ParseManifestAt(name, ref, content, instant)
				if err != nil {
					t.Fatal(err)
				}
				second, err := ParseManifestAt(name, ref, content, instant)
				if err != nil {
					t.Fatal(err)
				}
				if !first.ConsultadoEm.Equal(instant) || first.ConsultadoEm.Location() != time.UTC {
					t.Fatalf("incorrect supplied instant: %v", first.ConsultadoEm)
				}
				one, _ := json.Marshal(first)
				two, _ := json.Marshal(second)
				if string(one) != string(two) {
					t.Fatal("same bytes/ref/time produced different manifests")
				}
				first.ConsultadoEm = time.Time{}
				normalized, _ := json.Marshal(first)
				original, _ := json.Marshal(plain)
				if string(normalized) != string(original) {
					t.Fatal("collection metadata changed parsing")
				}
			}
		})
	}
}

func TestManifestAtPreservesParsingErrors(t *testing.T) {
	for _, sample := range []struct{ name, content string }{{"unknown", ""}, {"package.json", "{"}, {"composer.json", "{"}, {"vcpkg.json", "{"}} {
		original, err := ParseManifest(sample.name, domain.CommitRef{}, sample.content)
		current, atErr := ParseManifestAt(sample.name, domain.CommitRef{}, sample.content, time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC))
		if original != nil || current != nil || err == nil || atErr == nil || err.Error() != atErr.Error() {
			t.Fatalf("error changed for %s: %v / %v", sample.name, err, atErr)
		}
	}
}

func TestManifestObjectDependenciesHaveStableOrder(t *testing.T) {
	direct := func(name, version string) domain.Dependency {
		return domain.Dependency{Name: name, Version: version, Type: "direct"}
	}
	cases := []struct {
		name, content string
		want          []domain.Dependency
	}{
		{"package.json", `{"dependencies":{"z":"2","a":"1"},"devDependencies":{"y":"4","b":"3"},"peerDependencies":{"x":"6","c":"5"}}`, []domain.Dependency{direct("a", "1"), direct("z", "2"), {Name: "b", Version: "3", Type: "dev"}, {Name: "y", Version: "4", Type: "dev"}, {Name: "c", Version: "5", Type: "peer"}, {Name: "x", Version: "6", Type: "peer"}}},
		{"composer.json", `{"require":{"z":"2","php":"8","ext-json":"*","a":"1"},"require-dev":{"y":"4","b":"3"}}`, []domain.Dependency{direct("a", "1"), direct("z", "2"), {Name: "b", Version: "3", Type: "dev"}, {Name: "y", Version: "4", Type: "dev"}}},
		{"box.json", `{"dependencies":{"z":"2","a":"1","m":"3"}}`, []domain.Dependency{direct("a", "1"), direct("m", "3"), direct("z", "2")}},
		{"haxelib.json", `{"dependencies":{"z":"2","a":"1","m":"3"}}`, []domain.Dependency{direct("a", "1"), direct("m", "3"), direct("z", "2")}},
		{"dub.json", `{"dependencies":{"z":"2","a":{"version":"1"},"m":"3"}}`, []domain.Dependency{direct("a", ""), direct("m", ""), direct("z", "")}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected, _ := json.Marshal(tc.want)
			for attempt := 0; attempt < 32; attempt++ {
				result, err := ParseManifest(tc.name, domain.CommitRef{}, tc.content)
				if err != nil {
					t.Fatal(err)
				}
				actual, _ := json.Marshal(result.Dependencies)
				if string(actual) != string(expected) {
					t.Fatalf("unstable dependency order: %s, want %s", actual, expected)
				}
			}
		})
	}
}
