/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/implementation_retrieval_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"
)

type implementationContentStub struct{}

func (implementationContentStub) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, nil
}
func (implementationContentStub) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return nil, nil
}
func (implementationContentStub) ExtractRepoFiles(context.Context, domain.RepositoryCandidate, []string) (*domain.ExtractedRepository, error) {
	return nil, nil
}
func (implementationContentStub) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return c, nil
}

type implementationSearchStub struct {
	implementationContentStub
	report domain.ImplementationSearchReport
	err    error
}

func (s implementationSearchStub) SearchImplementations(context.Context, domain.RepositoryCandidate, string, domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	return s.report, s.err
}

func TestImplementationRetrievalOptionalAndError(t *testing.T) {
	r := usecase.RetrieveImplementationEvidence(context.Background(), implementationContentStub{}, domain.RepositoryCandidate{}, "query")
	if r.Status != "not_evaluated" || r.Reason == "" {
		t.Fatalf("optional provider unknown lost: %+v", r)
	}
	r = usecase.RetrieveImplementationEvidence(context.Background(), implementationSearchStub{err: errors.New("controlled failure")}, domain.RepositoryCandidate{}, "query")
	if r.Status == "complete" || (r.Reason == "" && len(r.Issues) == 0) {
		t.Fatalf("failure concealed: %+v", r)
	}
}

func TestImplementationRetrievalSelectionPrioritizesHit(t *testing.T) {
	tree := []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "go.mod", Type: "file"}, {Path: "queue.go", Type: "file"}, {Path: "neutral.go", Type: "file"}}
	report := domain.ImplementationSearchReport{Status: "complete", Hits: []domain.ImplementationHit{{Path: "neutral.go", Symbol: "QueueItems", Score: 100}}}
	paths := usecase.SelectImplementationPaths("queue", tree, report, 1)
	if len(paths) > 3 {
		t.Fatalf("limit exceeded: %v", paths)
	}
	found := map[string]bool{}
	for _, p := range paths {
		found[p] = true
	}
	if !found["neutral.go"] || !found["README.md"] || !found["go.mod"] {
		t.Fatalf("hit or foundational document lost: %v", paths)
	}
}

func TestImplementationRetrievalRejectsChangedBytes(t *testing.T) {
	old := "package sample\nfunc QueueItems() {}\n"
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(old)))
	report := domain.ImplementationSearchReport{Status: "complete", Hits: []domain.ImplementationHit{{Path: "neutral.go", Symbol: "QueueItems", SHA256: hash, Line: 2, EndLine: 2}}, Examined: []domain.ImplementationInput{{Path: "neutral.go", SHA256: hash, Bytes: int64(len(old))}}}
	unchanged := usecase.ValidateImplementationHits(report, []domain.ExtractedFile{{Path: "neutral.go", Content: old}})
	if len(unchanged.Hits) != 1 {
		t.Fatalf("unchanged definition discarded: %+v", unchanged)
	}
	changed := usecase.ValidateImplementationHits(report, []domain.ExtractedFile{{Path: "neutral.go", Content: old + "// mutation\n", Evidence: domain.FileEvidence{SHA256: hash}}})
	if len(changed.Hits) != 0 || changed.Status == "complete" || len(changed.Issues) == 0 {
		t.Fatalf("stale hash remained confirmed: %+v", changed)
	}
	if len(changed.Examined) != 1 || changed.Examined[0].SHA256 != hash {
		t.Fatalf("index provenance overwritten: %+v", changed)
	}
}
