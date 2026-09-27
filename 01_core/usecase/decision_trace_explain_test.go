/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_explain_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"strings"
	"testing"
)

func TestDecisionTraceExplainPointersAndTruncation(t *testing.T) {
	buildInput := domain.DecisionTraceBuildInput{
		SchemaVersion: "1",
		Report: domain.ReuseDecisionReport{
			Decision: domain.DecisionBuild,
			Reasons: []domain.DecisionReason{
				{
					RuleID:    "RULE-LONG",
					Candidate: "cand-1",
					Effect:    "BLOCK",
					Message:   strings.Repeat("m", 2000),
				},
			},
		},
	}
	built, err := usecase.BuildDecisionTrace(buildInput)
	if err != nil {
		t.Fatal(err)
	}

	explanation, err := usecase.ExplainDecisionTrace(domain.DecisionTraceExplainInput{
		SchemaVersion: "1",
		Graph:         built,
		BuildInput:    buildInput,
	})
	if err != nil {
		t.Fatal(err)
	}
	if explanation.SchemaVersion != domain.DecisionTraceExplanationSchemaVersion {
		t.Fatalf("schema version mismatch: %s", explanation.SchemaVersion)
	}
	if len(explanation.Entries) == 0 {
		t.Fatal("expected explanation entries")
	}
	foundTruncated := false
	for _, entry := range explanation.Entries {
		if entry.StatementTruncated {
			foundTruncated = true
		}
	}
	if !foundTruncated {
		t.Fatal("expected large reason message to be truncated")
	}
}
