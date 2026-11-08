/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery_plan_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func insightDiscoveryFixture() domain.InsightDiscoveryRequest {
	return domain.InsightDiscoveryRequest{SchemaVersion: "1", Question: "backpressure queue cancellation", FamiliarApproaches: []string{"queue", "thread pool"}}
}

func TestInsightDiscoveryPlanIndependentRoutesAndFamiliarAnchors(t *testing.T) {
	r := insightDiscoveryFixture()
	before, _ := json.Marshal(r)
	plan, err := PlanInsightDiscovery(r)
	if err != nil || len(plan.Queries) != 5 {
		t.Fatalf("plan missing: %v", err)
	}
	want := []string{"backpressure queue cancellation", "backpressure", "backpressure alternative", "cancellation tradeoffs", "backpressure failure"}
	ids := map[string]bool{}
	for i, q := range plan.Queries {
		if q.Query != want[i] || ids[q.ID] || q.Derivation != "LEXICAL_EXPLORATION_HEURISTIC" {
			t.Fatalf("invalid route: %+v", q)
		}
		ids[q.ID] = true
	}
	after, _ := json.Marshal(r)
	if string(before) != string(after) {
		t.Fatal("request mutated")
	}
	r.FamiliarApproaches[0], r.FamiliarApproaches[1] = r.FamiliarApproaches[1], r.FamiliarApproaches[0]
	again, err := PlanInsightDiscovery(r)
	if err != nil || !reflect.DeepEqual(plan, again) {
		t.Fatal("label ordering changes identity")
	}
	r.Question = "queue"
	r.FamiliarApproaches = []string{"queue"}
	fallback, err := PlanInsightDiscovery(r)
	if err != nil || fallback.Queries[0].Query != fallback.Queries[1].Query || fallback.Queries[0].ID == fallback.Queries[1].ID {
		t.Fatal("fallback erased independent route")
	}
	r.Question = "backpressure channel"
	changed, _ := PlanInsightDiscovery(r)
	if changed.ID == plan.ID {
		t.Fatal("question is not causal")
	}
}

func TestInsightDiscoveryPlanRejectsInvalidBudgetsBeforeEffects(t *testing.T) {
	for name, mutate := range map[string]func(*domain.InsightDiscoveryRequest){
		"schema":          func(r *domain.InsightDiscoveryRequest) { r.SchemaVersion = "2" },
		"blank":           func(r *domain.InsightDiscoveryRequest) { r.Question = " " },
		"no-terms":        func(r *domain.InsightDiscoveryRequest) { r.Question = "o que como" },
		"long":            func(r *domain.InsightDiscoveryRequest) { r.Question = strings.Repeat("a", 1001) },
		"per-query":       func(r *domain.InsightDiscoveryRequest) { r.MaxPerQuery = 6 },
		"refs":            func(r *domain.InsightDiscoveryRequest) { r.MaxReferences = 9 },
		"cards":           func(r *domain.InsightDiscoveryRequest) { r.MaxPassages = 17 },
		"negative":        func(r *domain.InsightDiscoveryRequest) { r.MaxReferences = -1 },
		"duration":        func(r *domain.InsightDiscoveryRequest) { r.TotalTimeout = time.Millisecond },
		"total":           func(r *domain.InsightDiscoveryRequest) { r.TotalTimeout = 6 * time.Minute },
		"attempt":         func(r *domain.InsightDiscoveryRequest) { r.AttemptTimeout = 31 * time.Second },
		"attempt-greater": func(r *domain.InsightDiscoveryRequest) { r.TotalTimeout = time.Second },
	} {
		t.Run(name, func(t *testing.T) {
			r := insightDiscoveryFixture()
			mutate(&r)
			if _, err := PlanInsightDiscovery(r); err == nil {
				t.Fatal("invalid request accepted")
			}
		})
	}
	r := insightDiscoveryFixture()
	r.MaxPerQuery = 5
	r.MaxReferences = 8
	r.MaxPassages = 16
	r.TotalTimeout = 5 * time.Minute
	r.AttemptTimeout = 30 * time.Second
	if _, err := PlanInsightDiscovery(r); err != nil {
		t.Fatal("exact budgets rejected", err)
	}
}
