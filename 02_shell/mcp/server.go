/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/server.md
 * @layer L2
 * @updated 2026-08-03
 */

package mcp

import (
	"bufio"
	"bytes"
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
	funnel      *usecase.SearchFunnel
	recommender RecommendationRecommender
	reader      io.Reader
	writer      io.Writer
}

type RecommendationRecommender interface {
	Recommend(context.Context, domain.RecommendationRequest) (domain.RecommendationResult, error)
}

func NewMCPServer(funnel *usecase.SearchFunnel, reader io.Reader, writer io.Writer, recommenders ...RecommendationRecommender) *MCPServer {
	server := &MCPServer{funnel: funnel, reader: reader, writer: writer}
	if len(recommenders) > 0 {
		server.recommender = recommenders[0]
	}
	return server
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
			recommendationTool(),
		}
		s.sendResult(req.ID, map[string]interface{}{"tools": tools})

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
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

func (s *MCPServer) handleToolCall(ctx context.Context, id interface{}, toolName string, rawArguments json.RawMessage) {
	if toolName == "bateia_recommend" {
		s.handleRecommendation(ctx, id, rawArguments)
		return
	}
	var args map[string]interface{}
	if err := json.Unmarshal(rawArguments, &args); err != nil {
		s.sendError(id, -32602, "Invalid params")
		return
	}
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

type recommendationIdeaArguments struct {
	Problem        string   `json:"problem"`
	Users          []string `json:"users,omitempty"`
	Capabilities   []string `json:"capabilities"`
	Constraints    []string `json:"constraints,omitempty"`
	NonGoals       []string `json:"non_goals,omitempty"`
	PreferredStack []string `json:"preferred_stack,omitempty"`
}

type recommendationPolicyArguments struct {
	SchemaVersion       *string   `json:"schema_version"`
	AllowedLicenses     *[]string `json:"allowed_licenses"`
	RequireLocal        *bool     `json:"require_local"`
	AllowArchived       *bool     `json:"allow_archived"`
	MaxComponents       *int      `json:"max_components"`
	ExtendMinCoverage   *float64  `json:"extend_min_coverage"`
	IntegrationMinValue *float64  `json:"integration_min_value"`
}

type recommendationTargetArguments struct {
	Name string `json:"name"`
}

type recommendationSearchArguments struct {
	Strategy      *string   `json:"strategy"`
	Language      string    `json:"language,omitempty"`
	Locales       *[]string `json:"locales"`
	MaxPerQuery   *int      `json:"max_per_query"`
	MaxCandidates *int      `json:"max_candidates"`
}

type recommendationArguments struct {
	Idea   *recommendationIdeaArguments   `json:"idea"`
	Policy *recommendationPolicyArguments `json:"policy"`
	Target *recommendationTargetArguments `json:"target"`
	Search *recommendationSearchArguments `json:"search"`
}

func decodeRecommendationArguments(raw json.RawMessage) (domain.RecommendationRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input recommendationArguments
	if err := decoder.Decode(&input); err != nil {
		return domain.RecommendationRequest{}, err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return domain.RecommendationRequest{}, fmt.Errorf("arguments devem conter um único objeto JSON")
	}
	if input.Idea == nil || input.Policy == nil || input.Target == nil || input.Search == nil ||
		input.Policy.SchemaVersion == nil || input.Policy.AllowedLicenses == nil || input.Policy.RequireLocal == nil || input.Policy.AllowArchived == nil || input.Policy.MaxComponents == nil || input.Policy.ExtendMinCoverage == nil || input.Policy.IntegrationMinValue == nil ||
		input.Search.Strategy == nil || input.Search.Locales == nil || input.Search.MaxPerQuery == nil || input.Search.MaxCandidates == nil {
		return domain.RecommendationRequest{}, fmt.Errorf("campos obrigatórios ausentes")
	}
	request := domain.RecommendationRequest{
		Idea: domain.IdeaSpec{
			Problem: input.Idea.Problem, Users: input.Idea.Users, Capabilities: input.Idea.Capabilities,
			Constraints: input.Idea.Constraints, NonGoals: input.Idea.NonGoals, PreferredStack: input.Idea.PreferredStack,
		},
		Policy: domain.ReuseDecisionPolicy{
			SchemaVersion: *input.Policy.SchemaVersion, AllowedLicenses: *input.Policy.AllowedLicenses,
			RequireLocal: *input.Policy.RequireLocal, AllowArchived: *input.Policy.AllowArchived,
			MaxComponents: *input.Policy.MaxComponents, ExtendMinCoverage: *input.Policy.ExtendMinCoverage,
			IntegrationMinValue: *input.Policy.IntegrationMinValue,
		},
		Target: domain.SoftwareAlternativeTarget{Name: input.Target.Name},
		Search: domain.RecommendationSearch{
			Strategy: *input.Search.Strategy, Language: input.Search.Language, Locales: *input.Search.Locales,
			MaxPerQuery: *input.Search.MaxPerQuery, MaxCandidates: *input.Search.MaxCandidates,
		},
	}
	return request, nil
}

func (s *MCPServer) handleRecommendation(ctx context.Context, id interface{}, rawArguments json.RawMessage) {
	request, err := decodeRecommendationArguments(rawArguments)
	if err != nil {
		s.sendError(id, -32602, "Invalid params")
		return
	}
	if _, err = usecase.ValidateRecommendationRequest(request); err != nil {
		s.sendError(id, -32602, "Invalid params")
		return
	}
	if s.recommender == nil {
		s.sendError(id, -32603, "Recommendation service unavailable")
		return
	}
	result, err := s.recommender.Recommend(ctx, request)
	if err != nil {
		s.sendError(id, -32603, "Recommendation failed")
		return
	}
	encoded, err := usecase.MarshalRecommendationResult(result)
	if err != nil {
		s.sendError(id, -32603, "Recommendation failed")
		return
	}
	s.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
}

func recommendationTool() map[string]interface{} {
	stringArray := func() map[string]interface{} {
		return map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}}
	}
	return map[string]interface{}{
		"name":        "bateia_recommend",
		"description": "Transforma uma ideia em decisão de reúso e alternativas de software com evidências rastreáveis",
		"inputSchema": map[string]interface{}{
			"type": "object", "additionalProperties": false,
			"properties": map[string]interface{}{
				"idea": map[string]interface{}{
					"type": "object", "additionalProperties": false,
					"properties": map[string]interface{}{
						"problem": map[string]interface{}{"type": "string", "minLength": 1}, "users": stringArray(),
						"capabilities": map[string]interface{}{"type": "array", "minItems": 1, "items": map[string]interface{}{"type": "string", "minLength": 1}},
						"constraints":  stringArray(), "non_goals": stringArray(), "preferred_stack": stringArray(),
					},
					"required": []string{"problem", "capabilities"},
				},
				"policy": map[string]interface{}{
					"type": "object", "additionalProperties": false,
					"properties": map[string]interface{}{
						"schema_version":   map[string]interface{}{"type": "string", "const": "1"},
						"allowed_licenses": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string", "minLength": 1}},
						"require_local":    map[string]interface{}{"type": "boolean"}, "allow_archived": map[string]interface{}{"type": "boolean"},
						"max_components":        map[string]interface{}{"type": "integer", "minimum": 2},
						"extend_min_coverage":   map[string]interface{}{"type": "number", "exclusiveMinimum": 0, "maximum": 1},
						"integration_min_value": map[string]interface{}{"type": "number", "exclusiveMinimum": 0, "maximum": 1},
					},
					"required": []string{"schema_version", "allowed_licenses", "require_local", "allow_archived", "max_components", "extend_min_coverage", "integration_min_value"},
				},
				"target": map[string]interface{}{"type": "object", "additionalProperties": false, "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string"}}},
				"search": map[string]interface{}{
					"type": "object", "additionalProperties": false,
					"properties": map[string]interface{}{
						"strategy": map[string]interface{}{"type": "string", "enum": []string{"literal", "lexical"}}, "language": map[string]interface{}{"type": "string"},
						"locales":       map[string]interface{}{"type": "array", "uniqueItems": true, "items": map[string]interface{}{"type": "string", "enum": []string{"pt-BR", "en-US", "es-ES", "fr-FR", "de-DE", "ja-JP", "ko-KR", "ru-RU", "ar-SA", "hi-IN"}}},
						"max_per_query": map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 25}, "max_candidates": map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 20},
					},
					"required": []string{"strategy", "locales", "max_per_query", "max_candidates"},
				},
			},
			"required": []string{"idea", "policy", "target", "search"},
		},
	}
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
