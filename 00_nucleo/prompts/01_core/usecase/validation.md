# Prompt proprietário — `01_core/usecase/validation.go`

Origem L0 exclusiva de `01_core/usecase/validation.go`. Este contrato foi separado de `00_nucleo/prompts/01_core/usecase/funnel.md` para preservar a biunivocidade entre prompt e consumidor.

## Obrigação

Materializar somente a responsabilidade desta unidade, preservando os contratos observáveis: `ValidateQueryFilter`. Manter a direção de dependências da camada e não incorporar responsabilidades pertencentes aos demais consumidores do prompt de origem.

## Restrições

- manter o arquivo na camada e no pacote atuais;
- não relaxar validações para satisfazer testes;
- mudanças futuras neste consumidor devem atualizar primeiro este prompt proprietário;
- consultar `00_nucleo/prompts/01_core/usecase/funnel.md` apenas como contexto de domínio relacionado, não como origem deste consumidor.
