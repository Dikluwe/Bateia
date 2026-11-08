/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_baseline.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"sort"
	"strings"
	"unicode/utf8"
)

func normalizeInsightBaseline(input *domain.InsightTargetBaseline, classes []string) (*domain.InsightTargetBaseline, *domain.LocatedInsightBaseline, error) {
	if input == nil {
		return nil, nil, nil
	}
	baseline := *input
	parts := strings.Split(baseline.Repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !validSynthesisText(baseline.Repository, 200) || strings.ContainsAny(baseline.Repository, " \\@:#?\t\r\n") || baseline.Repository != normalizeMechanismText(baseline.Repository) || len(baseline.Files) < 1 || len(baseline.Files) > 64 {
		return nil, nil, domain.ErrInvalidRunInput
	}
	baseline.Files = append([]domain.InsightBaselineFile(nil), baseline.Files...)
	sort.Slice(baseline.Files, func(i, j int) bool { return baseline.Files[i].Path < baseline.Files[j].Path })
	manifest := []struct {
		Path        string `json:"path"`
		ContentHash string `json:"content_hash"`
	}{}
	items, size := []domain.CollectedEvidence{}, 0
	for i := range baseline.Files {
		file := &baseline.Files[i]
		size += len(file.Content)
		if !validMechanismPath(file.Path) || (i > 0 && file.Path == baseline.Files[i-1].Path) || !utf8.ValidString(file.Content) || len(file.Content) > 256<<10 || size > 1<<20 {
			return nil, nil, domain.ErrInvalidRunInput
		}
		digest := contentHash(file.Content)
		if strings.TrimPrefix(file.ContentHash, "sha256:") != digest {
			return nil, nil, domain.ErrIntegrity
		}
		file.ContentHash = "sha256:" + digest
		manifest = append(manifest, struct {
			Path        string `json:"path"`
			ContentHash string `json:"content_hash"`
		}{file.Path, file.ContentHash})
		items = append(items, domain.CollectedEvidence{Path: file.Path, Content: file.Content, ContentHash: file.ContentHash})
	}
	digest := insightHash(manifest)
	if strings.TrimPrefix(baseline.SnapshotSHA256, "sha256:") != digest {
		return nil, nil, domain.ErrIntegrity
	}
	baseline.SnapshotSHA256 = "sha256:" + digest
	// The observation adapter uses the snapshot digest only as an internal identity key.
	// It does not classify the target or expose a snapshot as a Git commit.
	adapter := domain.InsightInvestigationRequest{MechanismClasses: classes, Sources: []domain.InsightEvidenceInput{{
		Bundle: domain.EvidenceBundle{Repository: domain.CommitRef{Owner: parts[0], Repo: parts[1], Commit: digest}, Items: items}, Observations: baseline.Observations,
	}}}
	normalized, err := normalizeSourceObservations(adapter)
	if err != nil {
		return nil, nil, err
	}
	baseline.Observations = normalized.Sources[0].Observations
	located := &domain.LocatedInsightBaseline{Repository: baseline.Repository, SnapshotSHA256: baseline.SnapshotSHA256, Scope: "SUPPLIED_FILES_ONLY", Origin: "DECLARED_NOT_REASSESSED", Observations: []domain.TargetInsightObservation{}}
	for _, class := range canonicalMechanismSet(classes) {
		for _, observation := range projectSourceObservations(normalized.Sources[0], class) {
			observation.ID = "target_" + observation.ID
			located.Observations = append(located.Observations, domain.TargetInsightObservation{MechanismClass: class, LocatedInsightObservation: observation})
		}
	}
	return &baseline, located, nil
}
