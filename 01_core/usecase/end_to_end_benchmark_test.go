/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/end_to_end_benchmark_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func benchmarkCaseFixture(id string) domain.EndToEndBenchmarkCase {
	idea := domain.IdeaSpec{Problem: "analisar dependências", Capabilities: []string{"dependency graph"}}
	ref := domain.CommitRef{Owner: "org", Repo: "graph", Commit: "abc123"}
	candidate := domain.RepositoryCandidate{Owner: "org", Repo: "graph", CommitID: "abc123", Description: "dependency graph"}
	metadata, _ := json.Marshal(domain.RepositoryCandidate{Owner: "org", Repo: "graph", CommitID: "abc123", UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)})
	items := []domain.CollectedEvidence{{Kind: domain.EvidenceRepository, Ref: ref, Path: "@repository-metadata", Content: string(metadata), ContentHash: contentHash(string(metadata))}, {Kind: domain.EvidenceFile, Ref: ref, Path: "graph.go", Content: "type DependencyGraph struct{}", ContentHash: contentHash("type DependencyGraph struct{}")}, {Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: "func TestDependencyGraph(){}", ContentHash: contentHash("func TestDependencyGraph(){}")}, {Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: "MIT License Permission is hereby granted, free of charge", ContentHash: contentHash("MIT License Permission is hereby granted, free of charge")}}
	rules := DefaultEvidenceRules()
	c := domain.EndToEndBenchmarkCase{SchemaVersion: "1", ID: domain.BenchmarkCaseID(id), Idea: idea, Discovery: domain.DiscoverySnapshot{SchemaVersion: "2", Idea: idea, Strategy: "literal", Candidates: []domain.DiscoveredCandidate{{Repository: candidate, ExpansionIDs: []string{"cap"}, Queries: []string{"dependency graph"}, Origins: []domain.CandidateOrigin{{ExpansionID: "cap", PlannedQuery: "dependency graph", EffectiveQuery: "dependency graph", Evidence: []domain.MatchEvidence{{MatchType: "content", Path: "graph.go", Snippet: "DependencyGraph"}}}}}}}, ScreeningModel: DefaultScreeningModel(), Bundles: []domain.EvidenceBundle{{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: items}}, EvidenceRules: rules, ReusePolicy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}, RecoveryPolicy: DefaultRecoveryPolicy(), ExpectedEvidence: []domain.EvidenceJudgment{{Repository: "org/graph", Capability: "dependency graph", Expected: domain.EvidenceProven}}, ExpectedDecision: domain.ReuseDecisionJudgment{Expected: domain.DecisionIntegrate, Forbidden: []domain.ReuseDecision{domain.DecisionBuild}, RequiredRules: []string{"DEC-INTEGRATE-001"}, RequiredSources: []string{"graph.go"}, Rationale: "componente delimitado"}, UtilityJudgment: domain.DossierUtilityJudgment{Label: domain.UtilityUseful, DecisionActionable: true, EvidenceSufficient: true, BlockersVisible: true, Rationale: "ponto de partida explícito"}, Provenance: domain.BenchmarkProvenance{CapturedAt: "2026-08-26T00:00:00Z", SourceMode: domain.IdeaRunReplay, OriginalRunID: id, OriginalRunHash: "runhash", ArtifactHashes: map[string]string{"bundle": "hash"}, CandidateCommits: map[string]string{"org/graph": "abc123"}, CollectorVersion: "1"}, Split: domain.SplitContract, Group: id}
	return c
}
func manifestForCases(cases ...domain.EndToEndBenchmarkCase) domain.BenchmarkManifest {
	m := domain.BenchmarkManifest{ProtocolVersion: "1", PromptPath: "step-0011", PromptHash: "prompt", BaselineCommit: "baseline", ContractPath: "benchmark.go", ContractHash: "contract", UnknownPolicyHash: "unknown", CaseHashes: map[domain.BenchmarkCaseID]string{}, SplitGroups: map[domain.BenchmarkCaseID]string{}, RoleCapabilities: []domain.RoleCapabilityReceipt{{Role: "ensaio", Executor: "same-session", Limitation: "sem atestação de isolamento"}}}
	for _, c := range cases {
		m.CaseHashes[c.ID] = BenchmarkCaseHash(c)
		m.SplitGroups[c.ID] = c.Group
	}
	return m
}
func TestValidCorrespondingCaseIsPreserved(t *testing.T) {
	c := benchmarkCaseFixture("valid")
	got := ValidateBenchmarkCase(manifestForCases(c), c)
	if got.Verdict != domain.CorrespondencePreserved {
		t.Fatalf("%s %#v", got.Verdict, got.Witnesses)
	}
}
func TestOpaqueHeadIsUnknownWhenManifestMatches(t *testing.T) {
	c := benchmarkCaseFixture("head")
	c.Bundles[0].Repository.Commit = "HEAD"
	for i := range c.Bundles[0].Items {
		c.Bundles[0].Items[i].Ref.Commit = "HEAD"
	}
	c.Discovery.Candidates[0].Repository.CommitID = "HEAD"
	c.Provenance.CandidateCommits["org/graph"] = "HEAD"
	c.ExpectedDecision.Expected = domain.DecisionUnknown
	got := ValidateBenchmarkCase(manifestForCases(c), c)
	if got.Verdict != domain.CorrespondenceUnknown {
		t.Fatalf("got %s", got.Verdict)
	}
}
func TestIdeaBundleCommitAndHashMutationsAreViolated(t *testing.T) {
	base := benchmarkCaseFixture("mut")
	mutations := []domain.EndToEndBenchmarkCase{}
	a := base
	a.Discovery.Idea.Problem = "outra"
	mutations = append(mutations, a)
	b := base
	b.Bundles[0].Repository.Repo = "other"
	mutations = append(mutations, b)
	c := base
	c.Bundles[0].Repository.Commit = "different"
	mutations = append(mutations, c)
	d := base
	d.Bundles[0].Items[0].Content = "tampered"
	mutations = append(mutations, d)
	for i, m := range mutations {
		if got := ValidateBenchmarkCase(manifestForCases(m), m); got.Verdict != domain.CorrespondenceViolated {
			t.Fatalf("mutation %d survived: %s", i, got.Verdict)
		}
	}
}
func TestCaseHashAndVerdictIgnoreOrdering(t *testing.T) {
	c := benchmarkCaseFixture("order")
	d := c
	d.Bundles = append([]domain.EvidenceBundle(nil), c.Bundles...)
	d.Bundles[0].Items = append([]domain.CollectedEvidence(nil), c.Bundles[0].Items...)
	for i, j := 0, len(d.Bundles[0].Items)-1; i < j; i, j = i+1, j-1 {
		d.Bundles[0].Items[i], d.Bundles[0].Items[j] = d.Bundles[0].Items[j], d.Bundles[0].Items[i]
	}
	if BenchmarkCaseHash(c) != BenchmarkCaseHash(d) {
		t.Fatal("order changed hash")
	}
	if !reflect.DeepEqual(ValidateBenchmarkCase(manifestForCases(c), c).Verdict, ValidateBenchmarkCase(manifestForCases(c), d).Verdict) {
		t.Fatal("order changed verdict")
	}
}
func TestTwelveSealedMutationsReachPerfectScore(t *testing.T) {
	base := benchmarkCaseFixture("sealed")
	manifest := manifestForCases(base)
	mutations := make([]domain.BenchmarkMutation, 12)
	for i := range mutations {
		m := base
		m.ExpectedDecision.Rationale += string(rune('a' + i))
		mutations[i] = domain.BenchmarkMutation{ID: string(rune('a' + i)), Case: m}
	}
	r := EvaluateBenchmarkMutations(manifest, mutations)
	if !r.Sealed || r.MutationScore == nil || *r.MutationScore != 1 {
		t.Fatalf("%#v", r)
	}
}
func TestEmptyBenchmarkMetricsRemainNull(t *testing.T) {
	r, err := RunEndToEndBenchmark(domain.BenchmarkManifest{}, domain.BenchmarkCorpus{}, domain.SplitTest)
	if err != nil {
		t.Fatal(err)
	}
	if r.Metrics.DecisionAccuracy != nil || r.Metrics.EvidencePrecision != nil {
		t.Fatalf("invented metrics %#v", r.Metrics)
	}
}
