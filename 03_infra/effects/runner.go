/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/effects/runner.md
 * @layer L3
 * @updated 2026-09-08
 */

package effects

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

const outputLimit = 1 << 20
const buildTimeout = 60 * time.Second

type Runner struct{}

func NewRunner() *Runner           { return &Runner{} }
func digest(data []byte) string    { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
func copyValues(v []int64) []int64 { return append([]int64{}, v...) }
func vectors(seed uint64) [][]int64 {
	result := [][]int64{{}, {7}, {2, 2, 1, 2}, {-1, -8, -3, 0}, {-4, 0, 9}, {9, 3, 1, -1}, {-1 << 63, 1<<63 - 1, 0, -1, 1}}
	// SplitMix64, specified here only for deterministic input generation, never as an oracle.
	next := func() uint64 {
		seed += 0x9e3779b97f4a7c15
		z := seed
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		return z ^ (z >> 31)
	}
	for i := 0; i < 32; i++ {
		v := make([]int64, int(next()%9))
		for j := range v {
			v[j] = int64(next()%17) - 8
		}
		result = append(result, v)
	}
	return result
}
func sourceDeclaration(src []byte, symbol string) (*domain.EffectDeclaration, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "candidate.go", src, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("unsupported source: %w", err)
	}
	if f.Name.Name != "candidate" {
		return nil, fmt.Errorf("package must be candidate")
	}
	for _, g := range f.Comments {
		for _, c := range g.List {
			if strings.HasPrefix(c.Text, "//go:") {
				return nil, fmt.Errorf("Go directives are outside scope")
			}
		}
	}
	var selected *ast.FuncDecl
	for _, d := range f.Decls {
		switch n := d.(type) {
		case *ast.GenDecl:
			if n.Tok != token.IMPORT {
				return nil, fmt.Errorf("only function declarations and permitted imports allowed")
			}
			for _, s := range n.Specs {
				i, ok := s.(*ast.ImportSpec)
				if !ok || i.Name != nil {
					return nil, fmt.Errorf("import aliases are outside scope")
				}
				name, e := strconv.Unquote(i.Path.Value)
				if e != nil || (name != "sort" && name != "slices") {
					return nil, fmt.Errorf("import outside scope")
				}
			}
		case *ast.FuncDecl:
			if n.Name.Name == "init" || n.Recv != nil || n.Body == nil {
				return nil, fmt.Errorf("init, methods and bodyless functions are outside scope")
			}
			if n.Name.Name == symbol {
				if selected != nil || !ast.IsExported(symbol) || n.Type.TypeParams != nil || !oneInt64Slice(n.Type.Params) || !oneInt64Slice(n.Type.Results) {
					return nil, fmt.Errorf("symbol must have exact func([]int64) []int64 signature")
				}
				selected = n
			}
		default:
			return nil, fmt.Errorf("declaration outside scope")
		}
	}
	if selected == nil {
		return nil, fmt.Errorf("exported symbol not found")
	}
	// File offsets ignore //line directives; line numbers count physical LF bytes.
	start := fset.PositionFor(selected.Pos(), false).Offset
	end := fset.PositionFor(selected.End(), false).Offset
	content := src[start:end]
	return &domain.EffectDeclaration{
		StartByte: start, EndByte: end,
		StartLine: bytes.Count(src[:start], []byte{'\n'}) + 1,
		EndLine:   bytes.Count(src[:end-1], []byte{'\n'}) + 1,
		Content:   string(content), SHA256: digest(content),
	}, nil
}
func allowedSource(src []byte, symbol string) error {
	_, err := sourceDeclaration(src, symbol)
	return err
}
func oneInt64Slice(list *ast.FieldList) bool {
	if list == nil || len(list.List) != 1 || len(list.List[0].Names) > 1 {
		return false
	}
	a, ok := list.List[0].Type.(*ast.ArrayType)
	if !ok || a.Len != nil {
		return false
	}
	id, ok := a.Elt.(*ast.Ident)
	return ok && id.Name == "int64"
}

const wrapperTemplate = `package main
import("encoding/json";"fmt";"os";"pilot/candidate")
type observation struct { InputBefore []int64 %s; InputAfter []int64 %s; Output []int64 %s; Completed bool %s; Failure string %s }
func clone(v []int64) []int64 { return append([]int64{},v...) }
func observe(input []int64)(o observation){
 argument:=clone(input);o.InputBefore=clone(argument);returned:=false
 defer func(){r:=recover();o.InputAfter=clone(argument);o.Completed=true;if !returned{o.Failure=fmt.Sprintf("panic: %%v",r)}}()
 result:=candidate.%s(argument);o.Output=clone(result);returned=true;return
}
func main(){var inputs [][]int64;if json.NewDecoder(os.Stdin).Decode(&inputs)!=nil{os.Exit(2)};e:=json.NewEncoder(os.Stdout);for _,v:=range inputs{if e.Encode(observe(v))!=nil{os.Exit(3)}}}
`

func wrapper(symbol string) []byte {
	return []byte(fmt.Sprintf(wrapperTemplate, "`json:\"input_before\"`", "`json:\"input_after\"`", "`json:\"output\"`", "`json:\"completed\"`", "`json:\"failure\"`", symbol))
}

type limitedOutput struct {
	bytes.Buffer
	exceeded bool
	cancel   context.CancelFunc
}

func (b *limitedOutput) Write(p []byte) (int, error) {
	if len(p) > outputLimit-b.Len() {
		remaining := outputLimit - b.Len()
		if remaining > 0 {
			b.Buffer.Write(p[:remaining])
		}
		b.exceeded = true
		b.cancel()
		return len(p), nil
	}
	return b.Buffer.Write(p)
}

type processResult struct {
	stdout, stderr []byte
	reason         string
	duration       time.Duration
	err            error
}

func run(ctx context.Context, timeout time.Duration, dir string, env []string, stdin []byte, program string, args ...string) processResult {
	child, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	o := &limitedOutput{cancel: cancel}
	e := &limitedOutput{cancel: cancel}
	cmd := exec.CommandContext(child, program, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = bytes.NewReader(stdin)
	cmd.Stdout = o
	cmd.Stderr = e
	cmd.WaitDelay = time.Second
	start := time.Now()
	err := cmd.Run()
	reason := "completed"
	if o.exceeded || e.exceeded {
		reason = "output_limit"
	} else if child.Err() != nil {
		if child.Err() == context.DeadlineExceeded {
			reason = "timeout"
		} else {
			reason = "canceled"
		}
	} else if err != nil {
		reason = "process_error"
	}
	return processResult{stdout: o.Bytes(), stderr: e.Bytes(), reason: reason, duration: time.Since(start), err: err}
}
func (r *Runner) Verify(ctx context.Context, req domain.EffectRequest) (report domain.EffectReport, err error) {
	startedAt := time.Now().UTC()
	defer func() {
		if err != nil {
			return
		}
		finishedAt := time.Now().UTC()
		// Clock adjustments must not produce an inverted observed interval.
		if finishedAt.Before(startedAt) {
			finishedAt = startedAt
		}
		report.Trace.StartedAt = startedAt.Format(time.RFC3339Nano)
		report.Trace.FinishedAt = finishedAt.Format(time.RFC3339Nano)
	}()
	if req.SourcePath == "" || !token.IsIdentifier(req.Symbol) || !ast.IsExported(req.Symbol) || req.Timeout <= 0 || req.Timeout > 30*time.Second {
		return report, fmt.Errorf("invalid effect request")
	}
	if req.ExpectedSourceSHA256 != "" {
		if _, e := hex.DecodeString(req.ExpectedSourceSHA256); len(req.ExpectedSourceSHA256) != 64 || e != nil {
			return report, fmt.Errorf("expected source SHA-256 must contain exactly 64 ASCII hexadecimal digits")
		}
		req.ExpectedSourceSHA256 = strings.ToLower(req.ExpectedSourceSHA256)
	}
	src, err := os.ReadFile(req.SourcePath)
	if err != nil {
		return report, err
	}
	actualSHA256 := digest(src)
	binding := domain.EffectSourceBinding{Status: "not_requested", ExpectedSHA256: req.ExpectedSourceSHA256, ActualSHA256: actualSHA256, Symbol: req.Symbol}
	if req.ExpectedSourceSHA256 != "" {
		binding.Status = "matched"
		if req.ExpectedSourceSHA256 != actualSHA256 {
			binding.Status = "mismatch"
		}
	}
	inputs := vectors(req.Seed)
	report = domain.EffectReport{SourceBinding: binding, SchemaVersion: "effect-run-v1", ContractID: "sort-preserve-v1", Status: "unknown", Reason: "not executed", SourceSHA256: actualSHA256, SourcePath: req.SourcePath, Symbol: req.Symbol, Seed: req.Seed, Cases: make([]domain.EffectCase, len(inputs)), Limitations: []string{
		"Evidence covers only the recorded int64 inputs; no general correctness or termination proof.",
		"Input values are checked after return; transient writes and memory sharing are not ruled out.",
		"A separate process is not a security sandbox or rigid memory limit. Do not execute hostile code.",
		"No claim about external effects; no automatic change to search scores or ActionAssessments.",
		"Duration is tool telemetry, not a benchmark of the candidate.",
		"Source binding covers the complete file only; no automatic link to search evidence, repository/commit identity or external dependencies.",
	}, Execution: map[string]any{"source_content": string(src), "source_sha256": actualSHA256, "goos": runtime.GOOS, "goarch": runtime.GOARCH, "execution_timeout": req.Timeout.String(), "build_timeout": buildTimeout.String(), "stdout_limit_bytes": outputLimit, "stderr_limit_bytes": outputLimit, "inputs": inputs, "duration_kind": "tool_telemetry", "stop_reason": "not_executed"}}
	report.Trace = domain.EffectTrace{
		SchemaVersion: "effect-trace-v1", Status: "unknown",
		Analyzer: analyzerIdentity(),
		Inputs: domain.EffectTraceInputs{
			ContractID: report.ContractID, SourceSHA256: actualSHA256,
			Symbol: req.Symbol, Seed: req.Seed, ExecutionTimeoutNS: int64(req.Timeout),
			BuildTimeoutNS: int64(buildTimeout), StdoutLimitBytes: outputLimit, StderrLimitBytes: outputLimit,
		},
	}
	report.Trace.Inputs.AnalyzerBinarySHA256 = report.Trace.Analyzer.BinarySHA256
	if report.Trace.Analyzer.Status != "recorded" {
		report.Trace.Reason = "analyzer identity unavailable: " + report.Trace.Analyzer.Reason
	} else if fingerprint := usecase.EffectTraceFingerprint(report.Trace.Inputs); fingerprint == "" {
		report.Trace.Reason = "effective request inputs cannot be fingerprinted"
	} else {
		report.Trace.Status = "recorded"
		report.Trace.Reason = "analyzer identity and effective request fingerprint recorded"
		report.Trace.RequestSHA256 = fingerprint
	}
	for i, v := range inputs {
		report.Cases[i] = domain.EffectCase{Observation: domain.EffectObservation{InputBefore: copyValues(v)}, Verdict: domain.EffectVerdict{Status: "unknown", Reason: "not executed"}}
	}
	stop := func(reason string) { report.Reason = reason; report.Execution["stop_reason"] = reason }
	if binding.Status == "mismatch" {
		stop("source_identity_mismatch")
		return report, nil
	}
	declaration, e := sourceDeclaration(src, req.Symbol)
	if e != nil {
		stop("outside_scope: " + e.Error())
		return report, nil
	}
	report.SourceBinding.Declaration = declaration
	root, err := os.MkdirTemp("", "bateia-effects-")
	if err != nil {
		return report, err
	}
	defer os.RemoveAll(root)
	for _, dir := range []string{"candidate", "run", "home", "tmp"} {
		if err = os.Mkdir(filepath.Join(root, dir), 0700); err != nil {
			return report, err
		}
	}
	wrap := wrapper(req.Symbol)
	report.Execution["wrapper_content"] = string(wrap)
	report.Execution["wrapper_sha256"] = digest(wrap)
	module := []byte("module pilot\n\ngo 1.23\n")
	report.Execution["module_content"] = string(module)
	report.Execution["module_sha256"] = digest(module)
	for name, data := range map[string][]byte{"go.mod": module, "candidate/candidate.go": src, "main.go": wrap} {
		if err = os.WriteFile(filepath.Join(root, name), data, 0600); err != nil {
			return report, err
		}
	}
	cacheDir := filepath.Join(os.TempDir(), "bateia-effects-build-cache-"+strconv.Itoa(os.Getuid()))
	if err = os.MkdirAll(cacheDir, 0700); err != nil {
		return report, err
	}
	cacheInfo, cacheErr := os.Lstat(cacheDir)
	if cacheErr != nil {
		return report, cacheErr
	}
	if !cacheInfo.IsDir() || cacheInfo.Mode().Perm()&0077 != 0 {
		return report, fmt.Errorf("build cache must be a private directory")
	}
	report.Execution["build_cache_policy"] = "reusable private tool cache under os.TempDir"
	goroot := runtime.GOROOT()
	goexe := filepath.Join(goroot, "bin", "go")
	env := []string{"PATH=" + filepath.Join(goroot, "bin"), "GOROOT=" + goroot, "HOME=" + filepath.Join(root, "home"), "TMPDIR=" + filepath.Join(root, "tmp"), "GOCACHE=" + cacheDir, "GOPATH=" + filepath.Join(root, "home", "go"), "GOPROXY=off", "GOSUMDB=off", "GOWORK=off", "GOENV=off", "GOTOOLCHAIN=local", "CGO_ENABLED=0"}
	version := run(ctx, 5*time.Second, root, env, nil, goexe, "version")
	report.Execution["go_version"] = strings.TrimSpace(string(version.stdout))
	if version.reason != "completed" {
		stop("toolchain_" + version.reason)
		return report, nil
	}
	binary := filepath.Join(root, "candidate-runner")
	args := []string{"build", "-trimpath", "-o", binary, "."}
	report.Execution["build_argv"] = append([]string{goexe}, args...)
	built := run(ctx, buildTimeout, root, env, nil, goexe, args...)
	report.Execution["build_duration_ms"] = built.duration.Milliseconds()
	report.Execution["build_stdout"] = string(built.stdout)
	report.Execution["build_stderr"] = string(built.stderr)
	if built.reason != "completed" {
		stop("build_" + built.reason)
		return report, nil
	}
	bin, err := os.ReadFile(binary)
	if err != nil {
		return report, err
	}
	report.Execution["binary_sha256"] = digest(bin)
	data, _ := json.Marshal(inputs)
	runEnv := []string{"HOME=" + filepath.Join(root, "run"), "TMPDIR=" + filepath.Join(root, "run"), "GOPROXY=off", "GOSUMDB=off", "GOWORK=off", "GOENV=off", "GOTOOLCHAIN=local"}
	report.Execution["argv"] = []string{binary}
	report.Execution["environment_policy"] = "explicit minimal environment, no inherited credentials"
	executed := run(ctx, req.Timeout, filepath.Join(root, "run"), runEnv, data, binary)
	report.Execution["duration_ms"] = executed.duration.Milliseconds()
	report.Execution["stderr"] = string(executed.stderr)
	report.Execution["stop_reason"] = executed.reason
	dec := json.NewDecoder(bytes.NewReader(executed.stdout))
	dec.DisallowUnknownFields()
	complete := 0
	malformed := false
	for i := range inputs {
		var o domain.EffectObservation
		if e := dec.Decode(&o); e != nil {
			if e != io.EOF {
				malformed = true
			}
			break
		}
		// Before snapshots must match the tool's input schedule; candidate output never supplies the oracle.
		if !sameValues(o.InputBefore, inputs[i]) || !o.Completed {
			malformed = true
			break
		}
		report.Cases[i] = domain.EffectCase{Observation: o, Verdict: usecase.EvaluateSortEffect(o)}
		complete++
	}
	var extra any
	if e := dec.Decode(&extra); e != io.EOF {
		malformed = true
	}
	reason := executed.reason
	if malformed {
		reason = "malformed_protocol"
	} else if complete != len(inputs) && reason == "completed" {
		reason = "incomplete_protocol"
	}
	report.Execution["stop_reason"] = reason
	for i := complete; i < len(inputs); i++ {
		report.Cases[i].Observation.Failure = reason
		report.Cases[i].Verdict = usecase.EvaluateSortEffect(report.Cases[i].Observation)
	}
	report.Status = "supported_in_scope"
	report.Reason = "all recorded cases completed and satisfy sort-preserve-v1"
	if reason != "completed" || complete != len(inputs) {
		report.Status = "unknown"
		report.Reason = reason
	}
	for _, c := range report.Cases {
		if c.Verdict.Status == "counterexample" {
			report.Status = "counterexample"
			report.Reason = c.Verdict.Reason
			break
		}
	}
	return report, nil
}
func sameValues(a, b []int64) bool {
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
