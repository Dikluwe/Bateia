# Especificação L0: Cycle Detector

- **Caminho**: `01_core/domain/cycle_detector.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na análise de convergência de anel da 10ª camada de dependências (`-depth 10`).
- **Objetivo**: Algoritmo DFS sobre o DependencyGraph para detecção de anéis e acoplamentos circulares de dependências.

## Contratos
```go
func DetectCycles(graph *DependencyGraph) [][]string
```
