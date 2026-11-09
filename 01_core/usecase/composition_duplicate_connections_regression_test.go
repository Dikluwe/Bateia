/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/composition_duplicate_connections_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestCompositionRejectsAmbiguousDuplicateConnectionsInEitherOrder(t *testing.T) {
	for name, mutate := range map[string]func(*domain.CompositionConnection){
		"cost": func(e *domain.CompositionConnection) {
			e.Cost = &domain.IntegrationCost{Unit: "hours", Low: 20, High: 22}
		},
		"unit":         func(e *domain.CompositionConnection) { e.Cost = &domain.IntegrationCost{Unit: "days", Low: 2, High: 4} },
		"unknown cost": func(e *domain.CompositionConnection) { e.Cost = nil },
		"status":       func(e *domain.CompositionConnection) { e.Status = domain.EvidencePartial },
		"source":       func(e *domain.CompositionConnection) { e.Sources = []domain.EvidenceSource{compositionTestSource("b")} },
	} {
		for _, reversed := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/forward", true: "/reversed"}[reversed], func(t *testing.T) {
				req, a, b := compositionFixture()
				req.Composition = verifiedComposition(a, b)
				left := req.Composition.Connections[0]
				right := left
				mutate(&right)
				req.Composition.Connections = []domain.CompositionConnection{left, right}
				if reversed {
					req.Composition.Connections[0], req.Composition.Connections[1] = right, left
				}
				before, _ := json.Marshal(req)
				got, err := DecideReuse(req)
				if err == nil || !reflect.DeepEqual(got, domain.ReuseDecisionReport{}) {
					t.Errorf("ambiguous duplicate admitted: %v %+v", err, got)
				}
				after, _ := json.Marshal(req)
				if string(before) != string(after) {
					t.Fatal("input mutated")
				}
			})
		}
	}
}

func TestCompositionIdenticalDuplicatePreservesReport(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		req, a, b := compositionFixture()
		req.Composition = verifiedComposition(a, b)
		if unknown {
			req.Composition.Connections[0].Cost = nil
		}
		control, err := DecideReuse(req)
		if err != nil || control.Decision != domain.DecisionCompose {
			t.Fatal(err, control)
		}
		raw, _ := json.Marshal(control)
		if !unknown && fmt.Sprintf("%x", sha256.Sum256(raw)) != "4d910f5b81417e53a84817568c03eb9d7784ab284d576a0aa4006d01e0661984" {
			t.Fatal("valid composition report changed")
		}
		duplicate := req.Composition.Connections[0]
		duplicate.Sources = append([]domain.EvidenceSource(nil), duplicate.Sources...)
		if duplicate.Cost != nil {
			cost := *duplicate.Cost
			duplicate.Cost = &cost
		}
		req.Composition.Connections = append(req.Composition.Connections, duplicate)
		got, err := DecideReuse(req)
		if err != nil || !reflect.DeepEqual(got, control) {
			t.Fatal("identical duplicate changed cost/evidence", err, got)
		}
		if unknown && got.Composition.Candidates[0].IntegrationCost != nil {
			t.Fatal("unknown cost became zero")
		}
	}
}
