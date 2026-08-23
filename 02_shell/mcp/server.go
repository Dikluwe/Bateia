/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/server.md
 * @layer L2
 * @updated 2026-08-03
 */

package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/report"
)

// JSONRPCRequest representa a mensagem enviada pelo cliente MCP.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse representa a resposta enviada pelo servidor MCP.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPServer lida com a comunicação JSON-RPC via stdio.
type MCPServer struct {
	funnel *usecase.SearchFunnel
	reader io.Reader
	writer io.Writer
}

func NewMCPServer(funnel *usecase.SearchFunnel, reader io.Reader, writer io.Writer) *MCPServer {
	return &MCPServer{
		funnel: funnel,
		reader: reader,
		writer: writer,
	}
}

func (s *MCPServer) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(s.reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(ctx, req)
	}

	return scanner.Err()
}

func (s *MCPServer) handleRequest(ctx context.Context, req JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		res := map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]string{
				"name":    "bateia-mcp",
				"version": "1.0.0",
			},
		}
		s.sendResult(req.ID, res)

	case "notifications/initialized":
		// Notificação sem resposta necessária.

	case "tools/list":
		tools := []map[string]interface{}{
			{
				"name":        "bateia_search",
				"description": "Executa busca arquitetural no GitHub/Sourcegraph reduzindo repositórios via funil de 3 estágios",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Query de busca (ex: 'hexagonal architecture' ou nome do pacote)",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"description": "Linguagem alvo (go, typescript, python, rust, php, etc.)",
							"default":     "go",
						},
						"depth": map[string]interface{}{
							"type":        "integer",
							"description": "Profundidade da busca recursiva em grafo de dependências",
							"default":     1,
						},
					},
					"required": []string{"query"},
				},
			},
			{
				"name":        "bateia_benchmark",
				"description": "Executa comparação arquitetural direta de uma lista de repositórios conhecidos",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"repos": map[string]interface{}{
							"type":        "string",
							"description": "Repositórios separados por vírgula no formato owner/repo (ex: 'traggo/server,solidtime-io/solidtime')",
						},
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Título ou tema da comparação de benchmark",
							"default":     "Benchmark de Repositórios",
						},
					},
					"required": []string{"repos"},
				},
			},
		}
		s.sendResult(req.ID, map[string]interface{}{"tools": tools})

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.sendError(req.ID, -32602, "Invalid params")
			return
		}

		s.handleToolCall(ctx, req.ID, params.Name, params.Arguments)

	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *MCPServer) handleToolCall(ctx context.Context, id interface{}, toolName string, args map[string]interface{}) {
	var filter domain.QueryFilter

	switch toolName {
	case "bateia_search":
		query, _ := args["query"].(string)
		lang, _ := args["language"].(string)
		if lang == "" {
			lang = "go"
		}
		depthFloat, _ := args["depth"].(float64)
		depth := int(depthFloat)
		if depth < 1 {
			depth = 1
		}

		filter = domain.QueryFilter{
			Query:      query,
			Language:   lang,
			MaxResults: 10,
			Depth:      depth,
		}

	case "bateia_benchmark":
		reposStr, _ := args["repos"].(string)
		query, _ := args["query"].(string)
		if query == "" {
			query = "Benchmark MCP"
		}

		var directRepos []string
		for _, r := range strings.Split(reposStr, ",") {
			trimmed := strings.TrimSpace(r)
			if trimmed != "" {
				directRepos = append(directRepos, trimmed)
			}
		}

		filter = domain.QueryFilter{
			Query:       query,
			Language:    "go",
			MaxResults:  len(directRepos),
			DirectRepos: directRepos,
		}

	default:
		s.sendError(id, -32601, fmt.Sprintf("Unknown tool: %s", toolName))
		return
	}
	if err := usecase.ValidateQueryFilter(filter); err != nil {
		s.sendError(id, -32602, fmt.Sprintf("Invalid params: %v", err))
		return
	}

	searchReport, err := s.funnel.Execute(ctx, filter)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Erro ao executar Bateia: %v", err))
		return
	}

	mdContent := report.GenerateMarkdownReport(searchReport)

	s.sendResult(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": mdContent,
			},
		},
	})
}

func (s *MCPServer) sendResult(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	bytes, _ := json.Marshal(resp)
	_, _ = s.writer.Write(append(bytes, '\n'))
}

func (s *MCPServer) sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	bytes, _ := json.Marshal(resp)
	_, _ = s.writer.Write(append(bytes, '\n'))
}
