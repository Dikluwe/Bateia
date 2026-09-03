/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/itunes/market_discovery.md
 * @layer L3
 * @updated 2026-09-01
 */

package itunes

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const connectorID = "apple-itunes-search-software"
const connectorVersion = "itunes-search-doc-1.0.3-response-v2"
const maximumBodyBytes = 2 << 20

type MarketDiscoveryConnector struct{ client *http.Client }

func NewMarketDiscoveryConnector(client *http.Client) (domain.MarketDiscoveryConnector, error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	clone := *client
	if clone.Timeout <= 0 || clone.Timeout > 30*time.Second {
		clone.Timeout = 15 * time.Second
	}
	clone.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirect forbidden") }
	return &MarketDiscoveryConnector{client: &clone}, nil
}

func (c *MarketDiscoveryConnector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: connectorID, Version: connectorVersion, Kind: domain.ConnectorApplicationCatalog,
		QueryLocales: []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"},
		Markets:      []string{"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US"}, SupportsMarketFilter: true,
		IdentityLevels: []string{"canonical_url", "provider_product_id"}, EvidenceSurfaces: []string{"name", "provider_declared_language_codes", "publisher", "storefront_observation", "summary"}, Replayable: true, PublicOnly: true}
}

func (c *MarketDiscoveryConnector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if ctx.Err() != nil {
		return domain.ConnectorDiscovery{Status: domain.AttemptCancelled, Candidates: []domain.ProductCandidate{}, ErrorClass: "CANCELLED"}, nil
	}
	if c == nil || c.client == nil || !validRequest(request) {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	endpoint := "https://itunes.apple.com/search?" + fixedQuery(request)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	httpRequest.Header.Set("Accept", "application/json")
	response, err := c.client.Do(httpRequest)
	if err != nil {
		if ctx.Err() != nil {
			return domain.ConnectorDiscovery{Status: domain.AttemptCancelled, Candidates: []domain.ProductCandidate{}, ErrorClass: "CANCELLED"}, nil
		}
		return unknownDiscovery(), nil
	}
	defer response.Body.Close()
	limited := io.LimitReader(response.Body, maximumBodyBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil || len(body) > maximumBodyBytes || response.StatusCode < 200 || response.StatusCode >= 300 || !utf8.Valid(body) || !strictJSON(body) {
		return unknownDiscovery(), nil
	}
	projection, err := parseResponse(body, request)
	if err != nil {
		return unknownDiscovery(), nil
	}
	digest := sha256.Sum256(body)
	projection.ResponseHash = "sha256:" + hex.EncodeToString(digest[:])
	if date, parseErr := http.ParseTime(response.Header.Get("Date")); parseErr == nil {
		projection.ObservedAt = date.UTC().Format(time.RFC3339)
	}
	return projection, nil
}

func validRequest(request domain.MarketDiscoveryRequest) bool {
	if request.QueryID == "" || request.QueryLocale == "" || strings.TrimSpace(request.Query) == "" || request.Market == "" || request.Limit < 1 || request.Limit > 10 {
		return false
	}
	allowedLocale := false
	for _, value := range []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"} {
		allowedLocale = allowedLocale || request.QueryLocale == value
	}
	allowedMarket := false
	for _, value := range []string{"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US"} {
		allowedMarket = allowedMarket || request.Market == value
	}
	if !allowedLocale || !allowedMarket || len([]rune(request.Query)) > 256 {
		return false
	}
	for _, r := range request.Query {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func fixedQuery(request domain.MarketDiscoveryRequest) string {
	encode := func(value string) string { return strings.ReplaceAll(url.QueryEscape(value), "+", "%20") }
	return "country=" + encode(request.Market) + "&entity=software&explicit=No&lang=en_us&limit=" + strconv.Itoa(request.Limit) + "&media=software&term=" + encode(request.Query) + "&version=2"
}

type responseEnvelope struct {
	ResultCount int            `json:"resultCount"`
	Results     []responseItem `json:"results"`
}
type responseItem struct {
	WrapperType       string      `json:"wrapperType"`
	Kind              string      `json:"kind"`
	TrackID           json.Number `json:"trackId"`
	TrackName         string      `json:"trackName"`
	TrackCensoredName string      `json:"trackCensoredName"`
	TrackViewURL      string      `json:"trackViewUrl"`
	ArtistName        string      `json:"artistName"`
	SellerName        string      `json:"sellerName"`
	BundleID          string      `json:"bundleId"`
	Description       string      `json:"description"`
	Genres            []string    `json:"genres"`
	PrimaryGenreName  string      `json:"primaryGenreName"`
	LanguageCodes     []string    `json:"languageCodesISO2A"`
	SellerURL         string      `json:"sellerUrl"`
}

func parseResponse(body []byte, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var envelope responseEnvelope
	if err := decoder.Decode(&envelope); err != nil || decoder.Decode(&struct{}{}) != io.EOF || envelope.ResultCount != len(envelope.Results) {
		return domain.ConnectorDiscovery{}, errors.New("invalid provider response")
	}
	if envelope.ResultCount == 0 {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, Candidates: []domain.ProductCandidate{}, ResultCount: 0, AcceptedCount: 0, ErrorClass: ""}, nil
	}
	candidates := []domain.ProductCandidate{}
	for _, item := range envelope.Results {
		externalID := item.TrackID.String()
		if externalID == "" || externalID == "0" {
			continue
		}
		if _, err := strconv.ParseUint(externalID, 10, 64); err != nil {
			continue
		}
		canonical := canonicalAppleURL(item.TrackViewURL, externalID)
		name := safeCatalogText(item.TrackName, 200)
		if name == "" {
			name = safeCatalogText(item.TrackCensoredName, 200)
		}
		if name == "" {
			continue
		}
		publisher := safeCatalogText(item.SellerName, 200)
		if publisher == "" {
			publisher = safeCatalogText(item.ArtistName, 200)
		}
		candidates = append(candidates, domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: connectorID, ExternalID: externalID, CanonicalURL: canonical}, Name: name, Summary: "", ProductKind: "application", Publisher: publisher, ListingLocales: []string{}, ObservedMarkets: []string{request.Market}, PublisherCountries: []string{}, RepositoryLinks: []domain.ProductRepositoryLink{}, Origins: []domain.MarketCandidateOrigin{}, Resolution: domain.IdentityDistinct, Witnesses: []domain.IdentityWitness{}})
	}
	status := domain.AttemptSucceeded
	if len(candidates) == 0 {
		status = domain.AttemptUnknownProvider
	}
	errorClass := ""
	if status == domain.AttemptUnknownProvider {
		errorClass = "UNKNOWN_PROVIDER"
	}
	return domain.ConnectorDiscovery{Status: status, Candidates: candidates, ResultCount: envelope.ResultCount, AcceptedCount: len(candidates), ErrorClass: errorClass}, nil
}

func canonicalAppleURL(raw, externalID string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || strings.ToLower(parsed.Hostname()) != "apps.apple.com" || parsed.Port() != "" || !strings.Contains(parsed.EscapedPath(), "/id"+externalID) {
		return ""
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.User = nil
	return parsed.String()
}

func safeCatalogText(value string, maximum int) string {
	value = strings.Join(strings.Fields(value), " ")
	lower := strings.ToLower(value)
	for _, bad := range []string{"canary", "<script", "authorization:", "bearer ", "ignore previous", "system prompt"} {
		if strings.Contains(lower, bad) {
			return ""
		}
	}
	runes := []rune(value)
	if len(runes) > maximum {
		runes = runes[:maximum]
	}
	for _, r := range runes {
		if unicode.IsControl(r) {
			return ""
		}
	}
	return string(runes)
}

func unknownDiscovery() domain.ConnectorDiscovery {
	return domain.ConnectorDiscovery{Status: domain.AttemptUnknownProvider, Candidates: []domain.ProductCandidate{}, ErrorClass: "UNKNOWN_PROVIDER"}
}

func strictJSON(data []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	stack := []map[string]bool{}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}
		switch value := token.(type) {
		case json.Delim:
			switch value {
			case '{':
				stack = append(stack, map[string]bool{})
			case '}':
				if len(stack) == 0 {
					return false
				}
				stack = stack[:len(stack)-1]
			}
		case string:
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				if top[value] {
					return false
				}
				top[value] = true
				var discard json.RawMessage
				if decoder.Decode(&discard) != nil {
					return false
				}
			}
		}
	}
	var value interface{}
	decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(&value) == nil && decoder.Decode(&struct{}{}) == io.EOF
}
