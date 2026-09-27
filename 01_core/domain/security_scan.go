/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/security_scan.md
 * @layer L1
 * @updated 2026-09-26
 */
package domain

type SecurityScanStatus string

const (
	SecurityNotScanned         SecurityScanStatus = "NOT_SCANNED"
	SecurityScannedNoMatch     SecurityScanStatus = "SCANNED_NO_MATCH"
	SecurityVulnerabilityFound SecurityScanStatus = "VULNERABILITY_FOUND"
)

// A completed heuristic scan is not a certificate that the software is safe.
type SecurityScanReport struct {
	Status            SecurityScanStatus `json:"status"`
	Scanner           string             `json:"scanner"`
	Scope             string             `json:"scope"`
	DatabaseUpdatedAt string             `json:"database_updated_at,omitempty"`
	Limitations       []string           `json:"limitations"`
}
