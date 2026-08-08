# Especificação L0: Symbol Table

- **Caminho**: `01_core/parser/symbol_table.go`
- **Camada**: $L_1$ (Núcleo Puro - ZERO I/O)
- **Origem & Rastreabilidade**: Descoberto na varredura de 6 camadas (`-depth 6`) pelo Bateia no repositório [`pingcap/tidb`](https://github.com/pingcap/tidb).
- **Objetivo**: Tabela de Símbolos de AST organizada hierarquicamente por escopos e tipos para acelerar a extração do Estágio 3.

## Contratos
```go
type SymbolTable struct {}
func NewSymbolTable() *SymbolTable
func (st *SymbolTable) Insert(scope string, sym domain.ASTSymbol)
func (st *SymbolTable) LookupScope(scope string) []domain.ASTSymbol
func (st *SymbolTable) LookupKind(kind string) []domain.ASTSymbol
```
