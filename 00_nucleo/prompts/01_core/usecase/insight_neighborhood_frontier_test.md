# insight_neighborhood_frontier_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_neighborhood_frontier`. **Materialização consumidora:** [insight_neighborhood_frontier_test.go](../../../../01_core/usecase/insight_neighborhood_frontier_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood_frontier.md) e as seções [foco e fronteira](../../../contratos/exploracao-de-vizinhancas-v1.md#foco-e-fronteira). Oráculos desta fatia: [SRC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Construir ciclo e ramo não alcançável a partir de foco conhecido; exigir término, distância mínima e exclusão do ramo externo. Descartar intermediário e verificar bloqueio. Alternar modo geral e foco; conservar identidade GitLab com subgrupo e fonte restrita. Controle de alias não pode liberar caminho conflitante. Ler alvo sem Lead não autoriza fabricar citação.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_frontier_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
