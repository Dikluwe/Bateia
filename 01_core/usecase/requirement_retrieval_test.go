/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/requirement_retrieval_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestRequirement0046RevalidationBeforeProjection(t *testing.T) {
	content := "package fixture\n\nfunc Queue() {}\n"
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	hit := domain.ImplementationHit{Retrieval: "lexical", Path: "queue.go", Symbol: "Queue", Line: 3, EndLine: 3, SHA256: hash, MatchedTerms: []string{"queue"}}
	old := domain.RequirementRetrieval{Requirement: "queue", Status: "evidence_located", MatchedTerms: []string{"queue"}, Sources: []domain.ImplementationReference{{Path: hit.Path, Symbol: hit.Symbol, Line: hit.Line, SHA256: hash}}}
	for _, tc := range []struct {
		name  string
		files []domain.ExtractedFile
		want  bool
	}{
		{"identical", []domain.ExtractedFile{{Path: "queue.go", Content: content, Size: int64(len(content))}}, true},
		{"changed", []domain.ExtractedFile{{Path: "queue.go", Content: "package fixture\n\nfunc Other() {}\n"}}, false},
		{"missing", nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			report := domain.ImplementationSearchReport{Status: "completed", Hits: []domain.ImplementationHit{hit}, Requirements: []domain.RequirementRetrieval{old}}
			got := usecase.ValidateImplementationRequirements("queue", report, tc.files)
			if len(got.Requirements) != 1 {
				t.Fatalf("requirement lost: %#v", got.Requirements)
			}
			r := got.Requirements[0]
			if tc.want {
				if r.Status != "evidence_located" || len(r.Sources) != 1 || r.Sources[0].SHA256 != hash {
					t.Fatalf("fresh source lost: %#v", r)
				}
			} else {
				if r.Status != "not_verified" || len(r.Sources) != 0 || len(r.MatchedTerms) != 0 {
					t.Fatalf("stale source survived: %#v", r)
				}
			}
		})
	}
}
