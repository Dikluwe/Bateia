/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_plan.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"sort"
	"strings"
)

func ValidateInsightSiteOptions(request domain.InsightSiteRequest) (domain.InsightSiteRequest, error) {
	base, err := PlanInsightDiscovery(request.Discovery)
	if err != nil {
		return request, err
	}
	request.Discovery = base.Request
	if request.FollowLinks == nil {
		n := 2
		request.FollowLinks = &n
	} else {
		n := *request.FollowLinks
		request.FollowLinks = &n
	}
	if *request.FollowLinks < 0 || *request.FollowLinks > 4 {
		return request, domain.ErrInvalidRunInput
	}
	if request.Memory != nil {
		memory, err := ValidateInsightExplorationMemory(request.Memory)
		if err != nil {
			return request, err
		}
		if err := validateInsightMemoryAdvance(memory); err != nil {
			return request, err
		}
		next := memory
		next.Generation++
		if len(canonicalJSON(sealInsightMemory(next))) > 1<<20 {
			return request, domain.ErrInvalidRunInput
		}
		request.Memory = &memory
	}
	if len(request.Sources) == 0 || len(request.Sources) > 12 {
		return request, domain.ErrInvalidRunInput
	}
	if request.MaxCalls == 0 {
		request.MaxCalls = 20
	}
	if request.MaxCalls < 1 || request.MaxCalls > 60 {
		return request, domain.ErrInvalidRunInput
	}
	if request.Locale == "" {
		request.Locale = "en-US"
	}
	if !validBCP47Locale(request.Locale) {
		return request, domain.ErrInvalidRunInput
	}
	return request, nil
}

func PlanInsightSites(request domain.InsightSiteRequest, bindings []domain.InsightSiteBinding) (domain.InsightSitePlan, error) {
	request, err := ValidateInsightSiteOptions(request)
	if err != nil {
		return domain.InsightSitePlan{}, err
	}
	if len(bindings) == 0 || len(bindings) > 12 {
		return domain.InsightSitePlan{}, domain.ErrInvalidRunInput
	}
	base, _ := PlanInsightDiscovery(request.Discovery)
	caps := map[string]domain.DiscoveryConnectorCapabilities{}
	for _, binding := range bindings {
		if binding.Connector == nil {
			return domain.InsightSitePlan{}, domain.ErrInvalidRunInput
		}
		capability, err := canonicalCapabilities(binding.Connector.Capabilities())
		if err != nil {
			return domain.InsightSitePlan{}, err
		}
		if _, exists := caps[capability.ID]; exists || (binding.FollowReader != nil && capability.ID != domain.ConnectorGitHubCatalog) || len(canonicalJSON(capability)) > 4096 || (binding.Reader != nil && capability.ID != domain.ConnectorGitHubCatalog) || (binding.RepositoryReader != nil && (binding.Reader != nil || (capability.ID != domain.ConnectorGiteeCatalog && capability.ID != domain.ConnectorAtomGitCatalog && capability.ID != domain.ConnectorCodebergCatalog && capability.ID != domain.ConnectorGitLabCatalog))) {
			return domain.InsightSitePlan{}, domain.ErrInvalidRunInput
		}
		caps[capability.ID] = capability
	}
	aliases := map[string]string{"github": domain.ConnectorGitHubCatalog, "sourcegraph": domain.ConnectorSourcegraph, "gitee": domain.ConnectorGiteeCatalog, "atomgit": domain.ConnectorAtomGitCatalog, "codeberg": domain.ConnectorCodebergCatalog, "gitlab": domain.ConnectorGitLabCatalog, "bitbucket": domain.ConnectorBitbucketCatalog, "joss": domain.ConnectorJOSS, "zbmath": domain.ConnectorZBMath, "itunes": domain.ConnectorAppleItunes, "playstore": domain.ConnectorGooglePlay, "fdroid": domain.ConnectorFDroid}
	selected := map[string]bool{}
	if len(request.Sources) == 1 && request.Sources[0] == "all" {
		for id := range caps {
			selected[id] = true
		}
	} else {
		for _, id := range request.Sources {
			if strings.TrimSpace(id) != id {
				return domain.InsightSitePlan{}, domain.ErrInvalidRunInput
			}
			if alias, ok := aliases[id]; ok {
				id = alias
			}
			if _, ok := caps[id]; !ok || selected[id] {
				return domain.InsightSitePlan{}, domain.ErrInvalidRunInput
			}
			selected[id] = true
		}
	}
	request.Sources = make([]string, 0, len(selected))
	for id := range selected {
		request.Sources = append(request.Sources, id)
	}
	sort.Strings(request.Sources)
	plan := domain.InsightSitePlan{SchemaVersion: "insight-site-plan-v1", SchedulingPolicy: insightSiteSchedulingPolicy, Request: request, Discovery: base, Sources: []domain.DiscoveryConnectorCapabilities{}}
	for _, id := range request.Sources {
		plan.Sources = append(plan.Sources, caps[id])
	}
	plan.ID = "site_plan_" + insightHash(plan)
	return plan, nil
}
