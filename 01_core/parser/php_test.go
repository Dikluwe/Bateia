/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/php_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseComposerJson_Basic(t *testing.T) {
	content := `{
		"name": "laravel/laravel",
		"type": "project",
		"require": {
			"php": "^8.2",
			"laravel/framework": "^10.10",
			"guzzlehttp/guzzle": "^7.2"
		},
		"require-dev": {
			"phpunit/phpunit": "^10.1",
			"pestphp/pest": "^2.0"
		}
	}`

	ref := domain.CommitRef{Owner: "laravel", Repo: "laravel"}
	manifest, err := ParseComposerJson(ref, content)
	if err != nil {
		t.Fatalf("ParseComposerJson falhou: %v", err)
	}

	if manifest.ModuleName != "laravel/laravel" {
		t.Errorf("Esperado ModuleName 'laravel/laravel', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "php" {
		t.Errorf("Esperado Language 'php', obtido '%s'", manifest.Language)
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
		t.Errorf("Esperado 2 dependências direct (laravel/framework, guzzlehttp/guzzle), obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (phpunit/phpunit, pestphp/pest), obtido %d", devCount)
	}
}
