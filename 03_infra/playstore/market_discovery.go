/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/playstore/market_discovery.md
 * @layer L3
 * @updated 2026-09-06
 */

package playstore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const connectorID = "google-play-store-apps"
const connectorVersion = "playstore-web-search-1.0.0"
const maximumBodyBytes = 2 << 20

type MarketDiscoveryConnector struct {
	client *http.Client
}

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
	return domain.DiscoveryConnectorCapabilities{
		ID:                   connectorID,
		Version:              connectorVersion,
		Kind:                 domain.ConnectorApplicationCatalog,
		QueryLocales: []string{
			"ar-EG", "ar-SA", "de-DE", "el-GR", "en-GB", "en-US", "en-ZA", "es-ES", "es-MX",
			"fr-FR", "hi-IN", "it-IT", "ja-JP", "kk-KZ", "ko-KR", "mn-MN", "nl-NL", "pl-PL",
			"pt-BR", "pt-PT", "ru-RU", "sv-SE", "th-TH", "tr-TR", "vi-VN", "zh-CN",
		},
		Markets: []string{
			"AE", "AR", "AU", "BR", "CA", "CH", "CL", "CN", "CO", "DE", "EG", "ES", "FR",
			"GB", "GR", "ID", "IE", "IL", "IN", "IT", "JP", "KR", "KZ", "MN", "MX", "MY",
			"NG", "NL", "NZ", "PE", "PH", "PL", "PT", "RO", "RU", "SA", "SE", "SG", "TH",
			"TR", "UA", "US", "VN", "ZA",
		},
		SupportsMarketFilter: true,
		IdentityLevels:       []string{"canonical_url", "provider_product_id"},
		EvidenceSurfaces:     []string{"name", "publisher", "storefront_observation", "summary"},
		Replayable:           true,
		PublicOnly:           true,
	}
}

func (c *MarketDiscoveryConnector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if !validRequest(request) {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	targetURL := fmt.Sprintf("https://play.google.com/store/search?q=%s&c=apps&gl=%s&hl=%s",
		url.QueryEscape(request.Query),
		url.QueryEscape(request.Market),
		url.QueryEscape(request.QueryLocale),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	req.Header.Set("User-Agent", "Bateia-PlayStore-Discovery-Agent/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return domain.ConnectorDiscovery{
			Status:     domain.AttemptEmptyObserved,
			ErrorClass: "TRANSPORT_ERROR",
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return domain.ConnectorDiscovery{
			Status:     domain.AttemptEmptyObserved,
			ErrorClass: fmt.Sprintf("HTTP_%d", resp.StatusCode),
		}, nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maximumBodyBytes+1))
	if err != nil || len(body) > maximumBodyBytes {
		return domain.ConnectorDiscovery{}, errors.New("response body exceeds budget")
	}

	return parsePlayStoreHTML(body, request)
}

var packageIDRegex = regexp.MustCompile(`/store/apps/details\?id=([a-zA-Z0-9_]+(?:\.[a-zA-Z0-9_]+)+)`)

func parsePlayStoreHTML(body []byte, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	text := string(body)
	matches := packageIDRegex.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
	}

	seen := map[string]bool{}
	candidates := []domain.ProductCandidate{}

	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		pkg := m[1]
		if seen[pkg] {
			continue
		}
		seen[pkg] = true

		canonicalURL := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s", pkg)
		name := safeCatalogText(pkg, 200)

		candidates = append(candidates, domain.ProductCandidate{
			Identity: domain.ProductIdentity{
				ConnectorID:  connectorID,
				ExternalID:   pkg,
				CanonicalURL: canonicalURL,
			},
			Name:            name,
			Summary:         "",
			ProductKind:     "application",
			Publisher:       "",
			ListingLocales:  []string{request.QueryLocale},
			ObservedMarkets: []string{request.Market},
			Resolution:      domain.IdentityDistinct,
		})

		if len(candidates) >= 10 {
			break
		}
	}

	return domain.ConnectorDiscovery{
		Status:        domain.AttemptSucceeded,
		Candidates:    candidates,
		ResultCount:   len(candidates),
		AcceptedCount: len(candidates),
	}, nil
}

func validRequest(request domain.MarketDiscoveryRequest) bool {
	if request.Query == "" || len(request.Query) > 256 || request.Market == "" || len(request.Market) != 2 || request.QueryLocale == "" {
		return false
	}
	return true
}

func safeCatalogText(value string, maximum int) string {
	if !utf8.ValidString(value) {
		return ""
	}
	var b strings.Builder
	for _, r := range value {
		if unicode.IsControl(r) && r != ' ' {
			continue
		}
		b.WriteRune(r)
		if b.Len() >= maximum {
			break
		}
	}
	return strings.TrimSpace(b.String())
}
