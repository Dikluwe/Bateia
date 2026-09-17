/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_step0072_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package domain_test

import (
	"encoding/json"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestStep0072CollectionCoveragePublicContract(t *testing.T) {
	if domain.RequirementTargetedExamined != "TARGETED_EXAMINED" ||
		domain.RequirementNoPlausiblePath != "NO_PLAUSIBLE_PATH" ||
		domain.RequirementCollectionIncomplete != "COLLECTION_INCOMPLETE" {
		t.Fatalf("valores públicos inesperados: %q %q %q", domain.RequirementTargetedExamined, domain.RequirementNoPlausiblePath, domain.RequirementCollectionIncomplete)
	}

	query := domain.EvidenceRequirementQuery{
		RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span",
		Terms: []string{"span", "basic_string"},
	}
	coverage := domain.RequirementCollectionCoverage{
		RequirementKind: "capability", Requirement: query.Requirement, RuleID: query.RuleID,
		Status: domain.RequirementCollectionIncomplete, SearchComplete: false,
		CandidatePaths: []string{"span"}, ExaminedPaths: []string{}, SkippedPaths: []string{"span"},
		ReasonCode: "COLLECTION_INCOMPLETE",
	}
	bundle := domain.EvidenceBundle{
		Plan: domain.EvidenceExtractionPlan{
			Budget:             domain.ExtractionBudget{CodeFiles: 2, RequirementFiles: 1},
			RequirementQueries: []domain.EvidenceRequirementQuery{query},
		},
		RequirementCoverage: []domain.RequirementCollectionCoverage{coverage},
	}
	encoded, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"requirement_files":1`, `"requirement_queries":[`, `"requirement_coverage":[`, `"candidate_paths":[`, `"examined_paths":[]`, `"skipped_paths":[`, `"reason_code":"COLLECTION_INCOMPLETE"`} {
		if !strings.Contains(string(encoded), field) {
			t.Fatalf("campo público ausente %s em %s", field, encoded)
		}
	}
	var roundTrip domain.EvidenceBundle
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if len(roundTrip.Plan.RequirementQueries) != 1 || len(roundTrip.RequirementCoverage) != 1 || roundTrip.RequirementCoverage[0].Status != domain.RequirementCollectionIncomplete {
		t.Fatalf("round-trip perdeu diagnóstico: %+v", roundTrip)
	}
}

func TestStep0072EmptyPublicArraysAreNeverNull(t *testing.T) {
	bundle := domain.EvidenceBundle{
		Plan: domain.EvidenceExtractionPlan{RequirementQueries: []domain.EvidenceRequirementQuery{}},
		RequirementCoverage: []domain.RequirementCollectionCoverage{{
			CandidatePaths: []string{}, ExaminedPaths: []string{}, SkippedPaths: []string{},
		}},
	}
	encoded, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{`"requirement_queries":null`, `"requirement_coverage":null`, `"candidate_paths":null`, `"examined_paths":null`, `"skipped_paths":null`} {
		if strings.Contains(string(encoded), forbidden) {
			t.Fatalf("array público serializado como null: %s", encoded)
		}
	}
}
