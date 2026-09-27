/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/improvement_synthesis.md
 * @layer L2
 * @updated 2026-09-02
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"bateia/01_core/domain"
)

func (s *MCPServer) WithImprovementSynthesis(starter domain.ImprovementSynthesisStarter) *MCPServer {
	if s != nil && starter != nil && s.improvementSynthesis == nil {
		s.improvementSynthesis = &improvementInterceptReader{server: s, starter: starter}
	}
	return s
}

type improvementInterceptReader struct {
	server  *MCPServer
	starter domain.ImprovementSynthesisStarter
}

func improvementToolsList() (map[string]interface{}, error) {
	result, err := marketToolsList()
	if err != nil {
		return nil, err
	}
	tools, ok := result["tools"].([]interface{})
	if !ok || len(tools) != 9 {
		return nil, errors.New("market tool list invalid")
	}
	artifactTool, ok := tools[5].(map[string]interface{})
	if !ok {
		return nil, errors.New("artifact tool invalid")
	}
	schema := artifactTool["inputSchema"].(map[string]interface{})
	properties := schema["properties"].(map[string]interface{})
	artifact := properties["artifact"].(map[string]interface{})
	artifact["enum"] = []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance", "connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps", "mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance"}
	tools = append(tools, map[string]interface{}{
		"name":        "bateia_start_improvement_synthesis",
		"description": "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis",
		"inputSchema": improvementStartSchema(),
	})
	return map[string]interface{}{"tools": tools}, nil
}

func improvementStartSchema() map[string]interface{} {
	controlledText := func(max int) map[string]interface{} {
		return map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(max), "pattern": `^[^\u0000-\u001F\u007F-\u009F]+$`}
	}
	textArray := func(max int) map[string]interface{} {
		return map[string]interface{}{"type": "array", "maxItems": float64(max), "items": controlledText(512)}
	}
	ideaTextArray := func(max int) map[string]interface{} {
		return map[string]interface{}{"type": "array", "maxItems": float64(max), "items": controlledText(200)}
	}
	idea := map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"problem": controlledText(1000), "users": ideaTextArray(20), "capabilities": ideaTextArray(20),
			"constraints": ideaTextArray(20), "non_goals": ideaTextArray(20), "preferred_stack": ideaTextArray(20),
		}, "required": []string{"problem", "capabilities"},
	}
	idea["properties"].(map[string]interface{})["capabilities"].(map[string]interface{})["minItems"] = float64(1)
	return map[string]interface{}{
		"type": "object", "additionalProperties": false,
		"properties": map[string]interface{}{
			"idempotency_key": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(128), "pattern": `^[\u0020-\u002E\u0030-\u005B\u005D-\u007E]+$`},
			"idea":            idea,
			"source_artifacts": map[string]interface{}{
				"type": "array", "minItems": float64(1), "maxItems": float64(12),
				"items": map[string]interface{}{
					"type": "object", "additionalProperties": false,
					"properties": map[string]interface{}{
						"run_id":   map[string]interface{}{"type": "string", "pattern": `^run_[0-9a-f]{32}$`},
						"run_kind": map[string]interface{}{"type": "string", "enum": []string{"RECOMMENDATION", "MARKET_DISCOVERY"}},
						"artifact": map[string]interface{}{"type": "string", "enum": []string{"mechanism_catalog", "mechanism_gaps", "market_candidates", "market_snapshot"}},
						"sha256":   map[string]interface{}{"type": "string", "pattern": `^[0-9a-f]{64}$`},
					}, "required": []string{"run_id", "run_kind", "artifact", "sha256"},
				},
			},
			"proposals": map[string]interface{}{
				"type": "array", "maxItems": float64(20),
				"items": map[string]interface{}{
					"type": "object", "additionalProperties": false,
					"properties": map[string]interface{}{
						"proposal_id": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(128), "pattern": `^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`},
						"title":       controlledText(512), "target_property": controlledText(512), "expected_effect": controlledText(512),
						"mechanism_refs":           map[string]interface{}{"type": "array", "maxItems": float64(4), "items": map[string]interface{}{"type": "string", "pattern": `^mech_[0-9a-f]{64}$`}},
						"product_observation_refs": map[string]interface{}{"type": "array", "maxItems": float64(20), "items": map[string]interface{}{"type": "string", "minLength": float64(1), "maxLength": float64(256), "pattern": `^[A-Za-z0-9][A-Za-z0-9._:-]{0,255}$`}},
						"assumptions":              textArray(20), "derivation": map[string]interface{}{"type": "string", "const": "EXTERNAL_AGENT_PROPOSAL"}, "rationale": controlledText(1000),
					}, "required": []string{"proposal_id", "title", "target_property", "expected_effect", "mechanism_refs", "product_observation_refs", "assumptions", "derivation", "rationale"},
				},
			},
			"policy": map[string]interface{}{
				"type": "object", "additionalProperties": false,
				"properties": map[string]interface{}{
					"max_hypotheses": integerRange(1, 12), "max_mechanisms_per_hypothesis": integerRange(1, 4), "require_counterevidence_search": map[string]interface{}{"type": "boolean", "const": true},
				}, "required": []string{"max_hypotheses", "max_mechanisms_per_hypothesis", "require_counterevidence_search"},
			},
		}, "required": []string{"idempotency_key", "idea", "source_artifacts", "proposals", "policy"},
	}
}

type improvementStartInput struct {
	IdempotencyKey  *string                     `json:"idempotency_key"`
	Idea            *domain.IdeaSpec            `json:"idea"`
	SourceArtifacts *[]improvementSourceInput   `json:"source_artifacts"`
	Proposals       *[]improvementProposalInput `json:"proposals"`
	Policy          *improvementPolicyInput     `json:"policy"`
}

type improvementSourceInput struct {
	RunID    *domain.RunID           `json:"run_id"`
	RunKind  *domain.RunKind         `json:"run_kind"`
	Artifact *domain.RunArtifactName `json:"artifact"`
	SHA256   *string                 `json:"sha256"`
}

type improvementProposalInput struct {
	ProposalID             *string                      `json:"proposal_id"`
	Title                  *string                      `json:"title"`
	TargetProperty         *string                      `json:"target_property"`
	ExpectedEffect         *string                      `json:"expected_effect"`
	MechanismRefs          *[]string                    `json:"mechanism_refs"`
	ProductObservationRefs *[]string                    `json:"product_observation_refs"`
	Assumptions            *[]string                    `json:"assumptions"`
	Derivation             *domain.HypothesisDerivation `json:"derivation"`
	Rationale              *string                      `json:"rationale"`
}

type improvementPolicyInput struct {
	MaxHypotheses                *int  `json:"max_hypotheses"`
	MaxMechanismsPerHypothesis   *int  `json:"max_mechanisms_per_hypothesis"`
	RequireCounterEvidenceSearch *bool `json:"require_counterevidence_search"`
}

func decodeImprovementStart(raw json.RawMessage) (string, domain.ImprovementSynthesisRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input improvementStartInput
	if decoder.Decode(&input) != nil || decoder.Decode(&struct{}{}) != io.EOF || input.IdempotencyKey == nil || input.Idea == nil || input.SourceArtifacts == nil || input.Proposals == nil || input.Policy == nil || !validDurableKey(*input.IdempotencyKey) {
		return "", domain.ImprovementSynthesisRequest{}, errors.New("invalid")
	}
	if input.Policy.MaxHypotheses == nil || input.Policy.MaxMechanismsPerHypothesis == nil || input.Policy.RequireCounterEvidenceSearch == nil {
		return "", domain.ImprovementSynthesisRequest{}, errors.New("invalid")
	}
	sources := make([]domain.SourceArtifactRef, len(*input.SourceArtifacts))
	for index, source := range *input.SourceArtifacts {
		if source.RunID == nil || source.RunKind == nil || source.Artifact == nil || source.SHA256 == nil {
			return "", domain.ImprovementSynthesisRequest{}, errors.New("invalid")
		}
		sources[index] = domain.SourceArtifactRef{RunID: *source.RunID, RunKind: *source.RunKind, Artifact: *source.Artifact, SHA256: *source.SHA256}
	}
	proposals := make([]domain.ImprovementProposal, len(*input.Proposals))
	for index, proposal := range *input.Proposals {
		if proposal.ProposalID == nil || proposal.Title == nil || proposal.TargetProperty == nil || proposal.ExpectedEffect == nil || proposal.MechanismRefs == nil || proposal.ProductObservationRefs == nil || proposal.Assumptions == nil || proposal.Derivation == nil || proposal.Rationale == nil {
			return "", domain.ImprovementSynthesisRequest{}, errors.New("invalid")
		}
		proposals[index] = domain.ImprovementProposal{ProposalID: *proposal.ProposalID, Title: *proposal.Title, TargetProperty: *proposal.TargetProperty, ExpectedEffect: *proposal.ExpectedEffect, MechanismRefs: *proposal.MechanismRefs, ProductObservationRefs: *proposal.ProductObservationRefs, Assumptions: *proposal.Assumptions, Derivation: *proposal.Derivation, Rationale: *proposal.Rationale}
	}
	policy := domain.ImprovementSynthesisPolicy{MaxHypotheses: *input.Policy.MaxHypotheses, MaxMechanismsPerHypothesis: *input.Policy.MaxMechanismsPerHypothesis, RequireCounterEvidenceSearch: *input.Policy.RequireCounterEvidenceSearch}
	return *input.IdempotencyKey, domain.ImprovementSynthesisRequest{Idea: *input.Idea, SourceArtifacts: sources, Proposals: proposals, Policy: policy}, nil
}

func (r *improvementInterceptReader) handleStart(ctx context.Context, id interface{}, raw json.RawMessage) {
	key, request, err := decodeImprovementStart(raw)
	if err != nil {
		r.server.sendError(id, -32602, "Invalid params")
		return
	}
	receipt, err := r.starter.StartImprovementSynthesis(ctx, key, request)
	if err != nil {
		code, message := durablePublicError(err)
		if errors.Is(err, domain.ErrInvalidRunInput) {
			code, message = -32602, "Invalid params"
		} else if errors.Is(err, domain.ErrSourceBasisUnavailable) {
			code, message = -32028, "SOURCE_BASIS_UNAVAILABLE"
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

func improvementArtifactRequested(raw json.RawMessage) bool {
	var value struct {
		Artifact string `json:"artifact"`
	}
	if json.Unmarshal(raw, &value) != nil {
		return false
	}
	for _, name := range []string{"mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance"} {
		if value.Artifact == name {
			return true
		}
	}
	return false
}

func decodeImprovementArtifact(raw json.RawMessage) (domain.ArtifactReadRequest, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input struct {
		RunID    *string `json:"run_id"`
		Artifact *string `json:"artifact"`
		Cursor   *string `json:"cursor"`
		Limit    *int    `json:"limit"`
	}
	if decoder.Decode(&input) != nil || decoder.Decode(&struct{}{}) != io.EOF || input.RunID == nil || input.Artifact == nil || input.Cursor == nil || input.Limit == nil || !durableRunIDPattern.MatchString(*input.RunID) || *input.Limit < 1 || *input.Limit > 50 || len(*input.Cursor) > 512 || (*input.Cursor != "" && !durableCursorPattern.MatchString(*input.Cursor)) || !improvementArtifactRequested(raw) {
		return domain.ArtifactReadRequest{}, errors.New("invalid")
	}
	return domain.ArtifactReadRequest{RunID: domain.RunID(*input.RunID), Artifact: domain.RunArtifactName(*input.Artifact), Cursor: *input.Cursor, Limit: *input.Limit}, nil
}

func (r *improvementInterceptReader) handleArtifact(ctx context.Context, id interface{}, raw json.RawMessage) {
	if r.server.durableRuns == nil || r.server.durableRuns.durable == nil {
		r.server.sendError(id, -32603, "Internal error")
		return
	}
	request, err := decodeImprovementArtifact(raw)
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
