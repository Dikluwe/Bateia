/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/effects/runner_test.md
 * @layer L3
 * @updated 2026-09-08
 */
package effects_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/effects"
)

const independentValidSource = `package candidate
import "sort"
func Sort(input []int64) []int64 { result:=clone(input); sort.Slice(result,func(i,j int)bool{return result[i]<result[j]}); return result }
func clone(input []int64) []int64 { return append([]int64(nil),input...) }
`

func effectFixture(t *testing.T, source string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "candidate.go")
	if err := os.WriteFile(p, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}
func effectDigest(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func effectRequest(p string) domain.EffectRequest {
	return domain.EffectRequest{SourcePath: p, Symbol: "Sort", Seed: 687, Timeout: 2 * time.Second}
}
func verifyEffect(t *testing.T, source string, timeout time.Duration) domain.EffectReport {
	t.Helper()
	r := effectRequest(effectFixture(t, source))
	r.Timeout = timeout
	report, err := effects.NewRunner().Verify(context.Background(), r)
	if err != nil {
		t.Fatalf("fixture evaluation: %v", err)
	}
	return report
}

// This oracle uses value counts and pair comparisons, not the candidate's sorting algorithm.
func independentObservationValid(o domain.EffectObservation) bool {
	if !o.Completed || o.Failure != "" || len(o.InputBefore) != len(o.InputAfter) || len(o.InputBefore) != len(o.Output) {
		return false
	}
	counts := map[int64]int{}
	for i, v := range o.InputBefore {
		if v != o.InputAfter[i] {
			return false
		}
		counts[v]++
	}
	for i, v := range o.Output {
		if i > 0 && o.Output[i-1] > v {
			return false
		}
		counts[v]--
	}
	for _, n := range counts {
		if n != 0 {
			return false
		}
	}
	return true
}
func executionStrings(v any) []string {
	switch x := v.(type) {
	case string:
		return []string{x}
	case []any:
		var r []string
		for _, e := range x {
			r = append(r, executionStrings(e)...)
		}
		return r
	case map[string]any:
		var r []string
		for _, e := range x {
			r = append(r, executionStrings(e)...)
		}
		return r
	}
	return nil
}
func TestEffectRunnerValidReceiptAndReproduction(t *testing.T) {
	const secret = "effect-test-secret-should-never-appear-0062"
	t.Setenv("BATEIA_TEST_CREDENTIAL", secret)
	sourcePath := effectFixture(t, independentValidSource)
	// Unrelated original-project files are deliberately invalid. Only source bytes belong in the isolated module.
	if err := os.WriteFile(filepath.Join(filepath.Dir(sourcePath), "unrelated.go"), []byte("this is not Go"), 0600); err != nil {
		t.Fatal(err)
	}
	req := effectRequest(sourcePath)
	runner := effects.NewRunner()
	first, err := runner.Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != "supported_in_scope" {
		t.Fatalf("valid fixture: %s: %s", first.Status, first.Reason)
	}
	if first.SchemaVersion != "effect-run-v1" || first.ContractID != "sort-preserve-v1" || first.SourceSHA256 != effectDigest(independentValidSource) || first.SourcePath != sourcePath || first.Symbol != "Sort" || first.Seed != req.Seed {
		t.Fatalf("receipt identity mismatch: %+v", first)
	}
	if first.Reason == "" || len(first.Limitations) == 0 || len(first.Execution) == 0 {
		t.Fatal("missing reason, limitations or execution evidence")
	}
	if len(first.Cases) < 39 {
		t.Fatalf("missing fixed cases or 32 generated cases: %d", len(first.Cases))
	}
	empty, single, duplicate, negative, sorted, reverse, extremeMin, extremeMax := false, false, false, false, false, false, false, false
	small := 0
	for i, c := range first.Cases {
		if c.Verdict.Status != "supported_in_scope" || !independentObservationValid(c.Observation) {
			t.Fatalf("case %d not independently valid: %+v", i, c)
		}
		a := c.Observation.InputBefore
		empty = empty || len(a) == 0
		single = single || len(a) == 1
		seen := map[int64]bool{}
		asc, desc := len(a) > 1, len(a) > 1
		bounded := len(a) <= 8
		for j, v := range a {
			duplicate = duplicate || seen[v]
			seen[v] = true
			negative = negative || v < 0
			extremeMin = extremeMin || v == math.MinInt64
			extremeMax = extremeMax || v == math.MaxInt64
			if v < -8 || v > 8 {
				bounded = false
			}
			if j > 0 {
				asc = asc && a[j-1] <= v
				desc = desc && a[j-1] > v
			}
		}
		sorted = sorted || asc
		reverse = reverse || desc
		if bounded {
			small++
		}
	}
	if !(empty && single && duplicate && negative && sorted && reverse && extremeMin && extremeMax) || small < 32 {
		t.Fatalf("input coverage missing: empty=%v single=%v duplicate=%v negative=%v sorted=%v reverse=%v min=%v max=%v bounded=%d", empty, single, duplicate, negative, sorted, reverse, extremeMin, extremeMax, small)
	}
	encoded, err := json.Marshal(first.Execution)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), secret) {
		t.Fatal("credential sentinel leaked into receipt")
	}
	var execution any
	if err := json.Unmarshal(encoded, &execution); err != nil {
		t.Fatal(err)
	}
	vals := executionStrings(execution)
	all := strings.Join(vals, "\n")
	if !strings.Contains(all, independentValidSource) || !strings.Contains(all, effectDigest(independentValidSource)) {
		t.Fatal("execution missing actual source content/hash")
	}
	wrapperFound := false
	hashes := map[string]bool{}
	for _, s := range vals {
		if len(s) == 64 {
			if b, e := hex.DecodeString(s); e == nil && len(b) == 32 {
				hashes[s] = true
			}
		}
		if strings.Contains(s, "package main") {
			wrapperFound = true
			if !strings.Contains(all, effectDigest(s)) {
				t.Error("wrapper content lacks matching SHA256")
			}
		}
		// Any absolute artifact path recorded beneath the system temp dir must be cleaned, except source and reusable cache paths.
		if filepath.IsAbs(s) && strings.HasPrefix(filepath.Clean(s), filepath.Clean(os.TempDir())+string(os.PathSeparator)) && s != sourcePath && s != filepath.Dir(sourcePath) && !strings.Contains(strings.ToLower(s), "cache") {
			if _, e := os.Stat(s); e == nil {
				t.Errorf("temporary execution artifact still exists: %s", s)
			}
		}
	}
	if !wrapperFound || len(hashes) < 3 {
		t.Fatalf("missing wrapper/source/binary evidence: wrapper=%v distinct hashes=%d", wrapperFound, len(hashes))
	}
	if !strings.Contains(all, runtime.GOOS) || !strings.Contains(all, runtime.GOARCH) || !strings.Contains(all, "go1.") {
		t.Error("missing Go/toolchain target metadata")
	}
	second, err := runner.Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if second.Status != first.Status || second.Seed != first.Seed || second.SourceSHA256 != first.SourceSHA256 || !reflect.DeepEqual(first.Cases, second.Cases) {
		t.Fatal("same source and seed did not reproduce vectors, observations and verdicts")
	}
	contents, err := os.ReadFile(sourcePath)
	if err != nil || string(contents) != independentValidSource {
		t.Fatal("source file changed or removed")
	}
}
func TestEffectRunnerCounterexamples(t *testing.T) {
	fixtures := map[string]string{
		"mutates_input": `package candidate
import "slices"
func Sort(v []int64) []int64 { slices.Sort(v); return v }
`,
		"panic": `package candidate
func Sort(v []int64) []int64 { panic("independent fixture panic") }
`,
	}
	for name, source := range fixtures {
		t.Run(name, func(t *testing.T) {
			r := verifyEffect(t, source, 2*time.Second)
			if r.Status != "counterexample" {
				t.Fatalf("%s: %s", r.Status, r.Reason)
			}
			witness := false
			for _, c := range r.Cases {
				if c.Verdict.Status == "counterexample" {
					witness = true
					if !c.Observation.Completed || independentObservationValid(c.Observation) {
						t.Fatalf("invalid counterexample witness: %+v", c)
					}
				}
			}
			if !witness {
				t.Fatal("global counterexample has no observable witness")
			}
		})
	}
}
func TestEffectRunnerTimeoutIsUnknown(t *testing.T) {
	started := time.Now()
	r := verifyEffect(t, `package candidate
func Sort(v []int64) []int64 { for {} }
`, 80*time.Millisecond)
	if r.Status != "unknown" || r.Reason == "" {
		t.Fatalf("timeout became %s: %s", r.Status, r.Reason)
	}
	if time.Since(started) > 65*time.Second {
		t.Fatal("build+run exceeded bounded contract budgets")
	}
	for _, c := range r.Cases {
		if c.Observation.Completed && c.Verdict.Status == "supported_in_scope" {
			t.Fatal("nonterminating source reported completed case")
		}
	}
}
func TestEffectRunnerCanceledContextCannotSucceed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	r, err := effects.NewRunner().Verify(ctx, effectRequest(effectFixture(t, independentValidSource)))
	if err == nil && r.Status != "unknown" {
		t.Fatalf("canceled operation produced %s", r.Status)
	}
	if time.Since(started) > 2*time.Second {
		t.Fatal("pre-canceled context did not stop promptly")
	}
}
func TestEffectRunnerOutOfScopeIsNotExecuted(t *testing.T) {
	// Prohibited declarations panic or loop if executed, exposing missing preflight validation without external effects.
	for name, source := range map[string]string{
		"package":    `package other; func Sort(v []int64) []int64 { for {} }`,
		"global":     `package candidate; var global = fail(); func fail() int { panic("executed global") }; func Sort(v []int64) []int64 { return v }`,
		"init":       `package candidate; func init() { panic("executed init") }; func Sort(v []int64) []int64 { return v }`,
		"type":       `package candidate; type V []int64; func Sort(v []int64) []int64 { for {} }`,
		"import":     `package candidate; import "fmt"; func Sort(v []int64) []int64 { fmt.Println("executed"); return v }`,
		"alias":      `package candidate; import s "sort"; func Sort(v []int64) []int64 { s.Slice(v,func(i,j int) bool { return v[i]<v[j] }); return v }`,
		"blank":      `package candidate; import _ "sort"; func Sort(v []int64) []int64 { for {} }`,
		"dot":        `package candidate; import . "sort"; func Sort(v []int64) []int64 { Slice(v,func(i,j int) bool { return v[i]<v[j] }); return v }`,
		"signature":  `package candidate; func Sort(v []int) []int { for {} }`,
		"generic":    `package candidate; func Sort[T any](v []int64) []int64 { for {} }`,
		"unexported": `package candidate; func sort(v []int64) []int64 { for {} }`,
		"directive":  "package candidate\n//go:noinline\nfunc Sort(v []int64) []int64 { for {} }\n",
	} {
		t.Run(name, func(t *testing.T) {
			r := verifyEffect(t, source, 80*time.Millisecond)
			if r.Status != "unknown" || r.Reason == "" {
				t.Fatalf("outside scope: %s: %s", r.Status, r.Reason)
			}
			for _, c := range r.Cases {
				if c.Observation.Completed {
					t.Fatalf("out-of-scope candidate executed: %+v", c)
				}
			}
		})
	}
}
