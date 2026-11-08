/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/fdroid/market_discovery_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package fdroid_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/fdroid"
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

func TestStep0034FDroidConnectorExtractsFOSSLinkage(t *testing.T) {
	mockClient := mockHTTPClient(func(req *http.Request) *http.Response {
		jsonResp := `[
			{
				"packageName": "com.github.food.app",
				"name": "Open Food Plan",
				"summary": "Self hosted meal planner",
				"sourceCode": "https://github.com/food-org/open-food-app",
				"license": "GPL-3.0-only",
				"webUrl": "https://f-droid.org/packages/com.github.food.app/",
				"issueTracker": "https://github.com/food-org/open-food-app/issues"
			}
		]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonResp)),
			Header:     make(http.Header),
		}
	})

	connector, err := fdroid.NewMarketDiscoveryConnector(mockClient)
	if err != nil {
		t.Fatalf("erro ao instanciar conector fdroid: %v", err)
	}

	caps := connector.Capabilities()
	if caps.ID != "fdroid-foss-apps" {
		t.Fatalf("ID esperado fdroid-foss-apps, obtido: %s", caps.ID)
	}
	if caps.SupportsMarketFilter {
		t.Fatalf("F-Droid é catálogo global FOSS e não filtra por storefront nacional")
	}

	req := domain.MarketDiscoveryRequest{
		Query:       "meal planner",
		QueryLocale: "en-US",
	}

	discovery, err := connector.Discover(context.Background(), req)
	if err != nil {
		t.Fatalf("erro ao executar Discover: %v", err)
	}
	if discovery.Status != domain.AttemptSucceeded {
		t.Fatalf("status esperado SUCCEEDED, obtido: %s", discovery.Status)
	}
	if len(discovery.Candidates) != 1 {
		t.Fatalf("esperado 1 candidato, obtidos: %d", len(discovery.Candidates))
	}
	cand := discovery.Candidates[0]
	if cand.FOSSLink == nil {
		t.Fatalf("FOSSLink não deve ser nulo para aplicativo F-Droid com código upstream")
	}
	if cand.FOSSLink.UpstreamURL != "https://github.com/food-org/open-food-app" {
		t.Fatalf("UpstreamURL incorreta: %s", cand.FOSSLink.UpstreamURL)
	}
	if cand.FOSSLink.LicenseSPDX != "GPL-3.0-only" {
		t.Fatalf("LicenseSPDX incorreta: %s", cand.FOSSLink.LicenseSPDX)
	}
	if len(cand.RepositoryLinks) != 1 || cand.RepositoryLinks[0].OwnerRepo != "food-org/open-food-app" {
		t.Fatalf("RepositoryLink causal incorreto: %+v", cand.RepositoryLinks)
	}
}
