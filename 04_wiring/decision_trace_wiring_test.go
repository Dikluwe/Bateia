/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/decision_trace_wiring_test.md
 * @layer L4
 * @updated 2026-09-26
 */
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestDecisionTraceWiringDoesNotLoadEnv(t *testing.T) {
	const helperKey = "BATEIA_TRACE_WIRING_HELPER"
	if os.Getenv(helperKey) == "1" {
		for index, arg := range os.Args {
			if arg == "--" {
				os.Args = append([]string{"bateia"}, os.Args[index+1:]...)
				main()
				os.Exit(0)
			}
		}
		os.Exit(2)
	}
	dir := t.TempDir()
	// If .env is loaded, the trace writer's secret guard rejects this declaration.
	// The token is a synthetic test marker, not an account credential.
	const marker = "synthetic-offline-trace-env-marker-0098"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("GITHUB_TOKEN="+marker+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	request := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{
		SchemaVersion: "1", Decision: domain.DecisionUnknown,
		Reasons: []domain.DecisionReason{{RuleID: marker, Message: "offline dispatch check"}},
	}}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "build.json")
	if err := os.WriteFile(input, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDecisionTraceWiringDoesNotLoadEnv$", "--", "idea", "trace", "build", "--input", input)
	command.Dir = dir
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key != helperKey && key != "GITHUB_TOKEN" {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, helperKey+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("offline trace dispatch failed: %v: %s", err, output)
	}
	var graph domain.DecisionTraceGraph
	if err := json.Unmarshal(output, &graph); err != nil {
		t.Fatalf("non-JSON output from real wiring: %v: %s", err, output)
	}
	if graph.Scope != domain.DecisionTraceScope || graph.Decision != domain.DecisionUnknown || !bytes.Contains(output, []byte(marker)) {
		t.Fatalf("offline declaration was unexpectedly filtered: %+v", graph)
	}
}
