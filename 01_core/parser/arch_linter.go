/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"strings"

	"bateia/01_core/domain"
)

// EvaluateArchViolations analisa o manifesto e a estrutura de diretórios em busca de anti-padrões arquiteturais.
func EvaluateArchViolations(manifest *domain.Manifest, tree []domain.DirectoryNode) []domain.ArchViolation {
	var violations []domain.ArchViolation

	if manifest == nil {
		return violations
	}

	// 1. Verificar acoplamento excessivo de dependências
	if len(manifest.Dependencies) > 35 {
		violations = append(violations, domain.ArchViolation{
			Rule:        "HIGH_DEPENDENCY_COUPLING",
			Severity:    "HIGH",
			Message:     "Repositório possui acoplamento elevado com mais de 35 dependências diretas.",
			Target:      manifest.ModuleName,
			Line:        1,
			CodeSnippet: "[dependencies] // > 35 dependências diretas",
			Suggestion:  "Divida as responsabilidades em múltiplos crates/módulos acoplados por interfaces genéricas.",
		})
	}

	// 2. Verificar vazamento de camadas em manifestos Go/Node
	for _, dep := range manifest.Dependencies {
		lowerDep := strings.ToLower(dep.Name)
		if strings.Contains(lowerDep, "gorm") || strings.Contains(lowerDep, "ent") || strings.Contains(lowerDep, "prisma") {
			hasInfra := false
			for _, node := range tree {
				if strings.Contains(node.Path, "infra") || strings.Contains(node.Path, "repository") || strings.Contains(node.Path, "database") {
					hasInfra = true
					break
				}
			}
			if !hasInfra && len(tree) > 0 {
				violations = append(violations, domain.ArchViolation{
					Rule:        "MISSING_INFRASTRUCTURE_ISOLATION",
					Severity:    "MEDIUM",
					Message:     "Dependência de ORM/Banco sem isolamento de camada de infraestrutura/repositório detectada.",
					Target:      dep.Name,
					Line:        10,
					CodeSnippet: dep.Name + " // Acesso direto ao ORM/Banco",
					Suggestion:  "Encapsule a dependência de banco dentro de uma camada dedicada (ex: 03_infra/ ou repository/).",
				})
			}
		}
	}

	if len(tree) > 0 {
		hasTests := false
		hasDevOps := false

		for _, node := range tree {
			lowerPath := strings.ToLower(node.Path)
			if strings.Contains(lowerPath, "test") || strings.HasSuffix(lowerPath, "_test.go") || strings.Contains(lowerPath, "spec") {
				hasTests = true
			}
			if strings.Contains(lowerPath, "docker") || strings.Contains(lowerPath, ".github") || strings.Contains(lowerPath, "makefile") {
				hasDevOps = true
			}
		}

		// 3. Ausência de Suíte de Testes
		if !hasTests {
			violations = append(violations, domain.ArchViolation{
				Rule:        "MISSING_TEST_SUITE",
				Severity:    "MEDIUM",
				Message:     "Nenhum diretório ou arquivo de teste unitário/integração detectado na estrutura do projeto.",
				Target:      manifest.ModuleName,
				Line:        1,
				CodeSnippet: "// Ausência de pasta tests/ ou arquivos *_test.go",
				Suggestion:  "Adicione arquivos de testes para garantir cobertura de regressão e integridade arquitetural.",
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
				Suggestion:  "Crie um pipeline de CI/CD (.github/workflows) e Dockerfile para padronizar o build de desenvolvimento.",
			})
		}
	}

	return violations
}
