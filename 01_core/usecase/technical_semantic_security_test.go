/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_semantic_security_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"strings"
	"testing"
)

func TestTechnicalLicenseInventoryIgnoresCallerKindForAmbiguity(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	for _, kind := range []domain.EvidenceKind{domain.EvidenceFile, domain.EvidenceDocumentation, domain.EvidenceManifest, domain.EvidenceLicense} {
		b := bundle
		b.Items = append(append([]domain.CollectedEvidence(nil), bundle.Items...), domain.CollectedEvidence{Kind: kind, Ref: ref, Path: "COPYing", Content: "Redistribution requires written consent.", ContentHash: contentHash("Redistribution requires written consent.")})
		if got := assessLicense(b); got.Status == domain.EvidenceProven {
			t.Fatalf("hidden root license with kind %s", kind)
		}
		review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{b}, asOf)
		if err != nil {
			t.Fatal(err)
		}
		if review.Verifications[0].Status == domain.TechnicalVerified {
			t.Fatal("inventory ambiguity promoted")
		}
		if _, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: ref}}, review); changed {
			t.Fatal("ambiguous gate filled")
		}
	}
	bundle.Items[0].Kind = domain.EvidenceFile
	if assessLicense(bundle).Status == domain.EvidenceProven {
		t.Fatal("mislabeled sole license became proof")
	}
}

func TestTechnicalLicenseCaseFoldConflictAndUnknownBundleSchema(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	refute := batch.Findings[0]
	refute.ID, refute.Claim, refute.Value, refute.Result = "refute", "different license claim", "mit", domain.TechnicalRefutes
	batch.Findings = append(batch.Findings, refute)
	review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil {
		t.Fatal(err)
	}
	for _, verification := range review.Verifications {
		if verification.Status != domain.TechnicalContradicted {
			t.Fatalf("case conflict lost: %+v", verification)
		}
	}
	if _, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: ref}}, review); changed {
		t.Fatal("conflict promoted")
	}
	batch.Findings = batch.Findings[:1]
	bundle.SchemaVersion = "invented"
	review, err = VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil || review.Verifications[0].Status == domain.TechnicalVerified {
		t.Fatalf("unknown bundle schema promoted: %v %+v", err, review)
	}
}

func TestTechnicalBundlesCannotSplitRootLicenseInventory(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	second := bundle
	second.Items = []domain.CollectedEvidence{{Kind: domain.EvidenceLicense, Ref: ref, Path: "COPYING", Content: "All rights reserved."}}
	if _, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle, second}, asOf); err == nil {
		t.Fatal("split same-revision bundles accepted")
	}
}

func TestTechnicalClaimKeysCannotCollideThroughNUL(t *testing.T) {
	batch, _, _, asOf := verificationFixture()
	a := batch.Findings[0]
	a.Candidate, a.Claim = "o/r\x00x", "y"
	b := a
	b.ID, b.Candidate, b.Claim, b.Result = "other", "o/r", "x\x00y", domain.TechnicalRefutes
	batch.Findings = []domain.TechnicalEvidence{a, b}
	review, err := ReviewTechnicalEvidence(batch, asOf)
	if err != nil || len(review.Claims) != 2 {
		t.Fatalf("claim identities collided: %v %+v", err, review)
	}
	for _, claim := range review.Claims {
		if claim.Status == domain.TechnicalClaimConflicting {
			t.Fatal("unrelated claims conflated")
		}
	}
}

func TestTechnicalSBOMSemanticsAndBudgetCorpus(t *testing.T) {
	for _, test := range []struct {
		name, content, target string
		want                  domain.TechnicalVerificationStatus
	}{
		{"declared", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"widget","version":"1.2.3"}]}`, "widget@1.2.3", domain.TechnicalVerified},
		{"multi_version", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"bom-ref":"a","name":"widget","version":"1.2.3"},{"bom-ref":"b","name":"widget","version":"9.9.9"}]}`, "widget@1.2.3", domain.TechnicalVerified},
		{"unknown_version", `{"bomFormat":"CycloneDX","specVersion":"invented","components":[{"name":"widget","version":"1.2.3"}]}`, "widget@1.2.3", domain.TechnicalUnsupported},
		{"duplicate_ref", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"bom-ref":"x","name":"widget","version":"1.2.3"},{"bom-ref":"x","name":"widget","version":"9.9.9"}]}`, "widget@1.2.3", domain.TechnicalUnverified},
		{"duplicate_name", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"other","name":"widget","version":"1.2.3"}]}`, "widget@1.2.3", domain.TechnicalUnverified},
		{"case_alias", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"Name":"other","name":"widget","version":"1.2.3"}]}`, "widget@1.2.3", domain.TechnicalUnverified},
		{"missing_version", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"widget"}]}`, "widget@", domain.TechnicalUnverified},
		{"missing_target", `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[]}`, "widget@1.2.3", domain.TechnicalUnverified},
	} {
		t.Run(test.name, func(t *testing.T) {
			batch, bundle, ref, asOf := verificationFixture()
			item := domain.CollectedEvidence{Kind: domain.EvidenceManifest, Ref: ref, Path: "sbom.json", Content: test.content, ContentHash: contentHash(test.content)}
			bundle.Items = []domain.CollectedEvidence{item}
			f := &batch.Findings[0]
			f.Technique, f.Predicate, f.Value = domain.TechniqueSBOM, domain.TechnicalSBOMComponent, test.target
			f.Source.Ref, f.Source.Digest = item.Path, "sha256:"+item.ContentHash
			review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
			if err != nil || review.Verifications[0].Status != test.want {
				t.Fatalf("want %s: %v %+v", test.want, err, review)
			}
			if !review.Verifications[0].IntegrityVerified || review.Verifications[0].BehaviorVerified {
				t.Fatal("integrity and interpretation conflated")
			}
		})
	}
	content := `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[` + strings.Repeat(`{"name":"widget","version":"1.2.3"},`, 2000) + `{}]}`
	if status, _ := verifyTechnicalSBOMDeclaration(content, "widget@1.2.3"); status == domain.TechnicalVerified {
		t.Fatal("component budget bypassed")
	}
}

func TestTechnicalPureAPIBudgets(t *testing.T) {
	batch, bundle, _, asOf := verificationFixture()
	unique := func(count int) []domain.EvidenceBundle {
		out := make([]domain.EvidenceBundle, count)
		for i := range out {
			out[i] = bundle
			out[i].Repository.Repo = fmt.Sprintf("r%d", i)
			out[i].Items = nil
		}
		return out
	}
	tooManyItems := bundle
	tooManyItems.Items = make([]domain.CollectedEvidence, 101)
	large := bundle
	large.Items = []domain.CollectedEvidence{{Content: strings.Repeat("x", maxTechnicalFileBytes+1)}}
	total := bundle
	total.Items = make([]domain.CollectedEvidence, 9)
	for i := range total.Items {
		total.Items[i].Content = strings.Repeat("x", maxTechnicalFileBytes)
	}
	items := unique(11)
	for i := range items {
		items[i].Items = make([]domain.CollectedEvidence, 100)
	}
	for _, bundles := range [][]domain.EvidenceBundle{unique(101), {tooManyItems}, {large}, {total}, items} {
		if _, err := VerifyTechnicalEvidence(batch, bundles, asOf); err == nil {
			t.Fatal("pure API work budget bypassed")
		}
	}
	for _, alter := range []func(*domain.TechnicalEvidence){
		func(f *domain.TechnicalEvidence) { f.Limitations = make([]string, 101) },
		func(f *domain.TechnicalEvidence) { f.Limitations = []string{strings.Repeat("x", 8193)} },
		func(f *domain.TechnicalEvidence) { f.ObservedAt = strings.Repeat("x", 65) },
		func(f *domain.TechnicalEvidence) {
			f.Technique = domain.TechnicalEvidenceTechnique(strings.Repeat("x", 8193))
		},
	} {
		b := batch
		b.Findings = append([]domain.TechnicalEvidence(nil), batch.Findings...)
		alter(&b.Findings[0])
		if _, err := ReviewTechnicalEvidence(b, asOf); err == nil {
			t.Fatal("finding metadata budget bypassed")
		}
	}
}
