/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/effect_verification.md
 * @layer L1
 * @updated 2026-09-08
 */

package usecase

import "bateia/01_core/domain"

func EvaluateSortEffect(o domain.EffectObservation) domain.EffectVerdict {
	if !o.Completed {
		return domain.EffectVerdict{Status: "unknown", Reason: "execution incomplete"}
	}
	fail := func(reason string) domain.EffectVerdict {
		return domain.EffectVerdict{Status: "counterexample", Reason: reason}
	}
	if o.Failure != "" {
		return fail(o.Failure)
	}
	if len(o.InputBefore) != len(o.InputAfter) {
		return fail("input length changed")
	}
	for i, v := range o.InputBefore {
		if o.InputAfter[i] != v {
			return fail("input values changed")
		}
	}
	if len(o.Output) != len(o.InputBefore) {
		return fail("output multiplicities changed")
	}
	counts := make(map[int64]int, len(o.InputBefore))
	for _, v := range o.InputBefore {
		counts[v]++
	}
	for i, v := range o.Output {
		if i > 0 && o.Output[i-1] > v {
			return fail("output is not nondecreasing")
		}
		if counts[v] == 0 {
			return fail("output multiplicities changed")
		}
		counts[v]--
	}
	return domain.EffectVerdict{Status: "supported_in_scope", Reason: "all observed sort-preserve-v1 conditions hold"}
}
