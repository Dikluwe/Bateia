/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/decision_trace_explain_integration_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

// Identical payloads in different capability rows must not collapse their keys.
func TestExplainTraceBindsEachCapabilityToItsOwnOriginalLocation(t *testing.T) {
	source := domain.EvidenceSource{Kind: domain.EvidenceFile, Ref: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}, Path: "literal~/name.go", RuleID: "rule", Strength: .4}
	candidate := domain.CandidateCapabilityStatus{Candidate: "org/repo@abc", Status: domain.EvidenceInconclusive, Evidence: []domain.EvidenceSource{source}}
	report := domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{
		{Capability: "alpha", Candidates: []domain.CandidateCapabilityStatus{candidate}},
		{Capability: "beta", Candidates: []domain.CandidateCapabilityStatus{candidate}},
	}}
	build := domain.DecisionTraceBuildInput{SchemaVersion: "1", Report: report}
	graph, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	request := domain.DecisionTraceExplainInput{SchemaVersion: "1", Graph: graph, BuildInput: build}
	input := filepath.Join(t.TempDir(), "explain.json")
	if err := writeJSONAtomic(input, request); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"trace", "explain", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var index domain.DecisionTraceExplanation
	if err := json.Unmarshal(stdout.Bytes(), &index); err != nil {
		t.Fatal(err)
	}
	if index.GraphHash != graph.GraphHash || index.Scope != domain.DecisionTraceScope || len(index.Entries) != len(graph.Nodes) {
		t.Fatalf("unbound index %+v", index)
	}
	document, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	claimPointers := map[string]string{}
	sourceLocations := 0
	for _, entry := range index.Entries {
		if len(entry.Locations) == 0 {
			t.Fatalf("missing locator %s", entry.NodeID)
		}
		for _, location := range entry.Locations {
			if location.Document != "report" {
				t.Fatalf("unexpected document %s", location.Document)
			}
			raw := traceIntegrationResolvePointer(t, document, location.Pointer)
			var compact bytes.Buffer
			if err := json.Compact(&compact, raw); err != nil {
				t.Fatal(err)
			}
			digest := sha256.Sum256(compact.Bytes())
			if hex.EncodeToString(digest[:]) != entry.Fingerprint {
				t.Fatalf("locator %s does not bind %s", location.Pointer, entry.NodeID)
			}
			if entry.Kind == domain.TraceClaim {
				claimPointers[entry.Label] = location.Pointer
			}
			if entry.Kind == domain.TraceSource {
				sourceLocations++
			}
		}
	}
	if claimPointers["Capability: alpha / org/repo@abc"] != "/matrix/0/candidates/0" || claimPointers["Capability: beta / org/repo@abc"] != "/matrix/1/candidates/0" {
		t.Fatalf("same payload confused distinct claims: %+v", claimPointers)
	}
	if sourceLocations != 2 {
		t.Fatalf("shared source lost occurrence: %d", sourceLocations)
	}
	after, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("explain mutated input")
	}
	// Asking for an index must leave the historical graph unchanged and usable.
	again, err := usecase.BuildDecisionTrace(build)
	if err != nil {
		t.Fatal(err)
	}
	originalBytes, _ := json.Marshal(graph)
	againBytes, _ := json.Marshal(again)
	if !bytes.Equal(originalBytes, againBytes) {
		t.Fatal("explain changed historical graph")
	}
	plan, err := usecase.PlanDecisionTraceImpact(domain.DecisionTraceImpactInput{SchemaVersion: "1", Graph: graph, Changes: []domain.DecisionTraceChange{}})
	if err != nil || plan.DecisionAfter != graph.Decision {
		t.Fatalf("index broke old impact: %v", err)
	}
}

func traceIntegrationResolvePointer(t *testing.T, document json.RawMessage, pointer string) json.RawMessage {
	t.Helper()
	if pointer == "" {
		return document
	}
	if !strings.HasPrefix(pointer, "/") {
		t.Fatalf("invalid pointer %q", pointer)
	}
	for _, part := range strings.Split(pointer[1:], "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		var object map[string]json.RawMessage
		if err := json.Unmarshal(document, &object); err == nil && object != nil {
			value, ok := object[part]
			if !ok {
				t.Fatalf("unresolved pointer %s", pointer)
			}
			document = value
			continue
		}
		var array []json.RawMessage
		if err := json.Unmarshal(document, &array); err != nil {
			t.Fatalf("non-container in pointer %s", pointer)
		}
		index := 0
		if part == "" {
			t.Fatal("empty array index")
		}
		for _, digit := range part {
			if digit < '0' || digit > '9' {
				t.Fatalf("invalid index %s", part)
			}
			index = index*10 + int(digit-'0')
		}
		if index < 0 || index >= len(array) {
			t.Fatalf("out of range pointer %s", pointer)
		}
		document = array[index]
	}
	return document
}
