/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery_step0068_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"sync"
	"testing"

	"bateia/01_core/domain"
)

const step0068SourceQuery = "predicados geometricos de orientacao e in-circle robustos a casos quase degenerados"

func step0068TechnicalExpansion(id, query, rule string, discriminators ...string) domain.SearchExpansion {
	return domain.SearchExpansion{
		ID: id, Facet: domain.FacetCapability, Query: query,
		Discriminators: append([]string(nil), discriminators...),
		Covers:         []string{"capabilities[0]"},
		Derivation: &domain.ExpansionDerivation{
			Strategy: "technical-lexicon-v1", RuleIDs: []string{rule}, SourceIDs: []string{"exp-source"},
		},
	}
}

func step0068DiscoveryPlan() domain.ExpansionPlan {
	return domain.ExpansionPlan{
		Idea: domain.IdeaSpec{Capabilities: []string{step0068SourceQuery}},
		Expansions: []domain.SearchExpansion{
			{ID: "exp-source", Facet: domain.FacetCapability, Query: step0068SourceQuery, Covers: []string{"capabilities[0]"}},
			step0068TechnicalExpansion("exp-translation", "robust geometric orientation and in-circle predicates for nearly degenerate cases", "robust-predicates.translation.en", "orientation", "in-circle"),
			step0068TechnicalExpansion("exp-symbol", "orient2d incircle", "robust-predicates.symbol-pair", "orient2d", "incircle"),
			step0068TechnicalExpansion("exp-reference", "Shewchuk exact predicates", "robust-predicates.reference-pair", "Shewchuk", "predicates"),
		},
	}
}

type step0068ReplayProvider struct {
	mu       sync.Mutex
	requests []domain.DiscoveryRequest
	mutate   bool
}

func step0068CloneDerivation(value *domain.ExpansionDerivation) *domain.ExpansionDerivation {
	if value == nil {
		return nil
	}
	return &domain.ExpansionDerivation{
		Strategy: value.Strategy, RuleIDs: append([]string(nil), value.RuleIDs...), SourceIDs: append([]string(nil), value.SourceIDs...),
	}
}

func (provider *step0068ReplayProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.mu.Lock()
	captured := request
	captured.RequirementIDs = append([]string(nil), request.RequirementIDs...)
	captured.Discriminators = append([]string(nil), request.Discriminators...)
	captured.Derivation = step0068CloneDerivation(request.Derivation)
	provider.requests = append(provider.requests, captured)
	provider.mu.Unlock()

	if provider.mutate && request.Derivation != nil {
		request.Derivation.Strategy = "mutated"
		request.Derivation.RuleIDs[0] = "mutated-rule"
		request.Derivation.SourceIDs[0] = "mutated-source"
		request.Discriminators[0] = "mutated-discriminator"
	}

	attempt := domain.ProviderAttempt{Query: request.Query, Status: "EMPTY"}
	result := domain.ProviderDiscovery{Provider: "step0068-replay", Attempts: []domain.ProviderAttempt{attempt}}
	var candidate domain.RepositoryCandidate
	switch captured.Query {
	case step0068SourceQuery:
		candidate = domain.RepositoryCandidate{Owner: "noise", Repo: "Whisk"}
	case "orient2d incircle", "Shewchuk exact predicates":
		candidate = domain.RepositoryCandidate{Owner: "geometry", Repo: "RobustKernel"}
	default:
		return result, nil
	}
	result.Candidates = []domain.RepositoryCandidate{candidate}
	result.Attempts[0].Status, result.Attempts[0].ResultCount = "SUCCEEDED", 1
	if captured.Derivation != nil {
		result.Matches = []domain.ProviderCandidateMatch{{
			Repository: candidate.Owner + "/" + candidate.Repo,
			Evidence:   []domain.MatchEvidence{{MatchType: "file", Path: "predicates.cpp", Snippet: captured.Discriminators[0] + " " + captured.Discriminators[1]}},
		}}
	}
	return result, nil
}

func (provider *step0068ReplayProvider) capturedRequests() []domain.DiscoveryRequest {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	return append([]domain.DiscoveryRequest(nil), provider.requests...)
}

func TestStep0068DiscoveryDeepCopiesDerivationIntoRequest(t *testing.T) {
	plan := step0068DiscoveryPlan()
	before, _ := json.Marshal(plan)
	provider := &step0068ReplayProvider{mutate: true}
	snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "lexical", "cpp", 5, 1, provider)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(plan)
	if string(before) != string(after) {
		t.Fatalf("provider conseguiu mutar plano de entrada:\n%s\n%s", before, after)
	}
	for _, expansion := range snapshot.PerExpansion {
		if expansion.Expansion.Derivation != nil {
			derivation := expansion.Expansion.Derivation
			if derivation.Strategy != "technical-lexicon-v1" || derivation.RuleIDs[0] == "mutated-rule" || derivation.SourceIDs[0] == "mutated-source" || expansion.Expansion.Discriminators[0] == "mutated-discriminator" {
				t.Fatalf("request compartilhou slices/pointer com snapshot: %+v", expansion.Expansion)
			}
		}
	}
	requests := provider.capturedRequests()
	if len(requests) != 4 {
		t.Fatalf("requests=%d, quero literal mais três técnicas", len(requests))
	}
	for index := 1; index < len(requests); index++ {
		if requests[index].Derivation == nil || requests[index].Derivation.Strategy != "technical-lexicon-v1" || len(requests[index].Derivation.RuleIDs) != 1 || len(requests[index].Derivation.SourceIDs) != 1 {
			t.Fatalf("request[%d] perdeu derivação: %+v", index, requests[index])
		}
	}
}

func TestStep0068CounterfactualUtilityExecutesAndUnionsEveryExpansion(t *testing.T) {
	provider := &step0068ReplayProvider{}
	snapshot, err := DiscoverExpansionPlan(context.Background(), step0068DiscoveryPlan(), "lexical", "cpp", 5, 1, provider)
	if err != nil {
		t.Fatal(err)
	}
	requests := provider.capturedRequests()
	if len(requests) != 4 || len(snapshot.PerExpansion) != 4 {
		t.Fatalf("custo observado incorreto: requests=%d per_expansion=%d", len(requests), len(snapshot.PerExpansion))
	}
	wantQueries := []string{step0068SourceQuery, "robust geometric orientation and in-circle predicates for nearly degenerate cases", "orient2d incircle", "Shewchuk exact predicates"}
	for index := range wantQueries {
		if requests[index].Query != wantQueries[index] {
			t.Fatalf("request[%d]=%q, quero %q", index, requests[index].Query, wantQueries[index])
		}
	}
	for _, forbidden := range []string{"predicados geometricos", "geometric predicates"} {
		for _, request := range requests[1:] {
			if request.Query == forbidden {
				t.Fatalf("consulta genérica proibida executada: %q", request.Query)
			}
		}
	}

	if len(snapshot.Candidates) != 2 {
		t.Fatalf("união deveria conter ruído literal e alvo: %+v", snapshot.Candidates)
	}
	var target *domain.DiscoveredCandidate
	for index := range snapshot.Candidates {
		if snapshot.Candidates[index].Repository.Owner == "geometry" && snapshot.Candidates[index].Repository.Repo == "RobustKernel" {
			target = &snapshot.Candidates[index]
		}
	}
	if target == nil || len(target.ExpansionIDs) != 2 || len(target.Queries) != 2 || len(target.Origins) != 2 {
		t.Fatalf("alvo não preservou duas origens técnicas: %+v", target)
	}
	for _, origin := range target.Origins {
		if origin.PlannedQuery == "" || origin.EffectiveQuery == "" || len(origin.Evidence) == 0 {
			t.Fatalf("origem técnica incompleta: %+v", origin)
		}
	}

	// Ground truth externo congelado: RobustKernel relevante, Whisk irrelevante.
	recallAt5 := 1.0
	contaminationAt5 := 1.0 / float64(len(snapshot.Candidates))
	addressable := 0
	if len(target.Origins) >= 2 {
		addressable = 1
	}
	if recallAt5 != 1 || contaminationAt5 != 0.5 || addressable != 1 {
		t.Fatalf("métricas candidatas inesperadas: recall=%v contamination=%v addressable=%d", recallAt5, contaminationAt5, addressable)
	}
	baselineRecall, baselineContamination, baselineAddressable, baselineAttempts := 0.0, 1.0, 0, 1
	if !(recallAt5 > baselineRecall && contaminationAt5 <= baselineContamination && addressable > baselineAddressable && len(requests) == 4 && baselineAttempts == 1) {
		t.Fatalf("delta A/B externo não observado")
	}
}

func TestStep0068CandidateUnionIsCanonicalAcrossTwentyInputPermutations(t *testing.T) {
	base := step0068DiscoveryPlan()
	var canonical []domain.DiscoveredCandidate
	for repetition := 0; repetition < 20; repetition++ {
		plan := base
		plan.Expansions = append([]domain.SearchExpansion(nil), base.Expansions...)
		rotation := repetition % len(plan.Expansions)
		plan.Expansions = append(plan.Expansions[rotation:], plan.Expansions[:rotation]...)
		if repetition%2 == 1 {
			for left, right := 0, len(plan.Expansions)-1; left < right; left, right = left+1, right-1 {
				plan.Expansions[left], plan.Expansions[right] = plan.Expansions[right], plan.Expansions[left]
			}
		}
		provider := &step0068ReplayProvider{}
		snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "lexical", "cpp", 5, 1, provider)
		if err != nil {
			t.Fatal(err)
		}
		if len(provider.capturedRequests()) != 4 {
			t.Fatalf("permutação %d não executou quatro expansões", repetition)
		}
		for index := range snapshot.Candidates {
			sort.Strings(snapshot.Candidates[index].ExpansionIDs)
			sort.Strings(snapshot.Candidates[index].Queries)
		}
		if repetition == 0 {
			canonical = snapshot.Candidates
			continue
		}
		if !reflect.DeepEqual(snapshot.Candidates, canonical) {
			t.Fatalf("permutação %d mudou união canônica:\n%+v\n%+v", repetition, canonical, snapshot.Candidates)
		}
	}
}
