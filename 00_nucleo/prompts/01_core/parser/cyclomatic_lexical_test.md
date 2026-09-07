# Prompt proprietário — `01_core/parser/cyclomatic_lexical_test.go`

Camada L1. Testes unitários de autoria do implementador, sem alegação de independência. Verificar o contrato lexical Go de `cyclomatic.md`: comentários e literais não alteram a estimativa; controles e operadores reais contribuem; for-range não duplica a contagem; entrada vazia preserva baseline; API legada delega ao método explícito. Não acessar filesystem, rede ou tempo.
