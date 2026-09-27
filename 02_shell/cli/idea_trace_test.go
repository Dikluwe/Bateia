/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_trace_test.md
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
)

func traceBuildFixture() domain.DecisionTraceBuildInput {
	return domain.DecisionTraceBuildInput{
		SchemaVersion: "1",
		Report:        domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown},
	}
}

func writeTraceFixture(t *testing.T, path string, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return data
}

func TestIdeaTraceBuildAndImpactOffline(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "build.json")
	original := writeTraceFixture(t, input, traceBuildFixture())
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "build", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var graph domain.DecisionTraceGraph
	if err := json.Unmarshal(stdout.Bytes(), &graph); err != nil {
		t.Fatal(err)
	}
	if graph.SchemaVersion != domain.DecisionTraceSchemaVersion || graph.Scope != domain.DecisionTraceScope || len(graph.GraphHash) != 64 || len(graph.Nodes) < 3 {
		t.Fatalf("invalid graph: %+v", graph)
	}
	output := filepath.Join(dir, "graph.json")
	if err := os.WriteFile(output, []byte("old result"), 0600); err != nil {
		t.Fatal(err)
	}
	var silent bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "build", "--input=" + input, "--output=" + output}, &silent); err != nil {
		t.Fatal(err)
	}
	onDisk, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(onDisk, stdout.Bytes()) || silent.Len() != 0 {
		t.Fatalf("disk/stdout mismatch or unexpected stdout: %v", err)
	}
	change := domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: graph.Nodes[0].ID, Reason: "CONTENT_CHANGED"}}}
	changePath := filepath.Join(dir, "change.json")
	writeTraceFixture(t, changePath, change)
	stdout.Reset()
	if err := RunIdeaCommand([]string{"trace", "impact", "--input", changePath}, &stdout); err != nil {
		t.Fatal(err)
	}
	var plan domain.DecisionTraceImpactPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.DecisionBefore != domain.DecisionUnknown || plan.DecisionAfter != domain.DecisionUnknown || len(plan.ChangedNodeIDs) != 1 || plan.GraphHash != graph.GraphHash {
		t.Fatalf("impact must not recalculate a decision: %+v", plan)
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(after, original) {
		t.Fatal("build input was modified")
	}
	temporary, err := filepath.Glob(filepath.Join(dir, ".bateia-trace-*"))
	if err != nil || len(temporary) != 0 {
		t.Fatalf("temporary results leaked: %v %v", temporary, err)
	}
}

func TestIdeaTraceRejectsAmbiguousArguments(t *testing.T) {
	cases := [][]string{
		{}, {"wat"}, {"build"}, {"build", "--input"}, {"build", "--input="},
		{"build", "--input", "a", "--input", "b"},
		{"build", "--input", "a", "--output", "x", "--output", "y"},
		{"build", "--input", "a", "positional"},
		{"build", "--input", "a", "--unknown", "x"},
		{"build", "--input", "--output", "x"},
		{"build", "--input", "a", "--output="},
	}
	for index, args := range cases {
		var out bytes.Buffer
		if err := RunIdeaCommand(append([]string{"trace"}, args...), &out); err == nil || out.Len() != 0 {
			t.Fatalf("case %d accepted ambiguous arguments: %v", index, args)
		}
	}
}

func TestIdeaTraceRejectsInvalidInputsWithoutOutput(t *testing.T) {
	cases := map[string]string{
		"unknown":   `{"schema_version":"1","report":{"decision":"INCONCLUSIVE"},"unknown":true}`,
		"duplicate": `{"schema_version":"1","schema_version":"1","report":{"decision":"INCONCLUSIVE"}}`,
		"casefold":  `{"schema_version":"1","SCHEMA_VERSION":"1","report":{"decision":"INCONCLUSIVE"}}`,
		"trailing":  `{"schema_version":"1","report":{"decision":"INCONCLUSIVE"}} {}`,
		"schema":    `{"schema_version":"999","report":{"decision":"INCONCLUSIVE"}}`,
		"decision":  `{"schema_version":"1","report":{"decision":"MAGIC"}}`,
		"depth":     strings.Repeat("[", 65) + "0" + strings.Repeat("]", 65),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			input, output := filepath.Join(dir, "input.json"), filepath.Join(dir, "output.json")
			if err := os.WriteFile(input, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"trace", "build", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("invalid input accepted")
			}
			if stdout.Len() != 0 {
				t.Fatal("invalid input produced stdout")
			}
			if _, err := os.Lstat(output); !os.IsNotExist(err) {
				t.Fatalf("invalid input produced output: %v", err)
			}
		})
	}
}

func TestIdeaTraceRejectsOutputAliases(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.json")
	original := writeTraceFixture(t, input, traceBuildFixture())
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
	parentAlias := filepath.Join(dir, "parent")
	if err := os.Symlink(dir, parentAlias); err != nil {
		t.Fatal(err)
	}
	for _, output := range []string{input, hard, symbolic, fifo, dir, filepath.Join(parentAlias, "input.json"), filepath.Join(parentAlias, "missing.json")} {
		var stdout bytes.Buffer
		if err := RunIdeaCommand([]string{"trace", "build", "--input", input, "--output", output}, &stdout); err == nil || stdout.Len() != 0 {
			t.Fatalf("unsafe output accepted: %s", output)
		}
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(original, after) {
		t.Fatal("input changed during failed output admission")
	}
	if _, err := os.Stat(filepath.Join(dir, "missing.json")); !os.IsNotExist(err) {
		t.Fatal("aliased output was created")
	}
}

func TestIdeaTraceInputFIFOAndTamperedGraphFailClosed(t *testing.T) {
	dir := t.TempDir()
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "build", "--input", fifo}, &out); err == nil || out.Len() != 0 {
		t.Fatal("FIFO input accepted")
	}
	input := filepath.Join(dir, "build.json")
	writeTraceFixture(t, input, traceBuildFixture())
	if err := RunIdeaCommand([]string{"trace", "build", "--input", input}, &out); err != nil {
		t.Fatal(err)
	}
	var graph domain.DecisionTraceGraph
	if err := json.Unmarshal(out.Bytes(), &graph); err != nil {
		t.Fatal(err)
	}
	graph.GraphHash = strings.Repeat("0", 64)
	impactInput := filepath.Join(dir, "impact.json")
	writeTraceFixture(t, impactInput, domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: graph.Nodes[0].ID, Reason: "CONTENT_CHANGED"}}})
	output := filepath.Join(dir, "output.json")
	const originalOutput = "previous result"
	if err := os.WriteFile(output, []byte(originalOutput), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := RunIdeaCommand([]string{"trace", "impact", "--input", impactInput, "--output", output}, &out); err == nil || out.Len() != 0 {
		t.Fatal("tampered graph accepted")
	}
	after, err := os.ReadFile(output)
	if err != nil || string(after) != originalOutput {
		t.Fatal("invalid impact replaced previous result")
	}
}

func TestDecisionTraceWriterDoesNotExportConfiguredSecret(t *testing.T) {
	const secret = "synthetic-secret-no-real-token"
	t.Setenv("GITHUB_TOKEN", secret)
	var stdout bytes.Buffer
	if err := writeDecisionTraceJSON(map[string]string{"label": secret}, "", &stdout); err == nil || stdout.Len() != 0 {
		t.Fatal("secret reached stdout")
	}
	output := filepath.Join(t.TempDir(), "output.json")
	if err := writeDecisionTraceJSON(map[string]string{"label": secret}, output, &stdout); err == nil {
		t.Fatal("secret artifact accepted")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("secret artifact created")
	}
}

func TestDecisionTracePreflightProtectsFutureOutputs(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "build.json")
	writeTraceFixture(t, input, traceBuildFixture())
	future := filepath.Join(dir, "new", "nested", "graph.json")
	if err := validateDecisionTraceOutput(future, input); err != nil {
		t.Fatalf("future dossier should be admissible before creation: %v", err)
	}
	if err := validateDecisionTraceOutput(future, future); err == nil {
		t.Fatal("future output alias accepted")
	}
	alias := filepath.Join(dir, "alias")
	if err := os.Symlink(dir, alias); err != nil {
		t.Fatal(err)
	}
	if err := validateDecisionTraceOutput(future, filepath.Join(alias, "new", "nested", "graph.json")); err == nil {
		t.Fatal("future output ancestor alias accepted")
	}
	dangling := filepath.Join(dir, "dangling-output.json")
	if err := os.Symlink(future, dangling); err != nil {
		t.Fatal(err)
	}
	if err := validateDecisionTraceOutput(future, dangling); err == nil {
		t.Fatal("dangling protected leaf alias accepted")
	}
	danglingParent := filepath.Join(dir, "dangling-parent")
	if err := os.Symlink(filepath.Join(dir, "new", "nested"), danglingParent); err != nil {
		t.Fatal(err)
	}
	if err := validateDecisionTraceOutput(future, filepath.Join(danglingParent, "graph.json")); err == nil {
		t.Fatal("dangling protected parent alias accepted")
	}
	if err := writeDecisionTraceJSON(map[string]string{"scope": "declared"}, future, nil, input); err == nil {
		t.Fatal("writer should not create unadmitted parent directories")
	}
	if _, err := os.Stat(filepath.Join(dir, "new")); !os.IsNotExist(err) {
		t.Fatal("writer created missing dossier")
	}
	output := filepath.Join(dir, "trace.json")
	if err := os.WriteFile(output, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeDecisionTraceJSON(map[string]string{"scope": "declared"}, output, nil, input, filepath.Join(dir, "normal-output.json")); err != nil {
		t.Fatalf("existing trace and distinct future normal output should work: %v", err)
	}
}

func TestIdeaTraceImpactRejectsUnknownReason(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "build.json")
	writeTraceFixture(t, input, traceBuildFixture())
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "build", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var graph domain.DecisionTraceGraph
	if err := json.Unmarshal(stdout.Bytes(), &graph); err != nil {
		t.Fatal(err)
	}
	impact := filepath.Join(dir, "impact.json")
	writeTraceFixture(t, impact, domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: graph.Nodes[0].ID, Reason: "approve automatically"}}})
	stdout.Reset()
	if err := RunIdeaCommand([]string{"trace", "impact", "--input", impact}, &stdout); err == nil || stdout.Len() != 0 {
		t.Fatal("unknown change reason accepted")
	}
}

func TestDecisionTraceWriterRefusesInvalidValueAndUnavailableStdout(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "trace.json")
	if err := os.WriteFile(output, []byte("previous"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeDecisionTraceJSON(make(chan int), output, nil); err == nil {
		t.Fatal("unserializable trace accepted")
	}
	data, err := os.ReadFile(output)
	if err != nil || string(data) != "previous" {
		t.Fatal("unserializable trace replaced output")
	}
	if err := writeDecisionTraceJSON(map[string]string{"scope": "declared"}, "", nil); err == nil {
		t.Fatal("nil stdout accepted")
	}
}
