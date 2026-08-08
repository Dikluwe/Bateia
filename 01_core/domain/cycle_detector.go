/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/cycle_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

// DetectCycles identifica acoplamentos circulares (anéis de dependência) no DependencyGraph usando DFS.
func DetectCycles(graph *DependencyGraph) [][]string {
	if graph == nil {
		return nil
	}

	graph.mu.RLock()
	defer graph.mu.RUnlock()

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles [][]string

	var dfs func(node string, path []string)
	dfs = func(node string, path []string) {
		visited[node] = true
		recStack[node] = true
		currentPath := append(path, node)

		for _, neighbor := range graph.edges[node] {
			if !visited[neighbor] {
				dfs(neighbor, currentPath)
			} else if recStack[neighbor] {
				// Ciclo encontrado! Extrair o sub-caminho do ciclo
				cycleStartIdx := -1
				for idx, pNode := range currentPath {
					if pNode == neighbor {
						cycleStartIdx = idx
						break
					}
				}
				if cycleStartIdx != -1 {
					cycle := append([]string(nil), currentPath[cycleStartIdx:]...)
					cycle = append(cycle, neighbor)
					cycles = append(cycles, cycle)
				}
			}
		}

		recStack[node] = false
	}

	for node := range graph.edges {
		if !visited[node] {
			dfs(node, nil)
		}
	}

	return cycles
}
