/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/effects/analyzer_identity_test.md
 * @layer L3
 * @updated 2026-09-08
 */
package effects_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/effects"
)

func trace64Digest(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }
func trace64Self(t *testing.T) string {
	t.Helper()
	path := "/proc/self/exe"
	if runtime.GOOS != "linux" {
		var err error
		path, err = os.Executable()
		if err != nil {
			t.Skipf("cannot observe executable: %v", err)
		}
	}
	f, err := os.Open(path)
	if err != nil {
		t.Skipf("cannot independently observe executable: %v", err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
func trace64InfraOracle(in domain.EffectTraceInputs) string {
	return trace64Digest([]byte(strings.Join([]string{"effect-request-v1", in.ContractID, strings.ToLower(in.AnalyzerBinarySHA256), strings.ToLower(in.SourceSHA256), in.Symbol, strconv.FormatUint(in.Seed, 10), strconv.FormatInt(in.ExecutionTimeoutNS, 10), strconv.FormatInt(in.BuildTimeoutNS, 10), strconv.FormatInt(in.StdoutLimitBytes, 10), strconv.FormatInt(in.StderrLimitBytes, 10)}, "\n") + "\n"))
}
func trace64CheckReport(t *testing.T, r domain.EffectReport, req domain.EffectRequest, source, self string) {
	t.Helper()
	tr := r.Trace
	if r.SchemaVersion != "effect-run-v1" || tr.SchemaVersion != "effect-trace-v1" {
		t.Fatalf("schema=%q trace=%q", r.SchemaVersion, tr.SchemaVersion)
	}
	if tr.Status != "recorded" || tr.Analyzer.Status != "recorded" {
		t.Fatalf("trace statuses: %+v", tr)
	}
	if tr.Analyzer.BinarySHA256 != self {
		t.Errorf("analyzer hash=%q want actual executable %q", tr.Analyzer.BinarySHA256, self)
	}
	if tr.Analyzer.GoVersion != runtime.Version() {
		t.Errorf("Go version=%q want %q", tr.Analyzer.GoVersion, runtime.Version())
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if tr.Analyzer.ModulePath != bi.Main.Path || tr.Analyzer.ModuleVersion != bi.Main.Version {
			t.Errorf("module metadata mismatch: %+v vs %+v", tr.Analyzer, bi.Main)
		}
		revision := ""
		var modified *bool
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				revision = s.Value
			}
			if s.Key == "vcs.modified" {
				if b, err := strconv.ParseBool(s.Value); err == nil {
					modified = &b
				}
			}
		}
		if tr.Analyzer.VCSRevision != revision {
			t.Errorf("VCS revision=%q want %q", tr.Analyzer.VCSRevision, revision)
		}
		if modified == nil {
			if tr.Analyzer.VCSModified != nil {
				t.Error("missing vcs.modified invented")
			}
		} else if tr.Analyzer.VCSModified == nil || *tr.Analyzer.VCSModified != *modified {
			t.Error("vcs.modified mismatch")
		}
	}
	start, err := time.Parse(time.RFC3339Nano, tr.StartedAt)
	if err != nil {
		t.Fatalf("started_at: %v", err)
	}
	end, err := time.Parse(time.RFC3339Nano, tr.FinishedAt)
	if err != nil {
		t.Fatalf("finished_at: %v", err)
	}
	if !strings.HasSuffix(tr.StartedAt, "Z") || !strings.HasSuffix(tr.FinishedAt, "Z") || start.After(end) {
		t.Errorf("invalid UTC interval %q .. %q", tr.StartedAt, tr.FinishedAt)
	}
	want := domain.EffectTraceInputs{ContractID: "sort-preserve-v1", AnalyzerBinarySHA256: self, SourceSHA256: trace64Digest([]byte(source)), Symbol: req.Symbol, Seed: req.Seed, ExecutionTimeoutNS: int64(req.Timeout), BuildTimeoutNS: int64(60 * time.Second), StdoutLimitBytes: 1 << 20, StderrLimitBytes: 1 << 20}
	if tr.Inputs != want {
		t.Errorf("inputs=%+v want %+v", tr.Inputs, want)
	}
	if tr.RequestSHA256 != trace64InfraOracle(want) {
		t.Errorf("request hash=%q want independent oracle %q", tr.RequestSHA256, trace64InfraOracle(want))
	}
	if r.SourceSHA256 != want.SourceSHA256 || r.SourceBinding.ActualSHA256 != want.SourceSHA256 {
		t.Error("trace source disagrees with source binding")
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]json.RawMessage
	if err = json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	var trace map[string]json.RawMessage
	if err = json.Unmarshal(doc["trace"], &trace); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"schema_version", "status", "reason", "started_at", "finished_at", "request_sha256", "analyzer", "inputs"} {
		if _, ok := trace[key]; !ok {
			t.Errorf("missing trace JSON key %q", key)
		}
	}
	var analyzer map[string]json.RawMessage
	if err = json.Unmarshal(trace["analyzer"], &analyzer); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"status", "reason", "binary_sha256", "module_path", "module_version", "go_version", "vcs_revision", "vcs_modified"} {
		if _, ok := analyzer[key]; !ok {
			t.Errorf("missing analyzer JSON key %q", key)
		}
	}
	var inputs map[string]json.RawMessage
	if err = json.Unmarshal(trace["inputs"], &inputs); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"contract_id", "analyzer_binary_sha256", "source_sha256", "symbol", "seed", "execution_timeout_ns", "build_timeout_ns", "stdout_limit_bytes", "stderr_limit_bytes"} {
		if _, ok := inputs[key]; !ok {
			t.Errorf("missing inputs JSON key %q", key)
		}
	}
}
func TestTrace64RunnerIdentityRepeatAndMismatch(t *testing.T) {
	self := trace64Self(t)
	const source = "package candidate\nimport \"sort\"\nfunc Sort(a []int64) []int64 { b := append([]int64(nil), a...); sort.Slice(b, func(i,j int) bool { return b[i]<b[j] }); return b }\n"
	dir := t.TempDir()
	path := filepath.Join(dir, "first.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	req := domain.EffectRequest{SourcePath: path, Symbol: "Sort", Seed: ^uint64(0), Timeout: 5 * time.Second}
	runner := effects.NewRunner()
	first, err := runner.Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	trace64CheckReport(t, first, req, source, self)
	if first.Status != "supported_in_scope" || first.SourceBinding.Status != "not_requested" || len(first.Cases) == 0 {
		t.Fatalf("first result: status=%q binding=%q cases=%d reason=%q", first.Status, first.SourceBinding.Status, len(first.Cases), first.Reason)
	}
	if first.Execution["build_timeout"] != "1m0s" || first.Execution["execution_timeout"] != req.Timeout.String() {
		t.Error("execution timeout metadata disagrees with trace")
	}
	for _, key := range []string{"stdout_limit_bytes", "stderr_limit_bytes"} {
		if fmt.Sprint(first.Execution[key]) != "1048576" {
			t.Errorf("%s=%v", key, first.Execution[key])
		}
	}
	req.SourcePath = filepath.Join(dir, "relocated.go")
	req.ExpectedSourceSHA256 = strings.ToUpper(trace64Digest([]byte(source)))
	if err := os.WriteFile(req.SourcePath, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	second, err := runner.Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	trace64CheckReport(t, second, req, source, self)
	if second.Status != "supported_in_scope" || second.SourceBinding.Status != "matched" {
		t.Fatalf("second result: %q/%q", second.Status, second.SourceBinding.Status)
	}
	if first.Trace.RequestSHA256 != second.Trace.RequestSHA256 {
		t.Error("path or expected hash changed fingerprint")
	}
	req.ExpectedSourceSHA256 = strings.Repeat("0", 64)
	mismatch, err := runner.Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	trace64CheckReport(t, mismatch, req, source, self)
	if mismatch.Status != "unknown" || mismatch.Reason != "source_identity_mismatch" || mismatch.SourceBinding.Status != "mismatch" || mismatch.SourceBinding.Declaration != nil || trace64HasCompletedCase(mismatch.Cases) {
		t.Fatalf("mismatch violated preflight: status=%q reason=%q binding=%+v cases=%d", mismatch.Status, mismatch.Reason, mismatch.SourceBinding, len(mismatch.Cases))
	}
	if mismatch.Execution["stop_reason"] != "source_identity_mismatch" {
		t.Error("missing mismatch stop reason")
	}
	if mismatch.Trace.RequestSHA256 != first.Trace.RequestSHA256 {
		t.Error("result/expected hash altered fingerprint")
	}
	for key, value := range mismatch.Execution {
		if (strings.Contains(key, "binary") || strings.Contains(key, "wrapper")) && value != nil && fmt.Sprint(value) != "" {
			t.Errorf("mismatch produced execution artifact %q=%v", key, value)
		}
	}
}
func TestTrace64RejectedSourceStillRecorded(t *testing.T) {
	self := trace64Self(t)
	const source = "package candidate\nimport \"os\"\nfunc Sort(a []int64) []int64 { os.Getenv(\"PATH\"); return a }\n"
	path := filepath.Join(t.TempDir(), "rejected.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	req := domain.EffectRequest{SourcePath: path, Symbol: "Sort", Seed: 7, Timeout: time.Second, ExpectedSourceSHA256: trace64Digest([]byte(source))}
	r, err := effects.NewRunner().Verify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	trace64CheckReport(t, r, req, source, self)
	if r.Status != "unknown" || r.SourceBinding.Status != "matched" || r.SourceBinding.Declaration != nil || trace64HasCompletedCase(r.Cases) {
		t.Fatalf("preflight result: status=%q binding=%+v cases=%d", r.Status, r.SourceBinding, len(r.Cases))
	}
}

func trace64HasCompletedCase(cases []domain.EffectCase) bool {
	for _, c := range cases {
		if c.Observation.Completed {
			return true
		}
	}
	return false
}
