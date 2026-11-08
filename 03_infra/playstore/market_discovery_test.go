/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/playstore/market_discovery_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package playstore_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/playstore"
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

func TestStep0034AndroidPlayStoreConnectorCapabilitiesAndDiscovery(t *testing.T) {
	mockClient := mockHTTPClient(func(req *http.Request) *http.Response {
		html := `<!DOCTYPE html><html><body>
		<a href="/store/apps/details?id=com.example.recipebook">Recipe Book</a>
		<a href="/store/apps/details?id=org.food.cookingdiary">Cooking Diary</a>
		</body></html>`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(html)),
			Header:     make(http.Header),
		}
	})

	connector, err := playstore.NewMarketDiscoveryConnector(mockClient)
	if err != nil {
		t.Fatalf("erro ao instanciar conector playstore: %v", err)
	}

	caps := connector.Capabilities()
	if caps.ID != "google-play-store-apps" {
		t.Fatalf("ID esperado google-play-store-apps, obtido: %s", caps.ID)
	}
	if caps.Kind != domain.ConnectorApplicationCatalog {
		t.Fatalf("Kind esperado APPLICATION_CATALOG, obtido: %s", caps.Kind)
	}
	if !caps.SupportsMarketFilter {
		t.Fatalf("Play Store deve suportar filtro por mercado (gl)")
	}

	req := domain.MarketDiscoveryRequest{
		Query:       "recipes cooking",
		QueryLocale: "pt-BR",
		Market:      "BR",
	}

	discovery, err := connector.Discover(context.Background(), req)
	if err != nil {
		t.Fatalf("erro ao executar Discover: %v", err)
	}
	if discovery.Status != domain.AttemptSucceeded {
		t.Fatalf("status esperado SUCCEEDED, obtido: %s", discovery.Status)
	}
	if len(discovery.Candidates) != 2 {
		t.Fatalf("esperados 2 candidatos, obtidos: %d", len(discovery.Candidates))
	}
	c1 := discovery.Candidates[0]
	if c1.Identity.ExternalID != "com.example.recipebook" {
		t.Fatalf("ExternalID incorreto: %s", c1.Identity.ExternalID)
	}
	if c1.Identity.CanonicalURL != "https://play.google.com/store/apps/details?id=com.example.recipebook" {
		t.Fatalf("CanonicalURL incorreta: %s", c1.Identity.CanonicalURL)
	}
}
