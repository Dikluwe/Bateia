/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_run_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
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

type maxCandidatesDiscoveryProvider struct {
	candidates []domain.RepositoryCandidate
}

func (p maxCandidatesDiscoveryProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	matches := make([]domain.ProviderCandidateMatch, 0, len(p.candidates))
	for _, candidate := range p.candidates {
		matches = append(matches, domain.ProviderCandidateMatch{
			Repository:     candidate.Owner + "/" + candidate.Repo,
			EffectiveQuery: request.Query,
			Evidence:       []domain.MatchEvidence{{MatchType: "symbol", Path: "graph.go", Snippet: "DependencyGraph"}},
		})
	}
	return domain.ProviderDiscovery{
		Provider:   "max-candidates-fake",
		Attempts:   []domain.ProviderAttempt{{Query: request.Query, Status: "SUCCEEDED", ResultCount: len(p.candidates)}},
		Candidates: append([]domain.RepositoryCandidate(nil), p.candidates...),
		Matches:    matches,
	}, nil
}

type maxCandidatesContentProvider struct {
	mu        sync.Mutex
	collected []string
}

func (*maxCandidatesContentProvider) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("manifest not used by idea run")
}

func (*maxCandidatesContentProvider) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{
		{Path: "README.md", Type: "file"},
		{Path: "LICENSE", Type: "file"},
		{Path: "graph.go", Type: "file"},
		{Path: "graph_test.go", Type: "file"},
	}, nil
}

func (p *maxCandidatesContentProvider) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	p.mu.Lock()
	p.collected = append(p.collected, candidate.Owner+"/"+candidate.Repo)
	p.mu.Unlock()
	return &domain.ExtractedRepository{Files: []domain.ExtractedFile{
		{Path: "README.md", Content: "dependency graph analysis"},
		{Path: "LICENSE", Content: "MIT License"},
		{Path: "graph.go", Content: "type DependencyGraph struct{}"},
		{Path: "graph_test.go", Content: "func TestDependencyGraph(){}"},
	}}, nil
}

func (*maxCandidatesContentProvider) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	if candidate.CommitID == "" || candidate.CommitID == "HEAD" {
		candidate.CommitID = "sha-" + candidate.Repo
	}
	return candidate, nil
}

func (p *maxCandidatesContentProvider) collectedRepositories() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.collected...)
}

func maxCandidatesIdea() domain.IdeaSpec {
	return domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
}

func maxCandidatesPolicy() domain.ReuseDecisionPolicy {
	return domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
}

func maxCandidatesRepository(repo string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{
		Owner:       "org",
		Repo:        repo,
		CommitID:    "sha-" + repo,
		Description: "dependency graph analysis",
	}
}

func maxCandidatesSnapshot(repositories ...string) domain.DiscoverySnapshot {
	idea := maxCandidatesIdea()
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal"}
	for _, repository := range repositories {
		snapshot.Candidates = append(snapshot.Candidates, domain.DiscoveredCandidate{
			Repository:   maxCandidatesRepository(repository),
			ExpansionIDs: []string{"dependency-graph"},
			Queries:      []string{"dependency graph"},
			Origins: []domain.CandidateOrigin{{
				ExpansionID:    "dependency-graph",
				PlannedQuery:   "dependency graph",
				EffectiveQuery: "dependency graph",
				Evidence:       []domain.MatchEvidence{{MatchType: "symbol", Path: "graph.go", Snippet: "DependencyGraph"}},
			}},
		})
	}
	return snapshot
}

func maxCandidatesBundle(repository string) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "org", Repo: repository, Commit: "sha-" + repository}
	fileContent := "type DependencyGraph struct{}"
	testContent := "func TestDependencyGraph(){}"
	licenseContent := "MIT License"
	return domain.EvidenceBundle{
		SchemaVersion: "1",
		Repository:    ref,
		Complete:      true,
		TreeComplete:  true,
		Items: []domain.CollectedEvidence{
			{Kind: domain.EvidenceFile, Ref: ref, Path: "graph.go", Content: fileContent, ContentHash: contentHash(fileContent)},
			{Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: testContent, ContentHash: contentHash(testContent)},
			{Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: licenseContent, ContentHash: contentHash(licenseContent)},
		},
	}
}

func maxCandidatesReplayRequest(snapshot domain.DiscoverySnapshot, limit int) domain.IdeaRunRequest {
	return domain.IdeaRunRequest{
		Idea:            snapshot.Idea,
		Policy:          maxCandidatesPolicy(),
		Mode:            domain.IdeaRunReplay,
		Strategy:        "literal",
		ReplayDirectory: "frozen-replay",
		OutputDirectory: "unused-output",
		MaxCandidates:   limit,
	}
}

func commitRefIdentity(ref domain.CommitRef) string {
	return ref.Owner + "/" + ref.Repo
}

func assertMaxCandidateRepositories(t *testing.T, result domain.IdeaRunResult, expected []string) {
	t.Helper()
	if len(result.Plans) != len(expected) || len(result.Bundles) != len(expected) || len(result.Assessment.Assessments) != len(expected) {
		t.Fatalf("limite não atravessou planos/bundles/assessments: plans=%d bundles=%d assessments=%d expected=%d", len(result.Plans), len(result.Bundles), len(result.Assessment.Assessments), len(expected))
	}
	for index, repository := range expected {
		want := "org/" + repository
		if got := commitRefIdentity(result.Plans[index].Repository); got != want {
			t.Errorf("plans[%d] = %s, want %s", index, got, want)
		}
		if got := commitRefIdentity(result.Bundles[index].Repository); got != want {
			t.Errorf("bundles[%d] = %s, want %s", index, got, want)
		}
		if got := commitRefIdentity(result.Assessment.Assessments[index].Repository); got != want {
			t.Errorf("assessments[%d] = %s, want %s", index, got, want)
		}
	}
}

func TestRunIdeaPipelineMaxCandidatesLimitsLiveCollectionAfterScreening(t *testing.T) {
	discovered := []domain.RepositoryCandidate{
		maxCandidatesRepository("delta"),
		maxCandidatesRepository("charlie"),
		maxCandidatesRepository("bravo"),
		maxCandidatesRepository("alpha"),
	}
	content := &maxCandidatesContentProvider{}
	req := domain.IdeaRunRequest{
		Idea:            maxCandidatesIdea(),
		Policy:          maxCandidatesPolicy(),
		Mode:            domain.IdeaRunLive,
		Strategy:        "literal",
		OutputDirectory: "unused-output",
		MaxPerQuery:     10,
		MaxCandidates:   2,
	}
	got, err := RunIdeaPipeline(context.Background(), req, IdeaRunDependencies{
		Discovery: maxCandidatesDiscoveryProvider{candidates: discovered},
		Content:   content,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Screening.Decisions) != len(discovered) {
		t.Fatalf("o limite foi aplicado antes da peneira: decisions=%d discovered=%d", len(got.Screening.Decisions), len(discovered))
	}
	assertMaxCandidateRepositories(t, got, []string{"alpha", "bravo"})
	collected := content.collectedRepositories()
	if fmt.Sprint(collected) != fmt.Sprint([]string{"org/alpha", "org/bravo"}) {
		t.Fatalf("provedor coletou candidatos fora do corte do planejador: got=%v", collected)
	}
}

func TestRunIdeaPipelineZeroMaxCandidatesDefaultsToFirstFiveInReplay(t *testing.T) {
	snapshot := maxCandidatesSnapshot("foxtrot", "bravo", "echo", "alpha", "delta", "charlie")
	bundles := []domain.EvidenceBundle{
		maxCandidatesBundle("foxtrot"),
		maxCandidatesBundle("charlie"),
		maxCandidatesBundle("alpha"),
		maxCandidatesBundle("echo"),
		maxCandidatesBundle("bravo"),
		maxCandidatesBundle("delta"),
	}
	got, err := RunIdeaPipeline(context.Background(), maxCandidatesReplayRequest(snapshot, 0), IdeaRunDependencies{
		ReplaySnapshot: &snapshot,
		ReplayBundles:  bundles,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Screening.Decisions) != 6 {
		t.Fatalf("o default foi aplicado antes da peneira: decisions=%d", len(got.Screening.Decisions))
	}
	assertMaxCandidateRepositories(t, got, []string{"alpha", "bravo", "charlie", "delta", "echo"})
}

func TestRunIdeaPipelineReplayExtraBundleCannotReintroduceCandidateOrPromoteUnknown(t *testing.T) {
	snapshot := maxCandidatesSnapshot("bravo", "alpha")
	unresolved := domain.EvidenceBundle{
		SchemaVersion: "1",
		Repository:    domain.CommitRef{Owner: "org", Repo: "alpha", Commit: "HEAD"},
		Plan: domain.EvidenceExtractionPlan{
			Repository: domain.CommitRef{Owner: "org", Repo: "alpha", Commit: "HEAD"},
			Decision:   domain.ScreeningDecision{Repository: "org/alpha"},
		},
		Failures: []domain.ExecutionFailure{{Scope: "org/alpha", Stage: "resolution", Message: "failed"}},
	}
	got, err := RunIdeaPipeline(context.Background(), maxCandidatesReplayRequest(snapshot, 1), IdeaRunDependencies{
		ReplaySnapshot: &snapshot,
		ReplayBundles: []domain.EvidenceBundle{
			maxCandidatesBundle("bravo"),
			unresolved,
			maxCandidatesBundle("intruder"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertMaxCandidateRepositories(t, got, []string{"alpha"})
	if got.Decision.Decision != domain.DecisionUnknown {
		t.Fatalf("bundle forte extra promoveu Unknown para %s", got.Decision.Decision)
	}
}

func TestRunIdeaPipelineReplayMissingBundleKeepsLimitedCandidateInconclusive(t *testing.T) {
	snapshot := maxCandidatesSnapshot("bravo", "alpha", "charlie")
	got, err := RunIdeaPipeline(context.Background(), maxCandidatesReplayRequest(snapshot, 2), IdeaRunDependencies{
		ReplaySnapshot: &snapshot,
		ReplayBundles:  []domain.EvidenceBundle{maxCandidatesBundle("alpha")},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertMaxCandidateRepositories(t, got, []string{"alpha", "bravo"})
	missing := got.Bundles[1]
	if missing.Complete || len(missing.Failures) == 0 {
		t.Fatalf("candidato sem bundle desapareceu como falha de replay: %+v", missing)
	}
	if got.Decision.Decision != domain.DecisionUnknown || IdeaRunStatusFor(got) != domain.IdeaRunInconclusive {
		t.Fatalf("bundle ausente não preservou Unknown/INCONCLUSIVE: decision=%s status=%s", got.Decision.Decision, IdeaRunStatusFor(got))
	}
}

func TestRunIdeaPipelineReplayRejectsSameRepositoryAtDifferentCommit(t *testing.T) {
	for _, wrongCommit := range []string{"sha-other", "HEAD"} {
		t.Run(wrongCommit, func(t *testing.T) {
			snapshot := maxCandidatesSnapshot("alpha")
			wrongBundle := maxCandidatesBundle("alpha")
			wrongBundle.Repository.Commit = wrongCommit
			for index := range wrongBundle.Items {
				wrongBundle.Items[index].Ref.Commit = wrongCommit
			}

			got, err := RunIdeaPipeline(context.Background(), maxCandidatesReplayRequest(snapshot, 1), IdeaRunDependencies{
				ReplaySnapshot: &snapshot,
				ReplayBundles:  []domain.EvidenceBundle{wrongBundle},
			})
			if err != nil {
				t.Fatal(err)
			}
			assertMaxCandidateRepositories(t, got, []string{"alpha"})

			planned := got.Plans[0].Repository
			if planned.Commit != "sha-alpha" {
				t.Fatalf("plano perdeu identidade imutável: %+v", planned)
			}
			if got.Bundles[0].Repository != planned || got.Assessment.Assessments[0].Repository != planned {
				t.Fatalf("candidato não permaneceu observável na identidade planejada: plan=%+v bundle=%+v assessment=%+v", planned, got.Bundles[0].Repository, got.Assessment.Assessments[0].Repository)
			}
			if got.Bundles[0].Complete || len(got.Bundles[0].Failures) == 0 {
				t.Fatalf("mismatch de commit não virou falha observável: %+v", got.Bundles[0])
			}
			if got.Decision.Decision != domain.DecisionUnknown || IdeaRunStatusFor(got) != domain.IdeaRunInconclusive {
				t.Fatalf("evidência de outro commit promoveu decisão: decision=%s status=%s", got.Decision.Decision, IdeaRunStatusFor(got))
			}
		})
	}
}

const m18WitnessCommit = "120aa5c8816f9a17e043d48c653a6ce200937ea9"

type m18ContentProvider struct {
	resolved      string
	resolutionErr error
}

func (m18ContentProvider) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("manifest not used by M18")
}

func (m18ContentProvider) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{
		{Path: "README.md", Type: "file"},
		{Path: "LICENSE", Type: "file"},
		{Path: "graph.go", Type: "file"},
		{Path: "graph_test.go", Type: "file"},
	}, nil
}

func (m18ContentProvider) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{Candidate: candidate, Files: []domain.ExtractedFile{
		{Path: "README.md", Content: "dependency graph analysis"},
		{Path: "LICENSE", Content: "MIT License"},
		{Path: "graph.go", Content: "type DependencyGraph struct{}"},
		{Path: "graph_test.go", Content: "func TestDependencyGraph(){}"},
	}}, nil
}

func (provider m18ContentProvider) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	if provider.resolutionErr != nil {
		return domain.RepositoryCandidate{}, provider.resolutionErr
	}
	if provider.resolved != "" {
		candidate.CommitID = provider.resolved
	}
	return candidate, nil
}

func m18Candidate(repo, commit string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{
		Owner: "julianpoy", Repo: repo, CommitID: commit,
		URL: "https://github.com/julianpoy/" + repo, Description: "dependency graph analysis",
	}
}

func m18Snapshot(candidates ...domain.RepositoryCandidate) domain.DiscoverySnapshot {
	idea := maxCandidatesIdea()
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal"}
	for _, candidate := range candidates {
		snapshot.Candidates = append(snapshot.Candidates, domain.DiscoveredCandidate{
			Repository:   candidate,
			ExpansionIDs: []string{"dependency-graph"},
			Queries:      []string{"dependency graph"},
			Origins: []domain.CandidateOrigin{{
				ExpansionID: "dependency-graph", PlannedQuery: "dependency graph", EffectiveQuery: "dependency graph",
				Evidence: []domain.MatchEvidence{{MatchType: "symbol", Path: "graph.go", Snippet: "DependencyGraph"}},
			}},
		})
	}
	return snapshot
}

func m18Bundle(repo, commit string) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "julianpoy", Repo: repo, Commit: commit}
	fileContent, testContent, licenseContent := "type DependencyGraph struct{}", "func TestDependencyGraph(){}", "MIT License"
	return domain.EvidenceBundle{
		SchemaVersion: "1", Repository: ref, Complete: commit != "" && commit != "HEAD", TreeComplete: true,
		Items: []domain.CollectedEvidence{
			{Kind: domain.EvidenceFile, Ref: ref, Path: "graph.go", Content: fileContent, ContentHash: contentHash(fileContent)},
			{Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: testContent, ContentHash: contentHash(testContent)},
			{Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: licenseContent, ContentHash: contentHash(licenseContent)},
		},
	}
}

func m18ReplayRun(t *testing.T, snapshot domain.DiscoverySnapshot, bundles []domain.EvidenceBundle) domain.IdeaRunResult {
	t.Helper()
	request := maxCandidatesReplayRequest(snapshot, len(snapshot.Candidates))
	result, err := RunIdeaPipeline(context.Background(), request, IdeaRunDependencies{ReplaySnapshot: &snapshot, ReplayBundles: bundles})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func m18RequireCandidateCommit(t *testing.T, result domain.IdeaRunResult, repository, commit string, inconclusive bool) domain.SoftwareAlternativeReport {
	t.Helper()
	for _, candidate := range result.Discovery.Candidates {
		if strings.EqualFold(candidate.Repository.Owner+"/"+candidate.Repository.Repo, repository) && candidate.Repository.CommitID != commit {
			t.Fatalf("discovery commit = %q, want %q", candidate.Repository.CommitID, commit)
		}
	}
	foundAssessment := false
	for _, assessment := range result.Assessment.Assessments {
		if strings.EqualFold(assessment.Repository.Owner+"/"+assessment.Repository.Repo, repository) {
			foundAssessment = true
			if assessment.Repository.Commit != commit {
				t.Fatalf("assessment commit = %q, want %q", assessment.Repository.Commit, commit)
			}
		}
	}
	if !foundAssessment {
		t.Fatalf("assessment missing for %s", repository)
	}
	report, err := RankSoftwareAlternatives(domain.SoftwareAlternativeTarget{Name: "Graph Tool"}, result.Discovery.Idea, result.Discovery, result.Assessment.Assessments, maxCandidatesPolicy())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, alternative := range append(report.Alternatives, report.Excluded...) {
		if !strings.EqualFold(alternative.Repository.Owner+"/"+alternative.Repository.Repo, repository) {
			continue
		}
		found = true
		if alternative.Repository.Commit != commit {
			t.Fatalf("ranking commit = %q, want %q", alternative.Repository.Commit, commit)
		}
		for _, reason := range alternative.Reasons {
			if strings.Contains(strings.ToLower(reason), "commit divergente") && !inconclusive {
				t.Fatalf("resolved candidate remained divergent: %+v", alternative.Reasons)
			}
		}
		if inconclusive && alternative.Fit != domain.AlternativeInconclusive {
			t.Fatalf("opaque resolution fit = %s, want INCONCLUSIVE", alternative.Fit)
		}
	}
	if !found {
		t.Fatalf("ranking candidate missing for %s", repository)
	}
	return report
}

func TestRunIdeaPipelineM18ResolvedCommitPropagation(t *testing.T) {
	cases := []struct {
		id  string
		run func(*testing.T)
	}{
		{"M18a-live-resolved-commit", testM18LiveResolvedCommit},
		{"M18b-replay-refines-head", testM18ReplayRefinesHead},
		{"M18c-immutable-divergence-preserved", testM18ImmutableDivergencePreserved},
		{"M18d-opaque-resolution-preserved", testM18OpaqueResolutionPreserved},
		{"M18e-reorder-deterministic", testM18ReorderDeterministic},
	}
	killed := 0
	for _, mutation := range cases {
		mutation := mutation
		if t.Run(mutation.id, mutation.run) {
			killed++
		}
	}
	if killed != len(cases) {
		t.Fatalf("M18 mutation score = %d/%d, want 5/5", killed, len(cases))
	}
}

func testM18LiveResolvedCommit(t *testing.T) {
	candidate := m18Candidate("RecipeSage", "HEAD")
	request := domain.IdeaRunRequest{
		Idea: maxCandidatesIdea(), Policy: maxCandidatesPolicy(), Mode: domain.IdeaRunLive,
		Strategy: "literal", OutputDirectory: "unused-output", MaxPerQuery: 5, MaxCandidates: 1,
	}
	result, err := RunIdeaPipeline(context.Background(), request, IdeaRunDependencies{
		Discovery: maxCandidatesDiscoveryProvider{candidates: []domain.RepositoryCandidate{candidate}},
		Content:   m18ContentProvider{resolved: m18WitnessCommit},
	})
	if err != nil {
		t.Fatal(err)
	}
	m18RequireCandidateCommit(t, result, "julianpoy/recipesage", m18WitnessCommit, false)
}

func testM18ReplayRefinesHead(t *testing.T) {
	snapshot := m18Snapshot(m18Candidate("RecipeSage", "HEAD"))
	result := m18ReplayRun(t, snapshot, []domain.EvidenceBundle{m18Bundle("RecipeSage", m18WitnessCommit)})
	m18RequireCandidateCommit(t, result, "julianpoy/recipesage", m18WitnessCommit, false)
}

func testM18ImmutableDivergencePreserved(t *testing.T) {
	const snapshotCommit = "1111111111111111111111111111111111111111"
	const bundleCommit = "2222222222222222222222222222222222222222"
	snapshot := m18Snapshot(m18Candidate("RecipeSage", snapshotCommit))
	result := m18ReplayRun(t, snapshot, []domain.EvidenceBundle{m18Bundle("RecipeSage", bundleCommit)})
	m18RequireCandidateCommit(t, result, "julianpoy/recipesage", snapshotCommit, true)
}

func testM18OpaqueResolutionPreserved(t *testing.T) {
	variants := []struct {
		name    string
		bundles []domain.EvidenceBundle
	}{
		{"empty", []domain.EvidenceBundle{m18Bundle("RecipeSage", "")}},
		{"head", []domain.EvidenceBundle{m18Bundle("RecipeSage", "HEAD")}},
		{"conflicting", []domain.EvidenceBundle{m18Bundle("RecipeSage", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), m18Bundle("RecipeSage", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")}},
	}
	for _, variant := range variants {
		variant := variant
		t.Run(variant.name, func(t *testing.T) {
			snapshot := m18Snapshot(m18Candidate("RecipeSage", "HEAD"))
			result := m18ReplayRun(t, snapshot, variant.bundles)
			m18RequireCandidateCommit(t, result, "julianpoy/recipesage", "HEAD", true)
		})
	}
}

func testM18ReorderDeterministic(t *testing.T) {
	alpha := m18Candidate("Alpha", "HEAD")
	bravo := m18Candidate("Bravo", "HEAD")
	alphaBundle := m18Bundle("Alpha", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	bravoBundle := m18Bundle("Bravo", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	first := m18ReplayRun(t, m18Snapshot(bravo, alpha), []domain.EvidenceBundle{alphaBundle, bravoBundle})
	second := m18ReplayRun(t, m18Snapshot(alpha, bravo), []domain.EvidenceBundle{bravoBundle, alphaBundle})
	firstRank, err := RankSoftwareAlternatives(domain.SoftwareAlternativeTarget{Name: "Graph Tool"}, first.Discovery.Idea, first.Discovery, first.Assessment.Assessments, maxCandidatesPolicy())
	if err != nil {
		t.Fatal(err)
	}
	secondRank, err := RankSoftwareAlternatives(domain.SoftwareAlternativeTarget{Name: "Graph Tool"}, second.Discovery.Idea, second.Discovery, second.Assessment.Assessments, maxCandidatesPolicy())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(struct {
		Discovery  domain.DiscoverySnapshot         `json:"discovery"`
		Assessment domain.AssessmentReport          `json:"assessment"`
		Ranking    domain.SoftwareAlternativeReport `json:"ranking"`
	}{first.Discovery, first.Assessment, firstRank})
	b, _ := json.Marshal(struct {
		Discovery  domain.DiscoverySnapshot         `json:"discovery"`
		Assessment domain.AssessmentReport          `json:"assessment"`
		Ranking    domain.SoftwareAlternativeReport `json:"ranking"`
	}{second.Discovery, second.Assessment, secondRank})
	if string(a) != string(b) {
		t.Fatalf("M18 reorder changed canonical output:\n%s\n%s", a, b)
	}
}
