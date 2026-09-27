/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_scan.md
 * @layer L2
 * @updated 2026-09-26
 */

package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"bateia/01_core/domain"
)

func RunSecurityScanCommand(args []string, stdout io.Writer, runner domain.BinaryVulnerabilityScanner) error {
	fs := flag.NewFlagSet("security scan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	binary := fs.String("binary", "", "Absolute path to an existing Go binary; never executed")
	output := fs.String("output", "", "Optional JSON report")
	timeout := fs.Duration("timeout", time.Minute, "Scanner deadline, at most 5m")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *binary == "" || !filepath.IsAbs(*binary) || fs.NArg() != 0 || *timeout <= 0 || *timeout > 5*time.Minute {
		return fmt.Errorf("security scan requires --binary ABSOLUTE_FILE [--output FILE] [--timeout 1m..5m]")
	}
	inputInfo, err := os.Lstat(*binary)
	if err != nil || !inputInfo.Mode().IsRegular() {
		return fmt.Errorf("--binary must identify a regular file, not a symlink")
	}
	if *output != "" && sameTechnicalOutputPath(*binary, *output) {
		return fmt.Errorf("report output must not overwrite the scanned binary")
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var result domain.VulnerabilityScanReport
	if runner == nil {
		result = domain.VulnerabilityScanReport{SchemaVersion: "1", DecisionImpact: "ADVISORY_ONLY", Binary: *binary,
			Scan: domain.SecurityScanReport{Status: domain.SecurityNotScanned, Scanner: "govulncheck", Scope: "Go binary only", Limitations: []string{"No scanner configured; no safety conclusion can be drawn."}}, Findings: []domain.BinaryVulnerabilityFinding{}, Failure: "scanner unavailable"}
	} else {
		var err error
		result, err = runner.ScanBinary(ctx, *binary)
		if err != nil {
			return err
		}
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	if len(encoded) > 16<<20 {
		return fmt.Errorf("scan report exceeds output byte budget")
	}
	if *output == "" {
		if stdout == nil {
			return fmt.Errorf("stdout unavailable")
		}
		_, err := stdout.Write(encoded)
		return err
	}
	// Atomic replacement prevents partial JSON and does not follow output links.
	if sameTechnicalOutputPath(*binary, *output) {
		return fmt.Errorf("report output must not overwrite the scanned binary")
	}
	return writeAlgorithmHistoryAtomic(*output, encoded, inputInfo)
}
