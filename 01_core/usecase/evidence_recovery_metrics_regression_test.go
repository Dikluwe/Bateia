/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_recovery_metrics_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestRecoveryDerivedMetricsAggregatePathsWithoutBlockingRetry(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: strings.Repeat("a", 40)}
	candidate := domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit}
	metadata, err := json.Marshal(candidate)
	if err != nil {
		t.Fatal(err)
	}
	first, second := "package fixture\nfunc Graph(){}\n", "package fixture\nfunc GraphOther(){}\n"
	for _, tc := range []struct {
		name         string
		fileBytes    int
		quota        int
		wantExecuted int
		wantFirst    string
		wantSecond   string
		firstStatus  domain.CollectionTargetStatus
		secondStatus domain.CollectionTargetStatus
	}{
		{"whole omission", 0, 2, 0, "", "", domain.TargetSkippedBudget, domain.TargetSkippedBudget},
		{"one byte below", len(first) - 1, 2, 1, first[:len(first)-1], "", domain.TargetSkippedBudget, domain.TargetSkippedBudget},
		{"exact first file", len(first), 2, 1, first, "", domain.TargetCollected, domain.TargetSkippedBudget},
		{"one byte above", len(first) + 1, 2, 1, first, second[:1], domain.TargetCollected, domain.TargetSkippedBudget},
		{"both files", len(first) + len(second), 2, 1, first, second, domain.TargetCollected, domain.TargetCollected},
		{"quantity omission", len(first) + len(second), 0, 0, "", "", domain.TargetSkippedBudget, domain.TargetSkippedBudget},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policy := DefaultRecoveryPolicy()
			policy.AdditionalCodeFiles = tc.quota
			policy.AdditionalBytes = int64(len(metadata) + tc.fileBytes)
			diagnosis := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "implementation", Candidate: ref.Key(), Capability: "dependency graph", Cause: domain.CauseMissingImplementation, Recoverable: true, RuleID: "dependency-graph"}}}
			plan, err := PlanEvidenceRecovery(diagnosis, policy, 1, nil)
			if err != nil || len(plan.Targets) != 1 || plan.Targets[0].Path != "" {
				t.Fatalf("derived fixture: %v %+v", err, plan)
			}
			// Preserve an original item matching the pattern: it cannot count as new execution.
			old := "historical graph bytes"
			oldHash := fmt.Sprintf("%x", sha256.Sum256([]byte(old)))
			original := []domain.EvidenceBundle{{SchemaVersion: "1", Repository: ref, Plan: domain.EvidenceExtractionPlan{Candidate: candidate}, Items: []domain.CollectedEvidence{{Ref: ref, Kind: domain.EvidenceFile, Path: "graph_old.go", Content: old, ContentHash: oldHash}}, BytesRead: int64(len(old))}}
			before, _ := json.Marshal(original)
			planBefore, _ := json.Marshal(plan)
			p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "graph_a.go", Type: "file"}, {Path: "graph_z.go", Type: "file"}}, content: map[string]string{"graph_a.go": first, "graph_z.go": second}}
			got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, original, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
			if err != nil || len(got.Bundles) != 1 {
				t.Fatalf("execute: %v %+v", err, got)
			}
			if got.Metrics.TargetsPlanned != 1 || got.Metrics.TargetsExecuted != tc.wantExecuted {
				t.Errorf("lost aggregate execution: %+v", got.Metrics)
			}
			wantCalls := 4 // preparation tree, metadata, collection tree, extraction
			if tc.quota == 0 {
				wantCalls = 3
			}
			wantBytes := int64(len(metadata) + len(tc.wantFirst) + len(tc.wantSecond))
			if len(p.calls) != wantCalls || got.Metrics.ProviderCalls != wantCalls || got.Metrics.AdditionalBytes != wantBytes || wantBytes > policy.AdditionalBytes {
				t.Fatalf("cost changed: calls=%v metrics=%+v wantBytes=%d", p.calls, got.Metrics, wantBytes)
			}
			b := got.Bundles[0]
			for i, path := range []string{"graph_a.go", "graph_z.go"} {
				want := []string{tc.wantFirst, tc.wantSecond}[i]
				status := []domain.CollectionTargetStatus{tc.firstStatus, tc.secondStatus}[i]
				result := completenessTarget(t, b, path)
				if result.Status != status || result.Bytes != int64(len(want)) {
					t.Fatalf("wrong result: %+v want=%s/%d", result, status, len(want))
				}
				found := false
				for _, item := range b.Items {
					if item.Path == path {
						found = true
						if item.Content != want || item.ContentHash != fmt.Sprintf("%x", sha256.Sum256([]byte(want))) {
							t.Fatal("admitted bytes/hash changed", item)
						}
					}
				}
				if found != (want != "") {
					t.Fatal("omission created item or lost admitted content", path)
				}
			}
			history := ExecutedRecoveryTargets(plan, got.Bundles)
			next, err := PlanEvidenceRecovery(diagnosis, policy, 2, history)
			if err != nil {
				t.Fatal(err)
			}
			if tc.secondStatus == domain.TargetSkippedBudget {
				if len(history) != 0 || len(next.Targets) != 1 {
					t.Fatal("aggregate metric blocked retry", history, next)
				}
			} else if len(history) != 1 || len(next.Targets) != 0 {
				t.Fatal("complete history lost", history, next)
			}
			after, _ := json.Marshal(original)
			planAfter, _ := json.Marshal(plan)
			if string(before) != string(after) || string(planBefore) != string(planAfter) {
				t.Fatal("historical input or plan mutated")
			}
		})
	}
}
