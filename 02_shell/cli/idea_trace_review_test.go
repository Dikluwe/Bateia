/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace_review_test.md
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

func traceReviewCLIFixture(t *testing.T) domain.DecisionTraceReviewInput {
	t.Helper()
	graph, err := usecase.BuildDecisionTrace(traceBuildFixture())
	if err != nil {
		t.Fatal(err)
	}
	return domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: graph.Nodes[0].ID, Reason: "INPUT_CHANGED"}}}}
}

func TestIdeaTraceReviewLedgerPreservesInputAndDecision(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	output := filepath.Join(dir, "output.json")
	request := traceReviewCLIFixture(t)
	writeTraceFixture(t, input, request)
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var seed domain.DecisionTraceReviewLedger
	if err := json.Unmarshal(stdout.Bytes(), &seed); err != nil {
		t.Fatal(err)
	}
	if seed.Status != "REVIEW_PENDING" || len(seed.Tasks) == 0 || seed.Counts.Pending != len(seed.Tasks) {
		t.Fatal("missing pending tasks")
	}
	for _, task := range seed.Tasks {
		request.Records = append(request.Records, domain.DecisionTraceReviewRecord{NodeID: task.Task.NodeID, PlanHash: seed.PlanHash, TaskHash: task.TaskHash, Status: "COMPLETED", Note: "declared result", Evidence: []domain.DecisionTraceReviewEvidence{{Ref: "logical-document", ContentHash: strings.Repeat("a", 64)}}})
	}
	original := writeTraceFixture(t, input, request)
	stdout.Reset()
	if err := RunIdeaCommand([]string{"trace", "review", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result domain.DecisionTraceReviewLedger
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "RECORDED_REVIEW_COMPLETE" || result.DecisionImpact != "NONE" || result.DecisionBefore != request.ImpactInput.Graph.Decision || result.DecisionAfter != result.DecisionBefore {
		t.Fatal("completed review promoted historical decision")
	}
	if err := os.WriteFile(output, []byte("previous ledger"), 0600); err != nil {
		t.Fatal(err)
	}
	var silent bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review", "--input=" + input, "--output=" + output}, &silent); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(data, stdout.Bytes()) || silent.Len() != 0 {
		t.Fatal("ledger disk/stdout mismatch")
	}
	unchanged, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(unchanged, original) {
		t.Fatal("input modified")
	}
}

func TestIdeaTraceReviewRejectsInvalidBeforeOutput(t *testing.T) {
	request := traceReviewCLIFixture(t)
	seed, err := usecase.TrackDecisionTraceReview(request)
	if err != nil {
		t.Fatal(err)
	}
	request.Records = []domain.DecisionTraceReviewRecord{{NodeID: seed.Tasks[0].Task.NodeID, PlanHash: seed.PlanHash, TaskHash: seed.Tasks[0].TaskHash, Status: "COMPLETED", Note: "declared result", Evidence: []domain.DecisionTraceReviewEvidence{{Ref: "logical-result", ContentHash: strings.Repeat("a", 64)}}}}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{}
	for name, addition := range map[string]string{"unknown": `,"unknown":true}`, "duplicate": `,"schema_version":"1"}`, "casefold": `,"SCHEMA_VERSION":"1"}`} {
		cases[name] = append(append([]byte{}, encoded[:len(encoded)-1]...), []byte(addition)...)
	}
	cases["trailing"] = append(append([]byte{}, encoded...), []byte(` {}`)...)
	cases["stale receipt"] = bytes.Replace(encoded, []byte(seed.PlanHash), []byte(strings.Repeat("b", 64)), 1)
	cases["completed no evidence"] = bytes.Replace(encoded, []byte(`"evidence":[{"ref":"logical-result","content_hash":"`+strings.Repeat("a", 64)+`"}]`), []byte(`"evidence":[]`), 1)
	cases["graph altered"] = bytes.Replace(encoded, []byte(request.ImpactInput.Graph.GraphHash), []byte(strings.Repeat("b", 64)), 1)
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input.json")
			output := filepath.Join(dir, "output.json")
			if err := os.WriteFile(input, data, 0600); err != nil {
				t.Fatal(err)
			}
			prior := []byte("keep historical ledger")
			if err := os.WriteFile(output, prior, 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"trace", "review", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
				t.Fatal("invalid receipt emitted result")
			}
			actual, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(actual, prior) {
				t.Fatal("invalid review replaced output")
			}
		})
	}
}

func TestIdeaTraceReviewProtectsAliasesAndFIFO(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	original := writeTraceFixture(t, input, traceReviewCLIFixture(t))
	alias := filepath.Join(dir, "hardlink.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "symlink.json")
	if err := os.Symlink(input, link); err != nil {
		t.Fatal(err)
	}
	for _, output := range []string{input, alias, link} {
		var stdout bytes.Buffer
		if err := RunIdeaCommand([]string{"trace", "review", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
			t.Fatal("protected output accepted")
		}
	}
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review", "--input", fifo}, &stdout); err == nil {
		t.Fatal("FIFO accepted")
	}
	actual, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(actual, original) {
		t.Fatal("input damaged")
	}
}
