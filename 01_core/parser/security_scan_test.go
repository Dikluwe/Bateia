/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security_scan_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package parser

import (
	"bateia/01_core/domain"
	"testing"
)

func TestSecurityScanDistinguishesNoAnalysisFromNoMatch(t *testing.T) {
	for _, tc := range []struct {
		manifest *domain.Manifest
		status   domain.SecurityScanStatus
	}{
		{nil, domain.SecurityNotScanned},
		{&domain.Manifest{}, domain.SecurityScannedNoMatch},
		{&domain.Manifest{Dependencies: []domain.Dependency{{Name: "unknown-package", Version: "1.0.0"}}}, domain.SecurityScannedNoMatch},
		{&domain.Manifest{Dependencies: []domain.Dependency{{Name: "github.com/dgrijalva/jwt-go", Version: "v3.2.0"}}}, domain.SecurityVulnerabilityFound},
	} {
		report, advisories := EvaluateSecurityScan(tc.manifest)
		if report.Status != tc.status || report.Scanner == "" || report.DatabaseUpdatedAt == "" || len(report.Limitations) == 0 {
			t.Fatalf("dishonest scope: %+v", report)
		}
		if (len(advisories) > 0) != (tc.status == domain.SecurityVulnerabilityFound) {
			t.Fatal("status does not match findings")
		}
	}
}
