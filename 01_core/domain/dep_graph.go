/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/dep_graph.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

import "sync"

// DependencyGraph representa um grafo dirigido de dependências em memória em L1.
type DependencyGraph struct {
	mu    sync.RWMutex
	edges map[string][]string
}

// NewDependencyGraph instancia um novo grafo de dependências.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		edges: make(map[string][]string),
	}
}

// AddEdge adiciona uma aresta dirigida de origem (from) para destino (to).
func (g *DependencyGraph) AddEdge(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, v := range g.edges[from] {
		if v == to {
			return
		}
	}
	g.edges[from] = append(g.edges[from], to)
}

// Reachable realiza uma busca em largura (BFS) para retornar todos os nós alcançáveis a partir da raiz (from).
func (g *DependencyGraph) Reachable(from string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	var queue []string

	queue = append(queue, from)
	visited[from] = true

	var result []string

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr != from {
			result = append(result, curr)
		}

		for _, neighbor := range g.edges[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}
