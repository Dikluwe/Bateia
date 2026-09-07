/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/forges.md
 * @layer L3
 * @updated 2026-09-06
 */

package forge

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

const maximumBodyBytes = 2 << 20

type BaseForgeConnector struct {
	id      string
	version string
	baseURL string
	client  *http.Client
}

func NewGiteeConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return newBaseForge("gitee-repository-catalog", "gitee-v5-1.0.0", "https://gitee.com/api/v5", client)
}

func NewCodebergConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return newBaseForge("codeberg-repository-catalog", "forgejo-v1-1.0.0", "https://codeberg.org/api/v1", client)
}

func NewGitLabConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return newBaseForge("gitlab-repository-catalog", "gitlab-v4-1.0.0", "https://gitlab.com/api/v4", client)
}

func newBaseForge(id, version, baseURL string, client *http.Client) domain.MarketDiscoveryConnector {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	clone := *client
	if clone.Timeout <= 0 || clone.Timeout > 30*time.Second {
		clone.Timeout = 15 * time.Second
	}
	clone.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirect forbidden") }
	return &BaseForgeConnector{id: id, version: version, baseURL: baseURL, client: &clone}
}

func (c *BaseForgeConnector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{
		ID:                   c.id,
		Version:              c.version,
		Kind:                 domain.ConnectorRepositoryCatalog,
		QueryLocales: []string{
			"ar-EG", "ar-SA", "de-DE", "el-GR", "en-GB", "en-US", "en-ZA", "es-ES", "es-MX",
			"fr-FR", "hi-IN", "it-IT", "ja-JP", "kk-KZ", "ko-KR", "mn-MN", "nl-NL", "pl-PL",
			"pt-BR", "pt-PT", "ru-RU", "sv-SE", "th-TH", "tr-TR", "vi-VN", "zh-CN",
		},
		Markets:              []string{},
		SupportsMarketFilter: false,
		IdentityLevels:       []string{"commit", "owner_repo"},
		EvidenceSurfaces:     []string{"code_occurrence", "immutable_commit", "repository_identity"},
		Replayable:           true,
		PublicOnly:           true,
	}
}

func (c *BaseForgeConnector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if request.Query == "" || len(request.Query) > 256 {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	var searchURL string
	switch c.id {
	case "gitee-repository-catalog":
		searchURL = fmt.Sprintf("%s/search/repositories?q=%s&per_page=10", c.baseURL, url.QueryEscape(request.Query))
	case "codeberg-repository-catalog":
		searchURL = fmt.Sprintf("%s/repos/search?q=%s&limit=10", c.baseURL, url.QueryEscape(request.Query))
	case "gitlab-repository-catalog":
		searchURL = fmt.Sprintf("%s/projects?search=%s&per_page=10", c.baseURL, url.QueryEscape(request.Query))
	default:
		return domain.ConnectorDiscovery{Status: domain.AttemptUnknownProvider}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	req.Header.Set("User-Agent", "Bateia-Forge-Discovery-Agent/1.0")
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

	return parseForgeResponse(c.id, body, request)
}

type genericForgeItem struct {
	FullName    string `json:"full_name"`
	PathWithNS  string `json:"path_with_namespace"`
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	WebURL      string `json:"web_url"`
	Description string `json:"description"`
	DefaultBranch string `json:"default_branch"`
}

func parseForgeResponse(connectorID string, body []byte, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	var items []genericForgeItem
	if err := json.Unmarshal(body, &items); err != nil {
		var envelope struct {
			Data []genericForgeItem `json:"data"`
		}
		if err2 := json.Unmarshal(body, &envelope); err2 != nil {
			return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
		}
		items = envelope.Data
	}

	if len(items) == 0 {
		return domain.ConnectorDiscovery{Status: domain.AttemptEmptyObserved, ResultCount: 0, AcceptedCount: 0}, nil
	}

	candidates := []domain.ProductCandidate{}
	for _, it := range items {
		ownerRepo := safeCatalogText(it.FullName, 200)
		if ownerRepo == "" {
			ownerRepo = safeCatalogText(it.PathWithNS, 200)
		}
		if ownerRepo == "" {
			ownerRepo = safeCatalogText(it.Name, 200)
		}
		if ownerRepo == "" {
			continue
		}

		canonicalURL := it.HTMLURL
		if canonicalURL == "" {
			canonicalURL = it.WebURL
		}

		candidates = append(candidates, domain.ProductCandidate{
			Identity: domain.ProductIdentity{
				ConnectorID:  connectorID,
				ExternalID:   ownerRepo,
				CanonicalURL: canonicalURL,
			},
			Name:            ownerRepo,
			Summary:         safeCatalogText(it.Description, 500),
			ProductKind:     "repository",
			ListingLocales:  []string{request.QueryLocale},
			ObservedMarkets: []string{},
			RepositoryLinks: []domain.ProductRepositoryLink{
				{
					OwnerRepo: ownerRepo,
					Commit:    it.DefaultBranch,
					Status:    domain.RepositoryLinkVerified,
					Source:    "forge_api_metadata",
				},
			},
			Resolution: domain.IdentityDistinct,
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
