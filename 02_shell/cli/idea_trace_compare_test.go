/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace_compare_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func traceCompareFixture(t *testing.T) domain.DecisionTraceCompareInput {
	t.Helper()
	build := traceBuildFixture()
	before, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	build.Report.Decision = domain.DecisionBuild
	build.Report.Reasons = []domain.DecisionReason{{RuleID: "new-assertion", Message: "New snapshot contains a different declared verdict"}}
	after, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	return domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after}
}

func TestIdeaTraceCompareOfflinePreservesHistoricalDecisionAndInput(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "compare.json"), filepath.Join(dir, "comparison.json")
	request := traceCompareFixture(t)
	original := writeTraceFixture(t, input, request)
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "compare", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var comparison domain.DecisionTraceComparison
	if err := json.Unmarshal(stdout.Bytes(), &comparison); err != nil {
		t.Fatal(err)
	}
	if comparison.BeforeGraphHash != request.Before.GraphHash || comparison.AfterGraphHash != request.After.GraphHash || comparison.HistoricalDecision != request.Before.Decision || comparison.CandidateDecision != request.After.Decision {
		t.Fatalf("comparison lost snapshot identity: %+v", comparison)
	}
	if comparison.DecisionImpact != "NONE" || comparison.Impact.DecisionBefore != request.Before.Decision || comparison.Impact.DecisionAfter != request.Before.Decision || comparison.Impact.Status != "REVIEW_REQUIRED" {
		t.Fatalf("comparison promoted supplied candidate verdict: %+v", comparison)
	}
	if err := os.WriteFile(output, []byte("old comparison"), 0600); err != nil {
		t.Fatal(err)
	}
	var silent bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "compare", "--input=" + input, "--output=" + output}, &silent); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(data, stdout.Bytes()) || silent.Len() != 0 {
		t.Fatal("comparison disk/stdout mismatch")
	}
	actual, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(actual, original) {
		t.Fatal("comparison changed input")
	}
	temporary, err := filepath.Glob(filepath.Join(dir, ".bateia-trace-*"))
	if err != nil || len(temporary) != 0 {
		t.Fatalf("temporary comparison results leaked: %v %v", temporary, err)
	}
	request.After = request.Before
	writeTraceFixture(t, input, request)
	stdout.Reset()
	if err := RunIdeaCommand([]string{"trace", "compare", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(stdout.Bytes(), &comparison); err != nil || comparison.Impact.Status != "NO_RECORDED_IMPACT" || len(comparison.Impact.Tasks) != 0 {
		t.Fatalf("identical snapshots scheduled review: %+v %v", comparison, err)
	}
}

func TestIdeaTraceCompareRejectsInvalidSnapshotsBeforeOutput(t *testing.T) {
	valid := traceCompareFixture(t)
	encoded, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{}
	mutate := func(name string, change func(*domain.DecisionTraceCompareInput)) {
		request := traceCompareFixture(t)
		change(&request)
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		cases[name] = data
	}
	mutate("before-hash-missing", func(request *domain.DecisionTraceCompareInput) { request.Before.GraphHash = "" })
	mutate("after-hash-missing", func(request *domain.DecisionTraceCompareInput) { request.After.GraphHash = "" })
	mutate("before-hash-tampered", func(request *domain.DecisionTraceCompareInput) { request.Before.GraphHash = strings.Repeat("0", 64) })
	mutate("after-hash-tampered", func(request *domain.DecisionTraceCompareInput) { request.After.GraphHash = strings.Repeat("0", 64) })
	mutate("declaration-tampered", func(request *domain.DecisionTraceCompareInput) { request.After.Nodes[0].Label += " changed" })
	mutate("schema", func(request *domain.DecisionTraceCompareInput) { request.SchemaVersion = "future" })
	cases["unknown-field"] = append(append([]byte(nil), encoded[:len(encoded)-1]...), []byte(`,"unknown":true}`)...)
	cases["duplicate-field"] = append(append([]byte(nil), encoded[:len(encoded)-1]...), []byte(`,"schema_version":"1"}`)...)
	cases["casefold-field"] = append(append([]byte(nil), encoded[:len(encoded)-1]...), []byte(`,"SCHEMA_VERSION":"1"}`)...)
	cases["trailing"] = append(append([]byte(nil), encoded...), []byte(` {}`)...)
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			input, output := filepath.Join(dir, "input.json"), filepath.Join(dir, "output.json")
			if err := os.WriteFile(input, data, 0600); err != nil {
				t.Fatal(err)
			}
			prior := []byte("previous comparison")
			if err := os.WriteFile(output, prior, 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"trace", "compare", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
				t.Fatal("invalid comparison accepted or emitted stdout")
			}
			actual, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(actual, prior) {
				t.Fatal("invalid comparison replaced output")
			}
			missing := filepath.Join(dir, "missing.json")
			if err := RunIdeaCommand([]string{"trace", "compare", "--input", input, "--output", missing}, &stdout); err == nil {
				t.Fatal("invalid comparison accepted missing output")
			}
			if _, err := os.Lstat(missing); !os.IsNotExist(err) {
				t.Fatal("invalid comparison created output")
			}
			stdout.Reset()
			if err := RunIdeaCommand([]string{"trace", "compare", "--input", input}, &stdout); err == nil || stdout.Len() != 0 {
				t.Fatal("invalid comparison emitted partial stdout")
			}
		})
	}
}

func TestIdeaTraceCompareRefusesAliasesFIFOAndAmbiguousArguments(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	original := writeTraceFixture(t, input, traceCompareFixture(t))
	hard, symbolic, fifo := filepath.Join(dir, "hard"), filepath.Join(dir, "symbolic"), filepath.Join(dir, "fifo")
	if err := os.Link(input, hard); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(input, symbolic); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	for _, output := range []string{input, hard, symbolic, fifo} {
		var stdout bytes.Buffer
		if err := RunIdeaCommand([]string{"trace", "compare", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
			t.Fatalf("unsafe comparison output accepted: %s", output)
		}
	}
	arguments := [][]string{
		{"compare"}, {"compare", "--input"}, {"compare", "--input="},
		{"compare", "--input", input, "--input", input},
		{"compare", "--input", input, "--output", "a", "--output", "b"},
		{"compare", "--input", input, "positional"},
		{"compare", "--input", input, "--unknown", "x"},
		{"compare", "--input", fifo},
	}
	for _, args := range arguments {
		var stdout bytes.Buffer
		if err := RunIdeaCommand(append([]string{"trace"}, args...), &stdout); err == nil || stdout.Len() != 0 {
			t.Fatalf("unsafe comparison arguments accepted: %v", args)
		}
	}
	actual, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(original, actual) {
		t.Fatal("failed comparison overwrote input")
	}
}
