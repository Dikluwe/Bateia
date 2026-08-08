# Especificacao L0: Dependency Graph

- **Caminho**: `01_core/domain/dep_graph.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Grafo de dependências dirigido em memória para análise transitiva.

## Contratos
```go
type DependencyGraph struct {}
func NewDependencyGraph() *DependencyGraph
func (g *DependencyGraph) AddEdge(from, to string)
func (g *DependencyGraph) Reachable(from string) []string
```
