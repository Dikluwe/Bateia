/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/recommendation_resume_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"reflect"
	"testing"
)

type observedResumeFiles struct {
	*memoryDossierFiles
	manifestReads int
	announcements []RecommendationDossierOutput
}

func (f *observedResumeFiles) ReadJSON(source DossierSource, path string, value any) error {
	if source == DossierOutputSource && path == "run.json" {
		f.manifestReads++
	}
	return f.memoryDossierFiles.ReadJSON(source, path, value)
}
func (f *observedResumeFiles) Announce(out RecommendationDossierOutput) {
	f.announcements = append(f.announcements, out)
}

func TestRecommendationResumeDiagnosticEndToEnd(t *testing.T) {
	request, deps := replayRunFixture()
	// Use a recorded coverage entry for byte comparison: the legacy encoder
	// alternates null/[] for an empty RequirementCoverage on a round trip.
	deps.ReplayBundles[0].RequirementCoverage = []domain.RequirementCollectionCoverage{{RequirementKind: "capability", Requirement: "dependency graph", CandidatePaths: []string{"graph.go"}, ExaminedPaths: []string{"graph.go"}}}
	files := &observedResumeFiles{memoryDossierFiles: &memoryDossierFiles{docs: map[string][]byte{}, snapshot: *deps.ReplaySnapshot, bundles: deps.ReplayBundles}}
	input := RecommendationDossierInput{Request: request, Target: domain.SoftwareAlternativeTarget{Name: "Graph"}, Rules: DefaultEvidenceRules(), Dependencies: deps}
	first, err := RunRecommendationDossier(context.Background(), input, files)
	if err != nil || first.Checkpoint != nil || files.manifestReads != 0 {
		t.Fatalf("fresh run: checkpoint=%+v err=%v", first.Checkpoint, err)
	}
	docs := map[string]string{}
	for name, data := range files.docs {
		docs[name] = string(data)
	}
	input.Request.Resume = true
	second, err := RunRecommendationDossier(context.Background(), input, files)
	if err != nil || second.Checkpoint == nil || !second.Checkpoint.Reusable || second.Checkpoint.Reason != CheckpointReusable || files.manifestReads != 1 {
		t.Fatalf("resume: checkpoint=%+v err=%v reads=%d", second.Checkpoint, err, files.manifestReads)
	}
	if first.Manifest.RunID != second.Manifest.RunID || !reflect.DeepEqual(first.Manifest.InputHashes, second.Manifest.InputHashes) || !reflect.DeepEqual(first.Manifest.RuleVersions, second.Manifest.RuleVersions) {
		t.Fatal("diagnostic altered causal identity")
	}
	if len(docs) != len(files.docs) {
		t.Fatal("diagnostic created new artifact")
	}
	for name, data := range files.docs {
		// run.json includes operational durations; all other artifact bytes stay frozen.
		if name != "run.json" && string(data) != docs[name] {
			t.Fatalf("diagnostic altered %s", name)
		}
	}
	if len(files.announcements) != 2 || files.announcements[1].Checkpoint != second.Checkpoint {
		t.Fatal("announcement did not receive the same assessment")
	}
}

func TestRecommendationResumeDiagnosticRejectionAndLaterFailure(t *testing.T) {
	for _, scenario := range []string{"input", "hash", "missing-manifest", "invalid-manifest", "later-failure"} {
		t.Run(scenario, func(t *testing.T) {
			request, deps := replayRunFixture()
			files := &observedResumeFiles{memoryDossierFiles: &memoryDossierFiles{docs: map[string][]byte{}, snapshot: *deps.ReplaySnapshot, bundles: deps.ReplayBundles}}
			input := RecommendationDossierInput{Request: request, Target: domain.SoftwareAlternativeTarget{Name: "Graph"}, Rules: DefaultEvidenceRules(), Dependencies: deps}
			if _, err := RunRecommendationDossier(context.Background(), input, files); err != nil {
				t.Fatal(err)
			}
			input.Request.Resume = true
			wantReason := CheckpointManifestUnavailable
			switch scenario {
			case "input":
				input.Request.Language = "rust"
				wantReason = CheckpointEvidenceChanged
			case "hash":
				files.docs["discovery-snapshot.json"] = []byte("tampered")
				wantReason = CheckpointArtifactMismatch
			case "missing-manifest":
				delete(files.docs, "run.json")
			default:
				files.docs["run.json"] = []byte("malformed /private/secret")
			}
			files.failRead = scenario == "later-failure"
			got, err := RunRecommendationDossier(context.Background(), input, files)
			if got.Checkpoint == nil || got.Checkpoint.Reusable || got.Checkpoint.Reason != wantReason || files.manifestReads != 1 {
				t.Fatalf("checkpoint=%+v err=%v reads=%d", got.Checkpoint, err, files.manifestReads)
			}
			if scenario == "later-failure" {
				if err == nil || got.Published || len(files.announcements) != 1 {
					t.Fatal("later failure promoted or diagnostic lost")
				}
			} else if err != nil || !got.Published || len(files.announcements) != 2 {
				t.Fatalf("rejected checkpoint prevented legitimate replay: %v", err)
			}
		})
	}
}
