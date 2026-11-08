/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_tree_hash_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

type hashOracleFile struct{ path, content string }

// Independent wire-format oracle; no production tree walkers or hash helpers.
func framedTreeOracle(files ...hashOracleFile) string {
	var preimage bytes.Buffer
	preimage.WriteString("bateia-local-tree-v2\x00")
	for _, f := range files {
		preimage.WriteByte(1)
		_ = binary.Write(&preimage, binary.BigEndian, uint64(len(f.path)))
		preimage.WriteString(f.path)
		_ = binary.Write(&preimage, binary.BigEndian, uint64(len(f.content)))
		preimage.WriteString(f.content)
	}
	preimage.WriteByte(0)
	sum := sha256.Sum256(preimage.Bytes())
	return hex.EncodeToString(sum[:])
}

func putHashFile(t *testing.T, root, path, content string) {
	t.Helper()
	name := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestTreeHashV2SeparatesLegacyStructuralCollision(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	putHashFile(t, a, "a", "Xb\x00Y")
	putHashFile(t, b, "a", "X")
	putHashFile(t, b, "b", "Y")
	// Both inventories used exactly this preimage in v1, not a SHA collision.
	if hashBytes([]byte("a\x00Xb\x00Y")) != "fcdcf93cd1ef536fba84fd05e3ef0fa57922c0577e69cb9214110e8ca306c889" {
		t.Fatal("legacy oracle changed")
	}
	ah, err := treeHash(a)
	if err != nil {
		t.Fatal(err)
	}
	bh, err := treeHash(b)
	if err != nil {
		t.Fatal(err)
	}
	if ah == bh || ah != framedTreeOracle(hashOracleFile{"a", "Xb\x00Y"}) || bh != framedTreeOracle(hashOracleFile{"a", "X"}, hashOracleFile{"b", "Y"}) {
		t.Fatal("v2 framing failed")
	}
}

func TestTreeHashV2PreservesLexicalDFSAndExclusions(t *testing.T) {
	root := t.TempDir()
	putHashFile(t, root, "a/z", "Z")
	putHashFile(t, root, "a.txt", "T")
	for _, dir := range []string{".git", "target", "node_modules", "dossier"} {
		putHashFile(t, root, dir+"/ignored", "ignored")
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, "external")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("a.txt", filepath.Join(root, "internal")); err != nil {
		t.Fatal(err)
	}
	got, err := treeHashExcluding(root, filepath.Join(root, "dossier"))
	if err != nil || got != framedTreeOracle(hashOracleFile{"a/z", "Z"}, hashOracleFile{"a.txt", "T"}) {
		t.Fatalf("DFS/exclusion: %s %v", got, err)
	}
	if got == framedTreeOracle(hashOracleFile{"a.txt", "T"}, hashOracleFile{"a/z", "Z"}) {
		t.Fatal("global path ordering used")
	}
	gitFile := t.TempDir()
	putHashFile(t, gitFile, ".git", "gitdir: external")
	if h, e := treeHash(gitFile); e != nil || h != framedTreeOracle(hashOracleFile{".git", "gitdir: external"}) {
		t.Fatal("regular .git not hashed")
	}
	namedTarget := filepath.Join(t.TempDir(), "target")
	putHashFile(t, namedTarget, "x", "X")
	if h, e := treeHash(namedTarget); e != nil || h != framedTreeOracle(hashOracleFile{"x", "X"}) {
		t.Fatal("root target skipped")
	}
}

func TestTreeHashSharedEntryAndByteBudgets(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	putHashFile(t, a, "a", "abc")
	putHashFile(t, b, "b", "de")
	for _, test := range []struct {
		entries   int
		bytes     int64
		wantError bool
	}{{2, 5, false}, {1, 5, true}, {2, 4, true}} {
		budget := &treeHashBudget{EntriesLeft: test.entries, BytesLeft: test.bytes}
		if _, e := treeHashExcludingContext(context.Background(), a, "", budget); e != nil {
			t.Fatal(e)
		}
		got, e := treeHashExcludingContext(context.Background(), b, "", budget)
		if (e != nil) != test.wantError || (e != nil && got != "") {
			t.Fatalf("budget %+v: %s %v", test, got, e)
		}
	}
	empty := t.TempDir()
	putHashFile(t, empty, "empty", "")
	if _, e := treeHashExcludingContext(context.Background(), empty, "", &treeHashBudget{EntriesLeft: 1}); e != nil {
		t.Fatalf("zero bytes should fit: %v", e)
	}
	// Excluded names still consume entries, excluded contents do not.
	if h, e := treeHashExcludingContext(context.Background(), a, filepath.Join(a, "missing"), &treeHashBudget{BytesLeft: 3}); e == nil || h != "" {
		t.Fatal("entry budget ignored")
	}
}

func TestTreeHashProductionByteLimits(t *testing.T) {
	root := t.TempDir()
	for i, name := range []string{"a", "b", "c", "d"} {
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
		if i == 0 {
			if _, e = treeHash(root); e != nil {
				t.Fatalf("exact per-file cap: %v", e)
			}
		}
	}
	if _, e := treeHash(root); e != nil {
		t.Fatalf("exact aggregate cap: %v", e)
	}
	putHashFile(t, root, "e", "x")
	if h, e := treeHash(root); e == nil || h != "" {
		t.Fatal("aggregate +1 accepted")
	}
	if e := os.Truncate(filepath.Join(root, "a"), maxIdeaInputBytes+1); e != nil {
		t.Fatal(e)
	}
	if h, e := treeHash(root); e == nil || h != "" {
		t.Fatal("per-file +1 accepted")
	}
}

func TestTreeHashRootLimitsOrderAndAliases(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	putHashFile(t, a, "a", "X")
	putHashFile(t, b, "b", "Y")
	first, e := localCandidates(domain.IdeaSpec{}, []string{a, b})
	if e != nil {
		t.Fatal(e)
	}
	second, e := localCandidates(domain.IdeaSpec{}, []string{b, a})
	if e != nil || !reflect.DeepEqual(first, second) {
		t.Fatal("local ordering changed identity")
	}
	if c, e := localCandidates(domain.IdeaSpec{}, []string{a, filepath.Join(a, ".")}); e == nil || c != nil {
		t.Fatal("duplicate roots accepted")
	}
	if c, e := localCandidates(domain.IdeaSpec{}, make([]string, maxHashLocalRoots+1)); e == nil || c != nil {
		t.Fatal("root cap ignored")
	}
	for _, out := range []string{a, filepath.Dir(a)} {
		if h, e := treeHashExcluding(a, out); e == nil || h != "" {
			t.Fatal("self/ancestor output accepted")
		}
	}
	parent := t.TempDir()
	alias := filepath.Join(parent, "alias")
	if e := os.Symlink(a, alias); e != nil {
		t.Fatal(e)
	}
	if h, e := treeHash(alias); e == nil || h != "" {
		t.Fatal("alias root accepted")
	}
	if h, e := treeHashExcluding(a, filepath.Join(alias, "dossier")); e == nil || h != "" {
		t.Fatal("aliased output ancestor accepted")
	}
	if h, e := treeHash(strings.Repeat("a", maxHashTreePath+1)); e == nil || h != "" {
		t.Fatal("oversized path accepted")
	}
}

func TestTreeHashDepthLimit(t *testing.T) {
	root := t.TempDir()
	path := root
	for i := 0; i < maxHashTreeDepth; i++ {
		path = filepath.Join(path, "d")
	}
	putHashFile(t, path, "x", "X")
	if _, e := treeHash(root); e != nil {
		t.Fatalf("exact depth: %v", e)
	}
	putHashFile(t, path, "d/y", "Y")
	if h, e := treeHash(root); e == nil || h != "" {
		t.Fatal("depth +1 accepted")
	}
}

func TestTreeHashRefusesNonregularAndHardlinksWithoutBlocking(t *testing.T) {
	for _, kind := range []string{"fifo", "hardlink"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			name := filepath.Join(root, "x")
			if kind == "fifo" {
				if e := syscall.Mkfifo(name, 0600); e != nil {
					t.Fatal(e)
				}
			} else {
				outside := t.TempDir()
				putHashFile(t, outside, "x", "secret fixture")
				if e := os.Link(filepath.Join(outside, "x"), name); e != nil {
					t.Fatal(e)
				}
			}
			done := make(chan error, 1)
			go func() {
				h, e := treeHash(root)
				if h != "" {
					e = errors.New("partial hash returned")
				}
				done <- e
			}()
			select {
			case e := <-done:
				if e == nil {
					t.Fatal("unsafe file accepted")
				}
			case <-time.After(time.Second):
				t.Fatal("hash blocked")
			}
		})
	}
}

func TestTreeHashRetainsRootDescriptor(t *testing.T) {
	parent := t.TempDir()
	path := filepath.Join(parent, "repo")
	putHashFile(t, path, "x", "original")
	root, e := commandHost().OpenRoot(path)
	if e != nil {
		t.Fatal(e)
	}
	defer root.Close()
	if e = os.Rename(path, filepath.Join(parent, "old")); e != nil {
		t.Fatal(e)
	}
	putHashFile(t, path, "x", "replacement")
	got, e := hashRootedTree(context.Background(), root, "", newTreeHashBudget())
	if e != nil || got != framedTreeOracle(hashOracleFile{"x", "original"}) {
		t.Fatalf("root switched: %s %v", got, e)
	}
}

type cancelHashReader struct{ cancel context.CancelFunc }

func (r cancelHashReader) Read(p []byte) (int, error) { r.cancel(); p[0] = 'x'; return 1, nil }
func TestTreeHashCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if h, e := treeHashExcludingContext(ctx, "/nonexistent-fixture", "", newTreeHashBudget()); !errors.Is(e, context.Canceled) || h != "" {
		t.Fatal("cancelled context touched root")
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	r := hashContextReader{ctx: ctx, reader: cancelHashReader{cancel}}
	p := make([]byte, 4)
	if n, e := r.Read(p); n != 1 || e != nil {
		t.Fatal(n, e)
	}
	if n, e := r.Read(p); n != 0 || !errors.Is(e, context.Canceled) {
		t.Fatal("stream cancellation ignored")
	}
}

func TestTreeHashCheckpointVersionRequired(t *testing.T) {
	root := t.TempDir()
	putHashFile(t, root, "discovery-snapshot.json", "{}")
	putHashFile(t, root, "evidence/bundles.json", "[]")
	artifacts := []domain.IdeaRunArtifact{{Path: "discovery-snapshot.json", Hash: hashBytes([]byte("{}"))}, {Path: filepath.Join("evidence", "bundles.json"), Hash: hashBytes([]byte("[]"))}}
	previous := domain.IdeaRunManifest{SchemaVersion: "1", Status: domain.IdeaRunComplete, InputHashes: map[string]string{"evidence": "e", "evidence_rules": "r"}, RuleVersions: map[string]string{}, Stages: []domain.IdeaRunStage{{Outputs: artifacts}}}
	versions := map[string]string{"local_snapshot": localTreeHashVersion}
	for _, version := range []string{"", "tree-sha256-v1", localTreeHashVersion} {
		previous.RuleVersions["local_snapshot"] = version
		got := reusableEvidenceCheckpoint(root, previous, "e", "r", versions)
		if got != (version == localTreeHashVersion) {
			t.Fatalf("version %q reused=%v", version, got)
		}
	}
}

func TestTreeHashRecommendAdmissionDoesNotPublish(t *testing.T) {
	idea, policy, _ := recommendReplayFixture(t)
	for _, kind := range []string{"fifo", "timeout-zero", "timeout-expired", "roots", "aggregate"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			out := filepath.Join(t.TempDir(), "out")
			args := []string{"--idea", idea, "--policy", policy, "--mode", "local", "--strategy", "literal", "--local", root, "--output", out}
			switch kind {
			case "fifo":
				if e := syscall.Mkfifo(filepath.Join(root, "fifo"), 0600); e != nil {
					t.Fatal(e)
				}
			case "timeout-zero":
				args = append(args, "--timeout", "0s")
			case "timeout-expired":
				args = append(args, "--timeout", "1ns")
			case "roots":
				for i := 0; i < maxHashLocalRoots; i++ {
					args = append(args, "--local", root)
				}
			case "aggregate":
				for _, name := range []string{"a", "b", "c", "d", "e"} {
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
			var stdout bytes.Buffer
			if e := RunIdeaRecommendCommand(args, &stdout, nil, nil, nil); e == nil {
				t.Fatal("unsafe admission accepted")
			}
			if stdout.Len() != 0 {
				t.Fatal("admission emitted output")
			}
			if _, e := os.Lstat(out); !os.IsNotExist(e) {
				t.Fatal("admission published dossier")
			}
		})
	}
}
