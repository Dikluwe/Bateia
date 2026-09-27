/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_bundle_limits_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"strings"
	"testing"
)

func TestTechnicalBundleLimitsEnforced(t *testing.T) {
	// Excesso de bundles (> 100)
	tooManyBundles := make([]domain.EvidenceBundle, 101)
	if err := validateTechnicalBundleLimits(tooManyBundles); err == nil || !strings.Contains(err.Error(), "100 bundles") {
		t.Fatalf("expected >100 bundles error, got %v", err)
	}

	// Excesso de itens no bundle (> 100)
	tooManyItems := []domain.EvidenceBundle{{
		Items: make([]domain.CollectedEvidence, 101),
	}}
	if err := validateTechnicalBundleLimits(tooManyItems); err == nil || !strings.Contains(err.Error(), "100 itens") {
		t.Fatalf("expected >100 items error, got %v", err)
	}

	// Arquivo excessivamente grande (> 2 MiB)
	hugeFile := []domain.EvidenceBundle{{
		Items: []domain.CollectedEvidence{{
			Content: strings.Repeat("x", (2<<20)+1),
		}},
	}}
	if err := validateTechnicalBundleLimits(hugeFile); err == nil || !strings.Contains(err.Error(), "2 MiB") {
		t.Fatalf("expected >2 MiB error, got %v", err)
	}

	// Bundle válido dentro dos orçamentos
	valid := []domain.EvidenceBundle{{
		Items: []domain.CollectedEvidence{{
			Path:    "go.mod",
			Content: "module test\ngo 1.25",
		}},
	}}
	if err := validateTechnicalBundleLimits(valid); err != nil {
		t.Fatalf("valid bundle rejected: %v", err)
	}
}
