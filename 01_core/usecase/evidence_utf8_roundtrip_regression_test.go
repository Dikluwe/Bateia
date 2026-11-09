/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_utf8_roundtrip_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func utf8RoundtripAssessment(t *testing.T, b domain.EvidenceBundle) {
	t.Helper()
	_, err := AssessEvidence(domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, []domain.EvidenceBundle{b}, DefaultEvidenceRules())
	if err != nil {
		t.Fatal("bundle failed integrity re-evaluation", err)
	}
}
func utf8RoundtripBundle(t *testing.T, b domain.EvidenceBundle) {
	t.Helper()
	utf8RoundtripAssessment(t, b)
	for round := 0; round < 3; round++ {
		before, err := json.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		var reopened domain.EvidenceBundle
		if err := json.Unmarshal(before, &reopened); err != nil {
			t.Fatal(err)
		}
		if len(b.Items) != len(reopened.Items) {
			t.Fatal("roundtrip lost items")
		}
		for i, item := range reopened.Items {
			if !utf8.ValidString(item.Content) || item.Content != b.Items[i].Content || item.ContentHash != b.Items[i].ContentHash {
				t.Fatalf("roundtrip changed source bytes/hash: round=%d path=%s", round, item.Path)
			}
		}
		after, err := json.Marshal(reopened)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatal("roundtrip changed serialized bytes", err)
		}
		utf8RoundtripAssessment(t, reopened)
		b = reopened
	}
}
func TestCollectionUTF8PrefixSurvivesEveryByteCutAndReplay(t *testing.T) {
	text := "package fixture\n// café € 😀 e\u0301\nfunc Graph(){}\n"
	if !utf8.ValidString(text) {
		t.Fatal("invalid source fixture")
	}
	for ceiling := 0; ceiling <= len(text); ceiling++ {
		t.Run(fmt.Sprint(ceiling), func(t *testing.T) {
			// Independent oracle: accumulate complete runes within the byte ceiling.
			want := ""
			for _, r := range text {
				piece := string(r)
				if len(want)+len(piece) > ceiling {
					break
				}
				want += piece
			}
			plan := evidenceIntegrityPlan()
			plan.Targets = []domain.ExtractionTarget{{Path: "text.go", Kind: domain.EvidenceFile}}
			metadata, _ := json.Marshal(plan.Candidate)
			plan.Budget.TotalBytes = int64(len(metadata) + ceiling)
			if ceiling == len(text) {
				plan.Budget.TotalBytes++
			}
			before, _ := json.Marshal(plan)
			p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "text.go", Type: "file", Size: int64(len(text))}}, content: map[string]string{"text.go": text}, sizes: map[string]int64{"text.go": int64(len(text))}}
			b, err := CollectEvidence(context.Background(), plan, p)
			if err != nil {
				t.Fatal(err)
			}
			result := completenessTarget(t, b, "text.go")
			expectedStatus := domain.TargetSkippedBudget
			if ceiling == len(text) {
				expectedStatus = domain.TargetCollected
			}
			if result.Status != expectedStatus || result.Bytes != int64(len(want)) || b.BytesRead != int64(len(metadata)+len(want)) || b.BytesRead > plan.Budget.TotalBytes {
				t.Fatalf("wrong bounded prefix: ceiling=%d result=%+v bytes=%d", ceiling, result, b.BytesRead)
			}
			if ceiling < len(text) && b.Complete {
				t.Fatal("fragment claimed completeness")
			}
			found := false
			for _, item := range b.Items {
				if item.Path == "text.go" {
					found = true
					digest := fmt.Sprintf("%x", sha256.Sum256([]byte(want)))
					if item.Content != want || !utf8.ValidString(item.Content) || item.ContentHash != digest {
						t.Fatalf("source prefix/hash changed: ceiling=%d got=%q want=%q", ceiling, item.Content, want)
					}
				}
			}
			if found != (want != "") {
				t.Fatal("empty fragment created an item or valid prefix was lost")
			}
			utf8RoundtripBundle(t, b)
			after, _ := json.Marshal(plan)
			if !bytes.Equal(before, after) {
				t.Fatal("caller input mutated")
			}
		})
	}
}
func TestCollectionDoesNotSpendBytesOnIncompleteLeadingCharacter(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = []domain.ExtractionTarget{{Path: "first.go", Kind: domain.EvidenceFile}, {Path: "next.go", Kind: domain.EvidenceFile}}
	metadata, _ := json.Marshal(plan.Candidate)
	plan.Budget.TotalBytes = int64(len(metadata) + 1)
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "first.go", Type: "file"}, {Path: "next.go", Type: "file"}}, content: map[string]string{"first.go": "€value", "next.go": "x"}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	first, next := completenessTarget(t, b, "first.go"), completenessTarget(t, b, "next.go")
	if first.Status != domain.TargetSkippedBudget || first.Bytes != 0 || next.Status != domain.TargetCollected || next.Bytes != 1 || b.BytesRead != plan.Budget.TotalBytes {
		t.Fatalf("discarded bytes consumed budget: first=%+v next=%+v", first, next)
	}
	for _, item := range b.Items {
		if item.Path == "first.go" {
			t.Fatal("created empty/invalid item")
		}
	}
	utf8RoundtripBundle(t, b)
}
func TestCollectionRejectsOriginallyInvalidUTF8WithoutReplacement(t *testing.T) {
	plan := evidenceIntegrityPlan()
	plan.Targets = []domain.ExtractionTarget{{Path: "text.go", Kind: domain.EvidenceFile}}
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "text.go", Type: "file"}}, content: map[string]string{"text.go": "before\xffafter"}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if b.Complete || len(b.Failures) == 0 || completenessTarget(t, b, "text.go").Status != domain.TargetFailed {
		t.Fatal("invalid source was not rejected")
	}
	for _, item := range b.Items {
		if item.Path == "text.go" {
			t.Fatal("invalid source admitted or silently sanitized")
		}
	}
	utf8RoundtripBundle(t, b)
}
func TestCoordinatedUTF8FragmentCanBePublishedAndReopened(t *testing.T) {
	text := "queue café mechanism"
	plan := evidenceIntegrityPlan()
	plan.Targets = []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation}}
	metadata, _ := json.Marshal(plan.Candidate)
	plan.Budget.TotalBytes = int64(len(metadata) + strings.Index(text, "é") + 1)
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, content: map[string]string{"README.md": text}}
	request := domain.InsightExplorationRequest{SchemaVersion: "1", Question: "queue"}
	got, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, p, request)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var reopened domain.InsightCollectionReport
	if err := json.Unmarshal(data, &reopened); err != nil {
		t.Fatal(err)
	}
	if len(reopened.Bundles) != 1 {
		t.Fatal("missing coordinated bundle")
	}
	utf8RoundtripBundle(t, reopened.Bundles[0])
	explored, err := ExploreInsights(context.Background(), domain.InsightExplorationRequest{SchemaVersion: "1", Question: "queue", ReviewAlternatives: true, Bundles: reopened.Bundles})
	if err != nil || explored.ID != got.Exploration.ID {
		t.Fatal("published source cannot regenerate context", err)
	}
}
func TestRecoveryUTF8FragmentRetainsReplayAndPendingHistory(t *testing.T) {
	base := evidenceIntegrityPlan()
	plan := budgetPlan(base.Repository)
	plan.Targets[0].Path = "text.go"
	plan.Targets[0].Kind = domain.EvidenceFile
	metadata, _ := json.Marshal(base.Candidate)
	plan.Policy.AdditionalBytes = int64(len(metadata) + 3)
	p := &completenessBatchPort{tree: []domain.DirectoryNode{{Path: "text.go", Type: "file"}}, content: map[string]string{"text.go": "a😀value"}}
	original := []domain.EvidenceBundle{{SchemaVersion: "1", Repository: base.Repository, Plan: base}}
	before, _ := json.Marshal(original)
	got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, original, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Bundles) != 1 {
		t.Fatal("missing recovery bundle")
	}
	result := completenessTarget(t, got.Bundles[0], "text.go")
	if result.Status != domain.TargetSkippedBudget || result.Bytes != 1 || got.Metrics.TargetsExecuted != 1 || len(ExecutedRecoveryTargets(plan, got.Bundles)) != 0 {
		t.Fatalf("wrong UTF-8 fragment state: %+v metrics=%+v", result, got.Metrics)
	}
	utf8RoundtripBundle(t, got.Bundles[0])
	after, _ := json.Marshal(original)
	if !bytes.Equal(before, after) {
		t.Fatal("original input changed")
	}
}
