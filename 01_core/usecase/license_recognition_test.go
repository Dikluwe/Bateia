/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/license_recognition_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"strings"
	"testing"
)

func TestLicenseRecognitionRequiresCompleteUnmodifiedGrant(t *testing.T) {
	for _, text := range []string{step0025MITLicense, canonicalMITTerms, "SPDX-License-Identifier: MIT\nCopyright (c) 2026 Someone\n" + canonicalMITTerms, strings.ReplaceAll(step0025MITLicense, "\n", "\r\n")} {
		if got := recognizedLicenseText(text); got != "MIT" {
			t.Fatalf("standard text not recognized: %q", got)
		}
	}
	for _, text := range []string{
		"MIT License\nPermission is hereby granted, free of charge",
		"This is not an MIT license grant. We quote: MIT License; Permission is hereby granted, free of charge.\nAll rights reserved. Redistribution prohibited.",
		"MIT License\nNot licensed for redistribution.\n" + canonicalMITTerms,
		"Copyright (c) 2026 Acme. Redistribution requires written consent.\n" + canonicalMITTerms,
		"Copyright (c) 2026 Acme Commercial use requires a separate paid agreement\n" + canonicalMITTerms,
		"Copyright (c) 2026 Acme The MIT text below is a quotation, not a license grant\n" + canonicalMITTerms,
		"Copyright (c) 2026 Acme Incorporated\n" + canonicalMITTerms,
		"Copyright (c) 2026 Usuário\n" + canonicalMITTerms,
		step0025MITLicense + "\nRedistribution prohibited.",
		strings.Replace(step0025MITLicense, "without restriction", "with restrictions", 1),
		"SPDX-License-Identifier: MIT", "Apache License\nVersion 2.0",
	} {
		b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: "LICENSE", Content: text})
		if got := assessLicense(b); got.Status == domain.EvidenceProven {
			t.Fatalf("unverified text promoted: %+v", got)
		}
	}
}

func TestLicenseAssessmentRejectsIncompleteAndAmbiguousRootDocuments(t *testing.T) {
	good := domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: "LICENSE", Content: step0025MITLicense}
	for _, b := range []domain.EvidenceBundle{
		testBundle(false, good),
		testBundle(true, good, domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: "COPYING", Content: step0025MITLicense}),
	} {
		if got := assessLicense(b); got.Status == domain.EvidenceProven {
			t.Fatalf("unsafe assessment: %+v", got)
		}
	}
}

func TestTechnicalLicenseDoesNotPromoteQuotedMarkers(t *testing.T) {
	batch, bundle, ref, asOf := verificationFixture()
	text := "This is not an MIT license grant. We quote: MIT License; Permission is hereby granted, free of charge.\nRedistribution prohibited."
	bundle.Items[0].Content = text
	bundle.Items[0].ContentHash = contentHash(text)
	batch.Findings[0].Source.Digest = "sha256:" + contentHash(text)
	review, err := VerifyTechnicalEvidence(batch, []domain.EvidenceBundle{bundle}, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if review.Verifications[0].Status == domain.TechnicalVerified {
		t.Fatal("quoted marker became proof")
	}
	if !review.Verifications[0].IntegrityVerified || review.Verifications[0].SourceAuthentication != "NOT_ESTABLISHED" || review.Verifications[0].BehaviorVerified {
		t.Fatalf("assurances conflated: %+v", review.Verifications[0])
	}
	if _, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: ref}}, review); changed {
		t.Fatal("license gate filled")
	}
}
