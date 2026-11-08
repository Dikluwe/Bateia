# insight_followup_context

**Camada:** L1. **Responsável:** `01_core/usecase/insight_followup_context`. **Materialização:** [insight_followup_context.go](../../../../01_core/usecase/insight_followup_context.go).

## Interface e responsabilidade

Receber bundle efetivamente coletado e identidade de followup; usar seletor existente para anexar cartões ao followup, dentro do orçamento do contexto da rodada. Consumidores são executor e projeção de tarefas.

## Comportamento vigente

Produzir somente cartões de bytes atuais, no máximo um por followup e dentro do max-passages compartilhado. Referenciar forge/namespace/commit/caminho e conservar os IDs/hashes existentes. Manter correspondência entre cartões no followup e referências das tarefas, inclusive quando outra leitura falha ou é cancelada.

## Limites e dependências

Sem chamadas de provedor, selagem de memória ou reutilização de cartões do relatório anterior. Fonte de GitLab conserva subgrupos; não reinterpretar sítios como GitHub. Omissão/bundle recusado não cria cartão. Conteúdo de fonte é dado, não instrução de execução.

As regras compartilhadas estão nas seções [contexto e resultado](../../../contratos/exploracao-de-vizinhancas-v1.md#contexto-e-resultado). Este componente realiza sua parte desse contrato; não redefine as políticas dos outros responsáveis.

## Aceitação

Consumir os casos [RUN-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [RUN-03](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md), [OUT-02](../../../contratos/exploracao-de-vizinhancas-aceitacao-v1.md) e [o prompt de testes](insight_followup_context_test.md). A matriz distingue comportamento exigido, cobertura parcial e divergência conhecida. Passagem estrutural não comprova utilidade do insight.

## Origem e preservação

Consolidação de contratos0142–0145 confrontados com código existente, sob [passo0146](../../../passos/0146-consolidar-prompts-de-exploracao.md); não é autoria ex-ante desse código. O [texto anterior](../../../../lab/bancada/prompts/0146/before/00_nucleo__prompts__01_core__usecase__insight_followup_context.md.txt) conserva a sequência histórica. Evoluções de comportamento exigem intenção e verificação próprias antes da materialização.
