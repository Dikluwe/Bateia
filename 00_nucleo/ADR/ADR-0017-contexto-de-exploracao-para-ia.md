# ADR-0017 — Contexto de exploração para IA

**Status:** ACEITO PARA O RECORTE 0136. **Data:** 2026-10-28. **Relacionados:** ADR-0008, ADR-0010, ADR-0015/0016; passo 0136.

## Contexto

O usuário esclareceu que o Bateia deve auxiliar a IA a sair das abordagens habituais, com menos trabalho manual. Analyze/compare conservam avaliações e observações fornecidas, mas exigem curadoria antes da exploração. Seu catálogo de quatro classes também não delimita todo o repertório que pode ser útil à IA.

## Decisão

Adicionar `idea insights explore` com pergunta direta, abordagens familiares opcionais e bundles coletados. Preparar cartões de leitura automaticamente, com bytes/localizadores e diversidade estrutural de repositórios e tipos. Não exigir classes, alvo, assessments ou observations. A seleção é lexical, limitada e determinística; não classifica novidade, relevância semântica ou comportamento.

Preparar tarefas para alternativas, inversão de premissas, analogias, recombinação e contraexemplos. Elas referenciam os cartões e exigem separar observação de hipótese, condições e experimento. As direções de pesquisa são planos ainda não executados. Os trechos são dados não confiáveis, nunca instruções. Sem corpus ou variedade suficiente, declarar lacunas em vez de inventar exemplos.

## Responsabilidade e consequências

L1 possui integridade/normalização, seleção e preparação de contexto; L2 adapta flags e apresentação. Descoberta/coleta existentes continuam responsáveis por trazer fontes. Não há novo modelo, chave/API, busca remota ou runtime neste comando. Analyze/compare e políticas de reúso permanecem intactos.

A IA pode começar a interpretar exemplos sem um pedido JSON de avaliações preparado à mão. O contexto não equivale a insights gerados por modelo, prova de criatividade ou efeito medido. Familiaridade declarada e matches lexicais são sinais limitados; não há tradução automática ou classificação de domínios. Avaliar utilidade da exploração requer experimentos próprios, além dos testes de integridade/seleção.
