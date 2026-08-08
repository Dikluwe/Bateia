# Especificação L0: FastCache Store (Cache In-Memory sem GC)

- **Caminho**: `03_infra/cache/fast_cache.go`
- **Camada**: $L_3$ (Infraestrutura)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`VictoriaMetrics/fastcache`](https://github.com/VictoriaMetrics/fastcache).
- **Objetivo**: Implementar `CacheStore` de alta performance em memória arena para zerar a pressão de alocação no heap do GC.

## Contratos
```go
type FastCacheStore struct {}
func NewFastCacheStore() *FastCacheStore
```
