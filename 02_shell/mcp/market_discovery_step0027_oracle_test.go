/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/market_discovery_step0027_oracle_test.md
 * @layer L2
 * @updated 2026-09-03
 */

package mcp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type step0027DurableStub struct{}

func (step0027DurableStub) Start(context.Context, string, domain.RecommendationRequest) (domain.RecommendationStartReceipt, error) {
	return domain.RecommendationStartReceipt{}, errors.New("unused")
}
func (step0027DurableStub) GetRun(context.Context, domain.RunID) (domain.RecommendationRunView, error) {
	return domain.RecommendationRunView{}, errors.New("unused")
}
func (step0027DurableStub) GetArtifact(context.Context, domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	return domain.ArtifactPage{}, errors.New("unused")
}
func (step0027DurableStub) Cancel(context.Context, domain.RunID) (domain.RecommendationRunView, error) {
	return domain.RecommendationRunView{}, errors.New("unused")
}

type step0027ConnectorStub struct {
	discoverCalls int
	cap           domain.DiscoveryConnectorCapabilities
}

func (c *step0027ConnectorStub) Capabilities() domain.DiscoveryConnectorCapabilities {
	if c.cap.ID != "" {
		return c.cap
	}
	return domain.DiscoveryConnectorCapabilities{
		ID: "apple-itunes-search-software", Version: "itunes-search-doc-1.0.3-response-v2", Kind: domain.ConnectorApplicationCatalog,
		QueryLocales: []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"},
		Markets:      []string{"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US"}, SupportsMarketFilter: true,
		IdentityLevels:   []string{"provider_product_id", "canonical_url"},
		EvidenceSurfaces: []string{"name", "summary", "publisher", "provider_declared_language_codes", "storefront_observation"},
		Replayable:       true, PublicOnly: true,
	}
}
func (c *step0027ConnectorStub) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	c.discoverCalls++
	return domain.ConnectorDiscovery{}, nil
}

type step0027StarterStub struct {
	calls int
	last  domain.MarketDiscoveryStartRequest
}

func (s *step0027StarterStub) StartMarketDiscovery(_ context.Context, key string, request domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartReceipt, error) {
	s.calls++
	s.last = request
	if key != "market-task-0027" {
		return domain.MarketDiscoveryStartReceipt{}, domain.ErrInvalidRunInput
	}
	return domain.MarketDiscoveryStartReceipt{SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef", RunKind: domain.RunKindMarketDiscovery, Lifecycle: domain.RunQueued, Version: 1}, nil
}

func step0027RunMCP(t *testing.T, withMarket bool, requests string, connector *step0027ConnectorStub, starter *step0027StarterStub) []JSONRPCResponse {
	t.Helper()
	var output bytes.Buffer
	server := NewMCPServer(nil, strings.NewReader(requests), &output).WithDurableRuns(step0027DurableStub{})
	if withMarket {
		server = server.WithMarketDiscovery([]domain.MarketDiscoveryConnector{connector}, starter)
	}
	if err := server.Serve(context.Background()); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	responses := make([]JSONRPCResponse, 0, len(lines))
	for _, line := range lines {
		var response JSONRPCResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("response JSON %q: %v", line, err)
		}
		responses = append(responses, response)
	}
	return responses
}

func step0027Tools(t *testing.T, response JSONRPCResponse) []map[string]any {
	t.Helper()
	encoded, err := json.Marshal(response.Result)
	if err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Tools []map[string]any `json:"tools"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope.Tools
}

func TestStep0027OracleToolsOneThroughSevenAreCertifiedAndEightNineAreExact(t *testing.T) {
	for path, want := range map[string]string{
		"server.go":       "e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42",
		"durable_runs.go": "436989b07761406908297d733e8c07cb890c12aae92613cc541a232d8e1ec29e",
		"../../00_nucleo/benchmark/0029/historical-hash-invalidation-o2-to-o3.json": "98c36700e72c434a2443ce4aa7e64c5b5c9b599f0eb0b5451cf07c9af46639ca",
		"../../00_nucleo/benchmark/0029/historical-hash-successor-o3.json":          "cf5da34d390b710425f71a5545d46978245db3581d5f2caacf762980238c4156",
		"../../00_nucleo/benchmark/0031/historical-contract-supersession.json":      "65598792dd329561e7d8a21bd1cf0bbd22c19d00876470361ac2ce6e4d6aae67",
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(content)
		if got := hex.EncodeToString(digest[:]); got != want {
			t.Fatalf("certified basal surface %s changed: %s", path, got)
		}
	}
	connector := &step0027ConnectorStub{}
	starter := &step0027StarterStub{}
	request := "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}\n"
	basal := step0027Tools(t, step0027RunMCP(t, false, request, connector, starter)[0])
	got := step0027Tools(t, step0027RunMCP(t, true, request, connector, starter)[0])
	if len(basal) != 7 || len(got) != 9 {
		t.Fatalf("tool refinement is not exact: basal=%d candidate=%d", len(basal), len(got))
	}
	for _, index := range []int{0, 1, 2, 3, 4, 6} {
		if !reflect.DeepEqual(got[index], basal[index]) {
			t.Fatalf("tool %d changed outside the additive contract", index+1)
		}
	}
	cloneTool := func(tool map[string]any) map[string]any {
		encoded, err := json.Marshal(tool)
		if err != nil {
			t.Fatal(err)
		}
		var cloned map[string]any
		if err := json.Unmarshal(encoded, &cloned); err != nil {
			t.Fatal(err)
		}
		return cloned
	}
	basalArtifactTool := cloneTool(basal[5])
	candidateArtifactTool := cloneTool(got[5])
	basalArtifactEnum := basalArtifactTool["inputSchema"].(map[string]any)["properties"].(map[string]any)["artifact"].(map[string]any)["enum"].([]any)
	candidateArtifactSchema := candidateArtifactTool["inputSchema"].(map[string]any)["properties"].(map[string]any)["artifact"].(map[string]any)
	candidateArtifactEnum := candidateArtifactSchema["enum"].([]any)
	wantArtifactEnum := append(append([]any{}, basalArtifactEnum...),
		"connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps")
	if !reflect.DeepEqual(candidateArtifactEnum, wantArtifactEnum) {
		t.Fatalf("tool 6 artifact enum = %v, want basal plus five ordered suffixes", candidateArtifactEnum)
	}
	candidateArtifactSchema["enum"] = append([]any{}, basalArtifactEnum...)
	if !reflect.DeepEqual(candidateArtifactTool, basalArtifactTool) {
		t.Fatal("tool 6 changed outside its five additive artifact enum suffixes")
	}
	wantNames := []string{"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run", "bateia_list_connectors", "bateia_start_market_discovery"}
	for i, want := range wantNames {
		if got[i]["name"] != want {
			t.Fatalf("tool[%d] = %v, want %s", i, got[i]["name"], want)
		}
	}
	wantListSchema := map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{}}
	if !reflect.DeepEqual(got[7]["inputSchema"], wantListSchema) {
		t.Fatalf("tool 8 schema = %#v", got[7]["inputSchema"])
	}
	schema := got[8]["inputSchema"].(map[string]any)
	if schema["type"] != "object" || schema["additionalProperties"] != false || !reflect.DeepEqual(schema["required"], []any{"idempotency_key", "idea", "query_proposals", "markets", "connector_ids", "budget"}) {
		t.Fatalf("tool 9 top-level schema changed: %#v", schema)
	}
	properties := schema["properties"].(map[string]any)
	locales := properties["query_proposals"].(map[string]any)["items"].(map[string]any)["properties"].(map[string]any)["locale"].(map[string]any)["enum"]
	markets := properties["markets"].(map[string]any)["items"].(map[string]any)["enum"]
	connectors := properties["connector_ids"].(map[string]any)["items"].(map[string]any)["enum"]
	if !reflect.DeepEqual(locales, []any{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"}) ||
		!reflect.DeepEqual(markets, []any{"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US"}) ||
		!reflect.DeepEqual(connectors, []any{"apple-itunes-search-software", "github-repository-catalog", "sourcegraph-repository-catalog"}) {
		t.Fatalf("tool 9 frozen enums changed: locales=%v markets=%v connectors=%v", locales, markets, connectors)
	}
	if connector.discoverCalls != 0 || starter.calls != 0 {
		t.Fatalf("tools/list had effects: discover=%d start=%d", connector.discoverCalls, starter.calls)
	}

	repoCapabilities := func(id string) domain.DiscoveryConnectorCapabilities {
		return domain.DiscoveryConnectorCapabilities{
			ID: id, Version: "baseline-21c48e357180", Kind: domain.ConnectorRepositoryCatalog,
			QueryLocales: []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"},
			Markets:      []string{}, SupportsMarketFilter: false,
			IdentityLevels: []string{"owner_repo", "commit"}, EvidenceSurfaces: []string{"repository_identity", "immutable_commit", "code_occurrence"},
			Replayable: true, PublicOnly: true,
		}
	}
	github := &step0027ConnectorStub{cap: repoCapabilities("github-repository-catalog")}
	sourcegraph := &step0027ConnectorStub{cap: repoCapabilities("sourcegraph-repository-catalog")}
	var output bytes.Buffer
	server := NewMCPServer(nil, strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":8,\"method\":\"tools/call\",\"params\":{\"name\":\"bateia_list_connectors\",\"arguments\":{}}}\n"), &output).
		WithDurableRuns(step0027DurableStub{}).
		WithMarketDiscovery([]domain.MarketDiscoveryConnector{sourcegraph, connector, github}, starter)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	var listed JSONRPCResponse
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &listed); err != nil || listed.Error != nil {
		t.Fatalf("list connectors response: %v %+v", err, listed.Error)
	}
	listedBytes, _ := json.Marshal(listed.Result)
	for _, exact := range []string{`\"id\":\"apple-itunes-search-software\"`, `\"id\":\"github-repository-catalog\"`, `\"id\":\"sourcegraph-repository-catalog\"`, `\"markets\":[]`, `\"public_only\":true`} {
		if !strings.Contains(string(listedBytes), exact) {
			t.Fatalf("connector list lacks %s: %s", exact, listedBytes)
		}
	}
	if strings.Index(string(listedBytes), "apple-itunes") > strings.Index(string(listedBytes), "github-repository") || strings.Index(string(listedBytes), "github-repository") > strings.Index(string(listedBytes), "sourcegraph-repository") {
		t.Fatalf("connector list is not canonical: %s", listedBytes)
	}
	if connector.discoverCalls+github.discoverCalls+sourcegraph.discoverCalls != 0 {
		t.Fatal("bateia_list_connectors performed discovery/network I/O")
	}
}

func TestStep0027OracleStartBoundaryIsClosedAndReturnsExactQueuedEnvelope(t *testing.T) {
	connector := &step0027ConnectorStub{}
	starter := &step0027StarterStub{}
	valid := `{"idempotency_key":"market-task-0027","idea":{"problem":"reduzir desperdício alimentar doméstico","capabilities":["inventário"]},"query_proposals":[{"locale":"en-US","query":"household food inventory","source_text":"inventário doméstico","derivation":"EXTERNAL_AGENT_PROPOSAL","rationale":"controle cruzado"}],"markets":["JP"],"connector_ids":["apple-itunes-search-software"],"budget":{"max_queries":1,"max_results_per_query":5,"max_total_candidates":5}}`
	requests := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_start_market_discovery","arguments":` + valid + `}}` + "\n"
	responses := step0027RunMCP(t, true, requests, connector, starter)
	if len(responses) != 1 || responses[0].Error != nil || starter.calls != 1 || connector.discoverCalls != 0 {
		t.Fatalf("valid start result=%+v start=%d discover=%d", responses, starter.calls, connector.discoverCalls)
	}
	encoded, _ := json.Marshal(responses[0].Result)
	if !strings.Contains(string(encoded), `\"run_kind\":\"MARKET_DISCOVERY\"`) || !strings.Contains(string(encoded), `\"lifecycle\":\"QUEUED\"`) {
		t.Fatalf("queued envelope changed: %s", encoded)
	}

	invalids := []string{
		strings.Replace(valid, `"budget":{`, `"unknown":true,"budget":{`, 1),
		valid + `{}`,
		strings.Replace(valid, `"query":"household food inventory"`, `"query":"x\u0000y"`, 1),
		strings.Replace(valid, `"markets":["JP"]`, `"markets":["JP","JP"]`, 1),
		strings.Replace(valid, `"connector_ids":["apple-itunes-search-software"]`, `"connector_ids":["unregistered"]`, 1),
	}
	for i, invalid := range invalids {
		before := starter.calls
		request := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_start_market_discovery","arguments":` + invalid + `}}` + "\n"
		response := step0027RunMCP(t, true, request, connector, starter)[0]
		if response.Error == nil || response.Error.Code != -32602 || starter.calls != before || connector.discoverCalls != 0 {
			t.Fatalf("invalid[%d] crossed boundary: %+v calls=%d", i, response, starter.calls)
		}
	}
}
