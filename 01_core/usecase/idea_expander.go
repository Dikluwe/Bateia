/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_expander.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

func normalizeIdeaText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func normalizeIdeaList(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]bool)
	for _, value := range values {
		normalized := normalizeIdeaText(value)
		key := strings.ToLower(normalized)
		if normalized == "" || seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, normalized)
	}
	return result
}

func normalizeIdea(spec domain.IdeaSpec) domain.IdeaSpec {
	return domain.IdeaSpec{
		Problem:        normalizeIdeaText(spec.Problem),
		Users:          normalizeIdeaList(spec.Users),
		Capabilities:   normalizeIdeaList(spec.Capabilities),
		Constraints:    normalizeIdeaList(spec.Constraints),
		NonGoals:       normalizeIdeaList(spec.NonGoals),
		PreferredStack: normalizeIdeaList(spec.PreferredStack),
	}
}

func expansionID(facet domain.ExpansionFacet, query string) string {
	hash := sha256.Sum256([]byte(string(facet) + "\x00" + strings.ToLower(query)))
	return fmt.Sprintf("exp-%s-%x", facet, hash[:6])
}

func appendExpansion(plan *domain.ExpansionPlan, seen map[string]int, facet domain.ExpansionFacet, query, rationale string, covers []string) {
	key := strings.ToLower(normalizeIdeaText(query))
	if key == "" {
		return
	}
	if index, exists := seen[key]; exists {
		plan.Expansions[index].Covers = append(plan.Expansions[index].Covers, covers...)
		return
	}
	seen[key] = len(plan.Expansions)
	plan.Expansions = append(plan.Expansions, domain.SearchExpansion{
		ID: expansionID(facet, query), Facet: facet, Query: query,
		Rationale: rationale, Covers: append([]string(nil), covers...),
	})
}

// ExpandIdea materializa perguntas rastreáveis sem executar descoberta ou inventar sinônimos.
func ExpandIdea(spec domain.IdeaSpec) (*domain.ExpansionPlan, error) {
	normalized := normalizeIdea(spec)
	if normalized.Problem == "" && len(normalized.Capabilities) == 0 {
		return nil, fmt.Errorf("ideia deve informar problem ou ao menos uma capability")
	}

	plan := &domain.ExpansionPlan{Idea: normalized, Expansions: make([]domain.SearchExpansion, 0)}
	seen := make(map[string]int)
	appendExpansion(plan, seen, domain.FacetProblem, normalized.Problem,
		"descobrir soluções para o problema central", []string{"problem"})
	for index, capability := range normalized.Capabilities {
		appendExpansion(plan, seen, domain.FacetCapability, capability,
			"buscar implementação da capacidade isoladamente", []string{fmt.Sprintf("capabilities[%d]", index)})
	}
	for index, constraint := range normalized.Constraints {
		appendExpansion(plan, seen, domain.FacetConstraint, constraint,
			"buscar soluções que declaram a restrição relevante", []string{fmt.Sprintf("constraints[%d]", index)})
	}
	return plan, nil
}
