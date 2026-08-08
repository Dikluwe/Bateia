# Especificacao L0: FastCache Store (Cache In-Memory sem GC)

- **Caminho**: `03_infra/cache/fast_cache.go`
- **Camada**: $L_3$ (Infraestrutura)
- **Objetivo**: Implementar CacheStore de alta performance em memória arena sem alocações no heap do GC.

## Contratos
```go
type FastCacheStore struct {}
func NewFastCacheStore() *FastCacheStore
```
