/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recover_history_regression_test.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

type historyRecoveryReader struct {
	fixtureContent
	calls int
}

func (p *historyRecoveryReader) FetchRepoMetadata(ctx context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.calls++
	return p.fixtureContent.FetchRepoMetadata(ctx, c)
}
func (p *historyRecoveryReader) FetchTree(ctx context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls++
	return p.fixtureContent.FetchTree(ctx, c)
}
func (p *historyRecoveryReader) ExtractRepoFiles(ctx context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.calls++
	return p.fixtureContent.ExtractRepoFiles(ctx, c, paths)
}
func TestRecoverCLIUsesObservedHistoryAcrossRounds(t *testing.T) {
	for _, status := range []domain.CollectionTargetStatus{domain.TargetSkippedBudget, domain.TargetCollected} {
		t.Run(string(status), func(t *testing.T) {
			run := recoveryDossierFixture(t)
			var originals []domain.EvidenceBundle
			originalPath := filepath.Join(run, "evidence", "bundles.json")
			if err := readJSONFile(originalPath, &originals); err != nil || len(originals) != 1 {
				t.Fatalf("fixture: %v", err)
			}
			ref := originals[0].Repository
			target := domain.ExtractionTarget{Path: "LICENSE", Kind: domain.EvidenceLicense}
			originals[0].Plan.Targets = []domain.ExtractionTarget{target}
			originals[0].Targets = []domain.CollectionTargetResult{{Target: target, Status: domain.TargetSkippedBudget}}
			kept := []domain.CollectedEvidence{}
			for _, item := range originals[0].Items {
				if item.Path != "LICENSE" {
					kept = append(kept, item)
				}
			}
			originals[0].Items = kept
			if err := writeJSONAtomic(originalPath, originals); err != nil {
				t.Fatal(err)
			}
			originalHash, err := fileHash(originalPath)
			if err != nil {
				t.Fatal(err)
			}
			diagnosis := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "license", Candidate: ref.Key(), Cause: domain.CauseUnknownLicense, RuleID: "GATE-LICENSE-001", Recoverable: true}}}
			policyPath := filepath.Join(t.TempDir(), "policy.json")
			if err := writeJSONAtomic(policyPath, usecase.DefaultRecoveryPolicy()); err != nil {
				t.Fatal(err)
			}
			diagnosisPath := filepath.Join(t.TempDir(), "diagnosis.json")
			if err := writeJSONAtomic(diagnosisPath, diagnosis); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			if err := RunIdeaRecoverCommand([]string{"plan", "--run", run, "--policy", policyPath, "--diagnosis", diagnosisPath}, &out, nil); err != nil {
				t.Fatal(err)
			}
			var first domain.EvidenceRecoveryPlan
			if err := json.Unmarshal(out.Bytes(), &first); err != nil {
				t.Fatal(err)
			}
			if len(first.Targets) != 1 || first.Targets[0].Path != "LICENSE" {
				t.Fatalf("original omission blocked target: %+v", first)
			}
			// A synthetic round result is a transport fixture, not a rewritten historical run.
			dir := filepath.Join(run, "recovery", "round-01")
			outcomes := append([]domain.EvidenceBundle(nil), originals...)
			outcomes[0].Targets = []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Path: "LICENSE", Kind: domain.EvidenceLicense, Reason: "recovery license"}, Status: status}}
			var before domain.ReuseDecisionReport
			if err := readJSONFile(filepath.Join(run, "reuse-decision.json"), &before); err != nil {
				t.Fatal(err)
			}
			for name, value := range map[string]any{"plan.json": first, "bundles.json": outcomes, "remaining-diagnosis.json": diagnosis, "reuse-decision.json": before} {
				if err := writeJSONAtomic(filepath.Join(dir, name), value); err != nil {
					t.Fatal(err)
				}
			}
			if err := writeJSONAtomic(filepath.Join(run, "recovery", "index.json"), domain.RecoveryIndex{SchemaVersion: "1", Rounds: []domain.RecoveryRoundIndex{{Round: 1}}}); err != nil {
				t.Fatal(err)
			}
			reader := &historyRecoveryReader{}
			if err := RunIdeaRecoverCommand([]string{"execute", "--resume", "--run", run, "--policy", policyPath}, &bytes.Buffer{}, reader); err != nil {
				t.Fatal(err)
			}
			var next domain.EvidenceRecoveryPlan
			if err := readJSONFile(filepath.Join(run, "recovery", "round-02", "plan.json"), &next); err != nil {
				t.Fatal(err)
			}
			want := "LICENSE"
			if status == domain.TargetCollected {
				want = "COPYING"
			}
			if len(next.Targets) != 1 || next.Targets[0].Path != want || reader.calls == 0 {
				t.Fatalf("round history wrong: status=%s plan=%+v calls=%d", status, next, reader.calls)
			}
			after, _ := fileHash(originalPath)
			if after != originalHash {
				t.Fatal("original bytes rewritten")
			}
		})
	}
}

func TestRecoverCLIMissingRoundOutcomeStopsBeforeProvider(t *testing.T) {
	run := recoveryDossierFixture(t)
	policyPath := filepath.Join(t.TempDir(), "policy.json")
	if err := writeJSONAtomic(policyPath, usecase.DefaultRecoveryPolicy()); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(run, "recovery", "round-01")
	if err := writeJSONAtomic(filepath.Join(dir, "remaining-diagnosis.json"), domain.InconclusiveDiagnosis{SchemaVersion: "1"}); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(filepath.Join(dir, "plan.json"), domain.EvidenceRecoveryPlan{SchemaVersion: "1", Round: 1}); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(filepath.Join(run, "recovery", "index.json"), domain.RecoveryIndex{SchemaVersion: "1", Rounds: []domain.RecoveryRoundIndex{{Round: 1}}}); err != nil {
		t.Fatal(err)
	}
	reader := &historyRecoveryReader{}
	err := RunIdeaRecoverCommand([]string{"execute", "--resume", "--run", run, "--policy", policyPath}, &bytes.Buffer{}, reader)
	if err == nil || !strings.Contains(err.Error(), "bundles.json") || reader.calls != 0 {
		t.Fatalf("missing outcome ignored: err=%v calls=%d", err, reader.calls)
	}
	if _, err := commandHost().Stat(filepath.Join(run, "recovery", "round-02")); err == nil {
		t.Fatal("published despite missing outcome")
	}
}
