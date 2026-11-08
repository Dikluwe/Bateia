/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/forges.md
 * @layer L3
 * @updated 2026-10-28
 */

package forge

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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

const maximumBodyBytes = 2 << 20

type BaseForgeConnector struct {
	id      string
	version string
	baseURL string
	client  *http.Client
	token   string
}

func NewGiteeConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return NewGiteeConnectorWithToken("", client)
}

// NewGiteeConnectorWithToken optionally authenticates public catalog searches.
func NewGiteeConnectorWithToken(token string, client *http.Client) domain.MarketDiscoveryConnector {
	c := newBaseForge(domain.ConnectorGiteeCatalog, "gitee-v5-2.0.0", "https://gitee.com/api/v5", client)
	c.token = strings.TrimSpace(token)
	return c
}

// NewAtomGitConnector requires a personal access token to search the public catalog.
func NewAtomGitConnector(token string, client *http.Client) domain.MarketDiscoveryConnector {
	c := newBaseForge(domain.ConnectorAtomGitCatalog, "atomgit-v5-1.0.0", "https://api.atomgit.com/api/v5", client)
	c.token = strings.TrimSpace(token)
	return c
}

func NewCodebergConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return newBaseForge("codeberg-repository-catalog", "forgejo-v1-2.0.0", "https://codeberg.org/api/v1", client)
}

func NewGitLabConnector(client *http.Client) domain.MarketDiscoveryConnector {
	return newBaseForge("gitlab-repository-catalog", "gitlab-v4-2.0.0", "https://gitlab.com/api/v4", client)
}

func newBaseForge(id, version, baseURL string, client *http.Client) *BaseForgeConnector {
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
		ID:      c.id,
		Version: c.version,
		Kind:    domain.ConnectorRepositoryCatalog,
		QueryLocales: []string{
			"ar-EG", "ar-SA", "de-DE", "el-GR", "en-GB", "en-US", "en-ZA", "es-ES", "es-MX",
			"fr-FR", "hi-IN", "it-IT", "ja-JP", "kk-KZ", "ko-KR", "mn-MN", "nl-NL", "pl-PL",
			"pt-BR", "pt-PT", "ru-RU", "sv-SE", "th-TH", "tr-TR", "vi-VN", "zh-CN",
		},
		Markets:              []string{},
		SupportsMarketFilter: false,
		IdentityLevels:       []string{"owner_repo"},
		EvidenceSurfaces:     []string{"repository_identity"},
		Replayable:           true,
		PublicOnly:           true,
	}
}

func (c *BaseForgeConnector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	query := strings.TrimSpace(request.Query)
	if query == "" || len(query) > 256 || !utf8.ValidString(query) {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	if ctx.Err() != nil {
		return forgeFailure(domain.AttemptCancelled, "CANCELLED"), nil
	}
	if c.id == domain.ConnectorAtomGitCatalog && c.token == "" {
		return forgeFailure(domain.AttemptUnknownProvider, "AUTH_REQUIRED"), nil
	}
	params := url.Values{"q": {query}, "per_page": {strconv.Itoa(forgeLimit(request.Limit))}, "page": {"1"}}
	var path string
	switch c.id {
	case domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog:
		path = "/search/repositories"
		if c.id == domain.ConnectorGiteeCatalog && c.token != "" {
			params.Set("access_token", c.token)
		}
	case domain.ConnectorCodebergCatalog:
		path = "/repos/search"
		params.Set("limit", params.Get("per_page"))
		params.Del("per_page")
	case domain.ConnectorGitLabCatalog:
		path = "/projects"
		params.Set("search", query)
		params.Set("visibility", "public")
		params.Del("q")
	default:
		return forgeFailure(domain.AttemptUnknownProvider, "UNKNOWN_PROVIDER"), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+params.Encode(), nil)
	if err != nil {
		return forgeFailure(domain.AttemptUnknownProvider, "INVALID_REQUEST"), nil
	}
	req.Header.Set("User-Agent", "Bateia-Forge-Discovery-Agent/1.0")
	req.Header.Set("Accept", "application/json")
	if c.id == domain.ConnectorAtomGitCatalog {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return forgeFailure(domain.AttemptCancelled, "CANCELLED"), nil
		}
		return forgeFailure(domain.AttemptUnknownProvider, "TRANSPORT_ERROR"), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return forgeFailure(domain.AttemptUnknownProvider, fmt.Sprintf("HTTP_%d", resp.StatusCode)), nil
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maximumBodyBytes+1))
	if ctx.Err() != nil {
		return forgeFailure(domain.AttemptCancelled, "CANCELLED"), nil
	}
	if err != nil || len(body) > maximumBodyBytes {
		return forgeFailure(domain.AttemptIntegrityError, "INVALID_RESPONSE"), nil
	}
	result, err := parseForgeResponse(c.id, body, request)
	result.ResponseHash = fmt.Sprintf("sha256:%x", sha256.Sum256(body))
	result.ObservedAt = time.Now().UTC().Format(time.RFC3339Nano)
	return result, err
}

func forgeLimit(limit int) int {
	if limit <= 0 || limit > 10 {
		return 10
	}
	return limit
}

func forgeFailure(status domain.MarketAttemptStatus, class string) domain.ConnectorDiscovery {
	return domain.ConnectorDiscovery{Status: status, Candidates: []domain.ProductCandidate{}, ErrorClass: class}
}

type genericForgeItem struct {
	FullName    string          `json:"full_name"`
	PathWithNS  string          `json:"path_with_namespace"`
	Name        string          `json:"name"`
	HTMLURL     string          `json:"html_url"`
	WebURL      string          `json:"web_url"`
	Description string          `json:"description"`
	Private     json.RawMessage `json:"private"`
	Public      json.RawMessage `json:"public"`
	Internal    json.RawMessage `json:"internal"`
	Visibility  string          `json:"visibility"`
}

func parseForgeResponse(connectorID string, body []byte, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	var items []genericForgeItem
	if err := json.Unmarshal(body, &items); err != nil {
		if connectorID != domain.ConnectorCodebergCatalog {
			return forgeFailure(domain.AttemptIntegrityError, "INVALID_RESPONSE"), nil
		}
		var envelope struct {
			Data []genericForgeItem `json:"data"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			return forgeFailure(domain.AttemptIntegrityError, "INVALID_RESPONSE"), nil
		}
		items = envelope.Data
	}
	if items == nil {
		return forgeFailure(domain.AttemptIntegrityError, "INVALID_RESPONSE"), nil
	}
	if len(items) == 0 {
		return forgeFailure(domain.AttemptEmptyObserved, ""), nil
	}
	candidates := []domain.ProductCandidate{}
	for _, it := range items {
		if !publicForgeItem(connectorID, it) {
			continue
		}
		ownerRepo := it.FullName
		if ownerRepo == "" {
			ownerRepo = it.PathWithNS
		}
		canonicalURL := it.HTMLURL
		if canonicalURL == "" {
			canonicalURL = it.WebURL
		}
		if !validForgeIdentity(connectorID, ownerRepo, canonicalURL) {
			continue
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
					Commit:    "",
					Status:    domain.RepositoryLinkUnverified,
					Source:    "forge_api_metadata",
				},
			},
			Resolution: domain.IdentityDistinct,
		})

		if len(candidates) >= forgeLimit(request.Limit) {
			break
		}
	}

	if len(candidates) == 0 {
		return forgeFailure(domain.AttemptIntegrityError, "NO_PUBLIC_IDENTITIES"), nil
	}
	return domain.ConnectorDiscovery{
		Status:        domain.AttemptSucceeded,
		Candidates:    candidates,
		ResultCount:   len(candidates),
		AcceptedCount: len(candidates),
	}, nil
}

// flagValue accepts the boolean examples and integer schema of the v5 APIs.
func flagValue(raw json.RawMessage, expected bool) bool {
	value := strings.TrimSpace(string(raw))
	if expected {
		return value == "true" || value == "1"
	}
	return value == "false" || value == "0"
}

func publicForgeItem(id string, item genericForgeItem) bool {
	if len(item.Private) > 0 && !flagValue(item.Private, false) {
		return false
	}
	if len(item.Internal) > 0 && !flagValue(item.Internal, false) {
		return false
	}
	if len(item.Public) > 0 && !flagValue(item.Public, true) {
		return false
	}
	if item.Visibility != "" && item.Visibility != "public" {
		return false
	}
	if id == domain.ConnectorGiteeCatalog || id == domain.ConnectorAtomGitCatalog {
		if !flagValue(item.Private, false) {
			return false
		}
	}
	return id != domain.ConnectorAtomGitCatalog || flagValue(item.Public, true)
}

func validForgeIdentity(id, ownerRepo, rawURL string) bool {
	if len(ownerRepo) > 200 || !utf8.ValidString(ownerRepo) {
		return false
	}
	parts := strings.Split(ownerRepo, "/")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return false
		}
		for _, r := range part {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != '.' {
				return false
			}
		}
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Port() != "" || u.Path != "/"+ownerRepo {
		return false
	}
	switch id {
	case domain.ConnectorGiteeCatalog:
		return u.Host == "gitee.com"
	case domain.ConnectorAtomGitCatalog:
		return u.Host == "atomgit.com" || u.Host == "gitcode.com"
	case domain.ConnectorCodebergCatalog:
		return u.Host == "codeberg.org"
	case domain.ConnectorGitLabCatalog:
		return u.Host == "gitlab.com"
	}
	return false
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
