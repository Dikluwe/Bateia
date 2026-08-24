/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func testBundle(complete bool, items ...domain.CollectedEvidence) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	for i := range items {
		items[i].Ref = ref
		if items[i].ContentHash == "" {
			items[i].ContentHash = contentHash(items[i].Content)
		}
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: complete, Items: items}
}

func TestDocumentationAloneDoesNotProveCapability(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Provides dependency graph analysis."})
	got := assessCapability("grafo de dependências", b, findRule("grafo de dependências", DefaultEvidenceRules()))
	if got.Status == domain.EvidenceProven {
		t.Fatalf("documentação isolada produziu PROVEN: %+v", got)
	}
}

func TestImplementedSymbolAndRelatedTestCanProve(t *testing.T) {
	b := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "graph.go", Content: "type DependencyGraph struct { Nodes []string }"},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "graph_test.go", Content: "func TestDependencyGraph(t *testing.T) {}"},
	)
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceProven {
		t.Fatalf("status = %s, esperado PROVEN", got.Status)
	}
	if len(got.Sources) == 0 || got.Sources[0].Ref.Commit == "" {
		t.Fatal("fonte não rastreável")
	}
}

func TestFailureNeverBecomesAbsent(t *testing.T) {
	b := testBundle(false)
	b.Failures = []domain.ExecutionFailure{{Scope: "org/repo@abc123", Stage: "content", Message: "timeout"}}
	got := assessCapability("static analysis", b, findRule("static analysis", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceInconclusive {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestVendorNoiseCannotProve(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "vendor/k8s/graph_test.go", Content: "func TestDependencyGraph(t *testing.T) {}"})
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status == domain.EvidenceProven {
		t.Fatalf("vendor produziu PROVEN: %+v", got)
	}
}

func TestAssessmentIsDeterministicAndDoesNotMutateBundle(t *testing.T) {
	itemA := domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "z_test.go", Content: "func TestStaticAnalyzer(t *testing.T) {}"}
	itemB := domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "static analysis"}
	a := testBundle(true, itemA, itemB)
	before := a
	b := testBundle(true, itemB, itemA)
	rule := findRule("static analysis", DefaultEvidenceRules())
	left := assessCapability("static analysis", a, rule)
	right := assessCapability("static analysis", b, rule)
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("ordem alterou resultado\n%+v\n%+v", left, right)
	}
	if !reflect.DeepEqual(a, before) {
		t.Fatal("entrada foi modificada")
	}
}

func TestPlanRequiresImmutableCommitAndSelectsDeterministically(t *testing.T) {
	report := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{
		{Repository: "z/drop@1", Decision: "DROP", RankingScore: .1}, {Repository: "a/drop@1", Decision: "DROP", RankingScore: .2},
		{Repository: "org/a@abc", Decision: "ADVANCE", RankingScore: .9},
	}}
	plans, err := PlanEvidence(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 3 || plans[1].Repository.Owner != "a" {
		t.Fatalf("seleção inesperada: %+v", plans)
	}
	report.Decisions[0].Repository = "org/mutable"
	if _, err := PlanEvidence(report); err == nil {
		t.Fatal("commit mutável foi aceito")
	}
}

type evidenceProvider struct {
	files   []domain.ExtractedFile
	treeErr error
}

func (p evidenceProvider) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("unused")
}
func (p evidenceProvider) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	if p.treeErr != nil {
		return nil, p.treeErr
	}
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "graph_test.go", Type: "file"}}, nil
}
func (p evidenceProvider) ExtractRepoFiles(_ context.Context, _ domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{Files: p.files}, nil
}
func (p evidenceProvider) FetchRepoMetadata(context.Context, domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return domain.RepositoryCandidate{}, nil
}

func TestCollectDeduplicatesByHashAndPreservesPartialFailure(t *testing.T) {
	plan := domain.EvidenceExtractionPlan{Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}, Candidate: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc"}, Budget: defaultEvidenceBudget}
	p := evidenceProvider{treeErr: errors.New("tree unavailable"), files: []domain.ExtractedFile{{Path: "README.md", Content: "same"}, {Path: "COPY.md", Content: "same"}}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if b.Complete || len(b.Failures) != 1 || len(b.Items) != 1 {
		t.Fatalf("bundle parcial inesperado: %+v", b)
	}
}
