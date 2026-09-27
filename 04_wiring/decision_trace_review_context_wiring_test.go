/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/decision_trace_review_context_wiring_test.md
 * @layer L4
 * @updated 2026-09-27
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

func TestDecisionTraceReviewContextWiringDoesNotLoadEnv(t *testing.T) {
	const helper = "BATEIA_TRACE_REVIEW_CONTEXT_WIRING_HELPER"
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
	const marker = "synthetic-trace-review-context-env-marker-0102"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("GITHUB_TOKEN="+marker+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown}}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	request := domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: graph.Nodes[0].ID, Reason: "INPUT_CHANGED"}}}}
	seed, err := usecase.TrackDecisionTraceReview(request)
	if err != nil {
		t.Fatal(err)
	}
	request.Records = []domain.DecisionTraceReviewRecord{{NodeID: seed.Tasks[0].Task.NodeID, PlanHash: seed.PlanHash, TaskHash: seed.Tasks[0].TaskHash, Status: "IN_PROGRESS", Note: marker}}
	encoded, err := json.Marshal(domain.DecisionTraceReviewContextInput{SchemaVersion: "1", ReviewInput: request, BuildInput: build})
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "review.json")
	if err := os.WriteFile(input, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDecisionTraceReviewContextWiringDoesNotLoadEnv$", "--", "idea", "trace", "review-context", "--input", input)
	command.Dir = dir
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key != helper && key != "GITHUB_TOKEN" && key != "SOURCEGRAPH_TOKEN" {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, helper+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("offline review context failed: %v %s", err, output)
	}
	var projection domain.DecisionTraceReviewContext
	if err := json.Unmarshal(output, &projection); err != nil {
		t.Fatalf("invalid JSON: %v %s", err, output)
	}
	if projection.PlanHash != seed.PlanHash || projection.GraphHash != graph.GraphHash || !strings.Contains(string(output), marker) || projection.Ledger.DecisionImpact != "NONE" {
		t.Fatal("offline receipt not preserved")
	}
}
