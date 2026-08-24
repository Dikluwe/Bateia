/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestDefaultEvidenceRules(t *testing.T) {
	rules := DefaultEvidenceRules()
	if len(rules.Rules) == 0 {
		t.Errorf("expected non-empty rules catalog")
	}
	if err := ValidateRuleCatalog(rules); err != nil {
		t.Errorf("expected valid rule catalog, got %v", err)
	}
}

func TestAssessEvidence(t *testing.T) {
	ref := domain.CommitRef{Owner: "owner", Repo: "repo", Commit: "1234567890123456789012345678901234567890"}
	spec := domain.IdeaSpec{
		Problem:      "Test Problem",
		Capabilities: []string{"grafo de dependencias"},
	}
	screening := domain.ScreeningReport{}
	bundles := []domain.EvidenceBundle{
		{
			Repository: ref,
			Complete:   true,
			Items: []domain.CollectedEvidence{
				{
					Kind:        domain.EvidenceTest,
					Ref:         ref,
					Path:        "file_test.go",
					Content:     "func TestDependencyGraph(t *testing.T) {}",
					ContentHash: "hash123",
				},
			},
		},
	}
	catalog := DefaultEvidenceRules()
	report, err := AssessEvidence(spec, screening, bundles, catalog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Assessments) == 0 {
		t.Errorf("expected candidate assessments")
	}
}
