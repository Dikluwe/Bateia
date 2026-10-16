/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/bitbucket/connector.md
 * @layer L3
 * @updated 2026-10-06
 */

package bitbucket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
)

const connectorID = "bitbucket-repository-catalog"
const maxResponseBytes = 2 << 20
const maxWorkspaces = 20

var workspacePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,99}$`)
var repoPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,62}$`)

type Connector struct {
	workspaces []string
	client     *http.Client
	baseURL    string
}

// NewConnector searches public repositories in explicitly configured workspaces.
func NewConnector(workspaces []string, client *http.Client) (*Connector, error) {
	if len(workspaces) == 0 || len(workspaces) > maxWorkspaces {
		return nil, fmt.Errorf("Bitbucket requires 1 to %d workspaces", maxWorkspaces)
	}
	seen := make(map[string]bool, len(workspaces))
	clean := make([]string, 0, len(workspaces))
	for _, workspace := range workspaces {
		workspace = strings.TrimSpace(workspace)
		if !workspacePattern.MatchString(workspace) {
			return nil, fmt.Errorf("invalid Bitbucket workspace %q", workspace)
		}
		if !seen[workspace] {
			seen[workspace] = true
			clean = append(clean, workspace)
		}
	}
	if client == nil {
		client = &http.Client{}
	}
	clone := *client
	if clone.Timeout <= 0 || clone.Timeout > 30*time.Second {
		clone.Timeout = 15 * time.Second
	}
	clone.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirect forbidden") }
	return &Connector{workspaces: clean, client: &clone, baseURL: "https://api.bitbucket.org/2.0"}, nil
}

func (c *Connector) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{
		ID: connectorID, Version: "bitbucket-v2-1.0.0", Kind: domain.ConnectorRepositoryCatalog,
		QueryLocales: []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"},
		Markets:      []string{}, SupportsMarketFilter: false,
		IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"repository_identity"},
		Replayable: true, PublicOnly: true,
	}
}

func (c *Connector) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	query := strings.TrimSpace(request.Query)
	if query == "" || len(query) > 256 {
		return domain.ConnectorDiscovery{}, domain.ErrInvalidRunInput
	}
	limit := request.Limit
	if limit <= 0 || limit > 10 {
		limit = 10
	}
	candidates := make([]domain.ProductCandidate, 0, limit)
	var failure string
	for _, workspace := range c.workspaces {
		if len(candidates) >= limit {
			break
		}
		endpoint := c.baseURL + "/repositories/" + url.PathEscape(workspace)
		params := url.Values{}
		params.Set("q", `(name ~ `+strconv.Quote(query)+` OR description ~ `+strconv.Quote(query)+`) AND is_private = false`)
		params.Set("pagelen", strconv.Itoa(limit-len(candidates)))
		params.Set("sort", "-updated_on")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
		if err != nil {
			return domain.ConnectorDiscovery{}, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Bateia-Bitbucket-Discovery-Agent/1.0")
		resp, err := c.client.Do(req)
		if err != nil {
			failure = "TRANSPORT_ERROR"
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			failure = fmt.Sprintf("HTTP_%d", resp.StatusCode)
			continue
		}
		if readErr != nil || len(body) > maxResponseBytes {
			failure = "INVALID_RESPONSE"
			continue
		}
		var page struct {
			Values []struct {
				FullName    string `json:"full_name"`
				Name        string `json:"name"`
				Description string `json:"description"`
				IsPrivate   bool   `json:"is_private"`
			} `json:"values"`
		}
		if err := json.Unmarshal(body, &page); err != nil {
			failure = "INVALID_RESPONSE"
			continue
		}
		for _, item := range page.Values {
			parts := strings.Split(item.FullName, "/")
			if item.IsPrivate || len(parts) != 2 || parts[0] != workspace || !repoPattern.MatchString(parts[1]) {
				continue
			}
			if len(candidates) >= limit {
				break
			}
			canonicalURL := "https://bitbucket.org/" + url.PathEscape(workspace) + "/" + url.PathEscape(parts[1])
			candidates = append(candidates, domain.ProductCandidate{
				Identity: domain.ProductIdentity{ConnectorID: connectorID, ExternalID: item.FullName, CanonicalURL: canonicalURL},
				Name:     item.FullName, Summary: safeText(item.Description, 500), ProductKind: "repository",
				ListingLocales: []string{request.QueryLocale}, ObservedMarkets: []string{},
				RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: item.FullName, Status: domain.RepositoryLinkVerified, Source: "bitbucket_api_metadata"}},
				Resolution:      domain.IdentityDistinct,
			})
		}
	}
	status := domain.AttemptEmptyObserved
	if len(candidates) > 0 {
		status = domain.AttemptSucceeded
	}
	return domain.ConnectorDiscovery{Status: status, Candidates: candidates, ResultCount: len(candidates), AcceptedCount: len(candidates), ErrorClass: failure}, nil
}

func safeText(value string, maximum int) string {
	if len(value) > maximum {
		value = value[:maximum]
	}
	return strings.TrimSpace(strings.ToValidUTF8(value, ""))
}
