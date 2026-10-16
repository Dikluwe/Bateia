/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/market_discovery.md
 * @layer L2
 * @updated 2026-10-06
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type MarketDiscoveryStarter interface {
	StartMarketDiscovery(context.Context, string, domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartReceipt, error)
}

func (s *MCPServer) WithMarketDiscovery(connectors []domain.MarketDiscoveryConnector, starter MarketDiscoveryStarter) *MCPServer {
	if s == nil || starter == nil || s.marketDiscovery != nil {
		return s
	}
	registry, err := usecase.NewDiscoveryConnectorRegistry(connectors)
	if err != nil {
		return s
	}
	s.marketDiscovery = &marketInterceptReader{server: s, market: starter, registry: registry}
	return s
}

type marketInterceptReader struct {
	server   *MCPServer
	market   MarketDiscoveryStarter
	registry *usecase.DiscoveryConnectorRegistry
}

func decodeMarketEnvelope(data []byte, request *JSONRPCRequest) (bool, bool) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if decoder.Decode(request) != nil {
		return false, false
	}
	return true, decoder.Decode(&struct{}{}) != io.EOF
}

func marketCallRequest(request JSONRPCRequest) bool {
	if request.Method != "tools/call" {
		return false
	}
	var params struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(request.Params, &params) != nil {
		return false
	}
	return params.Name == "bateia_start_market_discovery" || params.Name == "bateia_list_connectors"
}

var malformedMarketIDPattern = regexp.MustCompile(`"id"\s*:\s*([0-9]+)`)

func malformedMarketCallID(data []byte) (interface{}, bool) {
	text := string(data)
	if !strings.Contains(text, `"method":"tools/call"`) && !strings.Contains(text, `"method": "tools/call"`) {
		return nil, false
	}
	if !strings.Contains(text, `"name":"bateia_start_market_discovery"`) && !strings.Contains(text, `"name": "bateia_start_market_discovery"`) {
		return nil, false
	}
	match := malformedMarketIDPattern.FindStringSubmatch(text)
	if len(match) != 2 {
		return nil, false
	}
	id, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, false
	}
	return id, true
}

func marketArtifactRequested(raw json.RawMessage) bool {
	var value struct {
		Artifact string `json:"artifact"`
	}
	if json.Unmarshal(raw, &value) != nil {
		return false
	}
	for _, name := range []string{"connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps"} {
		if value.Artifact == name {
			return true
		}
	}
	return false
}

func marketToolsList() (map[string]interface{}, error) {
	return marketToolsListWithRegistry(nil)
}

func (r *marketInterceptReader) marketToolsList() (map[string]interface{}, error) {
	if r == nil {
		return marketToolsList()
	}
	return marketToolsListWithRegistry(r.registry)
}

func marketToolsListWithRegistry(registry *usecase.DiscoveryConnectorRegistry) (map[string]interface{}, error) {
	result, err := durableToolsList()
	if err != nil {
		return nil, err
	}
	tools, ok := result["tools"].([]interface{})
	if !ok || len(tools) != 7 {
		return nil, errors.New("durable tool list invalid")
	}
	artifactTool, ok := tools[5].(map[string]interface{})
	if !ok {
		return nil, errors.New("artifact tool invalid")
	}
	schema, ok := artifactTool["inputSchema"].(map[string]interface{})
	if !ok {
		return nil, errors.New("artifact schema invalid")
	}
	properties := schema["properties"].(map[string]interface{})
	artifact := properties["artifact"].(map[string]interface{})
	artifact["enum"] = []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance", "connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps"}
	tools = append(tools,
		map[string]interface{}{"name": "bateia_list_connectors", "description": "Lista capacidades estáticas de descoberta", "inputSchema": emptyObjectSchema()},
		map[string]interface{}{"name": "bateia_start_market_discovery", "description": "Inicia descoberta durável por idioma e mercado", "inputSchema": schemaForRegistry(registry)},
	)
	return map[string]interface{}{"tools": tools}, nil
}

func schemaForRegistry(registry *usecase.DiscoveryConnectorRegistry) map[string]interface{} {
	if registry == nil {
		return marketStartSchema()
	}
	caps := registry.Capabilities()
	if len(caps) == 3 {
		ids := []string{caps[0].ID, caps[1].ID, caps[2].ID}
		hasApple, hasGH, hasSG := false, false, false
		for _, id := range ids {
			if id == "apple-itunes-search-software" {
				hasApple = true
			} else if id == "github-repository-catalog" {
				hasGH = true
			} else if id == "sourcegraph-repository-catalog" {
				hasSG = true
			}
		}
		if hasApple && hasGH && hasSG {
			return marketStartSchema()
		}
	}
	ids := make([]string, 0, len(caps))
	for _, c := range caps {
		ids = append(ids, c.ID)
	}
	maxConnectors := len(ids)
	if maxConnectors < 3 {
		maxConnectors = 3
	}
	if maxConnectors > 10 {
		maxConnectors = 10
	}
	return buildMarketStartSchemaWithConnectors(ids, maxConnectors)
}

func emptyObjectSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object", "additionalProperties": false, "properties": map[string]interface{}{}}
}

func marketStartSchema() map[string]interface{} {
	textArray := func(min, max int) map[string]interface{} {
		value := map[string]interface{}{"type": "array", "maxItems": float64(max), "items": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(200)}}
		if min > 0 {
			value["minItems"] = float64(min)
		}
		return value
	}
	idea := map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"problem": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(1000)},
			"users":   textArray(0, 20), "capabilities": textArray(1, 20), "constraints": textArray(0, 20), "non_goals": textArray(0, 20), "preferred_stack": textArray(0, 20),
		}, "required": []string{"problem", "capabilities"},
	}
	proposal := map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"locale": map[string]interface{}{"type": "string", "pattern": `^[a-z]{2,3}(-[A-Z]{2})?$`, "minLength": float64(2), "maxLength": float64(6)},
			"query":  marketTextSchema(256), "source_text": marketTextSchema(512),
			"derivation": map[string]interface{}{"type": "string", "const": "EXTERNAL_AGENT_PROPOSAL"}, "rationale": marketTextSchema(512),
		}, "required": []string{"locale", "query", "source_text", "derivation", "rationale"},
	}
	return map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"idempotency_key": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(128), "pattern": `^[\u0020-\u002E\u0030-\u005B\u005D-\u007E]+$`},
			"idea":            idea,
			"query_proposals": map[string]interface{}{"type": "array", "minItems": float64(1), "maxItems": float64(20), "items": proposal},
			"markets":         map[string]interface{}{"type": "array", "minItems": float64(1), "maxItems": float64(50), "uniqueItems": true, "items": map[string]interface{}{"type": "string", "pattern": `^[A-Z]{2}$`, "minLength": float64(2), "maxLength": float64(2)}},
			"connector_ids":   map[string]interface{}{"type": "array", "minItems": float64(1), "maxItems": float64(3), "uniqueItems": true, "items": map[string]interface{}{"type": "string", "enum": []string{"apple-itunes-search-software", "github-repository-catalog", "sourcegraph-repository-catalog"}}},
			"budget":          map[string]interface{}{"type": "object", "additionalProperties": false, "properties": map[string]interface{}{"max_queries": integerRange(1, 20), "max_results_per_query": integerRange(1, 10), "max_total_candidates": integerRange(1, 100)}, "required": []string{"max_queries", "max_results_per_query", "max_total_candidates"}},
		}, "required": []string{"idempotency_key", "idea", "query_proposals", "markets", "connector_ids", "budget"},
	}
}

func buildMarketStartSchemaWithConnectors(connectorIDs []string, maxConnectors int) map[string]interface{} {
	schema := marketStartSchema()
	props := schema["properties"].(map[string]interface{})
	props["connector_ids"] = map[string]interface{}{
		"type":        "array",
		"minItems":    float64(1),
		"maxItems":    float64(maxConnectors),
		"uniqueItems": true,
		"items": map[string]interface{}{
			"type": "string",
			"enum": connectorIDs,
		},
	}
	return schema
}

func marketTextSchema(max int) map[string]interface{} {
	return map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(max), "pattern": `^[^\u0000-\u001F\u007F]+$`}
}
func integerRange(min, max int) map[string]interface{} {
	return map[string]interface{}{"type": "integer", "minimum": float64(min), "maximum": float64(max)}
}

func (r *marketInterceptReader) handleMarket(ctx context.Context, id interface{}, name string, arguments json.RawMessage) {
	if name == "bateia_list_connectors" {
		if !strictEmptyObject(arguments) {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
		encoded, err := json.Marshal(struct {
			SchemaVersion string                                  `json:"schema_version"`
			Connectors    []domain.DiscoveryConnectorCapabilities `json:"connectors"`
		}{"1", r.registry.Capabilities()})
		if err != nil {
			r.server.sendError(id, -32603, "Internal error")
			return
		}
		r.server.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
		return
	}
	key, request, err := decodeMarketStart(arguments, r.registry)
	if err != nil {
		r.server.sendError(id, -32602, "Invalid params")
		return
	}
	receipt, err := r.market.StartMarketDiscovery(ctx, key, request)
	if err != nil {
		code, message := durablePublicError(err)
		if errors.Is(err, domain.ErrInvalidRunInput) {
			code, message = -32602, "Invalid params"
		}
		r.server.sendError(id, code, message)
		return
	}
	encoded, err := json.Marshal(receipt)
	if err != nil {
		r.server.sendError(id, -32603, "Internal error")
		return
	}
	r.server.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
}

func (r *marketInterceptReader) handleMarketArtifact(ctx context.Context, id interface{}, arguments json.RawMessage) {
	request, err := decodeMarketArtifact(arguments)
	if err != nil {
		r.server.sendError(id, -32602, "Invalid params")
		return
	}
	page, err := r.server.durableRuns.durable.GetArtifact(ctx, request)
	if err != nil {
		code, message := durablePublicError(err)
		r.server.sendError(id, code, message)
		return
	}
	encoded, err := marshalDurableValue(page)
	if err != nil {
		r.server.sendError(id, -32603, "Internal error")
		return
	}
	r.server.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
}

func decodeMarketArtifact(raw json.RawMessage) (domain.ArtifactReadRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input struct {
		RunID    *string `json:"run_id"`
		Artifact *string `json:"artifact"`
		Cursor   *string `json:"cursor"`
		Limit    *int    `json:"limit"`
	}
	if decoder.Decode(&input) != nil || decoder.Decode(&struct{}{}) != io.EOF || input.RunID == nil || input.Artifact == nil || input.Cursor == nil || input.Limit == nil || !durableRunIDPattern.MatchString(*input.RunID) || *input.Limit < 1 || *input.Limit > 50 || len(*input.Cursor) > 512 || (*input.Cursor != "" && !durableCursorPattern.MatchString(*input.Cursor)) {
		return domain.ArtifactReadRequest{}, errors.New("invalid")
	}
	if !marketArtifactRequested(raw) {
		return domain.ArtifactReadRequest{}, errors.New("invalid")
	}
	return domain.ArtifactReadRequest{RunID: domain.RunID(*input.RunID), Artifact: domain.RunArtifactName(*input.Artifact), Cursor: *input.Cursor, Limit: *input.Limit}, nil
}

func strictEmptyObject(raw json.RawMessage) bool {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var value struct{}
	return decoder.Decode(&value) == nil && decoder.Decode(&struct{}{}) == io.EOF
}

type marketStartInput struct {
	IdempotencyKey *string          `json:"idempotency_key"`
	Idea           *domain.IdeaSpec `json:"idea"`
	QueryProposals *[]struct {
		Locale     string                       `json:"locale"`
		Query      string                       `json:"query"`
		SourceText string                       `json:"source_text"`
		Derivation domain.MarketQueryDerivation `json:"derivation"`
		Rationale  string                       `json:"rationale"`
	} `json:"query_proposals"`
	Markets      *[]string                     `json:"markets"`
	ConnectorIDs *[]string                     `json:"connector_ids"`
	Budget       *domain.MarketDiscoveryBudget `json:"budget"`
}

func decodeMarketStart(raw json.RawMessage, registry *usecase.DiscoveryConnectorRegistry) (string, domain.MarketDiscoveryStartRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input marketStartInput
	if decoder.Decode(&input) != nil || decoder.Decode(&struct{}{}) != io.EOF || input.IdempotencyKey == nil || input.Idea == nil || input.QueryProposals == nil || input.Markets == nil || input.ConnectorIDs == nil || input.Budget == nil || !validDurableKey(*input.IdempotencyKey) {
		return "", domain.MarketDiscoveryStartRequest{}, errors.New("invalid")
	}
	proposals := make([]domain.MarketQueryProposal, len(*input.QueryProposals))
	for i, value := range *input.QueryProposals {
		proposals[i] = domain.MarketQueryProposal{Locale: value.Locale, Query: value.Query, SourceText: value.SourceText, Derivation: value.Derivation, Rationale: value.Rationale}
	}
	request := domain.MarketDiscoveryStartRequest{Idea: *input.Idea, QueryProposals: proposals, Markets: *input.Markets, ConnectorIDs: *input.ConnectorIDs, Budget: *input.Budget}
	normalized, err := usecase.ValidateMarketDiscoveryStartRequest(request, registry)
	if err != nil {
		return "", domain.MarketDiscoveryStartRequest{}, err
	}
	return *input.IdempotencyKey, normalized, nil
}
