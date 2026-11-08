/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/algorithm_insight_test.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
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
)

func algorithmCLIRequest() domain.AlgorithmInsightRequest {
	mk := func(owner, text string) domain.EvidenceBundle {
		ref := domain.CommitRef{Owner: owner, Repo: "scheduler", Commit: strings.Repeat("a", 40)}
		h := sha256.Sum256([]byte(text))
		return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceFile, Ref: ref, Path: "logic.go", Content: text, ContentHash: "sha256:" + hex.EncodeToString(h[:])}}}
	}
	source, target := mk("reference", "ready.pop()"), mk("target", "scanAll()")
	ev := func(b domain.EvidenceBundle) domain.AlgorithmIdeaEvidence {
		return domain.AlgorithmIdeaEvidence{Source: domain.ImmutableCandidateRef{Repository: b.Repository.Owner + "/" + b.Repository.Repo, Commit: b.Repository.Commit}, LocatedInsightReference: domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: b.Items[0].Path, ContentHash: b.Items[0].ContentHash, StartLine: 1, EndLine: 1}, Excerpt: b.Items[0].Content}}
	}
	return domain.AlgorithmInsightRequest{SchemaVersion: "1", Question: "How to schedule ready work?", Sources: []domain.EvidenceBundle{source}, Target: target, Ideas: []domain.AlgorithmIdea{{Title: "Ready collection", Principle: "Maintain ready candidates.", Steps: []string{"Update readiness.", "Choose from maintained collection."}, Evidence: []domain.AlgorithmIdeaEvidence{ev(source)}, ApplicationConditions: []string{"Readiness changes are observable."}, FailureConditions: []string{"Missed updates leave stale entries."}, Decision: domain.AlgorithmDecision{Current: "Scan every entry.", Proposed: "Maintain candidates.", Benefit: "May reduce selection work.", Cost: "Bookkeeping."}, TargetRelation: "DIFFERENT_APPROACH", TargetEvidence: []domain.AlgorithmIdeaEvidence{ev(target)}, ValidationObservable: "Selections and update work.", RejectIf: "Ready work is lost."}}}
}
func TestAlgorithmInsightCLIValidatesBeforePublishing(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "request.json"), filepath.Join(dir, "report.json")
	request := algorithmCLIRequest()
	data, _ := json.Marshal(request)
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"insights", "formulate", "--input", input}
	var out bytes.Buffer
	if err := RunIdeaCommand(args, &out); err != nil {
		t.Fatal(err)
	}
	var report domain.AlgorithmInsightReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || len(report.Ideas) != 1 || report.EffectState != "UNVERIFIED" {
		t.Fatal("formulation missing", err)
	}
	expected := append([]byte{}, out.Bytes()...)
	out.Reset()
	if err := RunIdeaCommand(append(args, "--output", output), &out); err != nil {
		t.Fatal(err)
	}
	saved, _ := os.ReadFile(output)
	if !bytes.Equal(saved, expected) || out.Len() != 0 {
		t.Fatal("stdout/file differ")
	}
	if err := RunIdeaCommand(args, insightShortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatal("short write ignored", err)
	}
	if err := RunIdeaCommand(append(args, "--output", input), &out); err == nil {
		t.Fatal("input overwritten")
	}
	for _, bad := range [][]byte{bytes.Replace(data, []byte(`"schema_version":"1"`), []byte(`"extra":true,"schema_version":"1"`), 1), bytes.Replace(data, []byte(`"excerpt":"ready.pop()"`), []byte(`"excerpt":"invented"`), 1)} {
		if err := os.WriteFile(input, bad, 0600); err != nil {
			t.Fatal(err)
		}
		out.Reset()
		if err := RunIdeaCommand(append(args, "--output", output), &out); err == nil || out.Len() != 0 {
			t.Fatal("invalid input published")
		}
		after, _ := os.ReadFile(output)
		if !bytes.Equal(after, expected) {
			t.Fatal("failed validation replaced output")
		}
	}
}
