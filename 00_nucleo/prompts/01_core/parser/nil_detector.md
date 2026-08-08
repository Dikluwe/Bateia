# Especificação L0: Nil Pointer Detector

- **Caminho**: `01_core/parser/nil_detector.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 6 camadas (`-depth 6`) pelo Bateia no repositório [`uber-go/nilaway`](https://github.com/uber-go/nilaway).
- **Objetivo**: Análise estática pura em L1 para detecção da regra `NIL_POINTER_RISK` (risco de pânico por desreferenciação de ponteiro nulo sem checagem prévia).

## Contratos
```go
func EvaluateNilPointerRisk(filename string, content string) []domain.ArchViolation
```
