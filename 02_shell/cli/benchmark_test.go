/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/benchmark_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func writeBenchmarkCLIInputs(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.json")
	corpusPath := filepath.Join(dir, "corpus.json")
	if err := writeJSONAtomic(manifestPath, domain.BenchmarkManifest{ProtocolVersion: "1", CaseHashes: map[domain.BenchmarkCaseID]string{}}); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(corpusPath, domain.BenchmarkCorpus{SchemaVersion: "1", Cases: []domain.EndToEndBenchmarkCase{}}); err != nil {
		t.Fatal(err)
	}
	return manifestPath, corpusPath
}

func TestRunBenchmarkValidateIsOfflineAndEmitsReport(t *testing.T) {
	manifest, corpus := writeBenchmarkCLIInputs(t)
	var out bytes.Buffer
	if err := RunBenchmarkCommand([]string{"validate", "--manifest", manifest, "--corpus", corpus}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"valid": false`) {
		t.Fatalf("unexpected report: %s", out.String())
	}
}

func TestRunBenchmarkRunEmptySplitKeepsMetricsNullable(t *testing.T) {
	manifest, corpus := writeBenchmarkCLIInputs(t)
	var out bytes.Buffer
	if err := RunBenchmarkCommand([]string{"run", "--manifest", manifest, "--corpus", corpus, "--split", "TEST"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"decision_accuracy": null`) {
		t.Fatalf("empty denominator must be null: %s", out.String())
	}
}

func TestRunBenchmarkRejectsMissingInputsAndUnknownCommand(t *testing.T) {
	if err := RunBenchmarkCommand([]string{"validate"}, &bytes.Buffer{}); err == nil {
		t.Fatal("missing paths must fail")
	}
	if err := RunBenchmarkCommand([]string{"explode"}, &bytes.Buffer{}); err == nil {
		t.Fatal("unknown command must fail")
	}
}
