/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/investigation_plan_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/json"
	"sync"
	"testing"
)

func TestStep0043InvestigationExactRequirementBinding(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "render images", Capabilities: []string{"OpenGL OIT", "image export"}}
	snapshot := domain.DiscoverySnapshot{Idea: idea, PerExpansion: []domain.ExpansionDiscovery{{Expansion: domain.SearchExpansion{ID: "q-oit", Query: "OpenGL OIT", Covers: []string{"OpenGL OIT"}}, Provider: domain.ProviderDiscovery{Provider: "public", Attempts: []domain.ProviderAttempt{{Query: "OpenGL rendering", ResultCount: 0}}}}}}
	plan := usecase.BuildInvestigationPlan(idea, snapshot, nil)
	ids := map[string]string{}
	for _, r := range plan.Requirements {
		ids[r.Text] = r.ID
	}
	if len(plan.Searches) != 1 {
		t.Fatalf("search count: %d", len(plan.Searches))
	}
	idsUsed := plan.Searches[0].RequirementIDs
	found := false
	for _, id := range idsUsed {
		if id == ids["OpenGL OIT"] {
			found = true
		}
		if id == ids["image export"] {
			t.Fatal("query incorrectly attributed to unrelated mandatory capability")
		}
	}
	if !found {
		t.Fatalf("exact originating requirement not linked: %+v", plan.Searches[0])
	}
}

type step0043DiscoverySpy struct {
	mu       sync.Mutex
	requests []domain.DiscoveryRequest
}

func (s *step0043DiscoverySpy) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	s.mu.Lock()
	s.requests = append(s.requests, request)
	s.mu.Unlock()
	return domain.ProviderDiscovery{Provider: "public-spy", Attempts: []domain.ProviderAttempt{{Query: request.Query, Error: "too narrow"}, {Query: "rendering", ResultCount: 0}}}, nil
}
func TestStep0043DiscoveryEmitsRequirementIDsWithoutMutatingInput(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "render images", Capabilities: []string{"OpenGL OIT"}}
	input := domain.ExpansionPlan{Idea: idea, Expansions: []domain.SearchExpansion{{ID: "q-oit", Facet: domain.FacetCapability, Query: "OpenGL OIT", Covers: []string{"capabilities[0]"}}}}
	before, _ := json.Marshal(input)
	spy := &step0043DiscoverySpy{}
	snapshot, err := usecase.DiscoverExpansionPlan(context.Background(), input, "literal", "", 5, 1, spy)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(input)
	if string(before) != string(after) {
		t.Fatal("discovery mutated caller plan")
	}
	projected := usecase.BuildInvestigationPlan(idea, snapshot, nil)
	wanted := ""
	for _, r := range projected.Requirements {
		if r.Text == "OpenGL OIT" {
			wanted = r.ID
		}
	}
	if wanted == "" || len(spy.requests) != 1 {
		t.Fatalf("missing originating requirement or provider invocation: %q %+v", wanted, spy.requests)
	}
	matched := false
	for _, id := range spy.requests[0].RequirementIDs {
		if id == wanted {
			matched = true
		}
	}
	if !matched {
		t.Fatalf("provider received no correct obligation ID: %+v", spy.requests[0])
	}
	fallback := false
	for _, search := range projected.Searches {
		if search.EffectiveQuery == "rendering" {
			fallback = true
			if search.OriginalQuery != "OpenGL OIT" {
				t.Fatal("fallback changed original intent")
			}
		}
	}
	if !fallback {
		t.Fatal("fallback attempt not preserved")
	}
}
