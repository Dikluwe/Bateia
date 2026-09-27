/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_evidence.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"bateia/01_core/domain"
)

type technicalClaimIdentity struct {
	Candidate string
	Claim     string
}

// ReviewTechnicalEvidence reconciles optional specialized findings. Its output
// is advisory: it never changes the deterministic decision or EvidenceSource.
func ReviewTechnicalEvidence(batch domain.TechnicalEvidenceBatch, asOf time.Time) (domain.TechnicalEvidenceReview, error) {
	if len(batch.Findings) > 100 {
		return domain.TechnicalEvidenceReview{}, fmt.Errorf("lote técnico excede limite de 100 achados")
	}
	if batch.SchemaVersion != "1" {
		return domain.TechnicalEvidenceReview{}, fmt.Errorf("technical_evidence.schema_version deve ser 1")
	}
	asOf = asOf.UTC()
	review := domain.TechnicalEvidenceReview{ContentPolicy: "UNTRUSTED_DATA_NO_AUTHORITY", SchemaVersion: "1", AsOf: asOf.Format(time.RFC3339), DecisionImpact: "ADVISORY_ONLY", Claims: []domain.TechnicalClaimReview{}}
	groups := map[technicalClaimIdentity][]domain.TechnicalEvidence{}
	ids := map[string]bool{}
	for _, finding := range batch.Findings {
		if err := validateTechnicalFinding(finding, asOf); err != nil {
			return domain.TechnicalEvidenceReview{}, err
		}
		if ids[finding.ID] {
			return domain.TechnicalEvidenceReview{}, fmt.Errorf("technical_evidence.id duplicado: %s", finding.ID)
		}
		ids[finding.ID] = true
		key := technicalClaimIdentity{Candidate: finding.Candidate, Claim: finding.Claim}
		groups[key] = append(groups[key], finding)
	}
	keys := make([]technicalClaimIdentity, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Candidate != keys[j].Candidate {
			return keys[i].Candidate < keys[j].Candidate
		}
		return keys[i].Claim < keys[j].Claim
	})
	for _, key := range keys {
		findings := groups[key]
		claim := domain.TechnicalClaimReview{Candidate: findings[0].Candidate, Claim: findings[0].Claim, FindingIDs: []string{}}
		var supports, refutes bool
		limitations := map[string]bool{}
		for _, finding := range findings {
			claim.FindingIDs = append(claim.FindingIDs, finding.ID)
			for _, limitation := range finding.Limitations {
				if limitation = strings.TrimSpace(limitation); limitation != "" {
					limitations[limitation] = true
				}
			}
			if finding.ValidUntil != "" {
				until, _ := time.Parse(time.RFC3339, finding.ValidUntil)
				if asOf.After(until) {
					claim.StaleIDs = append(claim.StaleIDs, finding.ID)
					continue
				}
			}
			supports = supports || finding.Result == domain.TechnicalSupports
			refutes = refutes || finding.Result == domain.TechnicalRefutes
		}
		sort.Strings(claim.FindingIDs)
		sort.Strings(claim.StaleIDs)
		for limitation := range limitations {
			claim.Limitations = append(claim.Limitations, limitation)
		}
		sort.Strings(claim.Limitations)
		switch {
		case supports && refutes:
			claim.Status = domain.TechnicalClaimConflicting
			review.Warnings = append(review.Warnings, fmt.Sprintf("%s: %s tem evidências contraditórias", claim.Candidate, claim.Claim))
		case supports:
			claim.Status = domain.TechnicalClaimSupported
		case refutes:
			claim.Status = domain.TechnicalClaimRefuted
		case len(claim.StaleIDs) == len(claim.FindingIDs):
			claim.Status = domain.TechnicalClaimStale
		default:
			claim.Status = domain.TechnicalClaimInconclusive
		}
		if len(claim.StaleIDs) > 0 {
			review.Warnings = append(review.Warnings, fmt.Sprintf("%s: %s contém %d achado(s) vencido(s)", claim.Candidate, claim.Claim, len(claim.StaleIDs)))
		}
		review.Claims = append(review.Claims, claim)
	}
	return review, nil
}

func validateTechnicalFinding(f domain.TechnicalEvidence, asOf time.Time) error {
	if len(f.Limitations) > 100 || len(f.ObservedAt) > 64 || len(f.ValidUntil) > 64 {
		return fmt.Errorf("achado técnico excede orçamento de limitações ou datas")
	}
	for _, limitation := range f.Limitations {
		if len(limitation) > 8192 {
			return fmt.Errorf("limitação técnica excede tamanho máximo")
		}
	}
	for _, value := range []string{f.ID, f.Candidate, f.Claim, f.Source.Ref, f.Source.Revision, f.Source.Digest, f.Method, f.Value, string(f.Technique), string(f.Result), string(f.Predicate)} {
		if len(value) > 8192 {
			return fmt.Errorf("achado técnico excede limite de tamanho de campo")
		}
	}
	if strings.TrimSpace(f.ID) == "" || strings.TrimSpace(f.Candidate) == "" || strings.TrimSpace(f.Claim) == "" || strings.TrimSpace(f.Source.Ref) == "" || strings.TrimSpace(f.Method) == "" {
		return fmt.Errorf("achado técnico exige id, candidate, claim, source.ref e method")
	}
	if strings.TrimSpace(f.Source.Revision) == "" && strings.TrimSpace(f.Source.Digest) == "" {
		return fmt.Errorf("achado %s: source exige revision ou digest para reproduzir a origem", f.ID)
	}
	switch f.Technique {
	case domain.TechniqueCodeSearch, domain.TechniqueRepositoryMining, domain.TechniqueAPIRecommendation, domain.TechniqueSoftwareArchaeology, domain.TechniqueArchitectureRecovery, domain.TechniqueImpactAnalysis, domain.TechniqueSBOM, domain.TechniqueLicense, domain.TechniqueProjectHealth, domain.TechniqueBuildProvenance:
	default:
		return fmt.Errorf("achado %s: technique desconhecida %q", f.ID, f.Technique)
	}
	switch f.Result {
	case domain.TechnicalSupports, domain.TechnicalRefutes, domain.TechnicalInconclusive:
	default:
		return fmt.Errorf("achado %s: result inválido %q", f.ID, f.Result)
	}
	observed, err := time.Parse(time.RFC3339, f.ObservedAt)
	if err != nil {
		return fmt.Errorf("achado %s: observed_at exige RFC3339: %w", f.ID, err)
	}
	if observed.After(asOf) {
		return fmt.Errorf("achado %s: observed_at está no futuro", f.ID)
	}
	if f.ValidUntil != "" {
		until, err := time.Parse(time.RFC3339, f.ValidUntil)
		if err != nil {
			return fmt.Errorf("achado %s: valid_until exige RFC3339: %w", f.ID, err)
		}
		if until.Before(observed) {
			return fmt.Errorf("achado %s: valid_until anterior a observed_at", f.ID)
		}
	}
	return nil
}
