/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type countingDiscovery struct{ calls int }

func (c *countingDiscovery) Discover(context.Context, domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	c.calls++
	return domain.ProviderDiscovery{}, nil
}

func TestIdeaRecommendM20RendersUnknownLicenseWarningDeterministically(t *testing.T) {
	type result struct {
		json     []byte
		markdown []byte
		report   domain.SoftwareAlternativeReport
	}
	run := func(reordered bool) result {
		ideaPath, policyPath, source := recommendM20ReplayFixture(t, reordered)
		out := filepath.Join(t.TempDir(), "out")
		args := []string{
			"--idea", ideaPath,
			"--policy", policyPath,
			"--mode", "replay",
			"--strategy", "literal",
			"--replay", source,
			"--output", out,
			"--target-name", "Mealie",
			"--target-repo", "mealie-recipes/mealie",
		}
		if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
		jsonBytes, err := os.ReadFile(filepath.Join(out, "alternatives.json"))
		if err != nil {
			t.Fatal(err)
		}
		markdown, err := os.ReadFile(filepath.Join(out, "alternatives.md"))
		if err != nil {
			t.Fatal(err)
		}
		var report domain.SoftwareAlternativeReport
		if err := json.Unmarshal(jsonBytes, &report); err != nil {
			t.Fatal(err)
		}
		return result{json: jsonBytes, markdown: markdown, report: report}
	}

	first := run(false)
	second := run(true)
	if string(first.json) != string(second.json) {
		t.Fatalf("M20 input reorder changed alternatives JSON:\n%s\n%s", first.json, second.json)
	}
	if string(first.markdown) != string(second.markdown) {
		t.Fatalf("M20 input reorder changed alternatives Markdown:\n%s\n%s", first.markdown, second.markdown)
	}

	const warning = "licença não reconhecida; valide os termos antes da adoção"
	found := false
	for _, alternative := range first.report.Alternatives {
		if strings.EqualFold(alternative.Repository.Owner+"/"+alternative.Repository.Repo, "julianpoy/RecipeSage") {
			found = true
			if alternative.Fit != domain.AlternativeReference {
				t.Fatalf("RecipeSage fit = %s, want %s", alternative.Fit, domain.AlternativeReference)
			}
			if alternative.License.Status != domain.EvidenceInconclusive || alternative.License.Identifier != "" {
				t.Fatalf("RecipeSage license changed: %+v", alternative.License)
			}
			if !containsExactString(alternative.Reasons, warning) {
				t.Fatalf("RecipeSage warning missing from JSON reasons: %+v", alternative.Reasons)
			}
		}
	}
	if !found {
		t.Fatalf("RecipeSage did not migrate to alternatives: %+v", first.report)
	}
	if first.report.TopChoice != nil {
		t.Fatalf("unexpected top choice: %+v", first.report.TopChoice)
	}
	for _, line := range strings.Split(string(first.markdown), "\n") {
		if strings.Contains(strings.ToLower(line), "julianpoy/recipesage@") {
			if !strings.Contains(line, warning) {
				t.Fatalf("RecipeSage Markdown row hid warning: %s", line)
			}
			return
		}
	}
	t.Fatalf("RecipeSage row absent from Markdown: %s", first.markdown)
}

func recommendM20ReplayFixture(t *testing.T, reordered bool) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "source")
	ideaPath := filepath.Join(root, "idea.json")
	policyPath := filepath.Join(root, "policy.json")
	idea := domain.IdeaSpec{Problem: "organizar receitas culinárias para uma família", Capabilities: []string{"organização de receitas", "importação de receitas"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 4, ExtendMinCoverage: .66, IntegrationMinValue: .33}
	consulted := time.Date(2026, 9, 1, 12, 13, 12, 0, time.UTC)

	norish := domain.RepositoryCandidate{Owner: "norish-recipes", Repo: "norish", CommitID: "e3691fe08addabd42feb68e8cc942fcd6468634b", Stars: 10, Description: "organizar receitas culinárias para uma família; Imports recipes from Tandoor.", ConsultadoEm: consulted, UpdatedAt: consulted.Add(-time.Hour)}
	recipeSage := domain.RepositoryCandidate{Owner: "julianpoy", Repo: "RecipeSage", CommitID: "120aa5c8816f9a17e043d48c653a6ce200937ea9", Stars: 948, Description: "A Collaborative Recipe Keeper, Meal Planner, and Shopping List Organizer in PWA form. Imports recipes from Tandoor.", ConsultadoEm: consulted, UpdatedAt: consulted.Add(-time.Hour)}
	candidates := []domain.RepositoryCandidate{norish, recipeSage}
	bundles := []domain.EvidenceBundle{
		m20ReplayBundle(t, norish, true, true),
		m20ReplayBundle(t, recipeSage, false, false),
	}
	if reordered {
		candidates[0], candidates[1] = candidates[1], candidates[0]
		bundles[0], bundles[1] = bundles[1], bundles[0]
	}
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal"}
	for _, candidate := range candidates {
		snapshot.Candidates = append(snapshot.Candidates, domain.DiscoveredCandidate{Repository: candidate})
	}
	for path, value := range map[string]any{
		ideaPath:   idea,
		policyPath: policy,
		filepath.Join(source, "discovery-snapshot.json"):  snapshot,
		filepath.Join(source, "evidence", "bundles.json"): bundles,
	} {
		if err := writeJSONAtomic(path, value); err != nil {
			t.Fatal(err)
		}
	}
	return ideaPath, policyPath, source
}

func m20ReplayBundle(t *testing.T, candidate domain.RepositoryCandidate, includeStructuralProof, includeLicense bool) domain.EvidenceBundle {
	t.Helper()
	ref := domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}
	items := []domain.CollectedEvidence{}
	add := func(kind domain.EvidenceKind, path, content string) {
		items = append(items, domain.CollectedEvidence{Kind: kind, Ref: ref, Path: path, Content: content, ContentHash: hashBytes([]byte(content))})
	}
	add(domain.EvidenceDocumentation, "README.md", candidate.Description)
	if includeStructuralProof {
		add(domain.EvidenceFile, "recipes.py", "class RecipeController:")
		add(domain.EvidenceTest, "tandoor_test.go", "func TestTandoorImport(t *testing.T) {}")
	} else {
		add(domain.EvidenceDocumentation, "docs/recipes.md", "Recipe Keeper supports recipe organization and imports recipes for family meal planning.")
	}
	if includeLicense {
		add(domain.EvidenceLicense, "LICENSE", "MIT License\nPermission is hereby granted, free of charge")
	}
	metadata, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	add(domain.EvidenceRepository, "@repository-metadata", string(metadata))
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: items}
}

func containsExactString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
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
	if len(m.Stages) != 7 || m.DecisionPath != "reuse-decision.json" {
		t.Fatalf("manifest %#v", m)
	}
	assertRecommendManifestStagesHashed(t, m)
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

type checkpointContentSpy struct {
	manifestCalls int
	treeCalls     int
	extractCalls  int
	metadataCalls int
}

func (s *checkpointContentSpy) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	s.manifestCalls++
	return (fixtureContent{}).FetchManifest(ctx, candidate)
}

func (s *checkpointContentSpy) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	s.treeCalls++
	return (fixtureContent{}).FetchTree(ctx, candidate)
}

func (s *checkpointContentSpy) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	s.extractCalls++
	return (fixtureContent{}).ExtractRepoFiles(ctx, candidate, paths)
}

func (s *checkpointContentSpy) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	s.metadataCalls++
	return (fixtureContent{}).FetchRepoMetadata(ctx, candidate)
}

type checkpointFixture struct {
	ideaPath   string
	policyPath string
	outputPath string
	provider   *fixtureDiscovery
	content    *checkpointContentSpy
}

func newCheckpointFixture(t *testing.T) *checkpointFixture {
	t.Helper()
	root := t.TempDir()
	fixture := &checkpointFixture{
		ideaPath:   filepath.Join(root, "idea.json"),
		policyPath: filepath.Join(root, "policy.json"),
		outputPath: filepath.Join(root, "out"),
		provider:   &fixtureDiscovery{},
		content:    &checkpointContentSpy{},
	}
	idea := domain.IdeaSpec{Problem: "analisar dependencias", Capabilities: []string{"dependency graph"}}
	policy := checkpointPolicy("MIT")
	if err := writeJSONAtomic(fixture.ideaPath, idea); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(fixture.policyPath, policy); err != nil {
		t.Fatal(err)
	}
	fixture.run(t, "go", false)
	return fixture
}

func checkpointPolicy(license string) domain.ReuseDecisionPolicy {
	return domain.ReuseDecisionPolicy{
		SchemaVersion:       "1",
		AllowedLicenses:     []string{license},
		MaxComponents:       2,
		ExtendMinCoverage:   .6,
		IntegrationMinValue: .3,
	}
}

func (f *checkpointFixture) run(t *testing.T, language string, resume bool) {
	t.Helper()
	args := []string{
		"--idea", f.ideaPath,
		"--policy", f.policyPath,
		"--mode", "live",
		"--strategy", "literal",
		"--provider", "sourcegraph",
		"--language", language,
		"--max-per-query", "3",
		"--max-candidates", "2",
		"--output", f.outputPath,
	}
	if resume {
		args = append(args, "--resume")
	}
	if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, f.provider, f.content); err != nil {
		t.Fatal(err)
	}
}

func TestIdeaRecommendResumeRejectsTamperedDiscoverySnapshot(t *testing.T) {
	fixture := newCheckpointFixture(t)
	providerCalls := fixture.provider.calls

	path := filepath.Join(fixture.outputPath, "discovery-snapshot.json")
	var snapshot domain.DiscoverySnapshot
	if err := readJSONFile(path, &snapshot); err != nil {
		t.Fatal(err)
	}
	snapshot.Strategy = "tampered-after-run"
	if err := writeJSONAtomic(path, snapshot); err != nil {
		t.Fatal(err)
	}

	fixture.run(t, "go", true)
	if fixture.provider.calls <= providerCalls {
		t.Fatalf("tampered discovery checkpoint was reused: provider calls stayed at %d", providerCalls)
	}
}

func TestIdeaRecommendResumeRejectsTamperedEvidenceBundle(t *testing.T) {
	fixture := newCheckpointFixture(t)
	extractCalls := fixture.content.extractCalls

	path := filepath.Join(fixture.outputPath, "evidence", "bundles.json")
	var bundles []domain.EvidenceBundle
	if err := readJSONFile(path, &bundles); err != nil {
		t.Fatal(err)
	}
	if len(bundles) == 0 || len(bundles[0].Items) == 0 {
		t.Fatal("fixture produced no evidence to tamper")
	}
	bundles[0].Items[0].Content = "tampered after the recorded checkpoint hash"
	if err := writeJSONAtomic(path, bundles); err != nil {
		t.Fatal(err)
	}

	fixture.run(t, "go", true)
	if fixture.content.extractCalls <= extractCalls {
		t.Fatalf("tampered evidence checkpoint was reused: extraction calls stayed at %d", extractCalls)
	}
}

func TestIdeaRecommendResumeRejectsChangedCausalInput(t *testing.T) {
	fixture := newCheckpointFixture(t)
	providerCalls := fixture.provider.calls

	fixture.run(t, "rust", true)
	if fixture.provider.calls <= providerCalls {
		t.Fatalf("checkpoint created for language go was reused for rust: provider calls stayed at %d", providerCalls)
	}
}

func TestIdeaRecommendResumeReusesEvidenceForPolicyOnlyChangeAndRecalculatesDecision(t *testing.T) {
	fixture := newCheckpointFixture(t)
	providerCalls := fixture.provider.calls
	extractCalls := fixture.content.extractCalls
	before := readCheckpointDecision(t, fixture.outputPath)

	if err := writeJSONAtomic(fixture.policyPath, checkpointPolicy("Apache-2.0")); err != nil {
		t.Fatal(err)
	}
	fixture.run(t, "go", true)
	after := readCheckpointDecision(t, fixture.outputPath)

	if fixture.provider.calls != providerCalls {
		t.Fatalf("policy-only change repeated discovery: calls %d -> %d", providerCalls, fixture.provider.calls)
	}
	if fixture.content.extractCalls != extractCalls {
		t.Fatalf("policy-only change repeated evidence collection: calls %d -> %d", extractCalls, fixture.content.extractCalls)
	}
	if after.PolicyHash == before.PolicyHash {
		t.Fatal("policy-only resume kept the previous decision policy hash")
	}
	if len(after.Policy.AllowedLicenses) != 1 || after.Policy.AllowedLicenses[0] != "Apache-2.0" {
		t.Fatalf("policy-only resume kept the previous policy in the recalculated decision: %#v", after.Policy)
	}
}

func readCheckpointDecision(t *testing.T, outputPath string) domain.ReuseDecisionReport {
	t.Helper()
	var decision domain.ReuseDecisionReport
	if err := readJSONFile(filepath.Join(outputPath, "reuse-decision.json"), &decision); err != nil {
		t.Fatal(err)
	}
	return decision
}

func TestCheckpointFixtureWritesHashedRun(t *testing.T) {
	fixture := newCheckpointFixture(t)
	var manifest domain.IdeaRunManifest
	if err := readJSONFile(filepath.Join(fixture.outputPath, "run.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Status == domain.IdeaRunFailed || len(manifest.Stages) == 0 {
		t.Fatalf("checkpoint precondition: status=%s stages=%d", manifest.Status, len(manifest.Stages))
	}
	assertRecommendManifestStagesHashed(t, manifest)
	if _, err := os.Stat(filepath.Join(fixture.outputPath, "evidence", "bundles.json")); err != nil {
		t.Fatal(err)
	}
}

func assertRecommendManifestStagesHashed(t *testing.T, manifest domain.IdeaRunManifest) {
	t.Helper()
	foundAlternatives := false
	for _, stage := range manifest.Stages {
		if stage.Status != domain.IdeaStageComplete {
			t.Fatalf("incomplete stage: %#v", stage)
		}
		if stage.Name != "alternatives" {
			if len(stage.Outputs) != 1 || stage.Outputs[0].Hash == "" {
				t.Fatalf("stage without exactly one hashed output: %#v", stage)
			}
			continue
		}
		foundAlternatives = true
		if len(stage.Outputs) != 2 {
			t.Fatalf("alternatives stage outputs = %#v", stage.Outputs)
		}
		paths := map[string]bool{}
		for _, artifact := range stage.Outputs {
			if artifact.Hash == "" {
				t.Fatalf("unhashed alternatives artifact: %#v", artifact)
			}
			paths[filepath.ToSlash(artifact.Path)] = true
		}
		if !paths["alternatives.json"] || !paths["alternatives.md"] {
			t.Fatalf("alternatives stage lost sealed artifacts: %#v", stage.Outputs)
		}
	}
	if !foundAlternatives {
		t.Fatal("alternatives stage missing")
	}
}

func TestIdeaRecommendRunIDIgnoresOutputAndResumeOperationalFlags(t *testing.T) {
	ideaPath, policyPath, sourcePath := recommendReplayFixture(t)
	root := t.TempDir()
	outputA := filepath.Join(root, "output-a")
	outputB := filepath.Join(root, "output-b")
	run := func(outputPath string, resume bool) string {
		t.Helper()
		args := []string{
			"--idea", ideaPath,
			"--policy", policyPath,
			"--mode", "replay",
			"--strategy", "literal",
			"--replay", sourcePath,
			"--output", outputPath,
		}
		if resume {
			args = append(args, "--resume")
		}
		if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
		var manifest domain.IdeaRunManifest
		if err := readJSONFile(filepath.Join(outputPath, "run.json"), &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.RunID == "" {
			t.Fatal("run manifest has an empty run id")
		}
		return manifest.RunID
	}

	originalID := run(outputA, false)
	otherOutputID := run(outputB, false)
	if otherOutputID != originalID {
		t.Fatalf("output directory changed experiment identity: %q != %q", otherOutputID, originalID)
	}
	resumedID := run(outputA, true)
	if resumedID != originalID {
		t.Fatalf("resume changed experiment identity: %q != %q", resumedID, originalID)
	}
}

func TestIdeaRecommendTargetArtifactsRejectSealedRepositorySyntaxMutations(t *testing.T) {
	cases := []struct {
		mutation string
		target   string
	}{
		{"M06b-target-http-url", "https://github.com/ORG/REPO"},
		{"M06c-target-ssh", "git@github.com:org/repo"},
		{"M06d-target-dotgit", "org/repo.git"},
	}
	killed := 0
	for _, testCase := range cases {
		testCase := testCase
		if t.Run(testCase.mutation, func(t *testing.T) {
			idea, policy, source := recommendReplayFixture(t)
			out := filepath.Join(t.TempDir(), "out")
			args := []string{
				"--idea", idea,
				"--policy", policy,
				"--mode", "replay",
				"--strategy", "literal",
				"--replay", source,
				"--output", out,
				"--target-name", "Known Graph",
				"--target-repo", testCase.target,
			}
			if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
				t.Fatal(err)
			}

			for _, artifact := range []string{"alternatives.json", "alternatives.md"} {
				if _, err := os.Stat(filepath.Join(out, artifact)); err != nil {
					t.Fatalf("missing %s: %v", artifact, err)
				}
			}
			var report domain.SoftwareAlternativeReport
			if err := readJSONFile(filepath.Join(out, "alternatives.json"), &report); err != nil {
				t.Fatal(err)
			}
			if report.Target.Name != "Known Graph" || report.Target.Repository == nil || report.Target.Repository.Owner != "org" || report.Target.Repository.Repo != "repo" {
				t.Fatalf("target was not canonicalized: %+v", report.Target)
			}
			if len(report.Alternatives) != 0 || len(report.Excluded) != 1 || report.Excluded[0].Fit != domain.AlternativeNotAlternative || report.TopChoice != nil {
				t.Fatalf("target became its own alternative: %+v", report)
			}

			markdown, err := os.ReadFile(filepath.Join(out, "alternatives.md"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(markdown), "Known Graph") || !strings.Contains(string(markdown), "org/repo") {
				t.Fatalf("markdown lost target identity: %s", markdown)
			}

			var manifest domain.IdeaRunManifest
			if err := readJSONFile(filepath.Join(out, "run.json"), &manifest); err != nil {
				t.Fatal(err)
			}
			var alternativesStage *domain.IdeaRunStage
			for index := range manifest.Stages {
				if manifest.Stages[index].Name == "alternatives" {
					alternativesStage = &manifest.Stages[index]
					break
				}
			}
			if alternativesStage == nil || alternativesStage.Status != domain.IdeaStageComplete {
				t.Fatalf("alternatives stage missing/incomplete: %+v", manifest.Stages)
			}
			paths := map[string]bool{}
			for _, artifact := range alternativesStage.Outputs {
				if artifact.Hash == "" {
					t.Fatalf("unhashed alternatives artifact: %+v", artifact)
				}
				paths[filepath.ToSlash(artifact.Path)] = true
			}
			if !paths["alternatives.json"] || !paths["alternatives.md"] {
				t.Fatalf("alternatives stage outputs = %+v", alternativesStage.Outputs)
			}
		}) {
			killed++
		}
	}
	if killed != len(cases) {
		t.Fatalf("CLI mutation score = %d/%d, want 3/3", killed, len(cases))
	}
}

func TestIdeaRecommendLocalExcludesMatchingGitRemoteWithoutRemappingOtherWorktrees(t *testing.T) {
	root := t.TempDir()
	matching := newRecommendLocalGitFixture(t, filepath.Join(root, "unrelated-directory-a"), "https://github.com/Dikluwe/Bateia.git")
	nonMatching := newRecommendLocalGitFixture(t, filepath.Join(root, "unrelated-directory-b"), "git@github.com:other/tool.git")
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	ideaPath, policyPath, outputPath := filepath.Join(root, "idea.json"), filepath.Join(root, "policy.json"), filepath.Join(root, "out")
	if err := writeJSONAtomic(ideaPath, idea); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(policyPath, policy); err != nil {
		t.Fatal(err)
	}
	args := []string{
		"--idea", ideaPath,
		"--policy", policyPath,
		"--mode", "local",
		"--strategy", "literal",
		"--local", matching,
		"--local", nonMatching,
		"--output", outputPath,
		"--target-name", "Bateia",
		"--target-repo", "dikluwe/bateia",
	}
	if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	var report domain.SoftwareAlternativeReport
	if err := readJSONFile(filepath.Join(outputPath, "alternatives.json"), &report); err != nil {
		t.Fatal(err)
	}

	targetInAlternatives := 0
	targetInExcluded := 0
	nonTargetRecords := 0
	nonTargetTreatedAsTarget := false
	for _, alternative := range report.Alternatives {
		identity := strings.ToLower(alternative.Repository.Owner + "/" + alternative.Repository.Repo)
		if identity == "dikluwe/bateia" {
			targetInAlternatives++
		} else {
			nonTargetRecords++
		}
	}
	for _, alternative := range report.Excluded {
		identity := strings.ToLower(alternative.Repository.Owner + "/" + alternative.Repository.Repo)
		switch identity {
		case "dikluwe/bateia":
			targetInExcluded++
			if alternative.Fit != domain.AlternativeNotAlternative {
				t.Fatalf("matching worktree fit = %s, want NOT_ALTERNATIVE", alternative.Fit)
			}
		default:
			nonTargetRecords++
			if alternative.Fit == domain.AlternativeNotAlternative {
				nonTargetTreatedAsTarget = true
			}
		}
	}
	killed := 0
	if t.Run("M15-local-target-synthetic", func(t *testing.T) {
		if targetInAlternatives != 0 || targetInExcluded != 1 {
			t.Fatalf("matching remote placement: alternatives=%d excluded=%d report=%+v", targetInAlternatives, targetInExcluded, report)
		}
		if report.TopChoice != nil && strings.EqualFold(report.TopChoice.Owner+"/"+report.TopChoice.Repo, "dikluwe/bateia") {
			t.Fatalf("target became top choice: %+v", *report.TopChoice)
		}
	}) {
		killed++
	}
	if t.Run("M16-local-target-overreach", func(t *testing.T) {
		if nonTargetRecords != 1 || nonTargetTreatedAsTarget {
			t.Fatalf("non-matching remote was lost or remapped: count=%d treated_as_target=%v report=%+v", nonTargetRecords, nonTargetTreatedAsTarget, report)
		}
	}) {
		killed++
	}
	if killed != 2 {
		t.Fatalf("supplemental local mutation score = %d/2, want 2/2", killed)
	}
}

func newRecommendLocalGitFixture(t *testing.T, path, remote string) string {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"README.md":     "dependency graph analysis",
		"LICENSE":       "MIT License",
		"graph.go":      "package graph\n\ntype DependencyGraph struct{}\n",
		"graph_test.go": "package graph\n\nfunc TestDependencyGraph() {}\n",
	} {
		if err := os.WriteFile(filepath.Join(path, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	commands := [][]string{
		{"init", "-q"},
		{"config", "user.name", "Bateia Test"},
		{"config", "user.email", "bateia-test@example.invalid"},
		{"remote", "add", "origin", remote},
		{"add", "README.md", "LICENSE", "graph.go", "graph_test.go"},
		{"commit", "-q", "-m", "fixture"},
	}
	for _, arguments := range commands {
		command := exec.Command("git", append([]string{"-C", path}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", arguments, err, output)
		}
	}
	return path
}

type m21CLIRegionalDiscovery struct {
	mu       sync.Mutex
	requests []domain.DiscoveryRequest
}

func (provider *m21CLIRegionalDiscovery) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.mu.Lock()
	provider.requests = append(provider.requests, request)
	provider.mu.Unlock()
	candidate := domain.RepositoryCandidate{Owner: "regional", Repo: "recipes", CommitID: "abc123", Description: "dependency graph analysis"}
	return domain.ProviderDiscovery{
		Provider:   "m21-fixture",
		Attempts:   []domain.ProviderAttempt{{Query: request.Query, ResultCount: 1}},
		Candidates: []domain.RepositoryCandidate{candidate},
		Matches:    []domain.ProviderCandidateMatch{{Repository: "regional/recipes", EffectiveQuery: request.Query, Evidence: []domain.MatchEvidence{{MatchType: "content", Path: "README.md", Snippet: request.Query}}}},
	}, nil
}

func (provider *m21CLIRegionalDiscovery) regionalRequests() []domain.DiscoveryRequest {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	requests := []domain.DiscoveryRequest{}
	for _, request := range provider.requests {
		if request.SearchLocale != "" {
			requests = append(requests, request)
		}
	}
	return requests
}

func TestIdeaRecommendM21RegionalCLIReplayAndDeterminism(t *testing.T) {
	run := func(locales []string) (string, *m21CLIRegionalDiscovery) {
		root := t.TempDir()
		ideaPath := filepath.Join(root, "idea.json")
		policyPath := filepath.Join(root, "policy.json")
		idea := domain.IdeaSpec{Problem: "organizar receitas culinárias para uma família", Capabilities: []string{"dependency graph"}}
		policy := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
		if err := writeJSONAtomic(ideaPath, idea); err != nil {
			t.Fatal(err)
		}
		if err := writeJSONAtomic(policyPath, policy); err != nil {
			t.Fatal(err)
		}
		out := filepath.Join(root, "out")
		provider := &m21CLIRegionalDiscovery{}
		args := []string{"--idea", ideaPath, "--policy", policyPath, "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--output", out, "--target-name", "Mealie"}
		for _, locale := range locales {
			args = append(args, "--search-locale", locale)
		}
		if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, provider, fixtureContent{}); err != nil {
			t.Fatal(err)
		}
		return out, provider
	}

	firstOut, firstProvider := run([]string{"ja-JP", "pt-BR", "ja-JP"})
	secondOut, secondProvider := run([]string{"pt-BR", "ja-JP"})
	thirdOut, _ := run([]string{"en-US"})
	emptyOut, _ := run(nil)
	for name, provider := range map[string]*m21CLIRegionalDiscovery{"first": firstProvider, "second": secondProvider} {
		requests := provider.regionalRequests()
		if len(requests) != 2 {
			t.Fatalf("%s regional query budget count = %d, want 2: %+v", name, len(requests), requests)
		}
		wantQueries := map[string]string{"pt-BR": "gerenciador de receitas", "ja-JP": "レシピ管理"}
		seenRequests := map[string]bool{}
		for index, request := range requests {
			if request.Query != wantQueries[request.SearchLocale] || request.Language != "go" || !strings.Contains(request.RegionalProposalID, request.SearchLocale) {
				t.Fatalf("%s request[%d] lost locale/language/audit identity: %+v", name, index, request)
			}
			seenRequests[request.SearchLocale] = true
		}
		if !seenRequests["pt-BR"] || !seenRequests["ja-JP"] {
			t.Fatalf("%s requests lost a locale: %+v", name, requests)
		}
	}

	for _, artifact := range []string{"expansion-plan.json", "alternatives.json", "alternatives.md"} {
		left, err := os.ReadFile(filepath.Join(firstOut, artifact))
		if err != nil {
			t.Fatal(err)
		}
		right, err := os.ReadFile(filepath.Join(secondOut, artifact))
		if err != nil {
			t.Fatal(err)
		}
		if string(left) != string(right) {
			t.Fatalf("locale reorder/repeat changed %s:\n%s\n%s", artifact, left, right)
		}
	}
	var firstManifest, secondManifest, thirdManifest, emptyManifest domain.IdeaRunManifest
	if err := readJSONFile(filepath.Join(firstOut, "run.json"), &firstManifest); err != nil {
		t.Fatal(err)
	}
	if err := readJSONFile(filepath.Join(secondOut, "run.json"), &secondManifest); err != nil {
		t.Fatal(err)
	}
	if err := readJSONFile(filepath.Join(thirdOut, "run.json"), &thirdManifest); err != nil {
		t.Fatal(err)
	}
	if err := readJSONFile(filepath.Join(emptyOut, "run.json"), &emptyManifest); err != nil {
		t.Fatal(err)
	}
	if firstManifest.RunID == "" || firstManifest.RunID != secondManifest.RunID {
		t.Fatalf("equivalent locale sets changed execution identity: %q != %q", firstManifest.RunID, secondManifest.RunID)
	}
	if firstManifest.RunID == thirdManifest.RunID || firstManifest.InputHashes["evidence"] == thirdManifest.InputHashes["evidence"] {
		t.Fatalf("distinct canonical locale sets shared run/evidence identity: run=%q evidence=%q", firstManifest.RunID, firstManifest.InputHashes["evidence"])
	}
	if firstManifest.RunID == emptyManifest.RunID || firstManifest.InputHashes["evidence"] == emptyManifest.InputHashes["evidence"] {
		t.Fatalf("empty and non-empty locale sets shared run/evidence identity: non-empty=%+v empty=%+v", firstManifest.InputHashes, emptyManifest.InputHashes)
	}
	if _, ok := emptyManifest.InputHashes["search_locales"]; ok {
		t.Fatalf("zero-locale run materialized a search_locales input hash: %+v", emptyManifest.InputHashes)
	}

	var snapshot domain.DiscoverySnapshot
	if err := readJSONFile(filepath.Join(firstOut, "discovery-snapshot.json"), &snapshot); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(snapshot.SearchLocales, []string{"pt-BR", "ja-JP"}) || len(snapshot.Candidates) != 1 {
		t.Fatalf("snapshot regional projection = %+v", snapshot)
	}
	seen := map[string]bool{}
	for _, origin := range snapshot.Candidates[0].Origins {
		if origin.SearchLocale != "" {
			seen[origin.SearchLocale] = origin.PlannedQuery != "" && origin.EffectiveQuery != "" && origin.RegionalProposalID != ""
		}
	}
	if !seen["pt-BR"] || !seen["ja-JP"] {
		t.Fatalf("candidate lost regional origins: %+v", snapshot.Candidates[0].Origins)
	}

	var alternatives domain.SoftwareAlternativeReport
	if err := readJSONFile(filepath.Join(firstOut, "alternatives.json"), &alternatives); err != nil {
		t.Fatal(err)
	}
	values := append(append([]domain.SoftwareAlternative{}, alternatives.Alternatives...), alternatives.Excluded...)
	if len(values) != 1 || !reflect.DeepEqual(values[0].DiscoveryLocales, []string{"pt-BR", "ja-JP"}) {
		t.Fatalf("alternatives lost regional locales: %+v", alternatives)
	}
	markdown, err := os.ReadFile(filepath.Join(firstOut, "alternatives.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(markdown), "pt-BR") || !strings.Contains(string(markdown), "ja-JP") {
		t.Fatalf("Markdown hid regional locales: %s", markdown)
	}

	replayOut := filepath.Join(t.TempDir(), "replay")
	replayArgs := []string{"--idea", filepath.Join(firstOut, "idea.json"), "--policy", filepath.Join(firstOut, "reuse-policy.json"), "--mode", "replay", "--strategy", "literal", "--language", "go", "--replay", firstOut, "--output", replayOut, "--target-name", "Mealie", "--search-locale", "pt-BR", "--search-locale", "ja-JP"}
	if err := RunIdeaRecommendCommand(replayArgs, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatalf("equivalent canonical locale replay failed: %v", err)
	}
	mismatch := []string{"--idea", filepath.Join(firstOut, "idea.json"), "--policy", filepath.Join(firstOut, "reuse-policy.json"), "--mode", "replay", "--strategy", "literal", "--language", "go", "--replay", firstOut, "--output", filepath.Join(t.TempDir(), "mismatch"), "--target-name", "Mealie", "--search-locale", "en-US"}
	if err := RunIdeaRecommendCommand(mismatch, &bytes.Buffer{}, nil, nil, nil); err == nil {
		t.Fatal("replay accepted a different canonical locale set")
	}
	nonEmptyAsEmptyReplay := []string{"--idea", filepath.Join(firstOut, "idea.json"), "--policy", filepath.Join(firstOut, "reuse-policy.json"), "--mode", "replay", "--strategy", "literal", "--language", "go", "--replay", firstOut, "--output", filepath.Join(t.TempDir(), "non-empty-as-empty"), "--target-name", "Mealie"}
	if err := RunIdeaRecommendCommand(nonEmptyAsEmptyReplay, &bytes.Buffer{}, nil, nil, nil); err == nil {
		t.Fatal("replay accepted a non-empty checkpoint for an empty locale request")
	}
	emptyAsNonEmptyReplay := []string{"--idea", filepath.Join(emptyOut, "idea.json"), "--policy", filepath.Join(emptyOut, "reuse-policy.json"), "--mode", "replay", "--strategy", "literal", "--language", "go", "--replay", emptyOut, "--output", filepath.Join(t.TempDir(), "empty-as-non-empty"), "--target-name", "Mealie", "--search-locale", "pt-BR"}
	if err := RunIdeaRecommendCommand(emptyAsNonEmptyReplay, &bytes.Buffer{}, nil, nil, nil); err == nil {
		t.Fatal("replay accepted an empty checkpoint for a non-empty locale request")
	}

	resumeProvider := &m21CLIRegionalDiscovery{}
	resumeArgs := []string{"--idea", filepath.Join(firstOut, "idea.json"), "--policy", filepath.Join(firstOut, "reuse-policy.json"), "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--output", firstOut, "--target-name", "Mealie", "--resume", "--search-locale", "pt-BR", "--search-locale", "ja-JP"}
	if err := RunIdeaRecommendCommand(resumeArgs, &bytes.Buffer{}, nil, resumeProvider, fixtureContent{}); err != nil {
		t.Fatalf("equivalent canonical locale resume failed: %v", err)
	}
	if len(resumeProvider.requests) != 0 {
		t.Fatalf("equivalent locale resume repeated discovery: %+v", resumeProvider.requests)
	}
	differentProvider := &m21CLIRegionalDiscovery{}
	differentArgs := []string{"--idea", filepath.Join(firstOut, "idea.json"), "--policy", filepath.Join(firstOut, "reuse-policy.json"), "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--output", firstOut, "--target-name", "Mealie", "--resume", "--search-locale", "en-US"}
	if err := RunIdeaRecommendCommand(differentArgs, &bytes.Buffer{}, nil, differentProvider, fixtureContent{}); err == nil && len(differentProvider.requests) == 0 {
		t.Fatal("resume reused a checkpoint from a different canonical locale set")
	}
	nonEmptyAsEmptyProvider := &m21CLIRegionalDiscovery{}
	nonEmptyAsEmptyResume := []string{"--idea", filepath.Join(thirdOut, "idea.json"), "--policy", filepath.Join(thirdOut, "reuse-policy.json"), "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--output", thirdOut, "--target-name", "Mealie", "--resume"}
	if err := RunIdeaRecommendCommand(nonEmptyAsEmptyResume, &bytes.Buffer{}, nil, nonEmptyAsEmptyProvider, fixtureContent{}); err == nil && len(nonEmptyAsEmptyProvider.requests) == 0 {
		t.Fatal("resume reused a non-empty checkpoint for an empty locale request")
	}
	emptyAsNonEmptyProvider := &m21CLIRegionalDiscovery{}
	emptyAsNonEmptyResume := []string{"--idea", filepath.Join(emptyOut, "idea.json"), "--policy", filepath.Join(emptyOut, "reuse-policy.json"), "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--output", emptyOut, "--target-name", "Mealie", "--resume", "--search-locale", "pt-BR"}
	if err := RunIdeaRecommendCommand(emptyAsNonEmptyResume, &bytes.Buffer{}, nil, emptyAsNonEmptyProvider, fixtureContent{}); err == nil && len(emptyAsNonEmptyProvider.requests) == 0 {
		t.Fatal("resume reused an empty checkpoint for a non-empty locale request")
	}
}

func TestIdeaRecommendM21ZeroLocaleGoldenIdentity(t *testing.T) {
	const (
		legacyRunID      = "4dbe6d1bec49f5c0"
		wantEvidenceHash = "a92ab04e7601769fe61d3c9511bf0060e26069c23d777d5cfb946e04b83eb3a0"
		wantBundleHash   = "4ee79fb89117467e2eb70874c32fb77a922a3adb5807e675a0b1714e365c9912"
	)
	dossier := filepath.Join("..", "..", "lab", "bancada", "history", "benchmarks", "0024", "proof-recipes", "final-dossier")
	out := filepath.Join(t.TempDir(), "out")
	args := []string{
		"--idea", filepath.Join(dossier, "idea.json"),
		"--policy", filepath.Join(dossier, "reuse-policy.json"),
		"--mode", "replay",
		"--strategy", "literal",
		"--replay", dossier,
		"--output", out,
		"--target-name", "Mealie",
		"--target-repo", "mealie-recipes/mealie",
	}
	if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	var manifest domain.IdeaRunManifest
	if err := readJSONFile(filepath.Join(out, "run.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	wantVersions := map[string]string{"expansion": "1", "screening": "1", "evidence": usecase.DefaultEvidenceRules().Version, "decision": "0043", "alternatives": "0043", "analyzer": "0043"}
	if !reflect.DeepEqual(manifest.RuleVersions, wantVersions) {
		t.Fatalf("run version contract: got %+v want %+v", manifest.RuleVersions, wantVersions)
	}
	identity := struct {
		EvidenceHash string            `json:"evidence_hash"`
		PolicyHash   string            `json:"policy_hash"`
		TargetHash   string            `json:"target_hash"`
		RuleVersions map[string]string `json:"rule_versions"`
	}{manifest.InputHashes["evidence"], manifest.InputHashes["policy"], manifest.InputHashes["target"], wantVersions}
	identityBytes, err := json.Marshal(identity)
	if err != nil {
		t.Fatal(err)
	}
	wantRunID := fmt.Sprintf("%x", sha256.Sum256(identityBytes))[:16]
	if manifest.RunID != wantRunID || manifest.RunID == legacyRunID || manifest.InputHashes["evidence"] != wantEvidenceHash {
		t.Fatalf("versioned run identity or unchanged evidence violated: run=%q expected=%q evidence=%q", manifest.RunID, wantRunID, manifest.InputHashes["evidence"])
	}
	if _, ok := manifest.InputHashes["search_locales"]; ok {
		t.Fatalf("zero-locale legacy identity gained search_locales: %+v", manifest.InputHashes)
	}
	bundles, err := os.ReadFile(filepath.Join(out, "evidence", "bundles.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(bundles)); got != wantBundleHash {
		t.Fatalf("zero-locale evidence bytes changed: %s", got)
	}

	repeatedOut := filepath.Join(t.TempDir(), "repeat")
	repeatedArgs := append([]string(nil), args...)
	for i, arg := range repeatedArgs {
		if arg == "--output" {
			repeatedArgs[i+1] = repeatedOut
		}
	}
	if err := RunIdeaRecommendCommand(repeatedArgs, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	var repeated domain.IdeaRunManifest
	if err := readJSONFile(filepath.Join(repeatedOut, "run.json"), &repeated); err != nil {
		t.Fatal(err)
	}
	if repeated.RunID != manifest.RunID || !reflect.DeepEqual(repeated.InputHashes, manifest.InputHashes) {
		t.Fatal("identical replay changed versioned identity")
	}
}

func TestIdeaRecommendM21RejectsInvalidLocaleBeforeProviderAndNoFlagIsNeutral(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	provider := &m21CLIRegionalDiscovery{}
	args := []string{"--idea", idea, "--policy", policy, "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--output", filepath.Join(t.TempDir(), "invalid"), "--search-locale", "pt-br"}
	if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, provider, fixtureContent{}); err == nil {
		t.Fatal("invalid locale was accepted")
	}
	if len(provider.requests) != 0 {
		t.Fatalf("provider called before locale validation: %+v", provider.requests)
	}

	out := filepath.Join(t.TempDir(), "legacy")
	if err := RunIdeaRecommendCommand([]string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", out}, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range []string{"expansion-plan.json", "discovery-snapshot.json", "alternatives.json", "alternatives.md"} {
		data, err := os.ReadFile(filepath.Join(out, artifact))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "search_locale") || strings.Contains(string(data), "regional_proposal") || strings.Contains(string(data), "Locales de descoberta") {
			t.Fatalf("no-locale run changed legacy %s: %s", artifact, data)
		}
	}
}

func TestIdeaRecommendUsesExplicitScreeningModel(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	modelPath := filepath.Join("..", "..", "00_nucleo", "modelos", "screening-operacional-v1.json")
	modelBytes, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatal(err)
	}

	run := func(path string) domain.IdeaRunManifest {
		t.Helper()
		out := filepath.Join(t.TempDir(), "out")
		args := []string{
			"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal",
			"--replay", source, "--output", out, "--screening-model", path,
		}
		if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
		var manifest domain.IdeaRunManifest
		if err := readJSONFile(filepath.Join(out, "run.json"), &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.InputHashes["screening_model"] != hashBytes(mustReadFile(t, path)) {
			t.Fatalf("manifesto não identifica os bytes do modelo: %+v", manifest.InputHashes)
		}
		if manifest.RuleVersions["screening"] != "screening-opt-in-v1" {
			t.Fatalf("versão do screening opt-in ausente: %+v", manifest.RuleVersions)
		}
		var report domain.ScreeningReport
		if err := readJSONFile(filepath.Join(out, "screening-report.json"), &report); err != nil {
			t.Fatal(err)
		}
		for _, decision := range report.Decisions {
			if decision.Decision == "DROP" || decision.Probability != nil {
				t.Fatalf("modelo não calibrado ganhou autoridade destrutiva: %+v", decision)
			}
		}
		return manifest
	}

	first := run(modelPath)
	copyPath := filepath.Join(t.TempDir(), "model-with-extra-newline.json")
	if err := os.WriteFile(copyPath, append(append([]byte(nil), modelBytes...), '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	second := run(copyPath)
	if first.RunID == second.RunID {
		t.Fatal("mudança nos bytes do modelo não invalidou a identidade da execução")
	}
}

func TestIdeaRecommendRejectsInvalidScreeningModelBeforeOutput(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	for _, testCase := range []struct {
		name    string
		content string
	}{
		{name: "unknown-field", content: `{"schema_version":"1","weights":{},"advance_threshold":0.5,"drop_threshold":0.2,"audit_sample_rate":0,"audit_seed":1,"calibration_status":"uncalibrated","extra":true}`},
		{name: "calibrated", content: `{"schema_version":"1","weights":{},"advance_threshold":0.5,"drop_threshold":0.2,"audit_sample_rate":0,"audit_seed":1,"calibration_status":"calibrated"}`},
		{name: "trailing-value", content: `{"schema_version":"1","weights":{},"advance_threshold":0.5,"drop_threshold":0.2,"audit_sample_rate":0,"audit_seed":1,"calibration_status":"uncalibrated"} {}`},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			modelPath := filepath.Join(root, "model.json")
			if err := os.WriteFile(modelPath, []byte(testCase.content), 0o600); err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(root, "must-not-exist")
			err := RunIdeaRecommendCommand([]string{
				"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal",
				"--replay", source, "--output", out, "--screening-model", modelPath,
			}, &bytes.Buffer{}, nil, nil, nil)
			if err == nil || err.Error() != "modelo de screening inválido" {
				t.Fatalf("erro = %v", err)
			}
			if _, statErr := os.Stat(out); !os.IsNotExist(statErr) {
				t.Fatalf("modelo inválido deixou saída parcial: %v", statErr)
			}
		})
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
