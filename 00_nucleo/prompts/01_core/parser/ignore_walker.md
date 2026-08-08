# Especificação L0: Ignore Walker

- **Caminho**: `01_core/parser/ignore_walker.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 4 camadas (`-depth 4`) pelo Bateia no repositório [`boyter/gocodewalker`](https://github.com/boyter/gocodewalker).
- **Objetivo**: Filtro de nós de diretório usando padrões estilo `.gitignore` em memória.

## Contratos
```go
type IgnoreWalker struct {}
func NewIgnoreWalker(rules []string) *IgnoreWalker
func (w *IgnoreWalker) ShouldIgnore(path string) bool
```
