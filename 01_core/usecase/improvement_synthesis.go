/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_synthesis.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"bateia/01_core/domain"
)

var synthesisSHA256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
var proposalIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var mechanismIDPattern = regexp.MustCompile(`^mech_[0-9a-f]{64}$`)
var observationIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,255}$`)

type resolvedSynthesisSources struct {
	mechanisms    map[string]domain.ArchitecturalMechanism
	mechanismGaps map[string]domain.MechanismGap
	observations  map[string]bool
	degraded      bool
}

func (s *RecommendationRunService) ConfigureImprovementSynthesis(catalogSHA256 string) *RecommendationRunService {
	if s != nil && catalogSHA256 == activeMechanismRulesSHA256 {
		s.improvementEnabled = true
		s.mechanismCatalogSHA256 = catalogSHA256
	}
	return s
}

func (s *RecommendationRunService) ConfigureMechanismExtraction(extractor domain.MechanismExtractor, catalog []byte, catalogSHA256 string, baseline *domain.ConventionalBaseline) *RecommendationRunService {
	if s == nil || extractor == nil {
		return s
	}
	if _, err := decodeMechanismCatalog(catalog, catalogSHA256); err != nil {
		return s
	}
	s.mechanismExtractor = extractor
	s.mechanismCatalog = append([]byte(nil), catalog...)
	s.mechanismCatalogSHA256 = catalogSHA256
	s.mechanismBaseline = baseline
	s.improvementEnabled = true
	return s
}

func (s *RecommendationRunService) StartImprovementSynthesis(ctx context.Context, key string, request domain.ImprovementSynthesisRequest) (domain.ImprovementSynthesisStartReceipt, error) {
	if s == nil || s.store == nil || s.queue == nil || s.ids == nil || !s.improvementEnabled || s.mechanismCatalogSHA256 != activeMechanismRulesSHA256 || !validIdempotencyKey(key) {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrInvalidRunInput
	}
	normalized, err := normalizeImprovementSynthesisRequest(request)
	if err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrInvalidRunInput
	}
	resolved, err := s.resolveSynthesisSources(ctx, normalized.SourceArtifacts)
	if err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, err
	}
	if err := validateProposalReferences(normalized.Proposals, resolved); err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrInvalidRunInput
	}
	requestHash, err := synthesisRequestHash(normalized)
	if err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrInvalidRunInput
	}
	keyDigest := digestString(key)
	if record, ok, resolveErr := s.store.ResolveIdempotency(ctx, keyDigest); resolveErr != nil {
		return domain.ImprovementSynthesisStartReceipt{}, resolveErr
	} else if ok {
		if record.RequestHash != requestHash {
			return domain.ImprovementSynthesisStartReceipt{}, domain.ErrIdempotencyConflict
		}
		run, getErr := s.store.Get(ctx, record.RunID)
		if getErr != nil {
			return domain.ImprovementSynthesisStartReceipt{}, getErr
		}
		return improvementReceiptFor(run), nil
	}
	reservation, err := s.queue.Reserve(ctx)
	if err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, err
	}
	consumed := false
	defer func() {
		if !consumed {
			reservation.Release()
		}
	}()
	id, err := s.ids.NewRunID()
	if err != nil || !runIDPattern.MatchString(string(id)) {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	run := domain.RecommendationRun{
		SchemaVersion: "1", RunID: id, RunKind: domain.RunKindImprovementSynthesis,
		ImprovementRequest: &normalized, RequestHash: requestHash, Lifecycle: domain.RunQueued,
		Version: 1, Stage: domain.RunStageQueued, Progress: domain.RecommendationRunProgress{Planned: 5}, Artifacts: []domain.RunArtifactIndex{},
	}
	actual, created, err := s.store.CreateOrResolve(ctx, run, domain.IdempotencyRecord{KeyDigest: keyDigest, RequestHash: requestHash, RunID: id})
	if err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, err
	}
	if !created {
		if actual.RequestHash != requestHash {
			return domain.ImprovementSynthesisStartReceipt{}, domain.ErrIdempotencyConflict
		}
		return improvementReceiptFor(actual), nil
	}
	if err := reservation.Commit(actual.RunID); err != nil {
		return domain.ImprovementSynthesisStartReceipt{}, domain.ErrRunStoreUnavailable
	}
	consumed = true
	return improvementReceiptFor(actual), nil
}

func improvementReceiptFor(run domain.RecommendationRun) domain.ImprovementSynthesisStartReceipt {
	return domain.ImprovementSynthesisStartReceipt{SchemaVersion: "1", RunID: run.RunID, RunKind: domain.RunKindImprovementSynthesis, Lifecycle: run.Lifecycle, Version: run.Version}
}

func normalizeImprovementSynthesisRequest(request domain.ImprovementSynthesisRequest) (domain.ImprovementSynthesisRequest, error) {
	if !validSynthesisText(request.Idea.Problem, 1000) || !validSynthesisTextSet(request.Idea.Users, 20, 200) || !validSynthesisTextSet(request.Idea.Capabilities, 20, 200) || !validSynthesisTextSet(request.Idea.Constraints, 20, 200) || !validSynthesisTextSet(request.Idea.NonGoals, 20, 200) || !validSynthesisTextSet(request.Idea.PreferredStack, 20, 200) {
		return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
	}
	request.Idea = normalizeIdea(request.Idea)
	if request.Idea.Problem == "" || len(request.Idea.Capabilities) == 0 || len(request.SourceArtifacts) < 1 || len(request.SourceArtifacts) > 12 || len(request.Proposals) > 20 || request.Policy.MaxHypotheses < 1 || request.Policy.MaxHypotheses > 12 || request.Policy.MaxMechanismsPerHypothesis < 1 || request.Policy.MaxMechanismsPerHypothesis > 4 || !request.Policy.RequireCounterEvidenceSearch {
		return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
	}
	request.Idea.Users = canonicalSet(request.Idea.Users)
	request.Idea.Capabilities = canonicalSet(request.Idea.Capabilities)
	request.Idea.Constraints = canonicalSet(request.Idea.Constraints)
	request.Idea.NonGoals = canonicalSet(request.Idea.NonGoals)
	request.Idea.PreferredStack = canonicalSet(request.Idea.PreferredStack)
	seenSources := map[string]bool{}
	for index := range request.SourceArtifacts {
		source := &request.SourceArtifacts[index]
		if !runIDPattern.MatchString(string(source.RunID)) || !synthesisSHA256Pattern.MatchString(source.SHA256) || !allowedSynthesisSource(*source) {
			return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
		}
		key := string(source.RunID) + "\x00" + string(source.RunKind) + "\x00" + string(source.Artifact)
		if seenSources[key] {
			return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
		}
		seenSources[key] = true
	}
	sort.Slice(request.SourceArtifacts, func(i, j int) bool {
		a, b := request.SourceArtifacts[i], request.SourceArtifacts[j]
		return strings.Join([]string{string(a.RunID), string(a.RunKind), string(a.Artifact), a.SHA256}, "\x00") < strings.Join([]string{string(b.RunID), string(b.RunKind), string(b.Artifact), b.SHA256}, "\x00")
	})
	seenProposals := map[string]bool{}
	for index := range request.Proposals {
		proposal := &request.Proposals[index]
		if !validSynthesisText(proposal.ProposalID, 128) || !validSynthesisText(proposal.Title, 512) || !validSynthesisText(proposal.TargetProperty, 512) || !validSynthesisText(proposal.ExpectedEffect, 512) || !validSynthesisText(proposal.Rationale, 1000) || !validSynthesisTextSet(proposal.Assumptions, 20, 512) {
			return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
		}
		proposal.ProposalID = normalizeMechanismText(proposal.ProposalID)
		proposal.Title = normalizeMechanismText(proposal.Title)
		proposal.TargetProperty = normalizeMechanismText(proposal.TargetProperty)
		proposal.ExpectedEffect = normalizeMechanismText(proposal.ExpectedEffect)
		proposal.Rationale = normalizeMechanismText(proposal.Rationale)
		proposal.MechanismRefs = canonicalMechanismSet(proposal.MechanismRefs)
		proposal.ProductObservationRefs = canonicalMechanismSet(proposal.ProductObservationRefs)
		proposal.Assumptions = canonicalMechanismSet(proposal.Assumptions)
		if seenProposals[proposal.ProposalID] || !proposalIDPattern.MatchString(proposal.ProposalID) || proposal.Title == "" || len(proposal.Title) > 512 || proposal.TargetProperty == "" || len(proposal.TargetProperty) > 512 || proposal.ExpectedEffect == "" || len(proposal.ExpectedEffect) > 512 || proposal.Rationale == "" || len(proposal.Rationale) > 1000 || proposal.Derivation != domain.DerivationExternal || len(proposal.MechanismRefs) > 4 || len(proposal.ProductObservationRefs) > 20 || len(proposal.Assumptions) > 20 || len(proposal.MechanismRefs) > request.Policy.MaxMechanismsPerHypothesis {
			return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
		}
		for _, ref := range proposal.MechanismRefs {
			if !mechanismIDPattern.MatchString(ref) {
				return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
			}
		}
		for _, ref := range proposal.ProductObservationRefs {
			if !observationIDPattern.MatchString(ref) {
				return domain.ImprovementSynthesisRequest{}, domain.ErrInvalidRunInput
			}
		}
		seenProposals[proposal.ProposalID] = true
	}
	sort.Slice(request.Proposals, func(i, j int) bool { return request.Proposals[i].ProposalID < request.Proposals[j].ProposalID })
	return request, nil
}

func allowedSynthesisSource(source domain.SourceArtifactRef) bool {
	if source.RunKind == domain.RunKindRecommendation {
		return source.Artifact == domain.ArtifactMechanismCatalog || source.Artifact == domain.ArtifactMechanismGaps
	}
	if source.RunKind == domain.RunKindMarketDiscovery {
		return source.Artifact == domain.ArtifactMarketCandidates || source.Artifact == domain.ArtifactMarketSnapshot
	}
	return false
}

func synthesisRequestHash(request domain.ImprovementSynthesisRequest) (string, error) {
	envelope := struct {
		SchemaVersion        string                            `json:"schema_version"`
		Idea                 domain.IdeaSpec                   `json:"idea"`
		SourceArtifacts      []domain.SourceArtifactRef        `json:"source_artifacts"`
		Proposals            []domain.ImprovementProposal      `json:"proposals"`
		Policy               domain.ImprovementSynthesisPolicy `json:"policy"`
		MechanismRulesSHA256 string                            `json:"mechanism_rules_sha256"`
	}{"1", request.Idea, request.SourceArtifacts, request.Proposals, request.Policy, activeMechanismRulesSHA256}
	encoded, err := json.Marshal(envelope)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func (s *RecommendationRunService) resolveSynthesisSources(ctx context.Context, sources []domain.SourceArtifactRef) (resolvedSynthesisSources, error) {
	resolved := resolvedSynthesisSources{mechanisms: map[string]domain.ArchitecturalMechanism{}, mechanismGaps: map[string]domain.MechanismGap{}, observations: map[string]bool{}}
	for _, source := range sources {
		run, err := s.store.Get(ctx, source.RunID)
		if err != nil {
			return resolvedSynthesisSources{}, err
		}
		if effectiveRunKind(run) != source.RunKind {
			return resolvedSynthesisSources{}, domain.ErrInvalidRunInput
		}
		if !terminalLifecycle(run.Lifecycle) {
			return resolvedSynthesisSources{}, domain.ErrSourceBasisUnavailable
		}
		index, ok := publishedArtifactIndex(run, source.Artifact)
		if !ok {
			return resolvedSynthesisSources{}, domain.ErrSourceBasisUnavailable
		}
		hash := index.Hash
		if hash == "" {
			hash = index.ArtifactHash
		}
		if hash != "sha256:"+source.SHA256 {
			return resolvedSynthesisSources{}, domain.ErrIntegrity
		}
		entries, err := s.readAllArtifactEntries(ctx, source, hash)
		if err != nil {
			return resolvedSynthesisSources{}, err
		}
		if run.Lifecycle != domain.RunSucceeded {
			resolved.degraded = true
		}
		for _, entry := range entries {
			switch source.Artifact {
			case domain.ArtifactMechanismCatalog:
				var mechanism domain.ArchitecturalMechanism
				if json.Unmarshal(entry, &mechanism) != nil || !mechanismIDPattern.MatchString(mechanism.ID) || resolved.mechanisms[mechanism.ID].ID != "" {
					return resolvedSynthesisSources{}, domain.ErrIntegrity
				}
				resolved.mechanisms[mechanism.ID] = mechanism
			case domain.ArtifactMechanismGaps:
				var gap domain.MechanismGap
				if json.Unmarshal(entry, &gap) != nil || gap.Code == "" {
					return resolvedSynthesisSources{}, domain.ErrIntegrity
				}
				key := gap.Code + ":" + gap.RuleID + ":" + gap.WitnessHash
				resolved.mechanismGaps[key] = gap
			case domain.ArtifactMarketCandidates, domain.ArtifactMarketSnapshot:
				var identity struct {
					ID            string `json:"id"`
					CandidateID   string `json:"candidate_id"`
					ObservationID string `json:"observation_id"`
					QueryID       string `json:"query_id"`
					AttemptID     string `json:"attempt_id"`
					Identity      struct {
						ConnectorID string `json:"connector_id"`
						ExternalID  string `json:"external_id"`
					} `json:"identity"`
				}
				if json.Unmarshal(entry, &identity) != nil {
					return resolvedSynthesisSources{}, domain.ErrIntegrity
				}
				for _, id := range []string{identity.ID, identity.CandidateID, identity.ObservationID, identity.QueryID, identity.AttemptID, identity.Identity.ExternalID, identity.Identity.ConnectorID, identity.Identity.ConnectorID + ":" + identity.Identity.ExternalID} {
					if observationIDPattern.MatchString(id) {
						resolved.observations[id] = true
					}
				}
			}
		}
	}
	return resolved, nil
}

func publishedArtifactIndex(run domain.RecommendationRun, name domain.RunArtifactName) (domain.RunArtifactIndex, bool) {
	for _, index := range run.Artifacts {
		if index.Name == name && index.State == domain.ArtifactPublished {
			return index, true
		}
	}
	return domain.RunArtifactIndex{}, false
}

func (s *RecommendationRunService) readAllArtifactEntries(ctx context.Context, source domain.SourceArtifactRef, expectedHash string) ([][]byte, error) {
	cursor := ""
	entries := [][]byte{}
	for pageNumber := 0; pageNumber < 1000; pageNumber++ {
		page, err := s.store.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: source.RunID, Artifact: source.Artifact, Cursor: cursor, Limit: 50})
		if err != nil {
			return nil, err
		}
		if page.Artifact != source.Artifact || page.RunID != source.RunID || page.ArtifactHash != expectedHash {
			return nil, domain.ErrIntegrity
		}
		for _, entry := range page.Entries {
			if !json.Valid(entry) {
				return nil, domain.ErrIntegrity
			}
			entries = append(entries, append([]byte(nil), entry...))
		}
		if page.NextCursor == nil {
			return entries, nil
		}
		if *page.NextCursor == "" || *page.NextCursor == cursor {
			return nil, domain.ErrIntegrity
		}
		cursor = *page.NextCursor
	}
	return nil, domain.ErrIntegrity
}

func validateProposalReferences(proposals []domain.ImprovementProposal, resolved resolvedSynthesisSources) error {
	for _, proposal := range proposals {
		for _, ref := range proposal.MechanismRefs {
			if _, ok := resolved.mechanisms[ref]; !ok {
				return domain.ErrInvalidRunInput
			}
		}
		for _, ref := range proposal.ProductObservationRefs {
			if !resolved.observations[ref] {
				return domain.ErrInvalidRunInput
			}
		}
	}
	return nil
}

func (s *RecommendationRunService) ExecuteImprovementSynthesis(ctx context.Context, id domain.RunID) error {
	run, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if effectiveRunKind(run) != domain.RunKindImprovementSynthesis || run.ImprovementRequest == nil {
		return domain.ErrInvalidRunInput
	}
	if terminalLifecycle(run.Lifecycle) {
		return nil
	}
	if run.Lifecycle == domain.RunQueued {
		run, err = s.store.Transition(ctx, domain.RunTransition{RunID: id, RunKind: domain.RunKindImprovementSynthesis, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: domain.RunStageSynthesisCompare, Progress: domain.RecommendationRunProgress{Planned: 5, Completed: 1}})
		if err != nil {
			return err
		}
	}
	if run.Lifecycle != domain.RunRunning {
		return nil
	}
	run, err = s.ensureRequestArtifact(context.Background(), id)
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return s.confirmCancellation(context.Background(), id)
	}
	resolved, err := s.resolveSynthesisSources(ctx, run.ImprovementRequest.SourceArtifacts)
	if err != nil {
		return s.failRun(ctx, run, "SYNTHESIS_SOURCE_FAILED", "synthesis source resolution failed")
	}
	if err := validateProposalReferences(run.ImprovementRequest.Proposals, resolved); err != nil {
		return s.failRun(ctx, run, "SYNTHESIS_REFERENCE_FAILED", "synthesis reference validation failed")
	}
	run, err = s.advanceSynthesis(ctx, run, domain.RunStageSynthesisDerive, 2)
	if err != nil {
		return err
	}
	analysis := synthesizeImprovement(*run.ImprovementRequest, improvementAnalysisBasis{mechanisms: resolved.mechanisms, degraded: resolved.degraded})
	run, err = s.advanceSynthesis(ctx, run, domain.RunStageSynthesisCounterEvidence, 3)
	if err != nil {
		return err
	}
	run, err = s.advanceSynthesis(ctx, run, domain.RunStageSynthesisValidate, 4)
	if err != nil {
		return err
	}
	run, err = s.advanceSynthesis(ctx, run, domain.RunStageSynthesisPublish, 5)
	if err != nil {
		return err
	}
	documents := []struct {
		name    domain.RunArtifactName
		entries []any
	}{
		{domain.ArtifactArchitectureComparison, interfacesForSlice(analysis.comparisons)},
		{domain.ArtifactImprovementHypotheses, interfacesForSlice(analysis.hypotheses)},
		{domain.ArtifactHypothesisGaps, interfacesForSlice(analysis.gaps)},
		{domain.ArtifactValidationPlans, interfacesForSlice(analysis.plans)},
		{domain.ArtifactSynthesisProvenance, interfacesForSlice(analysis.provenance)},
	}
	for _, document := range documents {
		if ctx.Err() != nil {
			return s.confirmCancellation(context.Background(), id)
		}
		run, err = s.publishDocument(ctx, run, document.name, document.entries)
		if err != nil {
			return err
		}
	}
	_, err = s.store.Transition(ctx, domain.RunTransition{RunID: id, RunKind: domain.RunKindImprovementSynthesis, ExpectedVersion: run.Version, Lifecycle: domain.RunSucceeded, Stage: domain.RunStageTerminal, Progress: domain.RecommendationRunProgress{Planned: 5, Completed: 5}, ImprovementStatus: &analysis.status})
	return err
}

func (s *RecommendationRunService) advanceSynthesis(ctx context.Context, run domain.RecommendationRun, stage domain.RecommendationRunStage, completed int) (domain.RecommendationRun, error) {
	if ctx.Err() != nil {
		return run, ctx.Err()
	}
	if run.Stage == stage && run.Progress.Completed == completed {
		return run, nil
	}
	return s.store.Transition(ctx, domain.RunTransition{RunID: run.RunID, RunKind: domain.RunKindImprovementSynthesis, ExpectedVersion: run.Version, Lifecycle: domain.RunRunning, Stage: stage, Progress: domain.RecommendationRunProgress{Planned: 5, Completed: completed}})
}

func interfacesForSlice[T any](values []T) []any {
	out := make([]any, len(values))
	for index := range values {
		out[index] = values[index]
	}
	return out
}

func validSynthesisText(value string, max int) bool {
	if value == "" || len(value) > max || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f || r >= 0x80 && r <= 0x9f {
			return false
		}
	}
	return true
}

func validSynthesisTextSet(values []string, maxItems, maxBytes int) bool {
	if len(values) > maxItems {
		return false
	}
	for _, value := range values {
		if !validSynthesisText(value, maxBytes) {
			return false
		}
	}
	return true
}
