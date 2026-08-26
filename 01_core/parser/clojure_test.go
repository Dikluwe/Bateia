/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/clojure_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseProjectClj_Basic(t *testing.T) {
	content := `(defproject ring-app "0.1.0-SNAPSHOT"
  :description "Ring web application"
  :dependencies [[org.clojure/clojure "1.11.1"]
                 [ring/ring-core "1.9.5"]])`

	ref := domain.CommitRef{Owner: "ring-clojure", Repo: "ring-app"}
	manifest, err := ParseProjectClj(ref, content)
	if err != nil {
		t.Fatalf("ParseProjectClj falhou: %v", err)
	}

	if manifest.ModuleName != "ring-app" {
		t.Errorf("Esperado ModuleName 'ring-app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "clojure" {
		t.Errorf("Esperado Language 'clojure', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (org.clojure/clojure, ring/ring-core), obtido %d", len(manifest.Dependencies))
	}
}
