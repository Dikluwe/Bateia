/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/bitbucket/connector_test.md
 * @layer L3
 * @updated 2026-10-06
 */

package bitbucket

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestDiscoverPublicRepositoriesAcrossConfiguredWorkspaces(t *testing.T) {
	called := []string{}
	client := &http.Client{Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
		called = append(called, req.URL.Path)
		if req.URL.Host != "api.bitbucket.org" || !strings.Contains(req.URL.Query().Get("q"), `is_private = false`) || !strings.Contains(req.URL.Query().Get("q"), `name ~ "recipe"`) {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		body := `{"values":[]}`
		if strings.HasSuffix(req.URL.Path, "/second") {
			body = `{"values":[{"full_name":"second/cookbook","name":"cookbook","description":"Recipes","is_private":false},{"full_name":"second/.well-known","is_private":false},{"full_name":"second/secret","is_private":true},{"full_name":"other/foreign","is_private":false}]}`
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
	connector, err := NewConnector([]string{"first", "second"}, client)
	if err != nil {
		t.Fatal(err)
	}
	result, err := connector.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "recipe", QueryLocale: "pt-BR", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(called) != 2 || result.Status != domain.AttemptSucceeded || len(result.Candidates) != 2 {
		t.Fatalf("calls=%v result=%+v", called, result)
	}
	item := result.Candidates[0]
	if item.Identity.CanonicalURL != "https://bitbucket.org/second/cookbook" || item.Identity.ExternalID != "second/cookbook" {
		t.Fatalf("candidate=%+v", item)
	}
}

func TestConnectorRejectsInvalidWorkspace(t *testing.T) {
	for _, value := range []string{"", "../other", "https://bitbucket.org/team", "team/other", "Uppercase"} {
		if _, err := NewConnector([]string{value}, nil); err == nil {
			t.Fatalf("accepted %q", value)
		}
	}
}
