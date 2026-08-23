/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/server.md
 * @layer L2
 * @updated 2026-08-03
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPServer_Initialize(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n"
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	server := NewMCPServer(nil, reader, writer)
	err := server.Serve(context.Background())
	if err != nil {
		t.Fatalf("Serve falhou: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &resp); err != nil {
		t.Fatalf("Erro ao decodificar resposta JSON-RPC: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("Esperava jsonrpc '2.0', obteve '%s'", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Errorf("Esperava sem erro, obteve: %v", resp.Error)
	}
}

func TestMCPServer_ToolsList(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	server := NewMCPServer(nil, reader, writer)
	err := server.Serve(context.Background())
	if err != nil {
		t.Fatalf("Serve falhou: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &resp); err != nil {
		t.Fatalf("Erro ao decodificar resposta JSON-RPC: %v", err)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Resposta result inválida")
	}

	tools, ok := resultMap["tools"].([]interface{})
	if !ok || len(tools) != 2 {
		t.Fatalf("Esperava 2 ferramentas expostas no MCP, obteve %d", len(tools))
	}
}

func TestMCPServer_RejectsEmptySearchQuery(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_search","arguments":{"query":""}}}` + "\n"
	writer := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(input), writer)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatalf("Serve falhou: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &resp); err != nil {
		t.Fatalf("resposta inválida: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != -32602 {
		t.Fatalf("esperava Invalid params, recebeu: %+v", resp.Error)
	}
}

func TestMCPServer_RejectsMalformedBenchmarkRepository(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_benchmark","arguments":{"repos":"invalid"}}}` + "\n"
	writer := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(input), writer)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatalf("Serve falhou: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &resp); err != nil {
		t.Fatalf("resposta inválida: %v", err)
	}
	if resp.Error == nil || resp.Error.Code != -32602 {
		t.Fatalf("esperava Invalid params, recebeu: %+v", resp.Error)
	}
}
