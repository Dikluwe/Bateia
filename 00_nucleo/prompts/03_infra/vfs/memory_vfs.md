# Especificação L0: Memory VFS Store

- **Caminho**: `03_infra/vfs/memory_vfs.go`
- **Camada**: $L_3$ (Infraestrutura)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`hanwen/go-fuse/v2`](https://github.com/hanwen/go-fuse/v2) (`gostream`).
- **Objetivo**: Abstração de sistema de arquivos virtual mantido em memória em $L_3$.

## Contratos
```go
type MemoryVFS struct {}
func NewMemoryVFS() *MemoryVFS
```
