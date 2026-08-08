# Especificacao L0: Memory VFS Store

- **Caminho**: `03_infra/vfs/memory_vfs.go`
- **Camada**: $L_3$ (Infraestrutura)
- **Objetivo**: Abstracao de sistema de arquivos virtual em memoria.

## Contratos
```go
type MemoryVFS struct {}
func NewMemoryVFS() *MemoryVFS
```
