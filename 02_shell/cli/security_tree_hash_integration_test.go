/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_tree_hash_integration_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTreeHashProductionAggregateAcrossRoots(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	for _, root := range []string{a, b} {
		for _, name := range []string{"a", "b"} {
			f, e := os.Create(filepath.Join(root, name))
			if e != nil {
				t.Fatal(e)
			}
			if e = f.Truncate(maxIdeaInputBytes); e != nil {
				t.Fatal(e)
			}
			if e = f.Close(); e != nil {
				t.Fatal(e)
			}
		}
	}
	if _, e := localCandidates(domain.IdeaSpec{}, []string{a, b}); e != nil {
		t.Fatalf("exact aggregate: %v", e)
	}
	putHashFile(t, b, "c", "x")
	if c, e := localCandidates(domain.IdeaSpec{}, []string{a, b}); e == nil || c != nil {
		t.Fatal("aggregate +1 across roots accepted")
	}
	roots := make([]string, maxHashLocalRoots)
	for i := range roots {
		roots[i] = t.TempDir()
	}
	if c, e := localCandidates(domain.IdeaSpec{}, roots); e != nil || len(c) != maxHashLocalRoots {
		t.Fatal("exact root limit rejected", e)
	}
}

type hashCheckpointContent struct {
	fixtureContent
	calls int
}

func (c *hashCheckpointContent) FetchTree(ctx context.Context, repo domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	c.calls++
	return c.fixtureContent.FetchTree(ctx, repo)
}

func TestTreeHashLocalCheckpointVersionEndToEnd(t *testing.T) {
	// 0097: production always reads the private capture, never the injected local
	// fake. Assert actual collected bytes and observable checkpoint reuse instead.
	idea, policy, _ := recommendReplayFixture(t)
	root := t.TempDir()
	putHashFile(t, root, "x.go", "package fixture")
	putHashFile(t, root, "README.md", "Dependency graph fixture; captured documentation")
	putHashFile(t, root, "go.mod", "module fixture\n\ngo 1.25\n")
	output := filepath.Join(t.TempDir(), "out")
	args := []string{"--idea", idea, "--policy", policy, "--mode", "local", "--strategy", "literal", "--local", root, "--output", output}
	content := &hashCheckpointContent{}
	run := func(resume bool) {
		t.Helper()
		invocation := append([]string(nil), args...)
		if resume {
			invocation = append(invocation, "--resume")
		}
		if e := RunIdeaRecommendCommand(invocation, &bytes.Buffer{}, nil, nil, content); e != nil {
			t.Fatal(e)
		}
	}
	run(false)
	var manifest domain.IdeaRunManifest
	if e := readJSONFile(filepath.Join(output, "run.json"), &manifest); e != nil {
		t.Fatal(e)
	}
	if manifest.RuleVersions["local_snapshot"] != localTreeHashVersion {
		t.Fatal("snapshot version missing from run")
	}
	if manifest.RuleVersions["local_collection"] != localCollectionVersion || content.calls != 0 {
		t.Fatal("collection was not frozen or consulted original provider")
	}
	var bundles []domain.EvidenceBundle
	if err := readJSONFile(filepath.Join(output, "evidence", "bundles.json"), &bundles); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, bundle := range bundles {
		for _, item := range bundle.Items {
			if item.Path == "README.md" && item.Content == "Dependency graph fixture; captured documentation" && item.ContentHash == hashBytes([]byte(item.Content)) {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("positive control did not collect captured README bytes")
	}
	checkpointCount := func() int64 {
		for _, stage := range manifest.Stages {
			if stage.Name == "collect" {
				return stage.Counts["checkpoint_reused"]
			}
		}
		t.Fatal("collection stage missing")
		return -1
	}
	if checkpointCount() != 0 {
		t.Fatal("first run marked as reused")
	}
	run(true)
	if e := readJSONFile(filepath.Join(output, "run.json"), &manifest); e != nil {
		t.Fatal(e)
	}
	if checkpointCount() != 1 {
		t.Fatal("current frozen checkpoint not reused")
	}
	delete(manifest.RuleVersions, "local_collection")
	if e := writeJSONAtomic(filepath.Join(output, "run.json"), manifest); e != nil {
		t.Fatal(e)
	}
	run(true)
	if e := readJSONFile(filepath.Join(output, "run.json"), &manifest); e != nil {
		t.Fatal(e)
	}
	if checkpointCount() != 0 {
		t.Fatal("legacy local checkpoint reused")
	}
}
