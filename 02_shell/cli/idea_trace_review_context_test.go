/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace_review_context_test.md
 * @layer L2
 * @updated 2026-09-27
 */
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func traceReviewContextCLIFixture(t *testing.T) domain.DecisionTraceReviewContextInput {
	t.Helper()
	build := traceBuildFixture()
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	var reportID string
	for _, node := range graph.Nodes {
		if node.Kind == domain.TraceInput && node.Label == "report" {
			reportID = node.ID
		}
	}
	if reportID == "" {
		t.Fatal("fixture requires report input")
	}
	return domain.DecisionTraceReviewContextInput{SchemaVersion: "1", BuildInput: build, ReviewInput: domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: reportID, Reason: "INPUT_CHANGED"}}}}}
}

func TestIdeaTraceReviewContextBindsOriginalDocuments(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "request.json")
	output := filepath.Join(dir, "context.json")
	request := traceReviewContextCLIFixture(t)
	original := writeTraceFixture(t, input, request)
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review-context", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var context domain.DecisionTraceReviewContext
	if err := json.Unmarshal(stdout.Bytes(), &context); err != nil {
		t.Fatal(err)
	}
	ledger, err := usecase.TrackDecisionTraceReview(request.ReviewInput)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(context.Ledger, ledger) || context.PlanHash != ledger.PlanHash || context.GraphHash != ledger.GraphHash || len(context.TaskContexts) != 2 || len(context.Entries) != 2 {
		t.Fatal("context changed review ledger or projected unrelated nodes")
	}
	if err := os.WriteFile(output, []byte("old context"), 0600); err != nil {
		t.Fatal(err)
	}
	var silent bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review-context", "--input=" + input, "--output=" + output}, &silent); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(actual, stdout.Bytes()) || silent.Len() != 0 {
		t.Fatal("context disk/stdout mismatch")
	}
	unchanged, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(unchanged, original) {
		t.Fatal("context modified source input")
	}
	temporary, err := filepath.Glob(filepath.Join(dir, ".bateia-trace-*"))
	if err != nil || len(temporary) != 0 {
		t.Fatal("context temporary outputs leaked")
	}
}

func TestIdeaTraceReviewContextRejectsInvalidBeforeOutput(t *testing.T) {
	request := traceReviewContextCLIFixture(t)
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{}
	for name, addition := range map[string]string{"unknown": `,"node_ids":[]}`, "duplicate": `,"schema_version":"1"}`, "casefold": `,"SCHEMA_VERSION":"1"}`} {
		cases[name] = append(append([]byte{}, encoded[:len(encoded)-1]...), []byte(addition)...)
	}
	cases["trailing"] = append(append([]byte{}, encoded...), []byte(` {}`)...)
	cases["missing hash"] = bytes.Replace(encoded, []byte(request.ReviewInput.ImpactInput.Graph.GraphHash), []byte(""), 1)
	cases["schema"] = bytes.Replace(encoded, []byte(`"schema_version":"1"`), []byte(`"schema_version":"2"`), 1)
	altered := traceReviewContextCLIFixture(t)
	altered.BuildInput.Report.Decision = domain.DecisionBuild
	cases["altered original"], err = json.Marshal(altered)
	if err != nil {
		t.Fatal(err)
	}
	seed, err := usecase.TrackDecisionTraceReview(request.ReviewInput)
	if err != nil {
		t.Fatal(err)
	}
	const secret = "synthetic-review-context-configured-secret-0102"
	request.ReviewInput.Records = []domain.DecisionTraceReviewRecord{{NodeID: seed.Tasks[0].Task.NodeID, PlanHash: seed.PlanHash, TaskHash: seed.Tasks[0].TaskHash, Status: "IN_PROGRESS", Note: secret}}
	cases["configured secret"], err = json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "configured secret" {
				t.Setenv("GITHUB_TOKEN", secret)
			}
			dir := t.TempDir()
			input := filepath.Join(dir, "input.json")
			output := filepath.Join(dir, "output.json")
			if err := os.WriteFile(input, data, 0600); err != nil {
				t.Fatal(err)
			}
			prior := []byte("keep original context")
			if err := os.WriteFile(output, prior, 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"trace", "review-context", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
				t.Fatal("invalid context emitted output")
			}
			actual, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(actual, prior) {
				t.Fatal("invalid context replaced output")
			}
		})
	}
}

func TestIdeaTraceReviewContextProtectsInputAliasesAndFIFO(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	original := writeTraceFixture(t, input, traceReviewContextCLIFixture(t))
	alias := filepath.Join(dir, "alias.json")
	link := filepath.Join(dir, "symlink.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(input, link); err != nil {
		t.Fatal(err)
	}
	for _, output := range []string{input, alias, link} {
		var stdout bytes.Buffer
		if err := RunIdeaCommand([]string{"trace", "review-context", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
			t.Fatal("unsafe output accepted")
		}
	}
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "review-context", "--input", fifo}, &stdout); err == nil {
		t.Fatal("FIFO admitted")
	}
	actual, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(actual, original) {
		t.Fatal("source input damaged")
	}
	if err := RunIdeaCommand([]string{"trace", "review-context", "--input", input, "--input", input}, &stdout); err == nil {
		t.Fatal("duplicate argument accepted")
	}
	if strings.Contains(stdout.String(), "schema_version") {
		t.Fatal("failed request emitted JSON")
	}
}
