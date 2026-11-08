# insight_neighborhood_result

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood_result`. **Materialização:** [insight_neighborhood_result.go](../../../../01_core/usecase/insight_neighborhood_result.go).

## Interface e responsabilidade

Inicializar InsightSiteReport de expand, preparar contexto vazio para a rodada, projetar tarefas e concluir status/erro/ID. Consumido pela fachada e pela publicação CLI.

## Comportamento vigente

Manter coleções de catálogo/coleta antiga vazias, Calls0, origem em neighborhood e memória separada. Tarefas referenciam cartões atuais, memória e leads; catálogo vazio. Aplicar a sequência de estados do contrato, com FAILED/erro quando todas as tentativas falham e CANCELLED/TIMED_OUT prevalecendo. Selar resultado mesmo com erro operacional; nunca promover estado operacional a força de evidência.

## Limites e dependências

Sem provedor, leitura de arquivo, travessia ou mutação de aliases. NoveltyState NOT_ASSESSED e EffectState UNVERIFIED. FRONTIER_EXHAUSTED descreve o escopo conhecido, não inexistência de soluções. Cartões são fornecidos pelo proprietário de contexto, sem reciclagem de histórico.

As regras compartilhadas estão nas seções [contexto e resultado](../../../contratos/exploracao-de-vizinhancas-v1.md#contexto-e-resultado). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [OUT-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [OUT-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_result_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_result.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
