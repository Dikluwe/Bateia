# insight_neighborhood

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood`. **Materialização:** [insight_neighborhood.go](../../../../01_core/usecase/insight_neighborhood.go).

## Interface e responsabilidade

ExpandInsightNeighborhood recebe contexto, pedido de domínio e bindings; devolve relatório e erro. É a fachada L1 de expansão. ValidateInsightNeighborhood e as origens verificadas pertencem à admissão, não a uma interpretação paralela dentro da fachada.

## Comportamento vigente

Coordenar uma única admissão, plano herdado, deadline por runtime, bootstrap de memória, contexto vazio atual, leads, fronteira, executor compartilhado e conclusão. Reutilizar as origens da admissão; propagar erro em vez de ignorar uma segunda validação. Selecionar a fronteira antes de iniciar leituras. Retornar resultado selado para falhas de leitura e interrupções concluídas pelo executor e preservar o pedido do chamador.

## Limites e dependências

Nenhuma consulta de catálogo, pergunta nova ou I/O concreto. Delegar integridade a sources, alcance a frontier, aliases a selection, quotas a followups e estados/selo a result. Não alterar a jornada discover ao reutilizar o executor.

As regras compartilhadas estão nas seções [propósito e alcance](../../../contratos/exploracao-de-vizinhancas-v1.md#propósito-e-alcance), [admissão](../../../contratos/exploracao-de-vizinhancas-v1.md#admissão), [contexto e resultado](../../../contratos/exploracao-de-vizinhancas-v1.md#contexto-e-resultado). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [ADM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [OUT-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
