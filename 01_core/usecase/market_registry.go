/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_registry.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import "bateia/01_core/domain"

// NewMarketDiscoveryRegistry validates the configured sources before discovery
// is enabled. Source count is not a substitute for valid, unique capabilities.
func NewMarketDiscoveryRegistry(connectors []domain.MarketDiscoveryConnector) (*DiscoveryConnectorRegistry, error) {
	if len(connectors) == 0 {
		return nil, domain.ErrInvalidRunInput
	}
	return NewDiscoveryConnectorRegistry(connectors)
}
