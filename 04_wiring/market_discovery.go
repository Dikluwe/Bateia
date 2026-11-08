/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/market_discovery.md
 * @layer L4
 * @updated 2026-10-28
 */

package main

import (
	"context"
	"io"
	"os"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/bitbucket"
	"bateia/03_infra/fdroid"
	"bateia/03_infra/forge"
	"bateia/03_infra/itunes"
	"bateia/03_infra/playstore"
	"bateia/03_infra/research"
	"bateia/03_infra/runcontainer"
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
	connectors := []domain.MarketDiscoveryConnector{apple, repository("github-repository-catalog"), repository("sourcegraph-repository-catalog"), research.NewZBMath(nil), research.NewJOSS(nil), forge.NewGiteeConnectorWithToken(os.Getenv("BATEIA_GITEE_TOKEN"), nil), forge.NewAtomGitConnector(os.Getenv("BATEIA_ATOMGIT_TOKEN"), nil)}
	return withConfiguredBitbucket(connectors)
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
	gitee := forge.NewGiteeConnectorWithToken(os.Getenv("BATEIA_GITEE_TOKEN"), nil)
	atomgit := forge.NewAtomGitConnector(os.Getenv("BATEIA_ATOMGIT_TOKEN"), nil)
	codeberg := forge.NewCodebergConnector(nil)
	gitlab := forge.NewGitLabConnector(nil)

	connectors := []domain.MarketDiscoveryConnector{
		apple,
		play,
		fd,
		repository("github-repository-catalog"),
		repository("sourcegraph-repository-catalog"),
		gitee,
		atomgit,
		codeberg,
		gitlab,
		research.NewZBMath(nil),
		research.NewJOSS(nil),
	}
	return withConfiguredBitbucket(connectors)
}

func withConfiguredBitbucket(connectors []domain.MarketDiscoveryConnector) ([]domain.MarketDiscoveryConnector, error) {
	raw := strings.TrimSpace(os.Getenv("BATEIA_BITBUCKET_WORKSPACES"))
	if raw == "" {
		return connectors, nil
	}
	connector, err := bitbucket.NewConnector(strings.Split(raw, ","), nil)
	if err != nil {
		return nil, err
	}
	return append(connectors, connector), nil
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
	registry, err := usecase.NewMarketDiscoveryRegistry(connectors)
	if err != nil {
		return nil, nil, domain.ErrInvalidRunInput
	}
	durable, err := runcontainer.New(recommendationExecutor, root, cursorSecret, limits, workers, capacity)
	if err != nil {
		return nil, nil, err
	}
	service := durable.Service()
	service.ConfigureMarketDiscovery(usecase.NewMarketDiscoveryService(registry), registry)
	if err := durable.Recover(context.Background()); err != nil {
		return nil, nil, err
	}
	if err := durable.Start(); err != nil {
		return nil, nil, err
	}
	server := mcp.NewMCPServer(funnel, reader, writer, recommendationExecutor).
		WithDurableRuns(service).
		WithMarketDiscovery(connectors, service)
	return server, durable.Shutdown, nil
}
