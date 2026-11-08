# insight_followup_context_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_followup_context`. **Materialização consumidora:** [insight_followup_context_test.go](../../../../01_core/usecase/insight_followup_context_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_followup_context.md) e as seções [contexto e resultado](../../../contratos/exploracao-de-vizinhancas-v1.md#contexto-e-resultado). Oráculos desta fatia: [RUN-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [OUT-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Limitar max-passages com várias coletas e comparar número total e refs por destino. Cancelar após uma leitura válida e exigir cartão preservado apenas para a coleta atual. Usar identidade de outro sítio e verificar namespace/commit/path. Controle sem bundle e controle com cartão histórico devem produzir zero cartão novo; a matriz distingue cobertura existente e controles a ampliar.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followup_context_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
