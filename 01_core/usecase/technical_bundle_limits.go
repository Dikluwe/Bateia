/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_bundle_limits.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"strings"
	"unicode"
)

const (
	maxTechnicalBundles     = 100
	maxTechnicalBundleItems = 100
	maxTechnicalItems       = 1000
	maxTechnicalFileBytes   = 2 << 20
	maxTechnicalBundleBytes = 16 << 20
)

// Bound traversal/hashing in the pure API, not only at CLI transport ingress.
// This is a work budget; invalid identity/digests remain verification outcomes.
func validateTechnicalBundleLimits(bundles []domain.EvidenceBundle) error {
	if len(bundles) > maxTechnicalBundles {
		return fmt.Errorf("verificação excede limite de 100 bundles")
	}
	items, total := 0, 0
	identities := map[domain.CommitRef]bool{}
	for _, bundle := range bundles {
		if identities[bundle.Repository] {
			return fmt.Errorf("bundle técnico duplicado para a mesma identidade e revisão")
		}
		identities[bundle.Repository] = true
		if len(bundle.Items) > maxTechnicalBundleItems {
			return fmt.Errorf("bundle técnico excede limite de 100 itens")
		}
		items += len(bundle.Items)
		if items > maxTechnicalItems {
			return fmt.Errorf("verificação excede limite total de 1000 itens")
		}
		for _, item := range bundle.Items {
			if len(item.Content) > maxTechnicalFileBytes {
				return fmt.Errorf("arquivo técnico excede limite de 2 MiB")
			}
			total += len(item.Content)
			if total > maxTechnicalBundleBytes {
				return fmt.Errorf("conteúdo técnico excede limite total de 16 MiB")
			}
			for _, field := range []string{item.Path, item.ContentHash, item.Ref.Owner, item.Ref.Repo, item.Ref.Commit, string(item.Kind)} {
				if len(field) > 8192 {
					return fmt.Errorf("metadado técnico excede limite de campo")
				}
			}
		}
		for _, field := range []string{bundle.Repository.Owner, bundle.Repository.Repo, bundle.Repository.Commit, bundle.Plan.Repository.Owner, bundle.Plan.Repository.Repo, bundle.Plan.Repository.Commit} {
			if len(field) > 8192 {
				return fmt.Errorf("identidade de bundle excede limite de campo")
			}
		}
	}
	return nil
}

type technicalPredicateIdentity struct {
	Candidate string
	Predicate domain.TechnicalEvidencePredicate
	Value     string
}

func technicalPredicateKey(f domain.TechnicalEvidence) technicalPredicateIdentity {
	value := f.Value
	if f.Predicate == domain.TechnicalLicenseID {
		// Match the verifier's EqualFold semantics, including Unicode aliases.
		value = canonicalEvidenceFold(value)
	}
	return technicalPredicateIdentity{Candidate: f.Candidate, Predicate: f.Predicate, Value: value}
}

func canonicalEvidenceFold(value string) string {
	var out strings.Builder
	for _, r := range value {
		minimum := r
		for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
			if next < minimum {
				minimum = next
			}
		}
		out.WriteRune(minimum)
	}
	return out.String()
}
