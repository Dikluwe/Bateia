/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/durable_runs.md
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
	"regexp"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type DurableRecommendationRuns interface {
	Start(context.Context, string, domain.RecommendationRequest) (domain.RecommendationStartReceipt, error)
	GetRun(context.Context, domain.RunID) (domain.RecommendationRunView, error)
	GetArtifact(context.Context, domain.ArtifactReadRequest) (domain.ArtifactPage, error)
	Cancel(context.Context, domain.RunID) (domain.RecommendationRunView, error)
}

var (
	durableRunIDPattern  = regexp.MustCompile(`^run_[0-9a-f]{32}$`)
	durableCursorPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)
)

// WithDurableRuns registers tools behind the common request dispatcher.
func (s *MCPServer) WithDurableRuns(durable DurableRecommendationRuns) *MCPServer {
	if s != nil && durable != nil && s.durableRuns == nil {
		s.durableRuns = &durableInterceptReader{server: s, durable: durable}
	}
	return s
}

type durableInterceptReader struct {
	server  *MCPServer
	durable DurableRecommendationRuns
}

func durableToolsList() (map[string]interface{}, error) {
	writer := &bytes.Buffer{}
	probe := NewMCPServer(nil, strings.NewReader(""), writer)
	probe.handleRequest(context.Background(), JSONRPCRequest{JSONRPC: "2.0", ID: 1, Method: "tools/list"})
	var response JSONRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &response); err != nil || response.Error != nil {
		return nil, errors.New("legacy tool list unavailable")
	}
	encoded, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(encoded, &result); err != nil {
		return nil, err
	}
	tools, ok := result["tools"].([]interface{})
	if !ok || len(tools) != 3 {
		return nil, errors.New("legacy tool list invalid")
	}
	recommend := tools[2].(map[string]interface{})
	startSchemaBytes, _ := json.Marshal(recommend["inputSchema"])
	var startSchema map[string]interface{}
	_ = json.Unmarshal(startSchemaBytes, &startSchema)
	properties := startSchema["properties"].(map[string]interface{})
	properties["idempotency_key"] = map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(128), "pattern": `^[\u0020-\u002E\u0030-\u005B\u005D-\u007E]+$`}
	startSchema["required"] = []string{"idempotency_key", "idea"}
	tools = append(tools,
		map[string]interface{}{"name": "bateia_start_recommendation", "description": "Inicia recomendação durável assíncrona", "inputSchema": startSchema},
		map[string]interface{}{"name": "bateia_get_run", "description": "Consulta execução durável", "inputSchema": runIDSchema()},
		map[string]interface{}{"name": "bateia_get_artifact", "description": "Lê artefato seguro paginado", "inputSchema": artifactSchema()},
		map[string]interface{}{"name": "bateia_cancel_run", "description": "Solicita cancelamento da execução", "inputSchema": runIDSchema()},
	)
	return map[string]interface{}{"tools": tools}, nil
}

func runIDSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{"run_id": map[string]interface{}{"type": "string", "pattern": "^run_[0-9a-f]{32}$"}},
		"required":   []string{"run_id"},
	}
}

func artifactSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"run_id":   map[string]interface{}{"type": "string", "pattern": "^run_[0-9a-f]{32}$"},
			"artifact": map[string]interface{}{"type": "string", "enum": []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance"}},
			"cursor":   map[string]interface{}{"anyOf": []interface{}{map[string]interface{}{"const": ""}, map[string]interface{}{"type": "string", "maxLength": float64(512), "pattern": `^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`}}},
			"limit":    map[string]interface{}{"type": "integer", "minimum": float64(1), "maximum": float64(50)},
		},
		"required": []string{"run_id", "artifact", "cursor", "limit"},
	}
}

func (r *durableInterceptReader) handle(ctx context.Context, id interface{}, name string, arguments json.RawMessage) {
	var value interface{}
	var err error
	switch name {
	case "bateia_start_recommendation":
		key, request, decodeErr := decodeDurableStart(arguments)
		if decodeErr != nil {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
		value, err = r.durable.Start(ctx, key, request)
	case "bateia_get_run":
		runID, decodeErr := decodeDurableRunID(arguments)
		if decodeErr != nil {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
		value, err = r.durable.GetRun(ctx, runID)
	case "bateia_cancel_run":
		runID, decodeErr := decodeDurableRunID(arguments)
		if decodeErr != nil {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
		value, err = r.durable.Cancel(ctx, runID)
	case "bateia_get_artifact":
		read, decodeErr := decodeDurableArtifact(arguments)
		if decodeErr != nil {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
		value, err = r.durable.GetArtifact(ctx, read)
	}
	if err != nil {
		code, message := durablePublicError(err)
		r.server.sendError(id, code, message)
		return
	}
	encoded, err := marshalDurableValue(value)
	if err != nil {
		r.server.sendError(id, -32603, "Internal error")
		return
	}
	r.server.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
}

type durableStartArguments struct {
	IdempotencyKey *string                        `json:"idempotency_key"`
	Idea           *recommendationIdeaArguments   `json:"idea"`
	Policy         *recommendationPolicyArguments `json:"policy"`
	Target         *recommendationTargetArguments `json:"target"`
	Search         *recommendationSearchArguments `json:"search"`
}

func decodeDurableStart(raw json.RawMessage) (string, domain.RecommendationRequest, error) {
	if err := rejectRecommendationNulls(raw); err != nil {
		return "", domain.RecommendationRequest{}, errors.New("invalid")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input durableStartArguments
	if err := decoder.Decode(&input); err != nil || decoder.Decode(&struct{}{}) != io.EOF || input.IdempotencyKey == nil || input.Idea == nil {
		return "", domain.RecommendationRequest{}, errors.New("invalid")
	}
	key := *input.IdempotencyKey
	if !validDurableKey(key) {
		return "", domain.RecommendationRequest{}, errors.New("invalid")
	}
	request := materializeRecommendationArguments(recommendationArguments{Idea: input.Idea, Policy: input.Policy, Target: input.Target, Search: input.Search})
	if _, err := usecase.ValidateRecommendationRequest(request); err != nil {
		return "", domain.RecommendationRequest{}, errors.New("invalid")
	}
	return key, request, nil
}

func validDurableKey(key string) bool {
	if len(key) < 1 || len(key) > 128 || strings.Contains(key, "..") {
		return false
	}
	for _, value := range []byte(key) {
		if value < 0x20 || value > 0x7e || value == '/' || value == '\\' {
			return false
		}
	}
	return true
}

func decodeDurableRunID(raw json.RawMessage) (domain.RunID, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input struct {
		RunID *string `json:"run_id"`
	}
	if err := decoder.Decode(&input); err != nil || decoder.Decode(&struct{}{}) != io.EOF || input.RunID == nil || !durableRunIDPattern.MatchString(*input.RunID) {
		return "", errors.New("invalid")
	}
	return domain.RunID(*input.RunID), nil
}

func decodeDurableArtifact(raw json.RawMessage) (domain.ArtifactReadRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input struct {
		RunID    *string `json:"run_id"`
		Artifact *string `json:"artifact"`
		Cursor   *string `json:"cursor"`
		Limit    *int    `json:"limit"`
	}
	if err := decoder.Decode(&input); err != nil || decoder.Decode(&struct{}{}) != io.EOF || input.RunID == nil || input.Artifact == nil || input.Cursor == nil || input.Limit == nil || !durableRunIDPattern.MatchString(*input.RunID) || *input.Limit < 1 || *input.Limit > 50 || len(*input.Cursor) > 512 || (*input.Cursor != "" && !durableCursorPattern.MatchString(*input.Cursor)) {
		return domain.ArtifactReadRequest{}, errors.New("invalid")
	}
	name := domain.RunArtifactName(*input.Artifact)
	valid := false
	for _, candidate := range []domain.RunArtifactName{domain.ArtifactRequest, domain.ArtifactRecommendation, domain.ArtifactDecision, domain.ArtifactAlternatives, domain.ArtifactFailures, domain.ArtifactProvenance} {
		valid = valid || name == candidate
	}
	if !valid {
		return domain.ArtifactReadRequest{}, errors.New("invalid")
	}
	return domain.ArtifactReadRequest{RunID: domain.RunID(*input.RunID), Artifact: name, Cursor: *input.Cursor, Limit: *input.Limit}, nil
}

func marshalDurableValue(value interface{}) ([]byte, error) {
	page, ok := value.(domain.ArtifactPage)
	if !ok {
		return json.Marshal(value)
	}
	entries := make([]json.RawMessage, len(page.Entries))
	for index, entry := range page.Entries {
		if !json.Valid(entry) {
			return nil, errors.New("invalid artifact entry")
		}
		entries[index] = json.RawMessage(entry)
	}
	return json.Marshal(struct {
		SchemaVersion string                 `json:"schema_version"`
		RunID         domain.RunID           `json:"run_id"`
		Artifact      domain.RunArtifactName `json:"artifact"`
		ArtifactHash  string                 `json:"artifact_hash"`
		Entries       []json.RawMessage      `json:"entries"`
		NextCursor    *string                `json:"next_cursor"`
	}{page.SchemaVersion, page.RunID, page.Artifact, page.ArtifactHash, entries, page.NextCursor})
}

func durablePublicError(err error) (int, string) {
	for _, item := range []struct {
		target  error
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
	} {
		if errors.Is(err, item.target) {
			return item.code, item.message
		}
	}
	return -32603, "Internal error"
}
