# Execução 0047 — Pontuação pelo corpus, inspirada em busca com o Bateia

2026-09-07 (local). [Passo congelado](../passos/0047-ponderar-termos-pelo-corpus.md). **Implementado e promovido após ganho no corpus real e preservação dos controles declarados.** O ganho é de recuperação de arquivos na amostra, não prova de compreensão semântica geral.

## Como o Bateia ajudou a escolher

Usamos o binário0046 para `search -query 'BM25 code search' -lang go -max 3`. A primeira tentativa teve bloqueio de rede no sandbox; a tentativa autorizada completou. A consulta original no Sourcegraph retornou zero candidatos; o fallback explícito `lang:Go BM25 code` retornou `knowns-dev/knowns` e `vllm-project/semantic-router`. A extração trouxe README/manifestos, sem implementação de ranking suficiente para adotar como solução. Scores43 e0 são indícios de busca, não validação.

A análise local do Blevev2.6.0 com `BM25 idf scoring terms query` foi mais útil: localizou `newTermSearcherFromReader` e `docScore`. Inspeção confirmou frequência documental e normalização de comprimento em [scorer_term.go](https://github.com/blevesearch/bleve/blob/d8f2ab9a11166223bc4997143efda40ec98045e7/search/scorer/scorer_term.go) e ligação do scorer em [search_term.go](https://github.com/blevesearch/bleve/blob/d8f2ab9a11166223bc4997143efda40ec98045e7/search/searcher/search_term.go). Essa observação inspirou o experimento, sem copiar código nem instalar dependência.

[Consultas e resultados públicos](../../lab/bancada/history/benchmarks/0047/ideas/public-search-observations.json), [localizações fornecidas pelo Bateia](../../lab/bancada/history/benchmarks/0047/ideas/bleve-observations.json) e [arquivos/intervalos inspecionados](../../lab/bancada/history/benchmarks/0047/ideas/bleve-inspection.json). O fallback não foi confundido com sucesso da consulta original.

## Implementação

O parser faz uma coleta AST e calcula frequência dos termos nas declarações elegíveis, incluindo as que não correspondem à consulta. Nome e corpo têm contribuições separadas; corpo considera identificadores únicos, com normalização por comprimento. Comentários e literais isolados continuam fora da pontuação. Glossário e critérios de casamento lexical não foram ampliados para favorecer os benchmarks.

A heurística usa IDF `log1p((N-df+0.5)/(df+0.5))`, contribuições distintas para nome/corpo e saturação na conversão para Score inteiro. Preserva as faixas de prioridade de função pertinente e de identificador exato. `TermWeights` registra termo, campo, df, N e contribuição efetivamente usada; `Reason` registra fórmula, comprimento e média. **N é o recorte de declarações lidas, não o ecossistema.** Não apresentamos a heurística como implementação completa de BM25, probabilidade ou comprovação de capacidade.

Distribuição por requisitos usa as mesmas estatísticas. Chamadas diretas mantêm o contrato0045; helpers sem termos não ganham pesos artificiais. Releitura por hash e projeção de requisitos continuam ativas. Markdown mostra até cinco contribuições por hit e remete o restante ao JSON.

Método `go_ast_corpus_v1`, contrato `evidence-v9-corpus-retrieval`, checkpoint local `implementation_search: 0047`. Sem serviço, dependência ou execução externa nova; núcleo permanece puro.

## Avaliação congelada

Autor B preservou o corpus real0046 e acrescentou um holdout de outro projeto: `github.com/google/uuid` v1.6.0, commit `0f11ee6918f41a04c201eceeadf612a377bc7fbc`. As perguntas/esperados desse projeto não foram entregues ao autor do ranking. Baselines e holdout foram congelados antes da implementação; nenhuma revisão de esperado foi feita após resultados.

| Medida de recuperação | Antes | Depois |
| --- | --- | --- |
| Corpus real0046: arquivos-alvo nos primeiros3 caminhos | 0/6 | 2/6 |
| Corpus real0046: arquivos-alvo nos primeiros12 caminhos | 1/6 | 4/6 |
| Corpus real0046: negativas com algum indício | 1/3 | 1/3 |
| Holdout: arquivos-alvo nos cortes3 e12 | 2/2 | 2/2 |
| Holdout: negativa com algum indício | 1/1 | 1/1 |
| Controle sintético0045: positivos recuperados | 4/4 | 4/4 |

As métricas contam caminhos distintos entre até24 hits; são cortes de recuperação, não execução de capacidade. Precisão é contra um conjunto pequeno de arquivos rotulados, não julgamento exaustivo dos demais. O holdout preserva recall que já era perfeito e mantém indício na negativa: isso limita a força da evidência de generalização e não demonstra ausência de falsos positivos.

[Veredito B](../../lab/bancada/history/benchmarks/0047/evaluation/verdict.md) e [gate congelado](../../lab/bancada/history/benchmarks/0047/evaluation/candidate-comparison.json). A candidata passou na primeira revisão, sem ajuste aos esperados. Quatro testes novos exercitam raridade, normalização, contribuições finitas/determinísticas e preservação de exato/escopo AST; seis controles0045 também passaram.

## CLI completa e replay

As nove perguntas reais também foram executadas pela CLI, nos mesmos módulos. Todos os hashes Go varridos coincidiram com o manifesto congelado. O vínculo com commit de origem vem do manifesto dos módulos; a CLI lê distribuições sem checkout Git e não presume commit-base Git.

A coleta final passou a incluir:

- `db.go` na pergunta de leitura do bbolt;
- `transport.go` na pergunta sobre cabeçalho de autorização OAuth2;
- `oauth2.go` na pergunta sobre troca de código por token.

`db.go` na escrita já era coletado. As duas perguntas de Snappy continuam perdendo seus arquivos-alvo; ponderação não cria vocabulário inexistente. A coleta pela CLI confirma quatro de seis arquivos-alvo disponíveis, sem demonstrar que cada pergunta foi respondida completamente.

O controle `SanitizeAssessmentEvidence` continua recuperado. Na consulta natural sobre o próprio Bateia, `sem` recebeu contribuição por coincidir com identificador técnico do código. É uma ambiguidade lexical que IDF pode amplificar; a entrega não resolve interpretação de negação. A lista de fontes contém `idea_discovery.go` e `reuse_decision.go`, mas permanece incompleta para a pergunta. [Relatório natural preservado](../../lab/bancada/history/benchmarks/0047/natural.md).

A recomendação do Bateia manteve dez capacidades INCONCLUSIVE. Assessment, decisão, alternativas, bundles e `decision.md` foram reproduzidos byte a byte pelo replay. [Recibos, hashes, pesos e comparação CLI](../../lab/bancada/history/benchmarks/0047/cli-observations.json). Tempos são telemetria do Bateia, não benchmarks dos projetos analisados.

## Verificações e custódia

- `GOCACHE=/tmp/bateia-0047-cache go test ./...`: passou, incluindo gates existentes; testes HTTP locais executados com permissão correspondente.
- `crystalline-lint .`: `✓ No violations found`.
- `GOCACHE=/tmp/bateia-0047-cache go build -o /tmp/bateia-0047 ./04_wiring`: passou.
- `git diff --check`: passou.

Regime **A/B sem atestação física de isolamento**. Autoria de parser e de avaliação separadas; coordenador escreveu contrato, integração e relatório e consultou resultados públicos. Corpus adicional oculto do implementador, mas filesystem compartilhado impede alegar segregação atestada. Nenhum gate/certificado antigo ou esperado congelado foi reescrito. O contrato inicial e identidades anteriores permanecem em [contract-baseline.json](../../lab/bancada/history/benchmarks/0047/contract-baseline.json).

## Limites

Corpus pequeno e de conveniência, Go local, perguntas em português, lexicalidade e ambiguidades persistentes. O ganho no corpus conhecido não demonstra melhora universal. IDF pode destacar identificadores raros sem relação funcional com a frase; a normalização pode desfavorecer implementações longas pertinentes. Campos, contribuição e lacunas permitem auditar essas escolhas, mas não provam a solução. Tradução, negação, equivalência semântica, análise de outras linguagens e problemas anteriores de Bleve/gRPC ficam fora deste passo.
