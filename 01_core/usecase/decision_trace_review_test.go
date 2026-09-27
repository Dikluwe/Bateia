/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_review_test.md
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

func traceReviewFixture(t *testing.T) domain.DecisionTraceReviewInput {
	t.Helper()
	return domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: traceFixture(t), Changes: []domain.DecisionTraceChange{{NodeID: "source-a", Reason: "EXPIRED"}}}}
}

func traceReviewRun(t *testing.T, in domain.DecisionTraceReviewInput) domain.DecisionTraceReviewLedger {
	t.Helper()
	out, err := TrackDecisionTraceReview(in)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func traceReviewRecord(ledger domain.DecisionTraceReviewLedger, index int, status string) domain.DecisionTraceReviewRecord {
	return domain.DecisionTraceReviewRecord{NodeID: ledger.Tasks[index].Task.NodeID, PlanHash: ledger.PlanHash, TaskHash: ledger.Tasks[index].TaskHash, Status: status, Note: "declared review result", Evidence: []domain.DecisionTraceReviewEvidence{{Ref: "logical-result", ContentHash: strings.Repeat("a", 64)}}}
}

func traceReviewAssertRejected(t *testing.T, in domain.DecisionTraceReviewInput) {
	t.Helper()
	out, err := TrackDecisionTraceReview(in)
	if err == nil || !reflect.DeepEqual(out, domain.DecisionTraceReviewLedger{}) {
		t.Fatalf("accepted invalid review or leaked partial output: %v %+v", err, out)
	}
}

func TestDecisionTraceReviewStatusesAndHistoricalDecision(t *testing.T) {
	in := traceReviewFixture(t)
	seed := traceReviewRun(t, in)
	if len(seed.Tasks) != 3 || seed.Counts.Pending != 3 || seed.Status != "REVIEW_PENDING" {
		t.Fatalf("unexpected default: %+v", seed)
	}
	for _, task := range seed.Tasks {
		if task.Status != "PENDING" || task.Evidence == nil || !traceValidDigest(task.TaskHash) {
			t.Fatal("invalid default task")
		}
	}
	in.Records = []domain.DecisionTraceReviewRecord{traceReviewRecord(seed, 0, "IN_PROGRESS")}
	if got := traceReviewRun(t, in); got.Status != "REVIEW_IN_PROGRESS" || got.Counts.InProgress != 1 || got.Counts.Pending != 2 {
		t.Fatalf("progress counts: %+v", got)
	}
	in.Records[0].Status = "BLOCKED"
	if got := traceReviewRun(t, in); got.Status != "REVIEW_BLOCKED" || got.Counts.Blocked != 1 {
		t.Fatal("blocked not visible")
	}
	in.Records = nil
	for i := range seed.Tasks {
		in.Records = append(in.Records, traceReviewRecord(seed, i, "COMPLETED"))
	}
	complete := traceReviewRun(t, in)
	if complete.Status != "RECORDED_REVIEW_COMPLETE" || complete.Counts.Completed != 3 || complete.DecisionBefore != domain.DecisionUnknown || complete.DecisionAfter != complete.DecisionBefore || complete.DecisionImpact != "NONE" || complete.Scope != domain.DecisionTraceScope {
		t.Fatalf("completion promoted decision: %+v", complete)
	}
	if !strings.Contains(strings.Join(complete.Limitations, " "), "supplied declarations") {
		t.Fatal("completion missing declaration warning")
	}
	in.ImpactInput.Changes = nil
	in.Records = nil
	zero := traceReviewRun(t, in)
	if zero.Status != "NO_REVIEW_TASKS" || len(zero.Tasks) != 0 || zero.Tasks == nil || zero.Counts != (domain.DecisionTraceReviewCounts{}) {
		t.Fatal("zero tasks fabricated completion")
	}
}

func TestDecisionTraceReviewCanonicalBindingAndNonMutation(t *testing.T) {
	in := traceReviewFixture(t)
	in.ImpactInput.Changes = append(in.ImpactInput.Changes, domain.DecisionTraceChange{NodeID: "source-b", Reason: "REVOKED"})
	seed := traceReviewRun(t, in)
	for i := range seed.Tasks {
		in.Records = append(in.Records, traceReviewRecord(seed, i, "COMPLETED"))
	}
	before, _ := json.Marshal(in)
	first := traceReviewRun(t, in)
	for i, j := 0, len(in.Records)-1; i < j; i, j = i+1, j-1 {
		in.Records[i], in.Records[j] = in.Records[j], in.Records[i]
	}
	in.ImpactInput.Changes[0], in.ImpactInput.Changes[1] = in.ImpactInput.Changes[1], in.ImpactInput.Changes[0]
	for i, j := 0, len(in.ImpactInput.Graph.Nodes)-1; i < j; i, j = i+1, j-1 {
		in.ImpactInput.Graph.Nodes[i], in.ImpactInput.Graph.Nodes[j] = in.ImpactInput.Graph.Nodes[j], in.ImpactInput.Graph.Nodes[i]
	}
	if second := traceReviewRun(t, in); !reflect.DeepEqual(first, second) {
		t.Fatal("input order changed projection")
	}
	var original domain.DecisionTraceReviewInput
	if err := json.Unmarshal(before, &original); err != nil {
		t.Fatal(err)
	}
	first.Tasks[0].Evidence[0].Ref = "output-only mutation"
	if original.Records[0].Evidence[0].Ref != "logical-result" {
		t.Fatal("original receipt mutated")
	}
	for _, r := range in.Records {
		if r.Evidence[0].Ref != "logical-result" {
			t.Fatal("output aliases caller reference slice")
		}
	}
	unchanged, _ := json.Marshal(original)
	traceReviewRun(t, original)
	after, _ := json.Marshal(original)
	if string(unchanged) != string(after) {
		t.Fatal("input mutated")
	}
}

func TestDecisionTraceReviewInvalidatesReceiptsForNewTargets(t *testing.T) {
	in := traceReviewFixture(t)
	in.ImpactInput.Changes[0] = domain.DecisionTraceChange{NodeID: "source-a", Reason: "CONTENT_CHANGED", NewFingerprint: strings.Repeat("b", 64)}
	old := traceReviewRun(t, in)
	in.Records = []domain.DecisionTraceReviewRecord{traceReviewRecord(old, 0, "COMPLETED")}
	in.ImpactInput.Changes[0].NewFingerprint = strings.Repeat("c", 64)
	traceReviewAssertRejected(t, in)
	in.Records = nil
	next := traceReviewRun(t, in)
	if old.PlanHash == next.PlanHash || old.Tasks[0].TaskHash == next.Tasks[0].TaskHash {
		t.Fatal("new target reused receipt binding")
	}
	in.Records = []domain.DecisionTraceReviewRecord{traceReviewRecord(next, 0, "COMPLETED")}
	in.ImpactInput.Changes[0].Reason = "REVOKED"
	traceReviewAssertRejected(t, in)
}

func TestDecisionTraceReviewRejectsInvalidRecords(t *testing.T) {
	tests := map[string]func(*domain.DecisionTraceReviewInput){
		"schema":                func(in *domain.DecisionTraceReviewInput) { in.SchemaVersion = "2" },
		"graph hash":            func(in *domain.DecisionTraceReviewInput) { in.ImpactInput.Graph.GraphHash = "" },
		"unknown task":          func(in *domain.DecisionTraceReviewInput) { in.Records[0].NodeID = "claim-b" },
		"duplicate task":        func(in *domain.DecisionTraceReviewInput) { in.Records = append(in.Records, in.Records[0]) },
		"plan hash":             func(in *domain.DecisionTraceReviewInput) { in.Records[0].PlanHash = strings.Repeat("b", 64) },
		"task hash":             func(in *domain.DecisionTraceReviewInput) { in.Records[0].TaskHash = strings.Repeat("b", 64) },
		"status":                func(in *domain.DecisionTraceReviewInput) { in.Records[0].Status = "APPROVED" },
		"empty completion note": func(in *domain.DecisionTraceReviewInput) { in.Records[0].Note = " \t" },
		"no completion ref":     func(in *domain.DecisionTraceReviewInput) { in.Records[0].Evidence = nil },
		"empty blocked note":    func(in *domain.DecisionTraceReviewInput) { in.Records[0].Status = "BLOCKED"; in.Records[0].Note = "" },
		"empty ref":             func(in *domain.DecisionTraceReviewInput) { in.Records[0].Evidence[0].Ref = " " },
		"hash syntax": func(in *domain.DecisionTraceReviewInput) {
			in.Records[0].Evidence[0].ContentHash = "sha256:" + strings.Repeat("a", 64)
		},
		"uppercase hash": func(in *domain.DecisionTraceReviewInput) {
			in.Records[0].Evidence[0].ContentHash = strings.Repeat("A", 64)
		},
		"invalid UTF8": func(in *domain.DecisionTraceReviewInput) { in.Records[0].Note = string([]byte{255}) },
		"long note":    func(in *domain.DecisionTraceReviewInput) { in.Records[0].Note = strings.Repeat("x", 8193) },
		"21 references": func(in *domain.DecisionTraceReviewInput) {
			e := in.Records[0].Evidence[0]
			for len(in.Records[0].Evidence) < 21 {
				in.Records[0].Evidence = append(in.Records[0].Evidence, e)
			}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			in := traceReviewFixture(t)
			seed := traceReviewRun(t, in)
			in.Records = []domain.DecisionTraceReviewRecord{traceReviewRecord(seed, 0, "COMPLETED")}
			mutate(&in)
			traceReviewAssertRejected(t, in)
		})
	}
	in := traceReviewFixture(t)
	seed := traceReviewRun(t, in)
	r := traceReviewRecord(seed, 0, "COMPLETED")
	for len(r.Evidence) < 20 {
		r.Evidence = append(r.Evidence, r.Evidence[0])
	}
	in.Records = []domain.DecisionTraceReviewRecord{r}
	traceReviewRun(t, in)
}

func TestDecisionTraceReviewRecordBoundary(t *testing.T) {
	in := traceReviewFixture(t)
	g := in.ImpactInput.Graph
	g.GraphHash = ""
	g.Nodes = g.Nodes[:0]
	g.Edges = nil
	fp, _ := traceFingerprint(domain.DecisionUnknown)
	g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: "decision", Kind: domain.TraceDecision, Fingerprint: fp}, domain.DecisionTraceNode{ID: "source-a", Kind: domain.TraceSource, Fingerprint: strings.Repeat("a", 64)})
	g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: "source-a", To: "decision", Relation: domain.TraceRefutes})
	for i := 0; i < 998; i++ {
		id := fmt.Sprintf("claim-%04d", i)
		g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceClaim, Fingerprint: strings.Repeat("a", 64)})
		g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: "source-a", To: id, Relation: domain.TraceInforms}, domain.DecisionTraceEdge{From: id, To: "decision", Relation: domain.TraceDependsOn})
	}
	var err error
	in.ImpactInput.Graph, err = FreezeDecisionTrace(g)
	if err != nil {
		t.Fatal(err)
	}
	seed := traceReviewRun(t, in)
	if len(seed.Tasks) != 1000 {
		t.Fatal("boundary fixture not 1000 tasks")
	}
	for i := range seed.Tasks {
		r := traceReviewRecord(seed, i, "PENDING")
		r.Note = ""
		r.Evidence = nil
		in.Records = append(in.Records, r)
	}
	traceReviewRun(t, in)
	in.Records = append(in.Records, in.Records[0])
	traceReviewAssertRejected(t, in)
}

func TestDecisionTraceReviewExpandedOutputBudget(t *testing.T) {
	in := traceReviewFixture(t)
	g := in.ImpactInput.Graph
	g.GraphHash = ""
	g.Nodes = nil
	g.Edges = nil
	in.ImpactInput.Changes = nil
	fp, _ := traceFingerprint(domain.DecisionUnknown)
	g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: "decision", Kind: domain.TraceDecision, Fingerprint: fp}, domain.DecisionTraceNode{ID: "junction", Kind: domain.TraceClaim, Fingerprint: strings.Repeat("a", 64)})
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("%064x", i+1)
		g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceSource, Fingerprint: strings.Repeat("a", 64)})
		g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: id, To: "junction", Relation: domain.TraceInforms})
		in.ImpactInput.Changes = append(in.ImpactInput.Changes, domain.DecisionTraceChange{NodeID: id, Reason: "EXPIRED"})
	}
	for i := 0; i < 80; i++ {
		id := fmt.Sprintf("leaf-%03d", i)
		g.Nodes = append(g.Nodes, domain.DecisionTraceNode{ID: id, Kind: domain.TraceClaim, Fingerprint: strings.Repeat("a", 64)})
		g.Edges = append(g.Edges, domain.DecisionTraceEdge{From: "junction", To: id, Relation: domain.TraceDependsOn}, domain.DecisionTraceEdge{From: id, To: "decision", Relation: domain.TraceDependsOn})
	}
	var err error
	in.ImpactInput.Graph, err = FreezeDecisionTrace(g)
	if err != nil {
		t.Fatal(err)
	}
	seed := traceReviewRun(t, in)
	for i := range seed.Tasks {
		r := traceReviewRecord(seed, i, "IN_PROGRESS")
		r.Note = strings.Repeat("x", 8192)
		in.Records = append(in.Records, r)
	}
	if _, err := traceBoundedJSON(in); err != nil {
		t.Fatalf("input must fit: %v", err)
	}
	out, err := TrackDecisionTraceReview(in)
	if err == nil || !strings.Contains(err.Error(), "review output exceeds budget") || !reflect.DeepEqual(out, domain.DecisionTraceReviewLedger{}) {
		t.Fatalf("expanded output not controlled: %v", err)
	}
	for i := range in.Records {
		in.Records[i].Note = strings.Repeat("x", 4096)
	}
	traceReviewRun(t, in)
}
