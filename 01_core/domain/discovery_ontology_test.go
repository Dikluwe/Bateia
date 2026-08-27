/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/discovery_ontology_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestQueryQualityEnumIsStable(t *testing.T) {
	if QueryAccepted != "ACCEPTED" || QueryWeak != "WEAK" || QueryRejected != "REJECTED" {
		t.Fatal("query quality contract changed")
	}
}
