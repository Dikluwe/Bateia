/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_insights_test.md
 * @layer L2
 * @updated 2026-10-06
 */

package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func insightCLIRequest() domain.InsightInvestigationRequest {
	ref := domain.CommitRef{Owner: "illustrative", Repo: "library", Commit: strings.Repeat("a", 40)}
	content := "bounded queue with finite capacity"
	hash := sha256.Sum256([]byte(content))
	digest := "sha256:" + hex.EncodeToString(hash[:])
	item := domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Ref: ref, Path: "README.md", Content: content, ContentHash: digest}
	source := domain.EvidenceSource{Kind: item.Kind, Ref: ref, Path: item.Path, Line: 1, Snippet: content, RuleID: "bounded-queue-architecture-documentation", ContentHash: digest}
	return domain.InsightInvestigationRequest{SchemaVersion: "1", Question: "Como limitar admissão de trabalho?", Target: domain.InsightTarget{Name: "Bateia", Context: "Runtime local", Property: "trabalho pendente", DesiredEffect: "limitar sobrecarga"}, MechanismClasses: []string{"BOUNDED_QUEUE_BACKPRESSURE"}, Sources: []domain.InsightEvidenceInput{{Bundle: domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: []domain.CollectedEvidence{item}, Complete: true}, Assessments: []domain.CandidateAssessment{{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Status: domain.EvidencePartial, Sources: []domain.EvidenceSource{source}}}}}}}}
}

func TestInsightCLIProducesIndependentJSONAndProtectedOutput(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.json")
	output := filepath.Join(root, "output.json")
	data, _ := json.Marshal(insightCLIRequest())
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input}, &stdout); err != nil {
		t.Fatal(err)
	}
	var report domain.InsightInvestigationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || len(report.Hypotheses) != 1 || report.Hypotheses[0].Status != domain.HypothesisInconclusive {
		t.Fatalf("unexpected insight: %v %+v", err, report)
	}
	first := append([]byte(nil), stdout.Bytes()...)
	stdout.Reset()
	if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input, "--output", output}, &stdout); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(saved, first) || stdout.Len() != 0 {
		t.Fatal("file/stdout projection diverged")
	}
	before, _ := os.ReadFile(input)
	if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input, "--output", input}, &stdout); err == nil {
		t.Fatal("input overwritten")
	}
	alias := filepath.Join(root, "alias.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input, "--output", alias}, &stdout); err == nil {
		t.Fatal("input alias overwritten")
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(before, after) {
		t.Fatal("input mutated")
	}
}

type insightShortWriter struct{}

func (insightShortWriter) Write(data []byte) (int, error) { return len(data) - 1, nil }

func TestInsightCLIRejectsBadRequestsWithoutPublishing(t *testing.T) {
	data, _ := json.Marshal(insightCLIRequest())
	for name, payload := range map[string][]byte{
		"unknown-field":    append([]byte(`{"reuse_decision":"BUILD",`), data[1:]...),
		"extra-value":      append(append([]byte(nil), data...), []byte(` {}`)...),
		"duplicate-field":  append([]byte(`{"schema_version":"1",`), data[1:]...),
		"invalid-contract": []byte(`{"schema_version":"1"}`),
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			input := filepath.Join(root, "in.json")
			output := filepath.Join(root, "out.json")
			os.WriteFile(input, payload, 0600)
			os.WriteFile(output, []byte("preserve"), 0600)
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("invalid request accepted")
			}
			got, _ := os.ReadFile(output)
			if string(got) != "preserve" || stdout.Len() != 0 {
				t.Fatal("failed analysis published output")
			}
		})
	}
	input := filepath.Join(t.TempDir(), "valid.json")
	os.WriteFile(input, data, 0600)
	if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input}, insightShortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short writer accepted: %v", err)
	}
}
