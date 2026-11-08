/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_budget_revision_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type budgetRevisionPort struct {
	calls     []string
	refs      []domain.CommitRef
	hook      func(string)
	content   string
	emptyTree bool
}

func (p *budgetRevisionPort) call(stage string, c domain.RepositoryCandidate) {
	p.calls = append(p.calls, stage)
	p.refs = append(p.refs, domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID})
	if p.hook != nil {
		p.hook(stage)
	}
}
func (p *budgetRevisionPort) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("unused")
}
func (p *budgetRevisionPort) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.call("metadata", c)
	return c, nil
}
func (p *budgetRevisionPort) FetchTree(_ context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.call("tree", c)
	if p.emptyTree {
		return nil, nil
	}
	return []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}}, nil
}
func (p *budgetRevisionPort) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.call("extract", c)
	files := []domain.ExtractedFile{}
	for _, path := range paths {
		files = append(files, domain.ExtractedFile{Path: path, Content: p.content})
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}, Files: files}, nil
}
func (p *budgetRevisionPort) SearchImplementations(_ context.Context, c domain.RepositoryCandidate, _ string, _ domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	p.call("search", c)
	return domain.ImplementationSearchReport{Status: "completed"}, nil
}
func budgetRecovery(t *testing.T, plan domain.EvidenceRecoveryPlan, original []domain.EvidenceBundle, p *budgetRevisionPort) domain.RecoveryRoundResult {
	t.Helper()
	got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, original, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
	if err != nil {
		t.Fatal(err)
	}
	return got
}
func budgetPlan(refs ...domain.CommitRef) domain.EvidenceRecoveryPlan {
	p := domain.EvidenceRecoveryPlan{SchemaVersion: "1", Round: 1, Policy: DefaultRecoveryPolicy()}
	for i, ref := range refs {
		p.Targets = append(p.Targets, domain.RecoveryTarget{FindingID: fmt.Sprint(i), Candidate: ref, Path: "LICENSE", Kind: domain.EvidenceLicense, RuleID: "GATE-LICENSE-001"})
	}
	return p
}
func TestRecoverySeparatesRevisionsAndPreservesOriginals(t *testing.T) {
	a := evidenceIntegrityPlan().Repository
	b := a
	b.Commit = strings.Repeat("b", 40)
	diagnosis := domain.InconclusiveDiagnosis{SchemaVersion: "1"}
	for i, ref := range []domain.CommitRef{a, b} {
		diagnosis.Findings = append(diagnosis.Findings, domain.InconclusiveFinding{ID: fmt.Sprint(i), Candidate: ref.Key(), Cause: domain.CauseUnknownLicense, RuleID: "GATE-LICENSE-001", Recoverable: true})
	}
	plan, err := PlanEvidenceRecovery(diagnosis, DefaultRecoveryPolicy(), 1, nil)
	if err != nil || len(plan.Targets) != 2 {
		t.Fatalf("plan: %v %#v", err, plan)
	}
	for _, withOriginal := range []bool{false, true} {
		t.Run(fmt.Sprint(withOriginal), func(t *testing.T) {
			originals := []domain.EvidenceBundle{}
			if withOriginal {
				for _, ref := range []domain.CommitRef{a, b} {
					ep := evidenceIntegrityPlan()
					ep.Repository = ref
					ep.Candidate.CommitID = ref.Commit
					ep.Targets = []domain.ExtractionTarget{{Path: "LICENSE", Kind: domain.EvidenceLicense}}
					bundle, e := CollectEvidence(context.Background(), ep, &budgetRevisionPort{content: "identical content across commits"})
					if e != nil {
						t.Fatal(e)
					}
					originals = append(originals, bundle)
				}
			}
			before, _ := json.Marshal(originals)
			p := &budgetRevisionPort{content: "identical content across commits"}
			got := budgetRecovery(t, plan, originals, p)
			read := map[string]bool{}
			for i, stage := range p.calls {
				if stage == "extract" {
					read[p.refs[i].Key()] = true
				}
			}
			if len(read) != 2 || !read[a.Key()] || !read[b.Key()] || len(got.Bundles) != 2 || got.Metrics.TargetsExecuted != 2 {
				t.Fatalf("revision lost: reads=%v bundles=%d metrics=%+v", read, len(got.Bundles), got.Metrics)
			}
			for _, bundle := range got.Bundles {
				found := false
				for _, item := range bundle.Items {
					if item.Ref != bundle.Repository {
						t.Fatal("mixed revisions")
					}
					if item.Path == "LICENSE" {
						found = true
					}
				}
				if !found {
					t.Fatal("equal hash from another revision suppressed source")
				}
			}
			after, _ := json.Marshal(originals)
			if string(before) != string(after) {
				t.Fatal("original bundles mutated")
			}
		})
	}
}
func TestRecoveryDiagnosisDoesNotTransferRevisionFailures(t *testing.T) {
	before, assessment, bundles := recoveryDiagnosisFixture()
	a := assessment.Assessments[0].Repository
	b := a
	b.Commit = "def"
	other := assessment.Assessments[0]
	other.Repository = b
	assessment.Assessments = append(assessment.Assessments, other)
	failed := bundles[0]
	failed.Repository = b
	failed.Failures = []domain.ExecutionFailure{{Stage: "content", Message: "failure only at def"}}
	bundles = append(bundles, failed)
	got := DiagnoseInconclusive("audit", "digest", before, assessment, bundles, DefaultEvidenceRules())
	causes := map[string]domain.InconclusiveCause{}
	for _, f := range got.Findings {
		if f.Capability == "dependency graph" {
			causes[f.Candidate] = f.Cause
		}
	}
	if causes[a.Key()] == domain.CauseProviderFailure || causes[b.Key()] != domain.CauseProviderFailure {
		t.Fatalf("revision failures crossed: %v", causes)
	}
}
func TestRecoveryPreparationUsesTargetPin(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	plan := budgetPlan(ref)
	plan.Targets[0].Path = ""
	plan.Targets[0].Patterns = []string{"license"}
	base := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Plan: evidenceIntegrityPlan()}
	base.Plan.Candidate.CommitID = strings.Repeat("b", 40)
	p := &budgetRevisionPort{content: "license text"}
	budgetRecovery(t, plan, []domain.EvidenceBundle{base}, p)
	for _, seen := range p.refs {
		if seen != ref {
			t.Fatalf("preparation used wrong pin: %v", seen)
		}
	}
}
func TestGeneralCollectionStopsBetweenProviderPorts(t *testing.T) {
	for _, failure := range []error{context.Canceled, context.DeadlineExceeded} {
		for _, stage := range []string{"tree", "search"} {
			t.Run(failure.Error()+"/"+stage, func(t *testing.T) {
				ctx := &boundaryContext{Context: context.Background()}
				p := &budgetRevisionPort{content: "license"}
				p.hook = func(s string) {
					if s == stage {
						ctx.failure = failure
					}
				}
				_, err := CollectEvidence(ctx, evidenceIntegrityPlan(), p)
				want := []string{"metadata", "tree"}
				if stage == "search" {
					want = append(want, "search")
				}
				if !errors.Is(err, failure) || !reflect.DeepEqual(p.calls, want) {
					t.Fatalf("continued after interruption: err=%v calls=%v", err, p.calls)
				}
			})
		}
	}
	p := &budgetRevisionPort{content: "license"}
	if _, err := CollectEvidence(context.Background(), evidenceIntegrityPlan(), p); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.calls, []string{"metadata", "tree", "search", "extract"}) {
		t.Fatal("positive control lost ports", p.calls)
	}
}
func TestRecoverySharesRoundByteBudget(t *testing.T) {
	a := evidenceIntegrityPlan().Repository
	b := a
	b.Repo = "second"
	c := a
	c.Repo = "third"
	for _, limit := range []int64{400, 800, 4096} {
		for _, reverse := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/%v", limit, reverse), func(t *testing.T) {
				plan := budgetPlan(a, b, c)
				plan.Targets[2].Path = ""
				plan.Targets[2].Patterns = []string{"license"}
				plan.Policy.AdditionalBytes = limit
				if reverse {
					plan.Targets[0], plan.Targets[2] = plan.Targets[2], plan.Targets[0]
				}
				p := &budgetRevisionPort{content: strings.Repeat("x", 512)}
				got := budgetRecovery(t, plan, nil, p)
				var total int64
				for _, bundle := range got.Bundles {
					total += bundle.BytesRead
				}
				if total > limit || got.Metrics.AdditionalBytes != total {
					t.Fatalf("round exceeded: total=%d limit=%d metrics=%d", total, limit, got.Metrics.AdditionalBytes)
				}
				if limit == 400 {
					if got.Metrics.TargetsExecuted != 1 || len(p.calls) != 3 {
						t.Fatalf("exhausted groups executed: %+v calls=%v", got.Metrics, p.calls)
					}
					for _, ref := range p.refs {
						if ref != a {
							t.Fatalf("noncanonical allocation: %v", p.refs)
						}
					}
					for _, bundle := range got.Bundles {
						if bundle.Repository != a {
							if bundle.Complete || bundle.BytesRead != 0 || !bundleHasStatus(bundle, domain.TargetSkippedBudget) {
								t.Fatalf("omission not explicit: %+v", bundle)
							}
						}
					}
				}
				if limit == 4096 && got.Metrics.TargetsExecuted != 3 {
					t.Fatal("positive control skipped targets")
				}
				for _, assessment := range got.Assessment.Assessments {
					complete := false
					for _, bundle := range got.Bundles {
						if bundle.Repository == assessment.Repository {
							complete = bundle.Complete
						}
					}
					for _, cap := range assessment.Capabilities {
						if !complete && cap.Status == domain.EvidenceAbsent {
							t.Fatal("budget became absence")
						}
					}
				}
			})
		}
	}
}
func TestRecoveryDoesNotChargeOriginalBytesToRound(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	ep := evidenceIntegrityPlan()
	ep.Targets = []domain.ExtractionTarget{{Path: "LICENSE", Kind: domain.EvidenceLicense}}
	original, err := CollectEvidence(context.Background(), ep, &budgetRevisionPort{content: strings.Repeat("old bytes", 1024)})
	if err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(original)
	plan := budgetPlan(ref)
	plan.Policy.AdditionalBytes = 400
	p := &budgetRevisionPort{content: strings.Repeat("new bytes", 1024)}
	got := budgetRecovery(t, plan, []domain.EvidenceBundle{original}, p)
	if original.BytesRead <= 400 || got.Metrics.AdditionalBytes != 400 || len(got.Bundles) != 1 || got.Bundles[0].BytesRead <= 400 || len(p.calls) != 3 {
		t.Fatalf("historical bytes charged to round: metrics=%+v calls=%v", got.Metrics, p.calls)
	}
	after, _ := json.Marshal(original)
	if string(before) != string(after) {
		t.Fatal("historical input changed")
	}
}

func TestMetadataFitsTotalByteBudgetWithoutTruncation(t *testing.T) {
	ep := evidenceIntegrityPlan()
	ep.Targets = nil
	encoded, err := json.Marshal(ep.Candidate)
	if err != nil {
		t.Fatal(err)
	}
	size := int64(len(encoded))
	for _, limit := range []int64{1, size - 1, size, size + 1} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			ep.Budget.TotalBytes = limit
			p := &budgetRevisionPort{emptyTree: true}
			got, e := CollectEvidence(context.Background(), ep, p)
			if e != nil {
				t.Fatal(e)
			}
			if got.BytesRead > limit {
				t.Fatalf("metadata exceeded: bytes=%d limit=%d", got.BytesRead, limit)
			}
			if limit < size {
				if len(got.Items) != 0 || got.BytesRead != 0 || got.Complete || !bundleHasStatus(got, domain.TargetSkippedBudget) {
					t.Fatalf("oversized JSON admitted: %+v", got)
				}
			} else {
				if len(got.Items) != 1 || got.BytesRead != size || got.Items[0].Content != string(encoded) || !json.Valid([]byte(got.Items[0].Content)) || got.Items[0].ContentHash != fmt.Sprintf("%x", sha256.Sum256(encoded)) {
					t.Fatalf("metadata altered: %+v", got)
				}
			}
		})
	}
}
