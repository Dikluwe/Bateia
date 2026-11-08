/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_exploration_test.md
 * @layer L2
 * @updated 2026-10-28
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

func explorationCLIBundles(content string) []domain.EvidenceBundle {
	bundle := insightCLIRequest().Sources[0].Bundle
	bundle.Items[0].Content = content
	digest := sha256.Sum256([]byte(content))
	bundle.Items[0].ContentHash = "sha256:" + hex.EncodeToString(digest[:])
	return []domain.EvidenceBundle{bundle}
}

func TestInsightExplorationCLIQuestionAndCollectedBundlesWithoutCuratingJSON(t *testing.T) {
	root := t.TempDir()
	input, output := filepath.Join(root, "bundles.json"), filepath.Join(root, "out.json")
	data, _ := json.Marshal(explorationCLIBundles("queue cancellation example\n"))
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"insights", "explore", "--question", "How to rethink queue cancellation?", "--bundles", input, "--familiar", "bounded queue", "--familiar", "thread pool"}
	var stdout bytes.Buffer
	if err := RunIdeaCommand(args, &stdout); err != nil {
		t.Fatal(err)
	}
	var report domain.InsightExplorationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Cards) != 1 || len(report.Tasks) != 5 || report.NoveltyState != "NOT_ASSESSED" || len(report.FamiliarApproaches) != 2 || report.Cards[0].Reference.Excerpt != "queue cancellation example" {
		t.Fatal("context not prepared directly from collection")
	}
	expected := append([]byte(nil), stdout.Bytes()...)
	stdout.Reset()
	if err := RunIdeaCommand(append(args, "--output", output), &stdout); err != nil {
		t.Fatal(err)
	}
	saved, _ := os.ReadFile(output)
	if !bytes.Equal(saved, expected) || stdout.Len() != 0 {
		t.Fatal("file/stdout differs")
	}
	if err := RunIdeaCommand(args, insightShortWriter{}); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short write ignored: %v", err)
	}
	alias := filepath.Join(root, "alias.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand(append(args, "--output", alias), &stdout); err == nil {
		t.Fatal("input alias overwritten")
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(after, data) {
		t.Fatal("input mutated")
	}
	stdout.Reset()
	if err := RunIdeaCommand([]string{"insights", "explore", "--question", "--a question starting with dashes"}, &stdout); err != nil {
		t.Fatalf("flag value confused with flag: %v", err)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("SEARCH_NEEDED")) {
		t.Fatal("missing bundles promoted to context")
	}
}

func TestInsightExplorationCLIInvalidFlagsAndJSONNeverPublish(t *testing.T) {
	for _, args := range [][]string{
		{}, {"--question", ""}, {"--question", "question", "extra"}, {"--question", "question", "--unknown", "yes"},
		{"--question", "first", "-question", "second"}, {"--question=first", "--question=second"},
		{"--question", "question", "--max-passages", "17"}, {"--question", "question", "--max-passages", "bad"},
	} {
		var stdout bytes.Buffer
		if err := RunIdeaCommand(append([]string{"insights", "explore"}, args...), &stdout); err == nil || stdout.Len() != 0 {
			t.Fatal("invalid flags accepted")
		}
	}
	data, _ := json.Marshal(explorationCLIBundles("queue example"))
	for _, payload := range [][]byte{
		[]byte(strings.Replace(string(data), `"schema_version":"1"`, `"injected":true,"schema_version":"1"`, 1)),
		[]byte(strings.Replace(string(data), `"schema_version":"1"`, `"schema_version":"1","schema_version":"1"`, 1)),
		append(append([]byte(nil), data...), []byte(` {}`)...),
		[]byte(strings.Replace(string(data), "queue example", "tampered", 1)),
	} {
		root := t.TempDir()
		input, output := filepath.Join(root, "in.json"), filepath.Join(root, "out.json")
		if err := os.WriteFile(input, payload, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
			t.Fatal(err)
		}
		var stdout bytes.Buffer
		if err := RunIdeaCommand([]string{"insights", "explore", "--question", "queue", "--bundles", input, "--output", output}, &stdout); err == nil {
			t.Fatal("invalid corpus accepted")
		}
		saved, _ := os.ReadFile(output)
		if string(saved) != "preserve" || stdout.Len() != 0 {
			t.Fatal("failed exploration published")
		}
	}
}

func TestInsightExplorationCLIBlocksSecretInAutomaticExcerpt(t *testing.T) {
	const secret = "test-0136-not-a-real-credential"
	t.Setenv("GITHUB_TOKEN", secret)
	data, _ := json.Marshal(explorationCLIBundles("queue " + secret))
	input := filepath.Join(t.TempDir(), "in.json")
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err := RunIdeaCommand([]string{"insights", "explore", "--question", "queue", "--bundles", input}, &stdout)
	if err == nil || stdout.Len() != 0 || strings.Contains(err.Error(), secret) {
		t.Fatal("automatic reading exposed secret")
	}
}
