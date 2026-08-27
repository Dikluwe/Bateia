/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_run_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"testing"
)

func replayRunFixture() (domain.IdeaRunRequest, IdeaRunDependencies) {
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	candidate := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc123", Description: "dependency graph analysis"}
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal", Candidates: []domain.DiscoveredCandidate{{Repository: candidate, ExpansionIDs: []string{"capability-dependency-graph"}, Queries: []string{"dependency graph"}, Origins: []domain.CandidateOrigin{{ExpansionID: "capability-dependency-graph", PlannedQuery: "dependency graph", EffectiveQuery: "dependency graph", Evidence: []domain.MatchEvidence{{MatchType: "content", Path: "graph.go", Snippet: "DependencyGraph"}}}}}}}
	fileContent, testContent, licenseContent := "type DependencyGraph struct{}", "func TestDependencyGraph(){}", "MIT License"
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceFile, Ref: ref, Path: "graph.go", Content: fileContent, ContentHash: contentHash(fileContent)}, {Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: testContent, ContentHash: contentHash(testContent)}, {Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: licenseContent, ContentHash: contentHash(licenseContent)}}, Targets: []domain.CollectionTargetResult{}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	req := domain.IdeaRunRequest{Idea: idea, Policy: policy, Mode: domain.IdeaRunReplay, Strategy: "literal", ReplayDirectory: "source", OutputDirectory: "out"}
	return req, IdeaRunDependencies{ReplaySnapshot: &snapshot, ReplayBundles: []domain.EvidenceBundle{bundle}}
}
func TestRunIdeaPipelineReplayProducesWholeChainWithoutProviders(t *testing.T) {
	req, deps := replayRunFixture()
	got, err := RunIdeaPipeline(context.Background(), req, deps)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Expansion.Expansions) == 0 || len(got.Screening.Decisions) == 0 || len(got.Plans) == 0 || len(got.Assessment.Assessments) == 0 || got.Decision.Decision == "" {
		t.Fatalf("incomplete chain: %#v", got)
	}
}
func TestRunIdeaPipelineRejectsMissingReplayBundle(t *testing.T) {
	req, deps := replayRunFixture()
	deps.ReplayBundles = nil
	if _, err := RunIdeaPipeline(context.Background(), req, deps); err == nil {
		t.Fatal("missing bundle accepted")
	}
}
func TestIdeaRunStatusNeverPromotesUnknown(t *testing.T) {
	if got := IdeaRunStatusFor(domain.IdeaRunResult{Decision: domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}}); got != domain.IdeaRunInconclusive {
		t.Fatalf("got %s", got)
	}
}

func TestUnresolvedBundleBecomesInconclusiveInsteadOfBuild(t *testing.T) {
	req, deps := replayRunFixture()
	deps.ReplayBundles = []domain.EvidenceBundle{{SchemaVersion: "1", Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "HEAD"}, Plan: domain.EvidenceExtractionPlan{Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "HEAD"}, Decision: domain.ScreeningDecision{Repository: "org/repo"}}, Failures: []domain.ExecutionFailure{{Scope: "org/repo", Stage: "resolution", Message: "failed"}}}}
	got, err := RunIdeaPipeline(context.Background(), req, deps)
	if err != nil {
		t.Fatal(err)
	}
	if got.Decision.Decision != domain.DecisionUnknown {
		t.Fatalf("got %s", got.Decision.Decision)
	}
}
