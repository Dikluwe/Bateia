# Prompt proprietário — `01_core/usecase/probabilistic_screening_test.go`

Origem L0 exclusiva de `01_core/usecase/probabilistic_screening_test.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/probabilistic_screening.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Verificar de forma determinística o comportamento público do screening, cobrindo
ranking, migração, diagnóstico e a autoridade de calibração. Modelo não calibrado não
pode descartar; calibração válida preserva o descarte; estado contraditório ou não
finito deve resultar em `AUDIT` sem probabilidade. Os testes não dependem de rede,
relógio ou ordem acidental.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/probabilistic_screening.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
