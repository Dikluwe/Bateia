/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_step0071_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStep0071CatalogCoveragePublicSchema(t *testing.T) {
	if EvidenceRuleDeclared != "DECLARED" || EvidenceRuleHypothesis != "HYPOTHESIS" {
		t.Fatalf("bases instáveis: %q %q", EvidenceRuleDeclared, EvidenceRuleHypothesis)
	}
	if CatalogCoverageCovered != "COVERED" || CatalogCoverageGap != "CATALOG_GAP" || CatalogCoverageHypothesis != "HYPOTHESIS" {
		t.Fatalf("estados instáveis: %q %q %q", CatalogCoverageCovered, CatalogCoverageGap, CatalogCoverageHypothesis)
	}

	report := EvidenceCatalogCoverage{
		SchemaVersion: "1", RulesVersion: "custom-0071", RulesHash: strings.Repeat("a", 64),
		Entries: []CatalogCoverageEntry{{
			RequirementKind: "capability", Requirement: "quantum flux", Status: CatalogCoverageGap,
			ReasonCode: "CATALOG_GAP",
		}},
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, witness := range []string{`"schema_version":"1"`, `"rules_version":"custom-0071"`, `"status":"CATALOG_GAP"`, `"reason_code":"CATALOG_GAP"`} {
		if !strings.Contains(string(encoded), witness) {
			t.Fatalf("schema perdeu %s: %s", witness, encoded)
		}
	}
	var decoded EvidenceCatalogCoverage
	if err := json.Unmarshal(encoded, &decoded); err != nil || len(decoded.Entries) != 1 || decoded.Entries[0].Status != CatalogCoverageGap {
		t.Fatalf("round-trip inválido: err=%v report=%+v", err, decoded)
	}

	empty, err := json.Marshal(EvidenceCatalogCoverage{Entries: []CatalogCoverageEntry{}})
	if err != nil || !strings.Contains(string(empty), `"entries":[]`) {
		t.Fatalf("array público deve ser não nulo: err=%v json=%s", err, empty)
	}
	ruleJSON, _ := json.Marshal(CapabilityRule{ID: "hyp", Basis: EvidenceRuleHypothesis})
	if !strings.Contains(string(ruleJSON), `"basis":"HYPOTHESIS"`) {
		t.Fatalf("hipótese não explicitada: %s", ruleJSON)
	}
	evidenceJSON, _ := json.Marshal(CapabilityEvidence{Status: EvidenceInconclusive, ReasonCode: "CATALOG_GAP"})
	if !strings.Contains(string(evidenceJSON), `"reason_code":"CATALOG_GAP"`) {
		t.Fatalf("causa de Unknown perdida: %s", evidenceJSON)
	}
}
