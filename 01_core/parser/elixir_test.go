/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/elixir.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMixExs_Basic(t *testing.T) {
	content := `defmodule PhoenixApp.MixProject do
  use Mix.Project

  def project do
    [
      app: :phoenix_app,
      version: "0.1.0",
      deps: deps()
    ]
  end

  defp deps do
    [
      {:phoenix, "~> 1.7.0"},
      {:ecto_sql, "~> 3.10"}
    ]
  end
end`

	ref := domain.CommitRef{Owner: "phoenixframework", Repo: "phoenix_app"}
	manifest, err := ParseMixExs(ref, content)
	if err != nil {
		t.Fatalf("ParseMixExs falhou: %v", err)
	}

	if manifest.ModuleName != "phoenix_app" {
		t.Errorf("Esperado ModuleName 'phoenix_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "elixir" {
		t.Errorf("Esperado Language 'elixir', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (phoenix, ecto_sql), obtido %d", len(manifest.Dependencies))
	}
}
