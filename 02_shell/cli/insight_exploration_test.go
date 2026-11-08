/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_exploration_test.md
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
	"strconv"
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

func TestInsightExplorationCLIBooleanFlagsRejectDuplicatesBeforeReadingOrPublishing(t *testing.T) {
	for _, name := range []string{"algorithm-ideas", "review-alternatives"} {
		for _, flags := range [][]string{
			{"--" + name, "--" + name},
			{"--" + name + "=true", "--" + name + "=true"},
			{"--" + name + "=false", "--" + name + "=false"},
			{"--" + name, "-" + name + "=false"},
			{"-" + name + "=false", "--" + name},
		} {
			for _, existing := range []bool{false, true} {
				t.Run(strings.Join(flags, "/")+"/existing="+strconv.FormatBool(existing), func(t *testing.T) {
					root := t.TempDir()
					output := filepath.Join(root, "output.json")
					if existing {
						if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
							t.Fatal(err)
						}
					}
					args := []string{"insights", "explore", "--question", "queue", "--bundles", filepath.Join(root, "missing.json"), "--output", output}
					var stdout bytes.Buffer
					err := RunIdeaCommand(append(args, flags...), &stdout)
					if err == nil || !strings.Contains(err.Error(), "duplicada") || stdout.Len() != 0 {
						t.Fatalf("duplicate did not fail before reading/publication: %v", err)
					}
					data, readErr := os.ReadFile(output)
					if existing && (readErr != nil || string(data) != "preserve") || !existing && !os.IsNotExist(readErr) {
						t.Fatalf("rejected command changed output: %q, %v", data, readErr)
					}
				})
			}
		}
	}
}

func TestInsightExplorationCLIBooleanFormsPreserveModes(t *testing.T) {
	for _, name := range []string{"algorithm-ideas", "review-alternatives"} {
		for _, value := range []string{"absent", "implicit", "true", "false", "1", "0", "bad"} {
			t.Run(name+"/"+value, func(t *testing.T) {
				args := []string{"insights", "explore", "--question", "queue", "--familiar", "queue", "--familiar", "worker"}
				if value == "implicit" {
					args = append(args, "--"+name)
				} else if value != "absent" {
					args = append(args, "--"+name+"="+value)
				}
				var out bytes.Buffer
				err := RunIdeaCommand(args, &out)
				if value == "bad" {
					if err == nil || out.Len() != 0 {
						t.Fatal("invalid boolean accepted or published")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				var report domain.InsightExplorationReport
				if err := json.Unmarshal(out.Bytes(), &report); err != nil {
					t.Fatal(err)
				}
				want := 5
				if value == "implicit" || value == "true" || value == "1" {
					want = 7
					if name == "algorithm-ideas" {
						want = 10
					}
				}
				if len(report.Tasks) != want || len(report.FamiliarApproaches) != 2 || report.Status != "SEARCH_NEEDED" {
					t.Fatal("boolean/repeatable flag behavior changed")
				}
			})
		}
	}
}
