/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/observable_client_test.md
 * @layer L3
 * @updated 2026-08-26
 */
package web

import (
	"bateia/01_core/domain"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObservableWebClientPreservesPublicResult(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "find tools" {
			t.Errorf("query=%q", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"full_name":"acme/tool","html_url":"https://github.com/acme/tool","description":"Finds tools"}]}`))
	}))
	defer s.Close()
	c := NewObservableWebClient(s.Client())
	c.BaseURL = s.URL
	r, e := c.SearchOne(context.Background(), domain.ObservableWebRequest{RequestID: "r", QueryID: "q", Query: "find tools", MaxResults: 5})
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != "SUCCESS" || len(r.Results) != 1 || r.ResponseHash == "" || r.RequestHash == "" {
		t.Fatalf("%+v", r)
	}
}

func TestObservableWebClientDoesNotFabricateOnRateLimit(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusForbidden) }))
	defer s.Close()
	c := NewObservableWebClient(s.Client())
	c.BaseURL = s.URL
	r, e := c.SearchOne(context.Background(), domain.ObservableWebRequest{RequestID: "r", QueryID: "q", Query: "x"})
	if e == nil || r.Status != "UNKNOWN" || len(r.Results) != 0 {
		t.Fatalf("%+v %v", r, e)
	}
}
