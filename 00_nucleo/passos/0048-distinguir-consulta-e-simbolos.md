# Passo 0048 — Distinguir palavras da consulta de símbolos explícitos

2026-09-08. Contrato congelado antes da implementação; baseline worktree0047. Regime Tekt A/B sem atestação física.

## Problema e recorte

Na pergunta natural do Bateia, `sem` ganhou IDF por coincidir com identificador de código. Corrigir interpretação da consulta, sem remover identificadores válidos do corpus, alterar pesos0047, adicionar sinônimos ou fingir compreensão da negação. Go local, expansão de vocabulário e interpretação semântica completa ficam fora.

## Contrato observável

Adicionar parser puro InterpretImplementationQuery(query string) []domain.ImplementationQuery. Cada linha não vazia trim/deduplicada/ordenada gera {Original,Mode string; Terms,Symbols []string; Constraints []ImplementationQueryConstraint}. Constraint contém Marker,Text,Status strings. JSON snake_case. Original/Text preservam a linha original sem espaços externos; não parafrasear negação. Status da restrição sempre not_verified.

Uma linha inteira contendo apenas identificador Go válido é consulta symbol, inclusive `sem`, `not` ou `without`. Identificador válido entre backticks balanceados é explícito; identificadores CamelCase/snake_case no texto preservam o critério explícito0047. Linha com código explícito e texto é mixed; restante natural. Symbols contém grafia original dos identificadores explícitos, ordenados/deduplicados. Não inferir chamada/capacidade por isso. Backticks inválidos, não fechados ou conteúdo que não seja um identificador não autorizam ocultar operadores do texto; tratar conteúdo como linguagem natural. Não prometer suporte a seletores como símbolo único.

Em texto natural, palavras isoladas sem/não/nao/nunca/not/without/never (case insensitive) são marcadores de restrição; não entram em Terms, MatchedTerms ou TermWeights nem recebem bônus exato. Registrar marcador e linha integral, sem inferir escopo da negação nem assumir que termo seguinte está proibido/comprovado. Restante da frase continua busca lexical. Operador dentro de identificador explícito (por exemplo CacheWithoutLock) ou `sem` entre backticks permanece código pesquisável. Não mudar tokenização/glossário do código nem frequência documental0047. Terms usa expansão lexical existente, exceto separação desses operadores na consulta; identificador explícito preserva seus tokens.

RankGoImplementations usa Terms e Symbols da interpretação para ranking global e requisitos, com mesma ponderação/corpus/chamadas. ImplementationSearchReport.QueryInterpretation []ImplementationQuery omitempty. DescribeImplementationRequirements usa Terms interpretados e adiciona RequirementRetrieval.Constraints []ImplementationQueryConstraint omitempty, mantendo Status evidence_located como indício lexical e nunca satisfação. Sem fonte ou somente operadores, not_verified. Uma restrição registrada permanece not_verified mesmo com todos termos positivos encontrados. Projeção após revalidação refaz interpretação, para provedores sem suporte ou relatório antigo também. Query original enviada ao provedor/armazenada não é reescrita. Markdown apresenta modo, símbolos e restrições (limites de exibição), deixando explícito que escopo e cumprimento da negação não foram verificados. Versionar contrato/checkpoint/método.

## Aceitação congelada

Autor B congela pares pergunta natural / símbolo explícito antes da candidata: `sem` natural não pontua identificador homônimo; `sem` isolado e entre backticks continuam pesquisáveis; negação preservada como não verificada, inclusive só operador em frase; CamelCase, snake_case, backticks malformados, reordenação e duplicação. Testar projeção e não promoção de restrição. Não ler candidato para escrever casos. APIs antigas/IDF/declarações protegidas continuam sob testes anteriores.

Reexecutar corpus real0046 e holdout0047 nos mesmos snapshots/budgets e controles0045/46/47; exigir manutenção dos ganhos0047 (recall@3>=2/6,@12>=4/6 no real e2/2 holdout; negativas com indícios não aumentam). CLI: consulta natural original não apresenta peso por marcador `sem`; consulta explícita preserva correspondência; replay idêntico. Não exigir melhora artificial em recall nem mudar esperados. Se houver regressão, diagnosticar antes de ajustar; até duas revisões sem ganho por causa. Suite completa, lint, build ./04_wiring; prompts prévios, linhagem e zero IO L1. Não alterar gates/certificados ou commitar.
