/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-09-07
 */

package parser

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// CalculateHealthScore calcula a nota consolidada de saúde arquitetural (0 a 100%).
func CalculateHealthScore(violations []domain.ArchViolation) int {
	if len(violations) == 0 {
		return -1
	} // No findings do not establish health.
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
	directCount := 0
	for _, dep := range manifest.Dependencies {
		if !dep.Indirect {
			directCount++
		}
	}
	if directCount > 35 {
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
			isORM = DependencyIdentityMatch(lowerDep, "diesel") || DependencyIdentityMatch(lowerDep, "sqlx") || DependencyIdentityMatch(lowerDep, "sea-orm") || DependencyIdentityMatch(lowerDep, "rusqlite")
		case "python":
			isORM = DependencyIdentityMatch(lowerDep, "sqlalchemy") || DependencyIdentityMatch(lowerDep, "django") || DependencyIdentityMatch(lowerDep, "tortoise") || DependencyIdentityMatch(lowerDep, "peewee")
		case "typescript", "javascript":
			isORM = DependencyIdentityMatch(lowerDep, "prisma") || DependencyIdentityMatch(lowerDep, "typeorm") || DependencyIdentityMatch(lowerDep, "sequelize") || DependencyIdentityMatch(lowerDep, "drizzle")
		case "go":
			isORM = DependencyIdentityMatch(lowerDep, "gorm") || DependencyIdentityMatch(lowerDep, "entgo.io/ent") || DependencyIdentityMatch(lowerDep, "prisma") || DependencyIdentityMatch(lowerDep, "sqlx")
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
					Message:     fmt.Sprintf("Hipótese: dependência de ORM/Banco (%s) declarada; árvore não mostra diretório convencional. Isolamento real não verificado.", dep.Name),
					Target:      dep.Name,
					Line:        0,
					CodeSnippet: "",
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
			snippet := ""
			suggestion := "Verifique testes e o sistema de build da linguagem; a árvore observada pode estar incompleta."
			if lang == "go" {
				suggestion = "Verifique arquivos *_test.go e execute go test."
			}
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
