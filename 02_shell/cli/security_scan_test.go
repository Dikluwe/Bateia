/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_scan_test.md
 * @layer L2
 * @updated 2026-09-26
 */

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

type fakeSecurityScanner struct {
	called int
	binary string
}

func (s *fakeSecurityScanner) ScanBinary(ctx context.Context, binary string) (domain.VulnerabilityScanReport, error) {
	s.called++
	s.binary = binary
	return domain.VulnerabilityScanReport{SchemaVersion: "1", DecisionImpact: "ADVISORY_ONLY", Binary: binary, Scan: domain.SecurityScanReport{Status: domain.SecurityScannedNoMatch}, Complete: true, Findings: []domain.BinaryVulnerabilityFinding{}}, nil
}

func TestSecurityScanCLIStdoutAndAtomicFile(t *testing.T) {
	input := filepath.Join(t.TempDir(), "binary")
	if err := os.WriteFile(input, []byte("artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	runner := &fakeSecurityScanner{}
	var out bytes.Buffer
	if err := RunSecurityScanCommand([]string{"--binary", input}, &out, runner); err != nil {
		t.Fatal(err)
	}
	var result domain.VulnerabilityScanReport
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Scan.Status != domain.SecurityScannedNoMatch || runner.binary != input || runner.called != 1 {
		t.Fatal("unexpected CLI result")
	}
	output := filepath.Join(t.TempDir(), "result.json")
	out.Reset()
	if err := RunSecurityScanCommand([]string{"--binary", input, "--output", output}, &out, runner); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatal("file output should not duplicate stdout")
	}
	data, err := os.ReadFile(output)
	if err != nil || json.Unmarshal(data, &result) != nil {
		t.Fatal("invalid report file")
	}
	data, _ = os.ReadFile(input)
	if string(data) != "artifact" {
		t.Fatal("scanner CLI changed binary")
	}
}

func TestSecurityScanCLIRejectsOutputAliasesAndInvalidFlags(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "binary")
	if err := os.WriteFile(input, []byte("artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	link, hardlink := filepath.Join(dir, "symlink"), filepath.Join(dir, "hardlink")
	if err := os.Symlink(input, link); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(input, hardlink); err != nil {
		t.Fatal(err)
	}
	runner := &fakeSecurityScanner{}
	for _, args := range [][]string{
		nil, {"--binary", "relative"}, {"--binary", input, "extra"},
		{"--binary", input, "--timeout", "0s"}, {"--binary", input, "--timeout", "6m"},
		{"--binary", input, "--source", "."},
		{"--binary", input, "--output", input}, {"--binary", input, "--output", link}, {"--binary", input, "--output", hardlink},
	} {
		if err := RunSecurityScanCommand(args, &bytes.Buffer{}, runner); err == nil {
			t.Errorf("invalid args accepted: %q", args)
		}
	}
	if runner.called != 0 {
		t.Fatal("invalid args reached scanner")
	}
}

func TestSecurityScanCLINilScannerIsNotCleanScan(t *testing.T) {
	var out bytes.Buffer
	input := filepath.Join(t.TempDir(), "binary")
	if err := os.WriteFile(input, []byte("artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := RunSecurityScanCommand([]string{"--binary", input}, &out, nil); err != nil {
		t.Fatal(err)
	}
	var result domain.VulnerabilityScanReport
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Complete || result.Scan.Status != domain.SecurityNotScanned || result.Failure == "" || result.DecisionImpact != "ADVISORY_ONLY" {
		t.Fatal("nil scanner falsely claims coverage")
	}
}
