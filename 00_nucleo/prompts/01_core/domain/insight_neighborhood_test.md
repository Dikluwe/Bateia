# insight_neighborhood_test

**Camada:** L1. **Responsável:** aceitação de `01_core/domain/insight_neighborhood`. **Materialização consumidora:** [insight_neighborhood_test.go](../../../../01_core/domain/insight_neighborhood_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood.md) e as seções [interfaces e formatos](../../../contratos/exploracao-de-vizinhancas-v1.md#interfaces-e-formatos). Oráculos desta fatia: [WIRE-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Comparar bytes de relatório legado sem neighborhood; serializar expansão com origem/foco/fronteira e decodificar com campos estritos. Comparar ausência de neighbors/selection_policy com presença; distância0 deve permanecer distinta de nil. Verificar roundtrip de related_entry_id e tentativas anteriores sem interpretar esses campos como evidência.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__domain__insight_neighborhood_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
