/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_passages_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"strings"
	"testing"
)

func TestInsightPassagesDiversifyRepositoriesAndKinds(t *testing.T) {
	r := explorationFixture()
	r.MaxPassages = 3
	b := &r.Bundles[0]
	b.Items = append(b.Items, explorationItem(*b, "README.md", "queue cancellation architecture\n", domain.EvidenceDocumentation))
	other := renamedSourceFixture("other").Bundle
	other.Items = []domain.CollectedEvidence{explorationItem(other, "actor.go", "package actor\nfunc Dispatch() {}\n", domain.EvidenceFile)}
	r.Bundles = append(r.Bundles, other)
	got, err := ExploreInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	kinds, repos := map[domain.EvidenceKind]bool{}, map[string]bool{}
	for _, card := range got.Cards {
		kinds[card.Kind] = true
		repos[card.Source.Repository] = true
	}
	if len(got.Cards) != 3 || len(repos) != 2 || len(kinds) != 3 || got.Diagnostics.SelectedSources != 2 {
		t.Fatalf("structural variety lost: %+v", got.Cards)
	}
	for _, card := range got.Cards {
		if card.Source.Repository == "other/queue-library" && !containsText(card.SelectionReasons, "NO_QUESTION_TERM_MATCH") {
			t.Fatal("unmatched alternative presented as relevant")
		}
	}
}

func TestInsightPassagesFamiliarityIsLexicalAndDoesNotExcludeEvidence(t *testing.T) {
	r := explorationFixture()
	r.Question = "queue"
	r.MaxPassages = 1
	b := &r.Bundles[0]
	b.Items = []domain.CollectedEvidence{explorationItem(*b, "a.go", "queue Reserve\n", domain.EvidenceFile), explorationItem(*b, "b.go", "queue alternative\n", domain.EvidenceFile)}
	got, err := ExploreInsights(context.Background(), r)
	if err != nil || got.Cards[0].Reference.Path != "b.go" {
		t.Fatal("familiar match dominates tied alternative")
	}
	r.MaxPassages = 2
	got, err = ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 2 || len(got.Cards[1].FamiliarSubstringMatches) != 1 || got.Cards[1].FamiliarSubstringMatches[0] != "Reserve" {
		t.Fatal("familiar evidence discarded or lexical scope lost")
	}
}

func TestInsightPassagesExactLocationsAndNoDuplicateContext(t *testing.T) {
	r := explorationFixture()
	b := &r.Bundles[0]
	content := "\nqueue cancellation\nreturn false\n\n"
	b.Items = []domain.CollectedEvidence{explorationItem(*b, "queue.go", content, domain.EvidenceFile), explorationItem(*b, "copy.go", content, domain.EvidenceFile)}
	got, err := ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 1 || got.Diagnostics.DuplicatePassages != 1 {
		t.Fatal("copy consumes context budget")
	}
	ref := got.Cards[0].Reference
	if ref.StartLine != 2 || ref.EndLine != 3 || ref.Excerpt != "queue cancellation\nreturn false" || ref.ContentHash != "sha256:"+contentHash(content) {
		t.Fatalf("locator or bytes lost: %+v", ref)
	}
	if got.Diagnostics.SelectionLimited {
		t.Fatal("duplicates incorrectly presented as useful omitted context")
	}
}

func TestInsightPassagesPrefixAndTextBoundsAreExplicit(t *testing.T) {
	r := explorationFixture()
	r.Question = "queue"
	b := &r.Bundles[0]
	b.Items = []domain.CollectedEvidence{
		explorationItem(*b, "prefix.go", strings.Repeat("filler\n", 1024)+"queue\n", domain.EvidenceFile),
		explorationItem(*b, "binary.go", "a\x00b", domain.EvidenceFile),
		explorationItem(*b, "invalid.go", string([]byte{255}), domain.EvidenceFile),
		explorationItem(*b, "LICENSE", "queue", domain.EvidenceLicense),
		explorationItem(*b, "long.go", strings.Repeat("x", 2049), domain.EvidenceFile),
	}
	got, err := ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 1 || got.Diagnostics.PrefixLimitedFiles != 1 || got.Diagnostics.IgnoredItems != 4 {
		t.Fatalf("bounds/ignored items lost: %v %+v", err, got.Diagnostics)
	}
	if len(got.Cards[0].MatchedQuestionTerms) != 0 || got.Cards[0].Reference.EndLine > 1024 {
		t.Fatal("unscanned suffix used")
	}
	r.Bundles[0].Items = []domain.CollectedEvidence{explorationItem(r.Bundles[0], "limit.go", strings.Repeat("é", 1024), domain.EvidenceFile)}
	got, err = ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 1 || len(got.Cards[0].Reference.Excerpt) != 2048 {
		t.Fatal("exact UTF-8 byte boundary rejected or sliced")
	}
}

func TestInsightPassagesPreferAnotherRepositoryOverAnotherRevision(t *testing.T) {
	r := explorationFixture()
	r.MaxPassages = 2
	first := r.Bundles[0]
	first.Items = []domain.CollectedEvidence{explorationItem(first, "a.go", "queue first", domain.EvidenceFile)}
	revision := first
	revision.Repository.Commit = strings.Repeat("c", 40)
	revision.Items = []domain.CollectedEvidence{explorationItem(revision, "a.go", "queue second revision", domain.EvidenceFile)}
	other := renamedSourceFixture("other").Bundle
	other.Items = []domain.CollectedEvidence{explorationItem(other, "alternative.go", "different implementation", domain.EvidenceFile)}
	r.Bundles = []domain.EvidenceBundle{first, revision, other}
	got, err := ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 2 || got.Diagnostics.SelectedSources != 2 {
		t.Fatal("another revision consumed repository diversity slot")
	}
}

func TestInsightPassagesNoMatchUsesTextVarietyAndKeepsUncertainty(t *testing.T) {
	r := explorationFixture()
	r.Question = "vocabulário ausente"
	r.FamiliarApproaches = nil
	b := &r.Bundles[0]
	content := strings.Repeat("import repeated\n", 16) + strings.Repeat("worker capacity reserve rejection cancel context handler\n", 16)
	b.Items = []domain.CollectedEvidence{explorationItem(*b, "queue.go", content, domain.EvidenceFile)}
	got, err := ExploreInsights(context.Background(), r)
	if err != nil || len(got.Cards) != 1 || got.Cards[0].Reference.StartLine == 1 || !strings.Contains(got.Cards[0].Reference.Excerpt, "worker") || !containsText(got.Gaps, "QUESTION_TERMS_NOT_MATCHED") {
		t.Fatalf("no-match context is unhelpful or promoted: %v", err)
	}
	ref := got.Cards[0].Reference
	if ref.EndLine-ref.StartLine >= 16 || len(ref.Excerpt) > 2048 || ref.Excerpt != strings.Join(observationLines(content)[ref.StartLine-1:ref.EndLine], "\n") {
		t.Fatal("window bounds or bytes lost")
	}
}
