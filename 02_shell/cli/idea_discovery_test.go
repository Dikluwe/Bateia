/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_discovery.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

type cliDiscoveryProvider struct{ calls int }

func (provider *cliDiscoveryProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.calls++
	return domain.ProviderDiscovery{Provider: "fake", Attempts: []domain.ProviderAttempt{{Query: request.Query, ResultCount: 1}}, Candidates: []domain.RepositoryCandidate{{Owner: "org", Repo: request.ExpansionID}}}, nil
}

func TestRunIdeaDiscoverRequiresOutputBeforeProviderCall(t *testing.T) {
	provider := &cliDiscoveryProvider{}
	err := RunIdeaDiscoverCommand([]string{"--input", filepath.Join("testdata", "idea_spec.json")}, &bytes.Buffer{}, &bytes.Buffer{}, cliSemanticExpander{}, provider)
	if err == nil || provider.calls != 0 {
		t.Fatalf("output ausente acessou provedor: err=%v calls=%d", err, provider.calls)
	}
}

func TestRunIdeaDiscoverWritesVersionedSnapshot(t *testing.T) {
	provider := &cliDiscoveryProvider{}
	output := filepath.Join(t.TempDir(), "snapshot.json")
	var stdout, stderr bytes.Buffer
	err := RunIdeaDiscoverCommand([]string{
		"--input", filepath.Join("testdata", "idea_spec.json"), "--strategy", "literal",
		"--provider", "sourcegraph", "--max-per-query", "2", "--concurrency", "1", "--output", output,
	}, &stdout, &stderr, cliSemanticExpander{}, provider)
	if err != nil {
		t.Fatal(err)
	}
	if stdout.Len() != 0 || stderr.Len() == 0 {
		t.Fatalf("canais incorretos stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var snapshot domain.DiscoverySnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.SchemaVersion != "2" || provider.calls != len(snapshot.PerExpansion) {
		t.Fatalf("snapshot inválido: %+v calls=%d", snapshot, provider.calls)
	}
}

func TestRunIdeaDiscoveryEvaluateReplaysOffline(t *testing.T) {
	directory := t.TempDir()
	snapshotPath := filepath.Join(directory, "snapshot.json")
	judgmentsPath := filepath.Join(directory, "judgments.json")
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "1", Strategy: "literal", Candidates: []domain.DiscoveredCandidate{{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "repo"}}}}
	judgments := domain.DiscoveryJudgments{SchemaVersion: "1", Judgments: []domain.CandidateJudgment{{Repository: "org/repo", Label: domain.RelevanceRelevant, Reason: "match"}}}
	for path, value := range map[string]any{snapshotPath: snapshot, judgmentsPath: judgments} {
		data, _ := json.Marshal(value)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"discovery", "evaluate", "--snapshot", snapshotPath, "--judgments", judgmentsPath}, &stdout); err != nil {
		t.Fatal(err)
	}
	var evaluation domain.DiscoveryEvaluation
	if err := json.Unmarshal(stdout.Bytes(), &evaluation); err != nil {
		t.Fatal(err)
	}
	if evaluation.Metrics.JudgedPrecision == nil || *evaluation.Metrics.JudgedPrecision != 1 {
		t.Fatalf("avaliação incorreta: %+v", evaluation)
	}
}
