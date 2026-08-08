# Especificacao L0: Nil Pointer Detector

- **Caminho**: `01_core/parser/nil_detector.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Análise estática em L1 para detecção de risco de pânico por desreferenciação de ponteiro nulo.

## Contratos
```go
func EvaluateNilPointerRisk(filename string, content string) []domain.ArchViolation
```
