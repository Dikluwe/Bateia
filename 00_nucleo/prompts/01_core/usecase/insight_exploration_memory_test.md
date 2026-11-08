# insight_exploration_memory_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_exploration_memory`. **Materialização consumidora:** [insight_exploration_memory_test.go](../../../../01_core/usecase/insight_exploration_memory_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_exploration_memory.md) e as seções [memória](../../../contratos/exploracao-de-vizinhancas-v1.md#memória). Oráculos desta fatia: [MEM-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [MEM-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [MEM-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [MEM-04](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Validar e modificar a cópia sem alterar original. Adulterar hash/estado/observação/lead/duplicidade, inclusive com resselo externo; exigir rejeição. Descartar e reabrir com motivo, rejeitar EXAMINED declarado e razão vazia. Esgotar256 entradas sem remover história; priorizar candidato novo. No limite, aceitar snapshot999999 para validação, rejeitar seu avanço nos dois estados de feedback e em discover/expand sem efeitos, e aceitar o último avanço999998→999999 com saída que revalida. Comparar toda a entrada antes/depois. Para MEM-02, comparar todos os campos não alterados, incluindo Lead, LastRead e tentativas. Estes critérios antecedem a correção no [passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md).

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_exploration_memory_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
