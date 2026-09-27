/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace_explain_test.md
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

func traceExplainFixture(t *testing.T) domain.DecisionTraceExplainInput {
	t.Helper()
	build := traceBuildFixture()
	build.Report.Reasons = []domain.DecisionReason{{RuleID: "missing-proof", Message: "Evidence remains declarative"}}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	return domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build}
}

func TestIdeaTraceExplainOfflineBindsLocationsAndPreservesInput(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "explain.json"), filepath.Join(dir, "index.json")
	request := traceExplainFixture(t)
	original := writeTraceFixture(t, input, request)
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "explain", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var index domain.DecisionTraceExplanation
	if err := json.Unmarshal(stdout.Bytes(), &index); err != nil {
		t.Fatal(err)
	}
	if index.GraphHash != request.Graph.GraphHash || index.Scope != domain.DecisionTraceScope || len(index.Entries) != len(request.Graph.Nodes) {
		t.Fatalf("unbound explanation: %+v", index)
	}
	foundReason := false
	for _, entry := range index.Entries {
		for _, location := range entry.Locations {
			if location.Document == "report" && location.Pointer == "/reasons/0" {
				foundReason = true
			}
		}
	}
	if !foundReason {
		t.Fatal("reason location omitted")
	}
	var silent bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "explain", "--input", input, "--output", output}, &silent); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(data, stdout.Bytes()) || silent.Len() != 0 {
		t.Fatal("explain disk/stdout mismatch")
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(original, after) {
		t.Fatal("explanation changed input")
	}
	request.NodeIDs = []string{request.Graph.Nodes[0].ID}
	writeTraceFixture(t, input, request)
	stdout.Reset()
	if err := RunIdeaCommand([]string{"trace", "explain", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(stdout.Bytes(), &index); err != nil || len(index.Entries) != 1 || index.Entries[0].NodeID != request.NodeIDs[0] {
		t.Fatalf("invalid selection: %+v %v", index, err)
	}
}

func TestIdeaTraceExplainRejectsInvalidBindingsBeforeOutput(t *testing.T) {
	valid := traceExplainFixture(t)
	encoded, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{}
	mutate := func(name string, change func(*domain.DecisionTraceExplainInput)) {
		request := traceExplainFixture(t)
		change(&request)
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		cases[name] = data
	}
	mutate("report-mismatch", func(input *domain.DecisionTraceExplainInput) {
		input.BuildInput.Report.Reasons[0].Message = "different assertion"
	})
	mutate("unknown-node", func(input *domain.DecisionTraceExplainInput) { input.NodeIDs = []string{"not-recorded"} })
	mutate("duplicate-node", func(input *domain.DecisionTraceExplainInput) {
		input.NodeIDs = []string{input.Graph.Nodes[0].ID, input.Graph.Nodes[0].ID}
	})
	mutate("missing-hash", func(input *domain.DecisionTraceExplainInput) { input.Graph.GraphHash = "" })
	mutate("tampered-hash", func(input *domain.DecisionTraceExplainInput) { input.Graph.GraphHash = strings.Repeat("0", 64) })
	mutate("schema", func(input *domain.DecisionTraceExplainInput) { input.SchemaVersion = "future" })
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
			prior := []byte("previous result")
			if err := os.WriteFile(output, prior, 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"trace", "explain", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
				t.Fatal("invalid explanation accepted or emitted stdout")
			}
			after, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(prior, after) {
				t.Fatal("invalid explanation replaced output")
			}
			missing := filepath.Join(dir, "missing.json")
			if err := RunIdeaCommand([]string{"trace", "explain", "--input", input, "--output", missing}, &stdout); err == nil {
				t.Fatal("invalid explanation accepted missing output")
			}
			if _, err := os.Lstat(missing); !os.IsNotExist(err) {
				t.Fatal("invalid explanation created output")
			}
		})
	}
}

func TestIdeaTraceExplainRefusesAliasesAndFIFO(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	original := writeTraceFixture(t, input, traceExplainFixture(t))
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
		if err := RunIdeaCommand([]string{"trace", "explain", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
			t.Fatalf("unsafe explain output accepted: %s", output)
		}
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "explain", "--input", fifo}, &stdout); err == nil || stdout.Len() != 0 {
		t.Fatal("FIFO explanation input accepted")
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(original, after) {
		t.Fatal("failed explanation overwrote input")
	}
}
