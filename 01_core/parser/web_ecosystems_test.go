/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/web_ecosystems.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMixExs_Basic(t *testing.T) {
	content := `defmodule MyModule.MixProject do
  use Mix.Project

  def project do
    [
      app: :my_elixir_app,
      version: "0.1.0",
      deps: deps()
    ]
  end

  defp deps do
    [
      {:phoenix, "~> 1.7.10"},
      {:ecto_sql, "~> 3.10"}
    ]
  end
end`

	ref := domain.CommitRef{Owner: "elixir", Repo: "my_elixir_app"}
	manifest, err := ParseMixExs(ref, content)
	if err != nil {
		t.Fatalf("ParseMixExs falhou: %v", err)
	}

	if manifest.ModuleName != "my_elixir_app" {
		t.Errorf("Esperado ModuleName 'my_elixir_app', obtido '%s'", manifest.ModuleName)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}

func TestParseGemfile_Basic(t *testing.T) {
	content := `source 'https://rubygems.org'

gem 'rails', '~> 7.1.0'
gem 'pg', '~> 1.1'

group :development, :test do
  gem 'rspec-rails'
  gem 'rubocop'
end`

	ref := domain.CommitRef{Owner: "rails", Repo: "my_ruby_app"}
	manifest, err := ParseGemfile(ref, content)
	if err != nil {
		t.Fatalf("ParseGemfile falhou: %v", err)
	}

	if len(manifest.Dependencies) != 4 {
		t.Fatalf("Esperado 4 dependências, obtido %d", len(manifest.Dependencies))
	}
}
