# insight_exploration_memory_test

**Camada:** L1. **Responsável:** aceitação de `01_core/domain/insight_exploration_memory`. **Materialização consumidora:** [insight_exploration_memory_test.go](../../../../01_core/domain/insight_exploration_memory_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_exploration_memory.md) e as seções [interfaces e formatos](../../../contratos/exploracao-de-vizinhancas-v1.md#interfaces-e-formatos), [memória](../../../contratos/exploracao-de-vizinhancas-v1.md#memória). Oráculos desta fatia: [WIRE-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Construir relatório com entrada descartada/motivo e followup coletado/cartão; serializar e comparar roundtrip integral. Controlar campos opcionais ausentes e observação/lead presentes. Não usar uma entrada de memória para satisfazer a expectativa de bundle atual; validação semântica tem consumidor próprio.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__domain__insight_exploration_memory_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
