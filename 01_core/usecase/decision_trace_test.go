/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionTraceImpactRejectsExpandedOutputBeforeJSONAllocation(t *testing.T) {
	digest := strings.Repeat("a", 64)
	decisionDigest, err := traceFingerprint(domain.DecisionUnknown)
	if err != nil {
		t.Fatal(err)
	}
	graph := domain.DecisionTraceGraph{
		SchemaVersion: domain.DecisionTraceSchemaVersion, Scope: domain.DecisionTraceScope, Decision: domain.DecisionUnknown,
		Nodes: []domain.DecisionTraceNode{{ID: "decision", Kind: domain.TraceDecision, Fingerprint: decisionDigest}, {ID: "junction", Kind: domain.TraceClaim, Fingerprint: digest}},
	}
	changes := []domain.DecisionTraceChange{}
	for i := 1; i <= 100; i++ {
		id := fmt.Sprintf("%064x", i)
		graph.Nodes = append(graph.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceSource, Fingerprint: digest})
		graph.Edges = append(graph.Edges, domain.DecisionTraceEdge{From: id, To: "junction", Relation: domain.TraceInforms})
		changes = append(changes, domain.DecisionTraceChange{NodeID: id, Reason: "EXPIRED"})
	}
	for i := 0; i < 700; i++ {
		id := fmt.Sprintf("leaf-%03d", i)
		graph.Nodes = append(graph.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceClaim, Fingerprint: digest})
		graph.Edges = append(graph.Edges, domain.DecisionTraceEdge{From: "junction", To: id, Relation: domain.TraceDependsOn}, domain.DecisionTraceEdge{From: id, To: "decision", Relation: domain.TraceDependsOn})
	}
	graph, err = FreezeDecisionTrace(graph)
	if err != nil {
		t.Fatalf("bounded source graph rejected: %v", err)
	}
	input := domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: changes}
	encoded, err := traceBoundedJSON(input)
	if err != nil || len(encoded) > traceMaxJSON {
		t.Fatalf("fixture input must fit budget: %v", err)
	}
	plan, err := PlanDecisionTraceImpact(input)
	if err == nil || !strings.Contains(err.Error(), "impact output exceeds budget") {
		t.Fatalf("expanded output should fail under controlled quota: %v", err)
	}
	if !reflect.DeepEqual(plan, domain.DecisionTraceImpactPlan{}) {
		t.Fatal("partial oversized plan escaped")
	}
	// Same graph with one changed ancestor must remain usable: quota failure is
	// caused by the output expansion, not an unconditional ban on larger DAGs.
	input.Changes = changes[:1]
	plan, err = PlanDecisionTraceImpact(input)
	if err != nil || plan.Status != "REVIEW_REQUIRED" || len(plan.Tasks) != 703 {
		t.Fatalf("bounded positive plan rejected: tasks=%d err=%v", len(plan.Tasks), err)
	}
}

func traceFixture(t *testing.T) domain.DecisionTraceGraph {
	t.Helper()
	digest := strings.Repeat("a", 64)
	decisionDigest, _ := traceFingerprint(domain.DecisionUnknown)
	g, err := FreezeDecisionTrace(domain.DecisionTraceGraph{
		SchemaVersion: domain.DecisionTraceSchemaVersion, Scope: domain.DecisionTraceScope, Decision: domain.DecisionUnknown,
		Nodes: []domain.DecisionTraceNode{{ID: "decision", Kind: domain.TraceDecision, Fingerprint: decisionDigest}, {ID: "source-a", Kind: domain.TraceSource, Fingerprint: digest}, {ID: "claim-a", Kind: domain.TraceClaim, Fingerprint: digest}, {ID: "source-b", Kind: domain.TraceSource, Fingerprint: digest}, {ID: "claim-b", Kind: domain.TraceClaim, Fingerprint: digest}},
		Edges: []domain.DecisionTraceEdge{{From: "source-a", To: "claim-a", Relation: domain.TraceInforms}, {From: "source-b", To: "claim-b", Relation: domain.TraceRefutes}, {From: "claim-a", To: "decision", Relation: domain.TraceDependsOn}, {From: "claim-b", To: "decision", Relation: domain.TraceSupports}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func traceCloneGraph(t *testing.T, g domain.DecisionTraceGraph) domain.DecisionTraceGraph {
	t.Helper()
	b, _ := json.Marshal(g)
	var out domain.DecisionTraceGraph
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestDecisionTraceFreezeDeterministicNonMutating(t *testing.T) {
	g := traceFixture(t)
	original := traceCloneGraph(t, g)
	for i, j := 0, len(g.Nodes)-1; i < j; i, j = i+1, j-1 {
		g.Nodes[i], g.Nodes[j] = g.Nodes[j], g.Nodes[i]
	}
	for i, j := 0, len(g.Edges)-1; i < j; i, j = i+1, j-1 {
		g.Edges[i], g.Edges[j] = g.Edges[j], g.Edges[i]
	}
	unsorted := traceCloneGraph(t, g)
	got, err := FreezeDecisionTrace(g)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g, unsorted) {
		t.Fatal("caller graph mutated")
	}
	if !reflect.DeepEqual(got, original) {
		t.Fatal("canonical output depends on input order")
	}
	again, err := FreezeDecisionTrace(got)
	if err != nil || !reflect.DeepEqual(got, again) {
		t.Fatalf("not idempotent: %v", err)
	}
	if got.Nodes == nil || got.Edges == nil || got.Limitations == nil {
		t.Fatal("nil vectors")
	}
}

func TestDecisionTraceFreezeRejectsMalformedGraph(t *testing.T) {
	base := traceFixture(t)
	tests := map[string]func(*domain.DecisionTraceGraph){
		"tampered":       func(g *domain.DecisionTraceGraph) { g.Nodes[0].Label = "tampered" },
		"duplicate node": func(g *domain.DecisionTraceGraph) { g.Nodes = append(g.Nodes, g.Nodes[0]); g.GraphHash = "" },
		"duplicate edge": func(g *domain.DecisionTraceGraph) { g.Edges = append(g.Edges, g.Edges[0]); g.GraphHash = "" },
		"dangling":       func(g *domain.DecisionTraceGraph) { g.Edges[0].From = "unknown"; g.GraphHash = "" },
		"cycle": func(g *domain.DecisionTraceGraph) {
			g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: "decision", To: "source-a", Relation: domain.TraceRefutes})
			g.GraphHash = ""
		},
		"self":        func(g *domain.DecisionTraceGraph) { g.Edges[0].To = g.Edges[0].From; g.GraphHash = "" },
		"relation":    func(g *domain.DecisionTraceGraph) { g.Edges[0].Relation = "APPROVES"; g.GraphHash = "" },
		"node kind":   func(g *domain.DecisionTraceGraph) { g.Nodes[0].Kind = "PROOF"; g.GraphHash = "" },
		"decision":    func(g *domain.DecisionTraceGraph) { g.Decision = "APPROVED"; g.GraphHash = "" },
		"scope":       func(g *domain.DecisionTraceGraph) { g.Scope = "COMPLETE"; g.GraphHash = "" },
		"schema":      func(g *domain.DecisionTraceGraph) { g.SchemaVersion = "2"; g.GraphHash = "" },
		"fingerprint": func(g *domain.DecisionTraceGraph) { g.Nodes[0].Fingerprint = "a"; g.GraphHash = "" },
		"ID control":  func(g *domain.DecisionTraceGraph) { g.Nodes[0].ID = "bad\x00id"; g.GraphHash = "" },
		"no decision": func(g *domain.DecisionTraceGraph) {
			for i := range g.Nodes {
				if g.Nodes[i].Kind == domain.TraceDecision {
					g.Nodes[i].Kind = domain.TraceClaim
				}
			}
			g.GraphHash = ""
		},
		"two decisions": func(g *domain.DecisionTraceGraph) { g.Nodes[0].Kind = domain.TraceDecision; g.GraphHash = "" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			g := traceCloneGraph(t, base)
			mutate(&g)
			out, err := FreezeDecisionTrace(g)
			if err == nil || len(out.Nodes) != 0 || out.GraphHash != "" {
				t.Fatalf("malformed accepted/partial: %v", err)
			}
		})
	}
}

func traceBuildFixture(t *testing.T) (domain.DecisionTraceBuildInput, domain.DecisionTraceGraph) {
	t.Helper()
	a := domain.EvidenceSource{Kind: domain.EvidenceDocumentation, Path: "a.md", RuleID: "a", MatchClass: domain.EvidenceMatchPathOnly, Snippet: "declared", Strength: 0.1}
	b := domain.EvidenceSource{Kind: domain.EvidenceFile, Path: "b.go", RuleID: "b", Strength: 0.3}
	input := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{{Capability: "alpha", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "candidate-a", Status: domain.EvidenceInconclusive, Evidence: []domain.EvidenceSource{a}}}}, {Capability: "beta", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "candidate-b", Status: domain.EvidencePartial, Evidence: []domain.EvidenceSource{b}}}}}}}
	g, err := BuildDecisionTrace(input)
	if err != nil {
		t.Fatal(err)
	}
	return input, g
}

func TestDecisionTraceBuildPreservesDeclarationsAndReport(t *testing.T) {
	input, g := traceBuildFixture(t)
	before, _ := json.Marshal(input)
	again, err := BuildDecisionTrace(input)
	after, _ := json.Marshal(input)
	if err != nil || !reflect.DeepEqual(g, again) || string(before) != string(after) {
		t.Fatal("build mutated input or is nondeterministic")
	}
	if g.Decision != domain.DecisionUnknown {
		t.Fatal("historical decision changed")
	}
	sources := 0
	for _, node := range g.Nodes {
		if node.Kind == domain.TraceSource {
			sources++
			if node.Source == nil || node.Source.ContentHash != "" {
				t.Fatal("missing/forged declaration")
			}
			if node.Source.Path == "a.md" && (node.Source.MatchClass != domain.EvidenceMatchPathOnly || node.Source.Strength != 0.1) {
				t.Fatal("source was promoted")
			}
		}
	}
	if sources != 2 {
		t.Fatalf("sources=%d", sources)
	}
	for _, edge := range g.Edges {
		if edge.Relation == domain.TraceSupports || edge.Relation == domain.TraceRefutes {
			t.Fatal("support/refutation invented")
		}
	}
}

func TestDecisionTraceBuildRejectsAmbiguousSource(t *testing.T) {
	for _, field := range []string{"content", "snippet", "strength"} {
		t.Run(field, func(t *testing.T) {
			input, _ := traceBuildFixture(t)
			source := input.Report.Matrix[0].Candidates[0].Evidence[0]
			switch field {
			case "content":
				source.ContentHash = strings.Repeat("b", 64)
			case "snippet":
				source.Snippet = "different"
			case "strength":
				source.Strength = 0.9
			}
			input.Report.Matrix[1].Candidates[0].Evidence = append(input.Report.Matrix[1].Candidates[0].Evidence, source)
			if _, err := BuildDecisionTrace(input); err == nil {
				t.Fatal("ambiguous source accepted")
			}
		})
	}
}

func TestDecisionTraceImpactSelectiveAndNoRecalculation(t *testing.T) {
	g := traceFixture(t)
	before := traceCloneGraph(t, g)
	plan, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: []domain.DecisionTraceChange{{NodeID: "source-a", Reason: "EXPIRED"}}})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != "REVIEW_REQUIRED" || plan.DecisionImpact != "NONE" || plan.DecisionBefore != g.Decision || plan.DecisionAfter != g.Decision {
		t.Fatal("impact altered decision")
	}
	if !reflect.DeepEqual(plan.ImpactedNodeIDs, []string{"claim-a", "decision", "source-a"}) || !reflect.DeepEqual(plan.UnaffectedNodeIDs, []string{"claim-b", "source-b"}) {
		t.Fatalf("wrong branch: %+v", plan)
	}
	if !reflect.DeepEqual(g, before) {
		t.Fatal("historical graph mutated")
	}
	for _, task := range plan.Tasks {
		if !reflect.DeepEqual(task.ChangedAncestorIDs, []string{"source-a"}) {
			t.Fatal("missing ancestor")
		}
	}
}

func TestDecisionTraceImpactConfluenceAndRefutation(t *testing.T) {
	g := traceFixture(t)
	changes := []domain.DecisionTraceChange{{NodeID: "source-b", Reason: "REVOKED"}, {NodeID: "source-a", Reason: "CONTENT_CHANGED", NewFingerprint: strings.Repeat("b", 64)}}
	a, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: changes})
	if err != nil {
		t.Fatal(err)
	}
	changes[0], changes[1] = changes[1], changes[0]
	b, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: changes})
	if err != nil || !reflect.DeepEqual(a, b) {
		t.Fatal("change order affects plan")
	}
	if len(a.ImpactedNodeIDs) != 5 {
		t.Fatal("REFUTES did not propagate")
	}
	for _, task := range a.Tasks {
		if task.NodeID == "decision" && !reflect.DeepEqual(task.ChangedAncestorIDs, []string{"source-a", "source-b"}) {
			t.Fatal("convergent ancestors lost")
		}
	}
}

func TestDecisionTraceImpactNoopAndExpiration(t *testing.T) {
	g := traceFixture(t)
	for _, reason := range []string{"CONTENT_CHANGED", "EXPIRED", "REVOKED"} {
		plan, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: []domain.DecisionTraceChange{{NodeID: "source-a", Reason: reason, NewFingerprint: strings.Repeat("a", 64)}}})
		if err != nil {
			t.Fatal(err)
		}
		if reason == "CONTENT_CHANGED" {
			if plan.Status != "NO_RECORDED_IMPACT" || len(plan.Tasks) != 0 || len(plan.UnaffectedNodeIDs) != 5 {
				t.Fatal("identical content not noop")
			}
		} else if plan.Status != "REVIEW_REQUIRED" {
			t.Fatal("expiration/revocation became noop")
		}
	}
}

func TestDecisionTraceImpactRejectsInvalidChanges(t *testing.T) {
	g := traceFixture(t)
	for name, changes := range map[string][]domain.DecisionTraceChange{
		"unknown": {{NodeID: "unknown", Reason: "EXPIRED"}}, "duplicate": {{NodeID: "source-a", Reason: "EXPIRED"}, {NodeID: "source-a", Reason: "REVOKED"}}, "reason": {{NodeID: "source-a", Reason: "APPROVED"}}, "digest": {{NodeID: "source-a", Reason: "CONTENT_CHANGED", NewFingerprint: "sha256:a"}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: changes}); err == nil {
				t.Fatal("invalid change accepted")
			}
		})
	}
	g.GraphHash = ""
	if _, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g}); err == nil {
		t.Fatal("unfrozen graph accepted")
	}
}

func TestDecisionTraceLimitsBeforeSerialization(t *testing.T) {
	g := traceFixture(t)
	g.GraphHash = ""
	g.Nodes = make([]domain.DecisionTraceNode, 1001)
	if _, err := FreezeDecisionTrace(g); err == nil {
		t.Fatal("node budget ignored")
	}
	g = traceFixture(t)
	g.GraphHash = ""
	g.Edges = make([]domain.DecisionTraceEdge, 4001)
	if _, err := FreezeDecisionTrace(g); err == nil {
		t.Fatal("edge budget ignored")
	}
	input, _ := traceBuildFixture(t)
	input.Report.Gaps = []string{strings.Repeat("x", 8193)}
	if _, err := BuildDecisionTrace(input); err == nil {
		t.Fatal("string budget ignored")
	}
	input.Report.Gaps = make([]string, 4001)
	if _, err := BuildDecisionTrace(input); err == nil {
		t.Fatal("vector budget ignored")
	}
	input.Report.Gaps = make([]string, 300)
	for i := range input.Report.Gaps {
		input.Report.Gaps[i] = strings.Repeat("x", 8192)
	}
	if _, err := BuildDecisionTrace(input); err == nil {
		t.Fatal("aggregate budget ignored")
	}
	input.Report.Gaps = []string{string([]byte{255})}
	if _, err := BuildDecisionTrace(input); err == nil {
		t.Fatal("invalid UTF-8 accepted")
	}
	if _, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: traceFixture(t), Changes: make([]domain.DecisionTraceChange, 101)}); err == nil {
		t.Fatal("change budget ignored")
	}
}

func TestDecisionTraceSourceSnapshotAndDeclarationConsistency(t *testing.T) {
	_, graph := traceBuildFixture(t)
	frozen, err := FreezeDecisionTrace(graph)
	if err != nil {
		t.Fatal(err)
	}
	for i := range frozen.Nodes {
		if frozen.Nodes[i].Source != nil {
			frozen.Nodes[i].Source.Snippet = "mutated copy"
			if graph.Nodes[i].Source.Snippet == "mutated copy" {
				t.Fatal("returned source aliases input")
			}
			break
		}
	}
	graph.GraphHash = ""
	graph.Decision = domain.DecisionBuild
	if _, err = FreezeDecisionTrace(graph); err == nil {
		t.Fatal("decision fingerprint inconsistent with declaration accepted")
	}
	_, graph = traceBuildFixture(t)
	graph.GraphHash = ""
	for _, node := range graph.Nodes {
		if node.Source != nil {
			duplicate := node
			duplicate.ID = "different-id-same-source"
			graph.Nodes = append(graph.Nodes, duplicate)
			break
		}
	}
	if _, err = FreezeDecisionTrace(graph); err == nil {
		t.Fatal("duplicate source identity accepted under different IDs")
	}
}

func TestDecisionTraceSourcesNonFiniteAndStructuredIdentity(t *testing.T) {
	for _, strength := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		input, _ := traceBuildFixture(t)
		input.Report.Matrix[0].Candidates[0].Evidence[0].Strength = strength
		if _, err := BuildDecisionTrace(input); err == nil {
			t.Fatal("non-finite source accepted")
		}
	}
	input, _ := traceBuildFixture(t)
	a := input.Report.Matrix[0].Candidates[0].Evidence[0]
	b := a
	a.Path, a.Symbol = "x\x00y", "z"
	b.Path, b.Symbol = "x", "y\x00z"
	input.Report.Matrix[0].Candidates[0].Evidence = []domain.EvidenceSource{a, b}
	graph, err := BuildDecisionTrace(input)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, node := range graph.Nodes {
		if node.Source != nil && node.Source.RuleID == a.RuleID {
			ids[node.ID] = true
		}
	}
	if len(ids) != 2 {
		t.Fatal("structured source keys collapsed NUL-delimited tuple collision")
	}
}

func TestDecisionTraceDirectJudgmentChangesRemainHistorical(t *testing.T) {
	graph := traceFixture(t)
	for _, id := range []string{"claim-a", "decision"} {
		plan, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: id, Reason: "JUDGMENT_CHANGED"}}})
		if err != nil {
			t.Fatal(err)
		}
		if plan.Status != "REVIEW_REQUIRED" || plan.DecisionAfter != graph.Decision {
			t.Fatal("direct judgment change silently applied")
		}
	}
}

func TestDecisionTraceBuildCompositionTechnicalAndInputDependencies(t *testing.T) {
	input, _ := traceBuildFixture(t)
	input.RuleVersions = map[string]string{"evidence": "1"}
	input.Idea = &domain.IdeaSpec{}
	input.Report.Reasons = []domain.DecisionReason{{RuleID: "r", Message: "declared"}}
	input.Report.Blockers = []domain.DecisionReason{{RuleID: "b", Message: "blocked"}}
	input.Report.Components = []domain.ReusableComponent{{Name: "component", Kind: "library"}}
	input.Report.Composition = &domain.CompositionReport{Candidates: []domain.CompositionCandidate{{ID: "composition", Status: domain.CompositionProposed}}}
	input.Report.TechnicalEvidence = &domain.TechnicalEvidenceReview{Claims: []domain.TechnicalClaimReview{{Candidate: "a", Claim: "untrusted external declaration", Status: domain.TechnicalClaimSupported}}, Verifications: []domain.TechnicalEvidenceVerification{{FindingID: "f", Candidate: "a", Status: domain.TechnicalUnverified}}}
	g, err := BuildDecisionTrace(input)
	if err != nil {
		t.Fatal(err)
	}
	inputs, claims := 0, 0
	policyID := ""
	for _, n := range g.Nodes {
		if n.Kind == domain.TraceInput {
			inputs++
			if n.Label == "policy" {
				policyID = n.ID
			}
		}
		if n.Kind == domain.TraceClaim {
			claims++
		}
	}
	if inputs != 4 || claims != 10 {
		t.Fatalf("missing declarations inputs=%d claims=%d", inputs, claims)
	}
	plan, err := PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: g, Changes: []domain.DecisionTraceChange{{NodeID: policyID, Reason: "POLICY_CHANGED"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Tasks) != claims+2 || plan.DecisionAfter != domain.DecisionUnknown {
		t.Fatal("input change was not conservative or verdict changed")
	}
}
