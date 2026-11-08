# insight_neighborhood_plan_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_neighborhood_plan`. **Materialização consumidora:** [insight_neighborhood_plan_test.go](../../../../01_core/usecase/insight_neighborhood_plan_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood_plan.md) e as seções [admissão](../../../contratos/exploracao-de-vizinhancas-v1.md#admissão), [foco e fronteira](../../../contratos/exploracao-de-vizinhancas-v1.md#foco-e-fronteira), [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Oráculos desta fatia: [MEM-04](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

MEM-04 ([passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md)): relatório com memória efetiva999999 deve validar como histórico e ser rejeitado para expansão antes de resolução/leitura; manter o relatório e feedback originais intactos. Geração999998 admite uma rodada final com memória999999 revalidável.

Admitir entrada válida e guardar a representação normalizada. Mudar conteúdo e estado no original, depois exigir pedido/origens admitidos intactos e leitura do alvo original. Incluir controles de selo, foco desconhecido, conflito AllTrails/foco, limites e feedback que remove identidade/retrocede geração; todos antes das portas. Omitir override herda; ponteiro0 impede leitura.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_plan_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
