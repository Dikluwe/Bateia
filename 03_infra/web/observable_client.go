/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/observable_client.md
 * @layer L3
 * @updated 2026-08-26
 */
package web

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ObservableWebClient struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	Token      string
}

func NewObservableWebClient(client *http.Client, tokens ...string) *ObservableWebClient {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	token := ""
	if len(tokens) > 0 {
		token = tokens[0]
	}
	return &ObservableWebClient{BaseURL: "https://api.github.com/search/repositories", HTTPClient: client, UserAgent: "bateia/step-0018", Token: token}
}
func (c *ObservableWebClient) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "github-repository-search", SupportsLocale: false, PreservesRank: true, MaxResultsLimit: 100}
}

func (c *ObservableWebClient) SearchOne(ctx context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	if strings.TrimSpace(req.Query) == "" {
		return domain.RawWebResponse{}, fmt.Errorf("query vazia")
	}
	limit := req.MaxResults
	if limit < 1 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return domain.RawWebResponse{}, err
	}
	q := u.Query()
	q.Set("q", req.Query)
	q.Set("per_page", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return domain.RawWebResponse{}, err
	}
	hreq.Header.Set("Accept", "application/vnd.github+json")
	hreq.Header.Set("User-Agent", c.UserAgent)
	if c.Token != "" {
		hreq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	started := time.Now()
	hresp, err := c.HTTPClient.Do(hreq)
	elapsed := time.Since(started).Milliseconds()
	now := time.Now().UTC().Format(time.RFC3339)
	base := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "github-repository-search", RetrievedAt: now, Receipt: domain.WebReceipt{SanitizedURL: c.BaseURL, ResponseTime: elapsed}}
	if err != nil {
		base.Status = "UNKNOWN"
		base = usecase.SealWebResponse(req, base)
		return base, err
	}
	defer hresp.Body.Close()
	base.HTTPStatus = hresp.StatusCode
	if hresp.StatusCode != http.StatusOK {
		base.Status = "UNKNOWN"
		base = usecase.SealWebResponse(req, base)
		return base, fmt.Errorf("GitHub search HTTP %d", hresp.StatusCode)
	}
	var payload struct {
		Items []struct {
			FullName, HTMLURL, Description string `json:"-"`
		} `json:"items"`
	}
	var raw struct {
		Items []struct {
			FullName    string `json:"full_name"`
			HTMLURL     string `json:"html_url"`
			Description string `json:"description"`
		} `json:"items"`
	}
	if err := json.NewDecoder(hresp.Body).Decode(&raw); err != nil {
		return domain.RawWebResponse{}, err
	}
	_ = payload
	base.Status = "SUCCESS"
	for i, item := range raw.Items {
		snippet := strings.TrimSpace(item.Description)
		if snippet == "" {
			snippet = "Repository description unavailable"
		}
		base.Results = append(base.Results, domain.RawWebResult{Rank: i + 1, URL: item.HTMLURL, CanonicalURL: item.HTMLURL, Title: item.FullName, Snippet: snippet})
	}
	base = usecase.SealWebResponse(req, base)
	return base, nil
}
