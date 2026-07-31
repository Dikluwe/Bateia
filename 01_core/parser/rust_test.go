/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/rust.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseCargoToml_Basic(t *testing.T) {
	content := `
[package]
name = "axum-web-api"
version = "0.1.0"
edition = "2021"

[dependencies]
tokio = { version = "1.35.0", features = ["full"] }
axum = "0.7.4"
serde = { version = "1.0.195", features = ["derive"] }
tracing = "0.1.40"

[dev-dependencies]
criterion = "0.5.1"
tokio-test = "0.4.3"
`

	ref := domain.CommitRef{Owner: "tokio-rs", Repo: "axum-web-api"}
	manifest, err := ParseCargoToml(ref, content)
	if err != nil {
		t.Fatalf("ParseCargoToml falhou: %v", err)
	}

	if manifest.ModuleName != "axum-web-api" {
		t.Errorf("Esperado ModuleName 'axum-web-api', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "rust" {
		t.Errorf("Esperado Language 'rust', obtido '%s'", manifest.Language)
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

	if directCount != 4 {
		t.Errorf("Esperado 4 dependências direct (tokio, axum, serde, tracing), obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (criterion, tokio-test), obtido %d", devCount)
	}
}
