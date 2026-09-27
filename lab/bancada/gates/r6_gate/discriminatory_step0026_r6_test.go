/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package r6gate

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

var compiledVariants = []string{"M01", "M04", "M09", "M11", "M15", "M16", "M18", "M25", "M27", "M28", "M29", "M30", "M31", "M32", "M35", "M36", "M38", "M44", "M45", "M46", "M47", "M48", "M51", "M52", "M53", "M54", "M55", "M56", "M59", "M60"}

type processWitness struct{ FailedTests []string }

func TestStep0026R6CompiledVariantsUseOneBlindConformanceSuiteThreeByThree(t *testing.T) {
	repo, err := filepath.Abs(filepath.Join("..", "..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	assertCommonSuiteCopies(t, repo)
	assertNoRuntimeMutationSelector(t, repo)
	orders := [][]string{append([]string(nil), compiledVariants...), reverse(compiledVariants), permutation(compiledVariants)}
	baseline := map[string]processWitness{}
	for _, order := range orders {
		for repetition := 0; repetition < 3; repetition++ {
			if output, err := runGoTest(repo, "./lab/bancada/gates/r6_gate/reference"); err != nil {
				t.Fatalf("reference failed: %v\n%s", err, output)
			}
			results := runOrder(t, repo, order)
			for _, name := range order {
				w := results[name]
				if len(w.FailedTests) == 0 {
					t.Fatalf("compiled variant survived: %s", name)
				}
				if prior, ok := baseline[name]; ok && !equalStrings(prior.FailedTests, w.FailedTests) {
					t.Fatalf("unstable witness %s: %v != %v", name, prior.FailedTests, w.FailedTests)
				}
				baseline[name] = w
			}
		}
	}
	if len(baseline) != len(compiledVariants) {
		t.Fatalf("executed %d variants", len(baseline))
	}
}

func TestStep0026R6CertifiedPredecessorSurface(t *testing.T) {
	repo, err := filepath.Abs(filepath.Join("..", "..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	runCertifiedPredecessor(t, repo)
}

func runOrder(t *testing.T, repo string, order []string) map[string]processWitness {
	t.Helper()
	type result struct {
		name   string
		w      processWitness
		err    error
		output string
	}
	out := make(chan result, len(order))
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for _, name := range order {
		name := name
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			output, err := runGoTest(repo, "./lab/bancada/gates/r6_gate/testdata/mutants/"+name)
			<-sem
			if err == nil {
				out <- result{name: name, err: fmt.Errorf("process passed"), output: output}
				return
			}
			w := parseFailures(output)
			if len(w.FailedTests) == 0 {
				out <- result{name: name, err: fmt.Errorf("process failed without conformance witness: %v", err), output: output}
				return
			}
			out <- result{name: name, w: w}
		}()
	}
	wg.Wait()
	close(out)
	results := map[string]processWitness{}
	for r := range out {
		if r.err != nil {
			t.Fatalf("%s: %v\n%s", r.name, r.err, r.output)
		}
		results[r.name] = r.w
	}
	return results
}

func runGoTest(repo, pkg string) (string, error) {
	cmd := exec.Command("go", "test", pkg, "-run", "^TestCommonStep0026Conformance$", "-count=1", "-json")
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "GOCACHE=/tmp/bateia-go-cache-0026-r6")
	raw, err := cmd.CombinedOutput()
	return string(raw), err
}
func parseFailures(output string) processWitness {
	set := map[string]bool{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		var event struct{ Action, Test string }
		if json.Unmarshal(scanner.Bytes(), &event) == nil && event.Action == "fail" && strings.Contains(event.Test, "/") {
			set[event.Test] = true
		}
	}
	tests := make([]string, 0, len(set))
	for name := range set {
		tests = append(tests, name)
	}
	sort.Strings(tests)
	return processWitness{FailedTests: tests}
}

func assertCommonSuiteCopies(t *testing.T, repo string) {
	t.Helper()
	var want string
	paths := append([]string{"reference"}, variantPaths()...)
	for _, part := range paths {
		path := filepath.Join(repo, "lab", "bancada", "gates", "r6_gate", part, "conformance_test.go")
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		d := sha256.Sum256(raw)
		got := hex.EncodeToString(d[:])
		if want == "" {
			want = got
		} else if got != want {
			t.Fatalf("suite copy differs: %s", path)
		}
	}
}
func variantPaths() []string {
	out := make([]string, len(compiledVariants))
	for i, name := range compiledVariants {
		out[i] = filepath.Join("testdata", "mutants", name)
	}
	return out
}
func assertNoRuntimeMutationSelector(t *testing.T, repo string) {
	t.Helper()
	for _, part := range append([]string{"conformance", "reference"}, variantPaths()...) {
		entries, err := os.ReadDir(filepath.Join(repo, "lab", "bancada", "gates", "r6_gate", part))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".go" || entry.Name() == "conformance_test.go" {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(repo, "lab", "bancada", "gates", "r6_gate", part, entry.Name()))
			if err != nil {
				t.Fatal(err)
			}
			text := string(raw)
			if strings.Contains(text, "runtimeMutation") {
				t.Fatalf("runtime selector in %s/%s", part, entry.Name())
			}
		}
	}
}
func runCertifiedPredecessor(t *testing.T, repo string) {
	t.Helper()
	args := []string{"test", "./02_shell/mcp", "-run", "TestStep0026R3LegacyBenchmarkAndRecommendationUseIndependentByteGoldens|TestStep0026R3PredecessorServerHashesRemainCertified|TestStep0025LegacySuccessfulSearchMatchesFrozenBaselineBytes|TestStep0025LegacySuccessfulSearchUsesEvidenceContract", "-count=1"}
	cmd := exec.Command("go", args...)
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "GOCACHE=/tmp/bateia-go-cache-0026-r6")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("certified predecessor: %v\n%s", err, output)
	}
}
func reverse(in []string) []string {
	out := append([]string(nil), in...)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
func permutation(in []string) []string {
	out := append([]string(nil), in...)
	sort.Slice(out, func(i, j int) bool {
		a := sha256.Sum256([]byte("781fd3b7:05140182:" + out[i]))
		b := sha256.Sum256([]byte("781fd3b7:05140182:" + out[j]))
		return string(a[:]) < string(b[:])
	})
	return out
}
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
