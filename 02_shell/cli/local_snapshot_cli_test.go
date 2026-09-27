/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/local_snapshot_cli_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalCaptureCLIHashesConsultedMetadataAndCleans(t *testing.T) {
	idea, policy, _ := recommendReplayFixture(t)
	source := t.TempDir()
	storage := t.TempDir()
	t.Setenv("TMPDIR", storage)
	putHashFile(t, source, "README.md", "Dependency graph fixture")
	putHashFile(t, source, "go.mod", "module fixture\n\ngo 1.25\n")
	putHashFile(t, source, ".git/config", "[remote \"origin\"]\nurl = https://github.com/owner/old\n")
	output := filepath.Join(t.TempDir(), "out")
	run := func(resume bool) domain.IdeaRunManifest {
		t.Helper()
		args := []string{"--idea", idea, "--policy", policy, "--mode", "local", "--strategy", "literal", "--local", source, "--target-repo", "owner/old", "--output", output}
		if resume {
			args = append(args, "--resume")
		}
		if e := RunIdeaRecommendCommand(args, &bytes.Buffer{}, nil, nil, nil); e != nil {
			t.Fatal(e)
		}
		var manifest domain.IdeaRunManifest
		if e := readJSONFile(filepath.Join(output, "run.json"), &manifest); e != nil {
			t.Fatal(e)
		}
		entries, e := os.ReadDir(storage)
		if e != nil {
			t.Fatal(e)
		}
		// output's parent is allocated under TMPDIR by t.TempDir: only owned capture
		// directory names, not unrelated fixture directories, must be absent.
		for _, entry := range entries {
			if len(entry.Name()) >= len("bateia-local-snapshot-") && entry.Name()[:len("bateia-local-snapshot-")] == "bateia-local-snapshot-" {
				t.Fatal("successful CLI leaked capture")
			}
		}
		return manifest
	}
	first := run(false)
	if first.InputHashes["local_metadata"] == "" || first.RuleVersions["local_collection"] != localCollectionVersion {
		t.Fatal("capture identity not recorded")
	}
	putHashFile(t, source, ".git/config", "[remote \"origin\"]\nurl = https://github.com/owner/new\n")
	second := run(true)
	if first.InputHashes["evidence"] != second.InputHashes["evidence"] {
		t.Fatal("advisory Git hint changed file identity")
	}
	if first.InputHashes["local_metadata"] == second.InputHashes["local_metadata"] || first.RunID == second.RunID {
		t.Fatal("consulted metadata did not affect run identity")
	}
	found := false
	for _, stage := range second.Stages {
		if stage.Name == "collect" {
			found = true
			if stage.Counts["checkpoint_reused"] != 1 {
				t.Fatal("unchanged captured file evidence not reused")
			}
		}
	}
	if !found {
		t.Fatal("collection stage missing")
	}
}

func TestLocalCaptureCLICleansWhenOutputAdmissionFails(t *testing.T) {
	idea, policy, _ := recommendReplayFixture(t)
	source := t.TempDir()
	storage := t.TempDir()
	output := t.TempDir()
	putHashFile(t, source, "README.md", "fixture")
	putHashFile(t, output, "existing.txt", "preserve user fixture")
	t.Setenv("TMPDIR", storage)
	var stdout bytes.Buffer
	err := RunIdeaRecommendCommand([]string{"--idea", idea, "--policy", policy, "--mode", "local", "--strategy", "literal", "--local", source, "--output", output}, &stdout, nil, nil, nil)
	if err == nil || stdout.Len() != 0 {
		t.Fatal("nonempty output admitted")
	}
	entries, e := os.ReadDir(storage)
	if e != nil || len(entries) != 0 {
		t.Fatal("output refusal leaked capture", e)
	}
	data, e := os.ReadFile(filepath.Join(output, "existing.txt"))
	if e != nil || string(data) != "preserve user fixture" {
		t.Fatal("existing output changed")
	}
	if _, e := os.Lstat(filepath.Join(output, "run.json")); !os.IsNotExist(e) {
		t.Fatal("output refusal created manifest")
	}
}
