/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_registry_test.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"bateia/01_core/domain"
)

type registryConfigurationProbe struct {
	capability domain.DiscoveryConnectorCapabilities
}

func (c registryConfigurationProbe) Capabilities() domain.DiscoveryConnectorCapabilities {
	return c.capability
}
func (registryConfigurationProbe) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	panic("registry admission must not discover")
}
func configuredRegistryProbe(id string) registryConfigurationProbe {
	return registryConfigurationProbe{domain.DiscoveryConnectorCapabilities{ID: id, Version: "v1", Kind: domain.ConnectorProductCatalog, QueryLocales: []string{"pt-BR", "en-US"}, Markets: []string{}, IdentityLevels: []string{"catalog_id"}, EvidenceSurfaces: []string{}, Replayable: true, PublicOnly: true}}
}

func TestMarketRegistryAcceptsValidConfiguredSetsWithoutDiscovery(t *testing.T) {
	for _, count := range []int{1, 2, 3, 5, 11} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			connectors := []domain.MarketDiscoveryConnector{}
			for i := count - 1; i >= 0; i-- {
				connectors = append(connectors, configuredRegistryProbe(fmt.Sprintf("source-%02d", i)))
			}
			registry, err := NewMarketDiscoveryRegistry(connectors)
			if err != nil {
				t.Fatal(err)
			}
			caps := registry.Capabilities()
			if len(caps) != count || caps[0].ID != "source-00" || caps[0].QueryLocales[0] != "en-US" {
				t.Fatalf("registry not canonicalized: %+v", caps)
			}
			if _, ok := registry.Resolve("source-00"); !ok {
				t.Fatal("configured source not resolvable")
			}
			caps[0].QueryLocales[0] = "changed"
			if registry.Capabilities()[0].QueryLocales[0] != "en-US" {
				t.Fatal("projection mutates registry")
			}
		})
	}
}

func TestMarketRegistryRejectsInvalidConfigurationBeforeDiscovery(t *testing.T) {
	valid := configuredRegistryProbe("valid")
	bad := valid
	bad.capability.PublicOnly = false
	for name, connectors := range map[string][]domain.MarketDiscoveryConnector{
		"empty": nil, "nil": {nil}, "duplicate": {valid, valid}, "invalid": {bad},
	} {
		t.Run(name, func(t *testing.T) {
			registry, err := NewMarketDiscoveryRegistry(connectors)
			if !errors.Is(err, domain.ErrInvalidRunInput) || registry != nil {
				t.Fatalf("invalid configuration accepted: %v", err)
			}
		})
	}
	registry, err := NewDiscoveryConnectorRegistry(nil)
	if err != nil || len(registry.Capabilities()) != 0 {
		t.Fatal("generic empty registry contract changed")
	}
}
