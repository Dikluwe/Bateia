# Aceitação do plano multifonte

**Camada:** L1. **Consumidor:** [insight_site_plan_test.go](../../../../01_core/usecase/insight_site_plan_test.go).

## Interface e contrato

Consumir [o proprietário](insight_site_plan.md) e os casos ADM-03/MEM-04 da [matriz](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md). Testes usam valores e portas controladas; nenhum acesso à rede ou relógio concreto.

## Comportamento esperado

Pedidos com aliases e IDs equivalentes e bindings permutados devem gerar planos iguais, sem alterar entrada. Controlar fontes ausentes, desconhecidas e repetidas, combinação indevida de all, capacidades duplicadas/inválidas, leitores incompatíveis, locale e orçamento fora dos limites.

Para MEM-04, usar memória selada válida em999998 e999999. Ambas validam como história; somente999998 admite nova rodada, com ou sem tentativas de followup. Discover deve rejeitar999999 antes de buscar ou ler e sem produzir relatório publicável. Na última rodada válida, conferir geração999999 e revalidar todo o snapshot de saída. Comparar entrada antes/depois e não calcular o resultado esperado pela mesma guarda que se quer testar.

## Limites e dependências

A validação de opções não deve consultar capabilities para resolver uma rejeição de geração. Na execução pública, contar buscas, resoluções e leituras. A cobertura dos controles anteriores não substitui limites de bytes/fontes do proprietário sources. Preservar os testes já existentes e acrescentar as regressões de geração.

## Origem e preservação

Consolidação e intenção nova descritas no [passo0147](../../../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md). O [texto anterior](../../../../lab/bancada/prompts/0147/before/00_nucleo__prompts__01_core__usecase__insight_site_plan_test.md.txt) conserva o contrato histórico. Registrar falha antes da correção e passagem depois no Lab.
