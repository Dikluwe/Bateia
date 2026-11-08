/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/marketreplay/replay_test.md
 * @layer L3
 * @updated 2026-10-28
 */
package marketreplay_test

import (
	"bateia/01_core/domain"
	"bateia/03_infra/marketreplay"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func marketDiscoveryAppleRequest(query, market string) domain.MarketDiscoveryRequest {
	return domain.MarketDiscoveryRequest{
		QueryID: "q-cross", QueryLocale: "en-US", Query: query, Market: market, Limit: 5,
	}
}

func TestMarketDiscoveryReplayIsOfflineAndByteIdentical(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "lab", "bancada", "history", "benchmarks", "0027", "oracle", "protected", "fixtures", "apple", "success.json")
	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read frozen fixture: %v", err)
	}
	replay, err := marketreplay.NewConnector(fixture)
	if err != nil {
		t.Fatalf("new replay: %v", err)
	}
	first, err := replay.Discover(context.Background(), marketDiscoveryAppleRequest("household resource tracking", "BR"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := replay.Discover(context.Background(), marketDiscoveryAppleRequest("household resource tracking", "BR"))
	if err != nil {
		t.Fatal(err)
	}
	firstBytes, _ := json.Marshal(first)
	secondBytes, _ := json.Marshal(second)
	if !reflect.DeepEqual(firstBytes, secondBytes) {
		t.Fatalf("replay differs:\n%s\n%s", firstBytes, secondBytes)
	}

	corrupt := append([]byte(nil), fixture...)
	corrupt[len(corrupt)/2] ^= 1
	if broken, createErr := marketreplay.NewConnector(corrupt); createErr == nil {
		got, discoverErr := broken.Discover(context.Background(), marketDiscoveryAppleRequest("household resource tracking", "BR"))
		if discoverErr == nil && got.Status != domain.AttemptIntegrityError {
			t.Fatalf("corrupt replay = %+v, %v", got, discoverErr)
		}
	}
}
