/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ruby.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseGemfile_Basic(t *testing.T) {
	content := `source 'https://rubygems.org'

git_source(:github) { |repo_name| "https://github.com/#{repo_name}" }

gem 'rails', '~> 7.0.0'
gem 'pg', '~> 1.1'
gem 'puma', '~> 5.0'
`

	ref := domain.CommitRef{Owner: "rails", Repo: "sample_rails_app"}
	manifest, err := ParseGemfile(ref, content)
	if err != nil {
		t.Fatalf("ParseGemfile falhou: %v", err)
	}

	if manifest.Language != "ruby" {
		t.Errorf("Esperado Language 'ruby', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 3 {
		t.Fatalf("Esperado 3 gems (rails, pg, puma), obtido %d", len(manifest.Dependencies))
	}
}
