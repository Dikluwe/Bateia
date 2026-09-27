# Prompt proprietário — `01_core/parser/adr_generator_test.go`

Origem L0 exclusiva de `01_core/parser/adr_generator_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/parser/adr_generator.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do contrato relacionado, cobrindo sucesso, rejeição e limites já materializados nesta unidade. Preservar como observáveis: `TestGenerateADRDocument`. Os testes não podem depender de rede real, relógio não controlado ou ordem acidental.

Verificar status `PROPOSTO` tanto com candidatos quanto sem resultados. Preservar consulta, data fixa, contagens do funil, score, justificativa, URL, referência de commit e dependências diretas observadas (no máximo cinco). Nenhuma saída pode afirmar aceitação, testes em produção ou redução comprovada de dívida técnica. Resultados vazios não podem sugerir adoção, benefícios de uma seleção inexistente nem conclusão de rejeição de todo o universo de alternativas. Manifestos só com dependências indiretas não devem produzir uma seção vazia de dependências observadas.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/parser/adr_generator.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.

Cobrir referências de commit vazias e `HEAD`: o relatório deve declarar ausência de commit fixado, preservando a referência de repositório se disponível, sem rotulá-la como commit.
