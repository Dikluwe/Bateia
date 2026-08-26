/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_plan_test.md
 * @layer L2
 * @updated 2026-08-24
 */

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunIdeaPlanMatchesGoldenFile(t *testing.T) {
	input := filepath.Join("testdata", "idea_spec.json")
	want, err := os.ReadFile(filepath.Join("testdata", "idea_plan.golden.json"))
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}

	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"plan", "--input", input}, &stdout); err != nil {
		t.Fatalf("RunIdeaCommand: %v", err)
	}
	if stdout.String() != string(want) {
		t.Fatalf("unexpected plan\nwant:\n%s\ngot:\n%s", want, stdout.String())
	}
}

func TestRunIdeaPlanWritesOutputFileWithoutStdout(t *testing.T) {
	output := filepath.Join(t.TempDir(), "plan.json")
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{
		"plan",
		"--input", filepath.Join("testdata", "idea_spec.json"),
		"--output", output,
	}, &stdout); err != nil {
		t.Fatalf("RunIdeaCommand: %v", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should be empty when --output is used, got %q", stdout.String())
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "idea_plan.golden.json"))
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("output file differs from golden file")
	}
}

func TestRunIdeaPlanRejectsInvalidJSON(t *testing.T) {
	input := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(input, []byte(`{"problem":`), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	err := RunIdeaCommand([]string{"plan", "--input", input}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "IdeaSpec JSON inválida") {
		t.Fatalf("expected actionable decoding error, got %v", err)
	}
}

func TestRunIdeaPlanRejectsEmptyIdea(t *testing.T) {
	input := filepath.Join(t.TempDir(), "empty.json")
	if err := os.WriteFile(input, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	err := RunIdeaCommand([]string{"plan", "--input", input}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "problem ou ao menos uma capability") {
		t.Fatalf("expected empty idea error, got %v", err)
	}
}

func TestRunIdeaPlanRequiresInput(t *testing.T) {
	err := RunIdeaCommand([]string{"plan"}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--input") {
		t.Fatalf("expected missing input error, got %v", err)
	}
}
