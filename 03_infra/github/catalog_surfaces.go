/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/catalog_surfaces.md
 * @layer L3
 * @updated 2026-08-26
 */
package github

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SurfaceClient struct {
	client  *Client
	surface domain.CatalogSurfaceID
}

func NewReadmeSurface(client *Client) *SurfaceClient {
	return &SurfaceClient{client: client, surface: domain.SurfaceReadme}
}
func NewTopicsSurface(client *Client) *SurfaceClient {
	return &SurfaceClient{client: client, surface: domain.SurfaceTopics}
}
func (s *SurfaceClient) SurfaceCapabilities() domain.CatalogSurfaceCapabilities {
	return domain.CatalogSurfaceCapabilities{ID: s.surface, Provider: "github", Unit: map[domain.CatalogSurfaceID]string{domain.SurfaceReadme: "root_document", domain.SurfaceTopics: "structured_topic"}[s.surface], SupportsExactPhrase: true, ResolvesRepository: true, ResolvesRevision: true}
}
func (s *SurfaceClient) SearchSurface(ctx context.Context, r domain.SurfaceSearchRequest) (domain.SurfaceDiscovery, error) {
	if r.MaxResults < 1 || r.MaxResults > 100 || len(r.Query) > 4096 || strings.TrimSpace(r.Query) == "" {
		return domain.SurfaceDiscovery{}, unsafeRemote("consulta de superfície inválida")
	}
	if r.Surface != s.surface {
		return domain.SurfaceDiscovery{}, fmt.Errorf("superfície incompatível")
	}
	if s.surface == domain.SurfaceReadme {
		return s.searchReadme(ctx, r)
	}
	return s.searchTopics(ctx, r)
}
func (s *SurfaceClient) searchReadme(ctx context.Context, r domain.SurfaceSearchRequest) (domain.SurfaceDiscovery, error) {
	out := domain.SurfaceDiscovery{Capabilities: s.SurfaceCapabilities(), Matches: []domain.SurfaceMatch{}, Attempt: domain.CatalogAttempt{ExpansionID: r.ExpansionID, Query: r.Query, Status: "FAILED"}}
	started := time.Now()
	q := fmt.Sprintf("\"%s\" filename:README", r.Query)
	endpoint := fmt.Sprintf("%s/search/code?q=%s&per_page=%d", s.client.baseURL, url.QueryEscape(q), r.MaxResults)
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	req.Header.Set("Accept", "application/vnd.github.text-match+json")
	if s.client.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.client.token)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	out.Attempt.Duration = time.Since(started)
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		out.Attempt.Status = "RATE_LIMITED"
		out.Attempt.ErrorClass = "RATE_LIMIT"
		return out, fmt.Errorf("GitHub code search rate limited")
	}
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("GitHub code search status %d", resp.StatusCode)
	}
	var payload struct {
		Items []struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
			TextMatches []struct {
				Fragment string `json:"fragment"`
			} `json:"text_matches"`
		} `json:"items"`
	}
	if err := decodeBounded(resp, maxTreeEnvelopeBytes, &payload); err != nil {
		return out, err
	}
	if len(payload.Items) > r.MaxResults {
		return out, unsafeRemote("superfície retornou itens acima do limite solicitado")
	}
	for _, item := range payload.Items {
		if !isRootReadme(item.Path) {
			continue
		}
		parts := strings.SplitN(item.Repository.FullName, "/", 2)
		if len(parts) != 2 {
			continue
		}
		candidate := domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1]}
		if err := validateRemoteCandidate(candidate); err != nil {
			return out, err
		}
		repo, _ := s.client.FetchRepoMetadata(ctx, candidate)
		snippet := ""
		if len(item.TextMatches) > 0 {
			snippet = strings.Join(strings.Fields(item.TextMatches[0].Fragment), " ")
			if len([]rune(snippet)) > 240 {
				snippet = string([]rune(snippet)[:240])
			}
		}
		unknown := ""
		if snippet == "" {
			unknown = "snippet unavailable"
		}
		out.Matches = append(out.Matches, domain.SurfaceMatch{Surface: domain.SurfaceReadme, Provider: "github", Repository: item.Repository.FullName, Field: item.Path, Snippet: snippet, Discriminators: append([]string(nil), r.Discriminators...), Query: r.Query, ExpansionID: r.ExpansionID, CommitID: repo.CommitID, UnknownReason: unknown})
	}
	out.Attempt.Status = "SUCCEEDED"
	out.Attempt.ResultCount = len(out.Matches)
	return out, nil
}
func (s *SurfaceClient) searchTopics(ctx context.Context, r domain.SurfaceSearchRequest) (domain.SurfaceDiscovery, error) {
	out := domain.SurfaceDiscovery{Capabilities: s.SurfaceCapabilities(), Matches: []domain.SurfaceMatch{}, Attempt: domain.CatalogAttempt{ExpansionID: r.ExpansionID, Query: r.Query, Status: "FAILED"}}
	started := time.Now()
	topic := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(r.Query)), " ", "-")
	endpoint := fmt.Sprintf("%s/search/repositories?q=%s&per_page=%d", s.client.baseURL, url.QueryEscape("topic:"+topic), r.MaxResults)
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	if s.client.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.client.token)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	out.Attempt.Duration = time.Since(started)
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		out.Attempt.Status = "RATE_LIMITED"
		out.Attempt.ErrorClass = "RATE_LIMIT"
		return out, fmt.Errorf("GitHub topic search rate limited")
	}
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("GitHub topic search status %d", resp.StatusCode)
	}
	var payload catalogSearchResponse
	if err := decodeBounded(resp, maxTreeEnvelopeBytes, &payload); err != nil {
		return out, err
	}
	if len(payload.Items) > r.MaxResults {
		return out, unsafeRemote("superfície retornou itens acima do limite solicitado")
	}
	for _, item := range payload.Items {
		exact := false
		for _, t := range item.Topics {
			if strings.EqualFold(t, topic) {
				exact = true
			}
		}
		if !exact {
			continue
		}
		candidate := domain.RepositoryCandidate{Owner: item.Owner.Login, Repo: item.Name}
		if err := validateRemoteCandidate(candidate); err != nil {
			return out, err
		}
		repo, _ := s.client.FetchRepoMetadata(ctx, candidate)
		out.Matches = append(out.Matches, domain.SurfaceMatch{Surface: domain.SurfaceTopics, Provider: "github", Repository: item.Owner.Login + "/" + item.Name, Field: "topics", Snippet: topic, Discriminators: append([]string(nil), r.Discriminators...), Query: r.Query, ExpansionID: r.ExpansionID, CommitID: repo.CommitID})
	}
	out.Attempt.Status = "SUCCEEDED"
	out.Attempt.ResultCount = len(out.Matches)
	return out, nil
}
func isRootReadme(path string) bool {
	p := strings.ToLower(strings.TrimSpace(path))
	return p == "readme" || p == "readme.md" || p == "docs/readme" || p == "docs/readme.md" || p == "docs/index.md"
}

var _ domain.CatalogSurface = (*SurfaceClient)(nil)
