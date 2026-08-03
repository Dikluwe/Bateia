/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/arch_linter.md
 * @layer L1
 * @updated 2026-08-03
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
			Rule:     "HIGH_DEPENDENCY_COUPLING",
			Severity: "HIGH",
			Message:  "Repositório possui acoplamento elevado com mais de 35 dependências diretas.",
			Target:   manifest.ModuleName,
		})
	}

	// 2. Verificar vazamento de camadas em manifestos Go/Node
	for _, dep := range manifest.Dependencies {
		lowerDep := strings.ToLower(dep.Name)
		if strings.Contains(lowerDep, "gorm") || strings.Contains(lowerDep, "ent") || strings.Contains(lowerDep, "prisma") {
			// Se o banco for ORM pesado sem pasta de infraestrutura/repositório
			hasInfra := false
			for _, node := range tree {
				if strings.Contains(node.Path, "infra") || strings.Contains(node.Path, "repository") || strings.Contains(node.Path, "database") {
					hasInfra = true
					break
				}
			}
			if !hasInfra && len(tree) > 0 {
				violations = append(violations, domain.ArchViolation{
					Rule:     "MISSING_INFRASTRUCTURE_ISOLATION",
					Severity: "MEDIUM",
					Message:  "Dependência de ORM/Banco sem isolamento de camada de infraestrutura/repositório detectada.",
					Target:   dep.Name,
				})
			}
		}
	}

	return violations
}
