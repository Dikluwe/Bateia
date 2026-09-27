/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/multicriteria_composition.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import "bateia/01_core/domain"

func multicriteriaCompositions(request domain.MulticriteriaDecisionRequest) (domain.CompositionReport, error) {
	capabilities := normalizedSorted(request.Idea.Capabilities)
	candidates := make([]decisionCandidate, 0, len(request.Assessment.Assessments))
	for _, assessment := range request.Assessment.Assessments {
		candidates = append(candidates, inspectPreparedCandidate(SanitizeAssessmentEvidence(assessment), capabilities, request.Policy))
	}
	return analyzeCompositions(candidates, capabilities, request.Policy, request.Composition)
}

func compositionMemberMap(candidate domain.CompositionCandidate) map[string]domain.CompositionMember {
	members := make(map[string]domain.CompositionMember, len(candidate.Members))
	for _, member := range candidate.Members {
		members[member.ID] = member
	}
	return members
}
