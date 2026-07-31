/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/modern_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePubspecYaml_Basic(t *testing.T) {
	content := `name: my_flutter_app
description: A new Flutter project.

dependencies:
  flutter:
    sdk: flutter
  provider: ^6.0.0
  dio: ^5.0.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^2.0.0
`

	ref := domain.CommitRef{Owner: "flutter", Repo: "my_flutter_app"}
	manifest, err := ParsePubspecYaml(ref, content)
	if err != nil {
		t.Fatalf("ParsePubspecYaml falhou: %v", err)
	}

	if manifest.ModuleName != "my_flutter_app" {
		t.Errorf("Esperado ModuleName 'my_flutter_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "dart" {
		t.Errorf("Esperado Language 'dart', obtido '%s'", manifest.Language)
	}

	directCount := 0
	devCount := 0
	for _, dep := range manifest.Dependencies {
		if dep.Type == "direct" {
			directCount++
		} else if dep.Type == "dev" {
			devCount++
		}
	}

	if directCount != 2 {
		t.Errorf("Esperado 2 dependências direct (provider, dio), obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (flutter_test, flutter_lints), obtido %d", devCount)
	}
}

func TestParseBuildSbt_Basic(t *testing.T) {
	content := `
name := "my-scala-service"
version := "0.1.0"

libraryDependencies += "org.http4s" %% "http4s-dsl" % "0.23.18"
libraryDependencies += "org.scalatest" %% "scalatest" % "3.2.15" % Test
`

	ref := domain.CommitRef{Owner: "typelevel", Repo: "my-scala-service"}
	manifest, err := ParseBuildSbt(ref, content)
	if err != nil {
		t.Fatalf("ParseBuildSbt falhou: %v", err)
	}

	if manifest.ModuleName != "my-scala-service" {
		t.Errorf("Esperado ModuleName 'my-scala-service', obtido '%s'", manifest.ModuleName)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}
