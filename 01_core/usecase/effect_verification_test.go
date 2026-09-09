/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/effect_verification_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"reflect"
	"testing"
)

func TestEffectFrozenControls(t *testing.T) {
	cases := []struct {
		name        string
		observation domain.EffectObservation
		want        string
	}{
		{"empty", domain.EffectObservation{InputBefore: []int64{}, InputAfter: []int64{}, Output: []int64{}, Completed: true, Failure: ""}, "supported_in_scope"},
		{"nil-empty-equivalence", domain.EffectObservation{InputBefore: nil, InputAfter: []int64{}, Output: nil, Completed: true, Failure: ""}, "supported_in_scope"},
		{"singleton", domain.EffectObservation{InputBefore: []int64{7}, InputAfter: []int64{7}, Output: []int64{7}, Completed: true, Failure: ""}, "supported_in_scope"},
		{"negatives-duplicates", domain.EffectObservation{InputBefore: []int64{2, -3, 2, 0}, InputAfter: []int64{2, -3, 2, 0}, Output: []int64{-3, 0, 2, 2}, Completed: true, Failure: ""}, "supported_in_scope"},
		{"int64-extremes", domain.EffectObservation{InputBefore: []int64{9223372036854775807, -9223372036854775808, 0}, InputAfter: []int64{9223372036854775807, -9223372036854775808, 0}, Output: []int64{-9223372036854775808, 0, 9223372036854775807}, Completed: true, Failure: ""}, "supported_in_scope"},
		{"already-sorted", domain.EffectObservation{InputBefore: []int64{-2, 0, 0, 9}, InputAfter: []int64{-2, 0, 0, 9}, Output: []int64{-2, 0, 0, 9}, Completed: true, Failure: ""}, "supported_in_scope"},
		{"reversed-output", domain.EffectObservation{InputBefore: []int64{1, 2, 3}, InputAfter: []int64{1, 2, 3}, Output: []int64{3, 2, 1}, Completed: true, Failure: ""}, "counterexample"},
		{"duplicate-count", domain.EffectObservation{InputBefore: []int64{1, 1, 2}, InputAfter: []int64{1, 1, 2}, Output: []int64{1, 2, 2}, Completed: true, Failure: ""}, "counterexample"},
		{"missing-element", domain.EffectObservation{InputBefore: []int64{1, 2}, InputAfter: []int64{1, 2}, Output: []int64{1}, Completed: true, Failure: ""}, "counterexample"},
		{"extra-element", domain.EffectObservation{InputBefore: []int64{1}, InputAfter: []int64{1}, Output: []int64{1, 1}, Completed: true, Failure: ""}, "counterexample"},
		{"input-mutated", domain.EffectObservation{InputBefore: []int64{2, 1}, InputAfter: []int64{1, 2}, Output: []int64{1, 2}, Completed: true, Failure: ""}, "counterexample"},
		{"input-truncated", domain.EffectObservation{InputBefore: []int64{1, 2}, InputAfter: []int64{1}, Output: []int64{1, 2}, Completed: true, Failure: ""}, "counterexample"},
		{"noop", domain.EffectObservation{InputBefore: []int64{2, 1}, InputAfter: []int64{2, 1}, Output: []int64{2, 1}, Completed: true, Failure: ""}, "counterexample"},
		{"panic", domain.EffectObservation{InputBefore: []int64{1}, InputAfter: []int64{1}, Output: nil, Completed: true, Failure: "panic: fixture"}, "counterexample"},
		{"numeric-neighbor-near-max", domain.EffectObservation{InputBefore: []int64{9223372036854775807}, InputAfter: []int64{9223372036854775807}, Output: []int64{9223372036854775806}, Completed: true, Failure: ""}, "counterexample"},
		{"sum-collision", domain.EffectObservation{InputBefore: []int64{1, 3}, InputAfter: []int64{1, 3}, Output: []int64{2, 2}, Completed: true, Failure: ""}, "counterexample"},
		{"timeout", domain.EffectObservation{InputBefore: []int64{2, 1}, InputAfter: nil, Output: nil, Completed: false, Failure: "timeout"}, "unknown"},
		{"outside-scope", domain.EffectObservation{InputBefore: nil, InputAfter: nil, Output: nil, Completed: false, Failure: "forbidden import"}, "unknown"},
		{"incomplete-with-valid-looking-output", domain.EffectObservation{InputBefore: []int64{1}, InputAfter: []int64{1}, Output: []int64{1}, Completed: false, Failure: ""}, "unknown"},
		{"incomplete-with-negative-looking-output", domain.EffectObservation{InputBefore: []int64{2, 1}, InputAfter: []int64{1, 2}, Output: []int64{99}, Completed: false, Failure: "interrupted"}, "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := cloneEffectObservation(tc.observation)
			got := usecase.EvaluateSortEffect(tc.observation)
			if got.Status != tc.want {
				t.Fatalf("status=%q want=%q", got.Status, tc.want)
			}
			if got.Reason == "" {
				t.Error("verdict lacks reason")
			}
			if !reflect.DeepEqual(before, tc.observation) {
				t.Fatalf("evaluation mutated observation: before=%#v after=%#v", before, tc.observation)
			}
		})
	}
}

func cloneEffectSlice(in []int64) []int64 {
	if in == nil {
		return nil
	}
	out := make([]int64, len(in))
	copy(out, in)
	return out
}
func cloneEffectObservation(in domain.EffectObservation) domain.EffectObservation {
	in.InputBefore = cloneEffectSlice(in.InputBefore)
	in.InputAfter = cloneEffectSlice(in.InputAfter)
	in.Output = cloneEffectSlice(in.Output)
	return in
}

func TestEffectDoesNotMutateAliasedObservation(t *testing.T) {
	backing := []int64{3, 1, 2, 9}
	o := domain.EffectObservation{InputBefore: backing[:3], InputAfter: backing[:3], Output: backing[1:], Completed: true}
	original := cloneEffectSlice(backing)
	got := usecase.EvaluateSortEffect(o)
	if got.Status != "counterexample" {
		t.Fatalf("status=%q", got.Status)
	}
	if !reflect.DeepEqual(backing, original) {
		t.Fatal("aliased storage mutated")
	}
}

func TestEffectExhaustiveSmallOutputs(t *testing.T) {
	input := []int64{1, -1, 1}
	var visit func([]int64)
	visit = func(output []int64) {
		counts := map[int64]int{}
		for _, value := range input {
			counts[value]++
		}
		ordered := true
		for i, value := range output {
			counts[value]--
			if i > 0 && output[i-1] > value {
				ordered = false
			}
		}
		valid := ordered
		for _, n := range counts {
			if n != 0 {
				valid = false
			}
		}
		want := "counterexample"
		if valid {
			want = "supported_in_scope"
		}
		o := domain.EffectObservation{InputBefore: cloneEffectSlice(input), InputAfter: cloneEffectSlice(input), Output: cloneEffectSlice(output), Completed: true}
		before := cloneEffectObservation(o)
		got := usecase.EvaluateSortEffect(o)
		if got.Status != want {
			t.Errorf("output=%v got=%s want=%s", output, got.Status, want)
		}
		if !reflect.DeepEqual(before, o) {
			t.Fatal("evaluation mutated observation")
		}
		if len(output) < 4 {
			for value := int64(-1); value <= 1; value++ {
				next := append(cloneEffectSlice(output), value)
				visit(next)
			}
		}
	}
	visit(nil)
}
