/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/erlang.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseRebarConfig_Basic(t *testing.T) {
	content := `{erl_opts, [debug_info]}.
{deps, [
    {cowboy, "2.9.0"},
    {epgsql, "4.7.0"}
]}.`

	ref := domain.CommitRef{Owner: "erlang", Repo: "cowboy_app"}
	manifest, err := ParseRebarConfig(ref, content)
	if err != nil {
		t.Fatalf("ParseRebarConfig falhou: %v", err)
	}

	if manifest.Language != "erlang" {
		t.Errorf("Esperado Language 'erlang', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (cowboy, epgsql), obtido %d", len(manifest.Dependencies))
	}
}
