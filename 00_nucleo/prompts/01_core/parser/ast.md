# Specification: AST Parser (L1 Core)

## Purpose
Extract structural AST symbols (functions, structs, interfaces, types) from in-memory code files without performing any file or network I/O.

## Rules
- Pure L1 package (`01_core/parser`).
- ZERO I/O (no `os.ReadFile`, no network requests).
- Accepts filename and `[]byte` source content.
- Uses `go/ast` and `go/parser` for Go code.
- Returns a list of `domain.ASTSymbol`.
