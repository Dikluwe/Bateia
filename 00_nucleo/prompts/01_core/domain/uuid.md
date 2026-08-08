# Especificação L0: Deterministic UUID v5 Generator

- **Caminho**: `01_core/domain/uuid.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 10 camadas (`-depth 10`) no repositório [`google/uuid`](https://github.com/google/uuid).
- **Objetivo**: Geração de UUID v5 determinístico e idempotente por repo@commit para garantias de chave primária universal.

## Contratos
```go
func GenerateRepoUUID(repoKey, commitSHA string) string
```
