# Especificacao L0: Symbol Table

- **Caminho**: `01_core/parser/symbol_table.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Objetivo**: Tabela de símbolos de AST organizada por escopos e tipos.

## Contratos
```go
type SymbolTable struct {}
func NewSymbolTable() *SymbolTable
func (st *SymbolTable) Insert(scope string, sym domain.ASTSymbol)
func (st *SymbolTable) LookupScope(scope string) []domain.ASTSymbol
func (st *SymbolTable) LookupKind(kind string) []domain.ASTSymbol
```
