# Especificação L0: Cyclomatic & Cognitive Complexity Calculator

- **Caminho**: `01_core/parser/cyclomatic.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na análise de dívida técnica e métricas de qualidade de código na mineração do Bateia.
- **Objetivo**: Cálculo estático da Complexidade Ciclomática (v(G) de McCabe) de blocos de código em L1.

## Contratos
```go
func CalculateCyclomaticComplexity(code string) int
```
