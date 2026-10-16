/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0121/cmd/declarations/main_test.md
 * @layer lab
 * @updated 2026-10-07
 */
package main

import (
	concurrency "bateia/lab/bancada/concorrencia/0109"
	declarations "bateia/lab/bancada/concorrencia/0121"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestStrictInputAndExclusiveOutput(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "input.json")
	output := filepath.Join(directory, "output.json")
	code := "package p\n// New doc.\nfunc New(){}\n"
	request := declarations.Request{Source: concurrency.Request{SchemaVersion: "1", Origin: "CLI_FIXTURE", Files: []concurrency.SourceFile{{Path: "source.go", Content: code, SHA256: concurrency.Digest(code)}}}}
	raw, _ := json.Marshal(request)
	if err := os.WriteFile(input, raw, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"--input", input}, &out); err != nil {
		t.Fatal(err)
	}
	var result declarations.Report
	if err := json.Unmarshal(out.Bytes(), &result); err != nil || len(result.Declarations) != 1 {
		t.Fatal(err, result)
	}
	if err := run([]string{"--input", input, "--output", output}, &out); err != nil {
		t.Fatal(err)
	}
	original, _ := os.ReadFile(output)
	if err := run([]string{"--input", input, "--output", output}, &out); err == nil {
		t.Fatal("overwrote output")
	}
	after, _ := os.ReadFile(output)
	if !bytes.Equal(original, after) {
		t.Fatal("output changed")
	}
	info, _ := os.Stat(output)
	if info.Mode().Perm() != 0600 {
		t.Fatal(info.Mode())
	}
	for _, data := range []string{`{"source":{},"source":{}}`, `{"unknown":1}`, string(raw) + ` {}`} {
		os.WriteFile(input, []byte(data), 0600)
		if err := run([]string{"--input", input}, &out); err == nil {
			t.Fatal("invalid JSON accepted")
		}
	}
	for _, args := range [][]string{nil, {"--input", input, "extra"}, {"--unknown"}} {
		if err := run(args, &out); err == nil {
			t.Fatal(args)
		}
	}
}
