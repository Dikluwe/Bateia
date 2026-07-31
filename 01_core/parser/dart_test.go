/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/dart.md
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
version: 1.0.0+1

environment:
  sdk: ">=2.17.0 <3.0.0"

dependencies:
  flutter:
    sdk: flutter
  provider: ^6.0.0
  dio: ^5.0.0
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

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (provider, dio), obtido %d", len(manifest.Dependencies))
	}
}
