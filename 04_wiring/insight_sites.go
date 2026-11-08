/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_sites.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/03_infra/forge"
	"bateia/03_infra/github"
	"bateia/03_infra/sourcegraph"
	"os"
)

func insightDiscoverySites(client *github.Client) ([]domain.InsightSiteBinding, error) {
	connectors, err := federatedMarketDiscoveryConnectors()
	if err != nil {
		return nil, err
	}
	bindings := []domain.InsightSiteBinding{}
	for _, connector := range connectors {
		binding := domain.InsightSiteBinding{Connector: connector}
		switch connector.Capabilities().ID {
		case domain.ConnectorGitHubCatalog:
			binding.Connector = usecase.NewInsightGitHubConnector(client)
			binding.Reader = client
			binding.FollowReader = github.NewInsightLinkedReader(client)
		case domain.ConnectorGiteeCatalog:
			binding.RepositoryReader = forge.NewGiteeInsightReader(os.Getenv("BATEIA_GITEE_TOKEN"), nil)
		case domain.ConnectorAtomGitCatalog:
			binding.RepositoryReader = forge.NewAtomGitInsightReader(os.Getenv("BATEIA_ATOMGIT_TOKEN"), nil)
		case domain.ConnectorCodebergCatalog:
			binding.RepositoryReader = forge.NewCodebergInsightReader(nil)
		case domain.ConnectorGitLabCatalog:
			binding.RepositoryReader = forge.NewGitLabInsightReader(nil)
		case domain.ConnectorSourcegraph:
			binding.Connector = usecase.NewInsightSourcegraphConnector(sourcegraph.NewClient(os.Getenv("SOURCEGRAPH_TOKEN")))
		}
		bindings = append(bindings, binding)
	}
	return bindings, nil
}
