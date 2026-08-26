/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/python_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePyProjectToml_Basic(t *testing.T) {
	content := `
[tool.poetry]
name = "my-fastapi-service"
version = "0.1.0"

[tool.poetry.dependencies]
python = "^3.11"
fastapi = "^0.109.0"
sqlalchemy = "^2.0.25"

[tool.poetry.group.dev.dependencies]
pytest = "^8.0.0"
ruff = "^0.2.0"
`

	ref := domain.CommitRef{Owner: "tiangolo", Repo: "my-fastapi-service"}
	manifest, err := ParsePyProjectToml(ref, content)
	if err != nil {
		t.Fatalf("ParsePyProjectToml falhou: %v", err)
	}

	if manifest.ModuleName != "my-fastapi-service" {
		t.Errorf("Esperado ModuleName 'my-fastapi-service', obtido '%s'", manifest.ModuleName)
	}
	if manifest.Language != "python" {
		t.Errorf("Esperado Language 'python', obtido '%s'", manifest.Language)
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
		t.Errorf("Esperado 2 dependências direct (fastapi, sqlalchemy), obtido %d", directCount)
	}
	if devCount != 2 {
		t.Errorf("Esperado 2 dependências dev (pytest, ruff), obtido %d", devCount)
	}
}

func TestParseRequirementsTxt_Basic(t *testing.T) {
	content := `
# Requisitos do projeto
fastapi>=0.100.0
sqlalchemy==2.0.0 # ORM principal
psycopg2-binary
-e .
`

	ref := domain.CommitRef{Owner: "foo", Repo: "python-app"}
	manifest, err := ParseRequirementsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseRequirementsTxt falhou: %v", err)
	}

	if manifest.Language != "python" {
		t.Errorf("Esperado Language 'python', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 3 {
		t.Fatalf("Esperado 3 dependências, obtido %d", len(manifest.Dependencies))
	}

	if manifest.Dependencies[0].Name != "fastapi" {
		t.Errorf("Esperado 'fastapi', obtido '%s'", manifest.Dependencies[0].Name)
	}
}
