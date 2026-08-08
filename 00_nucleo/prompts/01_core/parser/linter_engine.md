# Specificacao L0: Linter Engine (Motor de Linteiro Concorrente)

- **Caminho**: `01_core/parser/linter_engine.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Executar regras de linter arquitetural de forma concorrente via worker pools em memória.

## Contratos e Interfaces
```go
type LinterEngine struct {
    workers int
}

func NewLinterEngine(workers int) *LinterEngine
func (e *LinterEngine) RunConcurrent(manifest *domain.Manifest, tree []domain.DirectoryNode) []domain.ArchViolation
```
