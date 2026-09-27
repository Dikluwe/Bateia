/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/advanced_intake.md
 * @layer L2
 * @updated 2026-09-24
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// WithAdvancedIntake adds the structured intake protocol without changing the
// three legacy tools exposed by NewMCPServer.
func (s *MCPServer) WithAdvancedIntake() *MCPServer {
	if s != nil && s.advancedIntake == nil {
		s.advancedIntake = &advancedIntakeInterceptReader{server: s}
	}
	return s
}

type advancedIntakeInterceptReader struct{ server *MCPServer }

func (r *advancedIntakeInterceptReader) handle(id interface{}, raw json.RawMessage) {
	var arguments struct {
		Input   domain.RecommendationIntakeDraft `json:"input"`
		Answers []domain.InformationAnswer       `json:"answers,omitempty"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if containsJSONNull(raw) || decoder.Decode(&arguments) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		r.server.sendError(id, -32602, "Invalid params")
		return
	}
	intake := usecase.PrepareRecommendationIntake(arguments.Input)
	if len(arguments.Answers) > 0 {
		var err error
		intake, err = usecase.ApplyInformationAnswers(intake, arguments.Answers)
		if err != nil {
			r.server.sendError(id, -32602, "Invalid params")
			return
		}
	}
	encoded, err := json.Marshal(intake)
	if err != nil {
		r.server.sendError(id, -32603, "Internal error")
		return
	}
	r.server.sendResult(id, map[string]interface{}{"content": []map[string]interface{}{{"type": "text", "text": string(encoded)}}})
}

func containsJSONNull(raw json.RawMessage) bool {
	var value interface{}
	if json.Unmarshal(raw, &value) != nil {
		return false
	}
	var walk func(interface{}) bool
	walk = func(current interface{}) bool {
		if current == nil {
			return true
		}
		switch typed := current.(type) {
		case []interface{}:
			for _, item := range typed {
				if walk(item) {
					return true
				}
			}
		case map[string]interface{}:
			for _, item := range typed {
				if walk(item) {
					return true
				}
			}
		}
		return false
	}
	return walk(value)
}

func baseToolsList() (map[string]interface{}, error) {
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
	return result, nil
}

func advancedIntakeTool() map[string]interface{} {
	text := map[string]interface{}{"type": "string"}
	texts := map[string]interface{}{"type": "array", "items": text}
	optional := func(properties map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"type": "object", "additionalProperties": false, "properties": properties}
	}
	input := optional(map[string]interface{}{
		"idea": optional(map[string]interface{}{"problem": text, "users": texts, "capabilities": texts, "qualities": texts, "operating_conditions": texts, "constraints": texts, "non_goals": texts, "preferred_stack": texts}),
		"policy": optional(map[string]interface{}{
			"schema_version": text, "allowed_licenses": texts, "require_local": map[string]interface{}{"type": "boolean"}, "allow_archived": map[string]interface{}{"type": "boolean"},
			"max_components": map[string]interface{}{"type": "integer"}, "extend_min_coverage": map[string]interface{}{"type": "number"}, "integration_min_value": map[string]interface{}{"type": "number"},
		}),
		"target": optional(map[string]interface{}{"name": text}),
		"search": optional(map[string]interface{}{"strategy": text, "language": text, "locales": texts, "max_per_query": map[string]interface{}{"type": "integer"}, "max_candidates": map[string]interface{}{"type": "integer"}}),
		"execution": optional(map[string]interface{}{
			"mode": map[string]interface{}{"type": "string", "enum": []string{"live", "local", "replay"}}, "provider": text, "local_paths": texts, "replay_directory": text, "output_directory": text,
			"resume": map[string]interface{}{"type": "boolean"}, "timeout_seconds": map[string]interface{}{"type": "integer"}, "allow_candidate_execution": map[string]interface{}{"type": "boolean"},
			"evidence_rules_ref": text, "screening_model_ref": text, "idempotency_key": text,
		}),
	})
	answer := map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"need_id": text, "answer_kind": map[string]interface{}{"type": "string", "enum": []string{"TEXT", "STRING_LIST", "BOOLEAN"}}, "value": map[string]interface{}{},
			"authority": map[string]interface{}{"const": "USER"}, "source_ref": text, "confirmed": map[string]interface{}{"const": true}, "supersedes": text,
		},
		"required": []string{"need_id", "answer_kind", "value", "authority", "confirmed"},
	}
	return map[string]interface{}{
		"name": "bateia_prepare_recommendation", "description": "Prepara entrada avançada, proveniência e perguntas tipadas sem executar a recomendação",
		"inputSchema": map[string]interface{}{
			"type": "object", "additionalProperties": false,
			"properties": map[string]interface{}{"input": input, "answers": map[string]interface{}{"type": "array", "items": answer}},
			"required":   []string{"input"},
		},
	}
}
