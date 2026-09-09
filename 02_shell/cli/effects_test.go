/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/effects_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/02_shell/cli"
)

type effectFake struct {
	calls   int
	request domain.EffectRequest
	report  domain.EffectReport
	err     error
}

func (f *effectFake) Verify(_ context.Context, r domain.EffectRequest) (domain.EffectReport, error) {
	f.calls++
	f.request = r
	return f.report, f.err
}
func TestEffectsInvalidArgumentsNeverExecute(t *testing.T) {
	for _, args := range [][]string{
		{}, {"other"}, {"verify"}, {"verify", "--source", "x.go"},
		{"verify", "--symbol", "Sort"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--timeout", "0s"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--timeout", "-1s"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--timeout", "31s"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--seed", "-1"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--seed", "18446744073709551616"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "--bad"},
		{"verify", "--source", "x.go", "--symbol", "Sort", "extra"},
	} {
		t.Run(joinEffectArgs(args), func(t *testing.T) {
			f := &effectFake{}
			if err := cli.RunEffectsCommand(args, io.Discard, f); err == nil {
				t.Fatal("invalid arguments accepted")
			}
			if f.calls != 0 {
				t.Fatalf("executed %d times before argument validation", f.calls)
			}
		})
	}
}
func joinEffectArgs(args []string) string { b, _ := json.Marshal(args); return string(b) }
func TestEffectsHelpDoesNotExecute(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"-h"}, {"verify", "--help"}, {"verify", "-h"}} {
		f := &effectFake{}
		var out bytes.Buffer
		if err := cli.RunEffectsCommand(args, &out, f); err != nil {
			t.Errorf("help %v: %v", args, err)
		}
		if f.calls != 0 || out.Len() == 0 {
			t.Errorf("help must be printable without execution: %v", args)
		}
	}
}
func TestEffectsCLIReportAndRequest(t *testing.T) {
	for _, fileOutput := range []bool{false, true} {
		t.Run(map[bool]string{false: "stdout", true: "file"}[fileOutput], func(t *testing.T) {
			f := &effectFake{report: domain.EffectReport{SchemaVersion: "effect-run-v1", ContractID: "sort-preserve-v1", Status: "unknown", Reason: "budget exhausted", SourceSHA256: "abc", SourcePath: "fixture.go", Symbol: "Sort", Seed: 1, Cases: []domain.EffectCase{}, Limitations: []string{"finite inputs"}, Execution: map[string]any{"stop_reason": "timeout"}}}
			args := []string{"verify", "--source", "fixture.go", "--symbol", "Sort"}
			p := filepath.Join(t.TempDir(), "receipt.json")
			if fileOutput {
				args = append(args, "--output", p, "--timeout", "30s", "--seed", "18446744073709551615")
			}
			var out bytes.Buffer
			if err := cli.RunEffectsCommand(args, &out, f); err != nil {
				t.Fatal(err)
			}
			if f.calls != 1 || f.request.SourcePath != "fixture.go" || f.request.Symbol != "Sort" {
				t.Fatalf("request/calls: %+v/%d", f.request, f.calls)
			}
			wantTimeout, wantSeed := 5*time.Second, uint64(1)
			if fileOutput {
				wantTimeout = 30 * time.Second
				wantSeed = ^uint64(0)
			}
			if f.request.Timeout != wantTimeout || f.request.Seed != wantSeed {
				t.Fatalf("defaults/options: %+v", f.request)
			}
			payload := out.Bytes()
			if fileOutput {
				var err error
				payload, err = os.ReadFile(p)
				if err != nil {
					t.Fatal(err)
				}
			}
			var got map[string]json.RawMessage
			if err := json.Unmarshal(payload, &got); err != nil {
				t.Fatalf("not one JSON report: %s: %v", payload, err)
			}
			for _, key := range []string{"schema_version", "contract_id", "status", "reason", "source_sha256", "source_path", "symbol", "seed", "cases", "limitations", "execution"} {
				if _, ok := got[key]; !ok {
					t.Errorf("missing public key %q", key)
				}
			}
			var status string
			if err := json.Unmarshal(got["status"], &status); err != nil || status != "unknown" {
				t.Fatalf("unknown lost: %s", payload)
			}
		})
	}
}

type effectFailWriter struct{}

func (effectFailWriter) Write([]byte) (int, error) { return 0, errors.New("output unavailable") }
func TestEffectsCLIErrorsPropagate(t *testing.T) {
	args := []string{"verify", "--source", "fixture.go", "--symbol", "Sort"}
	if err := cli.RunEffectsCommand(args, io.Discard, &effectFake{err: errors.New("source unreadable")}); err == nil {
		t.Fatal("verifier error hidden")
	}
	if err := cli.RunEffectsCommand(args, effectFailWriter{}, &effectFake{report: domain.EffectReport{Status: "unknown"}}); err == nil {
		t.Fatal("writer error hidden")
	}
	args = append(args, "--output", filepath.Join(t.TempDir(), "missing", "receipt.json"))
	if err := cli.RunEffectsCommand(args, io.Discard, &effectFake{report: domain.EffectReport{Status: "unknown"}}); err == nil {
		t.Fatal("output I/O error hidden")
	}
}
