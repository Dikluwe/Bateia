/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/security_scan_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain_test

import (
	"bateia/01_core/domain"
	"encoding/json"
	"testing"
)

func TestSecurityScanReportContract(t *testing.T) {
	rep := domain.SecurityScanReport{
		Status:            domain.SecurityScannedNoMatch,
		Scanner:           "govulncheck",
		Scope:             "BINARY",
		DatabaseUpdatedAt: "2026-09-26T00:00:00Z",
		Limitations:       []string{"heurística; sem garantia de segurança absoluta"},
	}
	data, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	var restored domain.SecurityScanReport
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Status != domain.SecurityScannedNoMatch || restored.Scanner != "govulncheck" {
		t.Fatalf("restored scan report mismatch: %+v", restored)
	}
}
