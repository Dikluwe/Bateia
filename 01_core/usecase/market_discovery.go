/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_discovery.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"bateia/01_core/domain"
)

// validISO3166Alpha2 valida um código de país conforme ISO 3166-1 alpha-2: exatamente 2 letras maiúsculas.
// Validação estrutural pura — zero I/O, zero rede. Não implica que o país tenha storefront ativo.
func validISO3166Alpha2(code string) bool {
	if len(code) != 2 {
		return false
	}
	return code[0] >= 'A' && code[0] <= 'Z' && code[1] >= 'A' && code[1] <= 'Z'
}

// validBCP47Locale valida uma tag de idioma/locale conforme BCP-47 (RFC 5646) no subconjunto
// "xx" (2 letras minúsculas de idioma) ou "xx-YY" (2 letras de idioma + 2 letras maiúsculas de região).
// Exemplos válidos: "en-US", "it-IT", "kk-KZ", "mn-MN". Inválidos: "en_US", "BR", "brBR".
// Validação estrutural pura — zero I/O, zero rede.
func validBCP47Locale(tag string) bool {
	if len(tag) == 2 {
		return tag[0] >= 'a' && tag[0] <= 'z' && tag[1] >= 'a' && tag[1] <= 'z'
	}
	if len(tag) == 5 {
		return tag[2] == '-' &&
			tag[0] >= 'a' && tag[0] <= 'z' && tag[1] >= 'a' && tag[1] <= 'z' &&
			tag[3] >= 'A' && tag[3] <= 'Z' && tag[4] >= 'A' && tag[4] <= 'Z'
	}
	return false
}

type DiscoveryConnectorRegistry struct {
	connectors map[string]domain.MarketDiscoveryConnector
	ordered    []domain.DiscoveryConnectorCapabilities
}

func NewDiscoveryConnectorRegistry(connectors []domain.MarketDiscoveryConnector) (*DiscoveryConnectorRegistry, error) {
	registry := &DiscoveryConnectorRegistry{connectors: map[string]domain.MarketDiscoveryConnector{}, ordered: []domain.DiscoveryConnectorCapabilities{}}
	for _, connector := range connectors {
		if connector == nil {
			return nil, domain.ErrInvalidRunInput
		}
		capabilities, err := canonicalCapabilities(connector.Capabilities())
		if err != nil {
			return nil, err
		}
		if _, exists := registry.connectors[capabilities.ID]; exists {
			return nil, domain.ErrInvalidRunInput
		}
		registry.connectors[capabilities.ID] = connector
		registry.ordered = append(registry.ordered, capabilities)
	}
	sort.Slice(registry.ordered, func(i, j int) bool { return registry.ordered[i].ID < registry.ordered[j].ID })
	return registry, nil
}

func (r *DiscoveryConnectorRegistry) Capabilities() []domain.DiscoveryConnectorCapabilities {
	if r == nil {
		return []domain.DiscoveryConnectorCapabilities{}
	}
	out := make([]domain.DiscoveryConnectorCapabilities, len(r.ordered))
	for i := range r.ordered {
		out[i] = cloneCapabilities(r.ordered[i])
	}
	return out
}

func (r *DiscoveryConnectorRegistry) Resolve(id string) (domain.MarketDiscoveryConnector, bool) {
	if r == nil {
		return nil, false
	}
	connector, ok := r.connectors[id]
	return connector, ok
}

func canonicalCapabilities(value domain.DiscoveryConnectorCapabilities) (domain.DiscoveryConnectorCapabilities, error) {
	value.ID = strings.TrimSpace(value.ID)
	value.Version = strings.TrimSpace(value.Version)
	if value.ID == "" || value.Version == "" || !value.PublicOnly || !value.Replayable || len(value.IdentityLevels) == 0 {
		return domain.DiscoveryConnectorCapabilities{}, domain.ErrInvalidRunInput
	}
	if value.Kind != domain.ConnectorRepositoryCatalog && value.Kind != domain.ConnectorApplicationCatalog && value.Kind != domain.ConnectorProductCatalog {
		return domain.DiscoveryConnectorCapabilities{}, domain.ErrInvalidRunInput
	}
	var err error
	if value.QueryLocales, err = canonicalValidatedSetBCP47(value.QueryLocales); err != nil || len(value.QueryLocales) == 0 {
		return domain.DiscoveryConnectorCapabilities{}, domain.ErrInvalidRunInput
	}
	if value.Markets, err = canonicalValidatedSetISO3166(value.Markets); err != nil {
		return domain.DiscoveryConnectorCapabilities{}, domain.ErrInvalidRunInput
	}
	if value.SupportsMarketFilter && len(value.Markets) == 0 {
		return domain.DiscoveryConnectorCapabilities{}, domain.ErrInvalidRunInput
	}
	value.IdentityLevels, err = canonicalNonEmptySet(value.IdentityLevels)
	if err != nil {
		return domain.DiscoveryConnectorCapabilities{}, err
	}
	value.EvidenceSurfaces, err = canonicalNonEmptySet(value.EvidenceSurfaces)
	if err != nil {
		return domain.DiscoveryConnectorCapabilities{}, err
	}
	return value, nil
}

// canonicalValidatedSetISO3166 deduplica e valida um conjunto de códigos de país ISO 3166-1 alpha-2.
func canonicalValidatedSetISO3166(values []string) ([]string, error) {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if !validISO3166Alpha2(value) || seen[value] {
			return nil, domain.ErrInvalidRunInput
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out, nil
}

// canonicalValidatedSetBCP47 deduplica e valida um conjunto de tags de locale BCP-47.
func canonicalValidatedSetBCP47(values []string) ([]string, error) {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if !validBCP47Locale(value) || seen[value] {
			return nil, domain.ErrInvalidRunInput
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out, nil
}

func canonicalNonEmptySet(values []string) ([]string, error) {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			return nil, domain.ErrInvalidRunInput
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	if out == nil {
		out = []string{}
	}
	return out, nil
}

func cloneCapabilities(value domain.DiscoveryConnectorCapabilities) domain.DiscoveryConnectorCapabilities {
	value.QueryLocales = append([]string{}, value.QueryLocales...)
	value.Markets = append([]string{}, value.Markets...)
	value.IdentityLevels = append([]string{}, value.IdentityLevels...)
	value.EvidenceSurfaces = append([]string{}, value.EvidenceSurfaces...)
	return value
}

func ValidateMarketDiscoveryStartRequest(request domain.MarketDiscoveryStartRequest, registry *DiscoveryConnectorRegistry) (domain.MarketDiscoveryStartRequest, error) {
	request, err := normalizeMarketDiscoveryStartRequest(request)
	if err != nil || registry == nil {
		return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
	}
	for _, id := range request.ConnectorIDs {
		if _, ok := registry.Resolve(id); !ok {
			return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
		}
	}
	return request, nil
}

func normalizeMarketDiscoveryStartRequest(request domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartRequest, error) {
	request.Idea = normalizeIdea(request.Idea)
	if request.Idea.Problem == "" || len(request.Idea.Problem) > 1000 || len(request.Idea.Capabilities) < 1 || len(request.Idea.Capabilities) > 20 || len(request.QueryProposals) < 1 || len(request.QueryProposals) > 20 || len(request.Markets) < 1 || len(request.Markets) > 50 || len(request.ConnectorIDs) < 1 || len(request.ConnectorIDs) > 10 {
		return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
	}
	if !validIdeaSlices(request.Idea) || request.Budget.MaxQueries < 1 || request.Budget.MaxQueries > 20 || request.Budget.MaxResultsPerQuery < 1 || request.Budget.MaxResultsPerQuery > 10 || request.Budget.MaxTotalCandidates < 1 || request.Budget.MaxTotalCandidates > 100 {
		return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
	}
	var err error
	request.Markets, err = canonicalValidatedSetISO3166(request.Markets)
	if err != nil {
		return domain.MarketDiscoveryStartRequest{}, err
	}
	request.ConnectorIDs, err = canonicalNonEmptySet(request.ConnectorIDs)
	if err != nil || len(request.ConnectorIDs) == 0 {
		return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
	}
	proposals := make([]domain.MarketQueryProposal, len(request.QueryProposals))
	seen := map[string]bool{}
	for i, proposal := range request.QueryProposals {
		proposal.Query = strings.TrimSpace(proposal.Query)
		proposal.SourceText = strings.TrimSpace(proposal.SourceText)
		proposal.Rationale = strings.TrimSpace(proposal.Rationale)
		if !validBCP47Locale(proposal.Locale) || proposal.Derivation != domain.QueryExternalAgent || !safeSemanticText(proposal.Query, 256) || !safeSemanticText(proposal.SourceText, 512) || !safeSemanticText(proposal.Rationale, 512) {
			return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
		}
		proposal.DerivationStatus = domain.QueryDerivationUnverified
		proposal.CatalogVersion, proposal.RuleID = "", ""
		proposal.InputHash = hashText(proposal.SourceText)
		proposal.ID = hashText(proposal.Locale + "\x00" + proposal.Query + "\x00" + proposal.SourceText + "\x00" + string(proposal.Derivation) + "\x00" + proposal.Rationale)
		if seen[proposal.ID] {
			return domain.MarketDiscoveryStartRequest{}, domain.ErrInvalidRunInput
		}
		seen[proposal.ID] = true
		proposals[i] = proposal
	}
	sort.Slice(proposals, func(i, j int) bool { return proposals[i].ID < proposals[j].ID })
	request.QueryProposals = proposals
	return request, nil
}

func validIdeaSlices(idea domain.IdeaSpec) bool {
	for _, values := range [][]string{idea.Users, idea.Capabilities, idea.Constraints, idea.NonGoals, idea.PreferredStack} {
		if len(values) > 20 {
			return false
		}
		for _, value := range values {
			if !safeSemanticText(value, 200) {
				return false
			}
		}
	}
	return true
}

func safeSemanticText(value string, maximum int) bool {
	if value == "" || len([]rune(value)) > maximum {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) || r == '\x7f' {
			return false
		}
	}
	lower := strings.ToLower(value)
	for _, forbidden := range []string{"http://", "https://", "callback=", "authorization:", "host:", "../", "\\.."} {
		if strings.Contains(lower, forbidden) {
			return false
		}
	}
	return true
}

func CanonicalMarketDiscoveryStartRequestHash(request domain.MarketDiscoveryStartRequest, registry *DiscoveryConnectorRegistry) (domain.MarketDiscoveryStartRequest, string, error) {
	normalized, err := ValidateMarketDiscoveryStartRequest(request, registry)
	if err != nil {
		return domain.MarketDiscoveryStartRequest{}, "", err
	}
	versions := []string{}
	if registry != nil {
		for _, capability := range registry.Capabilities() {
			versions = append(versions, capability.ID+"@"+capability.Version)
		}
		sort.Strings(versions)
	}
	encoded, err := json.Marshal(struct {
		Request  domain.MarketDiscoveryStartRequest `json:"request"`
		Versions []string                           `json:"connector_versions"`
	}{normalized, versions})
	if err != nil {
		return domain.MarketDiscoveryStartRequest{}, "", domain.ErrInvalidRunInput
	}
	return normalized, hashBytes(encoded), nil
}

func PlanMarketDiscovery(request domain.MarketDiscoveryStartRequest, registry *DiscoveryConnectorRegistry) ([]domain.MarketDiscoveryAttempt, error) {
	normalized, err := ValidateMarketDiscoveryStartRequest(request, registry)
	if err != nil {
		return nil, err
	}
	byID := map[string]domain.DiscoveryConnectorCapabilities{}
	for _, capability := range registry.Capabilities() {
		canonical, canonicalErr := canonicalCapabilities(capability)
		if canonicalErr != nil {
			return nil, canonicalErr
		}
		byID[canonical.ID] = canonical
	}
	attempts := []domain.MarketDiscoveryAttempt{}
	for _, proposal := range normalized.QueryProposals {
		for _, market := range normalized.Markets {
			for _, connectorID := range normalized.ConnectorIDs {
				capability, ok := byID[connectorID]
				if !ok {
					return nil, domain.ErrInvalidRunInput
				}
				status, reason, filtered := domain.AttemptPlanned, "PLANNED", capability.SupportsMarketFilter
				if !marketContains(capability.QueryLocales, proposal.Locale) || capability.SupportsMarketFilter && !marketContains(capability.Markets, market) {
					status, reason = domain.AttemptUnsupportedCapability, "UNSUPPORTED_CAPABILITY"
				} else if !capability.SupportsMarketFilter {
					status, reason, filtered = domain.AttemptMarketNotFiltered, "MARKET_NOT_FILTERED", false
				}
				queryHash := hashText(proposal.Query)
				requestHash := hashText(proposal.ID + "\x00" + market + "\x00" + connectorID + "\x00" + capability.Version)
				attempts = append(attempts, domain.MarketDiscoveryAttempt{
					AttemptID: hashText("attempt\x00" + requestHash), QueryID: proposal.ID, QueryLocale: proposal.Locale, Market: market,
					ConnectorID: connectorID, ConnectorVersion: capability.Version, MarketFilterApplied: filtered,
					QueryHash: queryHash, EffectiveQueryHash: queryHash, RequestHash: requestHash, PlannedCost: 0, Status: status, Reason: reason,
				})
			}
		}
	}
	sortAttempts(attempts)
	return attempts, nil
}

func AllocateMarketDiscoveryBudget(attempts []domain.MarketDiscoveryAttempt, budget domain.MarketDiscoveryBudget) []domain.MarketDiscoveryAttempt {
	out := append([]domain.MarketDiscoveryAttempt{}, attempts...)
	sortAttempts(out)
	groups := map[string][]int{}
	connectors := []string{}
	for i := range out {
		if out[i].Status != domain.AttemptPlanned && out[i].Status != domain.AttemptMarketNotFiltered {
			continue
		}
		if _, ok := groups[out[i].ConnectorID]; !ok {
			connectors = append(connectors, out[i].ConnectorID)
		}
		groups[out[i].ConnectorID] = append(groups[out[i].ConnectorID], i)
		if out[i].Status == domain.AttemptPlanned {
			out[i].Status, out[i].Reason = domain.AttemptBudgetExhausted, "BUDGET_EXHAUSTED"
		}
	}
	sort.Strings(connectors)
	used := 0
	for round := 0; used < budget.MaxQueries; round++ {
		progress := false
		for _, connector := range connectors {
			indices := groups[connector]
			if round >= len(indices) || used >= budget.MaxQueries {
				continue
			}
			index := indices[round]
			if out[index].Status == domain.AttemptBudgetExhausted {
				out[index].Status, out[index].Reason = domain.AttemptPlanned, "PLANNED"
			}
			out[index].PlannedCost = 1
			used++
			progress = true
		}
		if !progress {
			break
		}
	}
	sortAttempts(out)
	return out
}

func sortAttempts(values []domain.MarketDiscoveryAttempt) {
	sort.Slice(values, func(i, j int) bool {
		if values[i].QueryID != values[j].QueryID {
			return values[i].QueryID < values[j].QueryID
		}
		if values[i].Market != values[j].Market {
			return values[i].Market < values[j].Market
		}
		return values[i].ConnectorID < values[j].ConnectorID
	})
}

func ResolveProductIdentity(left, right domain.ProductCandidate) (domain.IdentityResolution, []domain.IdentityWitness) {
	if left.Identity.ConnectorID == right.Identity.ConnectorID {
		if left.Identity.ExternalID != "" && left.Identity.ExternalID == right.Identity.ExternalID {
			return domain.IdentitySame, []domain.IdentityWitness{{Kind: "provider_product_id", Value: left.Identity.ExternalID, ConnectorIDs: []string{left.Identity.ConnectorID}, Bilateral: true}}
		}
		return domain.IdentityDistinct, []domain.IdentityWitness{}
	}
	strong := map[string]bool{"official_shared_identifier": true, "reciprocal_official_url": true, "shared_package_id": true, "verified_official_repository": true, "bilateral_vendor_link": true}
	for _, leftLink := range left.RepositoryLinks {
		if leftLink.Status != domain.RepositoryLinkVerified || leftLink.OwnerRepo == "" || leftLink.Commit == "" {
			continue
		}
		for _, rightLink := range right.RepositoryLinks {
			if rightLink.Status == domain.RepositoryLinkVerified && rightLink.OwnerRepo == leftLink.OwnerRepo && rightLink.Commit == leftLink.Commit {
				witness := domain.IdentityWitness{Kind: "verified_official_repository", Value: leftLink.OwnerRepo + "@" + leftLink.Commit, ConnectorIDs: canonicalStrings([]string{left.Identity.ConnectorID, right.Identity.ConnectorID}), Bilateral: true}
				return domain.IdentitySame, []domain.IdentityWitness{witness}
			}
		}
	}
	witnesses := append(append([]domain.IdentityWitness{}, left.Witnesses...), right.Witnesses...)
	for _, witness := range witnesses {
		if witness.Bilateral && strong[witness.Kind] && witness.Value != "" && marketContains(witness.ConnectorIDs, left.Identity.ConnectorID) && marketContains(witness.ConnectorIDs, right.Identity.ConnectorID) {
			return domain.IdentitySame, []domain.IdentityWitness{witness}
		}
	}
	if strings.EqualFold(strings.TrimSpace(left.Name), strings.TrimSpace(right.Name)) {
		return domain.IdentityPossible, []domain.IdentityWitness{}
	}
	return domain.IdentityUnknown, []domain.IdentityWitness{}
}

func ConsolidateProductCandidates(values []domain.ProductCandidate) ([]domain.ProductCandidate, error) {
	byIdentity := map[string]domain.ProductCandidate{}
	for _, value := range values {
		if value.Identity.ConnectorID == "" || value.Identity.ExternalID == "" {
			return nil, domain.ErrInvalidRunInput
		}
		value = canonicalProduct(value)
		key := value.Identity.ConnectorID + "\x00" + value.Identity.ExternalID
		if prior, ok := byIdentity[key]; ok {
			prior.ObservedMarkets = canonicalStrings(append(prior.ObservedMarkets, value.ObservedMarkets...))
			prior.Origins = canonicalOrigins(append(prior.Origins, value.Origins...))
			prior.RepositoryLinks = canonicalRepoLinks(append(prior.RepositoryLinks, value.RepositoryLinks...))
			prior.Resolution = domain.IdentitySame
			byIdentity[key] = prior
		} else {
			if value.Resolution == "" {
				value.Resolution = domain.IdentityDistinct
			}
			byIdentity[key] = value
		}
	}
	out := make([]domain.ProductCandidate, 0, len(byIdentity))
	for _, value := range byIdentity {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Identity.ConnectorID != out[j].Identity.ConnectorID {
			return out[i].Identity.ConnectorID < out[j].Identity.ConnectorID
		}
		return out[i].Identity.ExternalID < out[j].Identity.ExternalID
	})
	return out, nil
}

func canonicalProduct(value domain.ProductCandidate) domain.ProductCandidate {
	value.ListingLocales = canonicalStrings(value.ListingLocales)
	value.ObservedMarkets = canonicalStrings(value.ObservedMarkets)
	value.PublisherCountries = canonicalStrings(value.PublisherCountries)
	value.RepositoryLinks = canonicalRepoLinks(value.RepositoryLinks)
	value.Origins = canonicalOrigins(value.Origins)
	if value.Witnesses == nil {
		value.Witnesses = []domain.IdentityWitness{}
	}
	return value
}

func canonicalStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func canonicalOrigins(values []domain.MarketCandidateOrigin) []domain.MarketCandidateOrigin {
	seen := map[string]bool{}
	out := []domain.MarketCandidateOrigin{}
	for _, value := range values {
		key, _ := json.Marshal(value)
		if !seen[string(key)] {
			seen[string(key)] = true
			out = append(out, value)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ConnectorID != out[j].ConnectorID {
			return out[i].ConnectorID < out[j].ConnectorID
		}
		if out[i].QueryID != out[j].QueryID {
			return out[i].QueryID < out[j].QueryID
		}
		return out[i].ObservedMarket < out[j].ObservedMarket
	})
	return out
}

func canonicalRepoLinks(values []domain.ProductRepositoryLink) []domain.ProductRepositoryLink {
	seen := map[string]bool{}
	out := []domain.ProductRepositoryLink{}
	for _, value := range values {
		key := value.OwnerRepo + "\x00" + value.Commit + "\x00" + string(value.Status) + "\x00" + value.Source
		if !seen[key] {
			seen[key] = true
			out = append(out, value)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].OwnerRepo != out[j].OwnerRepo {
			return out[i].OwnerRepo < out[j].OwnerRepo
		}
		return out[i].Commit < out[j].Commit
	})
	return out
}

func ProjectMarketDiscoveryArtifacts(capabilities []domain.DiscoveryConnectorCapabilities, plan []domain.MarketDiscoveryAttempt, discoveries []domain.ConnectorDiscovery) (domain.MarketDiscoveryArtifacts, error) {
	caps := append([]domain.DiscoveryConnectorCapabilities{}, capabilities...)
	sort.Slice(caps, func(i, j int) bool { return caps[i].ID < caps[j].ID })
	snapshots := []domain.MarketSnapshotEntry{}
	gaps := []domain.MarketGap{}
	candidates := []domain.ProductCandidate{}
	byAttempt := map[string]domain.ConnectorDiscovery{}
	for _, discovery := range discoveries {
		if discovery.AttemptID == "" {
			return domain.MarketDiscoveryArtifacts{}, domain.ErrInvalidRunInput
		}
		if _, exists := byAttempt[discovery.AttemptID]; exists {
			return domain.MarketDiscoveryArtifacts{}, domain.ErrInvalidRunInput
		}
		byAttempt[discovery.AttemptID] = discovery
	}
	for _, attempt := range plan {
		discovery, observed := byAttempt[attempt.AttemptID]
		status := attempt.Status
		if observed {
			status = discovery.Status
		}
		entryAttempt := attempt
		entryAttempt.Status = status
		entryAttempt.Reason = string(status)
		entry := domain.MarketSnapshotEntry{MarketDiscoveryAttempt: entryAttempt, ResultCount: len(discovery.Candidates), AcceptedCount: len(discovery.Candidates), ErrorClass: discovery.ErrorClass, ResponseHash: discovery.ResponseHash, ObservedAt: discovery.ObservedAt}
		snapshots = append(snapshots, entry)
		if status != domain.AttemptSucceeded {
			gaps = append(gaps, domain.MarketGap{AttemptID: attempt.AttemptID, QueryLocale: attempt.QueryLocale, Market: attempt.Market, ConnectorID: attempt.ConnectorID, Status: status, ReasonCode: string(status), ResponseHash: discovery.ResponseHash})
		}
		if attempt.Status == domain.AttemptMarketNotFiltered && status != domain.AttemptMarketNotFiltered {
			gaps = append(gaps, domain.MarketGap{AttemptID: attempt.AttemptID, QueryLocale: attempt.QueryLocale, Market: attempt.Market, ConnectorID: attempt.ConnectorID, Status: domain.AttemptMarketNotFiltered, ReasonCode: "MARKET_NOT_FILTERED", ResponseHash: discovery.ResponseHash})
		}
		observedMarket := ""
		if attempt.MarketFilterApplied {
			observedMarket = attempt.Market
		}
		for _, candidate := range discovery.Candidates {
			candidate.Origins = append(candidate.Origins, domain.MarketCandidateOrigin{ConnectorID: attempt.ConnectorID, ConnectorVersion: attempt.ConnectorVersion, QueryID: attempt.QueryID, QueryLocale: attempt.QueryLocale, ObservedMarket: observedMarket, ResponseHash: discovery.ResponseHash})
			if status == domain.AttemptSucceeded && attempt.MarketFilterApplied {
				candidate.ObservedMarkets = append(candidate.ObservedMarkets, attempt.Market)
			}
			candidates = append(candidates, candidate)
		}
	}
	sort.Slice(snapshots, func(i, j int) bool { return snapshots[i].AttemptID < snapshots[j].AttemptID })
	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].AttemptID != gaps[j].AttemptID {
			return gaps[i].AttemptID < gaps[j].AttemptID
		}
		return gaps[i].Status < gaps[j].Status
	})
	consolidated, err := ConsolidateProductCandidates(candidates)
	if err != nil {
		return domain.MarketDiscoveryArtifacts{}, err
	}
	return domain.MarketDiscoveryArtifacts{ConnectorCapabilities: caps, MarketPlan: append([]domain.MarketDiscoveryAttempt{}, plan...), MarketSnapshot: snapshots, MarketCandidates: consolidated, MarketGaps: gaps}, nil
}

func hashText(value string) string { return hashBytes([]byte(value)) }
func hashBytes(value []byte) string {
	digest := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(digest[:])
}
func marketContains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

type CapabilityOnlyConnector struct {
	Capability domain.DiscoveryConnectorCapabilities
}

func (c CapabilityOnlyConnector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return cloneCapabilities(c.Capability)
}
func (c CapabilityOnlyConnector) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	return domain.ConnectorDiscovery{Status: domain.AttemptMarketNotFiltered, Candidates: []domain.ProductCandidate{}, ErrorClass: "MARKET_NOT_FILTERED"}, nil
}

func NewRecommendationRunServiceWithMarketDiscovery(store domain.RecommendationRunStore, queue domain.RecommendationJobQueue, ids domain.RecommendationRunIDGenerator, executor domain.ObservableRecommendationExecutor, marketExecutor domain.ObservableMarketDiscoveryExecutor, registry *DiscoveryConnectorRegistry) *RecommendationRunService {
	service := NewRecommendationRunService(store, queue, ids, executor)
	return service.ConfigureMarketDiscovery(marketExecutor, registry)
}

func (s *RecommendationRunService) ConfigureMarketDiscovery(marketExecutor domain.ObservableMarketDiscoveryExecutor, registry *DiscoveryConnectorRegistry) *RecommendationRunService {
	if s != nil {
		s.market.registry = registry
		if marketExecutor == nil {
			marketExecutor = NewMarketDiscoveryService(registry)
		}
		s.market.executor = marketExecutor
	}
	return s
}

func (s *RecommendationRunService) StartMarketDiscovery(ctx context.Context, key string, request domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartReceipt, error) {
	if s == nil || s.store == nil || s.queue == nil || s.ids == nil || s.market.registry == nil || !validIdempotencyKey(key) {
		return domain.MarketDiscoveryStartReceipt{}, domain.ErrInvalidRunInput
	}
	normalized, requestHash, err := s.market.validate(request)
	if err != nil {
		return domain.MarketDiscoveryStartReceipt{}, err
	}

	keyDigest := digestString(key)
	if record, ok, resolveErr := s.store.ResolveIdempotency(ctx, keyDigest); resolveErr != nil {
		return domain.MarketDiscoveryStartReceipt{}, resolveErr
	} else if ok {
		if record.RequestHash != requestHash {
			return domain.MarketDiscoveryStartReceipt{}, domain.ErrIdempotencyConflict
		}
		run, getErr := s.store.Get(ctx, record.RunID)
		if getErr != nil {
			return domain.MarketDiscoveryStartReceipt{}, getErr
		}
		return marketReceiptFor(run), nil
	}
	reservation, err := s.queue.Reserve(ctx)
	if err != nil {
		return domain.MarketDiscoveryStartReceipt{}, err
	}
	consumed := false
	defer func() {
		if !consumed {
			reservation.Release()
		}
	}()
	id, err := s.ids.NewRunID()
	if err != nil || !runIDPattern.MatchString(string(id)) {
		return domain.MarketDiscoveryStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	run := domain.RecommendationRun{SchemaVersion: "1", RunID: id, RunKind: domain.RunKindMarketDiscovery, MarketRequest: &normalized, RequestHash: requestHash, Lifecycle: domain.RunQueued, Version: 1, Stage: domain.RunStageQueued, Progress: domain.RecommendationRunProgress{Planned: 4}, Artifacts: []domain.RunArtifactIndex{}}
	record := domain.IdempotencyRecord{KeyDigest: keyDigest, RequestHash: requestHash, RunID: id}
	actual, created, err := s.store.CreateOrResolve(ctx, run, record)
	if err != nil {
		return domain.MarketDiscoveryStartReceipt{}, err
	}
	if !created {
		if actual.RequestHash != requestHash {
			return domain.MarketDiscoveryStartReceipt{}, domain.ErrIdempotencyConflict
		}
		return marketReceiptFor(actual), nil
	}
	if err := reservation.Commit(actual.RunID); err != nil {
		return domain.MarketDiscoveryStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	consumed = true
	return marketReceiptFor(actual), nil
}

func marketReceiptFor(run domain.RecommendationRun) domain.MarketDiscoveryStartReceipt {
	return domain.MarketDiscoveryStartReceipt{SchemaVersion: "1", RunID: run.RunID, RunKind: domain.RunKindMarketDiscovery, Lifecycle: run.Lifecycle, Version: run.Version}
}

func (s *RecommendationRunService) executeMarketDiscovery(ctx context.Context, run domain.RecommendationRun) error {
	if run.MarketRequest == nil || s.market.registry == nil {
		return domain.ErrInvalidRunInput
	}
	var err error
	if run.Lifecycle == domain.RunQueued {
		run, err = s.store.Transition(ctx, domain.RunTransition{RunID: run.RunID, RunKind: domain.RunKindMarketDiscovery, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageMarketPlan, Progress: domain.RecommendationRunProgress{Planned: 4}})
		if err != nil {
			return err
		}
	}
	if run.Lifecycle != domain.RunRunning {
		return nil
	}
	run, err = s.ensureRequestArtifact(ctx, run.RunID)
	if err != nil {
		return err
	}
	selectedCaps := selectedMarketCapabilities(s.market.registry, *run.MarketRequest)
	run, err = s.publishDocument(ctx, run, domain.ArtifactConnectorCapabilities, interfacesFromCapabilities(selectedCaps))
	if err != nil {
		return err
	}
	executionCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	reporter := &marketRunProgressReporter{ctx: executionCtx, service: s, run: run, cancel: cancel}
	execution, executionErr := s.market.execute(executionCtx, *run.MarketRequest, reporter)
	if reporter.err != nil {
		return reporter.err
	}
	run = reporter.run
	if executionErr != nil {
		if ctx.Err() != nil {
			return s.confirmCancellation(context.Background(), run.RunID)
		}
		if failure, ok := executionErr.(marketExecutionFailure); ok {
			return s.failRun(ctx, run, failure.code, failure.message)
		}
		return s.failRun(ctx, run, "MARKET_DISCOVERY_FAILED", "market discovery failed")
	}
	artifacts := execution.Artifacts

	for _, document := range []struct {
		name    domain.RunArtifactName
		entries []any
	}{
		{domain.ArtifactMarketPlan, interfacesFromAttempts(artifacts.MarketPlan)}, {domain.ArtifactMarketSnapshot, interfacesFromSnapshots(artifacts.MarketSnapshot)}, {domain.ArtifactMarketCandidates, interfacesFromProducts(artifacts.MarketCandidates)}, {domain.ArtifactMarketGaps, interfacesFromGaps(artifacts.MarketGaps)},
	} {
		run, err = s.publishDocument(ctx, run, document.name, document.entries)
		if err != nil {
			return err
		}
	}
	status := execution.Result.Status
	_, err = s.store.Transition(ctx, domain.RunTransition{RunID: run.RunID, RunKind: domain.RunKindMarketDiscovery, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 4, Completed: 4}, MarketDiscoveryStatus: &status})
	if errors.Is(err, domain.ErrRunVersionConflict) {
		winner, getErr := s.store.Get(ctx, run.RunID)
		if getErr == nil && winner.Lifecycle == domain.RunCancelling {
			return s.confirmCancellation(context.Background(), run.RunID)
		}
	}
	return err
}

func interfacesFromCapabilities(values []domain.DiscoveryConnectorCapabilities) []any {
	out := make([]any, len(values))
	for i := range values {
		out[i] = values[i]
	}
	return out
}
func interfacesFromAttempts(values []domain.MarketDiscoveryAttempt) []any {
	out := make([]any, len(values))
	for i := range values {
		out[i] = values[i]
	}
	return out
}
func interfacesFromSnapshots(values []domain.MarketSnapshotEntry) []any {
	out := make([]any, len(values))
	for i := range values {
		out[i] = values[i]
	}
	return out
}
func interfacesFromProducts(values []domain.ProductCandidate) []any {
	out := make([]any, len(values))
	for i := range values {
		out[i] = values[i]
	}
	return out
}
func interfacesFromGaps(values []domain.MarketGap) []any {
	out := make([]any, len(values))
	for i := range values {
		out[i] = values[i]
	}
	return out
}

var _ = fmt.Sprintf

// ExportedValidISO3166Alpha2 expõe validISO3166Alpha2 para testes em L4.
// Uso exclusivo em testes de integração — não deve ser chamada em código de produção de outras camadas.
func ExportedValidISO3166Alpha2(code string) bool { return validISO3166Alpha2(code) }

// ExportedValidBCP47Locale expõe validBCP47Locale para testes em L4.
// Uso exclusivo em testes de integração — não deve ser chamada em código de produção de outras camadas.
func ExportedValidBCP47Locale(tag string) bool { return validBCP47Locale(tag) }

type marketRunProgressReporter struct {
	ctx     context.Context
	service *RecommendationRunService
	run     domain.RecommendationRun
	err     error
	cancel  context.CancelFunc
}

func (r *marketRunProgressReporter) MarketProgressError() error { return r.err }
func (r *marketRunProgressReporter) ReportMarketDiscoveryProgress(event domain.MarketDiscoveryProgressEvent) {
	if r.err != nil {
		return
	}
	updated, err := r.service.store.Transition(r.ctx, domain.RunTransition{RunID: r.run.RunID, RunKind: domain.RunKindMarketDiscovery, ExpectedVersion: r.run.Version, Lifecycle: domain.RunRunning, Stage: event.Stage, Progress: event.Progress})
	r.err = err
	if err != nil {
		r.cancel()
		return
	}
	r.run = updated
}
