# Especificação L0: Code Clone Detector

- **Caminho**: `01_core/parser/clone_detector.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas pelo Bateia (combinação de Fingerprints xxhash em L1).
- **Objetivo**: Comparador de similaridade de código para detecção de plágio e funções duplicadas em O(1).

## Contratos
```go
func CalculateSimilarityScore(fp1, fp2 uint64) float64
```
