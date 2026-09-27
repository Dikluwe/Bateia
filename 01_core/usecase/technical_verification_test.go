/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_verification_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"
	"time"

	"bateia/01_core/domain"
)

func verificationFixture() (domain.TechnicalEvidenceBatch, domain.EvidenceBundle, domain.CommitRef, time.Time) {
	ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: "abc123"}
	content := step0025MITLicense
	item := domain.CollectedEvidence{Kind: domain.EvidenceLicense, Ref: ref, Path: "LICENSE", Content: content, ContentHash: contentHash(content)}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Plan: domain.EvidenceExtractionPlan{Repository: ref}, Items: []domain.CollectedEvidence{item}, Complete: true}
	finding := domain.TechnicalEvidence{
		ID: "lic", Technique: domain.TechniqueLicense, Candidate: ref.Key(), Claim: "licença raiz identificada como MIT",
		Predicate: domain.TechnicalLicenseID, Value: "MIT", Result: domain.TechnicalSupports,
		Source: domain.TechnicalEvidenceSource{Ref: "LICENSE", Revision: ref.Commit, Digest: "sha256:" + item.ContentHash},
		Method: "comparar bytes coletados e detectar identificador", ObservedAt: "2026-09-20T00:00:00Z",
	}
	return domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{finding}}, bundle, ref, time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC)
}

func TestVerifyTechnicalLicenseAndApplyOnlyUnknownGate(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if len(review.Verifications) != 1 || review.Verifications[0].Status != domain.TechnicalVerified || review.Verifications[0].Source == nil {
		t.Fatalf("expected verified source: %+v", review)
	}
	assessments, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: ref, License: domain.LicenseAssessment{Status: domain.EvidenceInconclusive}}}, review)
	if !changed || assessments[0].License.Status != domain.EvidenceProven || assessments[0].License.Identifier != "MIT" {
		t.Fatalf("license gate not filled: %+v", assessments)
	}
	preexisting := domain.CandidateAssessment{Repository: ref, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "Apache-2.0", Sources: assessments[0].License.Sources}}
	_, changed = ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{preexisting}, review)
	if changed {
		t.Fatal("verified technical finding replaced an already proven assessment")
	}
}

func TestVerifyTechnicalEvidenceRejectsTamperedOrAmbiguousBundles(t *testing.T) {
	batch, bundle, _, asOf := verificationFixture()
	cases := []domain.EvidenceBundle{
		func() domain.EvidenceBundle {
			b := bundle
			b.Items = append([]domain.CollectedEvidence(nil), b.Items...)
			b.Items[0].Content += " changed"
			return b
		}(),
		func() domain.EvidenceBundle { b := bundle; b.Repository.Commit = "different"; return b }(),
		func() domain.EvidenceBundle { b := bundle; b.Complete = false; return b }(),
		func() domain.EvidenceBundle {
			b := bundle
			b.Items = append(append([]domain.CollectedEvidence(nil), b.Items...), b.Items[0])
			return b
		}(),
	}
	for _, b := range cases {
		review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{b}, asOf)
		if err != nil {
			t.Fatal(err)
		}
		if review.Verifications[0].Status == domain.TechnicalVerified {
			t.Fatalf("tampered or ambiguous bundle was verified: %+v", b)
		}
	}
}

func TestVerifyTechnicalSBOMDeclarationDoesNotFillLicenseGate(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	content := `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"widget","version":"1.2.3"}]}`
	bundle.Items = []domain.CollectedEvidence{{Kind: domain.EvidenceManifest, Ref: ref, Path: "sbom.json", Content: content, ContentHash: contentHash(content)}}
	f := batch.Findings[0]
	f.ID, f.Technique, f.Claim, f.Predicate, f.Value = "sbom", domain.TechniqueSBOM, "SBOM declara widget 1.2.3", domain.TechnicalSBOMComponent, "widget@1.2.3"
	f.Source.Ref, f.Source.Digest = "sbom.json", "sha256:"+contentHash(content)
	batch.Findings = []domain.TechnicalEvidence{f}
	review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if review.Verifications[0].Status != domain.TechnicalVerified || review.Verifications[0].Source.Kind != domain.EvidenceManifest {
		t.Fatalf("SBOM declaration not verified: %+v", review)
	}
	_, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: ref}}, review)
	if changed {
		t.Fatal("SBOM declaration changed license gate")
	}
}

func TestVerifyTechnicalEvidenceNeverPromotesConflictOrStale(t *testing.T) {
	batch, bundle, _, asOf := verificationFixture()
	refutation := batch.Findings[0]
	refutation.ID, refutation.Result = "refute", domain.TechnicalRefutes
	refutation.Claim = "texto diferente para a mesma licença"
	batch.Findings = append(batch.Findings, refutation)
	review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range review.Verifications {
		if v.Status != domain.TechnicalContradicted {
			t.Fatalf("conflicting claim was promoted: %+v", v)
		}
	}
	batch.Findings = batch.Findings[:1]
	batch.Findings[0].ValidUntil = "2026-09-21T00:00:00Z"
	review, err = VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil || review.Verifications[0].Status != domain.TechnicalStale {
		t.Fatalf("stale claim was promoted: %v %+v", err, review)
	}
}
