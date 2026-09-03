/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_discovery_step0027_oracle_test.md
 * @layer L1
 * @updated 2026-09-01
 */

package usecase_test

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type oracleConnector struct {
	cap   domain.DiscoveryConnectorCapabilities
	calls *int
}

type oracleRecoveryStore struct {
	domain.RecommendationRunStore
	called bool
}

func (s *oracleRecoveryStore) InterruptNonTerminal(context.Context) ([]domain.RecommendationRun, error) {
	s.called = true
	return []domain.RecommendationRun{{RunID: "run_0123456789abcdef0123456789abcdef", RunKind: domain.RunKindMarketDiscovery, Lifecycle: domain.RunInterrupted}}, nil
}

func (c oracleConnector) Capabilities() domain.DiscoveryConnectorCapabilities { return c.cap }
func (c oracleConnector) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if c.calls != nil {
		*c.calls++
	}
	return domain.ConnectorDiscovery{}, nil
}

func oracleCapabilities(id string, filter bool) domain.DiscoveryConnectorCapabilities {
	kind := domain.ConnectorRepositoryCatalog
	markets := []string{}
	identity := []string{"owner_repo", "commit"}
	if filter {
		kind = domain.ConnectorApplicationCatalog
		markets = []string{"BR", "JP", "US"}
		identity = []string{"provider_product_id", "canonical_url"}
	}
	return domain.DiscoveryConnectorCapabilities{
		ID: id, Version: "v1", Kind: kind,
		QueryLocales: []string{"en-US", "ja-JP", "pt-BR"}, Markets: markets,
		SupportsMarketFilter: filter, IdentityLevels: identity,
		EvidenceSurfaces: []string{"name", "publisher"}, Replayable: true, PublicOnly: true,
	}
}

func oracleStartRequest() domain.MarketDiscoveryStartRequest {
	return domain.MarketDiscoveryStartRequest{
		Idea: domain.IdeaSpec{Problem: "reduzir desperdício alimentar doméstico", Capabilities: []string{"inventário", "validade"}},
		QueryProposals: []domain.MarketQueryProposal{
			{ID: "q-en", Locale: "en-US", Query: "household food inventory expiry", SourceText: "inventário e validade", Derivation: domain.QueryExternalAgent, Rationale: "controle cruzado", InputHash: "sha256:01"},
			{ID: "q-ja", Locale: "ja-JP", Query: "家庭 食品 在庫 期限", SourceText: "inventário e validade", Derivation: domain.QueryExternalAgent, Rationale: "formulação funcional", InputHash: "sha256:02"},
		},
		Markets:      []string{"JP", "US"},
		ConnectorIDs: []string{"apple-itunes-search-software", "github-repository-catalog", "sourcegraph-repository-catalog"},
		Budget:       domain.MarketDiscoveryBudget{MaxQueries: 12, MaxResultsPerQuery: 5, MaxTotalCandidates: 30},
	}
}

func TestStep0027OracleRegistryValidationAndCanonicalPlanUseIndependentDimensions(t *testing.T) {
	calls := 0
	connectors := []domain.MarketDiscoveryConnector{
		oracleConnector{cap: oracleCapabilities("sourcegraph-repository-catalog", false), calls: &calls},
		oracleConnector{cap: oracleCapabilities("apple-itunes-search-software", true), calls: &calls},
		oracleConnector{cap: oracleCapabilities("github-repository-catalog", false), calls: &calls},
	}
	registry, err := usecase.NewDiscoveryConnectorRegistry(connectors)
	if err != nil {
		t.Fatalf("valid registry: %v", err)
	}
	capabilities := registry.Capabilities()
	if calls != 0 {
		t.Fatalf("Capabilities caused %d network/discovery calls", calls)
	}
	gotIDs := []string{capabilities[0].ID, capabilities[1].ID, capabilities[2].ID}
	wantIDs := []string{"apple-itunes-search-software", "github-repository-catalog", "sourcegraph-repository-catalog"}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("canonical connector order = %v, want %v", gotIDs, wantIDs)
	}

	request := oracleStartRequest()
	normalized, err := usecase.ValidateMarketDiscoveryStartRequest(request, registry)
	if err != nil {
		t.Fatalf("valid cross-locale request: %v", err)
	}
	plan, err := usecase.PlanMarketDiscovery(normalized, registry)
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	if len(plan) != 12 { // 2 proposals × 2 markets × 3 connectors
		t.Fatalf("plan length = %d, want 12", len(plan))
	}
	crossPairs := map[string]bool{}
	for _, attempt := range plan {
		crossPairs[attempt.QueryLocale+"/"+attempt.Market] = true
		if attempt.ConnectorID == "apple-itunes-search-software" && !attempt.MarketFilterApplied {
			t.Fatalf("Apple attempt lost explicit storefront: %+v", attempt)
		}
		if attempt.ConnectorID != "apple-itunes-search-software" && attempt.Status != domain.AttemptMarketNotFiltered {
			t.Fatalf("global connector claimed a regional filter: %+v", attempt)
		}
		if attempt.AttemptID == "" || attempt.QueryID == "" || attempt.RequestHash == "" || attempt.QueryHash == "" || attempt.EffectiveQueryHash == "" {
			t.Fatalf("attempt lacks causal provenance: %+v", attempt)
		}
	}
	if !crossPairs["en-US/JP"] || !crossPairs["ja-JP/US"] {
		t.Fatalf("locale/market controls disappeared: %v", crossPairs)
	}

	reordered := oracleStartRequest()
	sort.Slice(reordered.QueryProposals, func(i, j int) bool { return reordered.QueryProposals[i].ID > reordered.QueryProposals[j].ID })
	sort.Sort(sort.Reverse(sort.StringSlice(reordered.Markets)))
	sort.Sort(sort.Reverse(sort.StringSlice(reordered.ConnectorIDs)))
	normalized2, hash2, err := usecase.CanonicalMarketDiscoveryStartRequestHash(reordered, registry)
	if err != nil {
		t.Fatalf("canonical reordered request: %v", err)
	}
	_, hash1, err := usecase.CanonicalMarketDiscoveryStartRequestHash(request, registry)
	if err != nil || hash1 != hash2 {
		t.Fatalf("reordering changed idempotency hash: %q / %q (%v)", hash1, hash2, err)
	}
	plan2, err := usecase.PlanMarketDiscovery(normalized2, registry)
	if err != nil || !reflect.DeepEqual(plan, plan2) {
		t.Fatalf("reordering changed canonical plan: %v", err)
	}

	bad := oracleCapabilities("apple-itunes-search-software", true)
	bad.PublicOnly = false
	if _, err := usecase.NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{oracleConnector{cap: bad}}); err == nil {
		t.Fatal("PublicOnly=false registry was accepted")
	}
	if _, err := usecase.NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{connectors[0], connectors[0]}); err == nil {
		t.Fatal("duplicate connector id was accepted")
	}
}

func TestStep0027OracleBudgetIdentityDimensionsAndSafeArtifactsAreConservative(t *testing.T) {
	registry, err := usecase.NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{
		oracleConnector{cap: oracleCapabilities("sourcegraph-repository-catalog", false)},
		oracleConnector{cap: oracleCapabilities("apple-itunes-search-software", true)},
		oracleConnector{cap: oracleCapabilities("github-repository-catalog", false)},
	})
	if err != nil {
		t.Fatal(err)
	}
	normalized, err := usecase.ValidateMarketDiscoveryStartRequest(oracleStartRequest(), registry)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := usecase.PlanMarketDiscovery(normalized, registry)
	if err != nil {
		t.Fatal(err)
	}
	allocated := usecase.AllocateMarketDiscoveryBudget(plan, domain.MarketDiscoveryBudget{MaxQueries: 3, MaxResultsPerQuery: 5, MaxTotalCandidates: 15})
	firstRound := map[string]int{}
	for _, attempt := range allocated {
		if attempt.PlannedCost > 0 {
			firstRound[attempt.ConnectorID]++
		}
	}
	if len(firstRound) != 3 {
		t.Fatalf("first connector consumed budget before fair round: %v", firstRound)
	}

	base := domain.ProductCandidate{
		Identity: domain.ProductIdentity{ConnectorID: "apple-itunes-search-software", ExternalID: "1001", CanonicalURL: "https://apps.apple.com/app/id1001"},
		Name:     "Resource Ledger", ProductKind: "software", Publisher: "Publisher A",
		ListingLocales: []string{}, ObservedMarkets: []string{"JP"}, PublisherCountries: []string{},
		RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: "org/project", Commit: "0123456789abcdef0123456789abcdef01234567", Status: domain.RepositoryLinkVerified, Source: "official"}},
		Origins:         []domain.MarketCandidateOrigin{{ConnectorID: "apple-itunes-search-software", ConnectorVersion: "v1", QueryID: "q-en", QueryLocale: "en-US", ObservedMarket: "JP", ResponseHash: "sha256:body"}},
	}
	strong := base
	strong.Identity = domain.ProductIdentity{ConnectorID: "github-repository-catalog", ExternalID: "org/project"}
	resolution, witnesses := usecase.ResolveProductIdentity(base, strong)
	if resolution != domain.IdentitySame || len(witnesses) == 0 {
		t.Fatalf("bilateral verified owner/repo@commit not SAME: %s %v", resolution, witnesses)
	}
	homonym := base
	homonym.Identity = domain.ProductIdentity{ConnectorID: "other", ExternalID: "other-7", CanonicalURL: base.Identity.CanonicalURL}
	homonym.Publisher = "Publisher B"
	homonym.RepositoryLinks = []domain.ProductRepositoryLink{}
	resolution, _ = usecase.ResolveProductIdentity(base, homonym)
	if resolution == domain.IdentitySame {
		t.Fatal("name/canonical URL without bilateral witness caused SAME")
	}

	encoded, err := json.Marshal(base)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, exact := range []string{`"listing_locales":[]`, `"observed_markets":["JP"]`, `"publisher_countries":[]`, `"repository_links"`, `"origins"`} {
		if !strings.Contains(text, exact) {
			t.Fatalf("missing non-null independent dimension %s in %s", exact, text)
		}
	}
	for _, forbidden := range []string{"rating", "downloads", "review", "html", "top_choice", "PROVEN"} {
		if strings.Contains(strings.ToLower(text), strings.ToLower(forbidden)) {
			t.Fatalf("unsafe/evidence-promoting field %q in candidate JSON: %s", forbidden, text)
		}
	}
}

func TestStep0027OracleMarketLifecycleRecoveryAndStatusUseTheCertifiedRuntime(t *testing.T) {
	run := domain.RecommendationRun{
		SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef", RunKind: domain.RunKindMarketDiscovery,
		Lifecycle: domain.RunQueued, Version: 1, Stage: domain.RunStageQueued,
		Progress: domain.RecommendationRunProgress{Planned: 4, Completed: 0}, Artifacts: []domain.RunArtifactIndex{},
	}
	start := domain.RunTransition{
		RunID: run.RunID, ExpectedVersion: 1, Lifecycle: domain.RunRunning,
		Stage: domain.RunStageMarketPlan, Progress: domain.RecommendationRunProgress{Planned: 4, Completed: 0},
	}
	if err := usecase.ValidateRecommendationRunTransition(run, start); err != nil {
		t.Fatalf("shared lifecycle rejected market start: %v", err)
	}
	status := domain.MarketDiscoveryComplete
	running := run
	running.Lifecycle = domain.RunRunning
	running.Stage = domain.RunStageMarketPublish
	running.Progress = domain.RecommendationRunProgress{Planned: 4, Completed: 4}
	running.Artifacts = []domain.RunArtifactIndex{
		{Name: domain.ArtifactConnectorCapabilities, State: domain.ArtifactPublished},
		{Name: domain.ArtifactMarketPlan, State: domain.ArtifactPublished},
		{Name: domain.ArtifactMarketSnapshot, State: domain.ArtifactPublished},
		{Name: domain.ArtifactMarketCandidates, State: domain.ArtifactPublished},
		{Name: domain.ArtifactMarketGaps, State: domain.ArtifactPublished},
	}
	finish := domain.RunTransition{
		RunID: running.RunID, ExpectedVersion: running.Version, Lifecycle: domain.RunSucceeded,
		Stage: domain.RunStageTerminal, Progress: running.Progress, MarketDiscoveryStatus: &status,
	}
	if err := usecase.ValidateRecommendationRunTransition(running, finish); err != nil {
		t.Fatalf("shared lifecycle rejected market completion: %v", err)
	}
	recommendationStatus := domain.IdeaRunComplete
	finish.MarketDiscoveryStatus = nil
	finish.RecommendationStatus = &recommendationStatus
	if err := usecase.ValidateRecommendationRunTransition(running, finish); err == nil {
		t.Fatal("market run accepted recommendation_status instead of market_discovery_status")
	}

	store := &oracleRecoveryStore{}
	service := usecase.NewRecommendationRunService(store, nil, nil, nil)
	if err := service.Recover(context.Background()); err != nil || !store.called {
		t.Fatalf("market recovery did not use the certified store path: called=%v err=%v", store.called, err)
	}
}
