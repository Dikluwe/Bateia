# insight_neighborhood

**Camada:** L1. **Responsável:** `01_core/domain/insight_neighborhood`. **Materialização:** [insight_neighborhood.go](../../../../01_core/domain/insight_neighborhood.go).

## Interface e responsabilidade

Declarar InsightNeighborhoodRequest, InsightNeighborhood e InsightNeighborhoodNeighbor. O pedido interno transporta relatório, memória opcional, foco, AllTrails e overrides; o relatório projeta origem, foco e fronteira para CLI e consumidores de JSON.

## Comportamento vigente

Conservar todos os nomes, tipos e tags existentes. Memory/Focus são opcionais; FollowLinks é ponteiro para distinguir ausência de zero. Focus, SelectionPolicy e Neighbors opcionais preservam a serialização anterior. Distance é ponteiro: valor0 deve ser serializado quando presente. FrontierEntryIDs registra IDs, sem incorporar documentos.

## Limites e dependências

Tipos não validam pedidos nem executam efeitos. AllTrails pertence ao pedido interno; não acrescentar campo público indevido ao relatório. Estados e seleção pertencem aos casos de uso, não aos structs.

As regras compartilhadas estão nas seções [interfaces e formatos](../../../contratos/exploracao-de-vizinhancas-v1.md#interfaces-e-formatos). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [WIRE-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__domain__insight_neighborhood.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
