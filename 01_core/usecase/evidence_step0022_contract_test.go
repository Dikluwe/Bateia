/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0022_contract_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"testing"
)

func TestStep0022CatalogV2RequiresProvenanceAndRejectsIdentity(t *testing.T) {
	valid := DefaultEvidenceRules()
	if err := ValidateRuleCatalog(valid); err != nil {
		t.Fatalf("catálogo padrão inválido: %v", err)
	}

	missing := valid
	missing.Hash = ""
	missing.Rules = append([]domain.CapabilityRule(nil), valid.Rules...)
	missing.Rules[0].Provenance = ""
	if ValidateRuleCatalog(missing) == nil {
		t.Fatal("catálogo v2 aceitou regra sem proveniência")
	}

	identity := valid
	identity.Hash = ""
	identity.Rules = append([]domain.CapabilityRule(nil), valid.Rules...)
	identity.Rules[0].StrongEvidence = []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"github.com/acme/product"}}}
	if ValidateRuleCatalog(identity) == nil {
		t.Fatal("catálogo aceitou padrão dependente de identidade")
	}
}

func TestStep0022CatalogRejectsCrossRuleDuplicatesAndGenericPatterns(t *testing.T) {
	catalog := DefaultEvidenceRules()
	catalog.Hash = ""
	catalog.Rules = append([]domain.CapabilityRule(nil), catalog.Rules...)
	catalog.Rules[1].CapabilityTerms = append(catalog.Rules[1].CapabilityTerms, catalog.Rules[0].CapabilityTerms[0])
	if ValidateRuleCatalog(catalog) == nil {
		t.Fatal("catálogo aceitou termo duplicado entre regras")
	}

	catalog = DefaultEvidenceRules()
	catalog.Hash = ""
	catalog.Rules = append([]domain.CapabilityRule(nil), catalog.Rules...)
	catalog.Rules[0].StrongEvidence = []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"manager"}}}
	if ValidateRuleCatalog(catalog) == nil {
		t.Fatal("catálogo aceitou padrão excessivamente genérico")
	}
}

func TestStep0040IndependentLexicalSourcesRemainPartial(t *testing.T) {
	rule := findRule("compare software candidates", DefaultEvidenceRules())
	one := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "comparison.go", Content: "type CandidateComparison struct{}"})
	if got := assessCapability("compare software candidates", one, rule); got.Status == domain.EvidenceProven {
		t.Fatalf("fonte única produziu PROVEN: %#v", got)
	}

	two := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "comparison.go", Content: "type CandidateComparison struct{}"},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "comparison_test.go", Content: "func TestCandidateComparison(t *testing.T) {}"},
	)
	if got := assessCapability("compare software candidates", two, rule); got.Status != domain.EvidencePartial {
		t.Fatalf("fontes lexicais deveriam permanecer PARTIAL: %#v", got)
	}
}

func TestStep0022DuplicateContentDoesNotCreateIndependence(t *testing.T) {
	content := "type CandidateComparison struct{}"
	bundle := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "a.go", Content: content},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "copy.go", Content: content},
	)
	if got := assessCapability("compare software candidates", bundle, findRule("compare software candidates", DefaultEvidenceRules())); got.Status == domain.EvidenceProven {
		t.Fatalf("conteúdo duplicado fabricou independência: %#v", got)
	}
}

func TestStep0022RejectsMutableAndCrossRepositorySources(t *testing.T) {
	idea := domain.IdeaSpec{Capabilities: []string{"compare software candidates"}}
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo@abc", Decision: "ADVANCE"}}}
	mutable := testBundle(true)
	mutable.Repository.Commit = "HEAD"
	if _, err := AssessEvidence(idea, screening, []domain.EvidenceBundle{mutable}, DefaultEvidenceRules()); err == nil {
		t.Fatal("bundle em HEAD foi aceito")
	}

	cross := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "comparison.go", Content: "type CandidateComparison struct{}"})
	cross.Items[0].Ref.Repo = "other"
	if _, err := AssessEvidence(idea, screening, []domain.EvidenceBundle{cross}, DefaultEvidenceRules()); err == nil {
		t.Fatal("fonte de outro repositório foi aceita")
	}
}

func TestStep0022ExplicitNegativeConflictsAndRemovedSymbolDoesNotProve(t *testing.T) {
	rule := findRule("compare software candidates", DefaultEvidenceRules())
	negative := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "candidate comparison is not implemented"})
	if got := assessCapability("compare software candidates", negative, rule); got.Status != domain.EvidenceConflict {
		t.Fatalf("negação explícita não produziu conflito: %#v", got)
	}
	removed := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "comparison.go", Content: "// CandidateComparison was removed"})
	if got := assessCapability("compare software candidates", removed, rule); got.Status == domain.EvidenceProven {
		t.Fatalf("símbolo removido produziu PROVEN: %#v", got)
	}
}
