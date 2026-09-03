/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_runs_step0026_contract_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runstore"
)

func step0026L4Call(t *testing.T, writer io.Writer, scanner *bufio.Scanner, request string) mcp.JSONRPCResponse {
	t.Helper()
	if _, err := io.WriteString(writer, request+"\n"); err != nil {
		t.Fatal(err)
	}
	if !scanner.Scan() {
		t.Fatalf("no MCP response: %v", scanner.Err())
	}
	var response mcp.JSONRPCResponse
	if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
		t.Fatalf("invalid MCP response: %v: %s", err, scanner.Bytes())
	}
	return response
}

func step0026L4Text(t *testing.T, response mcp.JSONRPCResponse) map[string]interface{} {
	t.Helper()
	if response.Error != nil {
		t.Fatalf("MCP error: %#v", response.Error)
	}
	encoded, err := json.Marshal(response.Result)
	if err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil || len(envelope.Content) != 1 || envelope.Content[0].Type != "text" {
		t.Fatalf("invalid MCP envelope: %v %#v", err, response.Result)
	}
	var object map[string]interface{}
	if err := json.Unmarshal([]byte(envelope.Content[0].Text), &object); err != nil {
		t.Fatalf("invalid textual JSON: %v %s", err, envelope.Content[0].Text)
	}
	return object
}

func TestStep0026RealL4DurableCompositionStartGetPaginateCancelAndReopenOffline(t *testing.T) {
	root := t.TempDir()
	secret := []byte("0123456789abcdef0123456789abcdef")
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
	discovery := &step0025WiringDiscovery{}
	content := &step0025WiringContent{}
	recommendation := usecase.NewRecommendationService(nil, discovery, content)

	requestReader, requestWriter := io.Pipe()
	responseReader, responseWriter := io.Pipe()
	server, shutdown, err := NewDurableMCPServer((*usecase.SearchFunnel)(nil), recommendation, root, secret, limits, 1, 4, requestReader, responseWriter)
	if err != nil {
		t.Fatal(err)
	}
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(context.Background()) }()
	scanner := bufio.NewScanner(responseReader)

	listed := step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	listedJSON, _ := json.Marshal(listed.Result)
	if strings.Count(string(listedJSON), `"name":"bateia_`) != 7 {
		t.Fatalf("durable L4 did not inject seven tools: %s", listedJSON)
	}
	arguments := `{"idempotency_key":"l4-offline-1","idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":""},"search":{"strategy":"literal","language":"","locales":["pt-BR"],"max_per_query":1,"max_candidates":1}}`
	started := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":`+arguments+`}}`))
	runID, ok := started["run_id"].(string)
	if !ok || started["lifecycle"] != "QUEUED" || started["version"] != float64(1) {
		t.Fatalf("start was not immediate queued/v1: %#v", started)
	}

	var terminal map[string]interface{}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		terminal = step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+runID+`"}}}`))
		if lifecycle, _ := terminal["lifecycle"].(string); lifecycle == "SUCCEEDED" || lifecycle == "FAILED" {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if terminal["lifecycle"] != "SUCCEEDED" || terminal["recommendation_status"] != "INCONCLUSIVE" {
		t.Fatalf("offline run did not preserve successful inconclusion: %#v", terminal)
	}
	artifact := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_get_artifact","arguments":{"run_id":"`+runID+`","artifact":"recommendation","cursor":"","limit":1}}}`))
	if artifact["artifact"] != "recommendation" || !strings.HasPrefix(artifact["artifact_hash"].(string), "sha256:") {
		t.Fatalf("published artifact unavailable: %#v", artifact)
	}
	entries, ok := artifact["entries"].([]interface{})
	if !ok || len(entries) != 1 {
		t.Fatalf("recommendation document must contain one JSON entry: %#v", artifact)
	}
	cancelled := step0026L4Text(t, step0026L4Call(t, requestWriter, scanner, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+runID+`"}}}`))
	if cancelled["lifecycle"] != "SUCCEEDED" || cancelled["version"] != terminal["version"] {
		t.Fatalf("cancel rewrote terminal winner: before=%#v after=%#v", terminal, cancelled)
	}
	_ = requestWriter.Close()
	if err := <-serveDone; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	input := `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"` + runID + `"}}}` + "\n"
	output := &bytes.Buffer{}
	reopened, closeReopened, err := NewDurableMCPServer((*usecase.SearchFunnel)(nil), recommendation, root, secret, limits, 1, 4, strings.NewReader(input), output)
	if err != nil {
		t.Fatal(err)
	}
	if err := reopened.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := closeReopened(ctx); err != nil {
		t.Fatal(err)
	}
	var response mcp.JSONRPCResponse
	if err := json.Unmarshal(output.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	after := step0026L4Text(t, response)
	if after["lifecycle"] != "SUCCEEDED" || after["version"] != terminal["version"] || after["recommendation_status"] != "INCONCLUSIVE" {
		t.Fatalf("reopen changed terminal: before=%#v after=%#v", terminal, after)
	}
	if discovery.calls == 0 {
		t.Fatal("durable L4 did not inject the certified RecommendationService")
	}
}
