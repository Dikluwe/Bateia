# insight_neighborhood_plan

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood_plan`. **Materialização:** [insight_neighborhood_plan.go](../../../../01_core/usecase/insight_neighborhood_plan.go).

## Interface e responsabilidade

ValidateInsightNeighborhood preserva sua API pública e projeta a Request da admissão interna. admitInsightNeighborhood devolve pedido normalizado e Origins verificadas, consumidos pela fachada.

## Comportamento vigente

Copiar/validar relatório e memória original; obter origens congeladas antes de aplicar feedback. Feedback conserva identidades e não retrocede geração. Resolver foco explícito/herdado e AllTrails, inclusive aliases descartados. Rejeitar geração efetiva>=999999 usando a guarda pertencente ao container de memória, compartilhada com discover/feedback ([passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md)). Herdar opções e aplicar somente overrides permitidos, revalidando o plano; FollowLinks0 difere de ausência.

## Limites e dependências

Admissão não chama leitores nem relógio/rede concretos. Os limites e classes de erro são os do contrato compartilhado. Cópia profunda inclui dados de fontes e memória; fontes admitidas não podem apontar a armazenamento mutável do chamador. Integridade documental é delegada ao proprietário sources.

As regras compartilhadas estão nas seções [admissão](../../../contratos/exploracao-de-vizinhancas-v1.md#admissão), [foco e fronteira](../../../contratos/exploracao-de-vizinhancas-v1.md#foco-e-fronteira), [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [ADM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_plan_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_plan.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.

## 0159 — Guarda anterior de recibos

Delegar ao responsável sources a verificação um-para-um dos recibos primários antes de admitir as origens. A guarda nova não possui fontes, efeitos ou política de memória. Preservar o corpo histórico de insightNeighborhoodOrigins e o gate0145; o [passo0159](../../../passos/0159-fechar-cancelamento-escopo-e-recibos.md) acrescenta a obrigação sem reescrever o congelamento0145.
