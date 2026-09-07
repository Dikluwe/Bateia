/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0034_android_and_forges_test.md
 * @layer L4
 * @updated 2026-09-06
 */

package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/fdroid"
	"bateia/03_infra/forge"
	"bateia/03_infra/playstore"
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

func TestStep0034FederatedConnectorsComposition(t *testing.T) {
	connectors, err := federatedMarketDiscoveryConnectors()
	if err != nil {
		t.Fatalf("erro ao compor conectores federados: %v", err)
	}
	if len(connectors) != 8 {
		t.Fatalf("esperados 8 conectores federados, obtidos: %d", len(connectors))
	}

	foundIDs := map[string]bool{}
	for _, c := range connectors {
		foundIDs[c.Capabilities().ID] = true
	}

	expectedIDs := []string{
		"apple-itunes-search-software",
		"google-play-store-apps",
		"fdroid-foss-apps",
		"github-repository-catalog",
		"sourcegraph-repository-catalog",
		"gitee-repository-catalog",
		"codeberg-repository-catalog",
		"gitlab-repository-catalog",
	}

	for _, wantID := range expectedIDs {
		if !foundIDs[wantID] {
			t.Fatalf("conector esperado %s não encontrado no conjunto federado", wantID)
		}
	}
}
