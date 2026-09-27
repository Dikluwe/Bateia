/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_review_test.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func writeReviewFixture(t *testing.T, name string, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestIdeaReviewProtocolAndPrioritizeCommands(t *testing.T) {
	intake := domain.RecommendationIntake{
		SchemaVersion: "1", Status: domain.IntakeReady, InputHash: "intake-hash",
		Request: domain.RecommendationRequest{
			Idea:   domain.IdeaSpec{Problem: "search", Capabilities: []string{"index"}},
			Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, ExtendMinCoverage: .8, IntegrationMinValue: .5},
			Target: domain.SoftwareAlternativeTarget{Name: "Indexer"},
			Search: domain.RecommendationSearch{Strategy: "lexical", Locales: []string{"en"}, MaxPerQuery: 10, MaxCandidates: 20},
		},
		Execution: domain.RecommendationExecutionProfile{Mode: "replay", Provider: "fixture", LocalPaths: []string{}},
	}
	draft := domain.ReviewProtocolDraft{Intake: intake, ScreeningModel: usecase.DefaultScreeningModel(), EvidenceRulesVersion: "v1", EvidenceRulesHash: "hash"}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"review", "protocol", "--input", writeReviewFixture(t, "protocol.json", draft)}, &stdout); err != nil {
		t.Fatal(err)
	}
	var protocol domain.ReviewProtocol
	if err := json.Unmarshal(stdout.Bytes(), &protocol); err != nil {
		t.Fatal(err)
	}
	if protocol.ProtocolHash == "" || len(protocol.Eligibility) == 0 {
		t.Fatalf("protocolo incompleto: %+v", protocol)
	}

	stdout.Reset()
	request := domain.ActiveLearningRequest{BatchSize: 1, Screening: domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "acme/tool", RankingScore: .8, Signals: domain.ScreeningSignals{EvidenceCompleteness: .5}}}}}
	if err := RunIdeaCommand([]string{"review", "prioritize", "--input", writeReviewFixture(t, "priority.json", request)}, &stdout); err != nil {
		t.Fatal(err)
	}
	var plan domain.ActiveLearningPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.CanExclude || plan.Authority != "PRIORITIZATION_ONLY" || len(plan.Priorities) != 1 {
		t.Fatalf("plano inválido: %+v", plan)
	}
}

func TestIdeaReviewRejectsUnknownJSONFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(path, []byte(`{"batch_size":1,"screening":{},"hidden_gold":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
	err := RunIdeaCommand([]string{"review", "prioritize", "--input", path}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("campo desconhecido deveria ser rejeitado: %v", err)
	}
}
