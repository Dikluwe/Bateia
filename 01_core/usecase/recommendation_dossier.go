/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_dossier.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
)

type DossierSource string

const (
	DossierOutputSource DossierSource = "output"
	DossierReplaySource DossierSource = "replay"
)

type RecommendationDossierFiles interface {
	EvidenceCheckpointReader
	AlternativeIdentityReader
	Prepare(bool) error
	ReplayEvidenceIdentity(string) (DossierReplayEvidenceIdentity, error)
	CaptureLocal(context.Context, []string, bool) (domain.CapturedWorkspace, error)
	SnapshotProvider(domain.ContentProvider, []domain.SnapshotBinding) (domain.ContentProvider, error)
	HasArtifact(string) bool
	ReadJSON(DossierSource, string, any) error
	WriteJSON(string, any) error
	WriteDecision(domain.ReuseDecisionReport, []domain.EvidenceBundle) error
	WriteAlternatives(domain.SoftwareAlternativeReport) error
	WriteTrace(any, string, string) error
	Announce(RecommendationDossierOutput)
}
type RecommendationDossierInput struct {
	Request            domain.IdeaRunRequest
	Target             domain.SoftwareAlternativeTarget
	Rules              domain.CapabilityRuleCatalog
	CustomRules        bool
	ScreeningModel     *domain.ScreeningModel
	ScreeningModelHash string
	Trace              bool
	Dependencies       IdeaRunDependencies
	Recommendation     *RecommendationService
}
type RecommendationDossierOutput struct {
	Manifest              domain.IdeaRunManifest
	Result                domain.IdeaRunResult
	Alternatives          domain.SoftwareAlternativeReport
	Checkpoint            *EvidenceCheckpointAssessment
	StaleTrace, Published bool
}
type dossierWorkspace struct {
	domain.CapturedWorkspace
	domain.CapturedWorkspaceView
}
type dossierIdentityLookup func(string) []domain.RepositoryIdentity

func (f dossierIdentityLookup) RepositoryIdentities(url string) []domain.RepositoryIdentity {
	return f(url)
}
func RunRecommendationDossier(ctx context.Context, input RecommendationDossierInput, files RecommendationDossierFiles) (out RecommendationDossierOutput, returnErr error) {
	req := input.Request
	locales, localeErr := canonicalizeSearchLocales(req.SearchLocales)
	if localeErr != nil {
		return out, localeErr
	}
	req.SearchLocales = locales
	idea, policy, target := req.Idea, req.Policy, input.Target
	targetRepository := target.Repository
	evidenceRules := input.Rules
	content := input.Dependencies.Content
	if files == nil {
		return out, fmt.Errorf("dossiê exige porta de arquivos")
	}
	if input.ScreeningModel != nil && !IsValidOptInScreeningModel(*input.ScreeningModel) {
		return out, fmt.Errorf("modelo de screening inválido")
	}

	catalogCoverage, err := AnalyzeEvidenceCatalogCoverage(idea, evidenceRules)
	if err != nil {
		return out, err
	}
	deps := input.Dependencies
	deps.EvidenceRules = evidenceRules
	if input.ScreeningModel != nil {
		deps.ScreeningModel = *input.ScreeningModel
	}
	var localSnapshot *dossierWorkspace
	if req.Mode == domain.IdeaRunLocal {
		var resource domain.CapturedWorkspace
		resource, err = files.CaptureLocal(ctx, req.LocalPaths, targetRepository != nil)
		if err == nil {
			localSnapshot = &dossierWorkspace{CapturedWorkspace: resource, CapturedWorkspaceView: resource.View()}
		}
		if err != nil {
			return out, err
		}
		defer func() {
			if cleanupErr := localSnapshot.Close(); cleanupErr != nil {
				if returnErr == nil {
					returnErr = cleanupErr
				} else {
					returnErr = fmt.Errorf("%w; %v", returnErr, cleanupErr)
				}
			}
		}()
		deps.LocalCandidates = localSnapshot.Candidates
		deps.Content, err = files.SnapshotProvider(content, localSnapshot.Bindings)
		if err != nil {
			return out, err
		}
	}
	evidenceHash, err := DossierEvidenceInputHash(req, deps.LocalCandidates, files)
	if err != nil {
		return out, err
	}
	if err := ctx.Err(); err != nil {
		return out, err
	}
	if err := files.Prepare(req.Resume); err != nil {
		return out, err
	}
	inputHashes := map[string]string{"idea": dossierHashJSON(idea), "policy": dossierHashJSON(policy), "mode": dossierHashBytes([]byte(req.Mode)), "evidence": evidenceHash, "evidence_rules": catalogCoverage.RulesHash, "target": dossierHashJSON(target)}
	if input.ScreeningModel != nil {
		inputHashes["screening_model"] = input.ScreeningModelHash
	}
	if len(req.SearchLocales) > 0 {
		inputHashes["search_locales"] = dossierHashJSON(req.SearchLocales)
	}
	if localSnapshot != nil && targetRepository != nil {
		inputHashes["local_metadata"] = dossierHashJSON(localSnapshot.Identities)
	}
	ruleVersions := map[string]string{"expansion": "1", "screening": "1", "evidence": catalogCoverage.RulesVersion, "decision": "0043", "alternatives": "0043", "analyzer": "0043"}
	if input.ScreeningModel != nil {
		ruleVersions["screening"] = "screening-opt-in-v1"
	}
	if req.Mode == domain.IdeaRunLocal {
		ruleVersions["local_snapshot"] = "tree-sha256-v2"
		ruleVersions["local_collection"] = "frozen-local-v1"
	}
	if req.Mode == domain.IdeaRunLocal || req.Mode == domain.IdeaRunLive {
		ruleVersions["implementation_search"] = "0061"
	}
	if req.Mode == domain.IdeaRunLive {
		ruleVersions["discovery_query"] = "live0056-v1"
	}
	var previous domain.IdeaRunManifest
	staleTrace := false
	if req.Resume && !input.Trace {
		for _, path := range []string{"trace-graph.json", "trace-index.json"} {
			exists := files.HasArtifact(path)
			staleTrace = staleTrace || exists
		}
	}
	canReuseEvidence := false
	if req.Resume {
		assessment := EvidenceCheckpointAssessment{Reason: CheckpointManifestUnavailable, Subject: "run.json"}
		if e := files.ReadJSON(DossierOutputSource, "run.json", &previous); e == nil {
			assessment = EvaluateEvidenceCheckpoint(previous, EvidenceCheckpointCriteria{EvidenceHash: evidenceHash, EvidenceRulesHash: catalogCoverage.RulesHash, RuleVersions: ruleVersions, ScreeningModelHash: input.ScreeningModelHash}, files)
		}
		out.Checkpoint = &assessment
		canReuseEvidence = assessment.Reusable
	}
	runIdentity := struct {
		EvidenceHash       string            `json:"evidence_hash"`
		PolicyHash         string            `json:"policy_hash"`
		TargetHash         string            `json:"target_hash"`
		ScreeningModelHash string            `json:"screening_model_hash,omitempty"`
		LocalMetadataHash  string            `json:"local_metadata_hash,omitempty"`
		RuleVersions       map[string]string `json:"rule_versions"`
	}{EvidenceHash: evidenceHash, PolicyHash: inputHashes["policy"], TargetHash: inputHashes["target"], ScreeningModelHash: input.ScreeningModelHash, LocalMetadataHash: inputHashes["local_metadata"], RuleVersions: ruleVersions}
	manifest := domain.IdeaRunManifest{SchemaVersion: "1", RunID: StableRunID(dossierHashJSON(runIdentity)), Mode: req.Mode, Status: domain.IdeaRunRunning, InputHashes: inputHashes, RuleVersions: ruleVersions, Stages: []domain.IdeaRunStage{}}
	if err = files.WriteJSON("run.json", manifest); err != nil {
		return out, err
	}
	source := DossierReplaySource
	if canReuseEvidence {
		source = DossierOutputSource
	}
	if req.Mode == domain.IdeaRunReplay || canReuseEvidence {
		var snap domain.DiscoverySnapshot
		var bundles []domain.EvidenceBundle
		if e := files.ReadJSON(source, "discovery-snapshot.json", &snap); e != nil {
			manifest.Status = domain.IdeaRunFailed
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "replay", Message: e.Error()}}
			_ = files.WriteJSON("run.json", manifest)
			return out, e
		}
		if e := files.ReadJSON(source, "evidence/bundles.json", &bundles); e != nil {
			manifest.Status = domain.IdeaRunInconclusive
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "replay", Message: e.Error()}}
			_ = files.WriteJSON("run.json", manifest)
			return out, e
		}
		deps.ReplaySnapshot = &snap
		deps.ReplayBundles = bundles
	}
	var result domain.IdeaRunResult
	var alternatives domain.SoftwareAlternativeReport
	var runErr error
	usesRecommendationService := req.Mode == domain.IdeaRunLive && targetRepository == nil && !canReuseEvidence && input.ScreeningModel == nil
	if usesRecommendationService {
		recommendationService, constructorErr := NewRecommendationServiceWithEvidenceRules(deps.Semantic, deps.Discovery, content, evidenceRules)
		if constructorErr != nil {
			return out, constructorErr
		}
		if !input.CustomRules && input.Recommendation != nil {
			recommendationService = input.Recommendation
		}
		execution, recommendErr := recommendationService.RecommendWithRun(ctx, domain.RecommendationRequest{
			Idea: idea, Policy: policy, Target: target,
			Search: domain.RecommendationSearch{Strategy: req.Strategy, Language: req.Language, Locales: req.SearchLocales, MaxPerQuery: req.MaxPerQuery, MaxCandidates: req.MaxCandidates},
		})
		result, alternatives, runErr = execution.Run, execution.Recommendation.Alternatives, recommendErr
	} else {
		result, runErr = RunIdeaPipeline(ctx, req, deps)
	}
	if runErr != nil {
		manifest.Status = domain.IdeaRunFailed
		manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "pipeline", Message: runErr.Error()}}
		_ = files.WriteJSON("run.json", manifest)
		return out, runErr
	}
	if !usesRecommendationService {
		alternativesStarted := domain.RuntimeNow(ctx)
		identitiesFor := files.RepositoryIdentities
		if localSnapshot != nil {
			identitiesFor = func(candidateURL string) []domain.RepositoryIdentity { return localSnapshot.Identities[candidateURL] }
		}
		alternativeDiscovery, alternativeAssessments := AlternativeInputsForTarget(req.Mode, result.Discovery, result.Assessment.Assessments, targetRepository, dossierIdentityLookup(identitiesFor))
		alternatives, err = RankSoftwareAlternatives(target, idea, alternativeDiscovery, alternativeAssessments, policy)
		result.StageDurations["alternatives"] = domain.ElapsedMillis(domain.RuntimeNow(ctx), alternativesStarted)
		if err != nil {
			manifest.Status = domain.IdeaRunFailed
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "alternatives", Message: err.Error()}}
			_ = files.WriteJSON("run.json", manifest)
			return out, err
		}
	}
	var traceGraph *domain.DecisionTraceGraph
	var traceIndex domain.DecisionTraceExplanation
	if input.Trace {
		buildInput := domain.DecisionTraceBuildInput{SchemaVersion: "1", Idea: &idea, Report: result.Decision, RuleVersions: ruleVersions}
		graph, traceErr := BuildDecisionTrace(buildInput)
		if traceErr == nil {
			traceIndex, traceErr = ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: buildInput})
		}
		graphPath, indexPath := "trace-graph.json", "trace-index.json"
		if traceErr == nil {
			traceErr = files.WriteTrace(graph, graphPath, indexPath)
		}
		if traceErr == nil {
			traceErr = files.WriteTrace(traceIndex, indexPath, graphPath)
		}
		if traceErr != nil {
			MarkRecommendationTraceFailed(&manifest, traceErr)
			_ = files.WriteJSON("run.json", manifest)
			return out, traceErr
		}
		traceGraph = &graph
	}
	artifacts := []struct {
		path  string
		value any
	}{{"idea.json", idea}, {"reuse-policy.json", policy}, {"expansion-plan.json", result.Expansion}, {"discovery-snapshot.json", result.Discovery}, {"screening-report.json", result.Screening}, {"evidence/catalog-coverage.json", result.CatalogCoverage}, {"evidence/plans.json", result.Plans}, {"evidence/bundles.json", result.Bundles}, {"assessment-report.json", result.Assessment}, {"reuse-decision.json", result.Decision}, {"alternatives.json", alternatives}}
	for _, a := range artifacts {
		if err = files.WriteJSON(a.path, a.value); err != nil {
			return out, err
		}
	}
	if err = files.WriteDecision(result.Decision, result.Bundles); err != nil {
		return out, err
	}
	if err = files.WriteAlternatives(alternatives); err != nil {
		return out, err
	}
	manifest.Status = IdeaRunStatusFor(result)
	manifest.DecisionPath = "reuse-decision.json"
	manifest.Failures = append(manifest.Failures, result.Failures...)
	manifest.Stages = BuildRecommendationRunStages(files, result, alternatives)
	stageHashFailure := false
	for _, stage := range manifest.Stages {
		if stage.Status != domain.IdeaStageFailed {
			continue
		}
		stageHashFailure = true
		manifest.Status = domain.IdeaRunFailed
		manifest.Failures = append(manifest.Failures, domain.ExecutionFailure{Scope: "run", Stage: stage.Name, Message: stage.Message})
	}
	if stageHashFailure {
		if err = files.WriteJSON("run.json", manifest); err != nil {
			return out, err
		}
		out.Manifest = manifest
		return out, fmt.Errorf("dossiê sem hashes obrigatórios para publicação")
	}
	if traceGraph != nil {
		stage, err := BuildRecommendationTraceStage(files, *traceGraph, traceIndex, ruleVersions)
		if err != nil {
			MarkRecommendationTraceFailed(&manifest, err)
			_ = files.WriteJSON("run.json", manifest)
			return out, err
		}
		manifest.Stages = append(manifest.Stages, stage)
	}
	if localSnapshot != nil {
		for i := range manifest.Stages {
			stage := &manifest.Stages[i]
			if stage.Name == "discover" || stage.Name == "collect" {
				stage.Counts["checkpoint_reused"] = 0
				if canReuseEvidence && (stage.Name == "discover" || len(deps.ReplayBundles) > 0) {
					stage.Counts["checkpoint_reused"] = 1
				}
			}
		}
	}
	if err = files.WriteJSON("run.json", manifest); err != nil {
		return out, err
	}
	out = RecommendationDossierOutput{Manifest: manifest, Result: result, Alternatives: alternatives, Checkpoint: out.Checkpoint, StaleTrace: staleTrace, Published: true}
	files.Announce(out)
	return out, nil
}
