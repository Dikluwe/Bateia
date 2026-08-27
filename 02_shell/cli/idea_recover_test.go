/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recover_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func recoveryDossierFixture(t *testing.T) string {
	t.Helper()
	idea, policy, source := recommendReplayFixture(t)
	out := filepath.Join(t.TempDir(), "run")
	if err := RunIdeaRecommendCommand([]string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", out}, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	return out
}
func TestRecoverDiagnoseAndPlanAreOfflineAndTraceable(t *testing.T) {
	run := recoveryDossierFixture(t)
	var diagnosisOut bytes.Buffer
	if err := RunIdeaRecoverCommand([]string{"diagnose", "--run", run}, &diagnosisOut, nil); err != nil {
		t.Fatal(err)
	}
	var d domain.InconclusiveDiagnosis
	if err := json.Unmarshal(diagnosisOut.Bytes(), &d); err != nil {
		t.Fatal(err)
	}
	if len(d.Findings) == 0 {
		t.Fatal("expected findings")
	}
	diagnosisPath := filepath.Join(t.TempDir(), "diagnosis.json")
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	if err := writeJSONAtomic(diagnosisPath, d); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(policyPath, usecase.DefaultRecoveryPolicy()); err != nil {
		t.Fatal(err)
	}
	var planOut bytes.Buffer
	if err := RunIdeaRecoverCommand([]string{"plan", "--run", run, "--diagnosis", diagnosisPath, "--policy", policyPath}, &planOut, nil); err != nil {
		t.Fatal(err)
	}
	var p domain.EvidenceRecoveryPlan
	if err := json.Unmarshal(planOut.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	for _, target := range p.Targets {
		if target.FindingID == "" {
			t.Fatal("target without finding")
		}
	}
}
func TestRecoverExecutePreservesOriginalAndWritesHashedRound(t *testing.T) {
	run := recoveryDossierFixture(t)
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	if err := writeJSONAtomic(policyPath, usecase.DefaultRecoveryPolicy()); err != nil {
		t.Fatal(err)
	}
	original, err := fileHash(filepath.Join(run, "run.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaRecoverCommand([]string{"execute", "--run", run, "--policy", policyPath}, &bytes.Buffer{}, fixtureContent{}); err != nil {
		t.Fatal(err)
	}
	after, _ := fileHash(filepath.Join(run, "run.json"))
	if original != after {
		t.Fatal("original run mutated")
	}
	for _, path := range []string{"recovery/index.json", "recovery/round-01/diagnosis.json", "recovery/round-01/plan.json", "recovery/round-01/bundles.json", "recovery/round-01/assessment-report.json", "recovery/round-01/reuse-decision.json", "recovery/round-01/metrics.json"} {
		if _, err := os.Stat(filepath.Join(run, path)); err != nil {
			t.Fatalf("missing %s", path)
		}
	}
	var index domain.RecoveryIndex
	if err := readJSONFile(filepath.Join(run, "recovery", "index.json"), &index); err != nil {
		t.Fatal(err)
	}
	if len(index.Rounds) != 1 || index.Rounds[0].DecisionHash == "" {
		t.Fatalf("index %#v", index)
	}
}
func TestRecoverExecuteRequiresProviderBeforeMutation(t *testing.T) {
	run := recoveryDossierFixture(t)
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	_ = writeJSONAtomic(policyPath, usecase.DefaultRecoveryPolicy())
	if err := RunIdeaRecoverCommand([]string{"execute", "--run", run, "--policy", policyPath}, &bytes.Buffer{}, nil); err == nil {
		t.Fatal("missing provider accepted")
	}
	if _, err := os.Stat(filepath.Join(run, "recovery")); !os.IsNotExist(err) {
		t.Fatal("recovery mutated before provider validation")
	}
}
