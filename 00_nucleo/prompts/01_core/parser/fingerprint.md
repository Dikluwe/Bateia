# Especificação L0: Code Fingerprint

- **Caminho**: `01_core/parser/fingerprint.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) pelo Bateia nos repositórios [`dualeai/seek`](https://github.com/dualeai/seek) e [`cespare/xxhash/v2`](https://github.com/cespare/xxhash/v2).
- **Objetivo**: Geração de fingerprints de 64-bit em O(1) para comparação de similaridade sintática de código.

## Contratos
```go
func ComputeFingerprint(content string) uint64
func AreSimilar(fp1, fp2 uint64) bool
```
