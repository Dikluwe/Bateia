/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/investigation_boundaries_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type boundaryContext struct {
	context.Context
	failure error
}

func (c *boundaryContext) Err() error { return c.failure }

type boundaryRecheckPort struct {
	*technicalRecheckProvider
	afterTree, afterExtract func()
}

func (p *boundaryRecheckPort) FetchTree(ctx context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	v, e := p.technicalRecheckProvider.FetchTree(ctx, c)
	if p.afterTree != nil {
		p.afterTree()
	}
	return v, e
}
func (p *boundaryRecheckPort) ExtractRepoFiles(ctx context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	v, e := p.technicalRecheckProvider.ExtractRepoFiles(ctx, c, paths)
	if p.afterExtract != nil {
		p.afterExtract()
	}
	return v, e
}

func TestRecheckStopsAtProviderCancellationAndDeadline(t *testing.T) {
	for _, failure := range []error{context.Canceled, context.DeadlineExceeded} {
		for _, stage := range []string{"tree", "extract"} {
			t.Run(failure.Error()+"/"+stage, func(t *testing.T) {
				req, base := recheckFixture()
				_, _, _, asOf := verificationFixture()
				ctx := &boundaryContext{Context: context.Background()}
				p := &boundaryRecheckPort{technicalRecheckProvider: base}
				hook := func() { ctx.failure = failure }
				wantCalls := 1
				if stage == "tree" {
					p.afterTree = hook
				} else {
					p.afterExtract = hook
					wantCalls = 2
				}
				got, err := RecheckTechnicalSources(ctx, req, p, asOf)
				if !errors.Is(err, failure) || !reflect.DeepEqual(got, domain.TechnicalSourceRecheckReport{}) || p.calls != wantCalls {
					t.Fatalf("interruption returned evidence or reached another port: err=%v calls=%d", err, p.calls)
				}
			})
		}
	}
	req, p := recheckFixture()
	_, _, _, asOf := verificationFixture()
	for _, ctx := range []context.Context{nil, &boundaryContext{Context: context.Background(), failure: context.Canceled}} {
		func() {
			defer func() {
				if failure := recover(); failure != nil {
					t.Errorf("invalid context panicked: %v", failure)
				}
			}()
			got, err := RecheckTechnicalSources(ctx, req, p, asOf)
			if err == nil || p.calls != 0 || !reflect.DeepEqual(got, domain.TechnicalSourceRecheckReport{}) {
				t.Fatal("invalid context admitted")
			}
		}()
	}
}

type boundaryRecoveryPort struct {
	evidenceIntegrityPort
	requested []string
	extra     bool
}

func (p *boundaryRecoveryPort) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}, {Path: "selected.go", Type: "file"}, {Path: "outside.go", Type: "file"}}, nil
}
func (p *boundaryRecoveryPort) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.requested = append(p.requested, paths...)
	files := []domain.ExtractedFile{}
	for _, path := range paths {
		text := "package fixture\nfunc selected(){}"
		if path == "LICENSE" {
			text = step0025MITLicense
		}
		files = append(files, domain.ExtractedFile{Path: path, Content: text})
	}
	if p.extra {
		files = append(files, domain.ExtractedFile{Path: "unexpected.go", Content: "package unexpected"})
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}, Files: files}, nil
}
func boundaryRecover(t *testing.T, target domain.RecoveryTarget, p *boundaryRecoveryPort) (domain.RecoveryRoundResult, error) {
	t.Helper()
	plan := domain.EvidenceRecoveryPlan{SchemaVersion: "1", Round: 1, Policy: DefaultRecoveryPolicy(), Targets: []domain.RecoveryTarget{target}}
	return ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, nil, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
}
func TestRecoveryReadsOnlyExplicitOrDerivedFrozenTargets(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	for _, derived := range []bool{false, true} {
		target := domain.RecoveryTarget{FindingID: "finding", Candidate: ref, Path: "LICENSE", Kind: domain.EvidenceLicense, RuleID: "GATE-LICENSE-001"}
		want := []string{"LICENSE"}
		if derived {
			target.Path = ""
			target.Kind = domain.EvidenceFile
			target.Patterns = []string{"selected"}
			want = []string{"selected.go"}
		}
		p := &boundaryRecoveryPort{}
		got, err := boundaryRecover(t, target, p)
		if err != nil || !reflect.DeepEqual(p.requested, want) {
			t.Fatalf("frozen selection expanded: %v %v", err, p.requested)
		}
		if len(got.Bundles) != 1 || got.Bundles[0].Complete {
			t.Fatal("limited collection claimed full coverage")
		}
		for _, item := range got.Bundles[0].Items {
			if item.Path != "@repository-metadata" && item.Path != want[0] {
				t.Fatal("unplanned content admitted", item.Path)
			}
		}
	}
	p := &boundaryRecoveryPort{extra: true}
	got, err := boundaryRecover(t, domain.RecoveryTarget{FindingID: "finding", Candidate: ref, Path: "LICENSE", Kind: domain.EvidenceLicense}, p)
	if !errors.Is(err, domain.ErrIntegrity) || len(got.Bundles) != 0 {
		t.Fatalf("provider extras admitted: %v", err)
	}
	// The ordinary funnel remains free to sample the tree under its existing budgets.
	p = &boundaryRecoveryPort{}
	if _, err := CollectEvidence(context.Background(), evidenceIntegrityPlan(), p); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, path := range p.requested {
		found = found || path == "outside.go"
	}
	if !found {
		t.Fatal("general collection unexpectedly restricted")
	}
}

func TestNeighborhoodRejectsOrphanDuplicateAndInconsistentReceiptsBeforeEffects(t *testing.T) {
	for name, change := range map[string]func(*domain.InsightSiteReport){
		"orphan": func(r *domain.InsightSiteReport) {
			v := r.Readmes[0]
			v.CandidateID = "orphan"
			copy := *v.Receipt
			copy.ReferenceID = "orphan"
			v.Receipt = &copy
			r.Readmes = append(r.Readmes, v)
		},
		"duplicate": func(r *domain.InsightSiteReport) { r.Readmes = append(r.Readmes, r.Readmes[0]) },
		"nil receipt": func(r *domain.InsightSiteReport) {
			v := r.Readmes[0]
			v.CandidateID = "orphan"
			v.Receipt = nil
			r.Readmes = append(r.Readmes, v)
		},
		"outer status": func(r *domain.InsightSiteReport) {
			v := r.Readmes[0]
			v.CandidateID = "orphan"
			v.Status = "READ_FAILED"
			r.Readmes = append(r.Readmes, v)
		},
		"inner status": func(r *domain.InsightSiteReport) {
			v := r.Readmes[0]
			v.CandidateID = "orphan"
			copy := *v.Receipt
			copy.Status = "READ_FAILED"
			v.Receipt = &copy
			r.Readmes = append(r.Readmes, v)
		},
	} {
		t.Run(name, func(t *testing.T) {
			first, reader, bindings := neighborhoodFixture(t)
			change(&first)
			first.ID = ""
			first.ID = "site_discovery_" + insightHash(first)
			before, _ := json.Marshal(first)
			one := 1
			got, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: first, FollowLinks: &one}, bindings)
			if !errors.Is(err, domain.ErrIntegrity) || got.ID != "" || len(reader.reads) != 0 || len(reader.resolves) != 0 {
				t.Fatalf("inconsistent receipt reached effects: %v", err)
			}
			after, _ := json.Marshal(first)
			if string(before) != string(after) {
				t.Fatal("rejected input mutated")
			}
		})
	}
	first, _, _ := neighborhoodFixture(t)
	first.Readmes = append(first.Readmes, domain.InsightSiteReadme{CandidateID: "unread", ConnectorID: domain.ConnectorGitHubCatalog, Status: "READ_FAILED"})
	first.ID = ""
	first.ID = "site_discovery_" + insightHash(first)
	if _, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: first}); err != nil {
		t.Fatal("uncollected receipt rejected", err)
	}
}
