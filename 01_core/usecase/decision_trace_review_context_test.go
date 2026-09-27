/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_review_context_test.md
 * @layer L1
 * @updated 2026-09-27
 */
package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func traceReviewContextFixture(t *testing.T) domain.DecisionTraceReviewContextInput {
	t.Helper()
	source := domain.EvidenceSource{Kind: domain.EvidenceDocumentation, Ref: domain.CommitRef{Owner: "fixture", Repo: "app"}, Path: "README.md", RuleID: "docs-rule", Snippet: "declared capability", ContentHash: "sha256:" + strings.Repeat("a", 64)}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{{Capability: "read", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "app", Status: domain.EvidencePartial, Evidence: []domain.EvidenceSource{source}}}}, {Capability: "write", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "app", Status: domain.EvidencePartial, Evidence: []domain.EvidenceSource{source}}}}}, Reasons: []domain.DecisionReason{{RuleID: "review-reason", Message: strings.Repeat("ação", 500), Evidence: []domain.EvidenceSource{source}}}}}
	graph, err := BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	var sourceID string
	for _, node := range graph.Nodes {
		if node.Kind == domain.TraceSource {
			sourceID = node.ID
		}
	}
	if sourceID == "" {
		t.Fatal("fixture requires source")
	}
	return domain.DecisionTraceReviewContextInput{SchemaVersion: "1", BuildInput: build, ReviewInput: domain.DecisionTraceReviewInput{SchemaVersion: "1", ImpactInput: domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{{NodeID: sourceID, Reason: "EXPIRED"}}}}}
}

func traceReviewContextRun(t *testing.T, in domain.DecisionTraceReviewContextInput) domain.DecisionTraceReviewContext {
	t.Helper()
	out, err := BuildDecisionTraceReviewContext(in)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestDecisionTraceReviewContextJoinsOriginalLocationsAndRoots(t *testing.T) {
	in := traceReviewContextFixture(t)
	out := traceReviewContextRun(t, in)
	ledger := traceReviewRun(t, in.ReviewInput)
	if !reflect.DeepEqual(out.Ledger, ledger) || out.PlanHash != ledger.PlanHash || out.GraphHash != ledger.GraphHash || out.Scope != domain.DecisionTraceScope || out.ContextScope != domain.DecisionTraceReviewContextScope || len(out.TaskContexts) != len(ledger.Tasks) {
		t.Fatal("context changed review identity")
	}
	sourceID := in.ReviewInput.ImpactInput.Changes[0].NodeID
	selected := map[string]bool{}
	for _, context := range out.TaskContexts {
		if !sort.StringsAreSorted(context.ContextNodeIDs) {
			t.Fatal("context IDs not canonical")
		}
		expected := []string{context.NodeID}
		if sourceID != context.NodeID {
			expected = append(expected, sourceID)
		}
		sort.Strings(expected)
		if !reflect.DeepEqual(context.ContextNodeIDs, expected) {
			t.Fatal("context contains false/missing dependencies")
		}
		for _, id := range context.ContextNodeIDs {
			selected[id] = true
		}
	}
	ids := []string{}
	for id := range selected {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	explanation, err := ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: in.ReviewInput.ImpactInput.Graph, BuildInput: in.BuildInput, NodeIDs: ids})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out.Entries, explanation.Entries) || out.ReportFingerprint != explanation.ReportFingerprint || len(out.Entries) != len(ids) {
		t.Fatal("context differs from historical explanation")
	}
	foundSource, foundTruncation := false, false
	for _, entry := range out.Entries {
		if entry.NodeID == sourceID {
			foundSource = true
			if len(entry.Locations) != 3 {
				t.Fatal("shared source occurrences lost")
			}
			for _, location := range entry.Locations {
				if location.Document != "report" || !strings.HasSuffix(location.Pointer, "/evidence/0") {
					t.Fatal("incorrect source document pointer")
				}
			}
		}
		if entry.StatementTruncated {
			foundTruncation = true
			if !reflect.DeepEqual(entry.Locations, []domain.DecisionTraceLocation{{Document: "report", Pointer: "/reasons/0"}}) {
				t.Fatal("truncation lost original location")
			}
		}
	}
	if !foundSource || !foundTruncation {
		t.Fatal("missing source/truncation projection")
	}
}

func TestDecisionTraceReviewContextPreservesCompletionAndDoesNotMutate(t *testing.T) {
	in := traceReviewContextFixture(t)
	seed := traceReviewRun(t, in.ReviewInput)
	for i := range seed.Tasks {
		in.ReviewInput.Records = append(in.ReviewInput.Records, traceReviewRecord(seed, i, "COMPLETED"))
	}
	original, _ := json.Marshal(in)
	out := traceReviewContextRun(t, in)
	after, _ := json.Marshal(in)
	if string(original) != string(after) {
		t.Fatal("input mutated")
	}
	if out.Ledger.Status != "RECORDED_REVIEW_COMPLETE" || out.Ledger.DecisionBefore != domain.DecisionUnknown || out.Ledger.DecisionAfter != out.Ledger.DecisionBefore || out.Ledger.DecisionImpact != "NONE" {
		t.Fatal("document context promoted completion")
	}
	for i, j := 0, len(in.ReviewInput.Records)-1; i < j; i, j = i+1, j-1 {
		in.ReviewInput.Records[i], in.ReviewInput.Records[j] = in.ReviewInput.Records[j], in.ReviewInput.Records[i]
	}
	for i, j := 0, len(in.ReviewInput.ImpactInput.Graph.Nodes)-1; i < j; i, j = i+1, j-1 {
		in.ReviewInput.ImpactInput.Graph.Nodes[i], in.ReviewInput.ImpactInput.Graph.Nodes[j] = in.ReviewInput.ImpactInput.Graph.Nodes[j], in.ReviewInput.ImpactInput.Graph.Nodes[i]
	}
	if !reflect.DeepEqual(out, traceReviewContextRun(t, in)) {
		t.Fatal("input order changed context")
	}
	out.Ledger.Tasks[0].Evidence[0].Ref = "output-only"
	for _, record := range in.ReviewInput.Records {
		if record.Evidence[0].Ref != "logical-result" {
			t.Fatal("output aliases input references")
		}
	}
}

func TestDecisionTraceReviewContextEmptyStillBindsOriginalDocuments(t *testing.T) {
	in := traceReviewContextFixture(t)
	in.ReviewInput.ImpactInput.Changes = nil
	out := traceReviewContextRun(t, in)
	if out.Ledger.Status != "NO_REVIEW_TASKS" || len(out.TaskContexts) != 0 || len(out.Entries) != 0 || out.Entries == nil || out.TaskContexts == nil || out.ReportFingerprint == "" {
		t.Fatal("zero tasks projected all or fabricated completion")
	}
	in.BuildInput.Report.Decision = domain.DecisionBuild
	result, err := BuildDecisionTraceReviewContext(in)
	if err == nil || !reflect.DeepEqual(result, domain.DecisionTraceReviewContext{}) {
		t.Fatal("empty review accepted mismatching documents")
	}
}

func TestDecisionTraceReviewContextRejectsTampering(t *testing.T) {
	cases := map[string]func(*domain.DecisionTraceReviewContextInput){
		"schema": func(in *domain.DecisionTraceReviewContextInput) { in.SchemaVersion = "2" },
		"report": func(in *domain.DecisionTraceReviewContextInput) {
			in.BuildInput.Report.Reasons[0].Message = "different historical declaration"
		},
		"idea": func(in *domain.DecisionTraceReviewContextInput) {
			in.BuildInput.Idea = &domain.IdeaSpec{}
		},
		"rules": func(in *domain.DecisionTraceReviewContextInput) {
			in.BuildInput.RuleVersions = map[string]string{"rules": "v2"}
		},
		"hash": func(in *domain.DecisionTraceReviewContextInput) {
			in.ReviewInput.ImpactInput.Graph.GraphHash = strings.Repeat("b", 64)
		},
		"UTF8": func(in *domain.DecisionTraceReviewContextInput) {
			in.BuildInput.Report.Reasons[0].Message = string([]byte{255})
		},
		"receipt": func(in *domain.DecisionTraceReviewContextInput) {
			seed := traceReviewRun(t, in.ReviewInput)
			r := traceReviewRecord(seed, 0, "COMPLETED")
			r.PlanHash = strings.Repeat("b", 64)
			in.ReviewInput.Records = []domain.DecisionTraceReviewRecord{r}
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			in := traceReviewContextFixture(t)
			change(&in)
			out, err := BuildDecisionTraceReviewContext(in)
			if err == nil || !reflect.DeepEqual(out, domain.DecisionTraceReviewContext{}) {
				t.Fatalf("tampered context accepted or partial output: %v", err)
			}
		})
	}
}

func TestDecisionTraceReviewContextAggregateInputBudget(t *testing.T) {
	in := traceReviewContextFixture(t)
	in.BuildInput.Report.Reasons = nil
	for i := 0; i < 180; i++ {
		in.BuildInput.Report.Reasons = append(in.BuildInput.Report.Reasons, domain.DecisionReason{RuleID: fmt.Sprintf("large-%03d", i), Message: strings.Repeat("x", 7000)})
	}
	graph, err := BuildDecisionTrace(in.BuildInput)
	if err != nil {
		t.Fatal(err)
	}
	in.ReviewInput.ImpactInput.Graph = graph
	var reportID string
	for _, node := range graph.Nodes {
		if node.Kind == domain.TraceInput && node.Label == "report" {
			reportID = node.ID
		}
	}
	if reportID == "" {
		t.Fatal("fixture requires report input")
	}
	in.ReviewInput.ImpactInput.Changes = []domain.DecisionTraceChange{{NodeID: reportID, Reason: "INPUT_CHANGED"}}
	seed := traceReviewRun(t, in.ReviewInput)
	for i := range seed.Tasks {
		r := traceReviewRecord(seed, i, "IN_PROGRESS")
		r.Note = strings.Repeat("n", 8192)
		in.ReviewInput.Records = append(in.ReviewInput.Records, r)
	}
	// Each operand may fit while their combined request cannot.
	if _, err := traceBoundedJSON(in.ReviewInput); err != nil {
		t.Fatalf("review operand must fit: %v", err)
	}
	if _, err := traceBoundedJSON(in.BuildInput); err != nil {
		t.Fatalf("build operand must fit: %v", err)
	}
	out, err := BuildDecisionTraceReviewContext(in)
	if err == nil || !reflect.DeepEqual(out, domain.DecisionTraceReviewContext{}) {
		t.Fatalf("aggregate input escaped guard: %v", err)
	}
}
