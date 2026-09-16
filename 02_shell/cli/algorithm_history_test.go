/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/algorithm_history_test.md
 * @layer L2
 * @updated 2026-09-12
 */

package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"bateia/02_shell/cli"
)

const cliHistoryInput = `{
  "schema_version":"algorithm-history-input-v1",
  "algorithms":[
    {"id":"alpha","name":"Alpha","chronology":"KNOWN","year_start":1,"year_end":1,"chronology_reason":"dated","provenance_ids":["p"]},
    {"id":"beta","name":"Beta","chronology":"APPROXIMATE","year_start":2,"year_end":3,"chronology_reason":"bounded","provenance_ids":["p"]}
  ],
  "provenance":[{"id":"p","kind":"PRIMARY_SOURCE","origin_id":"archive:one","citation":"Synthetic citation","locator":"archive://one"}],
  "claims":[{"id":"beta-from-alpha","subject_id":"beta","relation":"DERIVED_FROM","object_id":"alpha","statement":"Beta derives from Alpha"}],
  "bases":[{"id":"e","claim_id":"beta-from-alpha","kind":"EVIDENCE","stance":"SUPPORT","provenance_id":"p","summary":"explicit","derived_from_ids":[]}]
}`

func writeCLIInput(t *testing.T, dir, name, contents string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func decodeCLIReport(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var report map[string]any
	if err := json.Unmarshal(raw, &report); err != nil {
		t.Fatalf("stdout/output is not one JSON report: %v\n%s", err, raw)
	}
	return report
}

func TestRunAlgorithmHistoryCommandCompleteAndScientificUnknownStdout(t *testing.T) {
	dir := t.TempDir()
	completePath := writeCLIInput(t, dir, "complete.json", cliHistoryInput)
	var stdout bytes.Buffer
	if err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", completePath}, &stdout); err != nil {
		t.Fatalf("complete command: %v", err)
	}
	if !bytes.HasSuffix(stdout.Bytes(), []byte("\n")) || !bytes.Contains(stdout.Bytes(), []byte("\n  \"")) {
		t.Fatalf("stdout must be two-space indented with final newline: %q", stdout.Bytes())
	}
	report := decodeCLIReport(t, stdout.Bytes())
	if report["status"] != "COMPLETE" || report["reason"] != "analyzed" {
		t.Fatalf("complete report: %v", report)
	}

	unknownPath := writeCLIInput(t, dir, "unknown.json", `{"schema_version":"algorithm-history-input-v1","algorithms":[],"provenance":[],"claims":[],"bases":[]}`)
	stdout.Reset()
	if err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", unknownPath}, &stdout); err != nil {
		t.Fatalf("decoded scientific UNKNOWN must serialize normally: %v", err)
	}
	report = decodeCLIReport(t, stdout.Bytes())
	if report["status"] != "UNKNOWN" || report["reason"] != "no_algorithms" {
		t.Fatalf("scientific UNKNOWN report: %v", report)
	}
	assessments, ok := report["assessments"].([]any)
	if !ok || len(assessments) != 0 {
		t.Fatalf("UNKNOWN assessments must be a non-null empty array: %v", report["assessments"])
	}
}

func TestRunAlgorithmHistoryCommandStrictUsageAndJSON(t *testing.T) {
	dir := t.TempDir()
	valid := writeCLIInput(t, dir, "valid.json", cliHistoryInput)
	unknownField := writeCLIInput(t, dir, "unknown-field.json", strings.Replace(cliHistoryInput, `"schema_version":`, `"unexpected":true,"schema_version":`, 1))
	trailing := writeCLIInput(t, dir, "trailing.json", cliHistoryInput+` {}`)
	malformed := writeCLIInput(t, dir, "malformed.json", `{`)
	cases := []struct {
		name string
		args []string
	}{
		{"missing history", []string{"analyze", "--input", valid}},
		{"missing analyze", []string{"history", "--input", valid}},
		{"missing input", []string{"history", "analyze"}},
		{"missing input value", []string{"history", "analyze", "--input"}},
		{"duplicate input", []string{"history", "analyze", "--input", valid, "--input", valid}},
		{"duplicate output", []string{"history", "analyze", "--input", valid, "--output", filepath.Join(dir, "a"), "--output", filepath.Join(dir, "b")}},
		{"unknown flag", []string{"history", "analyze", "--input", valid, "--wat"}},
		{"extra positional", []string{"history", "analyze", "--input", valid, "extra"}},
		{"unreadable input", []string{"history", "analyze", "--input", filepath.Join(dir, "missing")}},
		{"unknown json field", []string{"history", "analyze", "--input", unknownField}},
		{"trailing json", []string{"history", "analyze", "--input", trailing}},
		{"malformed json", []string{"history", "analyze", "--input", malformed}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			if err := cli.RunAlgorithmHistoryCommand(tc.args, &stdout); err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("failure wrote stdout: %q", stdout.Bytes())
			}
		})
	}
}

func TestRunAlgorithmHistoryCommandDecodeFailuresDoNotTouchDestination(t *testing.T) {
	dir := t.TempDir()
	inputs := []string{
		writeCLIInput(t, dir, "malformed-destination.json", `{`),
		writeCLIInput(t, dir, "unknown-destination.json", strings.Replace(cliHistoryInput, `"schema_version":`, `"unexpected":true,"schema_version":`, 1)),
		writeCLIInput(t, dir, "trailing-destination.json", cliHistoryInput+` {}`),
		filepath.Join(dir, "missing-input.json"),
	}
	for index, input := range inputs {
		t.Run(string(rune('a'+index)), func(t *testing.T) {
			output := writeCLIInput(t, dir, "protected-"+string(rune('a'+index))+".json", "original destination bytes")
			var stdout bytes.Buffer
			if err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("expected read/decode error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("failure wrote stdout: %q", stdout.Bytes())
			}
			got, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "original destination bytes" {
				t.Fatalf("failure changed destination: %q", got)
			}
		})
	}
}

func TestRunAlgorithmHistoryCommandAtomicOutputAndPreservesInput(t *testing.T) {
	dir := t.TempDir()
	inputBytes := []byte(cliHistoryInput)
	input := writeCLIInput(t, dir, "input.json", cliHistoryInput)
	output := writeCLIInput(t, dir, "report.json", "old destination")
	beforeInfo, err := os.Stat(output)
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", input, "--output", output}, &stdout); err != nil {
		t.Fatalf("output command: %v", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("output mode wrote stdout: %q", stdout.Bytes())
	}
	gotInput, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotInput, inputBytes) {
		t.Fatal("command changed input bytes")
	}
	outputBytes, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	decodeCLIReport(t, outputBytes)
	if !bytes.HasSuffix(outputBytes, []byte("\n")) {
		t.Fatal("output lacks final newline")
	}
	afterInfo, err := os.Stat(output)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		beforeStat, beforeOK := beforeInfo.Sys().(*syscall.Stat_t)
		afterStat, afterOK := afterInfo.Sys().(*syscall.Stat_t)
		if beforeOK && afterOK && beforeStat.Ino == afterStat.Ino {
			t.Fatal("destination inode did not change; output was not atomically renamed")
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".tmp") {
			t.Errorf("temporary output leaked: %s", entry.Name())
		}
	}
}

func TestRunAlgorithmHistoryCommandRejectsInputAliases(t *testing.T) {
	dir := t.TempDir()
	input := writeCLIInput(t, dir, "input.json", cliHistoryInput)
	original, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatal(err)
	}
	aliases := []struct {
		name string
		make func(t *testing.T) string
	}{
		{"lexical", func(t *testing.T) string { return filepath.Join(dir, "sub", "..", "input.json") }},
		{"symlink", func(t *testing.T) string {
			path := filepath.Join(dir, "input-link.json")
			if err := os.Symlink(input, path); err != nil {
				t.Skipf("symlink unavailable: %v", err)
			}
			return path
		}},
		{"hardlink", func(t *testing.T) string {
			path := filepath.Join(dir, "input-hard.json")
			if err := os.Link(input, path); err != nil {
				t.Skipf("hardlink unavailable: %v", err)
			}
			return path
		}},
	}
	for _, tc := range aliases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.make(t)
			var stdout bytes.Buffer
			if err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("expected input/output alias rejection")
			}
			if stdout.Len() != 0 {
				t.Fatalf("alias failure wrote stdout: %q", stdout.Bytes())
			}
			got, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, original) {
				t.Fatal("alias rejection changed input")
			}
		})
	}
}

func TestRunAlgorithmHistoryCommandWriteFailurePreservesDestination(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX directory permissions required")
	}
	dir := t.TempDir()
	input := writeCLIInput(t, dir, "input.json", cliHistoryInput)
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o700); err != nil {
		t.Fatal(err)
	}
	output := writeCLIInput(t, locked, "report.json", "preserve me")
	if err := os.Chmod(locked, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o700) })
	var stdout bytes.Buffer
	err := cli.RunAlgorithmHistoryCommand([]string{"history", "analyze", "--input", input, "--output", output}, &stdout)
	if os.Geteuid() == 0 {
		t.Skip("root can bypass the intended write failure")
	}
	if err == nil {
		t.Fatal("expected atomic write failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("write failure wrote stdout: %q", stdout.Bytes())
	}
	got, readErr := os.ReadFile(output)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != "preserve me" {
		t.Fatalf("write failure changed destination: %q", got)
	}
}
