/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/improvement_synthesis.md
 * @layer L4
 * @updated 2026-09-02
 */

package main

import (
	"context"
	_ "embed"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runstore"
)

func NewImprovementSynthesisMCPServer(
	funnel *usecase.SearchFunnel,
	recommendationExecutor domain.ObservableRecommendationExecutor,
	connectors []domain.MarketDiscoveryConnector,
	mechanismCatalog []byte,
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
	if _, err := extractor.ExtractMechanisms(context.Background(), domain.MechanismExtractionInput{RuleCatalogJSON: mechanismCatalog, RuleCatalogSHA256: mechanismRulesSHA256}); err != nil {
		return nil, nil, err
	}
	service, queue, store, err := newDurableComponents(recommendationExecutor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service.ConfigureMarketDiscovery(nil, registry)
	service.ConfigureMechanismExtraction(extractor, mechanismCatalog, mechanismRulesSHA256, baseline)
	service.ConfigureImprovementSynthesis(mechanismRulesSHA256)
	if err := service.Recover(context.Background()); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	if err := queue.Start(); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	server := mcp.NewMCPServer(funnel, reader, writer, recommendationExecutor).
		WithDurableRuns(service).
		WithMarketDiscovery(connectors, service).
		WithImprovementSynthesis(service)
	return server, durableShutdown(queue, store), nil
}

const mechanismRulesSHA256 = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"

//go:embed mechanism-rules.json
var embeddedMechanismCatalog []byte
