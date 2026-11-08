# insight_followup_selection

**Camada:** L1. **Responsável:** `01_core/usecase/insight_followup_selection`. **Materialização:** [insight_followup_selection.go](../../../../01_core/usecase/insight_followup_selection.go).

## Interface e responsabilidade

Ordenar fila de entradas pendentes recebida, usando plano/bindings, ordem das referências e distâncias opcionais. Mesma política atende executor discover e expand.

## Comportamento vigente

Comparar nesta ordem: elegibilidade por fonte selecionada e leitor; distância quando fornecida; tentativas crescentes; FirstGeneration crescente; presença/posição na ordem atual de referências; ID. Ordenar somente a fila, mantendo IDs e memória. Sem distâncias, conservar exatamente a política histórica de discover.

## Limites e dependências

Nenhum efeito, contador incrementado, resolução de alias ou julgamento de relevância. A disponibilidade é inferida da porta configurada; não consultar o provedor para saber se está acessível. O executor possui quotas e a projeção possui motivos/snapshot.

As regras compartilhadas estão nas seções [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [ORD-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ORD-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_followup_selection_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followup_selection.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
