/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/durable_runs_step0026_contract_test.md
 * @layer L2
 * @updated 2026-09-03
 */
package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type step0026DurableSpy struct {
	startCalls, getCalls, artifactCalls, cancelCalls int
	key                                              string
	request                                          domain.RecommendationRequest
	read                                             domain.ArtifactReadRequest
	receipt                                          domain.RecommendationStartReceipt
	run                                              domain.RecommendationRunView
	page                                             domain.ArtifactPage
	err                                              error
}

func (s *step0026DurableSpy) Start(_ context.Context, key string, request domain.RecommendationRequest) (domain.RecommendationStartReceipt, error) {
	s.startCalls++
	s.key, s.request = key, request
	return s.receipt, s.err
}
func (s *step0026DurableSpy) GetRun(_ context.Context, id domain.RunID) (domain.RecommendationRunView, error) {
	s.getCalls++
	s.run.RunID = id
	return s.run, s.err
}
func (s *step0026DurableSpy) GetArtifact(_ context.Context, request domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	s.artifactCalls++
	s.read = request
	return s.page, s.err
}
func (s *step0026DurableSpy) Cancel(_ context.Context, id domain.RunID) (domain.RecommendationRunView, error) {
	s.cancelCalls++
	s.run.RunID = id
	return s.run, s.err
}

func step0026Serve(t *testing.T, durable DurableRecommendationRuns, lines ...string) []JSONRPCResponse {
	t.Helper()
	writer := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(strings.Join(lines, "\n")+"\n"), writer)
	server.WithDurableRuns(durable)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var responses []JSONRPCResponse
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			return responses
		} else if err != nil {
			t.Fatalf("invalid response: %v\n%s", err, writer.Bytes())
		}
		responses = append(responses, response)
	}
}

func step0026Tools(t *testing.T, durable DurableRecommendationRuns) []interface{} {
	t.Helper()
	responses := step0026Serve(t, durable, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if len(responses) != 1 || responses[0].Error != nil {
		t.Fatalf("tools/list failed: %#v", responses)
	}
	result := decodeStep0025Result(t, responses[0].Result)
	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatalf("tools is not array: %#v", result)
	}
	return tools
}

func TestStep0026ToolsListAddsExactlyFourAndPreservesStep0025Deeply(t *testing.T) {
	legacy := step0026Tools(t, nil)
	tools := step0026Tools(t, &step0026DurableSpy{})
	if len(legacy) != 3 || len(tools) != 7 || !reflect.DeepEqual(legacy, tools[:3]) {
		t.Fatalf("non-additive MCP refinement: legacy=%#v durable=%#v", legacy, tools)
	}
	wantNames := []string{"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run"}
	for index, want := range wantNames {
		tool := tools[index].(map[string]interface{})
		if tool["name"] != want {
			t.Fatalf("tool[%d] = %v, want %s", index, tool["name"], want)
		}
	}

	recommend := tools[2].(map[string]interface{})["inputSchema"].(map[string]interface{})
	start := tools[3].(map[string]interface{})["inputSchema"].(map[string]interface{})
	recommendProperties := recommend["properties"].(map[string]interface{})
	startProperties := start["properties"].(map[string]interface{})
	for _, name := range []string{"idea", "policy", "target", "search"} {
		if !reflect.DeepEqual(recommendProperties[name], startProperties[name]) {
			t.Fatalf("start changed recommendation schema %s", name)
		}
	}
	wantKey := map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(128), "pattern": `^[\u0020-\u002E\u0030-\u005B\u005D-\u007E]+$`}
	if !reflect.DeepEqual(startProperties["idempotency_key"], wantKey) {
		t.Fatalf("idempotency key schema = %#v", startProperties["idempotency_key"])
	}
	if got := start["required"]; !reflect.DeepEqual(got, []interface{}{"idempotency_key", "idea"}) && !reflect.DeepEqual(got, []string{"idempotency_key", "idea"}) {
		t.Fatalf("start required = %#v", got)
	}

	wantSchemas := map[string]string{
		"bateia_get_run":      `{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"}},"required":["run_id"]}`,
		"bateia_get_artifact": `{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"},"artifact":{"type":"string","enum":["request","recommendation","decision","alternatives","failures","provenance"]},"cursor":{"anyOf":[{"const":""},{"type":"string","maxLength":512,"pattern":"^[A-Za-z0-9_-]+\\.[A-Za-z0-9_-]+$"}]},"limit":{"type":"integer","minimum":1,"maximum":50}},"required":["run_id","artifact","cursor","limit"]}`,
		"bateia_cancel_run":   `{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"}},"required":["run_id"]}`,
	}
	for _, index := range []int{4, 5, 6} {
		tool := tools[index].(map[string]interface{})
		var want interface{}
		if err := json.Unmarshal([]byte(wantSchemas[tool["name"].(string)]), &want); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tool["inputSchema"], want) {
			got, _ := json.Marshal(tool["inputSchema"])
			t.Fatalf("schema %s changed: %s", tool["name"], got)
		}
	}
}

func step0026StartArguments() string {
	return `{"idempotency_key":"agent-task-1","idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"pantry"},"search":{"strategy":"literal","language":"","locales":["en-US"],"max_per_query":1,"max_candidates":1}}`
}

func TestStep0026DurableToolsUseInjectedPortAndExactSuccessEnvelopes(t *testing.T) {
	id := domain.RunID("run_0123456789abcdef0123456789abcdef")
	status := domain.IdeaRunInconclusive
	spy := &step0026DurableSpy{
		receipt: domain.RecommendationStartReceipt{SchemaVersion: "1", RunID: id, Lifecycle: domain.RunQueued, Version: 1},
		run:     domain.RecommendationRunView{SchemaVersion: "1", Lifecycle: domain.RunSucceeded, Version: 9, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 8, Completed: 8}, RecommendationStatus: &status, Artifacts: []domain.RunArtifactIndex{}, Failure: nil},
		page:    domain.ArtifactPage{SchemaVersion: "1", RunID: id, Artifact: domain.ArtifactFailures, ArtifactHash: "sha256:" + strings.Repeat("a", 64), Entries: [][]byte{[]byte(`{"code":"SAFE"}`)}, NextCursor: nil},
	}
	responses := step0026Serve(t, spy,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":`+step0026StartArguments()+`}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+string(id)+`"}}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_get_artifact","arguments":{"run_id":"`+string(id)+`","artifact":"failures","cursor":"","limit":20}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+string(id)+`"}}}`,
	)
	if len(responses) != 4 || spy.startCalls != 1 || spy.getCalls != 1 || spy.artifactCalls != 1 || spy.cancelCalls != 1 {
		t.Fatalf("durable routing failed: responses=%#v calls=%d/%d/%d/%d", responses, spy.startCalls, spy.getCalls, spy.artifactCalls, spy.cancelCalls)
	}
	if spy.key != "agent-task-1" || spy.request.Idea.Problem != "food waste" || spy.read.Limit != 20 {
		t.Fatalf("arguments changed before port: key=%q request=%#v read=%#v", spy.key, spy.request, spy.read)
	}
	wantTexts := []string{
		`{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"QUEUED","version":1}`,
		`{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"SUCCEEDED","version":9,"stage":"terminal","progress":{"planned":8,"completed":8},"recommendation_status":"INCONCLUSIVE","artifacts":[],"failure":null}`,
		`{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","artifact":"failures","artifact_hash":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","entries":[{"code":"SAFE"}],"next_cursor":null}`,
		`{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"SUCCEEDED","version":9,"stage":"terminal","progress":{"planned":8,"completed":8},"recommendation_status":"INCONCLUSIVE","artifacts":[],"failure":null}`,
	}
	for index, response := range responses {
		if response.Error != nil {
			t.Fatalf("response[%d] error: %#v", index, response.Error)
		}
		result := decodeStep0025Result(t, response.Result)
		content := result["content"].([]interface{})
		item := content[0].(map[string]interface{})
		if item["type"] != "text" || item["text"] != wantTexts[index] {
			t.Fatalf("response[%d] envelope = %#v, want %s", index, item, wantTexts[index])
		}
	}
}

func TestStep0026StrictRuntimeValidationHappensBeforeDurablePort(t *testing.T) {
	id := "run_0123456789abcdef0123456789abcdef"
	base := step0026StartArguments()
	cases := map[string]string{
		"start-extra":       strings.TrimSuffix(base, "}") + `,"path":"/tmp/escape"}`,
		"start-missing-key": strings.Replace(base, `"idempotency_key":"agent-task-1",`, "", 1),
		"start-dotdot":      strings.Replace(base, "agent-task-1", "agent..task", 1),
		"start-slash":       strings.Replace(base, "agent-task-1", "agent/task", 1),
		"bad-run-id":        `{"run_id":"../escape"}`,
		"long-run-id":       `{"run_id":"run_` + strings.Repeat("a", 1024) + `"}`,
		"unknown-artifact":  `{"run_id":"` + id + `","artifact":"raw_readme","cursor":"","limit":20}`,
		"bad-limit-zero":    `{"run_id":"` + id + `","artifact":"failures","cursor":"","limit":0}`,
		"bad-limit-high":    `{"run_id":"` + id + `","artifact":"failures","cursor":"","limit":51}`,
		"cursor-syntax":     `{"run_id":"` + id + `","artifact":"failures","cursor":"not-a-signed-cursor","limit":1}`,
	}
	for name, arguments := range cases {
		t.Run(name, func(t *testing.T) {
			spy := &step0026DurableSpy{}
			tool := "bateia_start_recommendation"
			if strings.Contains(arguments, `"artifact"`) {
				tool = "bateia_get_artifact"
			} else if strings.Contains(arguments, `"run_id"`) {
				tool = "bateia_get_run"
			}
			responses := step0026Serve(t, spy, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"`+tool+`","arguments":`+arguments+`}}`)
			if len(responses) != 1 || responses[0].Error == nil || responses[0].Error.Code != -32602 || responses[0].Error.Message != "Invalid params" {
				t.Fatalf("invalid runtime response: %#v", responses)
			}
			if spy.startCalls+spy.getCalls+spy.artifactCalls+spy.cancelCalls != 0 {
				t.Fatal("invalid input reached durable port")
			}
		})
	}
}

func TestStep0026PublicErrorsAreExactAndNeverLeakCanaries(t *testing.T) {
	cases := []struct {
		err     error
		code    int
		message string
	}{
		{domain.ErrRunNotFound, -32020, "RUN_NOT_FOUND"},
		{domain.ErrArtifactNotReady, -32021, "ARTIFACT_NOT_READY"},
		{domain.ErrIdempotencyConflict, -32022, "IDEMPOTENCY_CONFLICT"},
		{domain.ErrQueueFull, -32023, "QUEUE_FULL"},
		{domain.ErrQuotaExceeded, -32024, "QUOTA_EXCEEDED"},
		{domain.ErrIntegrity, -32025, "INTEGRITY_ERROR"},
		{domain.ErrCursorInvalid, -32026, "CURSOR_INVALID"},
		{domain.ErrRunStoreUnavailable, -32027, "RUN_STORE_UNAVAILABLE"},
		{errors.New("TOKEN_CANARY /private/path Authorization: Bearer SECRET README_CANARY"), -32603, "Internal error"},
	}
	for _, item := range cases {
		spy := &step0026DurableSpy{err: item.err}
		responses := step0026Serve(t, spy, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"run_0123456789abcdef0123456789abcdef"}}}`)
		if len(responses) != 1 || responses[0].Error == nil || responses[0].Error.Code != item.code || responses[0].Error.Message != item.message {
			t.Fatalf("error %v mapped to %#v, want %d/%s", item.err, responses, item.code, item.message)
		}
		encoded, _ := json.Marshal(responses[0])
		for _, canary := range []string{"TOKEN_CANARY", "/private/path", "Authorization", "SECRET", "README_CANARY"} {
			if strings.Contains(string(encoded), canary) {
				t.Fatalf("canary %q leaked in %s", canary, encoded)
			}
		}
	}
}
