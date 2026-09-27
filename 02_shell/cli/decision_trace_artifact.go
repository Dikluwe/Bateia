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
const decisionTraceArtifactContractVersion = "decision-trace-projection-v2"

func buildDecisionTraceStage(root string, graph domain.DecisionTraceGraph, index domain.DecisionTraceExplanation, versions map[string]string) (domain.IdeaRunStage, error) {
	inputs := map[string]string{"rule_versions": hashJSON(versions)}
	for _, path := range []string{"idea.json", "reuse-decision.json"} {
		hash, err := fileHash(filepath.Join(root, path))
		if err != nil {
			return domain.IdeaRunStage{}, fmt.Errorf("hash de entrada trace: %w", err)
		}
		inputs[path] = hash
	}
	outputs := []domain.IdeaRunArtifact{}
	for _, path := range []string{"trace-graph.json", "trace-index.json"} {
		hash, err := fileHash(filepath.Join(root, path))
		if err != nil {
			return domain.IdeaRunStage{}, fmt.Errorf("hash de artefato trace: %w", err)
		}
		outputs = append(outputs, domain.IdeaRunArtifact{Path: path, Hash: hash})
	}
	locations := 0
	for _, entry := range index.Entries {
		locations += len(entry.Locations)
	}
	return domain.IdeaRunStage{
		Name: "trace", ContractVersion: decisionTraceArtifactContractVersion,
		Status: domain.IdeaStageComplete, InputHashes: inputs,
		Outputs: outputs,
		Counts:  map[string]int64{"nodes": int64(len(graph.Nodes)), "edges": int64(len(graph.Edges)), "entries": int64(len(index.Entries)), "locations": int64(locations)},
		Message: "Dependências registradas; hash não autentica fatos, autoria ou completude. Não altera o veredito.",
	}, nil
}

// A partial sidecar pair must not be advertised as a completed projection.
func markDecisionTraceFailed(manifest *domain.IdeaRunManifest, err error) {
	manifest.Status = domain.IdeaRunFailed
	manifest.Failures = append(manifest.Failures, domain.ExecutionFailure{Scope: "run", Stage: "trace", Message: err.Error()})
	manifest.Stages = append(manifest.Stages, domain.IdeaRunStage{
		Name: "trace", ContractVersion: decisionTraceArtifactContractVersion,
		Status:      domain.IdeaStageFailed,
		InputHashes: map[string]string{}, Outputs: []domain.IdeaRunArtifact{}, Counts: map[string]int64{},
		Message: "Projeção incompleta; artefatos trace anteriores ou parciais não devem ser consumidos como execução atual.",
	})
}
