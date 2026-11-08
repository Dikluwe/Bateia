# insight_neighborhood_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_neighborhood`. **Materialização consumidora:** [insight_neighborhood_test.go](../../../../01_core/usecase/insight_neighborhood_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood.md) e as seções [propósito e alcance](../../../contratos/exploracao-de-vizinhancas-v1.md#propósito-e-alcance), [admissão](../../../contratos/exploracao-de-vizinhancas-v1.md#admissão), [contexto e resultado](../../../contratos/exploracao-de-vizinhancas-v1.md#contexto-e-resultado). Oráculos desta fatia: [ADM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [OUT-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Usar conectores que falham se houver busca e leitores controlados. Encadear relatórios por três chamadas, verificando alvo efetivo, foco, memória e cartões de cada rodada. Adulterar entrada e exigir zero leituras; cancelar com quota de cartões e verificar resultado parcial rastreável. Reabrir alvo conhecido sem Lead e exigir leitura sem testemunho fabricado.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_neighborhood_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
