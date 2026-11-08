/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/project_catalog.md
 * @layer L3
 * @updated 2026-11-08
 */
package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bateia/01_core/domain"
)

type catalogSearchResponse struct {
	Items []struct {
		Name          string    `json:"name"`
		HTMLURL       string    `json:"html_url"`
		Description   string    `json:"description"`
		Stargazers    int       `json:"stargazers_count"`
		Archived      bool      `json:"archived"`
		DefaultBranch string    `json:"default_branch"`
		PushedAt      time.Time `json:"pushed_at"`
		Topics        []string  `json:"topics"`
		Owner         struct {
			Login string `json:"login"`
		} `json:"owner"`
		License *struct {
			SPDX string `json:"spdx_id"`
		} `json:"license"`
	} `json:"items"`
}

func (c *Client) Capabilities() domain.ProviderCapabilities {
	return domain.ProviderCapabilities{Provider: "github", Roles: []domain.DiscoveryProviderRole{domain.RoleProjectCatalog}, SearchUnits: []string{"repository_name", "description", "topics"}, IdentityLevels: []string{"owner/repo", "commit"}}
}

func (c *Client) SearchCatalog(ctx context.Context, request domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
	result := domain.CatalogDiscovery{Provider: "github", Role: domain.RoleProjectCatalog, SearchUnit: "repository", Candidates: []domain.CatalogCandidate{}}
	started := time.Now()
	query := strings.TrimSpace(request.Query)
	result.Attempt = domain.CatalogAttempt{ExpansionID: request.ExpansionID, Query: query, Status: "FAILED"}
	if query == "" || len(query) > 4096 || request.MaxResults < 1 || request.MaxResults > 100 {
		return result, fmt.Errorf("consulta de catálogo inválida")
	}
	searchCtx := ctx
	cancel := func() {}
	if request.AttemptTimeout > 0 {
		searchCtx, cancel = context.WithTimeout(ctx, request.AttemptTimeout)
	}
	defer cancel()
	scope := " in:name,description,readme"
	if request.MetadataOnly {
		scope = " in:name,description,topics"
	}
	endpoint := fmt.Sprintf("%s/search/repositories?q=%s&per_page=%d", c.baseURL, url.QueryEscape(query+scope), request.MaxResults)
	req, err := http.NewRequestWithContext(searchCtx, "GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		result.Attempt.Duration = time.Since(started)
		if searchCtx.Err() != nil {
			result.Attempt.Status = "TIMED_OUT"
			result.Attempt.ErrorClass = "TIMEOUT"
		}
		return result, fmt.Errorf("consulta GitHub falhou: %w", err)
	}
	defer resp.Body.Close()
	result.Attempt.Duration = time.Since(started)
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		result.Attempt.Status = "RATE_LIMITED"
		result.Attempt.ErrorClass = "RATE_LIMIT"
		return result, fmt.Errorf("catálogo GitHub limitado pela quota")
	}
	if resp.StatusCode != http.StatusOK {
		result.Attempt.ErrorClass = "PROVIDER"
		return result, fmt.Errorf("GitHub search retornou status %d", resp.StatusCode)
	}
	var payload catalogSearchResponse
	if err := decodeBounded(resp, maxTreeEnvelopeBytes, &payload); err != nil {
		return result, fmt.Errorf("resposta GitHub inválida: %w", err)
	}
	if len(payload.Items) > request.MaxResults {
		return result, unsafeRemote("catálogo retornou mais itens que o limite solicitado")
	}
	for _, item := range payload.Items {
		repo := domain.RepositoryCandidate{Owner: item.Owner.Login, Repo: item.Name, URL: item.HTMLURL, Description: item.Description, Stars: item.Stargazers, Archived: item.Archived, DefaultBranch: item.DefaultBranch, UpdatedAt: item.PushedAt, ConsultadoEm: time.Now().UTC()}
		if err := validateRemoteCandidate(repo); err != nil {
			return result, err
		}
		if !request.DeferRevisionResolution {
			resolved, _ := c.FetchRepoMetadata(searchCtx, repo)
			repo = resolved
		}
		fields := matchedCatalogFields(request.Discriminators, repo, item.Topics)
		license := ""
		if item.License != nil {
			license = item.License.SPDX
		}
		result.Candidates = append(result.Candidates, domain.CatalogCandidate{Repository: repo, Topics: append([]string(nil), item.Topics...), License: license, MatchedFields: fields})
	}
	result.Attempt.Status = "SUCCEEDED"
	result.Attempt.ResultCount = len(result.Candidates)
	result.Attempt.Duration = time.Since(started)
	return result, nil
}

func matchedCatalogFields(discriminators []string, repo domain.RepositoryCandidate, topics []string) []string {
	fields := []string{}
	matches := func(value string) bool {
		v := strings.Join(strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(strings.ToLower(value))), " ")
		for _, discriminator := range discriminators {
			d := strings.Join(strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(strings.ToLower(discriminator))), " ")
			if d != "" && strings.Contains(v, d) {
				return true
			}
		}
		return false
	}
	if matches(repo.Repo) {
		fields = append(fields, "name")
	}
	if matches(repo.Description) {
		fields = append(fields, "description")
	}
	if matches(strings.Join(topics, " ")) {
		fields = append(fields, "topics")
	}
	return fields
}

var _ domain.ProjectCatalog = (*Client)(nil)
