/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/request_dispatch.md
 * @layer L2
 * @updated 2026-10-06
 */
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Limits apply to the original transport, before any extension or JSON object
// materialization. Oversized frames terminate the session without draining an
// unbounded peer stream. No extension wraps or reads the transport anymore.
const maxMCPMessageBytes = 64 << 10
const maxMCPJSONDepth = 32
const maxMCPJSONTokens = 8192

func boundedMCPScanner(reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), maxMCPMessageBytes+1)
	return scanner
}

func validateMCPStructure(line []byte) error {
	if len(line) > maxMCPMessageBytes {
		return fmt.Errorf("message too large")
	}
	decoder := json.NewDecoder(bytes.NewReader(line))
	depth := 0
	for tokens := 0; ; tokens++ {
		token, err := decoder.Token()
		if err != nil {
			return nil
		} // Envelope decoding retains existing parse errors.
		if tokens >= maxMCPJSONTokens {
			return fmt.Errorf("too many JSON tokens")
		}
		if delimiter, ok := token.(json.Delim); ok {
			switch delimiter {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}
			if depth > maxMCPJSONDepth {
				return fmt.Errorf("JSON nesting too deep")
			}
		}
	}
}

func (s *MCPServer) extensionToolsList() (map[string]interface{}, error) {
	var result map[string]interface{}
	var err error
	switch {
	case s.improvementSynthesis != nil:
		result, err = s.improvementToolsList()
	case s.marketDiscovery != nil:
		result, err = s.marketDiscovery.marketToolsList()
	case s.durableRuns != nil:
		result, err = durableToolsList()
	default:
		result, err = baseToolsList()
	}
	if err != nil {
		return nil, err
	}
	if s.advancedIntake != nil {
		result["tools"] = append(result["tools"].([]interface{}), advancedIntakeTool())
	}
	return result, nil
}

// Called only after the single budget/context guard in handleRequest. The
// registration order cannot hide tools or bypass the guard.
func (s *MCPServer) handleExtensionRequest(ctx context.Context, request JSONRPCRequest) bool {
	if request.Method == "tools/list" && (s.advancedIntake != nil || s.durableRuns != nil || s.marketDiscovery != nil || s.improvementSynthesis != nil) {
		result, err := s.extensionToolsList()
		if err != nil {
			s.sendError(request.ID, -32603, "Internal error")
		} else {
			s.sendResult(request.ID, result)
		}
		return true
	}
	if request.Method != "tools/call" {
		return false
	}
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if json.Unmarshal(request.Params, &params) != nil {
		return false
	}
	switch params.Name {
	case "bateia_prepare_recommendation":
		if s.advancedIntake != nil {
			s.advancedIntake.handle(request.ID, params.Arguments)
			return true
		}
	case "bateia_start_improvement_synthesis":
		if s.improvementSynthesis != nil {
			s.improvementSynthesis.handleStart(ctx, request.ID, params.Arguments)
			return true
		}
	case "bateia_list_connectors", "bateia_start_market_discovery":
		if s.marketDiscovery != nil {
			s.marketDiscovery.handleMarket(ctx, request.ID, params.Name, params.Arguments)
			return true
		}
	case "bateia_get_artifact":
		if s.improvementSynthesis != nil && improvementArtifactRequested(params.Arguments) {
			s.improvementSynthesis.handleArtifact(ctx, request.ID, params.Arguments)
			return true
		}
		if s.marketDiscovery != nil && s.durableRuns != nil && marketArtifactRequested(params.Arguments) {
			s.marketDiscovery.handleMarketArtifact(ctx, request.ID, params.Arguments)
			return true
		}
		if s.durableRuns != nil {
			s.durableRuns.handle(ctx, request.ID, params.Name, params.Arguments)
			return true
		}
	case "bateia_start_recommendation", "bateia_get_run", "bateia_cancel_run":
		if s.durableRuns != nil {
			s.durableRuns.handle(ctx, request.ID, params.Name, params.Arguments)
			return true
		}
	}
	return false
}

// Historical market contract returns Invalid params for recognizable malformed
// calls. It still consumes the same shared call budget before responding.
func (s *MCPServer) handleMalformedMarketRequest(ctx context.Context, id interface{}) {
	if ctx.Err() != nil {
		s.sendError(id, -32800, "Request cancelled")
		return
	}
	err := s.toolCallBudget.Dispatch("tools/call", func() *JSONRPCError { s.sendError(id, -32602, "Invalid params"); return nil })
	if err != nil {
		s.sendError(id, err.Code, err.Message)
	}
}
