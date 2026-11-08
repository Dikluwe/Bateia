/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/federated_regression_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package forge_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/forge"
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func mockHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func TestStep0034FederatedForgesSearch(t *testing.T) {
	mockClient := mockHTTPClient(func(req *http.Request) *http.Response {
		jsonResp := `[
			{
				"full_name": "codeberg-user/culinary-ai",
				"name": "culinary-ai",
				"html_url": "https://codeberg.org/codeberg-user/culinary-ai",
				"description": "Smart recipe recommendations engine",
				"default_branch": "main"
			}
		]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonResp)),
			Header:     make(http.Header),
		}
	})

	gitee := forge.NewGiteeConnector(mockClient)
	codeberg := forge.NewCodebergConnector(mockClient)
	gitlab := forge.NewGitLabConnector(mockClient)

	for _, c := range []domain.MarketDiscoveryConnector{gitee, codeberg, gitlab} {
		caps := c.Capabilities()
		if caps.Kind != domain.ConnectorRepositoryCatalog {
			t.Fatalf("conector de forja %s deve ser REPOSITORY_CATALOG", caps.ID)
		}
		if caps.SupportsMarketFilter {
			t.Fatalf("forjas de código aberto não devem requerer filtro por país ISO 3166")
		}
	}

	req := domain.MarketDiscoveryRequest{
		Query:       "culinary-ai",
		QueryLocale: "en-US",
	}

	discovery, err := codeberg.Discover(context.Background(), req)
	if err != nil {
		t.Fatalf("erro ao executar Discover no Codeberg: %v", err)
	}
	if discovery.Status != domain.AttemptSucceeded {
		t.Fatalf("status esperado SUCCEEDED, obtido: %s", discovery.Status)
	}
	if len(discovery.Candidates) != 1 {
		t.Fatalf("esperado 1 candidato, obtido: %d", len(discovery.Candidates))
	}
	cand := discovery.Candidates[0]
	if cand.ProductKind != "repository" {
		t.Fatalf("ProductKind esperado repository, obtido: %s", cand.ProductKind)
	}
	if cand.Identity.ConnectorID != "codeberg-repository-catalog" {
		t.Fatalf("ConnectorID incorreto: %s", cand.Identity.ConnectorID)
	}
}
