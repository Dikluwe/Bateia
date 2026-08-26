/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_test.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	idea := domain.IdeaSpec{Problem: "análise de dependências", Capabilities: []string{"dependency graph"}, Constraints: []string{"execução local"}}
	ideaPath := writeEvidenceFixture(t, "idea.json", idea)
	var planOut bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "plan", "--idea", ideaPath, "--screening", screeningPath}, &planOut); err != nil {
		t.Fatal(err)
	}
	var plans []domain.EvidenceExtractionPlan
	if err := json.Unmarshal(planOut.Bytes(), &plans); err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].Budget.TotalBytes != 2*1024*1024 {
		t.Fatalf("plano inesperado: %+v", plans)
	}

	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	testContent := "func TestDependencyGraph(t *testing.T) {}"
	sum := sha256.Sum256([]byte(testContent))
	bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceTest, Ref: ref, Path: "graph_test.go", Content: testContent, ContentHash: hex.EncodeToString(sum[:])}}}
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
	assessmentPath := writeEvidenceFixture(t, "assessment.json", report)
	judgmentsPath := writeEvidenceFixture(t, "judgments.json", []domain.EvidenceJudgment{{Repository: "org/repo", Capability: "dependency graph", Expected: domain.EvidenceProven}})
	var metricsOut bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "evaluate", "--assessment", assessmentPath, "--evidence", bundlePath, "--judgments", judgmentsPath}, &metricsOut); err != nil {
		t.Fatal(err)
	}
	var metrics domain.EvidenceMetrics
	if err := json.Unmarshal(metricsOut.Bytes(), &metrics); err != nil {
		t.Fatal(err)
	}
	if metrics.EvidencePrecision == nil || *metrics.EvidencePrecision != 1 {
		t.Fatalf("métricas inesperadas: %+v", metrics)
	}
}

func TestEvidencePlanMarksMutableRepositoryForResolution(t *testing.T) {
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo", Decision: "ADVANCE"}}}
	path := writeEvidenceFixture(t, "screening.json", screening)
	ideaPath := writeEvidenceFixture(t, "idea.json", domain.IdeaSpec{Capabilities: []string{"static analysis"}})
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"evidence", "plan", "--idea", ideaPath, "--screening", path}, &out); err != nil {
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
