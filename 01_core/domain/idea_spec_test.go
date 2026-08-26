/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_spec_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIdeaSpecJSONKeepsDimensionsSeparate(t *testing.T) {
	spec := IdeaSpec{
		Problem: "impact analysis", Capabilities: []string{"find callers"},
		NonGoals: []string{"runtime tracing"}, PreferredStack: []string{"Rust"},
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	serialized := string(data)
	for _, field := range []string{`"problem"`, `"capabilities"`, `"non_goals"`, `"preferred_stack"`} {
		if !strings.Contains(serialized, field) {
			t.Fatalf("campo %s ausente: %s", field, serialized)
		}
	}
}
