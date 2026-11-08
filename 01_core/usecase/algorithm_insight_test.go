/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_insight_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func algorithmFixture() domain.AlgorithmInsightRequest {
	bundle := func(owner, content string) domain.EvidenceBundle {
		b := domain.EvidenceBundle{SchemaVersion: "1", Repository: domain.CommitRef{Owner: owner, Repo: "scheduler", Commit: strings.Repeat("a", 40)}}
		b.Items = []domain.CollectedEvidence{explorationItem(b, "schedule.go", content, domain.EvidenceFile)}
		return b
	}
	src := bundle("reference", "package schedule\nfunc next() {\n ready.pop()\n}\n")
	target := bundle("target", "package schedule\nfunc next() {\n scanAll()\n}\n")
	evidence := func(b domain.EvidenceBundle) domain.AlgorithmIdeaEvidence {
		return domain.AlgorithmIdeaEvidence{Source: immutableRef(b.Repository), LocatedInsightReference: domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: "schedule.go", ContentHash: b.Items[0].ContentHash, StartLine: 2, EndLine: 4}, Excerpt: b.Items[0].Content[strings.Index(b.Items[0].Content, "func") : len(b.Items[0].Content)-1]}}
	}
	return domain.AlgorithmInsightRequest{SchemaVersion: "1", Question: "How to choose ready work?", Target: target, Sources: []domain.EvidenceBundle{src}, Ideas: []domain.AlgorithmIdea{{Title: "Maintain ready work", Principle: "Maintain candidates when readiness changes instead of rediscovering them on every selection.", Steps: []string{"Update ready state when an input changes.", "Choose the next candidate from the maintained state."}, Evidence: []domain.AlgorithmIdeaEvidence{evidence(src)}, ApplicationConditions: []string{"Changes in readiness are observable."}, FailureConditions: []string{"A missed update leaves stale membership."}, Decision: domain.AlgorithmDecision{Current: "Scan every entry per selection.", Proposed: "Maintain a ready collection.", Benefit: "May reduce selection work.", Cost: "Update bookkeeping and stale-state handling."}, TargetRelation: "DIFFERENT_APPROACH", TargetEvidence: []domain.AlgorithmIdeaEvidence{evidence(target)}, ValidationObservable: "Selection work and stale selections in replay.", RejectIf: "Any ready work is lost or maintenance exceeds saved selection work."}}}
}
func TestAlgorithmIdeasBindDecisionsWithoutPromotingMeaning(t *testing.T) {
	r := algorithmFixture()
	before, _ := json.Marshal(r)
	got, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if got.Scope != "SOURCE_BOUND_DECLARED_REASONING" || got.InterpretationState != "DECLARED_NOT_REASSESSED" || got.TargetRelationState != "DECLARED_NOT_REASSESSED" || got.EffectState != "UNVERIFIED" || len(got.Ideas) != 1 {
		t.Fatal("meaning promoted", got)
	}
	after, _ := json.Marshal(r)
	if string(before) != string(after) {
		t.Fatal("input mutated")
	}
	r.Ideas[0].Evidence[0].ContentHash = strings.TrimPrefix(r.Ideas[0].Evidence[0].ContentHash, "sha256:")
	same, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil || !reflect.DeepEqual(got, same) {
		t.Fatal("hash prefix changed identity")
	}
	r.Ideas[0].Decision.Proposed = "Maintain a different ready representation."
	changed, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil || changed.ID == got.ID {
		t.Fatal("decision not causal")
	}
}
func TestAlgorithmIdeasRejectInventedEvidenceOrIncompleteReasoning(t *testing.T) {
	for name, mutate := range map[string]func(*domain.AlgorithmInsightRequest){
		"fabricated quote": func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence[0].Excerpt = "not in source" },
		"source role":      func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence = r.Ideas[0].TargetEvidence },
		"target role":      func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].TargetEvidence = r.Ideas[0].Evidence },
		"revision":         func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence[0].Source.Commit = "HEAD" },
		"path":             func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence[0].Path = "other.go" },
		"hash":             func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence[0].ContentHash = strings.Repeat("b", 64) },
		"range":            func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Evidence[0].EndLine = 100 },
		"bytes":            func(r *domain.AlgorithmInsightRequest) { r.Sources[0].Items[0].Content += "tamper" },
		"condition":        func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].ApplicationConditions = nil },
		"decision":         func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Decision.Proposed = "" },
		"validation":       func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].RejectIf = "" },
		"relation":         func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].TargetRelation = "PROVEN_NEW" },
		"steps":            func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Steps = r.Ideas[0].Steps[:1] },
		"duplicate":        func(r *domain.AlgorithmInsightRequest) { r.Ideas = append(r.Ideas, r.Ideas[0]) },
	} {
		t.Run(name, func(t *testing.T) {
			r := algorithmFixture()
			mutate(&r)
			if got, err := FormulateAlgorithmInsights(context.Background(), r); err == nil || got.ID != "" {
				t.Fatal("invalid request accepted")
			}
		})
	}
	if _, err := FormulateAlgorithmInsights(nil, algorithmFixture()); !errors.Is(err, domain.ErrInvalidRunInput) {
		t.Fatal("nil context")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FormulateAlgorithmInsights(ctx, algorithmFixture()); !errors.Is(err, context.Canceled) {
		t.Fatal("cancel ignored")
	}
}
func TestAlgorithmIdeasCanonicalOrderAndAgenda(t *testing.T) {
	r := algorithmFixture()
	second := r.Ideas[0]
	second.Title = "Another conditional path"
	second.Steps = append([]string{}, second.Steps...)
	r.Ideas = append(r.Ideas, second)
	first, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	r.Ideas[0], r.Ideas[1] = r.Ideas[1], r.Ideas[0]
	next, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil || !reflect.DeepEqual(first, next) {
		t.Fatal("idea order changes identity")
	}
	ex := explorationFixture()
	ex.AlgorithmIdeas = true
	report, err := ExploreInsights(context.Background(), ex)
	if err != nil || len(report.AlgorithmAgenda) != 6 || len(report.Tasks) != 10 {
		t.Fatal("algorithm agenda missing", err)
	}
	ex.Bundles = nil
	empty, err := ExploreInsights(context.Background(), ex)
	if err != nil || empty.Status != "SEARCH_NEEDED" || len(empty.Cards) != 0 {
		t.Fatal("empty corpus promoted")
	}
	for _, task := range empty.Tasks {
		if !containsText(task.Gaps, "NO_READING_CARDS") {
			t.Fatal("ungrounded algorithm task")
		}
	}
}
