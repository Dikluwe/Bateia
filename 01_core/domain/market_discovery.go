/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/market_discovery.md
 * @layer L1
 * @updated 2026-09-06
 */

package domain

import "context"

type RunKind string

const (
	RunKindRecommendation       RunKind = "RECOMMENDATION"
	RunKindMarketDiscovery      RunKind = "MARKET_DISCOVERY"
	RunKindImprovementSynthesis RunKind = "IMPROVEMENT_SYNTHESIS"
)

type MarketDiscoveryStatus string

const (
	MarketDiscoveryComplete     MarketDiscoveryStatus = "COMPLETE"
	MarketDiscoveryPartial      MarketDiscoveryStatus = "PARTIAL"
	MarketDiscoveryInconclusive MarketDiscoveryStatus = "INCONCLUSIVE"
)

type DiscoveryConnectorKind string

const (
	ConnectorRepositoryCatalog  DiscoveryConnectorKind = "REPOSITORY_CATALOG"
	ConnectorApplicationCatalog DiscoveryConnectorKind = "APPLICATION_CATALOG"
	ConnectorProductCatalog     DiscoveryConnectorKind = "PRODUCT_CATALOG"
)

type MarketQueryDerivation string

const (
	QueryControlledCatalog MarketQueryDerivation = "CONTROLLED_CATALOG"
	QueryExternalAgent     MarketQueryDerivation = "EXTERNAL_AGENT_PROPOSAL"
)

type QueryDerivationStatus string

const (
	QueryDerivationControlled QueryDerivationStatus = "CONTROLLED"
	QueryDerivationUnverified QueryDerivationStatus = "UNVERIFIED_PROPOSAL"
)

type RepositoryLinkStatus string

const (
	RepositoryLinkVerified   RepositoryLinkStatus = "VERIFIED"
	RepositoryLinkUnverified RepositoryLinkStatus = "UNVERIFIED"
	RepositoryLinkConflict   RepositoryLinkStatus = "CONFLICT"
)

type IdentityResolution string

const (
	IdentitySame     IdentityResolution = "SAME"
	IdentityDistinct IdentityResolution = "DISTINCT"
	IdentityPossible IdentityResolution = "POSSIBLE"
	IdentityUnknown  IdentityResolution = "UNKNOWN"
)

type MarketAttemptStatus string

const (
	AttemptPlanned                MarketAttemptStatus = "PLANNED"
	AttemptSucceeded              MarketAttemptStatus = "SUCCEEDED"
	AttemptEmptyObserved          MarketAttemptStatus = "EMPTY_OBSERVED"
	AttemptUnknownQueryDerivation MarketAttemptStatus = "UNKNOWN_QUERY_DERIVATION"
	AttemptUnsupportedCapability  MarketAttemptStatus = "UNSUPPORTED_CAPABILITY"
	AttemptMarketNotFiltered      MarketAttemptStatus = "MARKET_NOT_FILTERED"
	AttemptUnknownProvider        MarketAttemptStatus = "UNKNOWN_PROVIDER"
	AttemptBudgetExhausted        MarketAttemptStatus = "BUDGET_EXHAUSTED"
	AttemptCancelled              MarketAttemptStatus = "CANCELLED"
	AttemptIntegrityError         MarketAttemptStatus = "INTEGRITY_ERROR"
)

type DiscoveryConnectorCapabilities struct {
	ID                   string                 `json:"id"`
	Version              string                 `json:"version"`
	Kind                 DiscoveryConnectorKind `json:"kind"`
	QueryLocales         []string               `json:"query_locales"`
	Markets              []string               `json:"markets"`
	SupportsMarketFilter bool                   `json:"supports_market_filter"`
	IdentityLevels       []string               `json:"identity_levels"`
	EvidenceSurfaces     []string               `json:"evidence_surfaces"`
	Replayable           bool                   `json:"replayable"`
	PublicOnly           bool                   `json:"public_only"`
}

type MarketQueryProposal struct {
	ID               string                `json:"id"`
	Locale           string                `json:"locale"`
	Query            string                `json:"query"`
	SourceText       string                `json:"source_text"`
	Derivation       MarketQueryDerivation `json:"derivation"`
	DerivationStatus QueryDerivationStatus `json:"derivation_status"`
	CatalogVersion   string                `json:"catalog_version,omitempty"`
	RuleID           string                `json:"rule_id,omitempty"`
	Rationale        string                `json:"rationale"`
	InputHash        string                `json:"input_hash"`
}

type MarketDiscoveryBudget struct {
	MaxQueries         int `json:"max_queries"`
	MaxResultsPerQuery int `json:"max_results_per_query"`
	MaxTotalCandidates int `json:"max_total_candidates"`
}

type MarketDiscoveryStartRequest struct {
	Idea           IdeaSpec              `json:"idea"`
	QueryProposals []MarketQueryProposal `json:"query_proposals"`
	Markets        []string              `json:"markets"`
	ConnectorIDs   []string              `json:"connector_ids"`
	Budget         MarketDiscoveryBudget `json:"budget"`
}

type MarketDiscoveryRequest struct {
	QueryID     string `json:"query_id"`
	QueryLocale string `json:"query_locale"`
	Query       string `json:"query"`
	Market      string `json:"market"`
	Limit       int    `json:"limit"`
}

type ProductIdentity struct {
	ConnectorID  string `json:"connector_id"`
	ExternalID   string `json:"external_id"`
	CanonicalURL string `json:"canonical_url,omitempty"`
}

type FOSSLinkage struct {
	UpstreamURL  string               `json:"upstream_url"`
	LicenseSPDX  string               `json:"license_spdx,omitempty"`
	IssueTracker string               `json:"issue_tracker,omitempty"`
	Status       RepositoryLinkStatus `json:"status"`
}

const (
	ConnectorAppleItunes      = "apple-itunes-search-software"
	ConnectorGooglePlay       = "google-play-store-apps"
	ConnectorFDroid           = "fdroid-foss-apps"
	ConnectorGiteeCatalog     = "gitee-repository-catalog"
	ConnectorCodebergCatalog  = "codeberg-repository-catalog"
	ConnectorGitLabCatalog    = "gitlab-repository-catalog"
	ConnectorBitbucketCatalog = "bitbucket-repository-catalog"
	ConnectorGitHubCatalog    = "github-repository-catalog"
	ConnectorSourcegraph      = "sourcegraph-repository-catalog"
	ConnectorZBMath           = "zbmath-research-software"
	ConnectorJOSS             = "joss-reviewed-software"
)

// ScholarlyEvidence records the claim made by a research catalogue.
// INDEXED is bibliographic indexing; PEER_REVIEWED is JOSS acceptance.
type ScholarlyEvidence struct {
	Catalogue       string                     `json:"catalogue"`
	RecordURL       string                     `json:"record_url"`
	DOI             string                     `json:"doi,omitempty"`
	ReviewURL       string                     `json:"review_url,omitempty"`
	Assessment      string                     `json:"assessment"`
	ReferenceCount  int                        `json:"reference_count,omitempty"`
	RelatedSoftware []ScholarlyRelatedSoftware `json:"related_software,omitempty"`
}

// ScholarlyRelatedSoftware is a catalogue link, not a verified dependency or equivalent product.
type ScholarlyRelatedSoftware struct {
	Name      string `json:"name"`
	RecordURL string `json:"record_url"`
}

type ProductRepositoryLink struct {
	OwnerRepo string               `json:"owner_repo"`
	Commit    string               `json:"commit"`
	Status    RepositoryLinkStatus `json:"status"`
	Source    string               `json:"source"`
}

type MarketCandidateOrigin struct {
	ConnectorID      string `json:"connector_id"`
	ConnectorVersion string `json:"connector_version"`
	QueryID          string `json:"query_id"`
	QueryLocale      string `json:"query_locale"`
	ObservedMarket   string `json:"observed_market"`
	ResponseHash     string `json:"response_hash"`
}

type IdentityWitness struct {
	Kind         string   `json:"kind"`
	Value        string   `json:"value"`
	ConnectorIDs []string `json:"connector_ids"`
	Bilateral    bool     `json:"bilateral"`
}

type ProductCandidate struct {
	Identity           ProductIdentity         `json:"identity"`
	Name               string                  `json:"name"`
	Summary            string                  `json:"summary,omitempty"`
	ProductKind        string                  `json:"product_kind"`
	Publisher          string                  `json:"publisher,omitempty"`
	ListingLocales     []string                `json:"listing_locales"`
	ObservedMarkets    []string                `json:"observed_markets"`
	PublisherCountries []string                `json:"publisher_countries"`
	RepositoryLinks    []ProductRepositoryLink `json:"repository_links"`
	Origins            []MarketCandidateOrigin `json:"origins"`
	Resolution         IdentityResolution      `json:"resolution"`
	Witnesses          []IdentityWitness       `json:"witnesses"`
	FOSSLink           *FOSSLinkage            `json:"foss_link,omitempty"`
	ScholarlyEvidence  *ScholarlyEvidence      `json:"scholarly_evidence,omitempty"`
}

type MarketDiscoveryAttempt struct {
	AttemptID           string              `json:"attempt_id"`
	QueryID             string              `json:"query_id"`
	QueryLocale         string              `json:"query_locale"`
	Market              string              `json:"market"`
	ConnectorID         string              `json:"connector_id"`
	ConnectorVersion    string              `json:"connector_version"`
	MarketFilterApplied bool                `json:"market_filter_applied"`
	QueryHash           string              `json:"query_hash"`
	EffectiveQueryHash  string              `json:"effective_query_hash"`
	RequestHash         string              `json:"request_hash"`
	PlannedCost         int                 `json:"planned_cost"`
	Status              MarketAttemptStatus `json:"status"`
	Reason              string              `json:"reason"`
}

type ConnectorDiscovery struct {
	AttemptID     string              `json:"attempt_id"`
	Status        MarketAttemptStatus `json:"status"`
	Candidates    []ProductCandidate  `json:"candidates"`
	ResultCount   int                 `json:"result_count"`
	AcceptedCount int                 `json:"accepted_count"`
	ResponseHash  string              `json:"response_hash"`
	ErrorClass    string              `json:"error_class"`
	ObservedAt    string              `json:"observed_at"`
}

type MarketSnapshotEntry struct {
	MarketDiscoveryAttempt
	ResultCount   int    `json:"result_count"`
	AcceptedCount int    `json:"accepted_count"`
	ErrorClass    string `json:"error_class"`
	ResponseHash  string `json:"response_hash"`
	ObservedAt    string `json:"observed_at"`
}

type MarketDiscoveryGap struct {
	AttemptID    string              `json:"attempt_id"`
	QueryLocale  string              `json:"query_locale"`
	Market       string              `json:"market"`
	ConnectorID  string              `json:"connector_id"`
	Status       MarketAttemptStatus `json:"status"`
	ReasonCode   string              `json:"reason_code"`
	ResponseHash string              `json:"response_hash"`
}

type MarketGap = MarketDiscoveryGap

type MarketDiscoveryArtifacts struct {
	ConnectorCapabilities []DiscoveryConnectorCapabilities `json:"connector_capabilities"`
	MarketPlan            []MarketDiscoveryAttempt         `json:"market_plan"`
	MarketSnapshot        []MarketSnapshotEntry            `json:"market_snapshot"`
	MarketCandidates      []ProductCandidate               `json:"market_candidates"`
	MarketGaps            []MarketDiscoveryGap             `json:"market_gaps"`
}

type MarketDiscoveryResult struct {
	SchemaVersion string                           `json:"schema_version"`
	Status        MarketDiscoveryStatus            `json:"market_discovery_status"`
	Capabilities  []DiscoveryConnectorCapabilities `json:"connector_capabilities"`
	Plan          []MarketDiscoveryAttempt         `json:"market_plan"`
	Discoveries   []ConnectorDiscovery             `json:"market_snapshot"`
	Candidates    []ProductCandidate               `json:"market_candidates"`
	Gaps          []MarketDiscoveryGap             `json:"market_gaps"`
}

type MarketDiscoveryProgressEvent struct {
	Stage    RecommendationRunStage    `json:"stage"`
	Progress RecommendationRunProgress `json:"progress"`
}

type MarketDiscoveryProgressReporter interface {
	ReportMarketDiscoveryProgress(MarketDiscoveryProgressEvent)
}

type MarketDiscoveryStartReceipt struct {
	SchemaVersion string                  `json:"schema_version"`
	RunID         RunID                   `json:"run_id"`
	RunKind       RunKind                 `json:"run_kind"`
	Lifecycle     RecommendationLifecycle `json:"lifecycle"`
	Version       uint64                  `json:"version"`
}

type MarketDiscoveryConnector interface {
	Capabilities() DiscoveryConnectorCapabilities
	Discover(context.Context, MarketDiscoveryRequest) (ConnectorDiscovery, error)
}

type ObservableMarketDiscoveryExecutor interface {
	DiscoverMarket(context.Context, MarketDiscoveryStartRequest) (MarketDiscoveryResult, error)
	DiscoverMarketObserved(context.Context, MarketDiscoveryStartRequest, MarketDiscoveryProgressReporter) (MarketDiscoveryResult, error)
}
