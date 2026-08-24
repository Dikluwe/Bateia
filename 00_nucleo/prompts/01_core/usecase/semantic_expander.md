# Prompt Spec L0: Expansão Semântica e Avaliação

## Objetivo

Enriquecer um `ExpansionPlan` por uma estratégia injetada e comparar estratégias
num corpus versionado sem realizar I/O em L1.

## Restrições

- Camada L1, pacote `usecase`.
- A estratégia é uma porta; L1 não conhece vocabulário, arquivo ou rede.
- Expansões semânticas possuem estratégia, regras e expansões de origem.
- Preservar expansões literais, ideia-alvo, ordem e determinismo.
- Deduplicar queries e mesclar cobertura e proveniência.
- `NonGoals` bloqueiam promoções positivas.
- Avaliar cobertura, contaminação, novidade e duplicação sem divisão por zero.

## Critérios de Verificação

- A entrada não é modificada.
- Expansões sem proveniência são rejeitadas.
- Literal e lexical podem ser comparados pelo mesmo avaliador.
- Duas avaliações equivalentes produzem JSON idêntico.
