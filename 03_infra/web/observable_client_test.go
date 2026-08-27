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
	"testing"
)

func TestObservableWebClient(t *testing.T) {
	client := NewObservableWebClient("test-provider")
	resp, err := client.SearchOne(context.Background(), domain.ObservableWebRequest{
		RequestID:  "req1",
		QueryID:    "q1",
		Query:      "tool for evaluation",
		MaxResults: 5,
	})
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if resp.Status != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", resp.Status)
	}
}
