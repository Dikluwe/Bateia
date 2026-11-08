/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/resume_diagnostic_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestResumeDiagnosticMessagesAndPrivacy(t *testing.T) {
	for _, test := range []struct {
		reason usecase.EvidenceCheckpointReason
		text   string
	}{
		{usecase.CheckpointReusable, "compatível e hashes conferidos"},
		{usecase.CheckpointManifestUnavailable, "ler ou interpretar o manifesto"},
		{usecase.CheckpointStatusIneligible, "estado da execução"},
		{usecase.CheckpointSchemaUnsupported, "schema do manifesto"},
		{usecase.CheckpointEvidenceChanged, "entradas da evidência"},
		{usecase.CheckpointRulesChanged, "catálogo de regras"},
		{usecase.CheckpointModelChanged, "modelo de triagem"},
		{usecase.CheckpointVersionChanged, "versão de regra"},
		{usecase.CheckpointReaderUnavailable, "leitor dos artefatos"},
		{usecase.CheckpointArtifactUnreadable, "ler o artefato"},
		{usecase.CheckpointArtifactMismatch, "hash do artefato diverge"},
		{usecase.CheckpointArtifactUnrecorded, "registro com hash"},
	} {
		t.Run(string(test.reason), func(t *testing.T) {
			var stdout bytes.Buffer
			assessment := &usecase.EvidenceCheckpointAssessment{Reason: test.reason, Reusable: test.reason == usecase.CheckpointReusable, Subject: "/private/token\ninjected"}
			writeResumeDiagnostic(&stdout, assessment)
			got := stdout.String()
			if !strings.Contains(got, test.text) || !strings.Contains(got, string(test.reason)) || strings.Count(got, "\n") != 1 || strings.Contains(got, "/private/") || strings.Contains(got, "injected") {
				t.Fatalf("unsafe or unclear diagnostic: %q", got)
			}
			wantState := "não reutilizada"
			if assessment.Reusable {
				wantState = "anterior reutilizada"
			}
			if !strings.Contains(got, wantState) {
				t.Fatalf("state not presented: %q", got)
			}
		})
	}
	var stdout bytes.Buffer
	writeResumeDiagnostic(&stdout, nil)
	if stdout.Len() != 0 {
		t.Fatal("no assessment must be silent")
	}
	writeResumeDiagnostic(&stdout, &usecase.EvidenceCheckpointAssessment{Reason: "future-secret\nvalue", Subject: "evidence"})
	if !strings.Contains(stdout.String(), "UNKNOWN_CHECKPOINT_REASON") || strings.Contains(stdout.String(), "future-secret") || strings.Contains(stdout.String(), "; evidence") {
		t.Fatalf("unknown diagnostic leaked fields: %q", stdout.String())
	}
}

func TestIdeaRecommendResumeDiagnosticsWithReplay(t *testing.T) {
	for _, scenario := range []string{"reusable", "input", "hash", "unreadable", "missing-manifest", "invalid-manifest", "version"} {
		t.Run(scenario, func(t *testing.T) {
			idea, policy, source := recommendReplayFixture(t)
			output := filepath.Join(t.TempDir(), "out")
			args := []string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", output}
			var stdout bytes.Buffer
			if err := RunIdeaRecommendCommand(args, &stdout, nil, nil, nil); err != nil {
				t.Fatal(err)
			}
			firstLine := stdout.String()
			if !strings.HasPrefix(firstLine, "dossiê ") || strings.Contains(firstLine, "retomada:") || strings.Count(firstLine, "\n") != 1 {
				t.Fatalf("fresh output changed: %q", firstLine)
			}
			var before domain.IdeaRunManifest
			if err := readJSONFile(filepath.Join(output, "run.json"), &before); err != nil {
				t.Fatal(err)
			}
			entriesBefore, err := os.ReadDir(output)
			if err != nil {
				t.Fatal(err)
			}
			wantCode, wantSubject := "CHECKPOINT_REUSABLE", ""
			switch scenario {
			case "input":
				args = append(args, "--language", "rust")
				wantCode, wantSubject = "EVIDENCE_INPUT_CHANGED", "evidence"
			case "hash":
				if err := os.WriteFile(filepath.Join(output, "discovery-snapshot.json"), []byte("tampered /private/secret"), 0600); err != nil {
					t.Fatal(err)
				}
				wantCode, wantSubject = "ARTIFACT_HASH_MISMATCH", "discovery-snapshot.json"
			case "unreadable":
				if err := os.Remove(filepath.Join(output, "evidence", "bundles.json")); err != nil {
					t.Fatal(err)
				}
				wantCode, wantSubject = "ARTIFACT_UNREADABLE", "evidence/bundles.json"
			case "missing-manifest":
				if err := os.Remove(filepath.Join(output, "run.json")); err != nil {
					t.Fatal(err)
				}
				wantCode, wantSubject = "MANIFEST_UNAVAILABLE", "run.json"
			case "invalid-manifest":
				if err := os.WriteFile(filepath.Join(output, "run.json"), []byte("invalid /private/secret"), 0600); err != nil {
					t.Fatal(err)
				}
				wantCode, wantSubject = "MANIFEST_UNAVAILABLE", "run.json"
			case "version":
				old := before
				old.RuleVersions = map[string]string{}
				for key, value := range before.RuleVersions {
					old.RuleVersions[key] = value
				}
				old.RuleVersions["analyzer"] = "/private/secret"
				if err := writeJSONAtomic(filepath.Join(output, "run.json"), old); err != nil {
					t.Fatal(err)
				}
				wantCode, wantSubject = "RULE_VERSION_CHANGED", "analyzer"
			}
			stdout.Reset()
			if err := RunIdeaRecommendCommand(append(args, "--resume"), &stdout, nil, nil, nil); err != nil {
				t.Fatal(err)
			}
			got := stdout.String()
			if !strings.Contains(got, "retomada:") || !strings.Contains(got, wantCode) || (wantSubject != "" && !strings.Contains(got, "; "+wantSubject+")")) || strings.Contains(got, "/private/secret") || strings.Count(got, "\n") != 2 {
				t.Fatalf("resume diagnostic: %q", got)
			}
			if scenario != "reusable" && !strings.Contains(got, "não reutilizada") {
				t.Fatalf("rejection not explicit: %q", got)
			}
			var after domain.IdeaRunManifest
			if err := readJSONFile(filepath.Join(output, "run.json"), &after); err != nil {
				t.Fatal(err)
			}
			if scenario != "input" && (before.RunID != after.RunID || !reflect.DeepEqual(before.InputHashes, after.InputHashes) || !reflect.DeepEqual(before.RuleVersions, after.RuleVersions)) {
				t.Fatal("diagnostic changed causal identity or versions")
			}
			entriesAfter, err := os.ReadDir(output)
			if err != nil || len(entriesBefore) != len(entriesAfter) {
				t.Fatal("unexpected sidecar or missing artifact")
			}
		})
	}
}

func TestIdeaRecommendResumeDiagnosticMatchesLiveWork(t *testing.T) {
	fixture := newCheckpointFixture(t)
	args := []string{"--idea", fixture.ideaPath, "--policy", fixture.policyPath, "--mode", "live", "--strategy", "literal", "--provider", "sourcegraph", "--language", "go", "--max-per-query", "3", "--max-candidates", "2", "--output", fixture.outputPath, "--resume"}
	for _, corrupt := range []bool{false, true} {
		if !corrupt {
			if err := writeJSONAtomic(fixture.policyPath, checkpointPolicy("Apache-2.0")); err != nil {
				t.Fatal(err)
			}
		} else if err := os.WriteFile(filepath.Join(fixture.outputPath, "evidence", "bundles.json"), []byte("tampered"), 0600); err != nil {
			t.Fatal(err)
		}
		calls, extractions := fixture.provider.calls, fixture.content.extractCalls
		var stdout bytes.Buffer
		if err := RunIdeaRecommendCommand(args, &stdout, nil, fixture.provider, fixture.content); err != nil {
			t.Fatal(err)
		}
		if !corrupt {
			if !strings.Contains(stdout.String(), "CHECKPOINT_REUSABLE") || calls != fixture.provider.calls || extractions != fixture.content.extractCalls {
				t.Fatalf("policy-only change should reuse: %s", stdout.String())
			}
		} else if !strings.Contains(stdout.String(), "ARTIFACT_HASH_MISMATCH; evidence/bundles.json") || fixture.provider.calls <= calls || fixture.content.extractCalls <= extractions {
			t.Fatalf("corrupted evidence should be recollected: %s", stdout.String())
		}
	}
}
