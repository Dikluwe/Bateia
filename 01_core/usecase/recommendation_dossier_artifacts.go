/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_dossier_artifacts.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type evidenceInputs struct {
	Idea            domain.IdeaSpec                `json:"idea"`
	Mode            domain.IdeaRunMode             `json:"mode"`
	Strategy        string                         `json:"strategy"`
	Provider        string                         `json:"provider,omitempty"`
	Language        string                         `json:"language,omitempty"`
	SearchLocales   []string                       `json:"search_locales,omitempty"`
	MaxPerQuery     int                            `json:"max_per_query"`
	MaxCandidates   int                            `json:"max_candidates"`
	LocalCandidates []localEvidenceInput           `json:"local_candidates,omitempty"`
	Replay          *DossierReplayEvidenceIdentity `json:"replay,omitempty"`
}

type localEvidenceInput struct {
	URL      string `json:"url"`
	CommitID string `json:"commit_id"`
}

type DossierReplayEvidenceIdentity struct {
	Directory    string `json:"directory"`
	SnapshotHash string `json:"snapshot_hash"`
	BundlesHash  string `json:"bundles_hash"`
}

func DossierEvidenceInputHash(req domain.IdeaRunRequest, local []domain.RepositoryCandidate, files interface {
	ReplayEvidenceIdentity(string) (DossierReplayEvidenceIdentity, error)
}) (string, error) {
	inputs := evidenceInputs{Idea: req.Idea, Mode: req.Mode, Strategy: req.Strategy, Provider: req.Provider, Language: req.Language, SearchLocales: append([]string(nil), req.SearchLocales...), MaxPerQuery: req.MaxPerQuery, MaxCandidates: req.MaxCandidates}
	for _, candidate := range local {
		inputs.LocalCandidates = append(inputs.LocalCandidates, localEvidenceInput{URL: candidate.URL, CommitID: candidate.CommitID})
	}
	if req.Mode == domain.IdeaRunReplay {
		identity, err := files.ReplayEvidenceIdentity(req.ReplayDirectory)
		if err != nil {
			return "", err
		}
		inputs.Replay = &identity
	}
	return dossierHashJSON(inputs), nil
}

func BuildRecommendationRunStages(files EvidenceCheckpointReader, r domain.IdeaRunResult, alternatives ...domain.SoftwareAlternativeReport) []domain.IdeaRunStage {
	defs := []struct {
		name, path string
		count      int64
	}{{"expand", "expansion-plan.json", int64(len(r.Expansion.Expansions))}, {"discover", "discovery-snapshot.json", int64(len(r.Discovery.Candidates))}, {"screen", "screening-report.json", int64(len(r.Screening.Decisions))}, {"collect", "evidence/bundles.json", int64(len(r.Bundles))}, {"assess", "assessment-report.json", int64(len(r.Assessment.Assessments))}, {"decide", "reuse-decision.json", 1}}
	out := make([]domain.IdeaRunStage, 0, len(defs))
	ideaHash, _ := files.ArtifactSHA256("idea.json")
	prev := map[string]string{"idea.json": ideaHash}
	for _, d := range defs {
		hash, _ := files.ArtifactSHA256(d.path)
		if d.name == "decide" {
			policyHash, _ := files.ArtifactSHA256("reuse-policy.json")
			prev["reuse-policy.json"] = policyHash
		}
		out = append(out, domain.IdeaRunStage{Name: d.name, ContractVersion: "1", Status: domain.IdeaStageComplete, InputHashes: prev, Outputs: []domain.IdeaRunArtifact{{Path: d.path, Hash: hash}}, Counts: map[string]int64{"items": d.count}})
		out[len(out)-1].DurationMillis = r.StageDurations[d.name]
		prev = map[string]string{d.path: hash}
	}
	if len(alternatives) > 0 {
		jsonHash, _ := files.ArtifactSHA256("alternatives.json")
		markdownHash, _ := files.ArtifactSHA256("alternatives.md")
		assessmentHash, _ := files.ArtifactSHA256("assessment-report.json")
		discoveryHash, _ := files.ArtifactSHA256("discovery-snapshot.json")
		policyHash, _ := files.ArtifactSHA256("reuse-policy.json")
		out = append(out, domain.IdeaRunStage{
			Name:            "alternatives",
			ContractVersion: "1",
			Status:          domain.IdeaStageComplete,
			InputHashes: map[string]string{
				"assessment-report.json":  assessmentHash,
				"discovery-snapshot.json": discoveryHash,
				"reuse-policy.json":       policyHash,
				"target":                  dossierHashJSON(alternatives[0].Target),
			},
			Outputs:        []domain.IdeaRunArtifact{{Path: "alternatives.json", Hash: jsonHash}, {Path: "alternatives.md", Hash: markdownHash}},
			Counts:         map[string]int64{"items": int64(len(alternatives[0].Alternatives) + len(alternatives[0].Excluded))},
			DurationMillis: r.StageDurations["alternatives"],
		})
	}
	return out
}

// The trace is an optional projection, never a causal input or checkpoint gate.
const decisionTraceArtifactContractVersion = "decision-trace-projection-v2"

func BuildRecommendationTraceStage(files EvidenceCheckpointReader, graph domain.DecisionTraceGraph, index domain.DecisionTraceExplanation, versions map[string]string) (domain.IdeaRunStage, error) {
	inputs := map[string]string{"rule_versions": dossierHashJSON(versions)}
	for _, path := range []string{"idea.json", "reuse-decision.json"} {
		hash, err := files.ArtifactSHA256(path)
		if err != nil {
			return domain.IdeaRunStage{}, fmt.Errorf("hash de entrada trace: %w", err)
		}
		inputs[path] = hash
	}
	outputs := []domain.IdeaRunArtifact{}
	for _, path := range []string{"trace-graph.json", "trace-index.json"} {
		hash, err := files.ArtifactSHA256(path)
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
func MarkRecommendationTraceFailed(manifest *domain.IdeaRunManifest, err error) {
	manifest.Status = domain.IdeaRunFailed
	manifest.Failures = append(manifest.Failures, domain.ExecutionFailure{Scope: "run", Stage: "trace", Message: err.Error()})
	manifest.Stages = append(manifest.Stages, domain.IdeaRunStage{
		Name: "trace", ContractVersion: decisionTraceArtifactContractVersion,
		Status:      domain.IdeaStageFailed,
		InputHashes: map[string]string{}, Outputs: []domain.IdeaRunArtifact{}, Counts: map[string]int64{},
		Message: "Projeção incompleta; artefatos trace anteriores ou parciais não devem ser consumidos como execução atual.",
	})
}

func dossierHashJSON(value any) string { data, _ := json.Marshal(value); return dossierHashBytes(data) }
func dossierHashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
