/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func writeEvidenceFixture(t *testing.T, name string, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestEvidencePlanAndAssessOffline(t *testing.T) {
	screening := domain.ScreeningReport{SchemaVersion: "1", Decisions: []domain.ScreeningDecision{{Repository: "org/repo@abc123", Decision: "ADVANCE", RankingScore: .9}}}
	screeningPath := writeEvidenceFixture(t, "screening.json", screening)
	var planOut bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "plan", "--screening", screeningPath}, &planOut); err != nil {
		t.Fatal(err)
	}
	var plans []domain.EvidenceExtractionPlan
	if err := json.Unmarshal(planOut.Bytes(), &plans); err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].Budget.TotalBytes != 2*1024*1024 {
		t.Fatalf("plano inesperado: %+v", plans)
	}

	idea := domain.IdeaSpec{Problem: "análise de dependências", Capabilities: []string{"dependency graph"}, Constraints: []string{"execução local"}}
	ideaPath := writeEvidenceFixture(t, "idea.json", idea)
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: "func TestDependencyGraph(t *testing.T) {}", ContentHash: "fixed"}}}
	bundlePath := writeEvidenceFixture(t, "bundle.json", []domain.EvidenceBundle{bundle})
	var assessOut bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "assess", "--idea", ideaPath, "--screening", screeningPath, "--evidence", bundlePath}, &assessOut); err != nil {
		t.Fatal(err)
	}
	var report domain.AssessmentReport
	if err := json.Unmarshal(assessOut.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Assessments) != 1 || report.Assessments[0].Capabilities[0].Status != domain.EvidenceProven {
		t.Fatalf("assessment inesperado: %+v", report)
	}
}

func TestEvidencePlanMarksMutableRepositoryForResolution(t *testing.T) {
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo", Decision: "ADVANCE"}}}
	path := writeEvidenceFixture(t, "screening.json", screening)
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "plan", "--screening", path}, &out); err != nil {
		t.Fatal(err)
	}
	var plans []domain.EvidenceExtractionPlan
	if err := json.Unmarshal(out.Bytes(), &plans); err != nil {
		t.Fatal(err)
	}
	if plans[0].Repository.Commit != "HEAD" {
		t.Fatalf("commit pendente não declarado: %+v", plans[0])
	}
}
