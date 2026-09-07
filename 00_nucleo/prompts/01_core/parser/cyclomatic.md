# Especificação L0: estimativa lexical de ramificações Go

- Caminho: `01_core/parser/cyclomatic.go`.
- Camada: L1, função pura e sem I/O.
- Obrigação 0043 F: não apresentar uma contagem lexical como complexidade ciclomática derivada de CFG de McCabe.

## Contrato

Preservar `CalculateCyclomaticComplexity(code string) int` como API legada, delegando ao nome explícito `EstimateGoLexicalComplexity(code string) int`. Publicar `GoLexicalComplexityMethod` para que consumidores possam descrever o método. Atualmente não há consumidor de produção desta função.

A estimativa vale 1 mais a contagem de tokens Go `if`, `for`, `case`, `&&` e `||`, usando `go/scanner`. Comentários, strings interpretadas, strings raw e runas não contribuem. `range` não duplica `for`; `switch` não duplica seus `case`. Else isolado não incrementa a estimativa. Entrada vazia mantém 1 por compatibilidade. Operadores booleanos são contabilizados independentemente de fronteiras de palavras.

Escopo: convenções lexicais Go, sem construir CFG, resolver execução ou validar sintaxe completa. Não inferir suporte a outras linguagens nem equivalência com v(G). O scanner não pode escrever diagnósticos em stderr: usar callback sem I/O. Resultado sobre texto incompleto permanece somente contagem lexical.

Testes próprios desta mudança (não segregados): comentários, literais, controles reais, operadores, range e baseline. Não fornecer atestação independente para estes testes.
