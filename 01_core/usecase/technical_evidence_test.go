/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_evidence_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"
	"time"

	"bateia/01_core/domain"
)

func technicalFinding(id string, result domain.TechnicalEvidenceResult) domain.TechnicalEvidence {
	return domain.TechnicalEvidence{
		ID: id, Technique: domain.TechniqueLicense, Candidate: "o/r", Claim: "licença compatível com a composição", Result: result,
		Source: domain.TechnicalEvidenceSource{Ref: "LICENSE", Revision: "abc"}, Method: "inspeção do arquivo e modo de distribuição",
		ObservedAt: "2026-09-20T00:00:00Z",
	}
}

func TestReviewTechnicalEvidenceConflictAndStale(t *testing.T) {
	asOf := time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC)
	old := technicalFinding("old", domain.TechnicalSupports)
	old.ValidUntil = "2026-09-22T00:00:00Z"
	batch := domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{old, technicalFinding("a", domain.TechnicalSupports), technicalFinding("b", domain.TechnicalRefutes)}}
	review, err := ReviewTechnicalEvidence(batch, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if review.DecisionImpact != "ADVISORY_ONLY" || len(review.Claims) != 1 || review.Claims[0].Status != domain.TechnicalClaimConflicting {
		t.Fatalf("unexpected review: %+v", review)
	}
	if len(review.Claims[0].StaleIDs) != 1 || review.Claims[0].StaleIDs[0] != "old" || len(review.Warnings) != 2 {
		t.Fatalf("staleness/conflict not recorded: %+v", review)
	}
}

func TestReviewTechnicalEvidenceNeverTreatsExpiredSupportAsCurrent(t *testing.T) {
	f := technicalFinding("expired", domain.TechnicalSupports)
	f.ValidUntil = "2026-09-21T00:00:00Z"
	review, err := ReviewTechnicalEvidence(domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{f}}, time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if review.Claims[0].Status != domain.TechnicalClaimStale {
		t.Fatalf("expired support was promoted: %+v", review.Claims[0])
	}
}

func TestReviewTechnicalEvidenceRejectsInvalidProvenance(t *testing.T) {
	asOf := time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC)
	base := technicalFinding("x", domain.TechnicalSupports)
	cases := []domain.TechnicalEvidence{
		func() domain.TechnicalEvidence { f := base; f.Source.Ref = ""; return f }(),
		func() domain.TechnicalEvidence { f := base; f.Source.Revision = ""; return f }(),
		func() domain.TechnicalEvidence { f := base; f.ObservedAt = "2027-01-01T00:00:00Z"; return f }(),
		func() domain.TechnicalEvidence { f := base; f.Technique = "UNKNOWN"; return f }(),
		func() domain.TechnicalEvidence { f := base; f.ValidUntil = "2026-09-19T00:00:00Z"; return f }(),
	}
	for _, f := range cases {
		if _, err := ReviewTechnicalEvidence(domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{f}}, asOf); err == nil {
			t.Fatalf("accepted invalid finding: %+v", f)
		}
	}
	if _, err := ReviewTechnicalEvidence(domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{base, base}}, asOf); err == nil {
		t.Fatal("accepted duplicate finding ID")
	}
}
