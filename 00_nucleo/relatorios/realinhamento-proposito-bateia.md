# Realinhamento ao propósito do Bateia

Correção de direção solicitada pelo usuário em 2026-09-07. Este documento registra prioridades de produto; não declara novos recursos implementados nem invalida evidências históricas.

## Propósito preservado
D-001: transformar uma ideia em decisão fundamentada de reúso, evitando começar do zero quando existe base aproveitável. D-002: descobrir alternativas a um software ou necessidade e compará-las por capacidades e restrições comprovadas. A saída deve permitir escolher um ponto de partida, entender lacunas, bloqueios e incertezas.

Fontes: [D-001](../desejos/D-001-da-ideia-a-decisao-de-reuso.md), [D-002](../desejos/D-002-alternativas-de-software-com-evidencia.md), [passo 0001](../passos/0001-materializar-ideaspec-e-plano-de-expansao.md), [passo 0023](../passos/0023-validar-funil-completo-para-decisoes-reais-de-reuso.md) e [passo 0029](../passos/0029-medir-ganho-da-ia-usando-bateia.md). Estados divergentes entre índices e documentos não foram reconciliados nesta revisão; não se presume que experimentos históricos estejam aprovados.

## Desvio identificado
Os passos 0035–0038 concentraram implementação e avaliação no fluxo de análise local e na amostra Coin/C++. O adaptador Clang é injetado no SearchFunnel em 04_wiring/main.go. A jornada principal usa RecommendationService e IdeaRunDependencies, com expansão, descoberta, planos, bundles e assessment próprios. A existência de evidência AST/CFG no relatório local não comprova que o dossiê de recomendação melhorou.

As melhorias de rastreabilidade e os coletores continuam úteis como infraestrutura opcional. Coin permanece apenas um caso de regressão. Não ampliar especialização em renderização ou profundidade de análise por si só como próxima prioridade.

## Ordem de trabalho daqui em diante
1. Examinar a jornada `idea recommend` e MCP de ponta a ponta com ideias e candidatos de diferentes domínios/ecossistemas. Aproveitar fixtures públicas já existentes; preservar conjuntos ocultos e protocolos de avaliação.
2. Identificar onde a intenção original, capacidades obrigatórias, restrições e não objetivos deixam de orientar descoberta, triagem ou decisão. Expansão de busca nunca redefine o objetivo.
3. Mapear cada capacidade desejada para evidência, contraevidência ou estado desconhecido. Coletores especializados só devem alimentar essa avaliação por contrato explícito e sem converter estrutura em comportamento comprovado.
4. Comparar candidatos e possíveis composições: o que aproveitar, o que adaptar, bloqueios, lacunas e justificativa das alternativas preteridas. BUILD não decorre de falha de busca; INCONCLUSIVE deve dizer qual evidência falta.
5. Medir qualidade de descoberta, afirmações sustentadas, decisões proibidas, clareza da ação seguinte e custo da ferramenta. Repetir casos variados antes de alegar melhoria geral ou ganho para uma IA consumidora.

## Critério para aceitar a próxima melhoria
Deve corrigir uma perda observada na jornada principal e mostrar ganho em mais de um caso, sem regras vinculadas ao nome de um repositório. Testar também negativos e falhas de coleta. Aumentar profundidade de AST, número de grafos ou tamanho de relatório não é objetivo de produto.

## Alteração desta revisão
A abertura do README foi realinhada ao propósito e distingue infraestrutura de evidência e resultado de produto. Nenhum coletor foi removido, nenhum comportamento de produção foi alterado e nenhum novo experimento foi declarado executado.
