/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type countingDiscovery struct{ calls int }

func (c *countingDiscovery) Discover(context.Context, domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	c.calls++
	return domain.ProviderDiscovery{}, nil
}

type fixtureContent struct{}

func (fixtureContent) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, domain.ErrManifestNotFound
}
func (fixtureContent) FetchTree(_ context.Context, _ domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "LICENSE", Type: "file"}, {Path: "graph.go", Type: "file"}, {Path: "graph_test.go", Type: "file"}}, nil
}
func (fixtureContent) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	values := map[string]string{"README.md": "dependency graph analysis", "LICENSE": "MIT License", "graph.go": "type DependencyGraph struct{}", "graph_test.go": "func TestDependencyGraph(){}"}
	files := []domain.ExtractedFile{}
	for _, p := range paths {
		if v, ok := values[p]; ok {
			files = append(files, domain.ExtractedFile{Path: p, Content: v, Size: int64(len(v))})
		}
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}, Candidate: c, Files: files}, nil
}
func (fixtureContent) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	if c.CommitID == "" || c.CommitID == "HEAD" {
		c.CommitID = "abc123"
	}
	return c, nil
}

func recommendReplayFixture(t *testing.T) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "source")
	ideaPath := filepath.Join(root, "idea.json")
	policyPath := filepath.Join(root, "policy.json")
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	candidate := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc123", Description: "dependency graph analysis"}
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal", Candidates: []domain.DiscoveredCandidate{{Repository: candidate, ExpansionIDs: []string{"capability-dependency-graph"}, Queries: []string{"dependency graph"}, Origins: []domain.CandidateOrigin{{ExpansionID: "capability-dependency-graph", PlannedQuery: "dependency graph", EffectiveQuery: "dependency graph", Evidence: []domain.MatchEvidence{{MatchType: "content", Path: "graph.go", Snippet: "DependencyGraph"}}}}}}}
	fc, tc, lc := "type DependencyGraph struct{}", "func TestDependencyGraph(){}", "MIT License"
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceFile, Ref: ref, Path: "graph.go", Content: fc, ContentHash: hashBytes([]byte(fc))}, {Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: tc, ContentHash: hashBytes([]byte(tc))}, {Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: lc, ContentHash: hashBytes([]byte(lc))}}}
	for p, v := range map[string]any{ideaPath: idea, policyPath: policy, filepath.Join(source, "discovery-snapshot.json"): snapshot, filepath.Join(source, "evidence", "bundles.json"): []domain.EvidenceBundle{bundle}} {
		if err := writeJSONAtomic(p, v); err != nil {
			t.Fatal(err)
		}
	}
	return ideaPath, policyPath, source
}
func TestIdeaRecommendReplayCreatesDeterministicDossierWithoutNetwork(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	spy := &countingDiscovery{}
	run := func(out string) domain.ReuseDecisionReport {
		var stdout bytes.Buffer
		err := RunIdeaRecommendCommand([]string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", out}, &stdout, nil, spy, nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range []string{"run.json", "idea.json", "expansion-plan.json", "discovery-snapshot.json", "screening-report.json", "assessment-report.json", "reuse-decision.json", "decision.md", filepath.Join("evidence", "plans.json"), filepath.Join("evidence", "bundles.json")} {
			if _, err := os.Stat(filepath.Join(out, p)); err != nil {
				t.Fatalf("missing %s: %v", p, err)
			}
		}
		var decision domain.ReuseDecisionReport
		if err := readJSONFile(filepath.Join(out, "reuse-decision.json"), &decision); err != nil {
			t.Fatal(err)
		}
		return decision
	}
	a := run(filepath.Join(t.TempDir(), "a"))
	b := run(filepath.Join(t.TempDir(), "b"))
	if spy.calls != 0 {
		t.Fatalf("replay called provider %d times", spy.calls)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("replay decision changed")
	}
}
func TestIdeaRecommendValidatesOutputBeforeProvider(t *testing.T) {
	spy := &countingDiscovery{}
	err := RunIdeaRecommendCommand([]string{"--mode", "live"}, &bytes.Buffer{}, nil, spy, nil)
	if err == nil || spy.calls != 0 {
		t.Fatalf("err=%v calls=%d", err, spy.calls)
	}
}
func TestAtomicArtifactsRejectConfiguredSecret(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "super-secret-token")
	err := writeAtomic(filepath.Join(t.TempDir(), "x.json"), []byte(`{"value":"super-secret-token"}`), 0600)
	if err == nil || !strings.Contains(err.Error(), "segredo") {
		t.Fatalf("got %v", err)
	}
}
func TestManifestContainsHashChain(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	out := filepath.Join(t.TempDir(), "out")
	if err := RunIdeaRecommendCommand([]string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", out}, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(out, "run.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m domain.IdeaRunManifest
	if err = json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if len(m.Stages) != 6 || m.DecisionPath != "reuse-decision.json" {
		t.Fatalf("manifest %#v", m)
	}
	for _, s := range m.Stages {
		if len(s.Outputs) != 1 || s.Outputs[0].Hash == "" {
			t.Fatalf("stage without hash: %#v", s)
		}
	}
}

func TestResumeDoesNotReuseFailedRunWithoutArtifacts(t *testing.T) {
	root := t.TempDir()
	idea := domain.IdeaSpec{Problem: "x", Capabilities: []string{"dependency graph"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	ideaPath, policyPath, out := filepath.Join(root, "idea.json"), filepath.Join(root, "policy.json"), filepath.Join(root, "out")
	_ = writeJSONAtomic(ideaPath, idea)
	_ = writeJSONAtomic(policyPath, policy)
	_ = os.MkdirAll(out, 0755)
	_ = writeJSONAtomic(filepath.Join(out, "run.json"), domain.IdeaRunManifest{SchemaVersion: "1", Status: domain.IdeaRunFailed, InputHashes: map[string]string{"idea": hashJSON(idea)}})
	provider := &fixtureDiscovery{}
	if err := RunIdeaRecommendCommand([]string{"--idea", ideaPath, "--policy", policyPath, "--mode", "live", "--strategy", "literal", "--output", out, "--resume"}, &bytes.Buffer{}, nil, provider, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	if provider.calls == 0 {
		t.Fatal("failed checkpoint was reused")
	}
}

func TestIdeaRecommendLocalUsesNonPublicIdentityAndDirtyContentHash(t *testing.T) {
	root := t.TempDir()
	worktree := filepath.Join(root, "repo")
	if err := os.MkdirAll(worktree, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worktree, "dirty.go"), []byte("package dirty"), 0600); err != nil {
		t.Fatal(err)
	}
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	ideaPath, policyPath := filepath.Join(root, "idea.json"), filepath.Join(root, "policy.json")
	if err := writeJSONAtomic(ideaPath, idea); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(policyPath, policy); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "out")
	if err := RunIdeaRecommendCommand([]string{"--idea", ideaPath, "--policy", policyPath, "--mode", "local", "--strategy", "literal", "--local", worktree, "--output", out}, &bytes.Buffer{}, nil, nil, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	var snapshot domain.DiscoverySnapshot
	if err := readJSONFile(filepath.Join(out, "discovery-snapshot.json"), &snapshot); err != nil {
		t.Fatal(err)
	}
	candidate := snapshot.Candidates[0].Repository
	if candidate.Owner != "local" || !strings.HasPrefix(candidate.URL, "file://") || candidate.CommitID == "" {
		t.Fatalf("invented identity: %#v", candidate)
	}
	before := candidate.CommitID
	if err := os.WriteFile(filepath.Join(worktree, "dirty.go"), []byte("package changed"), 0600); err != nil {
		t.Fatal(err)
	}
	candidates, err := localCandidates(idea, []string{worktree})
	if err != nil {
		t.Fatal(err)
	}
	if candidates[0].CommitID == before {
		t.Fatal("unpublished change did not alter content hash")
	}
}

func TestIdeaRecommendLiveUsesInjectedProviders(t *testing.T) {
	root := t.TempDir()
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	ideaPath, policyPath := filepath.Join(root, "idea.json"), filepath.Join(root, "policy.json")
	_ = writeJSONAtomic(ideaPath, idea)
	_ = writeJSONAtomic(policyPath, policy)
	provider := &fixtureDiscovery{}
	err := RunIdeaRecommendCommand([]string{"--idea", ideaPath, "--policy", policyPath, "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--output", filepath.Join(root, "out")}, &bytes.Buffer{}, nil, provider, fixtureContent{})
	if err != nil {
		t.Fatal(err)
	}
	if provider.calls == 0 {
		t.Fatal("live did not call discovery")
	}
}

type fixtureDiscovery struct{ calls int }

func (f *fixtureDiscovery) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	f.calls++
	candidate := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc123", Description: "dependency graph analysis"}
	return domain.ProviderDiscovery{Provider: "fixture", Attempts: []domain.ProviderAttempt{{Query: request.Query, ResultCount: 1}}, Candidates: []domain.RepositoryCandidate{candidate}, Matches: []domain.ProviderCandidateMatch{{Repository: "org/repo", EffectiveQuery: request.Query, Evidence: []domain.MatchEvidence{{MatchType: "content", Path: "graph.go", Snippet: "DependencyGraph"}}}}}, nil
}
