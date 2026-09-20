/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/market_discovery_regression_test.md
 * @layer L4
 * @updated 2026-09-19
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/itunes"
	"bateia/03_infra/marketreplay"
)

type marketDiscoveryTransportProbe struct {
	mu       sync.Mutex
	requests []*http.Request
	status   int
	body     string
	err      error
}

func (p *marketDiscoveryTransportProbe) RoundTrip(request *http.Request) (*http.Response, error) {
	p.mu.Lock()
	p.requests = append(p.requests, request.Clone(request.Context()))
	p.mu.Unlock()
	if p.err != nil {
		return nil, p.err
	}
	status := p.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"text/javascript; charset=utf-8"}},
		Body:       io.NopCloser(strings.NewReader(p.body)),
		Request:    request,
	}, nil
}

func marketDiscoveryAppleRequest(query, market string) domain.MarketDiscoveryRequest {
	return domain.MarketDiscoveryRequest{
		QueryID: "q-cross", QueryLocale: "en-US", Query: query, Market: market, Limit: 5,
	}
}

func TestMarketDiscoveryAppleAllowlistParsingUnknownAndCancellation(t *testing.T) {
	const canary = "STEP0079_SUPERSESSION_PROMPT_HEADER_PATH"
	success := `{"resultCount":1,"results":[{"wrapperType":"software","kind":"software","trackId":1001,"trackName":"Resource Ledger","trackCensoredName":"Resource Ledger","artistName":"Publisher A","sellerName":"Publisher A","description":"` + canary + `","bundleId":"example.resource.ledger","trackViewUrl":"https://apps.apple.com/jp/app/resource-ledger/id1001?uo=4#fragment","languageCodesISO2A":["EN","JA"],"genres":["Utilities"],"primaryGenreName":"Utilities","artworkUrl100":"https://evil.invalid/promo.png","averageUserRating":5,"userReview":"private"}]}`
	probe := &marketDiscoveryTransportProbe{body: success}
	connector, err := itunes.NewMarketDiscoveryConnector(&http.Client{Transport: probe})
	if err != nil {
		t.Fatalf("new Apple connector: %v", err)
	}
	discovery, err := connector.Discover(context.Background(), marketDiscoveryAppleRequest("food & inventory?callback=evil", "JP"))
	if err != nil || discovery.Status != domain.AttemptSucceeded || len(discovery.Candidates) != 1 {
		t.Fatalf("success = %+v, %v", discovery, err)
	}
	probe.mu.Lock()
	requests := append([]*http.Request(nil), probe.requests...)
	probe.mu.Unlock()
	if len(requests) != 1 {
		t.Fatalf("HTTP calls = %d, want 1", len(requests))
	}
	request := requests[0]
	if request.Method != http.MethodGet || request.URL.Scheme != "https" || request.URL.Host != "itunes.apple.com" || request.URL.Path != "/search" {
		t.Fatalf("endpoint escaped allowlist: %s %s", request.Method, request.URL)
	}
	wantQuery := map[string][]string{
		"term": {"food & inventory?callback=evil"}, "country": {"JP"}, "media": {"software"},
		"entity": {"software"}, "limit": {"5"}, "version": {"2"}, "explicit": {"No"},
	}
	values := request.URL.Query()
	for key, want := range wantQuery {
		if !reflect.DeepEqual(values[key], want) {
			t.Fatalf("query[%s] = %v, want %v (%s)", key, values[key], want, request.URL.RawQuery)
		}
	}
	if len(values) != len(wantQuery)+1 || len(values["lang"]) != 1 || (values.Get("lang") != "en_us" && values.Get("lang") != "ja_jp") || values.Get("callback") != "" {
		t.Fatalf("unexpected parameter surface: %v", values)
	}
	candidate := discovery.Candidates[0]
	if candidate.Identity.ExternalID != "1001" || candidate.Identity.CanonicalURL != "https://apps.apple.com/jp/app/resource-ledger/id1001" ||
		!reflect.DeepEqual(candidate.ObservedMarkets, []string{"JP"}) || len(candidate.ListingLocales) != 0 || len(candidate.PublisherCountries) != 0 {
		t.Fatalf("identity or dimensions inferred incorrectly: %+v", candidate)
	}
	encoded, _ := json.Marshal(discovery)
	for _, forbidden := range []string{canary, "artwork", "rating", "review", "promo.png", "callback=evil"} {
		if strings.Contains(strings.ToLower(string(encoded)), strings.ToLower(forbidden)) {
			t.Fatalf("unsafe provider/query content %q crossed discovery: %s", forbidden, encoded)
		}
	}

	for name, outcome := range map[string]struct {
		status int
		body   string
		want   domain.MarketAttemptStatus
	}{
		"strict-empty": {200, `{"resultCount":0,"results":[]}`, domain.AttemptEmptyObserved},
		"malformed":    {200, `{"resultCount":0,"results":[}`, domain.AttemptUnknownProvider},
		"rate-limit":   {429, canary, domain.AttemptUnknownProvider},
		"missing-id":   {200, `{"resultCount":1,"results":[{"trackName":"No identity"}]}`, domain.AttemptUnknownProvider},
	} {
		t.Run(name, func(t *testing.T) {
			local := &marketDiscoveryTransportProbe{status: outcome.status, body: outcome.body}
			adapter, err := itunes.NewMarketDiscoveryConnector(&http.Client{Transport: local})
			if err != nil {
				t.Fatal(err)
			}
			got, callErr := adapter.Discover(context.Background(), marketDiscoveryAppleRequest("semantic query", "US"))
			if callErr != nil || got.Status != outcome.want {
				t.Fatalf("status = %s/%v, want %s", got.Status, callErr, outcome.want)
			}
			public, _ := json.Marshal(got)
			if strings.Contains(string(public), canary) {
				t.Fatalf("provider error leaked canary: %s", public)
			}
		})
	}

	cancelProbe := &marketDiscoveryTransportProbe{body: success}
	cancelConnector, err := itunes.NewMarketDiscoveryConnector(&http.Client{Transport: cancelProbe})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, callErr := cancelConnector.Discover(ctx, marketDiscoveryAppleRequest("semantic query", "US"))
	if callErr != nil && !errors.Is(callErr, context.Canceled) {
		t.Fatalf("cancel error = %v", callErr)
	}
	if result.Status != domain.AttemptCancelled {
		t.Fatalf("cancel status = %s", result.Status)
	}
	cancelProbe.mu.Lock()
	cancelledCalls := len(cancelProbe.requests)
	cancelProbe.mu.Unlock()
	if cancelledCalls != 0 {
		t.Fatalf("cancelled connector started %d calls", cancelledCalls)
	}
}

func TestMarketDiscoveryReplayIsOfflineAndByteIdentical(t *testing.T) {
	fixturePath := filepath.Join("..", "lab", "bancada", "history", "benchmarks", "0027", "oracle", "protected", "fixtures", "apple", "success.json")
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
