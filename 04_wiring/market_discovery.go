/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/market_discovery.md
 * @layer L4
 * @updated 2026-09-06
 */

package main

import (
	"context"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/fdroid"
	"bateia/03_infra/forge"
	"bateia/03_infra/itunes"
	"bateia/03_infra/playstore"
	"bateia/03_infra/runstore"
)

func defaultMarketDiscoveryConnectors() ([]domain.MarketDiscoveryConnector, error) {
	apple, err := itunes.NewMarketDiscoveryConnector(nil)
	if err != nil {
		return nil, err
	}
	locales := []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"}
	repository := func(id string) domain.MarketDiscoveryConnector {
		return usecase.CapabilityOnlyConnector{Capability: domain.DiscoveryConnectorCapabilities{
			ID: id, Version: "baseline-21c48e357180", Kind: domain.ConnectorRepositoryCatalog,
			QueryLocales: append([]string{}, locales...), Markets: []string{}, SupportsMarketFilter: false,
			IdentityLevels: []string{"commit", "owner_repo"}, EvidenceSurfaces: []string{"code_occurrence", "immutable_commit", "repository_identity"}, Replayable: true, PublicOnly: true,
		}}
	}
	return []domain.MarketDiscoveryConnector{apple, repository("github-repository-catalog"), repository("sourcegraph-repository-catalog")}, nil
}

func federatedMarketDiscoveryConnectors() ([]domain.MarketDiscoveryConnector, error) {
	apple, err := itunes.NewMarketDiscoveryConnector(nil)
	if err != nil {
		return nil, err
	}
	play, err := playstore.NewMarketDiscoveryConnector(nil)
	if err != nil {
		return nil, err
	}
	fd, err := fdroid.NewMarketDiscoveryConnector(nil)
	if err != nil {
		return nil, err
	}
	locales := []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"}
	repository := func(id string) domain.MarketDiscoveryConnector {
		return usecase.CapabilityOnlyConnector{Capability: domain.DiscoveryConnectorCapabilities{
			ID: id, Version: "baseline-21c48e357180", Kind: domain.ConnectorRepositoryCatalog,
			QueryLocales: append([]string{}, locales...), Markets: []string{}, SupportsMarketFilter: false,
			IdentityLevels: []string{"commit", "owner_repo"}, EvidenceSurfaces: []string{"code_occurrence", "immutable_commit", "repository_identity"}, Replayable: true, PublicOnly: true,
		}}
	}
	gitee := forge.NewGiteeConnector(nil)
	codeberg := forge.NewCodebergConnector(nil)
	gitlab := forge.NewGitLabConnector(nil)

	return []domain.MarketDiscoveryConnector{
		apple,
		play,
		fd,
		repository("github-repository-catalog"),
		repository("sourcegraph-repository-catalog"),
		gitee,
		codeberg,
		gitlab,
	}, nil
}

func NewMarketDiscoveryMCPServer(
	funnel *usecase.SearchFunnel,
	recommendationExecutor domain.ObservableRecommendationExecutor,
	connectors []domain.MarketDiscoveryConnector,
	root string,
	cursorSecret []byte,
	limits runstore.RunStoreLimits,
	workers int,
	capacity int,
	reader io.Reader,
	writer io.Writer,
) (*mcp.MCPServer, func(context.Context) error, error) {
	registry, err := usecase.NewDiscoveryConnectorRegistry(connectors)
	if err != nil || len(registry.Capabilities()) < 3 {
		return nil, nil, domain.ErrInvalidRunInput
	}
	service, queue, store, err := newDurableComponents(recommendationExecutor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service.ConfigureMarketDiscovery(nil, registry)
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
		WithMarketDiscovery(connectors, service)
	return server, durableShutdown(queue, store), nil
}
