/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_compare_source_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceCompareTypedSourcePayloadAndHash(t *testing.T) {
	g := traceFixture(t)
	g.GraphHash = ""
	var sourceIndex int
	for i := range g.Nodes {
		if g.Nodes[i].ID == "source-a" {
			sourceIndex = i
			g.Nodes[i].Source = &domain.EvidenceSource{Kind: domain.EvidenceDocumentation, Ref: domain.CommitRef{Owner: "fixture", Repo: "synthetic-doc"}, Snippet: "declared original content"}
			g.Nodes[i].Fingerprint, _ = traceFingerprint(g.Nodes[i].Source)
		}
	}
	before, err := FreezeDecisionTrace(g)
	if err != nil {
		t.Fatal(err)
	}
	after := traceCloneGraph(t, before)
	after.GraphHash = ""
	after.Nodes[sourceIndex].Source.Snippet = "declared changed content"
	// Even with no provided graph hash, a source payload with stale node digest
	// is not admitted as a fresh snapshot.
	if _, err := FreezeDecisionTrace(after); err == nil {
		t.Fatal("stale source digest accepted")
	}
	after.Nodes[sourceIndex].Fingerprint, _ = traceFingerprint(after.Nodes[sourceIndex].Source)
	after, err = FreezeDecisionTrace(after)
	if err != nil {
		t.Fatal(err)
	}
	original := traceCloneGraph(t, before)
	result, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.NodeChanges) != 1 || result.NodeChanges[0].Status != "DECLARATION_CHANGED" || result.NodeChanges[0].Before.Source.Snippet != "declared original content" || result.NodeChanges[0].After.Source.Snippet != "declared changed content" {
		t.Fatalf("wrong typed source delta: %+v", result.NodeChanges)
	}
	if result.Impact.Status != "REVIEW_REQUIRED" || result.Impact.DecisionBefore != result.Impact.DecisionAfter {
		t.Fatal("typed delta did not preserve historical review contract")
	}
	if !reflect.DeepEqual(before, original) {
		t.Fatal("source input mutated")
	}
}
