# insight_followup_selection_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_followup_selection`. **Materialização consumidora:** [insight_followup_selection_test.go](../../../../01_core/usecase/insight_followup_selection_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_followup_selection.md) e as seções [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução). Oráculos desta fatia: [ORD-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ORD-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Colocar foco sem leitor e destinos elegíveis perto/longe com contadores distintos: elegibilidade vence distância, distância vence tentativas. Repetir com distância ausente e exigir prioridade por tentativas. Comparar memória/contadores antes/depois e zero leituras. Para empates, variar uma chave por vez e permutar a entrada; No [passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md), completar ORD-02 com casos que isolem geração, presença de referência, posição de referência e ID, inclusive distância ausente e empatada. Permutar a fila em cada caso e exigir ordem esperada explícita, sem alterar memória/bindings ou chamar leitores.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followup_selection_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
