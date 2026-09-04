/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/temperature_pilot_budget_oracle_test.md
 * @layer L2
 * @updated 2026-09-03
 */
package mcp_test

import (
	"testing"

	mcp "bateia/02_shell/mcp"
)

func pilotLimit(value int) *int { return &value }

func pilotForward(counter *int) func() *mcp.JSONRPCError {
	return func() *mcp.JSONRPCError {
		*counter++
		return nil
	}
}

func requirePilotForwarded(t *testing.T, got *mcp.JSONRPCError) {
	t.Helper()
	if got != nil {
		t.Fatalf("request rejected unexpectedly: %+v", got)
	}
}

func requirePilotExhausted(t *testing.T, got *mcp.JSONRPCError) {
	t.Helper()
	if got == nil {
		t.Fatal("exhausted request was forwarded")
	}
	if got.Code != -32010 || got.Message != "Bateia tool call budget exhausted" {
		t.Fatalf("error=(%d,%q), want (-32010,%q)", got.Code, got.Message, "Bateia tool call budget exhausted")
	}
}

func TestToolCallBudgetLegacyUnlimited(t *testing.T) {
	budget, err := mcp.NewToolCallBudget(nil)
	if err != nil {
		t.Fatalf("nil legacy configuration rejected: %v", err)
	}
	forwarded := 0
	for i := 0; i < 64; i++ {
		requirePilotForwarded(t, budget.Dispatch("tools/call", pilotForward(&forwarded)))
	}
	if forwarded != 64 {
		t.Fatalf("forwarded=%d, want 64", forwarded)
	}
}

func TestToolCallBudgetForwardsExactlyNThenRejectsBeforeCallback(t *testing.T) {
	budget, err := mcp.NewToolCallBudget(pilotLimit(3))
	if err != nil {
		t.Fatal(err)
	}
	forwarded := 0
	for call := 1; call <= 3; call++ {
		requirePilotForwarded(t, budget.Dispatch("tools/call", pilotForward(&forwarded)))
		if forwarded != call {
			t.Fatalf("after call %d forwarded=%d", call, forwarded)
		}
	}
	decoded, handlerEffect := false, false
	requirePilotExhausted(t, budget.Dispatch("tools/call", func() *mcp.JSONRPCError {
		decoded = true
		handlerEffect = true
		return nil
	}))
	if decoded || handlerEffect || forwarded != 3 {
		t.Fatalf("N+1 crossed guard: decoded=%t handler_effect=%t forwarded=%d", decoded, handlerEffect, forwarded)
	}
}

func TestToolCallBudgetNonCallsDoNotConsumeAndToolLabelsShareCounter(t *testing.T) {
	budget, err := mcp.NewToolCallBudget(pilotLimit(2))
	if err != nil {
		t.Fatal(err)
	}
	nonCalls := 0
	forwardedByLabel := map[string]int{}
	steps := []struct {
		method    string
		toolLabel string
		exhausted bool
	}{
		{method: "tools/call", toolLabel: "bateia_recommend"},
		{method: "tools/list"},
		{method: "tools/call", toolLabel: "bateia_search"},
		{method: "notifications/initialized"},
		{method: "initialize"},
		{method: "tools/call", toolLabel: "bateia_benchmark", exhausted: true},
	}
	for _, step := range steps {
		result := budget.Dispatch(step.method, func() *mcp.JSONRPCError {
			if step.method == "tools/call" {
				forwardedByLabel[step.toolLabel]++
			} else {
				nonCalls++
			}
			return nil
		})
		if step.exhausted {
			requirePilotExhausted(t, result)
		} else {
			requirePilotForwarded(t, result)
		}
	}
	if nonCalls != 3 {
		t.Fatalf("interleaved non-call dispatches=%d, want 3", nonCalls)
	}
	if forwardedByLabel["bateia_recommend"] != 1 || forwardedByLabel["bateia_search"] != 1 || forwardedByLabel["bateia_benchmark"] != 0 {
		t.Fatalf("counter reset or split across interleaved methods/tool labels: %v", forwardedByLabel)
	}
}

func TestToolCallBudgetIsProcessInstanceLocal(t *testing.T) {
	first, err := mcp.NewToolCallBudget(pilotLimit(1))
	if err != nil {
		t.Fatal(err)
	}
	second, err := mcp.NewToolCallBudget(pilotLimit(1))
	if err != nil {
		t.Fatal(err)
	}
	firstForwarded, secondForwarded := 0, 0
	requirePilotForwarded(t, first.Dispatch("tools/call", pilotForward(&firstForwarded)))
	requirePilotForwarded(t, second.Dispatch("tools/call", pilotForward(&secondForwarded)))
	requirePilotExhausted(t, first.Dispatch("tools/call", pilotForward(&firstForwarded)))
	requirePilotExhausted(t, second.Dispatch("tools/call", pilotForward(&secondForwarded)))
	if firstForwarded != 1 || secondForwarded != 1 {
		t.Fatalf("independent counters forwarded=%d/%d, want 1/1", firstForwarded, secondForwarded)
	}
}

func TestToolCallBudgetRejectsNonPositiveConstruction(t *testing.T) {
	for _, value := range []int{0, -1} {
		t.Run(string(rune(value+2+'0')), func(t *testing.T) {
			budget, err := mcp.NewToolCallBudget(pilotLimit(value))
			if err == nil || budget != nil {
				t.Fatalf("limit %d did not fail closed: budget=%v err=%v", value, budget, err)
			}
		})
	}
}
