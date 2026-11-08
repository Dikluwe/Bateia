/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/alternative_target.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"strings"
)

// AlternativeIdentityReader provides advisory identities, including frozen capture hints.
type AlternativeIdentityReader interface {
	RepositoryIdentities(string) []domain.RepositoryIdentity
}
type localAlternativeBinding struct {
	repository *domain.RepositoryIdentity
}

func AlternativeInputsForTarget(mode domain.IdeaRunMode, discovery domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, target *domain.RepositoryIdentity, identitiesFor AlternativeIdentityReader) (domain.DiscoverySnapshot, []domain.CandidateAssessment) {
	discovery.Candidates = append([]domain.DiscoveredCandidate(nil), discovery.Candidates...)
	assessments = append([]domain.CandidateAssessment(nil), assessments...)
	if mode != domain.IdeaRunLocal || target == nil {
		return discovery, assessments
	}
	bindings := map[string]localAlternativeBinding{}
	for i := range discovery.Candidates {
		candidate := discovery.Candidates[i].Repository
		key := alternativeIdentityKey(candidate.Owner, candidate.Repo)
		identities := identitiesFor.RepositoryIdentities(candidate.URL)
		if len(identities) == 0 {
			discovery.Candidates[i].Repository.Stale = true
			continue
		}
		for _, identity := range identities {
			if alternativeIdentityKey(identity.Owner, identity.Repo) != alternativeIdentityKey(target.Owner, target.Repo) {
				continue
			}
			bound := domain.RepositoryIdentity{Owner: strings.ToLower(strings.TrimSpace(target.Owner)), Repo: strings.ToLower(strings.TrimSpace(target.Repo))}
			discovery.Candidates[i].Repository.Owner = bound.Owner
			discovery.Candidates[i].Repository.Repo = bound.Repo
			bindings[key] = localAlternativeBinding{repository: &bound}
			break
		}
	}
	for i := range assessments {
		key := alternativeIdentityKey(assessments[i].Repository.Owner, assessments[i].Repository.Repo)
		if binding := bindings[key]; binding.repository != nil {
			assessments[i] = rebindAlternativeAssessment(assessments[i], *binding.repository)
		}
	}
	return discovery, assessments
}

func alternativeIdentityKey(owner, repo string) string {
	return strings.ToLower(strings.TrimSpace(owner)) + "/" + strings.ToLower(strings.TrimSpace(repo))
}

func rebindAlternativeAssessment(assessment domain.CandidateAssessment, repository domain.RepositoryIdentity) domain.CandidateAssessment {
	assessment.Repository.Owner = repository.Owner
	assessment.Repository.Repo = repository.Repo
	assessment.ScreeningDecision.Repository = assessment.Repository.Key()
	assessment.ProblemAlignment.Sources = rebindAlternativeSources(assessment.ProblemAlignment.Sources, repository)
	assessment.Capabilities = append([]domain.CapabilityEvidence(nil), assessment.Capabilities...)
	for i := range assessment.Capabilities {
		assessment.Capabilities[i].Sources = rebindAlternativeSources(assessment.Capabilities[i].Sources, repository)
	}
	assessment.License.Sources = rebindAlternativeSources(assessment.License.Sources, repository)
	assessment.Maintenance.Sources = rebindAlternativeSources(assessment.Maintenance.Sources, repository)
	assessment.Operation.Sources = rebindAlternativeSources(assessment.Operation.Sources, repository)
	return assessment
}

func rebindAlternativeSources(sources []domain.EvidenceSource, repository domain.RepositoryIdentity) []domain.EvidenceSource {
	out := append([]domain.EvidenceSource(nil), sources...)
	for i := range out {
		out[i].Ref.Owner = repository.Owner
		out[i].Ref.Repo = repository.Repo
	}
	return out
}
