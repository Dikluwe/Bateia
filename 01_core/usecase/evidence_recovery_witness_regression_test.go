/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_recovery_witness_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"reflect"
	"testing"
)

func TestRecoveryMetadataDoesNotWitnessDerivedFiles(t *testing.T) {
	for _, kind := range []domain.EvidenceKind{domain.EvidenceFile, domain.EvidenceTest} {
		for _, treePath := range []string{"", "unrelated.go", "graph_test.go"} {
			t.Run(string(kind)+"/"+treePath, func(t *testing.T) {
				ref := evidenceIntegrityPlan().Repository
				cause := domain.CauseMissingImplementation
				if kind == domain.EvidenceTest {
					cause = domain.CauseMissingTest
				}
				d := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "graph", Candidate: ref.Key(), Capability: "graph", Cause: cause, Recoverable: true, RuleID: "graph-rule"}}}
				plan, err := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, nil)
				if err != nil || len(plan.Targets) != 1 {
					t.Fatal(err, plan)
				}
				p := &completenessBatchPort{content: map[string]string{"graph_test.go": "package fixture\nfunc TestGraph(){}"}}
				if treePath != "" {
					p.tree = []domain.DirectoryNode{{Path: treePath, Type: "file"}}
				}
				got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "graph", Capabilities: []string{"graph"}}, domain.ScreeningReport{}, nil, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), decisionPolicy())
				if err != nil {
					t.Fatal(err)
				}
				want := 0
				if treePath == "graph_test.go" {
					want = 1
				}
				history := ExecutedRecoveryTargets(plan, got.Bundles)
				next, err := PlanEvidenceRecovery(d, plan.Policy, 2, history)
				if err != nil || got.Metrics.TargetsExecuted != want || len(history) != want || len(next.Targets) != 1-want {
					t.Errorf("metadata fabricated execution/history: metrics=%+v history=%+v next=%+v err=%v", got.Metrics, history, next, err)
				}
				if want == 0 && (len(got.Bundles[0].Targets) != 0 || !reflect.DeepEqual(p.calls, []string{"tree", "metadata", "tree"})) {
					t.Fatal("unexpected file attempt", got.Bundles[0].Targets, p.calls)
				}
			})
		}
	}
}

func TestRecoveryWitnessFallbackKeepsRelevantItemsAndExplicitResults(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	plan := budgetPlan(ref)
	plan.Targets[0].Path, plan.Targets[0].Kind, plan.Targets[0].Patterns = "", domain.EvidenceFile, []string{"graph"}
	for _, tc := range []struct {
		name, path string
		kind       domain.EvidenceKind
		want       int
	}{
		{"metadata is not code", "@repository-metadata", domain.EvidenceFile, 0},
		{"metadata is not test", "@repository-metadata", domain.EvidenceTest, 0},
		{"legacy code", "graph.go", domain.EvidenceFile, 1},
		{"legacy test", "graph_test.go", domain.EvidenceTest, 1},
		{"explicit metadata", "@repository-metadata", domain.EvidenceRepository, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := plan
			p.Targets = append([]domain.RecoveryTarget(nil), plan.Targets...)
			p.Targets[0].Kind = tc.kind
			if tc.kind == domain.EvidenceRepository {
				p.Targets[0].Path = tc.path
			}
			b := domain.EvidenceBundle{Repository: ref, Items: []domain.CollectedEvidence{{Path: tc.path}}}
			if got := ExecutedRecoveryTargets(p, []domain.EvidenceBundle{b}); len(got) != tc.want {
				t.Errorf("wrong witness: %+v", got)
			}
		})
	}
	for _, status := range []domain.CollectionTargetStatus{domain.TargetCollected, domain.TargetNotFound, domain.TargetFailed, domain.TargetSkippedBudget, domain.TargetUnsupported} {
		b := domain.EvidenceBundle{Repository: ref, Targets: []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Path: "graph.go", Kind: domain.EvidenceFile, Reason: "recovery 0"}, Status: status}}}
		want := 0
		if status == domain.TargetCollected || status == domain.TargetNotFound || status == domain.TargetFailed {
			want = 1
		}
		if got := ExecutedRecoveryTargets(plan, []domain.EvidenceBundle{b}); len(got) != want {
			t.Fatal("explicit status changed", status, got)
		}
	}
}
