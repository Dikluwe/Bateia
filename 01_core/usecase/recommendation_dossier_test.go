/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_dossier_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type memoryDossierFiles struct {
	docs     map[string][]byte
	snapshot domain.DiscoverySnapshot
	bundles  []domain.EvidenceBundle
	failRead bool
}

func (f *memoryDossierFiles) Prepare(bool) error { return nil }
func (f *memoryDossierFiles) ArtifactSHA256(path string) (string, error) {
	data, ok := f.docs[path]
	if !ok {
		return "", errors.New("missing")
	}
	return dossierHashBytes(data), nil
}
func (f *memoryDossierFiles) RepositoryIdentities(string) []domain.RepositoryIdentity { return nil }
func (f *memoryDossierFiles) ReplayEvidenceIdentity(string) (DossierReplayEvidenceIdentity, error) {
	return DossierReplayEvidenceIdentity{Directory: "frozen-source", SnapshotHash: dossierHashJSON(f.snapshot), BundlesHash: dossierHashJSON(f.bundles)}, nil
}
func (f *memoryDossierFiles) CaptureLocal(context.Context, []string, bool) (domain.CapturedWorkspace, error) {
	return nil, errors.New("unexpected capture")
}
func (f *memoryDossierFiles) SnapshotProvider(domain.ContentProvider, []domain.SnapshotBinding) (domain.ContentProvider, error) {
	return nil, errors.New("unexpected provider")
}
func (f *memoryDossierFiles) HasArtifact(path string) bool { _, ok := f.docs[path]; return ok }
func (f *memoryDossierFiles) ReadJSON(source DossierSource, path string, value any) error {
	var data []byte
	if source == DossierReplaySource {
		if f.failRead {
			return errors.New("frozen source unavailable")
		}
		if path == "discovery-snapshot.json" {
			data, _ = json.Marshal(f.snapshot)
		} else {
			data, _ = json.Marshal(f.bundles)
		}
	} else {
		data = f.docs[path]
	}
	if len(data) == 0 {
		return errors.New("missing")
	}
	return json.Unmarshal(data, value)
}
func (f *memoryDossierFiles) WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err == nil {
		f.docs[path] = append(data, '\n')
	}
	return err
}
func (f *memoryDossierFiles) WriteDecision(domain.ReuseDecisionReport, []domain.EvidenceBundle) error {
	f.docs["decision.md"] = []byte("decision")
	return nil
}
func (f *memoryDossierFiles) WriteAlternatives(domain.SoftwareAlternativeReport) error {
	f.docs["alternatives.md"] = []byte("alternatives")
	return nil
}
func (f *memoryDossierFiles) WriteTrace(value any, path, _ string) error {
	return f.WriteJSON(path, value)
}
func (f *memoryDossierFiles) Announce(RecommendationDossierOutput) {}

func TestRecommendationDossierReplayOwnsIdentityAndPublication(t *testing.T) {
	request, deps := replayRunFixture()
	files := &memoryDossierFiles{docs: map[string][]byte{}, snapshot: *deps.ReplaySnapshot, bundles: deps.ReplayBundles}
	input := RecommendationDossierInput{Request: request, Target: domain.SoftwareAlternativeTarget{Name: "Graph"}, Rules: DefaultEvidenceRules(), Dependencies: deps}
	first, err := RunRecommendationDossier(context.Background(), input, files)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Published || first.Manifest.RunID == "" || len(first.Manifest.Stages) != 7 {
		t.Fatalf("incomplete publication: %+v", first.Manifest)
	}
	input.Request.OutputDirectory = "another-output"
	input.Request.Resume = true
	second, err := RunRecommendationDossier(context.Background(), input, files)
	if err != nil {
		t.Fatal(err)
	}
	if first.Manifest.RunID != second.Manifest.RunID {
		t.Fatal("operational flags changed causal identity")
	}
	files.failRead = true
	input.Request.Resume = false
	failed, err := RunRecommendationDossier(context.Background(), input, files)
	if err == nil || failed.Published {
		t.Fatal("read failure promoted")
	}
	var manifest domain.IdeaRunManifest
	if err := json.Unmarshal(files.docs["run.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Status != domain.IdeaRunFailed {
		t.Fatal("failure not recorded")
	}
}

func TestDossierRejectsInvalidOptInBeforePublication(t *testing.T) {
	request, deps := replayRunFixture()
	files := &memoryDossierFiles{docs: map[string][]byte{}, snapshot: *deps.ReplaySnapshot, bundles: deps.ReplayBundles}
	model := DefaultScreeningModel()
	model.CalibrationStatus = "calibrated"
	_, err := RunRecommendationDossier(context.Background(), RecommendationDossierInput{Request: request, Rules: DefaultEvidenceRules(), ScreeningModel: &model, Dependencies: deps}, files)
	if err == nil || len(files.docs) != 0 {
		t.Fatal("invalid opt-in reached publication")
	}
}
