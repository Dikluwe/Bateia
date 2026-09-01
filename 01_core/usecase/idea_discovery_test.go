/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"bateia/01_core/domain"
)

type fakeObservableProvider struct {
	mu                  sync.Mutex
	active, peak, calls int
	results             map[string][]domain.RepositoryCandidate
	failures            map[string]bool
}

type deadlineProvider struct{}

func (deadlineProvider) Discover(ctx context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	if request.ExpansionID == "exp-a" {
		return domain.ProviderDiscovery{Provider: "deadline", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "SUCCEEDED", ResultCount: 1}}, Candidates: []domain.RepositoryCandidate{candidate("org", "kept")}}, nil
	}
	<-ctx.Done()
	return domain.ProviderDiscovery{Provider: "deadline", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "TIMED_OUT", ErrorClass: "TIMEOUT", Error: ctx.Err().Error()}}}, ctx.Err()
}

func (provider *fakeObservableProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.mu.Lock()
	provider.active++
	provider.calls++
	if provider.active > provider.peak {
		provider.peak = provider.active
	}
	provider.mu.Unlock()
	defer func() { provider.mu.Lock(); provider.active--; provider.mu.Unlock() }()
	result := domain.ProviderDiscovery{Provider: "fake", Attempts: []domain.ProviderAttempt{{Query: request.Query}}, Candidates: provider.results[request.ExpansionID]}
	result.Attempts[0].ResultCount = len(result.Candidates)
	if provider.failures[request.ExpansionID] {
		result.Attempts[0].Error = "failed"
		return result, errors.New("provider failed")
	}
	return result, nil
}

func discoveryPlan() domain.ExpansionPlan {
	return domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "target"}, Expansions: []domain.SearchExpansion{
		{ID: "exp-a", Facet: domain.FacetProblem, Query: "query a", Covers: []string{"problem"}},
		{ID: "exp-b", Facet: domain.FacetCapability, Query: "query b", Covers: []string{"capabilities[0]"}},
		{ID: "exp-c", Facet: domain.FacetCapability, Query: "query c", Covers: []string{"capabilities[1]"}},
	}}
}

func candidate(owner, repo string) domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: owner, Repo: repo}
}

func TestDiscoverExpansionPlanDeduplicatesOriginsAndKeepsPartialFailure(t *testing.T) {
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{
		"exp-a": {candidate("Org", "A"), candidate("Org", "B")},
		"exp-b": {candidate("org", "b"), candidate("Org", "C")},
	}, failures: map[string]bool{"exp-c": true}}
	snapshot, err := DiscoverExpansionPlan(context.Background(), discoveryPlan(), "lexical", "go", 10, 2, provider)
	if err != nil {
		t.Fatal(err)
	}
	if provider.calls != 3 || provider.peak > 2 {
		t.Fatalf("chamadas/concorrência incorretas: calls=%d peak=%d", provider.calls, provider.peak)
	}
	if len(snapshot.Candidates) != 3 || snapshot.Candidates[0].Repository.Repo != "A" || snapshot.Candidates[1].Repository.Repo != "B" || snapshot.Candidates[2].Repository.Repo != "C" {
		t.Fatalf("ordem/deduplicação incorreta: %+v", snapshot.Candidates)
	}
	if len(snapshot.Candidates[1].ExpansionIDs) != 2 || len(snapshot.Candidates[1].Queries) != 2 {
		t.Fatalf("origens de B perdidas: %+v", snapshot.Candidates[1])
	}
	if len(snapshot.Failures) != 1 || snapshot.Failures[0].Scope != "exp-c" {
		t.Fatalf("falha parcial perdida: %+v", snapshot.Failures)
	}
}

func TestDiscoverExpansionPlanFailsWhenAllExpansionsFail(t *testing.T) {
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{}, failures: map[string]bool{"exp-a": true, "exp-b": true, "exp-c": true}}
	if _, err := DiscoverExpansionPlan(context.Background(), discoveryPlan(), "literal", "go", 10, 1, provider); err == nil {
		t.Fatal("falha total foi aceita")
	}
}

func TestDiscoverExpansionPlanDoesNotMutateInput(t *testing.T) {
	plan := discoveryPlan()
	original := clonePlan(plan)
	provider := &fakeObservableProvider{results: map[string][]domain.RepositoryCandidate{}, failures: map[string]bool{}}
	if _, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 10, 1, provider); err != nil {
		t.Fatal(err)
	}
	if plan.Expansions[0].Query != original.Expansions[0].Query || plan.Idea.Problem != original.Idea.Problem {
		t.Fatal("plano foi modificado")
	}
}

func TestBoundedDiscoveryPreservesCandidateAndMarksNotStarted(t *testing.T) {
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: 40 * time.Millisecond, AttemptTimeout: 20 * time.Millisecond, Concurrency: 1, MaxAttempts: 3}
	snapshot, err := DiscoverExpansionPlanWithPolicy(context.Background(), discoveryPlan(), "literal", "", 3, policy, deadlineProvider{}, nil)
	if err == nil {
		t.Fatal("deadline must be observable")
	}
	if len(snapshot.Candidates) != 1 || snapshot.Candidates[0].Repository.Repo != "kept" {
		t.Fatalf("completed candidate was lost: %#v", snapshot.Candidates)
	}
	if snapshot.Receipt == nil || snapshot.Receipt.Status != domain.DiscoveryRunPartial || snapshot.Receipt.CompletionCause != domain.CompletionTotalTimeout {
		t.Fatalf("unexpected receipt: %#v", snapshot.Receipt)
	}
	notStarted := 0
	for _, a := range snapshot.Receipt.Attempts {
		if a.Status == "NOT_STARTED" {
			notStarted++
		}
	}
	if notStarted != 1 || snapshot.Receipt.TimedOut != 1 {
		t.Fatalf("timeout/not-started misclassified: %#v", snapshot.Receipt)
	}
}

func TestDiscoveryPolicyRejectsInvalidTimeoutsBeforeProvider(t *testing.T) {
	provider := &fakeObservableProvider{}
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: time.Second, AttemptTimeout: 2 * time.Second, Concurrency: 1, MaxAttempts: 3}
	if _, err := DiscoverExpansionPlanWithPolicy(context.Background(), discoveryPlan(), "literal", "", 3, policy, provider, nil); err == nil {
		t.Fatal("attempt timeout greater than total survived")
	}
	if provider.calls != 0 {
		t.Fatal("provider was touched before policy validation")
	}
}

func TestM21RegionalCatalogAllIsExactAuditableAndCanonical(t *testing.T) {
	wantLocales := []string{"pt-BR", "en-US", "es-ES", "fr-FR", "de-DE", "ja-JP", "ko-KR", "ru-RU", "ar-SA", "hi-IN"}
	wantQueries := []string{"gerenciador de receitas", "recipe manager", "gestor de recetas", "gestionnaire de recettes", "Rezeptmanager", "レシピ管理", "레시피 관리자", "менеджер рецептов", "مدير الوصفات", "रेसिपी प्रबंधक"}
	plan, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"ja-JP", "all", "pt-BR", "all"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plan.Requested, wantLocales) || !reflect.DeepEqual(plan.Realized, wantLocales) || len(plan.Gaps) != 0 || len(plan.Proposals) != len(wantLocales) {
		t.Fatalf("canonical all plan = %+v", plan)
	}
	ids := map[string]bool{}
	for index, proposal := range plan.Proposals {
		if proposal.Query != wantQueries[index] || proposal.Regional == nil {
			t.Fatalf("proposal[%d] = %+v, want query %q", index, proposal, wantQueries[index])
		}
		regional := proposal.Regional
		if regional.Locale != wantLocales[index] || regional.CatalogVersion != "regional-recipes-v1" || regional.Concept != "recipe-manager" || regional.RuleID != "recipe-manager" || regional.SourceQuery != "gerenciador de receitas" || regional.DerivationStrategy != "controlled-catalog-exact-v1" {
			t.Fatalf("proposal[%d] lost audit provenance: %+v", index, regional)
		}
		if !strings.Contains(proposal.ID, wantLocales[index]) || proposal.ID != regional.ProposalID || ids[proposal.ID] {
			t.Fatalf("localized proposal id collided or hid locale: %+v", proposal)
		}
		ids[proposal.ID] = true
	}
	for _, asciiForbidden := range []struct {
		locale string
		query  string
	}{
		{"ja-JP", "recipe kanri"}, {"ko-KR", "resipi"}, {"ru-RU", "menedzher"}, {"ar-SA", "mudir"}, {"hi-IN", "recipe prabandhak"},
	} {
		for _, proposal := range plan.Proposals {
			if proposal.Regional != nil && proposal.Regional.Locale == asciiForbidden.locale && strings.EqualFold(proposal.Query, asciiForbidden.query) {
				t.Fatalf("%s was transliterated: %+v", asciiForbidden.locale, proposal)
			}
		}
	}
}

func TestM21RegionalInputRejectsInvalidAndPreservesExplicitUnknown(t *testing.T) {
	catalog := DefaultRegionalQueryCatalog()
	for _, locales := range [][]string{{"pt-br"}, {"en"}, {"zz-ZZ"}, {""}} {
		if _, err := PlanRegionalDiscovery(catalog, "recipe-manager", locales); err == nil {
			t.Fatalf("invalid locale accepted: %q", locales)
		}
	}

	missing := catalog
	missing.Concepts = append([]domain.RegionalCatalogConcept(nil), catalog.Concepts...)
	missing.Concepts[0].Queries = append([]domain.RegionalCatalogQuery(nil), catalog.Concepts[0].Queries[:len(catalog.Concepts[0].Queries)-1]...)
	plan, err := PlanRegionalDiscovery(missing, "recipe-manager", []string{"hi-IN"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Realized) != 0 || len(plan.Proposals) != 0 || len(plan.Gaps) != 1 || plan.Gaps[0].Locale != "hi-IN" || plan.Gaps[0].Concept != "recipe-manager" || plan.Gaps[0].Status != "Unknown" || strings.TrimSpace(plan.Gaps[0].Reason) == "" {
		t.Fatalf("missing controlled translation was not explicit Unknown: %+v", plan)
	}

	invalidCatalog := catalog
	invalidCatalog.Concepts = append([]domain.RegionalCatalogConcept(nil), catalog.Concepts...)
	invalidCatalog.Concepts[0].SourceQuery = ""
	if _, err := PlanRegionalDiscovery(invalidCatalog, "recipe-manager", []string{"pt-BR"}); err == nil {
		t.Fatal("catalog without source query was accepted")
	}
}

func TestM21RegionalPlanIsInvariantToReorderRepeatAndEmptyIsLegacyNeutral(t *testing.T) {
	catalog := DefaultRegionalQueryCatalog()
	first, err := PlanRegionalDiscovery(catalog, "recipe-manager", []string{"ja-JP", "pt-BR", "ja-JP", "en-US"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := PlanRegionalDiscovery(catalog, "recipe-manager", []string{"en-US", "pt-BR", "ja-JP"})
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Fatalf("locale reorder/repeat changed plan:\n%s\n%s", a, b)
	}
	empty, err := PlanRegionalDiscovery(catalog, "recipe-manager", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Requested) != 0 || len(empty.Realized) != 0 || len(empty.Gaps) != 0 || len(empty.Proposals) != 0 {
		t.Fatalf("empty locale set added regional work: %+v", empty)
	}
}

type m21RegionalProvider struct {
	requests   []domain.DiscoveryRequest
	reverse    bool
	failLocale string
}

func (provider *m21RegionalProvider) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	provider.requests = append(provider.requests, request)
	if request.SearchLocale == provider.failLocale {
		return domain.ProviderDiscovery{Provider: "regional", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "FAILED", Error: "regional failure"}}}, errors.New("regional failure")
	}
	candidates := []domain.RepositoryCandidate{candidate("regional", "shared"), candidate("regional", "other")}
	if provider.reverse {
		candidates[0], candidates[1] = candidates[1], candidates[0]
	}
	return domain.ProviderDiscovery{Provider: "regional", Attempts: []domain.ProviderAttempt{{Query: request.Query, ResultCount: len(candidates)}}, Candidates: candidates}, nil
}

func TestM21RegionalLocalePropagatesAndConsumesOrdinaryBudget(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"ja-JP", "pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	provider := &m21RegionalProvider{}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "organizar receitas culinárias"}, Expansions: regional.Proposals}
	snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 7, 1, provider)
	if err != nil {
		t.Fatal(err)
	}
	if len(provider.requests) != 2 || !reflect.DeepEqual(snapshot.SearchLocales, []string{"pt-BR", "ja-JP"}) {
		t.Fatalf("regional calls/locales = %d/%v", len(provider.requests), snapshot.SearchLocales)
	}
	for index, request := range provider.requests {
		if request.Language != "go" || request.SearchLocale == "" || request.RegionalProposalID == "" || request.SearchLocale == request.Language || request.MaxResults != 7 {
			t.Fatalf("request[%d] mixed language and locale: %+v", index, request)
		}
	}
	if len(snapshot.Candidates) != 2 {
		t.Fatalf("regional duplicate repository was not deduplicated: %+v", snapshot.Candidates)
	}
	for _, discovered := range snapshot.Candidates {
		if len(discovered.Origins) != 2 {
			t.Fatalf("regional origins lost: %+v", discovered)
		}
		for index, origin := range discovered.Origins {
			if origin.SearchLocale != []string{"pt-BR", "ja-JP"}[index] || origin.RegionalProposalID == "" || origin.PlannedQuery == "" || origin.EffectiveQuery == "" {
				t.Fatalf("origin[%d] lost provenance: %+v", index, origin)
			}
		}
	}

	reversed := &m21RegionalProvider{reverse: true}
	second, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 7, 1, reversed)
	if err != nil {
		t.Fatal(err)
	}
	left, _ := json.Marshal(snapshot)
	right, _ := json.Marshal(second)
	if string(left) != string(right) {
		t.Fatalf("provider response reorder changed snapshot:\n%s\n%s", left, right)
	}
}

func TestM21RegionalFailedAttemptPreservesLocaleProvenance(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"pt-BR", "ja-JP"})
	if err != nil {
		t.Fatal(err)
	}
	provider := &m21RegionalProvider{failLocale: "ja-JP"}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "organizar receitas culinárias"}, Expansions: regional.Proposals}
	snapshot, err := DiscoverExpansionPlan(context.Background(), plan, "literal", "go", 7, 1, provider)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, expansion := range snapshot.PerExpansion {
		for _, attempt := range expansion.Provider.Attempts {
			if attempt.SearchLocale == "ja-JP" {
				found = attempt.RegionalProposalID != "" && attempt.Query == "レシピ管理" && attempt.Error != ""
			}
		}
	}
	if !found || len(snapshot.Failures) != 1 {
		t.Fatalf("failed regional attempt lost locale provenance: %+v", snapshot)
	}
}

func TestM21RegionalPolicyReceiptAndOrdinaryLimitsAreEndToEnd(t *testing.T) {
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), "recipe-manager", []string{"pt-BR", "ja-JP"})
	if err != nil {
		t.Fatal(err)
	}
	provider := &m21RegionalProvider{}
	policy := domain.DiscoveryExecutionPolicy{SchemaVersion: "1", TotalTimeout: 5 * time.Second, AttemptTimeout: time.Second, Concurrency: 1, MaxAttempts: 3}
	plan := domain.ExpansionPlan{Idea: domain.IdeaSpec{Problem: "organizar receitas culinárias"}, Expansions: regional.Proposals}
	snapshot, err := DiscoverExpansionPlanWithPolicy(context.Background(), plan, "literal", "go", 7, policy, provider, nil)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Receipt == nil || len(snapshot.Receipt.Attempts) != 2 || snapshot.Receipt.Planned != 2 || snapshot.Receipt.Started != 2 {
		t.Fatalf("regional receipt counts = %+v", snapshot.Receipt)
	}
	wantLocales := []string{"pt-BR", "ja-JP"}
	for index, attempt := range snapshot.Receipt.Attempts {
		if attempt.SearchLocale != wantLocales[index] || attempt.RegionalProposalID == "" || attempt.PlannedQuery == "" || attempt.EffectiveQuery == "" {
			t.Fatalf("receipt attempt[%d] lost regional provenance: %+v", index, attempt)
		}
	}
	if len(provider.requests) != 2 {
		t.Fatalf("provider request count = %d, want 2", len(provider.requests))
	}
	for index, request := range provider.requests {
		if request.MaxResults != 7 || request.MaxAttempts != policy.MaxAttempts || request.AttemptTimeout != policy.AttemptTimeout || request.Language != "go" || request.SearchLocale == "" {
			t.Fatalf("regional request[%d] bypassed ordinary limits: %+v", index, request)
		}
	}
}
