/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/decision_trace_explain_wiring_test.md
 * @layer L4
 * @updated 2026-09-26
 */
package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func TestDecisionTraceExplainWiringDoesNotLoadEnv(t *testing.T) {
	const helper = "BATEIA_TRACE_EXPLAIN_WIRING_HELPER"
	if os.Getenv(helper) == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				os.Args = append([]string{"bateia"}, os.Args[i+1:]...)
				main()
				os.Exit(0)
			}
		}
		os.Exit(2)
	}
	dir := t.TempDir()
	const marker = "synthetic-trace-explain-env-marker-0099"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("GITHUB_TOKEN="+marker+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Reasons: []domain.DecisionReason{{RuleID: marker, Message: "recorded declaration"}}}}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	request := domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "explain.json")
	if err := os.WriteFile(input, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDecisionTraceExplainWiringDoesNotLoadEnv$", "--", "idea", "trace", "explain", "--input", input)
	command.Dir = dir
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key != helper && key != "GITHUB_TOKEN" {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, helper+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("offline explain failed: %v: %s", err, output)
	}
	var index domain.DecisionTraceExplanation
	if err := json.Unmarshal(output, &index); err != nil {
		t.Fatalf("invalid JSON: %v: %s", err, output)
	}
	if index.GraphHash != graph.GraphHash || !strings.Contains(string(output), marker) {
		t.Fatalf("offline declaration filtered: %+v", index)
	}
}
