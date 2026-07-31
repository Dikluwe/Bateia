/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/package_json.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePackageJson_Basic(t *testing.T) {
	content := `{
		"name": "express-app",
		"version": "1.0.0",
		"dependencies": {
			"express": "^4.18.2",
			"prisma": "^5.10.0"
		},
		"devDependencies": {
			"jest": "^29.7.0",
			"typescript": "^5.3.3"
		}
	}`

	ref := domain.CommitRef{Owner: "expressjs", Repo: "express-app", Commit: "abc123"}
	manifest, err := ParsePackageJson(ref, content)
	if err != nil {
		t.Fatalf("ParsePackageJson falhou: %v", err)
	}

	if manifest.ModuleName != "express-app" {
		t.Errorf("Esperado ModuleName 'express-app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "typescript" {
		t.Errorf("Esperado Language 'typescript', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 4 {
		t.Errorf("Esperado 4 dependências totais, obtido %d", len(manifest.Dependencies))
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
		t.Errorf("Esperado 2 dependências direct, obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev, obtido %d", devCount)
	}
}

func TestParsePackageJson_PeerDeps(t *testing.T) {
	content := `{
		"name": "ui-components",
		"peerDependencies": {
			"react": ">=18.0.0"
		}
	}`

	ref := domain.CommitRef{Owner: "ui", Repo: "components"}
	manifest, err := ParsePackageJson(ref, content)
	if err != nil {
		t.Fatalf("ParsePackageJson falhou: %v", err)
	}

	if len(manifest.Dependencies) != 1 {
		t.Fatalf("Esperado 1 dependência, obtido %d", len(manifest.Dependencies))
	}

	if manifest.Dependencies[0].Type != "peer" {
		t.Errorf("Esperado tipo 'peer', obtido '%s'", manifest.Dependencies[0].Type)
	}
}

func TestParsePackageJson_InvalidJSON(t *testing.T) {
	content := `{ invalid json }`
	ref := domain.CommitRef{Owner: "test", Repo: "invalid"}

	_, err := ParsePackageJson(ref, content)
	if err == nil {
		t.Error("Deveria retornar erro para JSON inválido")
	}
}
