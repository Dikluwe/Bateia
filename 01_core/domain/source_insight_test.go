/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/source_insight_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestSourceInsightOptionalObservationsPreserveLegacyBytes(t *testing.T) {
	source := InsightEvidenceInput{}
	absent, _ := json.Marshal(source)
	source.Observations = []SourceInsightObservation{}
	empty, _ := json.Marshal(source)
	if !bytes.Equal(absent, empty) || bytes.Contains(empty, []byte("observations")) {
		t.Fatal("optional observations changed legacy preimage")
	}
}

func TestSourceInsightReportRoundTripPreservesDeclaredOriginAndScope(t *testing.T) {
	report := SourceInsightInvestigationReport{SchemaVersion: "insight-source-analysis-v2", Status: ImprovementPartial,
		SourceAnalyses: []SourceInsightAnalysis{{
			CatalogDefinition: InsightCatalogDefinition{Origin: "CATALOG_TEMPLATE", Strategy: "generic"},
			ExtractionScope:   InsightExtractionScope{Languages: []string{"go"}},
			MechanismStatus:   MechanismPartial, AssessmentOrigin: "DECLARED_NOT_REASSESSED",
			Observations: []LocatedInsightObservation{{Kind: InsightPrecondition, Statement: "Configure the limit", Origin: "DECLARED_NOT_REASSESSED", References: []LocatedInsightReference{{InsightObservationReference: InsightObservationReference{Path: "pool.go", ContentHash: "sha256:fixture", StartLine: 1, EndLine: 1}, Excerpt: "setLimit(1)"}}}},
		}}, Hypotheses: []ImprovementHypothesis{{TargetState: TargetUnknown, Preconditions: []string{"Configure the limit"}, AdoptionDisposition: "NOT_AN_ADOPTION_DECISION"}}}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var again SourceInsightInvestigationReport
	if err := json.Unmarshal(encoded, &again); err != nil || !reflect.DeepEqual(report, again) {
		t.Fatalf("wire contract lost observation, scope or state: %v", err)
	}
}
