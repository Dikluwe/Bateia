/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/run_capabilities.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"sort"
)

type marketRunCapability struct {
	registry *DiscoveryConnectorRegistry
	executor domain.ObservableMarketDiscoveryExecutor
}

func (c *marketRunCapability) validate(request domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartRequest, string, error) {
	normalized, digest, err := CanonicalMarketDiscoveryStartRequestHash(request, c.registry)
	if err != nil {
		return domain.MarketDiscoveryStartRequest{}, "", domain.ErrInvalidRunInput
	}
	for _, id := range normalized.ConnectorIDs {
		if _, ok := c.registry.Resolve(id); !ok {
			return domain.MarketDiscoveryStartRequest{}, "", domain.ErrInvalidRunInput
		}
	}
	if _, err := PlanMarketDiscovery(normalized, c.registry); err != nil {
		return domain.MarketDiscoveryStartRequest{}, "", domain.ErrInvalidRunInput
	}
	return normalized, digest, nil
}

type synthesisRunCapability struct{ catalogSHA256 string }

func (c synthesisRunCapability) enabled() bool { return c.catalogSHA256 == activeMechanismRulesSHA256 }

type mechanismRunCapability struct {
	extractor     domain.MechanismExtractor
	catalog       []byte
	catalogSHA256 string
	baseline      *domain.ConventionalBaseline
}

func (c *mechanismRunCapability) configure(extractor domain.MechanismExtractor, catalog []byte, digest string, baseline *domain.ConventionalBaseline) bool {
	if extractor == nil {
		return false
	}
	if _, err := decodeMechanismCatalog(catalog, digest); err != nil {
		return false
	}
	*c = mechanismRunCapability{extractor: extractor, catalog: append([]byte(nil), catalog...), catalogSHA256: digest, baseline: baseline}
	return true
}
func (c mechanismRunCapability) enabled() bool { return c.extractor != nil }
func (c mechanismRunCapability) extract(ctx context.Context, run domain.IdeaRunResult) ([]domain.ArchitecturalMechanism, []domain.MechanismGap, error) {
	mechanisms := []domain.ArchitecturalMechanism{}
	gaps := []domain.MechanismGap{}
	for _, bundle := range run.Bundles {
		extracted, err := c.extractor.ExtractMechanisms(ctx, domain.MechanismExtractionInput{Bundle: bundle, Assessments: run.Assessment.Assessments, ConventionalBaseline: c.baseline, RuleCatalogJSON: c.catalog, RuleCatalogSHA256: c.catalogSHA256})
		if err != nil {
			return nil, nil, err
		}
		mechanisms = append(mechanisms, extracted.Mechanisms...)
		gaps = append(gaps, extracted.Gaps...)
	}
	sort.Slice(mechanisms, func(i, j int) bool { return mechanisms[i].ID < mechanisms[j].ID })
	sortMechanismGaps(gaps)
	return mechanisms, gaps, nil
}

func (c marketRunCapability) execute(ctx context.Context, request domain.MarketDiscoveryStartRequest, reporter domain.MarketDiscoveryProgressReporter) (MarketDiscoveryExecution, error) {
	if detailed, ok := c.executor.(interface {
		DiscoverMarketExecutionObserved(context.Context, domain.MarketDiscoveryStartRequest, domain.MarketDiscoveryProgressReporter) (MarketDiscoveryExecution, error)
	}); ok {
		return detailed.DiscoverMarketExecutionObserved(ctx, request, reporter)
	}
	if c.executor == nil {
		return MarketDiscoveryExecution{}, domain.ErrInvalidRunInput
	}
	result, err := c.executor.DiscoverMarketObserved(ctx, request, reporter)
	if err != nil {
		return MarketDiscoveryExecution{}, err
	}
	if result.SchemaVersion != "1" || (result.Status != domain.MarketDiscoveryComplete && result.Status != domain.MarketDiscoveryPartial && result.Status != domain.MarketDiscoveryInconclusive) {
		return MarketDiscoveryExecution{}, marketExecutionFailure{"RESULT_INVALID", "market discovery result was invalid"}
	}
	artifacts, err := ProjectMarketDiscoveryArtifacts(result.Capabilities, result.Plan, result.Discoveries)
	if err != nil {
		return MarketDiscoveryExecution{}, err
	}
	// Candidates and gaps are conclusions delivered by the configured investigation.
	// The runtime maps the recorded observations but does not derive these conclusions again.
	artifacts.MarketCandidates = append([]domain.ProductCandidate{}, result.Candidates...)
	artifacts.MarketGaps = append([]domain.MarketDiscoveryGap{}, result.Gaps...)
	return MarketDiscoveryExecution{Result: result, Artifacts: artifacts}, nil
}

// NewConfiguredMechanismExtractor admits configuration without running an extraction.
func NewConfiguredMechanismExtractor(catalog []byte, digest string) (domain.MechanismExtractor, error) {
	if _, err := decodeMechanismCatalog(catalog, digest); err != nil {
		return nil, domain.ErrIntegrity
	}
	return NewDeterministicMechanismExtractor(), nil
}
