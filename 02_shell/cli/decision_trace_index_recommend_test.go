/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_index_recommend_test.md
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
	"bateia/01_core/usecase"
)

func TestRecommendTraceIndexMatchesOriginalBuildAndStageHashes(t *testing.T) {
	ideaPath, policyPath, source := recommendReplayFixture(t)
	output := filepath.Join(t.TempDir(), "out")
	args := []string{"--idea", ideaPath, "--policy", policyPath, "--mode", "replay", "--strategy", "literal", "--replay", source, "--output", output, "--trace"}
	if err := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	var idea domain.IdeaSpec
	var report domain.ReuseDecisionReport
	var graph domain.DecisionTraceGraph
	var index domain.DecisionTraceExplanation
	var manifest domain.IdeaRunManifest
	for name, value := range map[string]any{"idea.json": &idea, "reuse-decision.json": &report, "trace-graph.json": &graph, "trace-index.json": &index, "run.json": &manifest} {
		if err := readJSONFile(filepath.Join(output, name), value); err != nil {
			t.Fatal(err)
		}
	}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Idea: &idea, Report: report, RuleVersions: manifest.RuleVersions}
	original, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	graphBytes, err := os.ReadFile(filepath.Join(output, "trace-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(append(encoded, '\n'), graphBytes) || graph.SchemaVersion != "decision-trace-v1" {
		t.Fatal("index changed graph-v1 bytes")
	}
	expected, err := usecase.ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: original, BuildInput: build})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(index, expected) {
		t.Fatal("index used other build declarations")
	}
	stage := manifest.Stages[len(manifest.Stages)-1]
	if stage.Name != "trace" || stage.Status != domain.IdeaStageComplete || stage.ContractVersion != "decision-trace-projection-v2" || len(stage.Outputs) != 2 {
		t.Fatalf("bad trace stage %+v", stage)
	}
	for i, path := range []string{"trace-graph.json", "trace-index.json"} {
		hash, err := fileHash(filepath.Join(output, path))
		if err != nil {
			t.Fatal(err)
		}
		if stage.Outputs[i].Path != path || stage.Outputs[i].Hash != hash {
			t.Fatalf("bad hash %+v", stage.Outputs[i])
		}
	}
	if stage.Counts["entries"] != int64(len(index.Entries)) {
		t.Fatal("missing index count")
	}
	// Historical index without graph is still outside a resumed, untraced run.
	if err := os.Remove(filepath.Join(output, "trace-graph.json")); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(output, "trace-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	args = args[:len(args)-1]
	if err := RunIdeaRecommendCommand(append(args, "--resume"), &stdout, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "não pertence à execução atual") {
		t.Fatal("orphan historical index did not warn")
	}
	after, err := os.ReadFile(filepath.Join(output, "trace-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("resume overwrote historical index")
	}
	if err := readJSONFile(filepath.Join(output, "run.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	for _, stage := range manifest.Stages {
		if stage.Name == "trace" {
			t.Fatal("untraced run consumed index")
		}
	}
}

func TestRecommendTraceIndexPreflightRejectsBothAliasesBeforeAnyWrite(t *testing.T) {
	for _, leaf := range []string{"trace-graph.json", "trace-index.json"} {
		t.Run(leaf, func(t *testing.T) {
			idea, policy, source := recommendReplayFixture(t)
			out := t.TempDir()
			if err := os.Link(idea, filepath.Join(out, leaf)); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(idea)
			if err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			args := []string{"--idea", idea, "--policy", policy, "--mode", "replay", "--replay", source, "--output", out, "--resume", "--trace"}
			if err := RunIdeaRecommendCommand(args, &stdout, nil, nil, nil); err == nil {
				t.Fatal("accepted sidecar/input alias")
			}
			if stdout.Len() != 0 {
				t.Fatal("bad preflight emitted bytes")
			}
			if _, err := os.Stat(filepath.Join(out, "run.json")); !os.IsNotExist(err) {
				t.Fatal("bad preflight wrote manifest")
			}
			after, err := os.ReadFile(idea)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, after) {
				t.Fatal("bad preflight overwrote idea")
			}
		})
	}
}
