/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/cmd/observer/main_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	concurrency "bateia/lab/bancada/concorrencia/0109"
)

func fixture(t *testing.T) (string, []byte) {
	t.Helper()
	code := "package p;func f(){_=make(chan int,2)}"
	r := concurrency.Request{SchemaVersion: "1", Origin: "CLI_FIXTURE", Files: []concurrency.SourceFile{{Path: "f.go", Content: code, SHA256: concurrency.Digest(code)}}}
	data, _ := json.Marshal(r)
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path, data
}

type shortWriter struct{}

func (shortWriter) Write(data []byte) (int, error) { return len(data) - 1, nil }

func TestObserveOutputAndProtectedPublication(t *testing.T) {
	input, original := fixture(t)
	var out bytes.Buffer
	if err := run([]string{"observe", "--input", input}, &out); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "output.json")
	if err := run([]string{"observe", "--input", input, "--output", output}, io.Discard); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(data, out.Bytes()) {
		t.Fatal("stdout/file differ")
	}
	info, _ := os.Stat(output)
	if info.Mode().Perm() != 0600 {
		t.Fatal("output permissions")
	}
	for _, alias := range []string{input, output} {
		before, _ := os.ReadFile(alias)
		if err := run([]string{"observe", "--input", input, "--output", alias}, io.Discard); err == nil {
			t.Fatal("existing file replaced")
		}
		after, _ := os.ReadFile(alias)
		if !bytes.Equal(before, after) {
			t.Fatal("existing bytes changed")
		}
	}
	link := filepath.Join(t.TempDir(), "alias.json")
	if err := os.Link(input, link); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"observe", "--input", input, "--output", link}, io.Discard); err == nil {
		t.Fatal("hardlink overwritten")
	}
	link = filepath.Join(t.TempDir(), "symlink.json")
	if err := os.Symlink(input, link); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"observe", "--input", input, "--output", link}, io.Discard); err == nil {
		t.Fatal("symlink overwritten")
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(after, original) {
		t.Fatal("input modified")
	}
	if err := run([]string{"observe", "--input", input}, shortWriter{}); err != io.ErrShortWrite {
		t.Fatalf("short write %v", err)
	}
}

func TestStrictInputHasNoPublicationOnError(t *testing.T) {
	_, valid := fixture(t)
	for _, data := range []string{
		`{"schema_version":"1","schema_version":"1"}`,
		`{"schema_version":"1","unexpected":true}`,
		string(valid) + "\n{}",
		strings.Repeat(" ", (2<<20)+1),
		`null`,
	} {
		dir := t.TempDir()
		input, output := filepath.Join(dir, "input.json"), filepath.Join(dir, "output.json")
		if err := os.WriteFile(input, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if err := run([]string{"observe", "--input", input, "--output", output}, io.Discard); err == nil {
			t.Fatal("invalid JSON accepted")
		}
		if _, err := os.Stat(output); !os.IsNotExist(err) {
			t.Fatal("output published after error")
		}
	}
}

func TestEvaluationDivergenceReportsAndFails(t *testing.T) {
	code := "package p;func f(){_=make(chan int,2)}"
	corpus := concurrency.Corpus{SchemaVersion: "1", ReferenceOrigin: "AUTHOR_REVIEWED_EXPLORATORY_NOT_INDEPENDENT", Cases: []concurrency.Case{{ID: "wrong", Family: "buffer", Request: concurrency.Request{SchemaVersion: "1", Origin: "CLI_FIXTURE", Files: []concurrency.SourceFile{{Path: "f.go", Content: code, SHA256: concurrency.Digest(code)}}}}}}
	data, _ := json.Marshal(corpus)
	input := filepath.Join(t.TempDir(), "corpus.json")
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"evaluate", "--corpus", input}, &out); err == nil {
		t.Fatal("negative control returned success")
	}
	var result concurrency.Evaluation
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Passed != 0 || len(result.Results[0].Mismatches) == 0 {
		t.Fatal("no discriminatory report")
	}
}

func TestFrozenEnvironmentCLIAndDependencyContrast(t *testing.T) {
	base := filepath.Join("..", "..", "..", "0111", "testdata")
	input := filepath.Join(base, "conc-input.json")
	envPath := filepath.Join(base, "environment.json")
	var out bytes.Buffer
	if err := run([]string{"observe", "--input", input, "--environment", envPath}, &out); err != nil {
		t.Fatal(err)
	}
	var report concurrency.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	resolved := 0
	for _, o := range report.Observations {
		if o.Property == "channel_type_identity" {
			resolved++
			if o.Resource != "limiter" || o.Authority != "GO_TYPES_SUPPLIED_ENVIRONMENT" || o.Location.Path != "pool/pool.go" {
				t.Fatal("unbound limiter resolution")
			}
		}
	}
	if resolved != 1 || report.EnvironmentSHA256 == "" || len(report.TypeDiagnostics) != 0 {
		t.Fatalf("wrong frozen result %+v", report)
	}
	out.Reset()
	if err := run([]string{"observe", "--input", input}, &out); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	for _, o := range report.Observations {
		if o.Property == "channel_type_identity" {
			t.Fatal("credit without environment")
		}
	}
	if len(report.TypeDiagnostics) == 0 {
		t.Fatal("dependency failure hidden")
	}
	bytesEnv, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	var env concurrency.Environment
	if err = json.Unmarshal(bytesEnv, &env); err != nil {
		t.Fatal(err)
	}
	env.Packages = env.Packages[:0] // Preserve real std data, remove source dependency bindings.
	data, _ := json.Marshal(env)
	missing := filepath.Join(t.TempDir(), "missing.json")
	if err = os.WriteFile(missing, data, 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err = run([]string{"observe", "--input", input, "--environment", missing}, &out); err != nil {
		t.Fatal(err)
	}
	var partial concurrency.Report
	if err = json.Unmarshal(out.Bytes(), &partial); err != nil {
		t.Fatal(err)
	}
	for _, o := range partial.Observations {
		if o.Property == "channel_type_identity" {
			t.Fatal("credit with unmapped source dependencies")
		}
	}
	// Integrity and strict JSON errors fail before output publication.
	env.Exports[0].SHA256 = "tampered"
	data, _ = json.Marshal(env)
	if err = os.WriteFile(missing, data, 0600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "result.json")
	if err = run([]string{"observe", "--input", input, "--environment", missing, "--output", output}, io.Discard); err == nil {
		t.Fatal("tampered environment accepted")
	}
	if _, err = os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("failed environment published output")
	}
	for _, bad := range []string{`{"schema_version":"1","schema_version":"1"}`, `{"unexpected":true}`} {
		if err = os.WriteFile(missing, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if err = run([]string{"observe", "--input", input, "--environment", missing}, io.Discard); err == nil {
			t.Fatal("non-strict environment accepted")
		}
	}
	if err = run([]string{"evaluate", "--corpus", input, "--environment", envPath}, io.Discard); err == nil {
		t.Fatal("evaluate accepted unsupported environment")
	}
}

func TestFrozenFlowRelations(t *testing.T) {
	base := filepath.Join("..", "..", "..", "0111", "testdata")
	var out bytes.Buffer
	if err := run([]string{"observe", "--input", filepath.Join(base, "conc-input.json"), "--environment", filepath.Join(base, "environment.json")}, &out); err != nil {
		t.Fatal(err)
	}
	var r concurrency.Report
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	sends, deferred, fields, launchCalls := 0, 0, 0, 0
	for _, e := range r.Relations {
		if e.Kind == "SEND_CASE_CONTAINS_GO_CALL" {
			sends++
			if e.From.Path != "pool/pool.go" || e.From.Line != 56 || e.To.Line != 59 || !strings.Contains(strings.Join(e.Conditions, ";"), "else(p.limiter == nil)") {
				t.Fatal("wrong reservation branch")
			}
		}
		if e.Kind == "DEFERRED_RECEIVER_RECEIVE" {
			deferred++
			if e.From.Line != 153 || e.To.Line != 172 {
				t.Fatal("wrong release relation")
			}
		}
		if e.Kind == "SAME_DECLARED_CHANNEL_FIELD" {
			fields++
		}
		if e.Kind == "STATIC_CALL_TARGET" && e.From.Path == "pool/pool.go" && e.To.Path == "waitgroup.go" && e.To.Line == 28 {
			launchCalls++
		}
	}
	if sends != 1 || deferred != 1 || fields != 1 || launchCalls != 2 {
		t.Fatalf("incorrect branches: %d %d %d %d", sends, deferred, fields, launchCalls)
	}
	gaps := map[string]bool{}
	for _, g := range r.Gaps {
		gaps[g.Code] = true
	}
	for _, gap := range []string{"WORKER_BOUND_UNRESOLVED", "ADMISSION_BOUND_UNRESOLVED", "CALLBACK_EXECUTION_UNRESOLVED", "CHANNEL_INSTANCE_IDENTITY_UNRESOLVED"} {
		if !gaps[gap] {
			t.Fatal("lost limitation", gap)
		}
	}
}

func TestFrozenCallbackChain(t *testing.T) {
	base := filepath.Join("..", "..", "..", "0111", "testdata")
	var out bytes.Buffer
	if err := run([]string{"observe", "--input", filepath.Join(base, "conc-input.json"), "--environment", filepath.Join(base, "environment.json")}, &out); err != nil {
		t.Fatal(err)
	}
	var r concurrency.Report
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	bindings := []concurrency.Relation{}
	var forwarded, invoked *concurrency.Relation
	for i := range r.CallbackRelations {
		e := &r.CallbackRelations[i]
		if e.Kind == "CALLBACK_ARGUMENT_BINDING" && e.To.Path == "waitgroup.go" && e.To.Line == 28 {
			bindings = append(bindings, *e)
		}
		if e.Kind == "CALLBACK_PARAMETER_FORWARDING" && e.From.Path == "waitgroup.go" {
			forwarded = e
		}
		if e.Kind == "CALLBACK_PARAMETER_INVOCATION" && e.From.Path == "panics/panics.go" {
			invoked = e
		}
	}
	if len(bindings) != 2 || forwarded == nil || invoked == nil || forwarded.ParameterDeclaration == nil || forwarded.ExecutionScope == nil {
		t.Fatal("callback chain incomplete")
	}
	for _, binding := range bindings {
		if binding.To != *forwarded.ParameterDeclaration {
			t.Fatal("wrong forwarded parameter")
		}
	}
	if forwarded.To != invoked.From || forwarded.ExecutionScope.Path != "waitgroup.go" || forwarded.ExecutionScope.Line != 30 {
		t.Fatal("wrong invocation/scope identity")
	}
	found := false
	for _, g := range r.Gaps {
		if g.Code == "CALLBACK_EXECUTION_UNRESOLVED" {
			found = true
		}
	}
	if !found {
		t.Fatal("structural chain promoted to executed callback")
	}
}
