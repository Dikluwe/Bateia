/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_collection_test.md
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

func TestReviewAlternativesPreservesFirstPassAndExposesDistantSourcePassage(t *testing.T) {
	r := explorationFixture()
	b := r.Bundles[0]
	lines := []string{}
	for i := 0; i < 64; i++ {
		lines = append(lines, "neutral explanatory material")
	}
	lines[0] = "queue cancellation coordinate first approach"
	lines[48] = "queue cancellation second competing approach"
	b.Items = []domain.CollectedEvidence{explorationItem(b, "README.md", strings.Join(lines, "\n"), domain.EvidenceDocumentation)}
	r.Bundles = []domain.EvidenceBundle{b}
	r.MaxPassages = 3
	before, _ := json.Marshal(r)
	standard, err := ExploreInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if len(standard.Cards) != 1 {
		t.Fatal("historic one-per-file behavior changed")
	}
	r.ReviewAlternatives = true
	reviewed, err := ExploreInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviewed.Cards) != 3 || reviewed.Cards[0].ID != standard.Cards[0].ID || reviewed.ID == standard.ID {
		t.Fatal("review lost first pass or causal identity")
	}
	for i, c := range reviewed.Cards {
		if c.Reference.Excerpt != strings.Join(lines[c.Reference.StartLine-1:c.Reference.EndLine], "\n") {
			t.Fatal("fabricated excerpt")
		}
		for _, prev := range reviewed.Cards[:i] {
			if c.Reference.StartLine <= prev.Reference.EndLine && c.Reference.EndLine >= prev.Reference.StartLine {
				t.Fatal("overlapping passages")
			}
		}
	}
	joined := ""
	for _, c := range reviewed.Cards {
		joined += c.Reference.Excerpt
	}
	if !strings.Contains(joined, "second competing approach") {
		t.Fatal("distant option still hidden")
	}
	if len(reviewed.Tasks) != 7 || reviewed.EffectState != "UNVERIFIED" {
		t.Fatal("missing review or promotion")
	}
	r.ReviewAlternatives = false
	after, _ := json.Marshal(r)
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	r.ReviewAlternatives = true
	r.MaxPassages = 1
	bounded, err := ExploreInsights(context.Background(), r)
	if err != nil || len(bounded.Cards) != 1 || !bounded.Diagnostics.SelectionLimited {
		t.Fatal("global ceiling ignored or loss not reported")
	}
}

func TestCollectInsightEvidenceConnectsNewBytesAndRejectsInvalidBeforeProvider(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: strings.Repeat("a", 40)}
	plan := domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit}, Budget: defaultEvidenceBudget, Targets: []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation}}}
	request := domain.InsightExplorationRequest{SchemaVersion: "1", Question: "How to coordinate queue cancellation?"}
	provider := evidenceProvider{files: []domain.ExtractedFile{{Path: "README.md", Content: "queue cancellation current evidence"}}}
	got, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, provider, request)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Bundles) != 1 || len(got.Exploration.Cards) == 0 || len(got.Exploration.Tasks) != 7 || got.Exploration.EffectState != "UNVERIFIED" {
		t.Fatalf("collection not explored: %+v", got)
	}
	if got.Exploration.Cards[0].Reference.Excerpt != "queue cancellation current evidence" {
		t.Fatal("context not derived from collected bytes")
	}
	request.Question = ""
	if _, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, nil, request); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatalf("invalid request reached collection: %v", err)
	}
	request.Question = "question"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := CollectInsightEvidence(ctx, []domain.EvidenceExtractionPlan{plan}, provider, request); !errors.Is(err, context.Canceled) {
		t.Fatal("cancel ignored")
	}
	plan.Repository.Commit = "HEAD"
	got, err = CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, nil, request)
	if err != nil || len(got.UnexploredRepositories) != 1 || len(got.Bundles[0].Failures) != 1 || got.Exploration.Status != "SEARCH_NEEDED" || !containsText(got.Exploration.Gaps, "UNRESOLVED_SOURCE_BUNDLES") {
		t.Fatalf("unresolved failure promoted or lost: %+v %v", got, err)
	}
}
