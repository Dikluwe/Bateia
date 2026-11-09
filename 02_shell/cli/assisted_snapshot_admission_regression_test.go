/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/assisted_snapshot_admission_regression_test.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestAssistedPrepareRejectsUnsealedAndTamperedInputWithoutOutput(t *testing.T) {
	for _, tamper := range []bool{false, true} {
		dir := t.TempDir()
		input, output := filepath.Join(dir, "snapshot.json"), filepath.Join(dir, "blind.json")
		request := domain.ObservableWebRequest{RequestID: "r", QueryID: "q", Query: "literal", MaxResults: 1}
		response := usecase.SealWebResponse(request, domain.RawWebResponse{RequestID: "r", QueryID: "q", Provider: "fixture", RetrievedAt: "2026-11-08T00:00:00Z", Status: "SUCCESS", Results: []domain.RawWebResult{{Rank: 1, URL: "https://github.com/acme/tool", Title: "literal", Snippet: "evidence"}}})
		s := domain.WebSearchSnapshot{SchemaVersion: "1", RunID: "fixture", Responses: []domain.RawWebResponse{response}}
		if tamper {
			s.SnapshotHash = usecase.WebSnapshotHash(s)
			s.Responses[0].Results[0].Title = "tampered"
		}
		writeStep0021JSON(t, input, s)
		err := RunAssistedJudgmentCommand([]string{"prepare", "--snapshot", input, "--output", output}, &bytes.Buffer{})
		if err == nil {
			t.Error("invalid snapshot published")
		}
		if _, e := os.Stat(output); !os.IsNotExist(e) {
			t.Error("invalid input created output", e)
		}
	}
}
