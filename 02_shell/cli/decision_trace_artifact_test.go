/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_artifact_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestRecommendTraceIsOptInAndPreservesCausalArtifacts(t *testing.T) {
	idea, policy, source := recommendReplayFixture(t)
	plain, traced := filepath.Join(t.TempDir(), "plain"), filepath.Join(t.TempDir(), "traced")
	var lastStdout bytes.Buffer
	run := func(output string, extra ...string) domain.IdeaRunManifest {
		t.Helper()
		args := []string{"--idea", idea, "--policy", policy, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", output}
		lastStdout.Reset()
		if err := RunIdeaRecommendCommand(append(args, extra...), &lastStdout, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
		var manifest domain.IdeaRunManifest
		if err := readJSONFile(filepath.Join(output, "run.json"), &manifest); err != nil {
			t.Fatal(err)
		}
		return manifest
	}
	a, b := run(plain), run(traced, "--trace")
	if a.RunID != b.RunID || !reflect.DeepEqual(a.InputHashes, b.InputHashes) || !reflect.DeepEqual(a.RuleVersions, b.RuleVersions) {
		t.Fatal("trace changed causal run identity")
	}
	for _, path := range []string{"reuse-decision.json", "decision.md", "assessment-report.json", "alternatives.json", "evidence/bundles.json"} {
		x, err := os.ReadFile(filepath.Join(plain, path))
		if err != nil {
			t.Fatal(err)
		}
		y, err := os.ReadFile(filepath.Join(traced, path))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(x, y) {
			t.Fatalf("trace changed %s", path)
		}
	}
	if _, err := os.Stat(filepath.Join(plain, "trace-graph.json")); !os.IsNotExist(err) {
		t.Fatal("default run emitted trace")
	}
	if _, err := os.Stat(filepath.Join(plain, "trace-index.json")); !os.IsNotExist(err) {
		t.Fatal("default run emitted trace index")
	}
	var graph domain.DecisionTraceGraph
	if err := readJSONFile(filepath.Join(traced, "trace-graph.json"), &graph); err != nil {
		t.Fatal(err)
	}
	if graph.Scope != domain.DecisionTraceScope || graph.GraphHash == "" {
		t.Fatalf("invalid trace %+v", graph)
	}
	if len(b.Stages) != len(a.Stages)+1 {
		t.Fatal("missing optional trace stage")
	}
	stage := b.Stages[len(b.Stages)-1]
	hash, err := fileHash(filepath.Join(traced, "trace-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	if stage.Name != "trace" || stage.Outputs[0].Hash != hash || stage.Counts["nodes"] != int64(len(graph.Nodes)) {
		t.Fatalf("invalid stage %+v", stage)
	}
	var index domain.DecisionTraceExplanation
	indexPath := filepath.Join(traced, "trace-index.json")
	if err := readJSONFile(indexPath, &index); err != nil {
		t.Fatal(err)
	}
	indexHash, err := fileHash(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	locations := 0
	for _, entry := range index.Entries {
		locations += len(entry.Locations)
	}
	if index.GraphHash != graph.GraphHash || len(index.Entries) != len(graph.Nodes) || len(stage.Outputs) != 2 || stage.Outputs[1].Path != "trace-index.json" || stage.Outputs[1].Hash != indexHash || stage.Counts["entries"] != int64(len(index.Entries)) || stage.Counts["locations"] != int64(locations) {
		t.Fatalf("unbound or unrecorded trace index: %+v / %+v", index, stage)
	}
	indexBefore, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	// Expiry asks for review; the historical decision and graph are not rewritten.
	before, err := os.ReadFile(filepath.Join(traced, "trace-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	var sourceID string
	for _, node := range graph.Nodes {
		if node.Kind == domain.TraceSource {
			sourceID = node.ID
			break
		}
	}
	if sourceID == "" {
		t.Fatal("recorded report sources were omitted")
	}
	request := domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: sourceID, Reason: "EXPIRED"}}}
	input := filepath.Join(t.TempDir(), "impact.json")
	if err := writeJSONAtomic(input, request); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "impact", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var plan domain.DecisionTraceImpactPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.Status != "REVIEW_REQUIRED" || plan.DecisionBefore != graph.Decision || plan.DecisionAfter != graph.Decision || plan.DecisionImpact != "NONE" || len(plan.Tasks) == 0 {
		t.Fatalf("unsafe impact %+v", plan)
	}
	after, err := os.ReadFile(filepath.Join(traced, "trace-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("impact rewrote historical graph")
	}
	// A resumed run without opt-in must not advertise or consume an older graph.
	c := run(traced, "--resume")
	if !strings.Contains(lastStdout.String(), "não pertence à execução atual") || !strings.Contains(lastStdout.String(), "trace-index.json") {
		t.Fatal("resume did not warn about historical trace")
	}
	indexAfter, err := os.ReadFile(indexPath)
	if err != nil || !bytes.Equal(indexBefore, indexAfter) {
		t.Fatal("untraced resume rewrote historical trace index")
	}
	for _, stage := range c.Stages {
		if stage.Name == "trace" {
			t.Fatal("resume consumed old optional trace")
		}
	}
}

func TestDecideTraceSidecarDoesNotChangeReportAndProtectsInputs(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, value any) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := writeJSONAtomic(path, value); err != nil {
			t.Fatal(err)
		}
		return path
	}
	idea := write("idea.json", domain.IdeaSpec{Problem: "inspect", Capabilities: []string{"core"}})
	policy := write("policy.json", domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3})
	assessment := write("assessment.json", domain.AssessmentReport{SchemaVersion: "1", Assessments: []domain.CandidateAssessment{}})
	args := []string{"decide", "--idea", idea, "--policy", policy, "--assessment", assessment}
	var plain, traced bytes.Buffer
	if err := RunIdeaCommand(args, &plain); err != nil {
		t.Fatal(err)
	}
	trace := filepath.Join(dir, "trace.json")
	if err := RunIdeaCommand(append(append([]string(nil), args...), "--trace-output", trace), &traced); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plain.Bytes(), traced.Bytes()) {
		t.Fatal("sidecar changed report bytes")
	}
	var graph domain.DecisionTraceGraph
	if err := readJSONFile(trace, &graph); err != nil {
		t.Fatal(err)
	}
	if graph.Decision != domain.DecisionUnknown {
		t.Fatalf("trace endorsed %+v", graph)
	}
	output := filepath.Join(dir, "report.json")
	if err := RunIdeaCommand(append(append([]string(nil), args...), "--trace-output", trace, "--output", output), nil); err != nil {
		t.Fatal(err)
	}
	// Even the legacy report destination must not overwrite an input in opt-in mode.
	beforeIdea, err := os.ReadFile(idea)
	if err != nil {
		t.Fatal(err)
	}
	beforeTrace, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand(append(append([]string(nil), args...), "--trace-output", trace, "--output", idea), nil); err == nil {
		t.Fatal("opt-in report destination accepted an input")
	}
	afterIdea, err := os.ReadFile(idea)
	if err != nil {
		t.Fatal(err)
	}
	afterTrace, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(beforeIdea, afterIdea) || !bytes.Equal(beforeTrace, afterTrace) {
		t.Fatal("invalid report destination changed an input or trace")
	}
	alias := filepath.Join(dir, "alias.json")
	if err := os.Link(idea, alias); err != nil {
		t.Fatal(err)
	}
	for _, destination := range []string{idea, alias, output} {
		prior, err := os.ReadFile(destination)
		if err != nil {
			t.Fatal(err)
		}
		var stdout bytes.Buffer
		unsafe := append(append([]string(nil), args...), "--trace-output", destination, "--output", output)
		if err := RunIdeaCommand(unsafe, &stdout); err == nil {
			t.Fatalf("accepted alias %s", destination)
		}
		if stdout.Len() != 0 {
			t.Fatal("invalid sidecar emitted stdout")
		}
		after, err := os.ReadFile(destination)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(prior, after) {
			t.Fatal("invalid sidecar overwrote input")
		}
	}
}
