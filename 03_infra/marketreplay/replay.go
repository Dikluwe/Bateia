/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/marketreplay/replay.md
 * @layer L3
 * @updated 2026-09-01
 */

package marketreplay

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const replayConnectorID = "apple-itunes-search-software"
const replayConnectorVersion = "itunes-search-doc-1.0.3-response-v2"

type fixtureRequest struct {
	Method          string            `json:"method"`
	Scheme          string            `json:"scheme"`
	Host            string            `json:"host"`
	Port            int               `json:"port"`
	Path            string            `json:"path"`
	FollowRedirects bool              `json:"follow_redirects"`
	QueryLocale     string            `json:"query_locale,omitempty"`
	Parameters      fixtureParameters `json:"parameters"`
}
type fixtureParameters struct {
	Term       string `json:"term,omitempty"`
	TermSHA256 string `json:"term_sha256"`
	Country    string `json:"country"`
	Media      string `json:"media"`
	Entity     string `json:"entity"`
	Limit      int    `json:"limit"`
	Lang       string `json:"lang"`
	Version    string `json:"version"`
	Explicit   string `json:"explicit"`
}
type fixtureCapture struct {
	HTTPStatusClass       string  `json:"http_status_class"`
	ProviderContentType   *string `json:"provider_content_type"`
	CapturedAt            string  `json:"captured_at"`
	ProviderBodySHA256    string  `json:"provider_body_sha256"`
	ProviderBodyRetained  bool    `json:"provider_body_retained"`
	ProjectionAllowlisted bool    `json:"projection_allowlisted"`
	ErrorClass            string  `json:"error_class,omitempty"`
}
type projectionItem struct {
	ArtistName        string   `json:"artistName,omitempty"`
	BundleID          string   `json:"bundleId,omitempty"`
	Description       string   `json:"description,omitempty"`
	Genres            []string `json:"genres,omitempty"`
	Kind              string   `json:"kind,omitempty"`
	LanguageCodes     []string `json:"languageCodesISO2A,omitempty"`
	PrimaryGenreName  string   `json:"primaryGenreName,omitempty"`
	SellerName        string   `json:"sellerName,omitempty"`
	SellerURL         string   `json:"sellerUrl,omitempty"`
	TrackCensoredName string   `json:"trackCensoredName,omitempty"`
	TrackID           uint64   `json:"trackId,omitempty"`
	TrackName         string   `json:"trackName,omitempty"`
	TrackViewURL      string   `json:"trackViewUrl,omitempty"`
	WrapperType       string   `json:"wrapperType,omitempty"`
}
type fixtureProjection struct {
	ResultCount int              `json:"resultCount"`
	Results     []projectionItem `json:"results"`
}
type fixture struct {
	SchemaVersion      string             `json:"schema_version"`
	FixtureID          string             `json:"fixture_id"`
	Visibility         string             `json:"visibility"`
	Synthetic          bool               `json:"synthetic"`
	ConnectorID        string             `json:"connector_id"`
	ConnectorVersion   string             `json:"connector_version"`
	Case               string             `json:"case"`
	ExpectedOutcome    string             `json:"expected_outcome"`
	Request            fixtureRequest     `json:"request"`
	RequestSHA256      string             `json:"request_sha256"`
	Capture            fixtureCapture     `json:"capture"`
	Projection         *fixtureProjection `json:"projection"`
	ExpectedProjection json.RawMessage    `json:"expected_projection,omitempty"`
}

type Connector struct{ fixture fixture }

func NewConnector(encoded []byte) (domain.MarketDiscoveryConnector, error) {
	if !utf8.Valid(encoded) {
		return nil, domain.ErrIntegrity
	}
	decoder := json.NewDecoder(strings.NewReader(string(encoded)))
	decoder.DisallowUnknownFields()
	var value fixture
	if err := decoder.Decode(&value); err != nil {
		return nil, domain.ErrIntegrity
	}
	var extra interface{}
	if decoder.Decode(&extra) != io.EOF {
		return nil, domain.ErrIntegrity
	}
	if value.SchemaVersion != "bateia.apple-offline-fixture.v1" || value.Visibility != "PROTECTED_ORACLE_ONLY" || value.ConnectorID != replayConnectorID || value.ConnectorVersion != replayConnectorVersion || value.FixtureID == "" || value.RequestSHA256 == "" || value.Request.Method != "GET" || value.Request.Scheme != "https" || value.Request.Host != "itunes.apple.com" || value.Request.Port != 443 || value.Request.Path != "/search" || value.Request.FollowRedirects || value.Request.Parameters.Media != "software" || value.Request.Parameters.Entity != "software" || value.Request.Parameters.Lang != "en_us" || value.Request.Parameters.Version != "2" || value.Request.Parameters.Explicit != "No" || value.Request.Parameters.Limit < 1 || value.Request.Parameters.Limit > 10 || value.Capture.ProviderBodyRetained || !value.Capture.ProjectionAllowlisted {
		return nil, domain.ErrIntegrity
	}
	canonical, err := canonicalJSON(value.Request)
	if err != nil || shaHex(canonical) != value.RequestSHA256 {
		return nil, domain.ErrIntegrity
	}
	if value.Request.Parameters.Term != "" && shaHex([]byte(value.Request.Parameters.Term)) != value.Request.Parameters.TermSHA256 {
		return nil, domain.ErrIntegrity
	}
	if value.Projection != nil && (value.Projection.ResultCount != len(value.Projection.Results) || value.Projection.Results == nil) {
		return nil, domain.ErrIntegrity
	}
	return &Connector{fixture: value}, nil
}

func (c *Connector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: replayConnectorID, Version: replayConnectorVersion, Kind: domain.ConnectorApplicationCatalog, QueryLocales: []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"}, Markets: []string{"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US"}, SupportsMarketFilter: true, IdentityLevels: []string{"canonical_url", "provider_product_id"}, EvidenceSurfaces: []string{"name", "provider_declared_language_codes", "publisher", "storefront_observation", "summary"}, Replayable: true, PublicOnly: true}
}

func (c *Connector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if c == nil || ctx.Err() != nil {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	parameters := c.fixture.Request.Parameters
	if request.QueryID == "" || request.QueryLocale == "" || request.Market != parameters.Country || request.Limit != parameters.Limit || c.fixture.Request.QueryLocale != "" && request.QueryLocale != c.fixture.Request.QueryLocale || shaHex([]byte(request.Query)) != parameters.TermSHA256 {
		return integrity(), nil
	}
	if parameters.Term != "" && request.Query != parameters.Term {
		return integrity(), nil
	}
	if c.fixture.ExpectedOutcome == "INTEGRITY_ERROR" {
		return integrity(), nil
	}
	if c.fixture.Projection == nil || c.fixture.Capture.ErrorClass != "" || !strings.HasPrefix(c.fixture.Capture.HTTPStatusClass, "2") {
		return domain.ConnectorDiscovery{Status: domain.AttemptUnknownProvider, Candidates: []domain.ProductCandidate{}, ErrorClass: "UNKNOWN_PROVIDER", ObservedAt: c.fixture.Capture.CapturedAt, ResponseHash: c.fixture.Capture.ProviderBodySHA256}, nil
	}
	if c.fixture.Projection.ResultCount == 0 {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, Candidates: []domain.ProductCandidate{}, ResultCount: 0, AcceptedCount: 0, ObservedAt: c.fixture.Capture.CapturedAt, ResponseHash: c.fixture.Capture.ProviderBodySHA256}, nil
	}
	candidates := []domain.ProductCandidate{}
	for _, item := range c.fixture.Projection.Results {
		if item.TrackID == 0 {
			continue
		}
		externalID := strconv.FormatUint(item.TrackID, 10)
		name := safeText(item.TrackName)
		if name == "" {
			name = safeText(item.TrackCensoredName)
		}
		if name == "" {
			continue
		}
		publisher := safeText(item.SellerName)
		if publisher == "" {
			publisher = safeText(item.ArtistName)
		}
		candidates = append(candidates, domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: replayConnectorID, ExternalID: externalID, CanonicalURL: canonicalURL(item.TrackViewURL, externalID)}, Name: name, Summary: "", ProductKind: "application", Publisher: publisher, ListingLocales: []string{}, ObservedMarkets: []string{request.Market}, PublisherCountries: []string{}, RepositoryLinks: []domain.ProductRepositoryLink{}, Origins: []domain.MarketCandidateOrigin{}, Resolution: domain.IdentityDistinct, Witnesses: []domain.IdentityWitness{}})
	}
	if len(candidates) == 0 {
		return integrity(), nil
	}
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: candidates, ResultCount: c.fixture.Projection.ResultCount, AcceptedCount: len(candidates), ResponseHash: c.fixture.Capture.ProviderBodySHA256, ObservedAt: c.fixture.Capture.CapturedAt}, nil
}

func integrity() domain.ConnectorDiscovery {
	return domain.ConnectorDiscovery{Status: domain.AttemptIntegrityError, Candidates: []domain.ProductCandidate{}, ErrorClass: "INTEGRITY_ERROR"}
}
func shaHex(value []byte) string {
	digest := sha256.Sum256(value)
	return hex.EncodeToString(digest[:])
}
func safeText(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	lower := strings.ToLower(value)
	for _, bad := range []string{"canary", "<script", "authorization:", "bearer ", "ignore previous", "system prompt"} {
		if strings.Contains(lower, bad) {
			return ""
		}
	}
	runes := []rune(value)
	if len(runes) > 200 {
		runes = runes[:200]
	}
	for _, r := range runes {
		if unicode.IsControl(r) {
			return ""
		}
	}
	return string(runes)
}
func canonicalURL(raw, id string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || strings.ToLower(parsed.Hostname()) != "apps.apple.com" || parsed.Port() != "" || !strings.Contains(parsed.EscapedPath(), "/id"+id) {
		return ""
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.User = nil
	return parsed.String()
}

func canonicalJSON(value interface{}) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var generic interface{}
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	if err := decoder.Decode(&generic); err != nil {
		return nil, err
	}
	builder := &strings.Builder{}
	if err := writeCanonical(builder, generic); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}
func writeCanonical(builder *strings.Builder, value interface{}) error {
	switch item := value.(type) {
	case nil:
		builder.WriteString("null")
	case bool:
		if item {
			builder.WriteString("true")
		} else {
			builder.WriteString("false")
		}
	case json.Number:
		builder.WriteString(item.String())
	case string:
		encoded, _ := json.Marshal(item)
		builder.Write(encoded)
	case []interface{}:
		builder.WriteByte('[')
		for i, entry := range item {
			if i > 0 {
				builder.WriteByte(',')
			}
			if err := writeCanonical(builder, entry); err != nil {
				return err
			}
		}
		builder.WriteByte(']')
	case map[string]interface{}:
		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		builder.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				builder.WriteByte(',')
			}
			encoded, _ := json.Marshal(key)
			builder.Write(encoded)
			builder.WriteByte(':')
			if err := writeCanonical(builder, item[key]); err != nil {
				return err
			}
		}
		builder.WriteByte('}')
	default:
		return errors.New("unsupported canonical value")
	}
	return nil
}
