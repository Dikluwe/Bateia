/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/discovery_canonical_response_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"testing"
)

type canonicalDiscoveryPort struct {
	mode    string
	reverse bool
}

func (p canonicalDiscoveryPort) Discover(_ context.Context, r domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	a := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "pin", Description: "Alpha"}
	b := a
	b.Description = "Zulu"
	candidates := []domain.RepositoryCandidate{a}
	matches := []domain.ProviderCandidateMatch{{Repository: "org/repo", AttemptIndex: 0, EffectiveQuery: r.Query, Evidence: []domain.MatchEvidence{{MatchType: "file", Path: "a.go", Snippet: "first"}}}}
	if p.mode == "metadata" {
		candidates = append(candidates, b)
		if p.reverse {
			candidates[0], candidates[1] = candidates[1], candidates[0]
		}
	}
	if p.mode == "matches" {
		second := matches[0]
		second.Evidence = []domain.MatchEvidence{{MatchType: "file", Path: "z.go", Snippet: "second"}}
		matches = append(matches, second)
		if p.reverse {
			matches[0], matches[1] = matches[1], matches[0]
		}
	}
	return domain.ProviderDiscovery{Provider: "fixture", Candidates: candidates, Matches: matches, Attempts: []domain.ProviderAttempt{{Query: r.Query, AttemptIndex: 0, ResultCount: len(candidates), Status: "SUCCEEDED"}}}, nil
}

func TestDiscoveryResponseVariantsCanonicalAndComplete(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "gerenciar receitas"}, Expansions: regional.Proposals}
	for _, mode := range []string{"control", "metadata", "matches"} {
		t.Run(mode, func(t *testing.T) {
			var first string
			for _, reverse := range []bool{false, true} {
				snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 5, 1, canonicalDiscoveryPort{mode, reverse})
				if err != nil {
					t.Fatal(err)
				}
				raw, _ := json.Marshal(snapshot)
				if !reverse {
					first = string(raw)
				} else if first != string(raw) {
					t.Fatalf("response permutation changed snapshot: %s", mode)
				}
				if mode == "matches" && len(snapshot.Candidates[0].Origins) != 2 {
					t.Fatal("matching evidence discarded")
				}
				if mode == "metadata" && snapshot.Candidates[0].Repository.Description != "Alpha" {
					t.Fatal("noncanonical representative")
				}
				if mode == "control" && dossierHashBytes(raw) != "806645965c200ec0ee5d72c523606b99a9b10ea56abc08a4c6094a0d24549c3b" {
					t.Fatal("historical canonical snapshot changed")
				}
			}
		})
	}
}
func TestRegionalCatalogRejectsPaddedKeysWithoutInventingGap(t *testing.T) {
	for _, key := range []string{" recipe-manager", "recipe-manager ", " recipe-manager "} {
		catalog := DefaultRegionalQueryCatalog()
		catalog.Concepts[0].Concept = key
		plan, err := PlanRegionalDiscovery(catalog, "recipe-manager", []string{"pt-BR"})
		if err == nil || len(plan.Proposals) > 0 || len(plan.Gaps) > 0 {
			t.Fatalf("inconsistent catalog accepted: %q %+v %v", key, plan, err)
		}
	}
	catalog := DefaultRegionalQueryCatalog()
	catalog.Concepts[0].Queries = nil
	plan, err := PlanRegionalDiscovery(catalog, "recipe-manager", []string{"pt-BR"})
	if err != nil || len(plan.Gaps) != 1 || plan.Gaps[0].Status != "Unknown" {
		t.Fatal("legitimate gap changed", err)
	}
}
func TestDiscoveryMetadataTotalOrder(t *testing.T) {
	a := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "pin"}
	b := a
	b.Archived = true
	var first string
	for _, values := range [][]domain.RepositoryCandidate{{a, b}, {b, a}} {
		snapshot := domain.DiscoverySnapshot{PerExpansion: []domain.ExpansionDiscovery{{Provider: domain.ProviderDiscovery{Candidates: append([]domain.RepositoryCandidate(nil), values...)}}}}
		canonicalizeDirectDiscoverySnapshot(&snapshot)
		raw, _ := json.Marshal(snapshot)
		if first == "" {
			first = string(raw)
		} else if first != string(raw) {
			t.Fatal("metadata tie lacks total order")
		}
	}
}

type retainedDiscoveryPort struct{ value domain.ProviderDiscovery }

func (p retainedDiscoveryPort) Discover(context.Context, domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	return p.value, nil
}

func TestDiscoveryOwnsCanonicalCopiesAndDeduplicatesOrigins(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "gerenciar receitas"}, Expansions: regional.Proposals}
	var first string
	for _, reverse := range []bool{false, true} {
		value, _ := (canonicalDiscoveryPort{mode: "matches", reverse: reverse}).Discover(context.Background(), domain.DiscoveryRequest{Query: regional.Proposals[0].Query})
		// Same compared fields, distinct remaining metadata: exercise the full tie-break.
		a := value.Candidates[0]
		b := a
		b.Archived = true
		value.Candidates = []domain.RepositoryCandidate{a, b}
		if reverse {
			value.Candidates[0], value.Candidates[1] = value.Candidates[1], value.Candidates[0]
		}
		for i := range value.Matches {
			value.Matches[i].Evidence = []domain.MatchEvidence{{Path: "z.go", Snippet: "last"}, {Path: "a.go", Snippet: "first"}}
			if reverse {
				e := value.Matches[i].Evidence
				e[0], e[1] = e[1], e[0]
			}
		}
		value.Matches = append(value.Matches, value.Matches[0])
		before, _ := json.Marshal(value)
		snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 5, 1, retainedDiscoveryPort{value})
		if err != nil {
			t.Fatal(err)
		}
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("provider response mutated")
		}
		if len(snapshot.Candidates) != 1 || len(snapshot.Candidates[0].Origins) != 1 || len(snapshot.Candidates[0].Origins[0].Evidence) != 2 {
			t.Fatal("origin union lost evidence or retained duplicate")
		}
		raw, _ := json.Marshal(snapshot)
		if first == "" {
			first = string(raw)
		} else if first != string(raw) {
			t.Fatal("full declaration/evidence permutation changed snapshot")
		}
	}
	for _, value := range []domain.ProviderDiscovery{{}, {Candidates: []domain.RepositoryCandidate{}, Matches: []domain.ProviderCandidateMatch{}}} {
		before, _ := json.Marshal(value)
		canonicalizeDiscoveryProvider(&value)
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("nil/empty representation changed")
		}
	}
}
