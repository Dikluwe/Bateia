/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_limits_profile_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type limitsProfilePort struct {
	evidenceIntegrityPort
	tree      []domain.DirectoryNode
	requested []string
	extra     []domain.ExtractedFile
	extracts  int
}

func (p *limitsProfilePort) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return p.tree, nil
}
func (p *limitsProfilePort) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.extracts++
	p.requested = append(p.requested, paths...)
	files := []domain.ExtractedFile{}
	for _, path := range paths {
		files = append(files, domain.ExtractedFile{Path: path, Content: "same content"})
	}
	files = append(files, p.extra...)
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}, Files: files}, nil
}
func TestExplicitCollectionRespectsPathQuotas(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		paths                    []string
		kinds                    []domain.EvidenceKind
		code, docs, tests, depth int
		want                     []string
		omitted                  bool
	}{
		{"zero", []string{"a.go"}, []domain.EvidenceKind{domain.EvidenceFile}, 0, 4, 4, 4, nil, true},
		{"one", []string{"a.go", "b.go"}, []domain.EvidenceKind{domain.EvidenceFile, domain.EvidenceFile}, 1, 4, 4, 4, []string{"a.go"}, true},
		{"depth", []string{"src/deep/a.go"}, []domain.EvidenceKind{domain.EvidenceFile}, 4, 4, 4, 1, nil, true},
		{"exact", []string{"a.go"}, []domain.EvidenceKind{domain.EvidenceFile}, 1, 4, 4, 1, []string{"a.go"}, false},
		{"duplicate", []string{"a.go", "a.go"}, []domain.EvidenceKind{domain.EvidenceFile, domain.EvidenceFile}, 1, 4, 4, 1, []string{"a.go"}, false},
		{"documents", []string{"README.md", "LICENSE"}, []domain.EvidenceKind{domain.EvidenceDocumentation, domain.EvidenceLicense}, 4, 1, 4, 4, []string{"README.md"}, true},
		{"tests examples", []string{"main_test.go", "examples/demo.go"}, []domain.EvidenceKind{domain.EvidenceTest, domain.EvidenceExample}, 4, 4, 1, 4, []string{"main_test.go"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := evidenceIntegrityPlan()
			plan.Targets = nil
			plan.RequirementQueries = nil
			plan.Budget.CodeFiles = tc.code
			plan.Budget.Documents = tc.docs
			plan.Budget.Tests = tc.tests
			plan.Budget.TreeDepth = tc.depth
			p := &limitsProfilePort{}
			for i, path := range tc.paths {
				plan.Targets = append(plan.Targets, domain.ExtractionTarget{Path: path, Kind: tc.kinds[i], Capabilities: []string{"dependency graph"}})
				p.tree = append(p.tree, domain.DirectoryNode{Path: path, Type: "file"})
			}
			before, _ := json.Marshal(plan)
			got, err := CollectEvidence(context.Background(), plan, p)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(p.requested, tc.want) || (len(tc.want) == 0 && p.extracts != 0) {
				t.Fatalf("paths outside quota: requested=%v want=%v extracts=%d", p.requested, tc.want, p.extracts)
			}
			if tc.omitted && (got.Complete || !bundleHasStatus(got, domain.TargetSkippedBudget)) {
				t.Fatalf("omission not explicit: %+v", got)
			}
			if !tc.omitted && !got.Complete {
				t.Fatalf("positive control incomplete: %+v", got)
			}
			for _, item := range got.Items {
				if item.Kind == domain.EvidenceRepository {
					continue
				}
				found := false
				for _, path := range tc.want {
					found = found || item.Path == path
				}
				if !found {
					t.Fatal("omitted content admitted", item.Path)
				}
			}
			after, _ := json.Marshal(plan)
			if string(before) != string(after) {
				t.Fatal("plan mutated")
			}
		})
	}
}
func TestGeneralResponseCannotReadmitBudgetOmittedPath(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = []domain.ExtractionTarget{{Path: "a.go", Kind: domain.EvidenceFile}, {Path: "b.go", Kind: domain.EvidenceFile}}
	plan.Budget.CodeFiles = 1
	p := &limitsProfilePort{tree: []domain.DirectoryNode{{Path: "a.go", Type: "file"}, {Path: "b.go", Type: "file"}}, extra: []domain.ExtractedFile{{Path: "b.go", Content: "unrequested bytes"}}}
	got, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.requested, []string{"a.go"}) {
		t.Fatal("provider asked for omitted path", p.requested)
	}
	for _, item := range got.Items {
		if item.Path == "b.go" {
			t.Fatal("response readmitted omitted path")
		}
	}
	if got.Complete || !bundleHasStatus(got, domain.TargetSkippedBudget) {
		t.Fatal("omission lost")
	}
}
func TestDerivedRecoveryRespectsCodeQuota(t *testing.T) {
	for _, quota := range []int{0, 1, 2} {
		t.Run(fmt.Sprint(quota), func(t *testing.T) {
			ref := evidenceIntegrityPlan().Repository
			policy := DefaultRecoveryPolicy()
			policy.AdditionalCodeFiles = quota
			diagnosis := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "implementation", Candidate: ref.Key(), Capability: "dependency graph", Cause: domain.CauseMissingImplementation, RuleID: "dependency-graph", Recoverable: true}}}
			plan, err := PlanEvidenceRecovery(diagnosis, policy, 1, nil)
			if err != nil {
				t.Fatal(err)
			}
			p := &limitsProfilePort{tree: []domain.DirectoryNode{{Path: "graph1.go", Type: "file"}, {Path: "graph2.go", Type: "file"}}}
			got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, nil, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
			if err != nil {
				t.Fatal(err)
			}
			if len(p.requested) != quota || (quota == 0 && p.extracts != 0) {
				t.Fatalf("derived paths ignored quota: quota=%d paths=%v extracts=%d", quota, p.requested, p.extracts)
			}
			if len(got.Bundles) != 1 {
				t.Fatal("missing bundle")
			}
			if quota < 2 && (got.Bundles[0].Complete || !bundleHasStatus(got.Bundles[0], domain.TargetSkippedBudget)) {
				t.Fatal("derived omission not explicit")
			}
			wantExecuted := 1
			if quota == 0 {
				wantExecuted = 0
			}
			if got.Metrics.TargetsExecuted != wantExecuted {
				t.Fatal("omitted target counted as executed", got.Metrics)
			}
		})
	}
}
func TestRecheckRejectsDeclaredSizeDivergence(t *testing.T) {
	for _, variant := range []string{"unknown", "exact", "tree larger", "file larger", "file smaller", "file negative", "tree negative", "tree conflict"} {
		t.Run(variant, func(t *testing.T) {
			req, p := recheckFixture()
			_, _, _, asOf := verificationFixture()
			size := int64(len(p.extracted.Files[0].Content))
			switch variant {
			case "exact":
				p.tree[0].Size = size
				p.extracted.Files[0].Size = size
			case "tree larger":
				p.tree[0].Size = size + 100
			case "file larger":
				p.extracted.Files[0].Size = size + 100
			case "file smaller":
				p.extracted.Files[0].Size = size - 1
			case "file negative":
				p.extracted.Files[0].Size = -1
			case "tree negative":
				p.tree[0].Size = -1
			case "tree conflict":
				p.tree[0].Size = size
				p.tree = append(p.tree, domain.DirectoryNode{Path: "LICENSE", Type: "file", Size: size + 1})
			}
			got, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
			if err != nil {
				t.Fatal(err)
			}
			valid := variant == "unknown" || variant == "exact"
			if len(got.Bundles) != 1 || len(got.Review.Verifications) != 1 {
				t.Fatal("missing result")
			}
			if valid {
				if !got.Bundles[0].Complete || got.Review.Verifications[0].Status != domain.TechnicalVerified {
					t.Fatal("valid size control rejected", got)
				}
			} else {
				if got.Bundles[0].Complete || len(got.Bundles[0].Items) != 0 || len(got.Failures) == 0 || got.Review.Verifications[0].Status == domain.TechnicalVerified {
					t.Fatalf("incomplete size promoted: variant=%s got=%+v", variant, got)
				}
			}
		})
	}
}

func TestDirectedSelectionCannotBypassTestQuota(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = nil
	plan.Budget.Tests = 0
	plan.RequirementQueries = []domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "graph", RuleID: "fixture", Terms: []string{"graph"}}}
	p := &limitsProfilePort{tree: []domain.DirectoryNode{{Path: "graph_test.go", Type: "file"}}}
	got, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if p.extracts != 0 || got.Complete || !bundleHasStatus(got, domain.TargetSkippedBudget) {
		t.Fatalf("directed selection ignored quota: paths=%v bundle=%+v", p.requested, got)
	}
}
