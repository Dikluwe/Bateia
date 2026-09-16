/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0071_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func step0071Rule(id string, basis domain.EvidenceRuleBasis, terms, patterns []string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, Basis: basis, Family: "software-capability", Provenance: "step-0071-contract",
		TargetHints: []string{"source symbol"}, CapabilityTerms: terms,
		StrongEvidence:     []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: patterns}},
		SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceDocumentation, Patterns: patterns}},
	}
}

func step0071Catalog(rules ...domain.CapabilityRule) domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "custom-0071", Rules: rules}
}

func step0071Bundle(items ...domain.CollectedEvidence) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "org", Repo: "candidate", Commit: "abc123"}
	for i := range items {
		items[i].Ref = ref
		sum := sha256.Sum256([]byte(items[i].Content))
		items[i].ContentHash = hex.EncodeToString(sum[:])
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: items, Targets: []domain.CollectionTargetResult{}}
}

func step0071CoverageEntry(t *testing.T, report domain.EvidenceCatalogCoverage, requirement string) domain.CatalogCoverageEntry {
	t.Helper()
	for _, entry := range report.Entries {
		if entry.Requirement == requirement {
			return entry
		}
	}
	t.Fatalf("cobertura ausente para %q: %+v", requirement, report)
	return domain.CatalogCoverageEntry{}
}

func step0071Capability(t *testing.T, assessment domain.CandidateAssessment, capability string) domain.CapabilityEvidence {
	t.Helper()
	for _, evidence := range assessment.Capabilities {
		if evidence.Capability == capability {
			return evidence
		}
	}
	t.Fatalf("assessment ausente para %q: %+v", capability, assessment)
	return domain.CapabilityEvidence{}
}

func TestStep0071CatalogGapIsUnknownNotCandidateAbsence(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "catalog routing", Capabilities: []string{"quantum flux"}}
	catalog := step0071Catalog(step0071Rule("other", domain.EvidenceRuleDeclared, []string{"other capability"}, []string{"OtherEngine"}))
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		t.Fatal(err)
	}
	entry := step0071CoverageEntry(t, coverage, "quantum flux")
	if entry.Status != domain.CatalogCoverageGap || entry.ReasonCode != "CATALOG_GAP" {
		t.Fatalf("lacuna não preservada: %+v", entry)
	}

	bundle := step0071Bundle(domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "engine.go", Content: "type QuantumFlux struct{}"})
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	got := step0071Capability(t, report.Assessments[0], "quantum flux")
	if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "CATALOG_GAP" || len(got.Sources) != 0 {
		t.Fatalf("catálogo foi confundido com candidato: %+v", got)
	}
	if len(report.Assessments[0].Gaps) != 0 {
		t.Fatalf("lacuna do catálogo atribuída ao candidato: %+v", report.Assessments[0].Gaps)
	}
	plan := usecase.BuildInvestigationPlan(idea, domain.DiscoverySnapshot{}, report.Assessments)
	for _, action := range plan.Actions {
		if action.Classification == "NOT_VERIFIED" {
			t.Fatalf("CATALOG_GAP degradado para NOT_VERIFIED: %+v", action)
		}
	}
}

func TestStep0071UnsafeFragmentAndHypothesisRemainUnknown(t *testing.T) {
	fragmentIdea := domain.IdeaSpec{Problem: "routing", Capabilities: []string{"graph"}}
	declared := step0071Catalog(step0071Rule("dependency-graph", domain.EvidenceRuleDeclared, []string{"dependency graph"}, []string{"DependencyGraph"}))
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(fragmentIdea, declared)
	if err != nil {
		t.Fatal(err)
	}
	if got := step0071CoverageEntry(t, coverage, "graph"); got.Status != domain.CatalogCoverageGap {
		t.Fatalf("fragmento reverso casou regra insegura: %+v", got)
	}

	hypothesis := step0071Catalog(step0071Rule("flux-hypothesis", domain.EvidenceRuleHypothesis, []string{"quantum flux"}, []string{"FluxEngine"}))
	idea := domain.IdeaSpec{Problem: "quantum flux", Capabilities: []string{"quantum flux"}}
	coverage, err = usecase.AnalyzeEvidenceCatalogCoverage(idea, hypothesis)
	if err != nil {
		t.Fatal(err)
	}
	if got := step0071CoverageEntry(t, coverage, "quantum flux"); got.Status != domain.CatalogCoverageHypothesis || got.RuleBasis != domain.EvidenceRuleHypothesis {
		t.Fatalf("hipótese não ficou explícita: %+v", got)
	}
	bundle := step0071Bundle(domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "flux.go", Content: "type FluxEngine struct{}"})
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, hypothesis)
	if err != nil {
		t.Fatal(err)
	}
	got := step0071Capability(t, report.Assessments[0], "quantum flux")
	if got.Status == domain.EvidenceProven || got.Status == domain.EvidenceAbsent || got.ReasonCode != "HYPOTHESIS_ONLY" || len(report.Assessments[0].ReusableParts) != 0 {
		t.Fatalf("hipótese lexical ganhou autoridade: %+v parts=%+v", got, report.Assessments[0].ReusableParts)
	}
}

func TestStep0071CoverageAndAssessmentAreCanonicalAcrossReordering(t *testing.T) {
	rules := []domain.CapabilityRule{
		step0071Rule("alpha", domain.EvidenceRuleDeclared, []string{"alpha capability", "alpha alias"}, []string{"AlphaEngine"}),
		step0071Rule("beta", domain.EvidenceRuleDeclared, []string{"beta capability", "beta alias"}, []string{"BetaEngine"}),
	}
	items := []domain.CollectedEvidence{
		{Kind: domain.EvidenceFile, Path: "alpha.go", Content: "type AlphaEngine struct{}"},
		{Kind: domain.EvidenceFile, Path: "beta.go", Content: "type BetaEngine struct{}"},
	}
	var canonicalCoverage, canonicalAssessment []byte
	for iteration := 0; iteration < 20; iteration++ {
		idea := domain.IdeaSpec{Problem: "alpha capability", Capabilities: []string{"alpha capability", "beta capability"}}
		catalogRules := append([]domain.CapabilityRule(nil), rules...)
		bundleItems := append([]domain.CollectedEvidence(nil), items...)
		if iteration%2 == 1 {
			idea.Capabilities[0], idea.Capabilities[1] = idea.Capabilities[1], idea.Capabilities[0]
			catalogRules[0], catalogRules[1] = catalogRules[1], catalogRules[0]
			bundleItems[0], bundleItems[1] = bundleItems[1], bundleItems[0]
		}
		coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, step0071Catalog(catalogRules...))
		if err != nil {
			t.Fatal(err)
		}
		report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{step0071Bundle(bundleItems...)}, step0071Catalog(catalogRules...))
		if err != nil {
			t.Fatal(err)
		}
		coverageJSON, _ := json.Marshal(coverage)
		assessmentJSON, _ := json.Marshal(report)
		if iteration == 0 {
			canonicalCoverage, canonicalAssessment = coverageJSON, assessmentJSON
			continue
		}
		if string(coverageJSON) != string(canonicalCoverage) || string(assessmentJSON) != string(canonicalAssessment) {
			t.Fatalf("reordenação %d alterou saída\ncoverage=%s\nassessment=%s", iteration, coverageJSON, assessmentJSON)
		}
	}
}
