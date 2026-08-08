/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/dep_graph.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

import "testing"

func TestDependencyGraph_Reachable(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddEdge("bateia", "linter")
	graph.AddEdge("linter", "ast")
	graph.AddEdge("ast", "parser")

	reachable := graph.Reachable("bateia")
	if len(reachable) != 3 {
		t.Fatalf("Esperava 3 nós alcançáveis, obteve %d", len(reachable))
	}

	expectedMap := map[string]bool{"linter": true, "ast": true, "parser": true}
	for _, node := range reachable {
		if !expectedMap[node] {
			t.Errorf("Nó inesperado no caminho: %s", node)
		}
	}
}
