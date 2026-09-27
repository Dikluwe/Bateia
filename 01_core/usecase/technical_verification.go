/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_verification.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"sort"
	"strings"
	"time"

	"bateia/01_core/domain"
)

// VerifyTechnicalEvidence binds structured findings to collected bytes at a
// pinned repository revision. It verifies internal bundle integrity, not the
// authenticity of an arbitrary JSON bundle's origin.
func VerifyTechnicalEvidence(batch domain.TechnicalEvidenceBatch, bundles []domain.EvidenceBundle, asOf time.Time) (domain.TechnicalEvidenceReview, error) {
	if err := validateTechnicalBundleLimits(bundles); err != nil {
		return domain.TechnicalEvidenceReview{}, err
	}
	review, err := ReviewTechnicalEvidence(batch, asOf)
	if err != nil {
		return domain.TechnicalEvidenceReview{}, err
	}
	claimStatus := map[technicalClaimIdentity]domain.TechnicalClaimStatus{}
	for _, claim := range review.Claims {
		claimStatus[technicalClaimIdentity{Candidate: claim.Candidate, Claim: claim.Claim}] = claim.Status
	}
	structuredResults := map[technicalPredicateIdentity]uint8{}
	for _, finding := range batch.Findings {
		if finding.Predicate == "" {
			continue
		}
		if finding.ValidUntil != "" {
			until, _ := time.Parse(time.RFC3339, finding.ValidUntil)
			if asOf.After(until) {
				continue
			}
		}
		key := technicalPredicateKey(finding)
		if finding.Result == domain.TechnicalSupports {
			structuredResults[key] |= 1
		} else if finding.Result == domain.TechnicalRefutes {
			structuredResults[key] |= 2
		}
	}
	findings := append([]domain.TechnicalEvidence(nil), batch.Findings...)
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	for _, finding := range findings {
		v := domain.TechnicalEvidenceVerification{SourceAuthentication: "NOT_ESTABLISHED", FindingID: finding.ID, Candidate: finding.Candidate, Predicate: finding.Predicate, Value: finding.Value, Status: domain.TechnicalUnverified}
		if finding.ValidUntil != "" {
			until, _ := time.Parse(time.RFC3339, finding.ValidUntil)
			if asOf.After(until) {
				v.Status, v.Message = domain.TechnicalStale, "achado vencido; nova coleta necessária"
				review.Verifications = append(review.Verifications, v)
				continue
			}
		}
		structuredKey := technicalPredicateKey(finding)
		status := claimStatus[technicalClaimIdentity{Candidate: finding.Candidate, Claim: finding.Claim}]
		if status == domain.TechnicalClaimConflicting || structuredResults[structuredKey] == 3 {
			v.Status, v.Message = domain.TechnicalContradicted, "alegação possui suporte e refutação; requer resolução"
		} else if finding.Result != domain.TechnicalSupports || status != domain.TechnicalClaimSupported {
			v.Message = "somente achado atual de suporte pode ser verificado como fonte positiva"
		} else {
			v = verifyStructuredFinding(finding, bundles)
		}
		review.Verifications = append(review.Verifications, v)
	}
	return review, nil
}

func verifyStructuredFinding(f domain.TechnicalEvidence, bundles []domain.EvidenceBundle) domain.TechnicalEvidenceVerification {
	v := domain.TechnicalEvidenceVerification{SourceAuthentication: "NOT_ESTABLISHED", FindingID: f.ID, Candidate: f.Candidate, Predicate: f.Predicate, Value: f.Value, Status: domain.TechnicalUnverified}
	if (f.Predicate == domain.TechnicalLicenseID && f.Technique != domain.TechniqueLicense) || (f.Predicate == domain.TechnicalSBOMComponent && f.Technique != domain.TechniqueSBOM) {
		v.Message = "predicado não corresponde à técnica declarada"
		return v
	}
	if f.Predicate != domain.TechnicalLicenseID && f.Predicate != domain.TechnicalSBOMComponent {
		v.Status, v.Message = domain.TechnicalUnsupported, "técnica ainda sem verificador determinístico"
		return v
	}
	if strings.TrimSpace(f.Value) == "" || strings.TrimSpace(f.Source.Revision) == "" {
		v.Message = "verificação exige value e source.revision fixada"
		return v
	}
	digest, valid := evidenceDigest(f.Source.Digest)
	if !valid {
		v.Message = "verificação exige source.digest SHA-256 válido"
		return v
	}
	var matches []domain.CollectedEvidence
	var matchedBundle *domain.EvidenceBundle
	for i := range bundles {
		bundle := &bundles[i]
		if bundle.SchemaVersion != "1" || bundle.Repository.Key() != f.Candidate || bundle.Repository.Commit != f.Source.Revision || bundle.Plan.Repository != bundle.Repository || !bundle.Complete || len(bundle.Failures) != 0 {
			continue
		}
		for _, item := range bundle.Items {
			if item.Path == f.Source.Ref {
				matches = append(matches, item)
				matchedBundle = bundle
			}
		}
	}
	if len(matches) != 1 || matchedBundle == nil {
		v.Message = "fonte ausente, duplicada ou bundle incompleto na revisão fixada"
		return v
	}
	item := matches[0]
	if item.Ref != matchedBundle.Repository || !validMechanismPath(item.Path) || item.ContentHash != contentHash(item.Content) || digest != contentHash(item.Content) {
		v.Message = "identidade ou digest dos bytes coletados não confere"
		return v
	}
	v.IntegrityVerified = true
	switch f.Predicate {
	case domain.TechnicalLicenseID:
		if item.Kind != domain.EvidenceLicense || !authoritativeLicensePath(item.Path) {
			v.Message = "identificação de licença exige documento raiz de licença"
			return v
		}
		license := assessLicense(*matchedBundle)
		if license.Status != domain.EvidenceProven || !strings.EqualFold(license.Identifier, f.Value) || !EvidenceSourcesAdmissible(item.Ref, license.Sources, false) {
			v.Message = "identificador da licença não confere com o documento"
			return v
		}
		v.InterpretationBasis = license.RecognitionBasis
		v.Status, v.Source, v.Message = domain.TechnicalVerified, &license.Sources[0], "identificador de licença verificado no documento raiz; compatibilidade jurídica não inferida"
	case domain.TechnicalSBOMComponent:
		if item.Kind != domain.EvidenceManifest || !strings.HasSuffix(strings.ToLower(item.Path), ".json") {
			v.Message = "declaração SBOM exige manifesto JSON"
			return v
		}
		status, message := verifyTechnicalSBOMDeclaration(item.Content, f.Value)
		if status != domain.TechnicalVerified {
			v.Status, v.Message = status, message
			return v
		}
		first := truncateSnippet(strings.Split(item.Content, "\n")[0])
		source := sourceFromItem(item, domain.EvidenceManifest, "technical-sbom-declaration", .5, first)
		if !EvidenceSourcesAdmissible(item.Ref, []domain.EvidenceSource{source}, false) {
			v.Message = "fonte SBOM não admissível"
			return v
		}
		v.InterpretationBasis = "DECLARATION_ONLY"
		v.Status, v.Source, v.Message = domain.TechnicalVerified, &source, message
	}
	return v
}

// ApplyVerifiedTechnicalLicenses only fills an inconclusive license gate. It
// does not replace an already proven assessment or infer license compatibility.
func ApplyVerifiedTechnicalLicenses(assessments []domain.CandidateAssessment, review domain.TechnicalEvidenceReview) ([]domain.CandidateAssessment, bool) {
	out := append([]domain.CandidateAssessment(nil), assessments...)
	byCandidate := map[string][]domain.TechnicalEvidenceVerification{}
	for _, v := range review.Verifications {
		if v.Status == domain.TechnicalVerified && v.IntegrityVerified && v.InterpretationBasis == "FULL_TEXT_TEMPLATE" && v.Predicate == domain.TechnicalLicenseID && v.Source != nil {
			byCandidate[v.Candidate] = append(byCandidate[v.Candidate], v)
		}
	}
	changed := false
	for i, assessment := range out {
		verified := byCandidate[assessment.Repository.Key()]
		if len(verified) == 0 || SanitizeAssessmentEvidence(assessment).License.Status != domain.EvidenceInconclusive {
			continue
		}
		identifier := verified[0].Value
		consistent := true
		for _, v := range verified {
			consistent = consistent && strings.EqualFold(v.Value, identifier) && v.Source.Ref == assessment.Repository
		}
		if !consistent {
			continue
		}
		candidate := assessment
		candidate.License = domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: identifier, RecognitionBasis: verified[0].InterpretationBasis, Sources: []domain.EvidenceSource{*verified[0].Source}, Caveat: "texto completo reconhecido; compatibilidade e escopo jurídico não inferidos"}
		if SanitizeAssessmentEvidence(candidate).License.Status != domain.EvidenceProven {
			continue
		}
		out[i] = candidate
		changed = true
	}
	return out, changed
}
