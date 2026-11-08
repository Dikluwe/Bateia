/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_sources.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"encoding/json"

	"bateia/01_core/domain"
)

// synthesisArtifactReader grants source resolution only read access to frozen runs
// and their artifacts. It has no lifecycle, queue or publication operations.
type synthesisArtifactReader interface {
	Get(context.Context, domain.RunID) (domain.RecommendationRun, error)
	ReadArtifact(context.Context, domain.ArtifactReadRequest) (domain.ArtifactPage, error)
}

type resolvedSynthesisSources struct {
	mechanisms    map[string]domain.ArchitecturalMechanism
	mechanismGaps map[string]domain.MechanismGap
	observations  map[string]bool
	degraded      bool
}

func resolveSynthesisSources(ctx context.Context, reader synthesisArtifactReader, sources []domain.SourceArtifactRef) (resolvedSynthesisSources, error) {
	resolved := resolvedSynthesisSources{mechanisms: map[string]domain.ArchitecturalMechanism{}, mechanismGaps: map[string]domain.MechanismGap{}, observations: map[string]bool{}}
	for _, source := range sources {
		run, err := reader.Get(ctx, source.RunID)
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
		entries, err := readAllSynthesisArtifactEntries(ctx, reader, source, hash)
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

func readAllSynthesisArtifactEntries(ctx context.Context, reader synthesisArtifactReader, source domain.SourceArtifactRef, expectedHash string) ([][]byte, error) {
	cursor := ""
	entries := [][]byte{}
	for pageNumber := 0; pageNumber < 1000; pageNumber++ {
		page, err := reader.ReadArtifact(ctx, domain.ArtifactReadRequest{RunID: source.RunID, Artifact: source.Artifact, Cursor: cursor, Limit: 50})
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
