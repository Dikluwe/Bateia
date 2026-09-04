/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/temperature_pilot_config_oracle_test.md
 * @layer L4
 * @updated 2026-09-03
 */
package main

import "testing"

func TestParseMCPMaxToolCallsAbsentPreservesUnlimitedLegacy(t *testing.T) {
	limit, err := parseMCPMaxToolCalls("ignored-placeholder", false)
	if err != nil || limit != nil {
		t.Fatalf("absent config=(%v,%v), want (nil,nil)", limit, err)
	}
}

func TestParseMCPMaxToolCallsCanonicalPositive(t *testing.T) {
	limit, err := parseMCPMaxToolCalls("3", true)
	if err != nil || limit == nil || *limit != 3 {
		t.Fatalf("positive config=(%v,%v), want pointer to 3", limit, err)
	}
}

func TestParseMCPMaxToolCallsInvalidFailsClosed(t *testing.T) {
	for _, raw := range []string{"", "0", "-1", "+3", " 3", "3 ", "3.0", "three", "999999999999999999999999999999999999999"} {
		t.Run(raw, func(t *testing.T) {
			limit, err := parseMCPMaxToolCalls(raw, true)
			if err == nil || limit != nil {
				t.Fatalf("invalid %q did not fail closed: limit=%v err=%v", raw, limit, err)
			}
		})
	}
}
