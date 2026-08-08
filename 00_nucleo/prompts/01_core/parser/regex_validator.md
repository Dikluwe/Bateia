# Especificação L0: Regex Validator (ReDoS Prevention)

- **Caminho**: `01_core/parser/regex_validator.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) pelo Bateia nos repositórios [`mintoolkit/mint`](https://github.com/mintoolkit/mint) e [`bmatcuk/doublestar`](https://github.com/bmatcuk/doublestar).
- **Objetivo**: Validador preventivo contra expressões regulares com retrocesso catastrófico (ReDoS).

## Contratos
```go
func IsSafeRegex(pattern string) bool
```
