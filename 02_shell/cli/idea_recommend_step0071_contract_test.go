/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend_step0071_contract_test.md
 * @layer L2
 * @updated 2026-09-16
 */

package cli

import (
	"bytes"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestStep0071IdeaRecommendUsesProvidedCatalogEndToEnd(t *testing.T) {
	ideaPath, policyPath, source := recommendReplayFixture(t)
	rulesPath := filepath.Join(t.TempDir(), "rules.json")
	catalog := domain.CapabilityRuleCatalog{
		SchemaVersion: "2", Version: "provided-0071",
		Rules: []domain.CapabilityRule{{
			ID: "provided-dependency-graph", Basis: domain.EvidenceRuleDeclared,
			Family: "software-capability", Provenance: "step-0071-contract", TargetHints: []string{"source symbol"},
			CapabilityTerms:    []string{"analisar dependências", "dependency graph"},
			StrongEvidence:     []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"DependencyGraph"}}},
			SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceDocumentation, Patterns: []string{"dependency graph"}}},
		}},
	}
	if err := writeJSONAtomic(rulesPath, catalog); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "dossier")
	err := RunIdeaRecommendCommand([]string{
		"--idea", ideaPath, "--policy", policyPath, "--mode", "replay", "--strategy", "literal",
		"--replay", source, "--output", output, "--rules", rulesPath,
	}, &bytes.Buffer{}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var manifest domain.IdeaRunManifest
	if err := readJSONFile(filepath.Join(output, "run.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.RuleVersions["evidence"] != "provided-0071" || manifest.InputHashes["evidence_rules"] == "" {
		t.Fatalf("identidade causal não usa catálogo fornecido: %+v %+v", manifest.RuleVersions, manifest.InputHashes)
	}
	var assessment domain.AssessmentReport
	if err := readJSONFile(filepath.Join(output, "assessment-report.json"), &assessment); err != nil {
		t.Fatal(err)
	}
	if assessment.RulesVersion != "provided-0071" || assessment.CatalogCoverage.RulesVersion != "provided-0071" {
		t.Fatalf("assessment recuou ao catálogo embutido: %+v", assessment)
	}
	var coverage domain.EvidenceCatalogCoverage
	if err := readJSONFile(filepath.Join(output, "evidence", "catalog-coverage.json"), &coverage); err != nil {
		t.Fatal(err)
	}
	if coverage.RulesVersion != "provided-0071" || coverage.RulesHash == "" || len(coverage.Entries) == 0 {
		t.Fatalf("artifact de cobertura inválido: %+v", coverage)
	}
}
