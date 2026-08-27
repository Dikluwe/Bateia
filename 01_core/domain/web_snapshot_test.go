/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/web_snapshot_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestWebSnapshotTypes(t *testing.T) {
	s := WebSearchSnapshot{SchemaVersion: "1", Arm: "F2"}
	if s.Arm != "F2" {
		t.Errorf("expected F2, got %s", s.Arm)
	}
}
