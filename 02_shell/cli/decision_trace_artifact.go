/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_artifact.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"fmt"
	"path/filepath"

	"bateia/01_core/domain"
)

// The trace is an optional projection, never a causal input or checkpoint gate.
func buildDecisionTraceStage(root string, graph domain.DecisionTraceGraph, versions map[string]string) (domain.IdeaRunStage, error) {
	inputs := map[string]string{"rule_versions": hashJSON(versions)}
	for _, path := range []string{"idea.json", "reuse-decision.json"} {
		hash, err := fileHash(filepath.Join(root, path))
		if err != nil {
			return domain.IdeaRunStage{}, fmt.Errorf("hash de entrada trace: %w", err)
		}
		inputs[path] = hash
	}
	hash, err := fileHash(filepath.Join(root, "trace-graph.json"))
	if err != nil {
		return domain.IdeaRunStage{}, fmt.Errorf("hash de artefato trace: %w", err)
	}
	return domain.IdeaRunStage{
		Name: "trace", ContractVersion: domain.DecisionTraceSchemaVersion,
		Status: domain.IdeaStageComplete, InputHashes: inputs,
		Outputs: []domain.IdeaRunArtifact{{Path: "trace-graph.json", Hash: hash}},
		Counts:  map[string]int64{"nodes": int64(len(graph.Nodes)), "edges": int64(len(graph.Edges))},
		Message: "Dependências registradas; hash não autentica fatos, autoria ou completude. Não altera o veredito.",
	}, nil
}
