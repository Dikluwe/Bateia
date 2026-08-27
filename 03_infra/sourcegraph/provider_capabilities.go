/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/provider_capabilities.md
 * @layer L3
 * @updated 2026-08-26
 */
package sourcegraph

import "bateia/01_core/domain"

func (c *Client) Capabilities() domain.ProviderCapabilities {
	return domain.ProviderCapabilities{Provider: "sourcegraph", Roles: []domain.DiscoveryProviderRole{domain.RoleCodeProof}, SearchUnits: []string{"code_occurrence"}, IdentityLevels: []string{"repository", "revision", "path", "line"}}
}

var _ domain.ProviderCapabilityReporter = (*Client)(nil)
