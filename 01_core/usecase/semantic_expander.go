/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/semantic_expander.md
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

type SemanticExpander interface {
	Expand(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error)
}

type EnrichmentResult struct {
	Plan              domain.ExpansionPlan
	ProposedCount     int
	DuplicatesRemoved int
}

func cloneStrings(values []string) []string { return append([]string(nil), values...) }

func cloneDerivation(value *domain.ExpansionDerivation) *domain.ExpansionDerivation {
	if value == nil {
		return nil
	}
	return &domain.ExpansionDerivation{
		Strategy: value.Strategy, RuleIDs: cloneStrings(value.RuleIDs), SourceIDs: cloneStrings(value.SourceIDs),
	}
}

func clonePlan(plan domain.ExpansionPlan) domain.ExpansionPlan {
	cloned := domain.ExpansionPlan{Idea: normalizeIdea(plan.Idea), Expansions: make([]domain.SearchExpansion, len(plan.Expansions))}
	for index, expansion := range plan.Expansions {
		cloned.Expansions[index] = expansion
		cloned.Expansions[index].Covers = cloneStrings(expansion.Covers)
		cloned.Expansions[index].Derivation = cloneDerivation(expansion.Derivation)
	}
	return cloned
}

func semanticExpansionID(strategy string, facet domain.ExpansionFacet, query string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(strategy) + "\x00" + string(facet) + "\x00" + strings.ToLower(query)))
	return fmt.Sprintf("exp-%s-%x", facet, hash[:6])
}

func appendUnique(values []string, additions ...string) []string {
	seen := make(map[string]bool, len(values)+len(additions))
	for _, value := range values {
		seen[value] = true
	}
	for _, value := range additions {
		if value != "" && !seen[value] {
			seen[value] = true
			values = append(values, value)
		}
	}
	return values
}

func containsPhrase(text, phrase string) bool {
	text = strings.ToLower(normalizeIdeaText(text))
	phrase = strings.ToLower(normalizeIdeaText(phrase))
	return phrase != "" && strings.Contains(" "+text+" ", " "+phrase+" ")
}

func blockedByNonGoal(query string, nonGoals []string) bool {
	for _, nonGoal := range nonGoals {
		if containsPhrase(query, nonGoal) || containsPhrase(nonGoal, query) {
			return true
		}
	}
	return false
}

func validateSemanticExpansion(expansion domain.SearchExpansion, sourceIDs map[string]bool) error {
	if normalizeIdeaText(expansion.Query) == "" {
		return fmt.Errorf("expansão semântica possui query vazia")
	}
	if expansion.Derivation == nil || strings.TrimSpace(expansion.Derivation.Strategy) == "" || len(expansion.Derivation.RuleIDs) == 0 || len(expansion.Derivation.SourceIDs) == 0 {
		return fmt.Errorf("expansão semântica %q não possui estratégia, regra e origem", expansion.Query)
	}
	for _, sourceID := range expansion.Derivation.SourceIDs {
		if !sourceIDs[sourceID] {
			return fmt.Errorf("expansão semântica %q referencia origem desconhecida %q", expansion.Query, sourceID)
		}
	}
	return nil
}

// EnrichExpansionPlan preserva o plano literal e incorpora propostas auditáveis.
func EnrichExpansionPlan(plan domain.ExpansionPlan, expander SemanticExpander) (EnrichmentResult, error) {
	result := EnrichmentResult{Plan: clonePlan(plan)}
	if expander == nil {
		return result, nil
	}
	proposals, err := expander.Expand(clonePlan(plan))
	if err != nil {
		return EnrichmentResult{}, err
	}
	result.ProposedCount = len(proposals)

	seen := make(map[string]int, len(result.Plan.Expansions))
	sourceIDs := make(map[string]bool, len(result.Plan.Expansions))
	for index, expansion := range result.Plan.Expansions {
		seen[strings.ToLower(normalizeIdeaText(expansion.Query))] = index
		sourceIDs[expansion.ID] = true
	}
	for _, proposal := range proposals {
		proposal.Query = normalizeIdeaText(proposal.Query)
		if err := validateSemanticExpansion(proposal, sourceIDs); err != nil {
			return EnrichmentResult{}, err
		}
		if blockedByNonGoal(proposal.Query, result.Plan.Idea.NonGoals) {
			continue
		}
		key := strings.ToLower(proposal.Query)
		if index, exists := seen[key]; exists {
			result.DuplicatesRemoved++
			existing := &result.Plan.Expansions[index]
			existing.Covers = appendUnique(existing.Covers, proposal.Covers...)
			if existing.Derivation == nil {
				existing.Derivation = cloneDerivation(proposal.Derivation)
			} else {
				existing.Derivation.RuleIDs = appendUnique(existing.Derivation.RuleIDs, proposal.Derivation.RuleIDs...)
				existing.Derivation.SourceIDs = appendUnique(existing.Derivation.SourceIDs, proposal.Derivation.SourceIDs...)
			}
			continue
		}
		proposal.ID = semanticExpansionID(proposal.Derivation.Strategy, proposal.Facet, proposal.Query)
		proposal.Covers = appendUnique(nil, proposal.Covers...)
		proposal.Derivation = cloneDerivation(proposal.Derivation)
		seen[key] = len(result.Plan.Expansions)
		result.Plan.Expansions = append(result.Plan.Expansions, proposal)
	}
	return result, nil
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

func foundConcepts(expansions []domain.SearchExpansion, concepts []string) []string {
	queries := make(map[string]bool, len(expansions))
	for _, expansion := range expansions {
		queries[strings.ToLower(normalizeIdeaText(expansion.Query))] = true
	}
	found := make([]string, 0)
	for _, concept := range concepts {
		if queries[strings.ToLower(normalizeIdeaText(concept))] {
			found = append(found, normalizeIdeaText(concept))
		}
	}
	return found
}

// EvaluateExpansionStrategy mede uma estratégia sobre um corpus sem executar I/O.
func EvaluateExpansionStrategy(corpus domain.EvaluationCorpus, strategy string, expander SemanticExpander) (domain.EvaluationReport, error) {
	if len(corpus.Cases) == 0 {
		return domain.EvaluationReport{}, fmt.Errorf("corpus deve conter ao menos um caso")
	}
	report := domain.EvaluationReport{Strategy: strategy, Cases: make([]domain.CaseEvaluation, 0, len(corpus.Cases))}
	var expectedTotal, expectedFound, forbiddenTotal, forbiddenFound, expansionTotal, newTotal, proposedTotal, duplicateTotal int
	for _, evaluationCase := range corpus.Cases {
		if strings.TrimSpace(evaluationCase.ID) == "" {
			return domain.EvaluationReport{}, fmt.Errorf("caso de avaliação sem id")
		}
		literal, err := ExpandIdea(evaluationCase.Idea)
		if err != nil {
			return domain.EvaluationReport{}, fmt.Errorf("caso %s inválido: %w", evaluationCase.ID, err)
		}
		enriched, err := EnrichExpansionPlan(*literal, expander)
		if err != nil {
			return domain.EvaluationReport{}, fmt.Errorf("caso %s: %w", evaluationCase.ID, err)
		}
		foundExpected := foundConcepts(enriched.Plan.Expansions, evaluationCase.ExpectedConcepts)
		foundForbidden := foundConcepts(enriched.Plan.Expansions, evaluationCase.ForbiddenConcepts)
		newCount := len(enriched.Plan.Expansions) - len(literal.Expansions)
		caseResult := domain.CaseEvaluation{
			ID: evaluationCase.ID, FoundExpected: foundExpected, FoundForbidden: foundForbidden,
			ExpansionCount: len(enriched.Plan.Expansions), NewExpansionCount: newCount, ProposedExpansionCount: enriched.ProposedCount,
			Metrics: domain.ExpansionMetrics{
				Coverage:      ratio(len(foundExpected), len(evaluationCase.ExpectedConcepts)),
				Contamination: ratio(len(foundForbidden), len(evaluationCase.ForbiddenConcepts)),
				Novelty:       ratio(newCount, len(enriched.Plan.Expansions)),
				Duplication:   ratio(enriched.DuplicatesRemoved, enriched.ProposedCount),
			},
		}
		report.Cases = append(report.Cases, caseResult)
		expectedTotal += len(evaluationCase.ExpectedConcepts)
		expectedFound += len(foundExpected)
		forbiddenTotal += len(evaluationCase.ForbiddenConcepts)
		forbiddenFound += len(foundForbidden)
		expansionTotal += len(enriched.Plan.Expansions)
		newTotal += newCount
		proposedTotal += enriched.ProposedCount
		duplicateTotal += enriched.DuplicatesRemoved
	}
	report.Metrics = domain.ExpansionMetrics{
		Coverage: ratio(expectedFound, expectedTotal), Contamination: ratio(forbiddenFound, forbiddenTotal),
		Novelty: ratio(newTotal, expansionTotal), Duplication: ratio(duplicateTotal, proposedTotal),
	}
	return report, nil
}
