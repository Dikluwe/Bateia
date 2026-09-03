/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/durable_runs_step0026_r2_test.md
 * @layer L2
 * @updated 2026-09-01
 */
package mcp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

const step0026LegacySearchTool = `{"name":"bateia_search","description":"Executa busca arquitetural no GitHub/Sourcegraph reduzindo repositórios via funil de 3 estágios","inputSchema":{"type":"object","properties":{"query":{"type":"string","description":"Query de busca (ex: 'hexagonal architecture' ou nome do pacote)"},"language":{"type":"string","description":"Linguagem alvo (go, typescript, python, rust, php, etc.)","default":"go"},"depth":{"type":"integer","description":"Profundidade da busca recursiva em grafo de dependências","default":1}},"required":["query"]}}`
const step0026LegacyBenchmarkTool = `{"name":"bateia_benchmark","description":"Executa comparação arquitetural direta de uma lista de repositórios conhecidos","inputSchema":{"type":"object","properties":{"repos":{"type":"string","description":"Repositórios separados por vírgula no formato owner/repo (ex: 'traggo/server,solidtime-io/solidtime')"},"query":{"type":"string","description":"Título ou tema da comparação de benchmark","default":"Benchmark de Repositórios"}},"required":["repos"]}}`

func TestStep0026R2LegacyGoldenIsManifestBoundAndAllThreeToolsReallyExecute(t *testing.T) {
	for path, want := range map[string]string{
		"../../00_nucleo/benchmark/0025/certificate.json": "bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d",
		"recommendation_step0025_contract_test.go":        "cc88dfe582c5e61688468c89b94a173e0f23aea4d0e4105ba2ea6749c595dea6",
		"server_test.go": "f7d15d6e29b66c6591e91d832e5c0e1c155641ca345a150610470133a3dd098d",
		"../../00_nucleo/prompts/02_shell/mcp/server.md":                                "82fbdf503a125a50da615fb85534440200593bebfca43b86ec5e7629be839f1e",
		"../../00_nucleo/prompts/02_shell/mcp/recommendation_step0025_contract_test.md": "04bfd375118c94df889ebed6aaf01bc105e2df806b0e30b93c6b506682335767",
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(data)
		if got := hex.EncodeToString(digest[:]); got != want {
			t.Fatalf("immutable Step 0025 artifact changed: %s got=%s want=%s", path, got, want)
		}
	}

	recommendFixture := domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunComplete,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionBuild, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{}, Warnings: []string{},
	}
	recommender := &step0025RecommenderSpy{result: recommendFixture}
	durable := &step0026DurableSpy{}
	lines := []string{
		`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":40,"method":"tools/call","params":{"name":"bateia_search","arguments":{"query":"legacy success"}}}`,
		`{"jsonrpc":"2.0","id":41,"method":"tools/call","params":{"name":"bateia_benchmark","arguments":{"repos":"acme/example","query":"legacy benchmark"}}}`,
		`{"jsonrpc":"2.0","id":42,"method":"tools/call","params":{"name":"bateia_recommend","arguments":` + step0025MCPRequest() + `}}`,
	}
	writer := &bytes.Buffer{}
	server := NewMCPServer(step0025LegacyFunnel(), strings.NewReader(strings.Join(lines, "\n")+"\n"), writer, recommender)
	server.WithDurableRuns(durable)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var responses []JSONRPCResponse
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		responses = append(responses, response)
	}
	if len(responses) != 4 || recommender.calls != 1 || durable.startCalls+durable.getCalls+durable.artifactCalls+durable.cancelCalls != 0 {
		t.Fatalf("legacy routing acquired async behavior: responses=%d sync=%d async=%d/%d/%d/%d", len(responses), recommender.calls, durable.startCalls, durable.getCalls, durable.artifactCalls, durable.cancelCalls)
	}
	toolsObject := decodeStep0025Result(t, responses[0].Result)
	tools := toolsObject["tools"].([]interface{})
	wantTools := []interface{}{}
	for _, raw := range []string{step0026LegacySearchTool, step0026LegacyBenchmarkTool, step0025ExactToolJSON} {
		var value interface{}
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			t.Fatal(err)
		}
		wantTools = append(wantTools, value)
	}
	if !reflect.DeepEqual(tools[:3], wantTools) {
		got, _ := json.Marshal(tools[:3])
		want, _ := json.Marshal(wantTools)
		t.Fatalf("certified first-three schemas changed:\ngot  %s\nwant %s", got, want)
	}
	search := step0025CanonicalizeLegacyGeneratedAt(t, mustStep0025JSON(t, responses[1]))
	if search != step0025BaselineSearchSuccess {
		t.Fatalf("certified search golden changed:\n%s", search)
	}
	benchmarkGolden := `{"jsonrpc":"2.0","id":41,"error":{"code":-32602,"message":"Invalid params: depth deve estar entre 1 e 5"}}`
	if got := string(mustStep0025JSON(t, responses[2])); got != benchmarkGolden {
		t.Fatalf("certified benchmark behavior changed:\ngot  %s\nwant %s", got, benchmarkGolden)
	}
	for _, index := range []int{1, 3} {
		if responses[index].Error != nil {
			t.Fatalf("legacy tool %d did not execute: %#v", index, responses[index].Error)
		}
		encoded, _ := json.Marshal(responses[index])
		for _, forbidden := range []string{`"run_id"`, `"lifecycle"`, `"idempotency_key"`} {
			if strings.Contains(string(encoded), forbidden) {
				t.Fatalf("legacy response acquired durable field %s: %s", forbidden, encoded)
			}
		}
	}
	result := decodeStep0025Result(t, responses[3].Result)
	item := result["content"].([]interface{})[0].(map[string]interface{})
	if item["text"] != step0026R3LegacyRecommendationGolden {
		t.Fatalf("synchronous recommendation golden changed: %v", item["text"])
	}
}

func TestStep0026R2EveryNewToolObjectIsClosedAtRuntime(t *testing.T) {
	id := "run_0123456789abcdef0123456789abcdef"
	start := step0026StartArguments()
	cases := map[string][]string{
		"bateia_start_recommendation": {
			strings.TrimSuffix(start, "}") + `,"extra":true}`,
			strings.Replace(start, `"idea":{`, `"idea":{"extra":true,`, 1),
			strings.Replace(start, `"policy":{`, `"policy":{"extra":true,`, 1),
			strings.Replace(start, `"target":{`, `"target":{"extra":true,`, 1),
			strings.Replace(start, `"search":{`, `"search":{"extra":true,`, 1),
			strings.Replace(start, `"idempotency_key":"agent-task-1",`, "", 1),
			strings.Replace(start, `"idea":`, `"missing_idea":`, 1),
			strings.Replace(start, `"policy":`, `"missing_policy":`, 1),
			strings.Replace(start, `"target":`, `"missing_target":`, 1),
			strings.Replace(start, `"search":`, `"missing_search":`, 1),
		},
		"bateia_get_run": {
			`{}`, `{"run_id":"` + id + `","extra":true}`, `{"run_id":1}`,
		},
		"bateia_cancel_run": {
			`{}`, `{"run_id":"` + id + `","extra":true}`, `{"run_id":null}`,
		},
		"bateia_get_artifact": {
			`{"artifact":"failures","cursor":"","limit":1}`,
			`{"run_id":"` + id + `","cursor":"","limit":1}`,
			`{"run_id":"` + id + `","artifact":"failures","limit":1}`,
			`{"run_id":"` + id + `","artifact":"failures","cursor":""}`,
			`{"run_id":"` + id + `","artifact":"failures","cursor":"","limit":1,"extra":true}`,
			`{"run_id":"` + id + `","artifact":"failures","cursor":"","limit":1.5}`,
		},
	}
	for tool, arguments := range cases {
		for index, argument := range arguments {
			t.Run(tool+string(rune('a'+index)), func(t *testing.T) {
				spy := &step0026DurableSpy{}
				responses := step0026Serve(t, spy, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"`+tool+`","arguments":`+argument+`}}`)
				if len(responses) != 1 || responses[0].Error == nil || responses[0].Error.Code != -32602 || responses[0].Error.Message != "Invalid params" {
					t.Fatalf("closed input accepted: %#v", responses)
				}
				if spy.startCalls+spy.getCalls+spy.artifactCalls+spy.cancelCalls != 0 {
					t.Fatal("invalid closed input touched port")
				}
			})
		}
	}
}

func TestStep0026R2EveryLifecycleProjectsExplicitNullAndEmptyArrays(t *testing.T) {
	id := domain.RunID("run_0123456789abcdef0123456789abcdef")
	for _, lifecycle := range []domain.RecommendationLifecycle{domain.RunQueued, domain.RunRunning, domain.RunCancelling, domain.RunSucceeded, domain.RunFailed, domain.RunCancelled, domain.RunInterrupted} {
		t.Run(string(lifecycle), func(t *testing.T) {
			view := domain.RecommendationRunView{SchemaVersion: "1", RunID: id, Lifecycle: lifecycle, Version: 3, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8}, Artifacts: []domain.RunArtifactIndex{}}
			if lifecycle == domain.RunSucceeded {
				status := domain.IdeaRunPartial
				view.RecommendationStatus = &status
			}
			if lifecycle == domain.RunFailed {
				view.Failure = &domain.RecommendationRunFailure{Code: "RECOMMENDATION_FAILED", Stage: domain.RunStageCollect, Message: "recommendation execution failed"}
			}
			spy := &step0026DurableSpy{run: view}
			response := step0026Serve(t, spy, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+string(id)+`"}}}`)[0]
			encoded := step0026R2ResultText(t, response)
			for _, required := range []string{`"recommendation_status":`, `"artifacts":[]`, `"failure":`} {
				if !strings.Contains(encoded, required) {
					t.Fatalf("%s omitted closed field %s: %s", lifecycle, required, encoded)
				}
			}
			if lifecycle != domain.RunSucceeded && strings.Contains(encoded, `"recommendation_status":"`) {
				t.Fatalf("%s invented recommendation status: %s", lifecycle, encoded)
			}
		})
	}
}

func TestStep0026R2MissingAndTerminalFailuresRemainDistinctAndSanitized(t *testing.T) {
	id := "run_0123456789abcdef0123456789abcdef"
	missing := &step0026DurableSpy{err: domain.ErrRunNotFound}
	responses := step0026Serve(t, missing,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+id+`"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+id+`"}}}`,
	)
	if missing.getCalls != 1 || missing.cancelCalls != 1 {
		t.Fatalf("missing lookup not routed exactly once: %#v", missing)
	}
	for _, response := range responses {
		if response.Error == nil || response.Error.Code != -32020 || response.Error.Message != "RUN_NOT_FOUND" {
			t.Fatalf("missing projected incorrectly: %#v", response)
		}
	}
	for _, failure := range []domain.RecommendationRunFailure{
		{Code: "RECOMMENDATION_FAILED", Stage: domain.RunStageDiscover, Message: "recommendation execution failed"},
		{Code: "RESULT_INVALID", Stage: domain.RunStageAssess, Message: "recommendation result was invalid"},
		{Code: "ARTIFACT_LIMIT_EXCEEDED", Stage: domain.RunStagePublish, Message: "artifact publication exceeded configured quota"},
	} {
		spy := &step0026DurableSpy{run: domain.RecommendationRunView{SchemaVersion: "1", Lifecycle: domain.RunFailed, Version: 5, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 2}, Artifacts: []domain.RunArtifactIndex{}, Failure: &failure}}
		response := step0026Serve(t, spy, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+id+`"}}}`)[0]
		encoded := step0026R2ResultText(t, response)
		if !strings.Contains(encoded, `"code":"`+failure.Code+`"`) || !strings.Contains(encoded, `"message":"`+failure.Message+`"`) {
			t.Fatalf("controlled terminal failure changed: %s", encoded)
		}
		for _, forbidden := range []string{"TOKEN", "Bearer", "/private", "idempotency"} {
			if strings.Contains(encoded, forbidden) {
				t.Fatalf("terminal failure leaked %s: %s", forbidden, encoded)
			}
		}
	}
	malformedPage := &step0026DurableSpy{page: domain.ArtifactPage{SchemaVersion: "1", RunID: domain.RunID(id), Artifact: domain.ArtifactFailures, ArtifactHash: "sha256:" + strings.Repeat("a", 64), Entries: [][]byte{[]byte(`{"raw":"TOKEN_CANARY"`)}}}
	response := step0026Serve(t, malformedPage, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_get_artifact","arguments":{"run_id":"`+id+`","artifact":"failures","cursor":"","limit":1}}}`)[0]
	if response.Error == nil || response.Error.Code != -32603 || strings.Contains(response.Error.Message, "TOKEN_CANARY") {
		t.Fatalf("malformed artifact bytes leaked/became success: %#v", response)
	}
}

func step0026R2ResultText(t *testing.T, response JSONRPCResponse) string {
	t.Helper()
	result := decodeStep0025Result(t, response.Result)
	content, ok := result["content"].([]interface{})
	if !ok || len(content) != 1 {
		t.Fatalf("MCP result has no single content item: %#v", result)
	}
	item, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("MCP content item is not an object: %#v", content[0])
	}
	text, ok := item["text"].(string)
	if !ok {
		t.Fatalf("MCP content text absent: %#v", item)
	}
	var closed interface{}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	if err := decoder.Decode(&closed); err != nil {
		t.Fatalf("MCP content text is not JSON: %v text=%q", err, text)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("MCP content text has trailing JSON: %v text=%q", err, text)
	}
	return text
}
