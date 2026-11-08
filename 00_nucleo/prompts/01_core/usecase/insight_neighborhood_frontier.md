# insight_neighborhood_frontier

**Camada:** L1. **Responsável:** `01_core/usecase/insight_neighborhood_frontier`. **Materialização:** [insight_neighborhood_frontier.go](../../../../01_core/usecase/insight_neighborhood_frontier.go).

## Interface e responsabilidade

Receber origens já verificadas, memória, foco normalizado e plano; preparar leads e calcular entradas pendentes alcançáveis com distâncias opcionais para seleção/executor.

## Comportamento vigente

Reutilizar derivação de referências documentais e registrar primeiro testemunho sem duplicar leitura. Percorrer citações conhecidas com visitados; foco tem distância0 e só alcança descendentes válidos. Bloquear descartes e conflitos de alias pela projeção proprietária; conservar namespace GitLab. Sem foco, usar fronteira geral sem inventar distância. Permitir pendência reaberta sem Lead em expand.

## Limites e dependências

Integridade de READMEs pertence a sources, aliases a selection e ordenação a followup_selection. Esta fronteira não revalida documentos nem chama provedor. Respeitar os limites de extração/memória do contrato; links obtidos nesta leitura só serão aprofundados em outra chamada. Primeiro testemunho não representa todos os caminhos alternativos.

As regras compartilhadas estão nas seções [foco e fronteira](../../../contratos/exploracao-de-vizinhancas-v1.md#foco-e-fronteira). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [SRC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_neighborhood_frontier_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_frontier.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
