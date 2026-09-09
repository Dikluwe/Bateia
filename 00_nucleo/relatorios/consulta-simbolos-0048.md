# Consulta natural, símbolos e restrições — execução 0048

O [passo 0048](../passos/0048-distinguir-consulta-e-simbolos.md) foi implementado em 2026-09-08. A recuperação local de implementações Go distingue linguagem natural, símbolos explícitos e consultas mistas. Negação identificada na pergunta fica registrada como restrição `not_verified`, sem contribuir como palavra positiva para o ranking. Isto corrige uma ambiguidade lexical observada; não implementa compreensão semântica geral.

## Implementação

`InterpretImplementationQuery` preserva cada linha original, termos usados, símbolos e marcadores de restrição. Um identificador Go isolado, ou entre backticks simples válidos, continua pesquisável, inclusive `sem`, `s` e `use`. Nomes com CamelCase e underscores preservam a heurística anterior de símbolo explícito. Em texto natural, `sem`, `não`, `nao`, `nunca`, `not`, `without` e `never` viram restrições não verificadas. O texto integral acompanha cada marcador: o analisador não inventa o escopo da negação.

O ranking usa essa interpretação; a projeção por requisito e a revalidação a reconstruem da consulta original, inclusive para relatórios antigos. JSON e Markdown apresentam as restrições. Um resumo anterior ao limite de exibição evita ocultar a existência de uma restrição em uma linha posterior. O método passa a `go_ast_query_v1` e a versão de evidência a `evidence-v10-query-interpretation`.

## Evidência observada

Os [controles da CLI](../benchmark/0048/cli-observations.json) usam uma fixture sintética do integrador, separada dos casos do avaliador B. Em `buscar fila sem cache`, a baseline 0047 atribuía peso ao identificador `sem`; a candidata deixa de pontuá-lo e registra a restrição. `sem` isolado e explicitado entre backticks continuam localizando sua definição. Delimitadores inválidos não ocultam a negação. Uma consulta composta somente por marcadores naturais não produz hits nesses controles.

Na [autoanálise do Bateia](../benchmark/0048/bateia.md), a consulta `como a busca preserva requisitos originais e rejeita capacidades sem evidência` deixou de usar `sem` como termo positivo na recuperação de implementações. A cobertura permanece insuficiente e o cumprimento da restrição não foi demonstrado. A [decisão local](../benchmark/0048/decision.md) conserva as dez capacidades como `INCONCLUSIVE`. Os cinco artefatos comparados entre execução local e replay são idênticos byte a byte; hashes e snapshots estão no recibo da CLI.

## Regressões e verificação

A [comparação final B](../benchmark/0048/evaluation/final-comparison.json) preserva as métricas da baseline 0047:

| Conjunto | Arquivos-alvo recuperados @3 | @12 | Requisitos negativos com indícios |
| --- | --- | --- | --- |
| Real 0046 | 2/6 | 4/6 | 1/3 |
| Holdout 0047 | 2/2 | 2/2 | 1/1 |

Esses limites contam caminhos distintos nos primeiros 24 hits; não equivalem aos arquivos efetivamente extraídos. Não houve ganho de recall demonstrado neste passo, nem eliminação dos falsos indícios preexistentes. Os controles 0045, 0046 e 0047 passaram.

Build de `./04_wiring`, `go test ./...` e `crystalline-lint .` passaram na revisão final. A suíte completa usou permissão para servidores HTTP locais dos testes. O [recibo de verificação](../benchmark/0048/verification.json) registra comandos, hashes e logs. O corpus, sua ponderação e o mecanismo de relações funcionais não foram alterados.

## Autoria e limites

Foi aplicado Tekt A/B, sem atestação de isolamento físico: A implementou o interpretador e a integração no ranking; B escreveu e executou os casos independentes; o integrador modificou a projeção, os relatórios e a CLI. O [veredito B](../benchmark/0048/evaluation/verdict.md) registra congelamentos e revisões. A revisão r1 substituiu um import proibido pelo linteiro sem alterar expectativas. A extensão r2 acrescentou os símbolos `s` e `use` após um contraexemplo público; foi congelada antes da execução, mas a correção já estava pronta, portanto não há execução falha registrada desse caso contra a candidata inicial. A baseline não possui a API nova: sua falha de compilação não demonstra força contra mutações semânticas.

O catálogo de marcadores é finito. Não há resolução de escopo, contradição, ausência de efeitos, fluxo de dados ou prova comportamental. O critério de identificador isolado também pode interpretar uma palavra natural única como símbolo; essa é uma convenção explícita, não uma inferência de intenção garantida. A pontuação de descoberta por README e as consultas enviadas aos provedores não foram modificadas. As medições de duração nos artefatos são telemetria do Bateia, não benchmarks dos projetos analisados.
