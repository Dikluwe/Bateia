/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_integrity_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package domain_test

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEvidencePopulatedSerializationCompatibility(t *testing.T) {
	b := domain.EvidenceBundle{SchemaVersion: "1", Repository: domain.CommitRef{Owner: "example", Repo: "project", Commit: strings.Repeat("a", 40)}, RequirementCoverage: []domain.RequirementCollectionCoverage{{RequirementKind: "document", Requirement: "root README", RuleID: "readme", Status: domain.RequirementTargetedExamined, CandidatePaths: []string{"README.md"}, ExaminedPaths: []string{"README.md"}}}}
	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	const want = "6fc91c3dfcf827ad4bcbf0820cc6fbe71fd82d2d58bac8f4c9919d5e4e81598f"
	if got != want {
		t.Fatalf("populated JSON changed: %s", got)
	}
}

func TestEvidenceCoverageEmptySerializationIsStableAndNeverNull(t *testing.T) {
	for _, empty := range []struct {
		name     string
		coverage []domain.RequirementCollectionCoverage
	}{{"nil", nil}, {"empty", []domain.RequirementCollectionCoverage{}}} {
		t.Run(empty.name, func(t *testing.T) {
			b := domain.EvidenceBundle{RequirementCoverage: empty.coverage}
			original := b
			first, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(original, b) {
				t.Fatal("marshal mutated caller")
			}
			for round := 0; round < 5; round++ {
				raw, err := json.Marshal(b)
				if err != nil {
					t.Fatal(err)
				}
				if strings.Contains(string(raw), `"requirement_coverage":null`) || !strings.Contains(string(raw), `"requirement_coverage":[]`) {
					t.Fatalf("coverage must be array: %s", raw)
				}
				if string(raw) != string(first) {
					t.Fatal("wire bytes changed after roundtrip")
				}
				if err := json.Unmarshal(raw, &b); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	var legacy domain.EvidenceBundle
	if err := json.Unmarshal([]byte(`{"requirement_coverage":null}`), &legacy); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(legacy)
	if err != nil || !strings.Contains(string(raw), `"requirement_coverage":[]`) {
		t.Fatal("null legacy input not canonicalized")
	}
}
