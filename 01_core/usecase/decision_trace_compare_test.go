/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_compare_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func compareFixture(t *testing.T) domain.DecisionTraceGraph {
	t.Helper()
	digest := strings.Repeat("a", 64)
	decisionDigest, err := traceFingerprint(domain.DecisionUnknown)
	if err != nil {
		t.Fatal(err)
	}
	return compareRefreeze(t, domain.DecisionTraceGraph{
		SchemaVersion: domain.DecisionTraceSchemaVersion, Scope: domain.DecisionTraceScope, Decision: domain.DecisionUnknown,
		Nodes: []domain.DecisionTraceNode{
			{ID: "decision", Kind: domain.TraceDecision, Fingerprint: decisionDigest},
			{ID: "source", Kind: domain.TraceSource, Fingerprint: digest},
			{ID: "claim", Kind: domain.TraceClaim, Fingerprint: digest},
			{ID: "unrelated", Kind: domain.TraceClaim, Fingerprint: digest},
		},
		Edges: []domain.DecisionTraceEdge{
			{From: "source", To: "claim", Relation: domain.TraceRefutes},
			{From: "claim", To: "decision", Relation: domain.TraceDependsOn},
		},
	})
}

func compareRefreeze(t *testing.T, graph domain.DecisionTraceGraph) domain.DecisionTraceGraph {
	t.Helper()
	graph.GraphHash = ""
	frozen, err := FreezeDecisionTrace(graph)
	if err != nil {
		t.Fatal(err)
	}
	return frozen
}

func compareClone(t *testing.T, graph domain.DecisionTraceGraph) domain.DecisionTraceGraph {
	t.Helper()
	data, err := json.Marshal(graph)
	if err != nil {
		t.Fatal(err)
	}
	var clone domain.DecisionTraceGraph
	if err = json.Unmarshal(data, &clone); err != nil {
		t.Fatal(err)
	}
	return clone
}

func compareRun(t *testing.T, before, after domain.DecisionTraceGraph) domain.DecisionTraceComparison {
	t.Helper()
	result, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestDecisionTraceCompareIdenticalCanonicalAndNonMutating(t *testing.T) {
	before := compareFixture(t)
	after := compareClone(t, before)
	for i, j := 0, len(after.Nodes)-1; i < j; i, j = i+1, j-1 {
		after.Nodes[i], after.Nodes[j] = after.Nodes[j], after.Nodes[i]
	}
	for i, j := 0, len(after.Edges)-1; i < j; i, j = i+1, j-1 {
		after.Edges[i], after.Edges[j] = after.Edges[j], after.Edges[i]
	}
	for i, j := 0, len(after.Limitations)-1; i < j; i, j = i+1, j-1 {
		after.Limitations[i], after.Limitations[j] = after.Limitations[j], after.Limitations[i]
	}
	initialBefore, initialAfter := compareClone(t, before), compareClone(t, after)
	result := compareRun(t, before, after)
	if result.SchemaVersion != domain.DecisionTraceComparisonSchemaVersion || result.Scope != domain.DecisionTraceScope || result.Impact.Status != "NO_RECORDED_IMPACT" || result.BeforeGraphHash != result.AfterGraphHash || result.MetadataChanged {
		t.Fatalf("canonical identical comparison: %#v", result)
	}
	if len(result.NodeChanges) != 0 || len(result.AddedEdges) != 0 || len(result.RemovedEdges) != 0 || len(result.Impact.Tasks) != 0 {
		t.Fatal("input ordering fabricated differences")
	}
	if result.NodeChanges == nil || result.AddedEdges == nil || result.RemovedEdges == nil || result.Limitations == nil {
		t.Fatal("nil output vectors")
	}
	if !reflect.DeepEqual(before, initialBefore) || !reflect.DeepEqual(after, initialAfter) {
		t.Fatal("comparison mutated caller snapshots")
	}
	if !reflect.DeepEqual(result, compareRun(t, before, after)) {
		t.Fatal("nondeterministic comparison")
	}
}

func TestDecisionTraceCompareLabelOnlyAndKindOnlyChangesRequireReview(t *testing.T) {
	for _, kindOnly := range []bool{false, true} {
		before := compareFixture(t)
		after := compareClone(t, before)
		for i := range after.Nodes {
			if after.Nodes[i].ID == "source" {
				if kindOnly {
					after.Nodes[i].Kind = domain.TraceInput
				} else {
					after.Nodes[i].Label = "different declared label"
				}
			}
		}
		after = compareRefreeze(t, after)
		result := compareRun(t, before, after)
		if len(result.NodeChanges) != 1 || result.NodeChanges[0].Status != "DECLARATION_CHANGED" || result.NodeChanges[0].NodeID != "source" {
			t.Fatalf("missing whole-node difference: %#v", result.NodeChanges)
		}
		if result.NodeChanges[0].Before.Fingerprint != result.NodeChanges[0].After.Fingerprint {
			t.Fatal("fixture must preserve fingerprint")
		}
		if result.Impact.Status != "REVIEW_REQUIRED" || !reflect.DeepEqual(result.Impact.ImpactedNodeIDs, []string{"claim", "decision", "source"}) || !reflect.DeepEqual(result.Impact.UnaffectedNodeIDs, []string{"unrelated"}) {
			t.Fatalf("whole-node declaration review missing: %#v", result.Impact)
		}
		for _, task := range result.Impact.Tasks {
			if task.NodeID == "source" && task.Reason != "INPUT_CHANGED" {
				t.Fatalf("snapshot event inferred as %s", task.Reason)
			}
		}
	}
}

func TestDecisionTraceCompareEdgeOnlyIncludesRefutations(t *testing.T) {
	before := compareFixture(t)
	after := compareClone(t, before)
	for i := range after.Edges {
		if after.Edges[i].From == "source" {
			after.Edges[i].Relation = domain.TraceSupports
		}
	}
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if len(result.NodeChanges) != 0 || len(result.AddedEdges) != 1 || len(result.RemovedEdges) != 1 || result.RemovedEdges[0].Relation != domain.TraceRefutes || result.AddedEdges[0].Relation != domain.TraceSupports {
		t.Fatalf("missing edge-only difference: %#v", result)
	}
	if !reflect.DeepEqual(result.Impact.ChangedNodeIDs, []string{"claim"}) || !reflect.DeepEqual(result.Impact.ImpactedNodeIDs, []string{"claim", "decision"}) {
		t.Fatalf("old consumer should trigger review: %#v", result.Impact)
	}
}

func TestDecisionTraceCompareFullSourceAttachmentAndStalePayload(t *testing.T) {
	t.Helper()
	before := compareFixture(t)
	source := domain.EvidenceSource{Kind: domain.EvidenceDocumentation, Ref: domain.CommitRef{Owner: "fixture", Repo: "local"}, Path: "README.md", Snippet: "declared source", Strength: 0.2}
	digest, err := traceFingerprint(source)
	if err != nil {
		t.Fatal(err)
	}
	for i := range before.Nodes {
		if before.Nodes[i].ID == "source" {
			before.Nodes[i].Fingerprint = digest
		}
	}
	before = compareRefreeze(t, before)
	after := compareClone(t, before)
	for i := range after.Nodes {
		if after.Nodes[i].ID == "source" {
			after.Nodes[i].Source = &source
		}
	}
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if len(result.NodeChanges) != 1 || result.NodeChanges[0].Status != "DECLARATION_CHANGED" || result.NodeChanges[0].Before.Source != nil || !reflect.DeepEqual(result.NodeChanges[0].After.Source, &source) || result.NodeChanges[0].Before.Fingerprint != result.NodeChanges[0].After.Fingerprint || result.Impact.Status != "REVIEW_REQUIRED" {
		t.Fatalf("full source attachment omitted despite equal fingerprint: %#v", result)
	}
	// Output source records must not alias the caller's mutable source pointer.
	result.NodeChanges[0].After.Source.Snippet = "output-only mutation"
	for _, node := range after.Nodes {
		if node.ID == "source" && node.Source.Snippet != "declared source" {
			t.Fatal("source payload aliases caller input")
		}
	}
	stale := compareClone(t, after)
	for i := range stale.Nodes {
		if stale.Nodes[i].Source != nil {
			stale.Nodes[i].Source.Snippet = "changed source without re-fingerprinting"
		}
	}
	got, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: stale})
	if err == nil || !reflect.DeepEqual(got, domain.DecisionTraceComparison{}) {
		t.Fatalf("stale full source payload accepted: %v", err)
	}
}

func TestDecisionTraceCompareAddedAndRemovedSnapshotsAreNotSourceEvents(t *testing.T) {
	before := compareFixture(t)
	after := compareClone(t, before)
	kept := []domain.DecisionTraceNode{}
	for _, node := range after.Nodes {
		if node.ID != "source" {
			kept = append(kept, node)
		}
	}
	after.Nodes = append(kept, domain.DecisionTraceNode{ID: "new-source", Kind: domain.TraceSource, Fingerprint: strings.Repeat("b", 64)})
	for i := range after.Edges {
		if after.Edges[i].From == "source" {
			after.Edges[i].From = "new-source"
		}
	}
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if len(result.NodeChanges) != 2 || result.NodeChanges[0].NodeID != "new-source" || result.NodeChanges[0].Status != "ADDED" || result.NodeChanges[0].Before != nil || result.NodeChanges[0].After == nil || result.NodeChanges[1].Status != "REMOVED" || result.NodeChanges[1].Before == nil || result.NodeChanges[1].After != nil {
		t.Fatalf("bad added/removed snapshot declarations: %#v", result.NodeChanges)
	}
	if !reflect.DeepEqual(result.Impact.ChangedNodeIDs, []string{"claim", "decision", "source"}) {
		t.Fatalf("new IDs must not enter old impact planner: %#v", result.Impact)
	}
	for _, task := range result.Impact.Tasks {
		if task.Reason == "REVOKED" || task.Reason == "EXPIRED" || task.NodeID == "new-source" {
			t.Fatal("snapshot difference fabricated a source event or old node")
		}
	}
	if !strings.Contains(strings.Join(result.Limitations, " "), "not source expiration, revocation") {
		t.Fatal("missing snapshot-event limitation")
	}
}

func TestDecisionTraceCompareCandidateDecisionIsNotPromoted(t *testing.T) {
	before := compareFixture(t)
	after := compareClone(t, before)
	after.Decision = domain.DecisionBuild
	for i := range after.Nodes {
		if after.Nodes[i].Kind == domain.TraceDecision {
			after.Nodes[i].Fingerprint, _ = traceFingerprint(after.Decision)
		}
	}
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if result.CandidateDecision != domain.DecisionBuild || result.HistoricalDecision != domain.DecisionUnknown || result.DecisionImpact != "NONE" || result.Impact.DecisionImpact != "NONE" || result.Impact.DecisionBefore != domain.DecisionUnknown || result.Impact.DecisionAfter != domain.DecisionUnknown {
		t.Fatalf("candidate declaration promoted: %#v", result)
	}
	if result.Impact.Status != "REVIEW_REQUIRED" {
		t.Fatal("supplied different decision must require historical review")
	}
}

func TestDecisionTraceCompareMetadataChangeRequiresConservativeReview(t *testing.T) {
	before := compareFixture(t)
	after := compareClone(t, before)
	after.Limitations = append(after.Limitations, "Additional declared snapshot limitation.")
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if !result.MetadataChanged || len(result.NodeChanges) != 0 || !reflect.DeepEqual(result.Impact.ChangedNodeIDs, []string{"decision"}) || result.Impact.Status != "REVIEW_REQUIRED" {
		t.Fatalf("metadata change disappeared: %#v", result)
	}
}

func TestDecisionTraceCompareRejectsInvalidAndAggregateBudgetWithoutPartialOutput(t *testing.T) {
	before := compareFixture(t)
	tests := map[string]func(*domain.DecisionTraceCompareInput){
		"schema":              func(in *domain.DecisionTraceCompareInput) { in.SchemaVersion = "2" },
		"missing before hash": func(in *domain.DecisionTraceCompareInput) { in.Before.GraphHash = "" },
		"missing after hash":  func(in *domain.DecisionTraceCompareInput) { in.After.GraphHash = "" },
		"tampered before":     func(in *domain.DecisionTraceCompareInput) { in.Before.Nodes[0].Label = "tampered" },
		"tampered after":      func(in *domain.DecisionTraceCompareInput) { in.After.Nodes[0].Label = "tampered" },
		"invalid UTF8":        func(in *domain.DecisionTraceCompareInput) { in.After.Nodes[0].Label = string([]byte{255}) },
		"overlong text":       func(in *domain.DecisionTraceCompareInput) { in.After.Nodes[0].Label = strings.Repeat("a", 8193) },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			in := domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: compareClone(t, before), After: compareClone(t, before)}
			mutate(&in)
			got, err := CompareDecisionTraces(in)
			if err == nil || !reflect.DeepEqual(got, domain.DecisionTraceComparison{}) {
				t.Fatalf("invalid input accepted or partial output returned: %v", err)
			}
		})
	}
	large := compareClone(t, before)
	for i := 0; i < 150; i++ {
		large.Nodes = append(large.Nodes, domain.DecisionTraceNode{ID: fmt.Sprintf("large-%03d", i), Kind: domain.TraceSource, Fingerprint: strings.Repeat("a", 64), Label: strings.Repeat("x", 8192)})
	}
	large = compareRefreeze(t, large)
	got, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: large, After: large})
	if err == nil || !reflect.DeepEqual(got, domain.DecisionTraceComparison{}) {
		t.Fatalf("combined input budget bypassed: %v", err)
	}
}

func TestDecisionTraceCompareRejectsMoreThan100HistoricalRoots(t *testing.T) {
	before := compareFixture(t)
	for i := 0; i < 101; i++ {
		before.Nodes = append(before.Nodes, domain.DecisionTraceNode{ID: fmt.Sprintf("root-%03d", i), Kind: domain.TraceSource, Fingerprint: strings.Repeat("a", 64)})
	}
	before = compareRefreeze(t, before)
	after := compareClone(t, before)
	for i := range after.Nodes {
		if strings.HasPrefix(after.Nodes[i].ID, "root-") {
			after.Nodes[i].Label = "changed declaration"
		}
	}
	after = compareRefreeze(t, after)
	got, err := CompareDecisionTraces(domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after})
	if err == nil || !strings.Contains(err.Error(), "review root budget") || !reflect.DeepEqual(got, domain.DecisionTraceComparison{}) {
		t.Fatalf("roots silently dropped or partial result: %v", err)
	}
	// The exact neighboring limit remains usable.
	for i := range after.Nodes {
		if after.Nodes[i].ID == "root-100" {
			after.Nodes[i].Label = ""
		}
	}
	after = compareRefreeze(t, after)
	result := compareRun(t, before, after)
	if len(result.NodeChanges) != 100 || len(result.Impact.ChangedNodeIDs) != 100 {
		t.Fatal("100-root positive neighbor incorrectly bounded")
	}
}

func TestDecisionTraceCompareRejectsExpandedOutputBeforeReturn(t *testing.T) {
	before := compareFixture(t)
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("%064x", i+1)
		before.Nodes = append(before.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceSource, Fingerprint: strings.Repeat("a", 64), Label: strings.Repeat("a", 8192)})
		before.Edges = append(before.Edges, domain.DecisionTraceEdge{From: id, To: "claim", Relation: domain.TraceInforms})
	}
	for i := 0; i < 80; i++ {
		id := fmt.Sprintf("leaf-%03d", i)
		before.Nodes = append(before.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceClaim, Fingerprint: strings.Repeat("a", 64)})
		before.Edges = append(before.Edges, domain.DecisionTraceEdge{From: "claim", To: id, Relation: domain.TraceDependsOn}, domain.DecisionTraceEdge{From: id, To: "decision", Relation: domain.TraceDependsOn})
	}
	before = compareRefreeze(t, before)
	after := compareClone(t, before)
	for i := range after.Nodes {
		if len(after.Nodes[i].Label) == 8192 {
			after.Nodes[i].Label = strings.Repeat("b", 8192)
		}
	}
	after = compareRefreeze(t, after)
	in := domain.DecisionTraceCompareInput{SchemaVersion: "1", Before: before, After: after}
	if _, err := traceBoundedJSON(in); err != nil {
		t.Fatalf("fixture aggregate input must fit: %v", err)
	}
	got, err := CompareDecisionTraces(in)
	if err == nil || !strings.Contains(err.Error(), "comparison output exceeds budget") || !reflect.DeepEqual(got, domain.DecisionTraceComparison{}) {
		t.Fatalf("expanded comparison must fail controlled quota without partial result: %v", err)
	}
}
