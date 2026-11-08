# insight_neighborhood_test

**Camada:** L2. **Responsável:** aceitação de `02_shell/cli/insight_neighborhood`. **Materialização consumidora:** [insight_neighborhood_test.go](../../../../02_shell/cli/insight_neighborhood_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_neighborhood.md) e as seções [transporte cli](../../../contratos/exploracao-de-vizinhancas-v1.md#transporte-cli). Oráculos desta fatia: [CLI-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [CLI-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [ADM-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Usar factory contada e leitores controlados. Variar flags/JSON/selo/saída-alias e exigir zero factory/publicação em rejeição. Encadear foco herdado, all-trails e troca explícita, com conflito rejeitado antes da factory. Falha operacional com ID publica relatório útil; falha sem ID não publica. Controlar segredo em conteúdo/erro e exigir mensagem sanitizada.

## Limites e dependências

Usar host de teste e portas/factory controladas; não depender de rede pública. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__02_shell__cli__insight_neighborhood_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
