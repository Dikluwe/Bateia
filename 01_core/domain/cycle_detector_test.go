/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/cycle_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

import "testing"

func TestDetectCycles(t *testing.T) {
	graph := NewDependencyGraph()
	// Ciclo: A -> B -> C -> A
	graph.AddEdge("pkgA", "pkgB")
	graph.AddEdge("pkgB", "pkgC")
	graph.AddEdge("pkgC", "pkgA")

	cycles := DetectCycles(graph)
	if len(cycles) == 0 {
		t.Fatalf("Esperava pelo menos 1 ciclo detectado (pkgA -> pkgB -> pkgC -> pkgA)")
	}

	foundCycle := cycles[0]
	if len(foundCycle) < 3 {
		t.Errorf("Ciclo incompleto: %v", foundCycle)
	}
}
