/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_comparison_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func comparisonCLIRequest(content string) domain.InsightComparisonRequest {
	digest := sha256.Sum256([]byte(content))
	hash := "sha256:" + hex.EncodeToString(digest[:])
	manifest := `[{"path":"target.go","content_hash":"` + hash + `"}]`
	snapshot := sha256.Sum256([]byte(manifest))
	r := domain.InsightComparisonRequest{SchemaVersion: "1", Investigation: sourceInsightCLIRequest(), TargetBaseline: &domain.InsightTargetBaseline{
		Repository: "org/target", SnapshotSHA256: "sha256:" + hex.EncodeToString(snapshot[:]),
		Files: []domain.InsightBaselineFile{{Path: "target.go", Content: content, ContentHash: hash}},
	}}
	r.TargetBaseline.Observations = []domain.SourceInsightObservation{{MechanismClass: r.Investigation.MechanismClasses[0], Kind: domain.InsightPrecondition, Statement: "A configuração do alvo está declarada neste trecho", References: []domain.InsightObservationReference{{Path: "target.go", ContentHash: hash, StartLine: 1, EndLine: 1}}}}
	return r
}

func TestInsightComparisonCLIRoutesAndProtectsTransport(t *testing.T) {
	directory := t.TempDir()
	input, output := filepath.Join(directory, "input.json"), filepath.Join(directory, "output.json")
	data, _ := json.Marshal(comparisonCLIRequest("configure target capacity\n"))
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"insights", "compare", "--input", input}
	var stdout bytes.Buffer
	if err := RunIdeaCommand(args, &stdout); err != nil {
		t.Fatal(err)
	}
	var report domain.InsightComparisonReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "insight-baseline-comparison-v1" || len(report.Comparisons) != 5 || report.TargetBaseline.Observations[0].References[0].Excerpt != "configure target capacity" || report.EffectState != "UNVERIFIED" {
		t.Fatal("comparison not projected with scoped baseline")
	}
	expected := append([]byte(nil), stdout.Bytes()...)
	stdout.Reset()
	if err := RunIdeaCommand(append(args, "--output", output), &stdout); err != nil {
		t.Fatal(err)
	}
	saved, _ := os.ReadFile(output)
	if !bytes.Equal(saved, expected) || stdout.Len() != 0 {
		t.Fatal("file/stdout differ")
	}
	if err := RunIdeaCommand(args, insightShortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short write ignored: %v", err)
	}
	alias := filepath.Join(directory, "alias.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand(append(args, "--output", alias), &stdout); err == nil {
		t.Fatal("input alias overwritten")
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(after, data) {
		t.Fatal("input mutated")
	}
}

func TestInsightComparisonCLIRejectsInvalidJSONAndKeepsPriorOutput(t *testing.T) {
	request := comparisonCLIRequest("configure target capacity\n")
	data, _ := json.Marshal(request)
	request.TargetBaseline.Files[0].Content += "altered"
	invalid, _ := json.Marshal(request)
	for name, payload := range map[string][]byte{
		"unknown":          append([]byte(`{"effect_verified":true,`), data[1:]...),
		"duplicate":        append([]byte(`{"schema_version":"1",`), data[1:]...),
		"nested-unknown":   []byte(strings.Replace(string(data), `"snapshot_sha256":`, `"commit":"invented","snapshot_sha256":`, 1)),
		"nested-duplicate": []byte(strings.Replace(string(data), `"repository":"org/target"`, `"repository":"org/target","repository":"org/target"`, 1)),
		"trailing":         append(append([]byte(nil), data...), []byte(` {}`)...),
		"schema":           []byte(strings.Replace(string(data), `"schema_version":"1"`, `"schema_version":"3"`, 1)),
		"tampered":         invalid,
	} {
		t.Run(name, func(t *testing.T) {
			directory := t.TempDir()
			input, output := filepath.Join(directory, "in.json"), filepath.Join(directory, "out.json")
			if err := os.WriteFile(input, payload, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"insights", "compare", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("invalid request accepted")
			}
			saved, _ := os.ReadFile(output)
			if string(saved) != "preserve" || stdout.Len() != 0 {
				t.Fatal("invalid request published output")
			}
		})
	}
}

func TestInsightComparisonCLIProtectsSecretsInTargetExcerpts(t *testing.T) {
	const secret = "test-0135-not-a-real-credential"
	t.Setenv("GITHUB_TOKEN", secret)
	data, _ := json.Marshal(comparisonCLIRequest("configuration " + secret + "\n"))
	input := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"insights", "compare", "--input", input}, &stdout)
	if err == nil || stdout.Len() != 0 || strings.Contains(err.Error(), secret) {
		t.Fatal("secret published from target excerpt")
	}
}
