# Especificação L0: Coupling & Instability Calculator

- **Caminho**: `01_core/parser/coupling.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) nos repositórios [`openshift/telemeter`](https://github.com/openshift/telemeter) e [`efficientgo/core`](https://github.com/efficientgo/core).
- **Objetivo**: Cálculo do Acoplamento Aferente (Ca), Eferente (Ce) e Índice de Instabilidade (I) de Robert C. Martin em L1.

## Contratos
```go
type CouplingMetrics struct {
    Afferent  int     `json:"afferent"`
    Efferent  int     `json:"efferent"`
    Instability float64 `json:"instability"`
}

func CalculateCoupling(afferent, efferent int) CouplingMetrics
```
