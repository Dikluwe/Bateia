/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/catalog_surfaces_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"testing"
)

type surfaceStub struct {
	caps    domain.CatalogSurfaceCapabilities
	matches []domain.SurfaceMatch
	calls   int
}

func (s *surfaceStub) SurfaceCapabilities() domain.CatalogSurfaceCapabilities { return s.caps }
func (s *surfaceStub) SearchSurface(context.Context, domain.SurfaceSearchRequest) (domain.SurfaceDiscovery, error) {
	s.calls++
	return domain.SurfaceDiscovery{Capabilities: s.caps, Matches: s.matches, Attempt: domain.CatalogAttempt{Status: "SUCCEEDED"}}, nil
}
func TestSurfacePlanPreservesFullDiscriminator(t *testing.T) {
	p := PlanCatalogSurfaces(domain.ExpansionPlan{Expansions: []domain.SearchExpansion{{ID: "x", Query: "open source project discovery", Discriminators: []string{"project discovery"}}}}, []domain.CatalogSurfaceCapabilities{{ID: domain.SurfaceReadme}})
	if !p[0].Valid || p[0].Effective != "project discovery" {
		t.Fatalf("bad plan: %+v", p)
	}
}
func TestTwoSurfaceMatchesAdvanceButRepeatedSurfaceDoesNot(t *testing.T) {
	candidate := func(s domain.CatalogSurfaceID) domain.SurfaceMatch {
		return domain.SurfaceMatch{Surface: s, Repository: "o/r", CommitID: "abc", Discriminators: []string{"software reuse"}}
	}
	readme := &surfaceStub{caps: domain.CatalogSurfaceCapabilities{ID: domain.SurfaceReadme}, matches: []domain.SurfaceMatch{candidate(domain.SurfaceReadme)}}
	topics := &surfaceStub{caps: domain.CatalogSurfaceCapabilities{ID: domain.SurfaceTopics}, matches: []domain.SurfaceMatch{candidate(domain.SurfaceTopics)}}
	plans := []domain.SurfaceQueryPlan{{ExpansionID: "a", Surface: domain.SurfaceReadme, Effective: "software reuse", PreservedDiscriminators: []string{"software reuse"}, Valid: true}, {ExpansionID: "b", Surface: domain.SurfaceTopics, Effective: "software reuse", PreservedDiscriminators: []string{"software reuse"}, Valid: true}}
	snap, err := DiscoverCatalogSurfaces(t.Context(), plans, []domain.CatalogSurface{readme, topics}, 3)
	if err != nil {
		t.Fatal(err)
	}
	if ScreenMultiView(snap).Decisions[0].Verdict != domain.CatalogAdvance {
		t.Fatal("independent views did not advance")
	}
}
func TestOpaqueMatchRemainsUnknown(t *testing.T) {
	s := &surfaceStub{caps: domain.CatalogSurfaceCapabilities{ID: domain.SurfacePackage}, matches: []domain.SurfaceMatch{{Surface: domain.SurfacePackage, Package: "x", UnknownReason: "source missing"}}}
	p := []domain.SurfaceQueryPlan{{ExpansionID: "x", Surface: domain.SurfacePackage, Effective: "reuse", PreservedDiscriminators: []string{"reuse"}, Valid: true}}
	got, _ := DiscoverCatalogSurfaces(t.Context(), p, []domain.CatalogSurface{s}, 1)
	if len(got.Candidates) != 0 || len(got.Unknown) == 0 {
		t.Fatal("unknown promoted")
	}
}

func TestResolvedRepositoryWithoutSnippetRemainsUnknown(t *testing.T) {
	s := &surfaceStub{caps: domain.CatalogSurfaceCapabilities{ID: domain.SurfaceReadme}, matches: []domain.SurfaceMatch{{Surface: domain.SurfaceReadme, Repository: "o/r", CommitID: "abc", UnknownReason: "snippet unavailable"}}}
	p := []domain.SurfaceQueryPlan{{ExpansionID: "x", Surface: domain.SurfaceReadme, Effective: "software reuse", PreservedDiscriminators: []string{"software reuse"}, Valid: true}}
	got, _ := DiscoverCatalogSurfaces(t.Context(), p, []domain.CatalogSurface{s}, 1)
	if len(got.Candidates) != 0 || len(got.Unknown) == 0 {
		t.Fatal("opaque snippet promoted")
	}
}
