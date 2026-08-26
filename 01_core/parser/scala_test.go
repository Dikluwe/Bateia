/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/scala_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseBuildSbt_Basic(t *testing.T) {
	content := `name := "HttpServerApp"
version := "0.1.0"
scalaVersion := "2.13.8"

libraryDependencies ++= Seq(
  "org.http4s" %% "http4s-dsl" % "0.23.12",
  "org.http4s" %% "http4s-ember-server" % "0.23.12"
)`

	ref := domain.CommitRef{Owner: "typelevel", Repo: "HttpServerApp"}
	manifest, err := ParseBuildSbt(ref, content)
	if err != nil {
		t.Fatalf("ParseBuildSbt falhou: %v", err)
	}

	if manifest.ModuleName != "HttpServerApp" {
		t.Errorf("Esperado ModuleName 'HttpServerApp', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "scala" {
		t.Errorf("Esperado Language 'scala', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (http4s-dsl, http4s-ember-server), obtido %d", len(manifest.Dependencies))
	}
}
