/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/emergent_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseNimble_Basic(t *testing.T) {
	content := `
version       = "0.1.0"
author        = "Nim Developer"
description   = "Sample Nimble App"
license       = "MIT"

requires "nim >= 1.6.0"
requires "jester"
requires "norm"
`

	ref := domain.CommitRef{Owner: "nim-lang", Repo: "sample_app"}
	manifest, err := ParseNimble(ref, content)
	if err != nil {
		t.Fatalf("ParseNimble falhou: %v", err)
	}

	if manifest.Language != "nim" {
		t.Errorf("Esperado Language 'nim', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (jester, norm), obtido %d", len(manifest.Dependencies))
	}
}

func TestParseShardYaml_Basic(t *testing.T) {
	content := `
name: kemal_app
version: 0.1.0

dependencies:
  kemal:
    github: kemalcr/kemal
  pg:
    github: will/crystal-pg
`

	ref := domain.CommitRef{Owner: "kemalcr", Repo: "kemal_app"}
	manifest, err := ParseShardYaml(ref, content)
	if err != nil {
		t.Fatalf("ParseShardYaml falhou: %v", err)
	}

	if manifest.ModuleName != "kemal_app" {
		t.Errorf("Esperado ModuleName 'kemal_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "crystal" {
		t.Errorf("Esperado Language 'crystal', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (kemal, pg), obtido %d", len(manifest.Dependencies))
	}
}

func TestParseVMod_Basic(t *testing.T) {
	content := `Module {
	name: 'vweb_app'
	description: 'Web app written in V'
	version: '0.1.0'
	license: 'MIT'
	deps: ['vweb', 'sqlite']
}`

	ref := domain.CommitRef{Owner: "vlang", Repo: "vweb_app"}
	manifest, err := ParseVMod(ref, content)
	if err != nil {
		t.Fatalf("ParseVMod falhou: %v", err)
	}

	if manifest.Language != "v" {
		t.Errorf("Esperado Language 'v', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vweb, sqlite), obtido %d", len(manifest.Dependencies))
	}
}

func TestParseRebarConfig_Basic(t *testing.T) {
	content := `{erl_opts, [debug_info]}.
{deps, [
    {cowboy, "2.9.0"},
    {jsx, "~> 3.0"}
]}.`

	ref := domain.CommitRef{Owner: "erlang", Repo: "rebar_app"}
	manifest, err := ParseRebarConfig(ref, content)
	if err != nil {
		t.Fatalf("ParseRebarConfig falhou: %v", err)
	}

	if manifest.Language != "erlang" {
		t.Errorf("Esperado Language 'erlang', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (cowboy, jsx), obtido %d", len(manifest.Dependencies))
	}
}
