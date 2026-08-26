# Prompt proprietário — `01_core/parser/ast_fuzzer_test.go`

Origem L0 exclusiva de `01_core/parser/ast_fuzzer_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/parser/ast.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Exercitar de forma determinística a robustez do contrato relacionado, incluindo entradas válidas, malformadas e limites que não podem causar panic. Preservar como observáveis: `FuzzExtractASTSymbols`. Seeds regressivas devem permanecer explícitas e qualquer falha deve ser reproduzível.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/parser/ast.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
