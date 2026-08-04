# Especificação L0 — Analisador Sintático de AST (01_core/parser/ast_analyzer.go)

## Responsabilidade
Implementar funções puras em Go em L1 para analisar código Go e extrair símbolos estruturais refinados (structs, interfaces, métodos, funções exportadas e comentários de documentação) utilizando o pacote nativo `go/ast` e `go/parser`.

## Restrições Tekt
- **Zero I/O**: Operações estritamente em memória a partir de `string` ou `[]byte`.
- **Camada**: L1 (Núcleo Puro).
- Sem dependências de L2, L3 ou L4.
