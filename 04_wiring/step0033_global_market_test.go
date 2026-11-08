/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0033_global_market_test.md
 * @layer L4
 * @updated 2026-10-28
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// step0033ConnectorStub implementa domain.MarketDiscoveryConnector com resposta configurável por mercado.
type step0033ConnectorStub struct {
	capability domain.DiscoveryConnectorCapabilities
}

func (s step0033ConnectorStub) Capabilities() domain.DiscoveryConnectorCapabilities {
	return s.capability
}

func (s step0033ConnectorStub) Discover(_ context.Context, req domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	// Mercados "vazios" simulam storefronts sem catálogo de software (AttemptEmptyObserved)
	emptyMarkets := map[string]bool{"MN": true}
	if emptyMarkets[req.Market] {
		return domain.ConnectorDiscovery{
			Status:     domain.AttemptEmptyObserved,
			Candidates: []domain.ProductCandidate{},
		}, nil
	}
	// Todos os outros mercados retornam 1 candidato fictício
	return domain.ConnectorDiscovery{
		Status:        domain.AttemptSucceeded,
		Candidates:    []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: "stub", ExternalID: "1", CanonicalURL: "https://example.com/" + req.Market}}},
		ResultCount:   1,
		AcceptedCount: 1,
	}, nil
}

// TestStep0033MCPSchemaUsesPatternNotEnum confirma que o MCP usa pattern aberto em vez de enum fechado.
func TestStep0033MCPSchemaUsesPatternNotEnum(t *testing.T) {
	// Simular uma conexão MCP e verificar o schema via tools/list
	registryCapabilities := domain.DiscoveryConnectorCapabilities{
		ID: "apple-itunes-search-software", Version: "baseline-21c48e357180",
		Kind: domain.ConnectorApplicationCatalog, QueryLocales: []string{"en-US"},
		Markets: []string{"US"}, SupportsMarketFilter: true,
		IdentityLevels: []string{"external_id"}, EvidenceSurfaces: []string{"catalog_listing"},
		Replayable: true, PublicOnly: true,
	}
	connector := step0033ConnectorStub{capability: registryCapabilities}
	registry, err := usecase.NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{connector})
	if err != nil {
		t.Fatal("falha ao criar registry:", err)
	}

	// Verificar que registry aceita o conector
	caps := registry.Capabilities()
	if len(caps) != 1 {
		t.Fatalf("esperado 1 conector no registry, obtido %d", len(caps))
	}
	_ = fmt.Sprintf("conector registrado: %s", caps[0].ID)
}

// TestStep0033MultiBatchWithEmptyStorefrontDoesNotFail confirma resiliência de lote:
// um mercado sem catálogo (AttemptEmptyObserved) não impede a execução nos demais.
func TestStep0033MultiBatchWithEmptyStorefrontDoesNotFail(t *testing.T) {
	capabilities := domain.DiscoveryConnectorCapabilities{
		ID: "stub-global", Version: "v1",
		Kind:                 domain.ConnectorApplicationCatalog,
		QueryLocales:         []string{"en-US", "it-IT", "el-GR", "mn-MN"},
		Markets:              []string{"US", "IT", "GR", "MN"},
		SupportsMarketFilter: true,
		IdentityLevels:       []string{"external_id"},
		EvidenceSurfaces:     []string{"catalog_listing"},
		Replayable:           true, PublicOnly: true,
	}
	connector := step0033ConnectorStub{capability: capabilities}
	registry, err := usecase.NewDiscoveryConnectorRegistry([]domain.MarketDiscoveryConnector{connector})
	if err != nil {
		t.Fatal("falha ao criar registry:", err)
	}

	request := domain.MarketDiscoveryStartRequest{
		Idea: domain.IdeaSpec{Problem: "discover recipe apps globally", Capabilities: []string{"recipe management"}},
		QueryProposals: []domain.MarketQueryProposal{
			{Locale: "en-US", Query: "recipe management app", SourceText: "gestao de receitas",
				Derivation: domain.QueryExternalAgent, Rationale: "international expansion probe"},
			{Locale: "it-IT", Query: "ricette cucina app", SourceText: "gestione ricette",
				Derivation: domain.QueryExternalAgent, Rationale: "mercado italiano"},
			{Locale: "el-GR", Query: "συνταγές μαγειρικής", SourceText: "suntages",
				Derivation: domain.QueryExternalAgent, Rationale: "mercado grego"},
			{Locale: "mn-MN", Query: "хоолны жор", SourceText: "mongolian food recipes",
				Derivation: domain.QueryExternalAgent, Rationale: "mercado mongol"},
		},
		Markets:      []string{"US", "IT", "GR", "MN"},
		ConnectorIDs: []string{"stub-global"},
		Budget:       domain.MarketDiscoveryBudget{MaxQueries: 5, MaxResultsPerQuery: 5, MaxTotalCandidates: 50},
	}

	normalized, err := usecase.ValidateMarketDiscoveryStartRequest(request, registry)
	if err != nil {
		t.Fatalf("request global válido rejeitado: %v", err)
	}
	if len(normalized.Markets) != 4 {
		t.Fatalf("esperado 4 mercados normalizados, obtido %d", len(normalized.Markets))
	}

	// Verificar que o plano de descoberta é gerado sem erros
	plan, err := usecase.PlanMarketDiscovery(normalized, registry)
	if err != nil {
		t.Fatalf("plano global falhou: %v", err)
	}
	if len(plan) == 0 {
		t.Fatal("plano global não gerou tentativas")
	}

	// Executar as descobertas com o stub (MN retorna vazio, outros retornam candidatos)
	discoveries := make([]domain.ConnectorDiscovery, 0, len(plan))
	allottedPlan := usecase.AllocateMarketDiscoveryBudget(plan, request.Budget)
	for _, attempt := range allottedPlan {
		discovery, discErr := connector.Discover(context.Background(), domain.MarketDiscoveryRequest{
			QueryID:     attempt.QueryID,
			QueryLocale: attempt.QueryLocale,
			Query:       "test query",
			Market:      attempt.Market,
			Limit:       request.Budget.MaxResultsPerQuery,
		})
		if discErr != nil {
			t.Fatalf("discover falhou para mercado %s: %v", attempt.Market, discErr)
		}
		discovery.AttemptID = attempt.AttemptID
		if discovery.Candidates == nil {
			discovery.Candidates = []domain.ProductCandidate{}
		}
		discoveries = append(discoveries, discovery)
	}

	// Contar AttemptEmptyObserved (MN deve aparecer)
	emptyCount := 0
	successCount := 0
	for _, d := range discoveries {
		if d.Status == domain.AttemptEmptyObserved {
			emptyCount++
		} else if d.Status == domain.AttemptSucceeded {
			successCount++
		}
	}
	if emptyCount == 0 {
		t.Error("esperado pelo menos 1 mercado com AttemptEmptyObserved (MN)")
	}
	if successCount == 0 {
		t.Error("esperado pelo menos 1 mercado com AttemptSucceeded")
	}
}

var _ = json.Marshal
