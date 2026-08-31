/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-08-26
 */

package parser

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// CalculateHealthScore calcula a nota consolidada de saúde arquitetural (0 a 100%).
func CalculateHealthScore(violations []domain.ArchViolation) int {
	score := 100
	for _, v := range violations {
		switch v.Severity {
		case "CRITICAL", "HIGH":
			score -= 20
		case "MEDIUM":
			score -= 10
		case "LOW":
			score -= 5
		}
	}
	if score < 0 {
		return 0
	}
	return score
}

// EvaluateArchViolations analisa o manifesto e a estrutura de diretórios em busca de anti-padrões arquiteturais.
func EvaluateArchViolations(manifest *domain.Manifest, tree []domain.DirectoryNode) []domain.ArchViolation {
	var violations []domain.ArchViolation

	if manifest == nil {
		return violations
	}

	lang := strings.ToLower(manifest.Language)

	// 1. Verificar acoplamento excessivo de dependências
	if len(manifest.Dependencies) > 35 {
		snippet := "[dependencies] // > 35 dependências diretas"
		suggestion := "Divida as responsabilidades em múltiplos módulos/crates acoplados por interfaces genéricas."
		if lang == "rust" {
			snippet = "[dependencies] // > 35 crates no Cargo.toml"
			suggestion = "Divida o projeto em um Cargo Workspace com sub-crates especializados."
		} else if lang == "python" {
			snippet = "dependencies = [...] // > 35 pacotes no pyproject.toml"
			suggestion = "Desacople pacotes utilitários em módulos reutilizáveis."
		}

		violations = append(violations, domain.ArchViolation{
			Rule:        "HIGH_DEPENDENCY_COUPLING",
			Severity:    "HIGH",
			Message:     "Repositório possui acoplamento elevado com mais de 35 dependências diretas.",
			Target:      manifest.ModuleName,
			Line:        1,
			CodeSnippet: snippet,
			Suggestion:  suggestion,
		})
	}

	// 2. Verificar vazamento de camadas em manifestos
	for _, dep := range manifest.Dependencies {
		lowerDep := strings.ToLower(dep.Name)

		isORM := false
		switch lang {
		case "rust":
			isORM = strings.Contains(lowerDep, "diesel") || strings.Contains(lowerDep, "sqlx") || strings.Contains(lowerDep, "sea-orm") || strings.Contains(lowerDep, "rusqlite")
		case "python":
			isORM = strings.Contains(lowerDep, "sqlalchemy") || strings.Contains(lowerDep, "django") || strings.Contains(lowerDep, "tortoise") || strings.Contains(lowerDep, "peewee")
		case "typescript", "javascript":
			isORM = strings.Contains(lowerDep, "prisma") || strings.Contains(lowerDep, "typeorm") || strings.Contains(lowerDep, "sequelize") || strings.Contains(lowerDep, "drizzle")
		default: // go
			isORM = strings.Contains(lowerDep, "gorm") || strings.Contains(lowerDep, "ent") || strings.Contains(lowerDep, "prisma") || strings.Contains(lowerDep, "sqlx")
		}

		if isORM {
			hasInfra := false
			for _, node := range tree {
				lowerPath := strings.ToLower(node.Path)
				if strings.Contains(lowerPath, "infra") || strings.Contains(lowerPath, "repository") || strings.Contains(lowerPath, "database") || strings.Contains(lowerPath, "db") {
					hasInfra = true
					break
				}
			}
			if !hasInfra && len(tree) > 0 {
				suggestion := "Encapsule a dependência de banco dentro de uma camada dedicada (ex: 03_infra/ ou repository/)."
				if lang == "rust" {
					suggestion = "Encapsule o acesso ao banco/ORM dentro de um módulo/trait de infraestrutura (ex: src/infra/ ou src/db/)."
				} else if lang == "python" {
					suggestion = "Isole as chamadas do ORM dentro de uma camada de infraestrutura/repositório (ex: app/infra/ ou app/repositories/)."
				}

				violations = append(violations, domain.ArchViolation{
					Rule:        "MISSING_INFRASTRUCTURE_ISOLATION",
					Severity:    "MEDIUM",
					Message:     fmt.Sprintf("Dependência de ORM/Banco (%s) sem isolamento de camada de infraestrutura/repositório detectada.", dep.Name),
					Target:      dep.Name,
					Line:        10,
					CodeSnippet: dep.Name + " // Acesso direto ao ORM/Banco",
					Suggestion:  suggestion,
				})
			}
		}
	}

	if len(tree) > 0 {
		hasTests := false
		hasDevOps := false

		for _, node := range tree {
			lowerPath := strings.ToLower(node.Path)
			if strings.Contains(lowerPath, "test") || strings.Contains(lowerPath, "benches") || strings.Contains(lowerPath, "spec") ||
				MatchGlobRecursive("**/*test*", lowerPath) || MatchGlobRecursive("**/*spec*", lowerPath) || MatchGlobRecursive("**/tests/**", lowerPath) || MatchGlobRecursive("**/__tests__/**", lowerPath) {
				hasTests = true
			}
			if strings.Contains(lowerPath, "docker") || strings.Contains(lowerPath, ".github") || strings.Contains(lowerPath, "makefile") || strings.Contains(lowerPath, "ci.yml") {
				hasDevOps = true
			}
		}

		// 3. Ausência de Suíte de Testes
		if !hasTests {
			snippet := "// Ausência de pasta tests/ ou arquivos *_test.go"
			suggestion := "Adicione arquivos de testes (*_test.go) para garantir cobertura de regressão."
			if lang == "rust" {
				snippet = "// Ausência de pasta tests/ ou módulos #[cfg(test)]"
				suggestion = "Crie uma pasta tests/ ou adicione suítes de teste com cargo test."
			} else if lang == "python" {
				snippet = "// Ausência de pasta tests/ ou arquivos test_*.py"
				suggestion = "Adicione suítes de teste com pytest ou unittest na pasta tests/."
			}

			violations = append(violations, domain.ArchViolation{
				Rule:        "MISSING_TEST_SUITE",
				Severity:    "MEDIUM",
				Message:     "Nenhum diretório ou arquivo de teste unitário/integração detectado na estrutura do projeto.",
				Target:      manifest.ModuleName,
				Line:        1,
				CodeSnippet: snippet,
				Suggestion:  suggestion,
			})
		}

		// 4. Ausência de Automação DevOps
		if !hasDevOps {
			violations = append(violations, domain.ArchViolation{
				Rule:        "MISSING_DEVOPS_AUTOMATION",
				Severity:    "LOW",
				Message:     "Nenhum arquivo de automação CI/CD, Dockerfile ou Makefile detectado na raiz do projeto.",
				Target:      manifest.ModuleName,
				Line:        1,
				CodeSnippet: "// Ausência de Dockerfile / Makefile / .github/workflows",
				Suggestion:  "Crie um pipeline de CI/CD (.github/workflows) e Dockerfile para padronizar o build.",
			})
		}
	}

	return violations
}
