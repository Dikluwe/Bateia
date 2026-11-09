/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_completeness_batch_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type completenessBatchPort struct {
	calls   []string
	tree    []domain.DirectoryNode
	content map[string]string
	sizes   map[string]int64
}

func (p *completenessBatchPort) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	p.calls = append(p.calls, "manifest")
	return nil, errors.New("unused")
}
func (p *completenessBatchPort) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.calls = append(p.calls, "metadata")
	if mutableIdeaRunCommit(c.CommitID) {
		c.CommitID = strings.Repeat("a", 40)
	}
	return c, nil
}
func (p *completenessBatchPort) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls = append(p.calls, "tree")
	return p.tree, nil
}
func (p *completenessBatchPort) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.calls = append(p.calls, "extract")
	files := []domain.ExtractedFile{}
	for _, path := range paths {
		if text, ok := p.content[path]; ok {
			files = append(files, domain.ExtractedFile{Path: path, Content: text, Size: p.sizes[path]})
		}
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}, Files: files}, nil
}
func completenessTarget(t *testing.T, b domain.EvidenceBundle, path string) domain.CollectionTargetResult {
	t.Helper()
	for _, result := range b.Targets {
		if result.Target.Path == path {
			return result
		}
	}
	t.Fatal("missing target", path)
	return domain.CollectionTargetResult{}
}
func TestGenericCollectionRejectsDeclaredSizeDivergence(t *testing.T) {
	for _, variant := range []string{"unknown", "exact", "file larger", "file smaller", "file negative", "tree larger", "tree negative", "tree conflict"} {
		t.Run(variant, func(t *testing.T) {
			plan := evidenceIntegrityPlan()
			plan.Targets = []domain.ExtractionTarget{{Path: "LICENSE", Kind: domain.EvidenceLicense}}
			size := int64(len(step0025MITLicense))
			p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}}, content: map[string]string{"LICENSE": step0025MITLicense}, sizes: map[string]int64{}}
			switch variant {
			case "exact":
				p.tree[0].Size = size
				p.sizes["LICENSE"] = size
			case "file larger":
				p.sizes["LICENSE"] = size + 100
			case "file smaller":
				p.sizes["LICENSE"] = size - 1
			case "file negative":
				p.sizes["LICENSE"] = -1
			case "tree larger":
				p.tree[0].Size = size + 100
			case "tree negative":
				p.tree[0].Size = -1
			case "tree conflict":
				p.tree[0].Size = size
				p.tree = append(p.tree, domain.DirectoryNode{Path: "LICENSE", Type: "file", Size: size + 1})
			}
			before, _ := json.Marshal(plan)
			bundle, err := CollectEvidence(context.Background(), plan, p)
			if err != nil {
				t.Fatal(err)
			}
			assessed, err := AssessEvidence(domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, DefaultEvidenceRules())
			if err != nil {
				t.Fatal(err)
			}
			valid := variant == "unknown" || variant == "exact"
			result := completenessTarget(t, bundle, "LICENSE")
			if valid {
				if !bundle.Complete || assessed.Assessments[0].License.Status != domain.EvidenceProven || result.Status != domain.TargetCollected {
					t.Fatal("valid size control lost")
				}
			} else {
				if bundle.Complete || len(bundle.Failures) == 0 || assessed.Assessments[0].License.Status == domain.EvidenceProven || result.Status != domain.TargetFailed {
					t.Fatalf("inconsistent size promoted: %s %+v", variant, bundle)
				}
				for _, item := range bundle.Items {
					if item.Path == "LICENSE" {
						t.Fatal("inconsistent source admitted")
					}
				}
			}
			after, _ := json.Marshal(plan)
			if string(before) != string(after) {
				t.Fatal("caller plan mutated")
			}
		})
	}
}
func TestGenericSizeFailureDoesNotCompleteDirectedCoverage(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = nil
	plan.RequirementQueries = []domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "graph", RuleID: "fixture", Terms: []string{"graph"}}}
	text := "package fixture\nfunc Graph(){}"
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "graph.go", Type: "file"}}, content: map[string]string{"graph.go": text}, sizes: map[string]int64{"graph.go": int64(len(text) + 1)}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if b.Complete || len(b.RequirementCoverage) != 1 || b.RequirementCoverage[0].Status != domain.RequirementCollectionIncomplete || len(b.RequirementCoverage[0].ExaminedPaths) != 0 {
		t.Fatalf("rejected file counted as examined: %+v", b)
	}
}
func TestByteOmissionIsBudgetAndAllowsLicenseRecovery(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			plan := evidenceIntegrityPlan()
			plan.Targets = []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation, Capabilities: []string{"dependency graph"}}, {Path: "LICENSE", Kind: domain.EvidenceLicense, Capabilities: []string{"dependency graph"}}}
			metadata, _ := json.Marshal(plan.Candidate)
			plan.Budget.TotalBytes = int64(len(metadata) + 200)
			p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "LICENSE", Type: "file"}}, content: map[string]string{"README.md": strings.Repeat("x", 200), "LICENSE": step0025MITLicense}}
			if missing {
				delete(p.content, "LICENSE")
			}
			b, err := CollectEvidence(context.Background(), plan, p)
			if err != nil {
				t.Fatal(err)
			}
			result := completenessTarget(t, b, "LICENSE")
			want := domain.TargetSkippedBudget
			if missing {
				want = domain.TargetNotFound
			}
			if result.Status != want || result.Bytes != 0 || b.Complete {
				t.Fatalf("wrong byte omission: %+v", result)
			}
			if !missing {
				if strings.Contains(result.Message, "não devolveu") {
					t.Fatal("budget blamed on provider")
				}
				assessed, err := AssessEvidence(domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, []domain.EvidenceBundle{b}, DefaultEvidenceRules())
				if err != nil {
					t.Fatal(err)
				}
				diagnosis := DiagnoseInconclusive("test", "digest", domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, assessed, []domain.EvidenceBundle{b}, DefaultEvidenceRules())
				budget := false
				for _, f := range diagnosis.Findings {
					budget = budget || f.Cause == domain.CauseBudgetExhausted
				}
				if !budget {
					t.Fatal("byte budget cause lost", diagnosis)
				}
				onlyLicense := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "license", Candidate: plan.Repository.Key(), Cause: domain.CauseUnknownLicense, RuleID: "GATE-LICENSE-001", Recoverable: true}}}
				next, err := PlanEvidenceRecovery(onlyLicense, DefaultRecoveryPolicy(), 1, AttemptedOriginalRecoveryTargets([]domain.EvidenceBundle{b}))
				if err != nil || len(next.Targets) != 1 || next.Targets[0].Path != "LICENSE" {
					t.Fatalf("budget blocked resume: err=%v plan=%+v", err, next)
				}
			}
		})
	}
}
func TestTruncatedFileIsBudgetLimitedAfterRawSizeCheck(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = []domain.ExtractionTarget{{Path: "LICENSE", Kind: domain.EvidenceLicense}}
	metadata, _ := json.Marshal(plan.Candidate)
	plan.Budget.TotalBytes = int64(len(metadata) + 20)
	size := int64(len(step0025MITLicense))
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "LICENSE", Type: "file", Size: size}}, content: map[string]string{"LICENSE": step0025MITLicense}, sizes: map[string]int64{"LICENSE": size}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	result := completenessTarget(t, b, "LICENSE")
	if result.Status != domain.TargetSkippedBudget || result.Bytes != 20 || b.Complete || len(b.Failures) != 0 {
		t.Fatalf("truncation misclassified as size/provider failure: %+v", b)
	}
	found := false
	for _, item := range b.Items {
		if item.Path == "LICENSE" {
			found = true
			if item.Content != step0025MITLicense[:20] || item.ContentHash != contentHash(item.Content) {
				t.Fatal("fragment bytes/hash altered")
			}
		}
	}
	if !found {
		t.Fatal("fragment not retained")
	}
	if len(AttemptedOriginalRecoveryTargets([]domain.EvidenceBundle{b})) != 0 {
		t.Fatal("partial target blocked in active history")
	}
}
func TestRecoveryCountsPartialReadWithoutBlockingRetry(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	plan := budgetPlan(ref)
	plan.Policy.AdditionalBytes = 400
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}}, content: map[string]string{"LICENSE": step0025MITLicense}}
	got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, nil, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
	if err != nil {
		t.Fatal(err)
	}
	if got.Metrics.TargetsExecuted != 1 || got.Metrics.ProviderCalls != 3 || len(ExecutedRecoveryTargets(plan, got.Bundles)) != 0 {
		t.Fatalf("partial read/history mismatch: %+v", got.Metrics)
	}
}
func TestCoordinatedBatchRejectsKnownDuplicateBeforeEffects(t *testing.T) {
	for _, variant := range []string{"exact", "case", "unresolved"} {
		for _, reverse := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%v", variant, reverse), func(t *testing.T) {
				a := evidenceIntegrityPlan()
				b := a
				switch variant {
				case "case":
					b.Repository.Owner = strings.ToUpper(b.Repository.Owner)
					b.Candidate.Owner = b.Repository.Owner
					b.Repository.Repo = strings.ToUpper(b.Repository.Repo)
					b.Candidate.Repo = b.Repository.Repo
				case "unresolved":
					a.Repository.Commit = ""
					a.Candidate.CommitID = ""
					b.Repository.Commit = "HEAD"
					b.Candidate.CommitID = "HEAD"
				}
				plans := []domain.EvidenceExtractionPlan{a, b}
				if reverse {
					plans[0], plans[1] = plans[1], plans[0]
				}
				before, _ := json.Marshal(plans)
				p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, content: map[string]string{"README.md": "queue mechanism"}}
				got, err := CollectInsightEvidence(context.Background(), plans, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "queue"})
				if !errors.Is(err, domain.ErrInvalidRunInput) || len(p.calls) != 0 || !reflect.DeepEqual(got, domain.InsightCollectionReport{}) {
					t.Fatalf("known invalid batch reached effects: err=%v calls=%v", err, p.calls)
				}
				after, _ := json.Marshal(plans)
				if string(before) != string(after) {
					t.Fatal("batch mutated")
				}
			})
		}
	}
	for _, count := range []int{1, 2} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			a := evidenceIntegrityPlan()
			b := a
			b.Repository.Commit = strings.Repeat("b", 40)
			b.Candidate.CommitID = b.Repository.Commit
			plans := []domain.EvidenceExtractionPlan{a}
			if count == 2 {
				plans = append(plans, b)
			}
			p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, content: map[string]string{"README.md": "queue mechanism"}}
			got, err := CollectInsightEvidence(context.Background(), plans, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "queue"})
			if err != nil || len(got.Bundles) != count || len(p.calls) != 3*count {
				t.Fatalf("distinct origins rejected: %v calls=%v", err, p.calls)
			}
		})
	}
}
