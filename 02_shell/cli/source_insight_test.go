/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/source_insight_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sourceInsightCLIRequest() domain.InsightInvestigationRequest {
	r := insightCLIRequest()
	r.SchemaVersion = "2"
	item := r.Sources[0].Bundle.Items[0]
	r.Sources[0].Observations = []domain.SourceInsightObservation{{MechanismClass: r.MechanismClasses[0], Kind: domain.InsightPrecondition, Statement: "A capacidade finita é declarada nesta documentação", References: []domain.InsightObservationReference{{Path: item.Path, ContentHash: item.ContentHash, StartLine: 1, EndLine: 1}}}}
	return r
}

func TestSourceInsightCLIProjectsSpecificObservationsAndProtectsInput(t *testing.T) {
	root := t.TempDir()
	input, output := filepath.Join(root, "in.json"), filepath.Join(root, "out.json")
	data, err := json.Marshal(sourceInsightCLIRequest())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	args := []string{"insights", "analyze", "--input", input}
	if err := RunIdeaCommand(args, &stdout); err != nil {
		t.Fatal(err)
	}
	var report domain.SourceInsightInvestigationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "insight-source-analysis-v2" || len(report.SourceAnalyses) != 1 || len(report.Hypotheses) != 1 {
		t.Fatalf("v2 not dispatched: %+v", report)
	}
	analysis := report.SourceAnalyses[0]
	if analysis.Observations[0].References[0].Excerpt != "bounded queue with finite capacity" || analysis.Observations[0].Origin != "DECLARED_NOT_REASSESSED" || analysis.CatalogDefinition.Origin != "CATALOG_TEMPLATE" {
		t.Fatal("source observation or catalog origin lost")
	}
	if len(report.Hypotheses[0].Preconditions) != 1 || report.Hypotheses[0].Preconditions[0] == "go" || report.Hypotheses[0].TargetState != domain.TargetUnknown || report.Hypotheses[0].Status != domain.HypothesisInconclusive {
		t.Fatal("source condition confused with scope or target effect")
	}
	first := append([]byte(nil), stdout.Bytes()...)
	stdout.Reset()
	if err := RunIdeaCommand(append(args, "--output", output), &stdout); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(first, saved) || stdout.Len() != 0 {
		t.Fatal("file/stdout projection differs")
	}
	alias := filepath.Join(root, "alias.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{input, alias} {
		if err := RunIdeaCommand(append(args, "--output", path), &stdout); err == nil {
			t.Fatal("input alias accepted")
		}
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(data, after) {
		t.Fatal("input overwritten")
	}
	if err := RunIdeaCommand(args, insightShortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short writer accepted: %v", err)
	}
}

func TestSourceInsightCLIRejectsInvalidAndForgedObservationsWithoutPublishing(t *testing.T) {
	valid := sourceInsightCLIRequest()
	data, _ := json.Marshal(valid)
	payloads := map[string][]byte{
		"extra-value":             append(append([]byte(nil), data...), []byte(" {}")...),
		"duplicate-schema":        append([]byte(`{"schema_version":"2",`), data[1:]...),
		"forged-origin":           bytes.Replace(data, []byte(`"statement":`), []byte(`"origin":"VERIFIED","statement":`), 1),
		"unknown-reference-field": bytes.Replace(data, []byte(`"start_line":`), []byte(`"verified":true,"start_line":`), 1),
	}
	for _, change := range []struct {
		name  string
		apply func(*domain.InsightInvestigationRequest)
	}{
		{"bad-hash", func(r *domain.InsightInvestigationRequest) {
			r.Sources[0].Observations[0].References[0].ContentHash = strings.Repeat("0", 64)
		}},
		{"bad-line", func(r *domain.InsightInvestigationRequest) { r.Sources[0].Observations[0].References[0].EndLine = 2 }},
		{"legacy-with-observations", func(r *domain.InsightInvestigationRequest) { r.SchemaVersion = "1" }},
		{"unknown-schema", func(r *domain.InsightInvestigationRequest) { r.SchemaVersion = "3" }},
	} {
		r := sourceInsightCLIRequest()
		change.apply(&r)
		payloads[change.name], _ = json.Marshal(r)
	}
	for name, payload := range payloads {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			input, output := filepath.Join(root, "in.json"), filepath.Join(root, "out.json")
			if err := os.WriteFile(input, payload, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := RunIdeaCommand([]string{"insights", "analyze", "--input", input, "--output", output}, &stdout); err == nil {
				t.Fatal("invalid or promoted observation accepted")
			}
			saved, _ := os.ReadFile(output)
			if string(saved) != "preserve" || stdout.Len() != 0 {
				t.Fatal("failed request published output")
			}
		})
	}
}

func TestSourceInsightCLIKeepsSecretProtection(t *testing.T) {
	const secret = "test-0134-not-a-real-credential"
	t.Setenv("GITHUB_TOKEN", secret)
	r := sourceInsightCLIRequest()
	r.Sources[0].Observations[0].Statement = "Não publicar " + secret
	data, _ := json.Marshal(r)
	input := filepath.Join(t.TempDir(), "in.json")
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"insights", "analyze", "--input", input}, &stdout)
	if err == nil || stdout.Len() != 0 || strings.Contains(err.Error(), secret) {
		t.Fatal("secret exposed or published")
	}
}
