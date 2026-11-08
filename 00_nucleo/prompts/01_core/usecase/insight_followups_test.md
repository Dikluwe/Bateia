# insight_followups_test

**Camada:** L1. **Responsável:** aceitação de `01_core/usecase/insight_followups`. **Materialização consumidora:** [insight_followups_test.go](../../../../01_core/usecase/insight_followups_test.go).

## Interface e contrato

Consumir [o prompt proprietário](insight_followups.md) e as seções [ordenação e execução](../../../contratos/exploracao-de-vizinhancas-v1.md#ordenação-e-execução), [memória](../../../contratos/exploracao-de-vizinhancas-v1.md#memória). Oráculos desta fatia: [RUN-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [FOC-01](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Nomes dos testes atuais localizam cobertura; não definem a regra de negócio.

## Comportamento esperado

Executar duas rodadas com leitores em memória e verificar avanço sem repetir examinados. Descartar/reabrir; falhar e retentar; omitir fonte/leitor e contar efeitos. Combinar coleta inicial com followups para esgotar orçamento agregado, inclusive retorno excedente após fetch. Cancelar durante leitura e preservar apenas conteúdo coletado. Testar omissão antes do efeito e diferença entre expand sem Lead e discover com trilha.

## Limites e dependências

Testes do core usam apenas dados e portas controladas, sem I/O concreto. Declarar expectativas antes de executar a ação, contar efeitos quando a garantia for rejeição antecipada e comparar o estado original quando houver risco de alias. A [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) distingue assertivas existentes de controles ainda pendentes; esta consolidação não declara os últimos implementados.

## Origem e preservação

Consolidação documental posterior ao código, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followups_test.md.txt) preserva a história0142–0145. A preservação de nomes/corpos exigida na refatoração0145 foi um critério daquela mudança, não substitui estes oráculos para evoluções futuras.
