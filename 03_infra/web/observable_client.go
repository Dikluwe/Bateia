/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/observable_client.md
 * @layer L3
 * @updated 2026-08-26
 */
package web

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
	"time"
)

type ObservableWebClient struct {
	ProviderName string
}

func NewObservableWebClient(provider string) *ObservableWebClient {
	if provider == "" {
		provider = "web-mock"
	}
	return &ObservableWebClient{ProviderName: provider}
}

func (c *ObservableWebClient) Capabilities(ctx context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{
		Provider:        c.ProviderName,
		SupportsLocale:  true,
		PreservesRank:   true,
		MaxResultsLimit: 10,
	}
}

func (c *ObservableWebClient) SearchOne(ctx context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	if req.Query == "" {
		return domain.RawWebResponse{}, fmt.Errorf("query vazia")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	results := []domain.RawWebResult{
		{
			Rank:         1,
			URL:          fmt.Sprintf("https://github.com/example/%s-tool", req.QueryID),
			CanonicalURL: fmt.Sprintf("https://github.com/example/%s-tool", req.QueryID),
			Title:        fmt.Sprintf("Tool for %s", req.Query),
			Snippet:      fmt.Sprintf("Open source tool implementing %s", req.Query),
			IsAd:         false,
		},
	}

	return domain.RawWebResponse{
		RequestID:   req.RequestID,
		QueryID:     req.QueryID,
		Provider:    c.ProviderName,
		RetrievedAt: now,
		Status:      "SUCCESS",
		HTTPStatus:  200,
		Results:     results,
		Receipt: domain.WebReceipt{
			SanitizedURL: "https://api.search.provider/v1/search",
			ResponseTime: 42,
		},
	}, nil
}
