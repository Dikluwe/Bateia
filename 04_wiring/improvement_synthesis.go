/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/improvement_synthesis.md
 * @layer L4
 * @updated 2026-09-19
 */

package main

import (
	"context"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runcontainer"
	"bateia/03_infra/runstore"
)

func NewImprovementSynthesisMCPServer(
	funnel *usecase.SearchFunnel,
	recommendationExecutor domain.ObservableRecommendationExecutor,
	connectors []domain.MarketDiscoveryConnector,
	mechanismCatalog []byte,
	mechanismCatalogSHA256 string,
	baseline *domain.ConventionalBaseline,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
	reader io.Reader,
	writer io.Writer,
) (*mcp.MCPServer, func(context.Context) error, error) {
	registry, err := usecase.NewDiscoveryConnectorRegistry(connectors)
	if err != nil || len(registry.Capabilities()) != 3 {
		return nil, nil, domain.ErrInvalidRunInput
	}
	extractor := usecase.NewDeterministicMechanismExtractor()
	if _, err := extractor.ExtractMechanisms(context.Background(), domain.MechanismExtractionInput{RuleCatalogJSON: mechanismCatalog, RuleCatalogSHA256: mechanismCatalogSHA256}); err != nil {
		return nil, nil, err
	}
	durable, err := runcontainer.New(recommendationExecutor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service := durable.Service()
	service.ConfigureMarketDiscovery(nil, registry)
	service.ConfigureMechanismExtraction(extractor, mechanismCatalog, mechanismCatalogSHA256, baseline)
	service.ConfigureImprovementSynthesis(mechanismCatalogSHA256)
	if err := durable.Recover(context.Background()); err != nil {
		return nil, nil, err
	}
	if err := durable.Start(); err != nil {
		return nil, nil, err
	}
	server := mcp.NewMCPServer(funnel, reader, writer, recommendationExecutor).
		WithDurableRuns(service).
		WithMarketDiscovery(connectors, service).
		WithImprovementSynthesis(service)
	return server, durable.Shutdown, nil
}
