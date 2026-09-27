/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_advanced_input_test.md
 * @layer L2
 * @updated 2026-09-26
 */

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestAdvancedIdeaCommandsRejectDuplicateJSONBeforeOutput(t *testing.T) {
	for _, payload := range []string{
		`{"schema_version":"1","schema_version":"2"}`,
		`{"schema_version":"1","SCHEMA_VERSION":"2"}`,
		`{"schema_version":"1","schema_\u0076ersion":"2"}`,
		`{"item":{"name":"first","Name":"second"}}`,
	} {
		t.Run(payload, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input.json")
			if err := os.WriteFile(input, []byte(payload), 0600); err != nil {
				t.Fatal(err)
			}
			commands := []struct {
				name string
				run  func(string, *bytes.Buffer) error
			}{
				{"plan", func(output string, out *bytes.Buffer) error {
					return RunIdeaCommand([]string{"plan", "--input", input, "--output", output}, out)
				}},
				{"screen", func(output string, out *bytes.Buffer) error {
					return RunIdeaCommand([]string{"screen", "--snapshot", input, "--output", output}, out)
				}},
				{"evaluate", func(output string, out *bytes.Buffer) error {
					return RunIdeaCommand([]string{"evaluate", "--corpus", input, "--output", output}, out)
				}},
				{"catalog", func(output string, out *bytes.Buffer) error {
					return runIdeaCatalogScreen([]string{"--snapshot", input, "--output", output}, out)
				}},
				{"radar", func(output string, out *bytes.Buffer) error {
					return RunIdeaRadarExportCommand([]string{"export", "--input", input, "--output", output}, out)
				}},
				{"evidence_collect", func(output string, out *bytes.Buffer) error {
					return RunIdeaEvidenceCollectCommand([]string{"--plan", input, "--output", output}, out, nil)
				}},
				{"discover", func(output string, out *bytes.Buffer) error {
					return RunIdeaDiscoverCommand([]string{"--plan", input, "--output", output}, out, &bytes.Buffer{}, nil, nil)
				}},
				{"recommend", func(output string, out *bytes.Buffer) error {
					return RunIdeaRecommendCommand([]string{"--idea", input, "--policy", input, "--replay", dir, "--output", output}, out, nil, nil, nil)
				}},
			}
			for _, cmd := range commands {
				t.Run(cmd.name, func(t *testing.T) {
					output := filepath.Join(dir, cmd.name+"-out")
					var out bytes.Buffer
					err := cmd.run(output, &out)
					if err == nil || !strings.Contains(err.Error(), "duplicad") {
						t.Fatalf("duplicate guard missing: %v", err)
					}
					if out.Len() != 0 {
						t.Fatal("rejected input emitted stdout")
					}
					if _, err := os.Lstat(output); !os.IsNotExist(err) {
						t.Fatal("rejected input created output")
					}
				})
			}
		})
	}
}

func TestRecommendReadersRejectOversizedAndFIFOInputs(t *testing.T) {
	dir := t.TempDir()
	large := filepath.Join(dir, "oversized")
	f, err := os.Create(large)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxIdeaInputBytes + 1); err != nil {
		t.Fatal(err)
	}
	f.Close()
	fifo := filepath.Join(dir, "fifo")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{large, fifo} {
		for _, run := range []func() error{
			func() error { _, err := fileHash(input); return err },
			func() error { _, _, err := readOptInScreeningModel(input); return err },
			func() error { var v map[string]any; return readJSONFile(input, &v) },
			func() error { return RunIdeaRadarExportCommand([]string{"export", "--input", input}, &bytes.Buffer{}) },
		} {
			done := make(chan error, 1)
			go func() { done <- run() }()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("unsafe input accepted")
				}
			case <-time.After(time.Second):
				t.Fatal("input reader blocked on FIFO")
			}
		}
	}
}

func TestRecommendScreeningModelRejectsDuplicateJSON(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "model.json")
	if err := os.WriteFile(input, []byte(`{"schema_version":"1","SCHEMA_VERSION":"1"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readOptInScreeningModel(input); err == nil {
		t.Fatal("duplicate screening model accepted")
	}
	output := filepath.Join(dir, "dossier")
	err := RunIdeaRecommendCommand([]string{"--idea", input, "--policy", input, "--replay", dir, "--screening-model", input, "--output", output}, &bytes.Buffer{}, nil, nil, nil)
	if err == nil {
		t.Fatal("command accepted ambiguous model")
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Fatal("bad model created dossier")
	}
}

// 0096 supersedes the ambiguous legacy tree hash; file-byte hashes stay stable.
func TestRecommendFileHashStableAndTreeHashV2(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("alpha"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("beta"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := fileHash(filepath.Join(dir, "a.txt"))
	if err != nil || got != hashBytes([]byte("alpha")) {
		t.Fatal("bounded read changed file hash")
	}
	got, err = treeHash(dir)
	if err != nil || got != framedTreeOracle(hashOracleFile{"a.txt", "alpha"}, hashOracleFile{"b.txt", "beta"}) {
		t.Fatal("tree hash differs from v2 framing")
	}
}

func TestLocalGitMetadataReadersBoundInput(t *testing.T) {
	dir := t.TempDir()
	if err := syscall.Mkfifo(filepath.Join(dir, ".git"), 0600); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := readLocalGitConfig(dir); done <- err }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("FIFO Git pointer accepted")
		}
	case <-time.After(time.Second):
		t.Fatal("Git metadata reader blocked")
	}
}
