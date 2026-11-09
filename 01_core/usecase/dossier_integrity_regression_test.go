/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/dossier_integrity_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type failingDossierFiles struct {
	*memoryDossierFiles
	fail      string
	empty     bool
	announced bool
	prepared  bool
}

func (f *failingDossierFiles) ArtifactSHA256(path string) (string, error) {
	if path == f.fail {
		if f.empty {
			return "", nil
		}
		return "", errors.New("hash unavailable")
	}
	return f.memoryDossierFiles.ArtifactSHA256(path)
}
func (f *failingDossierFiles) Announce(RecommendationDossierOutput) { f.announced = true }
func (f *failingDossierFiles) Prepare(resume bool) error {
	f.prepared = true
	return f.memoryDossierFiles.Prepare(resume)
}
func dossierIntegrityFixture() (RecommendationDossierInput, *failingDossierFiles) {
	req, deps := replayRunFixture()
	return RecommendationDossierInput{Request: req, Target: domain.SoftwareAlternativeTarget{Name: "Graph"}, Rules: DefaultEvidenceRules(), Dependencies: deps}, &failingDossierFiles{memoryDossierFiles: &memoryDossierFiles{docs: map[string][]byte{}, snapshot: *deps.ReplaySnapshot, bundles: deps.ReplayBundles}}
}
func TestDossierHashFailureBlocksPublication(t *testing.T) {
	for _, path := range []string{"idea.json", "reuse-policy.json", "expansion-plan.json", "discovery-snapshot.json", "screening-report.json", "evidence/bundles.json", "assessment-report.json", "reuse-decision.json", "alternatives.json", "alternatives.md"} {
		for _, empty := range []bool{false, true} {
			t.Run(path+string(rune('0'+boolInt(empty))), func(t *testing.T) {
				input, files := dossierIntegrityFixture()
				files.fail = path
				files.empty = empty
				out, err := RunRecommendationDossier(context.Background(), input, files)
				if err == nil || out.Published || files.announced {
					t.Fatal("hash failure published", path, err)
				}
				var manifest domain.IdeaRunManifest
				if e := json.Unmarshal(files.docs["run.json"], &manifest); e != nil {
					t.Fatal(e)
				}
				if manifest.Status != domain.IdeaRunFailed {
					t.Fatal("failure not persisted")
				}
				failed := 0
				for _, stage := range manifest.Stages {
					if stage.Status == domain.IdeaStageFailed {
						failed++
						if len(stage.Outputs) > 0 {
							t.Fatal("failed stage advertises outputs")
						}
					}
				}
				if failed == 0 {
					t.Fatal("no failed stage")
				}
			})
		}
	}
}
func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
func TestDossierLocaleIdentityAndResumeCanonical(t *testing.T) {
	for _, variants := range [][][]string{{{"pt-BR", "en-US"}, {"en-US", "pt-BR"}, {"pt-BR", "en-US", "pt-BR"}}, {{"all"}, append([]string(nil), canonicalSearchLocales...), {"all", "en-US"}}} {
		input, files := dossierIntegrityFixture()
		files.snapshot.SearchLocales = append([]string(nil), variants[0]...)
		var first domain.IdeaRunManifest
		for i, locales := range variants {
			input.Request.SearchLocales = locales
			input.Request.Resume = i > 0
			before, _ := json.Marshal(locales)
			out, err := RunRecommendationDossier(context.Background(), input, files)
			if err != nil {
				t.Fatal(err)
			}
			after, _ := json.Marshal(locales)
			if string(before) != string(after) {
				t.Fatal("caller locales mutated")
			}
			if i == 0 {
				first = out.Manifest
			} else {
				if first.RunID != out.Manifest.RunID || first.InputHashes["evidence"] != out.Manifest.InputHashes["evidence"] || first.InputHashes["search_locales"] != out.Manifest.InputHashes["search_locales"] {
					t.Fatal("equivalent locales changed identity")
				}
				if out.Checkpoint == nil || !out.Checkpoint.Reusable {
					t.Fatalf("checkpoint not reused: %+v", out.Checkpoint)
				}
			}
		}
	}
	input, files := dossierIntegrityFixture()
	input.Request.SearchLocales = []string{"invalid"}
	if _, err := RunRecommendationDossier(context.Background(), input, files); err == nil || files.prepared || len(files.docs) > 0 {
		t.Fatal("invalid locale reached effects")
	}
	if _, err := DossierEvidenceInputHash(input.Request, nil, files); err == nil {
		t.Fatal("hash accepted invalid locale")
	}
}
func TestDossierCanonicalControlsKeepHashes(t *testing.T) {
	input, files := dossierIntegrityFixture()
	hash, err := DossierEvidenceInputHash(input.Request, nil, files)
	if err != nil {
		t.Fatal(err)
	}
	out, err := RunRecommendationDossier(context.Background(), input, files)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "f1f129c66936aa0fe0899490117925b465c606d26afef0c963b14717a33cc93e" || out.Manifest.RunID != "630123f56d98a495" || dossierHashJSON(out.Manifest.Stages) != "cf6598de9b77c061f854b64b56d4dac0f523634b94831ef7935724c16ae7034f" {
		t.Fatal("canonical historical control changed")
	}
	if len(out.Manifest.InputHashes["search_locales"]) > 0 || !out.Published || !files.announced {
		t.Fatal("legacy profile changed")
	}
	for _, s := range out.Manifest.Stages {
		if s.Status != domain.IdeaStageComplete {
			t.Fatal("valid stage failed")
		}
		for _, a := range s.Outputs {
			if len(a.Hash) != 64 || strings.TrimSpace(a.Hash) == "" {
				t.Fatal("missing hash")
			}
		}
	}
}
