/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/observable_client.md
 * @layer L3
 * @updated 2026-08-31
 */
package web

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/json"
	"errors"
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
	return domain.WebSearchCapabilities{Provider: "github-repository-search", SupportsLocale: false, PreservesRank: true, MaxResultsLimit: 100, Authenticated: c.Token != "", PublicOnly: true}
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
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			base.Status = "UNKNOWN_TIMEOUT"
		} else {
			base.Status = "UNKNOWN_PROVIDER_RESPONSE"
		}
		base = usecase.SealWebResponse(req, base)
		return base, err
	}
	defer hresp.Body.Close()
	base.HTTPStatus = hresp.StatusCode
	base.Receipt.RateLimit = readRateLimit(hresp.Header)
	if hresp.StatusCode != http.StatusOK {
		if hresp.StatusCode == http.StatusForbidden || hresp.StatusCode == http.StatusTooManyRequests {
			base.Status = "UNKNOWN_RATE_LIMIT"
		} else {
			base.Status = "UNKNOWN_PROVIDER_RESPONSE"
		}
		base = usecase.SealWebResponse(req, base)
		return base, fmt.Errorf("github repository search indisponível (HTTP %d)", hresp.StatusCode)
	}
	var raw struct {
		Items []struct {
			FullName    string `json:"full_name"`
			HTMLURL     string `json:"html_url"`
			Description string `json:"description"`
			Private     *bool  `json:"private"`
			Visibility  string `json:"visibility"`
		} `json:"items"`
	}
	if err := json.NewDecoder(hresp.Body).Decode(&raw); err != nil {
		base.Status = "UNKNOWN_PROVIDER_RESPONSE"
		base = usecase.SealWebResponse(req, base)
		return base, fmt.Errorf("resposta inválida do provedor")
	}
	if len(raw.Items) == 0 {
		base.Status = "EMPTY"
		base = usecase.SealWebResponse(req, base)
		return base, nil
	}
	for _, item := range raw.Items {
		if item.Private != nil && *item.Private {
			base.Status = "REJECTED_PRIVATE"
			base.Results = nil
			base = usecase.SealWebResponse(req, base)
			return base, fmt.Errorf("resultado não público rejeitado")
		}
		if item.Private == nil && !strings.EqualFold(item.Visibility, "public") {
			base.Status = "UNKNOWN_PROVIDER_RESPONSE"
			base.Results = nil
			base = usecase.SealWebResponse(req, base)
			return base, fmt.Errorf("visibilidade do resultado não comprovada")
		}
		if item.Visibility != "" && !strings.EqualFold(item.Visibility, "public") {
			base.Status = "REJECTED_PRIVATE"
			base.Results = nil
			base = usecase.SealWebResponse(req, base)
			return base, fmt.Errorf("resultado não público rejeitado")
		}
	}
	base.Status = "SUCCESS"
	for i, item := range raw.Items {
		snippet := strings.TrimSpace(item.Description)
		if snippet == "" {
			snippet = "Repository description unavailable"
		}
		base.Results = append(base.Results, domain.RawWebResult{Rank: i + 1, URL: item.HTMLURL, CanonicalURL: item.HTMLURL, Title: item.FullName, Snippet: snippet, Visibility: "public"})
	}
	base = usecase.SealWebResponse(req, base)
	if containsSensitiveWebValue(base, c.Token) {
		return domain.RawWebResponse{}, fmt.Errorf("resposta rejeitada por política de segurança")
	}
	return base, nil
}
