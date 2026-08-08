# Especificação L0: JSONPath Metadata Extractor

- **Caminho**: `01_core/parser/jsonpath_extractor.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) no repositório [`PaesslerAG/jsonpath`](https://github.com/PaesslerAG/jsonpath) (via [`bradleyjkemp/sigma-go`](https://github.com/bradleyjkemp/sigma-go)).
- **Objetivo**: Extrator declarativo de campos e metadados de manifestos JSON/YAML em L1.

## Contratos
```go
func QueryJSONPath(jsonContent string, path string) (string, bool)
```
