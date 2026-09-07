/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/fdroid/market_discovery.md
 * @layer L3
 * @updated 2026-09-06
 */

package fdroid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const connectorID = "fdroid-foss-apps"
const connectorVersion = "fdroid-index-v2-1.0.0"
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
		Markets:              []string{},
		SupportsMarketFilter: false,
		IdentityLevels:       []string{"canonical_url", "provider_product_id", "verified_official_repository"},
		EvidenceSurfaces:     []string{"name", "summary", "repository_identity", "license"},
		Replayable:           true,
		PublicOnly:           true,
	}
}

func (c *MarketDiscoveryConnector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if request.Query == "" || len(request.Query) > 256 {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	targetURL := fmt.Sprintf("https://search.f-droid.org/api/v1/packages?q=%s", url.QueryEscape(request.Query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	req.Header.Set("User-Agent", "Bateia-FDroid-Discovery-Agent/1.0")
	req.Header.Set("Accept", "application/json")

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

	return parseFDroidResponse(body, request)
}

type fdroidPackageItem struct {
	PackageName string `json:"packageName"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	SourceCode  string `json:"sourceCode"`
	License     string `json:"license"`
	WebURL      string `json:"webUrl"`
	Tracker     string `json:"issueTracker"`
}

func parseFDroidResponse(body []byte, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	var items []fdroidPackageItem
	if err := json.Unmarshal(body, &items); err != nil {
		var envelope struct {
			Packages []fdroidPackageItem `json:"packages"`
		}
		if err2 := json.Unmarshal(body, &envelope); err2 != nil {
			return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
		}
		items = envelope.Packages
	}

	if len(items) == 0 {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
	}

	candidates := []domain.ProductCandidate{}
	for _, it := range items {
		pkg := safeCatalogText(it.PackageName, 100)
		if pkg == "" {
			continue
		}
		name := safeCatalogText(it.Name, 200)
		if name == "" {
			name = pkg
		}
		summary := safeCatalogText(it.Summary, 500)
		canonicalURL := it.WebURL
		if canonicalURL == "" {
			canonicalURL = fmt.Sprintf("https://f-droid.org/packages/%s/", pkg)
		}

		repoLinks := []domain.ProductRepositoryLink{}
		var fossLink *domain.FOSSLinkage
		if it.SourceCode != "" {
			ownerRepo := extractOwnerRepo(it.SourceCode)
			if ownerRepo != "" {
				repoLinks = append(repoLinks, domain.ProductRepositoryLink{
					OwnerRepo: ownerRepo,
					Commit:    "",
					Status:    domain.RepositoryLinkVerified,
					Source:    "fdroid_source_code_metadata",
				})
			}
			fossLink = &domain.FOSSLinkage{
				UpstreamURL:  it.SourceCode,
				LicenseSPDX:  safeCatalogText(it.License, 50),
				IssueTracker: it.Tracker,
				Status:       domain.RepositoryLinkVerified,
			}
		}

		candidates = append(candidates, domain.ProductCandidate{
			Identity: domain.ProductIdentity{
				ConnectorID:  connectorID,
				ExternalID:   pkg,
				CanonicalURL: canonicalURL,
			},
			Name:            name,
			Summary:         summary,
			ProductKind:     "application",
			Publisher:       "F-Droid Community",
			ListingLocales:  []string{request.QueryLocale},
			ObservedMarkets: []string{},
			RepositoryLinks: repoLinks,
			FOSSLink:        fossLink,
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

func extractOwnerRepo(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) >= 2 {
		return parts[0] + "/" + strings.TrimSuffix(parts[1], ".git")
	}
	return ""
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
