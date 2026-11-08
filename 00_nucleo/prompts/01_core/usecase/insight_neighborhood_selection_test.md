# insight_neighborhood_selection_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_neighborhood_selection`. **Materialização consumidora:** [insight_neighborhood_selection_test.go](../../../../01_core/usecase/insight_neighborhood_selection_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood_selection.md) e as seções [aliases e vizinhos](../../../contratos/exploracao-de-vizinhancas-v1.md#aliases-e-vizinhos), [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Oráculos desta fatia: [ALS-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ALS-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ORD-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Combinar estados de base/alias: pendente, examinado, descartado e reaberto com motivo. Exigir código de supressão correto, ausência de leitura/consumo extra e memória anterior intacta. Variar host/forge e base ausente para impedir fusão indevida. Comparar neighbors anterior à execução com followups posterior; testar distância0 e ausência no modo geral.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_selection_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
