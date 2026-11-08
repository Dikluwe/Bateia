/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_plan_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestNeighborhoodAdmissionOwnsFrozenSourcesIndependentlyOfCaller(t *testing.T) {
	input, reader, bindings := neighborhoodFixture(t)
	one := 1
	admission, err := admitInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: input, Focus: "https://github.com/branch/alpha", FollowLinks: &one})
	if err != nil || len(admission.Origins) != 1 || len(reader.resolves) != 0 {
		t.Fatal("admission lost verified sources or performed effects", err)
	}
	before, _ := json.Marshal(admission.Request)
	content := admission.Origins[0].bundle.Items[0].Content
	input.Bundles[0].Items[0].Content = "caller changed the input after admission"
	input.ExplorationMemory.Entries[0].State = "DISCARDED"
	after, _ := json.Marshal(admission.Request)
	if string(before) != string(after) || admission.Origins[0].bundle.Items[0].Content != content {
		t.Fatal("admitted request/sources alias caller-owned state")
	}
	result, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), admission.Request, bindings)
	if err != nil || len(result.Followups) != 1 || result.Followups[0].Status != "COLLECTED" || len(reader.resolves) != 1 || reader.resolves[0] != "branch/alpha" {
		t.Fatal("caller mutation invalidated the independently admitted investigation", err)
	}
}

func TestNeighborhoodGenerationLimitPrecedesEffects(t *testing.T) {
	for _, generation := range []int{999998, 999999} {
		t.Run(fmt.Sprint(generation), func(t *testing.T) {
			input, reader, bindings := neighborhoodFixture(t)
			input.ExplorationMemory.Generation = generation
			memory := sealInsightMemory(*input.ExplorationMemory)
			input.ExplorationMemory = &memory
			input.ID = ""
			input.ID = "site_discovery_" + insightHash(input)
			before, _ := json.Marshal(input)
			zero := 0
			result, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: input, FollowLinks: &zero}, bindings)
			after, _ := json.Marshal(input)
			if string(before) != string(after) || len(reader.reads) != 0 || len(reader.resolves) != 0 {
				t.Fatal("generation admission mutated input or read a source")
			}
			if generation == 999999 {
				if !errors.Is(err, domain.ErrInvalidRunInput) || result.ID != "" {
					t.Fatal("exhausted generation admitted", err)
				}
				return
			}
			if err != nil || result.ExplorationMemory == nil || result.ExplorationMemory.Generation != 999999 {
				t.Fatal("last valid expansion failed", err)
			}
			if _, err := ValidateInsightExplorationMemory(result.ExplorationMemory); err != nil {
				t.Fatal("last expansion memory cannot be reused as history", err)
			}
		})
	}
}
