# Prompt Spec L0: Vocabulário de Ideia e Plano de Expansão

## Objetivo

Definir em L1 os tipos puros que separam a intenção do usuário das consultas de
descoberta: `IdeaSpec`, `ExpansionFacet`, `SearchExpansion` e `ExpansionPlan`.

## Restrições

- Camada L1, pacote `domain`.
- Zero I/O e zero dependências externas.
- Estruturas serializáveis em JSON.
- Problema, capacidades, restrições, não objetivos e stack permanecem dimensões
  distintas.
- A representação não executa busca nem contém detalhes de Sourcegraph/GitHub.

## Critérios de Verificação

- A forma pública é coberta por testes de serialização.
- `NonGoals` e `PreferredStack` não são confundidos com capacidades.
- Os tipos podem ser consumidos por L1 e L2 sem importar infraestrutura.
