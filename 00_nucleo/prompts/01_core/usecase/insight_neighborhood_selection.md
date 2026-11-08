# insight_neighborhood_selection

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood_selection`. **Materialização:** [insight_neighborhood_selection.go](../../../../01_core/usecase/insight_neighborhood_selection.go).

## Interface e responsabilidade

Resolver relações conservadoras de alias conhecido e projetar InsightNeighborhoodNeighbor a partir da fronteira ordenada. Consumidores são admissão, travessia e relatório de expand.

## Comportamento vigente

Reconhecer X/X.git somente com as duas identidades existentes no mesmo connector/host/URL correspondente. Aplicar a tabela compartilhada de supressões, conflitos e reabertura explícita, sem fundir entradas. Snapshot dos vizinhos antecede leituras: registrar tentativas anteriores, distância opcional, motivo e related_entry_id. Fila e projeção usam o mesmo ordenador; supressões vêm depois da fila.

## Limites e dependências

Não remover .git globalmente, mudar memória/contadores nem executar leituras. Conflito bloqueia automação e gera lacuna; exame prévio não comprova mecanismo. A projeção não altera a política sem distância de discover e não mede novidade/relevância.

As regras compartilhadas estão nas seções [aliases e vizinhos](../../../contratos/exploracao-de-vizinhancas-v1.md#aliases-e-vizinhos), [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [ALS-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ALS-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ORD-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_selection_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_selection.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
