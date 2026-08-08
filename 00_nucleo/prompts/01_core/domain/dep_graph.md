# Especificação L0: Dependency Graph

- **Caminho**: `01_core/domain/dep_graph.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`TransformerOptimus/SuperCoder`](https://github.com/TransformerOptimus/SuperCoder) (`FalkorDB`).
- **Objetivo**: Grafo de dependências dirigido em memória para análise transitiva de pacotes.

## Contratos
```go
type DependencyGraph struct {}
func NewDependencyGraph() *DependencyGraph
func (g *DependencyGraph) AddEdge(from, to string)
func (g *DependencyGraph) Reachable(from string) []string
```
