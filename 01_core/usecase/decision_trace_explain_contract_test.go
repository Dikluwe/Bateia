/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_explain_contract_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func explainFixture(t *testing.T) domain.DecisionTraceExplainInput {
	t.Helper()
	source := domain.EvidenceSource{Kind: domain.EvidenceFile, Path: "unsafe-looking/~/data.json", Snippet: "declared observation", Strength: 0.2}
	component := domain.ReusableComponent{Name: "adapter", Kind: "module", Path: "adapter.go", Evidence: []domain.EvidenceSource{source}}
	candidate := domain.CandidateCapabilityStatus{Candidate: "candidate-a", Status: domain.EvidenceInconclusive, Evidence: []domain.EvidenceSource{source, source}}
	build := domain.DecisionTraceBuildInput{
		SchemaVersion: "1", Idea: &domain.IdeaSpec{Problem: "declarations only", Capabilities: []string{"alpha", "beta"}}, RuleVersions: map[string]string{"algorithm~/name": "v1"},
		Report: domain.ReuseDecisionReport{
			SchemaVersion: "1", Decision: domain.DecisionUnknown, Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1"},
			Matrix:            []domain.CapabilityComparison{{Capability: "alpha", Candidates: []domain.CandidateCapabilityStatus{candidate}}, {Capability: "beta", Candidates: []domain.CandidateCapabilityStatus{candidate}}},
			Reasons:           []domain.DecisionReason{{RuleID: "r1", Effect: "INFORMS", Message: "unverified", Evidence: []domain.EvidenceSource{source}}},
			Blockers:          []domain.DecisionReason{{RuleID: "r2", Effect: "BLOCK", Message: "unknown", Evidence: []domain.EvidenceSource{source}}},
			Components:        []domain.ReusableComponent{component, component},
			Composition:       &domain.CompositionReport{SchemaVersion: "1", Candidates: []domain.CompositionCandidate{{ID: "c1", Status: domain.CompositionProposed, Members: []domain.CompositionMember{{ID: "m1", Component: component}}, Evidence: []domain.EvidenceSource{source}}}},
			TechnicalEvidence: &domain.TechnicalEvidenceReview{SchemaVersion: "1", DecisionImpact: "NONE", Claims: []domain.TechnicalClaimReview{{Candidate: "candidate-a", Claim: "compatibility", Status: domain.TechnicalClaimInconclusive}}, Verifications: []domain.TechnicalEvidenceVerification{{FindingID: "f1", Candidate: "candidate-a", Status: domain.TechnicalUnverified, Source: &source, Message: "not authenticated", SourceAuthentication: "NOT_AUTHENTICATED"}}},
		},
	}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	return domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build}
}

func explainClone(t *testing.T, input domain.DecisionTraceExplainInput) domain.DecisionTraceExplainInput {
	t.Helper()
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var cloned domain.DecisionTraceExplainInput
	if err := json.Unmarshal(b, &cloned); err != nil {
		t.Fatal(err)
	}
	return cloned
}

// Resolve without decoding leaf records into maps, preserving exact Marshal
// object field order and number representations for the record fingerprint.
func explainResolve(t *testing.T, document []byte, pointer string) []byte {
	t.Helper()
	if pointer != "" {
		if !strings.HasPrefix(pointer, "/") {
			t.Fatalf("not an RFC 6901 pointer: %q", pointer)
		}
		for _, raw := range strings.Split(pointer[1:], "/") {
			token := strings.ReplaceAll(strings.ReplaceAll(raw, "~1", "/"), "~0", "~")
			trimmed := bytes.TrimSpace(document)
			if len(trimmed) == 0 {
				t.Fatal("empty pointer document")
			}
			if trimmed[0] == '[' {
				var children []json.RawMessage
				if err := json.Unmarshal(document, &children); err != nil {
					t.Fatal(err)
				}
				index, err := strconv.Atoi(token)
				if err != nil || index < 0 || index >= len(children) {
					t.Fatalf("invalid array index %q", token)
				}
				document = children[index]
			} else {
				var children map[string]json.RawMessage
				if err := json.Unmarshal(document, &children); err != nil {
					t.Fatal(err)
				}
				child, ok := children[token]
				if !ok {
					t.Fatalf("missing pointer key %q", token)
				}
				document = child
			}
		}
	}
	var compact bytes.Buffer
	if err := json.Compact(&compact, document); err != nil {
		t.Fatal(err)
	}
	return compact.Bytes()
}

func explainDigest(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func TestDecisionTraceExplainEveryLocatorBindsExactOriginalRecord(t *testing.T) {
	input := explainFixture(t)
	before, _ := json.Marshal(input)
	result, err := usecase.ExplainDecisionTrace(input)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(input)
	if !bytes.Equal(before, after) {
		t.Fatal("explanation mutated historical graph/build declarations")
	}
	if len(result.Entries) != len(input.Graph.Nodes) || result.Scope != domain.DecisionTraceScope || result.GraphHash != input.Graph.GraphHash {
		t.Fatal("missing graph binding or scope")
	}
	documents := map[string][]byte{}
	for name, value := range map[string]any{"report": input.BuildInput.Report, "idea": input.BuildInput.Idea, "rule_versions": input.BuildInput.RuleVersions} {
		documents[name], err = json.MarshalIndent(value, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
	}
	if result.ReportFingerprint != explainDigest(explainResolve(t, documents["report"], "")) {
		t.Fatal("incorrect report fingerprint")
	}
	nodes := map[string]domain.DecisionTraceNode{}
	for _, node := range input.Graph.Nodes {
		nodes[node.ID] = node
	}
	ids := []string{}
	sourceLocations := []string{}
	capabilityEntries := []domain.DecisionTraceExplanationEntry{}
	for _, entry := range result.Entries {
		ids = append(ids, entry.NodeID)
		node := nodes[entry.NodeID]
		if node.Kind != entry.Kind || node.Label != entry.Label || node.Fingerprint != entry.Fingerprint {
			t.Fatal("changed node declaration")
		}
		if len(entry.Locations) == 0 {
			t.Fatalf("entry has no locations: %s", entry.NodeID)
		}
		if len(entry.Statement) > 1024 || !utf8.ValidString(entry.Statement) {
			t.Fatal("unbounded display statement")
		}
		previous := ""
		for _, location := range entry.Locations {
			doc, ok := documents[location.Document]
			if !ok {
				t.Fatalf("unknown logical document %q", location.Document)
			}
			key := location.Document + ":" + location.Pointer
			if key <= previous {
				t.Fatal("locations are duplicate or unsorted")
			}
			previous = key
			if explainDigest(explainResolve(t, doc, location.Pointer)) != entry.Fingerprint {
				t.Fatalf("pointer %s mismatches record fingerprint", key)
			}
			if entry.Kind == domain.TraceSource {
				sourceLocations = append(sourceLocations, location.Pointer)
			}
		}
		if strings.HasPrefix(entry.Statement, "Declared capability ") {
			capabilityEntries = append(capabilityEntries, entry)
		}
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatal("entries are not canonically sorted")
	}
	wantedSourceLocations := []string{"/blockers/0/evidence/0", "/components/0/evidence/0", "/components/1/evidence/0", "/composition/candidates/0/evidence/0", "/composition/candidates/0/members/0/component/evidence/0", "/matrix/0/candidates/0/evidence/0", "/matrix/0/candidates/0/evidence/1", "/matrix/1/candidates/0/evidence/0", "/matrix/1/candidates/0/evidence/1", "/reasons/0/evidence/0", "/technical_evidence/verifications/0/source"}
	sort.Strings(wantedSourceLocations)
	if !reflect.DeepEqual(sourceLocations, wantedSourceLocations) {
		t.Fatalf("source occurrences missing: %v", sourceLocations)
	}
	if len(capabilityEntries) != 2 || capabilityEntries[0].Fingerprint != capabilityEntries[1].Fingerprint || capabilityEntries[0].NodeID == capabilityEntries[1].NodeID || reflect.DeepEqual(capabilityEntries[0].Locations, capabilityEntries[1].Locations) {
		t.Fatal("equal candidate records merged distinct capability identities")
	}
	again, err := usecase.ExplainDecisionTrace(input)
	if err != nil || !reflect.DeepEqual(result, again) {
		t.Fatal("nondeterministic explanation")
	}
	rebuilt, err := usecase.BuildDecisionTrace(input.BuildInput)
	if err != nil || !reflect.DeepEqual(rebuilt, input.Graph) || rebuilt.SchemaVersion != "decision-trace-v1" {
		t.Fatal("historical graph changed after explaining")
	}
}

func TestDecisionTraceExplainSelectionAndTampering(t *testing.T) {
	base := explainFixture(t)
	for name, mutate := range map[string]func(*domain.DecisionTraceExplainInput){
		"schema":             func(i *domain.DecisionTraceExplainInput) { i.SchemaVersion = "2" },
		"no frozen hash":     func(i *domain.DecisionTraceExplainInput) { i.Graph.GraphHash = "" },
		"hash":               func(i *domain.DecisionTraceExplainInput) { i.Graph.GraphHash = strings.Repeat("a", 64) },
		"graph declaration":  func(i *domain.DecisionTraceExplainInput) { i.Graph.Nodes[0].Label = "altered" },
		"report declaration": func(i *domain.DecisionTraceExplainInput) { i.BuildInput.Report.Reasons[0].Message = "altered" },
		"rule versions":      func(i *domain.DecisionTraceExplainInput) { i.BuildInput.RuleVersions["algorithm~/name"] = "v2" },
		"unknown selector":   func(i *domain.DecisionTraceExplainInput) { i.NodeIDs = []string{"missing"} },
		"duplicate selector": func(i *domain.DecisionTraceExplainInput) {
			i.NodeIDs = []string{i.Graph.Nodes[0].ID, i.Graph.Nodes[0].ID}
		},
		"selection budget": func(i *domain.DecisionTraceExplainInput) { i.NodeIDs = make([]string, 1001) },
		"string budget":    func(i *domain.DecisionTraceExplainInput) { i.BuildInput.Idea.Problem = strings.Repeat("x", 8193) },
		"UTF8":             func(i *domain.DecisionTraceExplainInput) { i.BuildInput.Idea.Problem = string([]byte{0xff}) },
	} {
		t.Run(name, func(t *testing.T) {
			input := explainClone(t, base)
			mutate(&input)
			result, err := usecase.ExplainDecisionTrace(input)
			if err == nil || !reflect.DeepEqual(result, domain.DecisionTraceExplanation{}) {
				t.Fatalf("invalid input returned success/partial result: %v", err)
			}
		})
	}
	base.NodeIDs = []string{base.Graph.Nodes[len(base.Graph.Nodes)-1].ID, base.Graph.Nodes[0].ID}
	result, err := usecase.ExplainDecisionTrace(base)
	if err != nil || len(result.Entries) != 2 || result.Entries[0].NodeID != base.Graph.Nodes[0].ID {
		t.Fatalf("selection failed: %v", err)
	}
}

func TestDecisionTraceExplainMultibyteStatementTruncatesButLocatorRemainsExact(t *testing.T) {
	input := explainFixture(t)
	input.BuildInput.Report.Reasons[0].Message = strings.Repeat("界", 700)
	graph, err := usecase.BuildDecisionTrace(input.BuildInput)
	if err != nil {
		t.Fatal(err)
	}
	input.Graph = graph
	result, err := usecase.ExplainDecisionTrace(input)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range result.Entries {
		if !strings.HasPrefix(entry.Statement, "Recorded reason ") {
			continue
		}
		if !entry.StatementTruncated || len(entry.Statement) > 1024 || len(entry.Statement) < 1022 || !utf8.ValidString(entry.Statement) {
			t.Fatal("multibyte statement truncation invalid")
		}
		report, _ := json.Marshal(input.BuildInput.Report)
		if len(explainResolve(t, report, entry.Locations[0].Pointer)) < 2100 {
			t.Fatal("locator points to truncated projection")
		}
		return
	}
	t.Fatal("reason projection missing")
}

func TestDecisionTraceExplainRejectsExpandedLocatorVector(t *testing.T) {
	// Individual input vectors fit their quota and declarations coalesce into
	// a tiny graph, but their union exceeds the output locator vector quota.
	source := domain.EvidenceSource{}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}}
	build.Report.Matrix = []domain.CapabilityComparison{{Capability: "alpha", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "same", Status: domain.EvidenceInconclusive, Evidence: make([]domain.EvidenceSource, 4000)}}}}
	build.Report.Reasons = []domain.DecisionReason{{RuleID: "extra", Evidence: []domain.EvidenceSource{source}}}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatalf("fixture build must fit input budget: %v", err)
	}
	result, err := usecase.ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build})
	if err == nil || !strings.Contains(err.Error(), "explanation output exceeds budget") || !reflect.DeepEqual(result, domain.DecisionTraceExplanation{}) {
		t.Fatalf("expanded locator output must fail closed, got %v", err)
	}
	build.Report.Matrix[0].Candidates[0].Evidence = build.Report.Matrix[0].Candidates[0].Evidence[:3999]
	graph, err = usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	result, err = usecase.ExplainDecisionTrace(domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build})
	if err != nil {
		t.Fatalf("bounded neighboring output rejected: %v", err)
	}
	for _, entry := range result.Entries {
		if entry.Kind == domain.TraceSource && len(entry.Locations) == 4000 {
			return
		}
	}
	t.Fatal("expected 4000 bound source occurrences")
}

func TestDecisionTraceExplainRejectsAggregateTypedInputBeforeBuilding(t *testing.T) {
	base := explainFixture(t)
	for _, name := range []string{"aggregate text", "vector"} {
		t.Run(name, func(t *testing.T) {
			input := explainClone(t, base)
			if name == "aggregate text" {
				input.BuildInput.RuleVersions = map[string]string{}
				for i := 0; i < 300; i++ {
					input.BuildInput.RuleVersions[strconv.Itoa(i)] = strings.Repeat("x", 8192)
				}
			} else {
				input.BuildInput.Report.Matrix = make([]domain.CapabilityComparison, 4001)
			}
			result, err := usecase.ExplainDecisionTrace(input)
			if err == nil || !strings.Contains(err.Error(), "budget") || !reflect.DeepEqual(result, domain.DecisionTraceExplanation{}) {
				t.Fatalf("oversized typed input not rejected under quota: %v", err)
			}
		})
	}
}
