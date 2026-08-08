# Especificacao L0: Ignore Walker

- **Caminho**: `01_core/parser/ignore_walker.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Filtro de nos de diretório estilo .gitignore em memoria.

## Contratos
```go
type IgnoreWalker struct {}
func NewIgnoreWalker(rules []string) *IgnoreWalker
func (w *IgnoreWalker) ShouldIgnore(path string) bool
```
