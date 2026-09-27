/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/systematic_review_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func reviewFixture(t *testing.T) (domain.ReviewProtocol, domain.ReviewSynthesisInput) {
	t.Helper()
	intake := domain.RecommendationIntake{
		SchemaVersion: "1", Status: domain.IntakeReady, InputHash: "intake-hash",
		Request: domain.RecommendationRequest{
			Idea:   domain.IdeaSpec{Problem: "index code", Capabilities: []string{"index"}},
			Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, ExtendMinCoverage: .8, IntegrationMinValue: .5},
			Target: domain.SoftwareAlternativeTarget{Name: "Indexer"},
			Search: domain.RecommendationSearch{Strategy: "lexical", Language: "go", Locales: []string{"en", "pt-BR"}, MaxPerQuery: 10, MaxCandidates: 20},
		},
		Execution: domain.RecommendationExecutionProfile{Mode: "replay", Provider: "fixture", LocalPaths: []string{}},
	}
	protocol, err := FreezeReviewProtocol(domain.ReviewProtocolDraft{Intake: intake, ScreeningModel: DefaultScreeningModel(), EvidenceRulesVersion: "rules-v1", EvidenceRulesHash: "rules-hash"})
	if err != nil {
		t.Fatal(err)
	}
	a := domain.RepositoryCandidate{Owner: "Acme", Repo: "Alpha", CommitID: "abc"}
	b := domain.RepositoryCandidate{Owner: "Acme", Repo: "Beta", CommitID: "def"}
	discovery := domain.DiscoverySnapshot{
		SchemaVersion: "2", Idea: intake.Request.Idea, Strategy: "lexical", SearchLocales: []string{"en", "pt-BR"},
		PerExpansion: []domain.ExpansionDiscovery{
			{Expansion: domain.SearchExpansion{ID: "q1", Query: "code index"}, Provider: domain.ProviderDiscovery{Provider: "fixture", Candidates: []domain.RepositoryCandidate{a, b}}},
			{Expansion: domain.SearchExpansion{ID: "q2", Query: "source index"}, Provider: domain.ProviderDiscovery{Provider: "fixture", Candidates: []domain.RepositoryCandidate{a}}},
		},
		Candidates: []domain.DiscoveredCandidate{{Repository: a}, {Repository: b}},
	}
	source := domain.EvidenceSource{Kind: domain.EvidenceTestImplementation, Ref: domain.CommitRef{Owner: "acme", Repo: "alpha", Commit: "abc"}, Path: "index_test.go", RuleID: "cap-index", MatchClass: domain.EvidenceMatchStrongContent, Strength: 1}
	assessment := domain.CandidateAssessment{
		Repository:       domain.CommitRef{Owner: "acme", Repo: "alpha", Commit: "abc"},
		ProblemAlignment: domain.CapabilityEvidence{Capability: "problem", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}},
		Capabilities:     []domain.CapabilityEvidence{{Capability: "index", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}}},
		License:          domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{{Kind: domain.EvidenceLicense, Ref: source.Ref, Path: "LICENSE", RuleID: "license", MatchClass: domain.EvidenceMatchStrongContent, Strength: 1}}},
		Maintenance:      domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{{Kind: domain.EvidenceRepository, Ref: source.Ref, RuleID: "repo", MatchClass: domain.EvidenceMatchStrongContent, Strength: 1}}},
		Operation:        domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}},
	}
	alternative := domain.SoftwareAlternative{Repository: assessment.Repository, Fit: domain.AlternativeRecommended, License: assessment.License, Maintenance: assessment.Maintenance, Operation: assessment.Operation}
	synthesis := domain.ReviewSynthesisInput{
		Protocol: protocol, Discovery: discovery,
		Screening:       domain.ScreeningReport{SchemaVersion: "1", Strategy: "probabilistic", Decisions: []domain.ScreeningDecision{{Repository: "acme/alpha", RankingScore: .9, Decision: "ADVANCE"}, {Repository: "acme/beta", RankingScore: .1, Decision: "DROP"}}},
		Assessment:      domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{assessment}},
		EvidenceBundles: []domain.EvidenceBundle{{SchemaVersion: "1", Repository: assessment.Repository, Complete: true, TreeComplete: true, RequirementCoverage: []domain.RequirementCollectionCoverage{{SearchComplete: true}}}},
		Decision:        domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionExtend},
		Alternatives:    domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{alternative}, Excluded: []domain.SoftwareAlternative{}},
		SearchedThrough: "2026-09-24",
	}
	return protocol, synthesis
}

func TestSystematicReviewBuildsConservativeLedgerAndFlow(t *testing.T) {
	_, input := reviewFixture(t)
	pack, err := SynthesizeSystematicReview(input)
	if err != nil {
		t.Fatal(err)
	}
	if pack.Flow.IdentifiedRecords != 3 || pack.Flow.DuplicateRecords != 1 || pack.Flow.UniqueCandidates != 2 {
		t.Fatalf("fluxo inesperado: %+v", pack.Flow)
	}
	if !pack.Flow.ConservationChecked || pack.Flow.Included != 1 || pack.Flow.Excluded != 1 || pack.Flow.Inconclusive != 0 {
		t.Fatalf("disposições não fecham: %+v", pack.Flow)
	}
	if pack.Ledger.Candidates[0].ReasonCode != domain.ReasonIncludedRecommended || pack.Ledger.Candidates[1].ReasonCode != domain.ReasonScreeningDrop {
		t.Fatalf("reason codes não são estáveis: %+v", pack.Ledger.Candidates)
	}
	if len(pack.Certainty.Claims) == 0 || pack.Certainty.Claims[0].ObservedStatus == "" || pack.Certainty.Claims[0].Level == "" {
		t.Fatalf("certeza por claim ausente: %+v", pack.Certainty)
	}
}

func TestSystematicReviewRejectsMutatedProtocol(t *testing.T) {
	_, input := reviewFixture(t)
	input.Protocol.Search.MaxCandidates++
	if _, err := SynthesizeSystematicReview(input); err == nil {
		t.Fatal("protocolo alterado depois do congelamento foi aceito")
	}
}

func TestLivingReviewAppendsAndSeparatesProtocolLineages(t *testing.T) {
	_, input := reviewFixture(t)
	first, err := SynthesizeSystematicReview(input)
	if err != nil {
		t.Fatal(err)
	}
	index, err := UpdateLivingReview(domain.LivingReviewUpdateInput{CurrentPackage: first})
	if err != nil {
		t.Fatal(err)
	}
	input.SearchedThrough = "2026-10-01"
	input.Discovery.Candidates = append(input.Discovery.Candidates, domain.DiscoveredCandidate{Repository: domain.RepositoryCandidate{Owner: "acme", Repo: "gamma", CommitID: "ghi"}})
	second, err := SynthesizeSystematicReview(input)
	if err != nil {
		t.Fatal(err)
	}
	index, err = UpdateLivingReview(domain.LivingReviewUpdateInput{Index: index, PreviousPackage: &first, CurrentPackage: second})
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Lineages) != 1 || len(index.Lineages[0].Versions) != 2 || len(index.Lineages[0].Versions[1].Delta.AddedCandidates) != 1 {
		t.Fatalf("índice vivo inesperado: %+v", index)
	}
	third := second
	third.Protocol.Question.Problem = "changed question"
	third.Protocol.ProtocolHash = protocolDigest(third.Protocol)
	third.PackageHash = packageDigest(third)
	index, err = UpdateLivingReview(domain.LivingReviewUpdateInput{Index: index, CurrentPackage: third})
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Lineages) != 2 {
		t.Fatalf("mudança de protocolo deveria abrir linhagem: %+v", index.Lineages)
	}
}

func TestActiveLearningOnlyPrioritizesAndReplansFromRevealedLabels(t *testing.T) {
	report := activeLearningReport()
	initial, err := PrioritizeReviewCandidates(domain.ActiveLearningRequest{Screening: report, BatchSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if initial.CanExclude || initial.Authority != "PRIORITIZATION_ONLY" || len(initial.Priorities) != 2 {
		t.Fatalf("autoridade indevida: %+v", initial)
	}
	after, err := PrioritizeReviewCandidates(domain.ActiveLearningRequest{Screening: report, Reviewed: []domain.CandidateJudgment{{Repository: initial.Priorities[0].Repository, Label: domain.RelevanceRelevant}}, BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Priorities) != 1 || after.Priorities[0].Repository == initial.Priorities[0].Repository {
		t.Fatalf("replanejamento não retirou item revelado: %+v", after)
	}
}

func TestActiveLearningEvaluationNeverClaimsStoppingAuthority(t *testing.T) {
	report := activeLearningReport()
	gold := []domain.CandidateJudgment{{Repository: "acme/a", Label: domain.RelevanceRelevant}, {Repository: "acme/b", Label: domain.RelevanceIrrelevant}, {Repository: "acme/c", Label: domain.RelevanceRelevant}, {Repository: "acme/d", Label: domain.RelevanceIrrelevant}}
	evaluation, err := EvaluateActiveLearning(domain.ActiveLearningCorpus{SchemaVersion: "1", Cases: []domain.ActiveLearningCase{{ID: "independent", Screening: report, Gold: gold}}})
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.StoppingAuthority || len(evaluation.Cases) != 1 {
		t.Fatalf("avaliação ganhou autoridade operacional: %+v", evaluation)
	}
}

func activeLearningReport() domain.ScreeningReport {
	decisions := []domain.ScreeningDecision{}
	for index, score := range []float64{.9, .6, .4, .1} {
		name := string(rune('a' + index))
		decisions = append(decisions, domain.ScreeningDecision{Repository: "acme/" + name, RankingScore: score, Decision: "AUDIT", Signals: domain.ScreeningSignals{DescriptionAlignment: score, CapabilityCoverage: score, EvidenceCompleteness: .5}, Origins: []domain.CandidateOrigin{{ExpansionID: "q" + name}}})
	}
	return domain.ScreeningReport{SchemaVersion: "1", Strategy: "fixture", Decisions: decisions}
}
