/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/nim.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseNimble_Basic(t *testing.T) {
	content := `version       = "0.1.0"
author        = "Nim Developer"
description   = "Jester web application"
license       = "MIT"
packageName   = "jester_app"

requires "nim >= 1.6.0"
requires "jester >= 0.5.0"
requires "norm >= 2.2.0"
`

	ref := domain.CommitRef{Owner: "nim-lang", Repo: "jester_app"}
	manifest, err := ParseNimble(ref, content)
	if err != nil {
		t.Fatalf("ParseNimble falhou: %v", err)
	}

	if manifest.ModuleName != "jester_app" {
		t.Errorf("Esperado ModuleName 'jester_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "nim" {
		t.Errorf("Esperado Language 'nim', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (jester, norm), obtido %d", len(manifest.Dependencies))
	}
}
