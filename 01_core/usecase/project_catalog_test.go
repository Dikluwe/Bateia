/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/project_catalog_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"testing"
)

type catalogStub struct {
	caps   domain.ProviderCapabilities
	calls  int
	result domain.CatalogDiscovery
}

func (s *catalogStub) Capabilities() domain.ProviderCapabilities { return s.caps }
func (s *catalogStub) SearchCatalog(context.Context, domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
	s.calls++
	return s.result, nil
}

func TestCatalogRejectsWrongRoleBeforeNetwork(t *testing.T) {
	s := &catalogStub{caps: domain.ProviderCapabilities{Provider: "sourcegraph", Roles: []domain.DiscoveryProviderRole{domain.RoleCodeProof}}}
	_, err := DiscoverProjectCatalog(t.Context(), domain.ExpansionPlan{Expansions: []domain.SearchExpansion{{ID: "x", Query: "software reuse", Discriminators: []string{"software reuse"}}}}, 3, domain.DiscoveryExecutionPolicy{}, s)
	if err == nil || s.calls != 0 {
		t.Fatalf("wrong role touched network: %v calls=%d", err, s.calls)
	}
}

func TestCatalogDeduplicatesAndPreservesOrigins(t *testing.T) {
	c := domain.CatalogCandidate{Repository: domain.RepositoryCandidate{Owner: "Org", Repo: "Repo", CommitID: "abc"}, MatchedFields: []string{"description"}}
	s := &catalogStub{caps: domain.ProviderCapabilities{Provider: "github", Roles: []domain.DiscoveryProviderRole{domain.RoleProjectCatalog}}, result: domain.CatalogDiscovery{Provider: "github", Role: domain.RoleProjectCatalog, SearchUnit: "repository", Candidates: []domain.CatalogCandidate{c}, Attempt: domain.CatalogAttempt{Query: "software reuse", Status: "SUCCEEDED"}}}
	plan := domain.ExpansionPlan{Expansions: []domain.SearchExpansion{{ID: "a", Query: "software reuse", Discriminators: []string{"software reuse"}}, {ID: "b", Query: "project discovery", Discriminators: []string{"project discovery"}}}}
	snapshot, err := DiscoverProjectCatalog(t.Context(), plan, 3, domain.DiscoveryExecutionPolicy{}, s)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Candidates) != 1 || len(snapshot.Candidates[0].Origins) != 2 {
		t.Fatalf("provenance lost: %+v", snapshot.Candidates)
	}
}

func TestCatalogScreeningKeepsUnknownAndStarsDoNotAdvance(t *testing.T) {
	r := ScreenProjectCatalog(domain.CatalogSnapshot{Candidates: []domain.CatalogCandidate{{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "popular", Stars: 999999}, MatchedFields: []string{"readme"}}}})
	if r.Decisions[0].Verdict == domain.CatalogAdvance {
		t.Fatal("stars promoted opaque candidate")
	}
}

func TestMultiQueryWithoutDiscriminatorDoesNotAdvance(t *testing.T) {
	c := domain.CatalogCandidate{Repository: domain.RepositoryCandidate{Owner: "o", Repo: "r", CommitID: "abc"}, Origins: []domain.CandidateOrigin{{ExpansionID: "a"}, {ExpansionID: "b"}}}
	r := ScreenProjectCatalog(domain.CatalogSnapshot{Candidates: []domain.CatalogCandidate{c}})
	if r.Decisions[0].Verdict == domain.CatalogAdvance {
		t.Fatal("multi-query noise advanced without catalog-field discriminator")
	}
}

func TestProofPlanScopesEveryQueryToAdvancedRepository(t *testing.T) {
	source := domain.ExpansionPlan{Expansions: []domain.SearchExpansion{{ID: "x", Query: "software reuse", Discriminators: []string{"software reuse"}}}}
	report := domain.CatalogScreeningReport{Decisions: []domain.CatalogDecision{{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "repo"}, Verdict: domain.CatalogAdvance, Origins: []domain.CandidateOrigin{{ExpansionID: "x"}}}, {Repository: domain.RepositoryCandidate{Owner: "bad", Repo: "noise"}, Verdict: domain.CatalogReview, Origins: []domain.CandidateOrigin{{ExpansionID: "x"}}}}}
	got := BuildCatalogProofPlan(report, source)
	if len(got.Expansions) != 1 || got.Expansions[0].RepositoryScope != "org/repo" {
		t.Fatalf("unscoped proof: %+v", got)
	}
}
