# Contrato de evidência

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.

## Passo 0037
Aplicar o contrato de resposta por mecanismo em `00_nucleo/passos/0037-resposta-fundamentada-por-mecanismo.md`. Contexto único de compilação, trechos exatos, limites explícitos e estados desconhecidos preservados.

## Passo 0038
Aplicar contrato do protótipo AST+CFG: adaptador opcional explicitamente configurado, hashes e limites, fluxo intraprocedural e estados não avaliados.

## 0044
Versão evidence-v6-go-retrieval identifica inclusão de recuperação Go antes da seleção; não muda autoridade de comprovação.

0045: versionar EvidenceVersion para evidence-v7-functional-retrieval pela alteração na ordenação e nas referências estáticas conferidas.

0046: EvidenceVersion evidence-v8-requirement-retrieval para projeção por requisito após releitura.

0047: se candidato for promovido, versão evidence-v9-corpus-retrieval para pesos pelo corpus em hits.

0048: contrato evidence-v10-query-interpretation registra separação de consulta/símbolos e restrições sem alterar linguagem do corpus.

0050: EvidenceVersion evidence-v11-remote-retrieval identifica recuperação GitHub por commit/blob verificados e limites explícitos. Ranking local preservado.

0051: EvidenceVersion evidence-v12-declaration-documentation para canal documental separado e intervalos revalidados.

0052: EvidenceVersion evidence-v13-query-expansions para recibos de expansão lexical query-only.

0052 r2: expansão por catálogo só participa de hits admissíveis por ação/contexto no canal examinado; adendo 0052-adendo-admissao-por-acao.md prevalece sobre admissão irrestrita inicial. Não alterar versão do catálogo de pares; recibo do passo registra revisão da política.

## Passo 0054
Aplicar o contrato de `00_nucleo/passos/0054-preservar-intencao-na-descoberta.md`: consulta integral única, sem fallback destrutivo, erros explícitos e versão de evidência nova para invalidar cache anterior.

## Passo 0055 — hipótese rastreável

Aplicar o contrato congelado 00_nucleo/passos/0055-planejar-consultas-de-codigo.md, que substitui a tentativa única por original seguido de no máximo uma hipótese apenas após vazio válido. QueryPlan, RawCandidates e RejectedCandidates registram somente a derivada; contagens finais são admitidos. Plano preserva pergunta/restrições não verificadas. Ambas portas usam catálogo puro; linguagem e escopo regex literal persistem. Admitir somente FileMatch com caminho e ação+objeto em snippets efetivamente retidos do mesmo arquivo; reter todos os previews limitados individualmente a 240 runes, sem descrição/nome. Não combinar arquivos. Preservar limites, erros, timeouts, cancelamento e rastros nos candidatos. Validar envelope GraphQL antes de considerá-lo vazio. EvidenceVersion muda para evidence-v15-planned-discovery.

## Revisão de sintaxe 0055-v2
Aplicar 00_nucleo/passos/0055-adendo-sintaxe.md: emitir `patterntype:regexp (a|b) AND (x|y)`, sem aspas/content. O catálogo, filtros, bloqueios e admissão permanecem. Versão do plano 0055-v2 e EvidenceVersion evidence-v16-planned-discovery-syntax. Esta revisão substitui a serialização anterior do contrato.

## Passo 0056 — cobertura observada

Aplicar o contrato congelado `00_nucleo/passos/0056-explicar-cobertura-e-perdas.md`. Somente na derivada: count limitado20..200 sem overflow, coleta de até N repositórios distintos na ordem da resposta e evidências posteriores desses mesmos repositórios, admissão lexical antes do limite final. Cobertura unknown, matchCount nullable, contagens observadas e diagnósticos com motivos, no máximo32 previews retidos de240 runes normalizadas por repositório; omissões e truncamentos explícitos. Diagnósticos usam somente evidência de arquivo retida, exigindo ação e objeto no mesmo arquivo. Selected/insufficient_evidence/candidate_limit distinguem admissão e seleção. Propagar Coverage e CandidateDiagnostics nas tentativas e rastros; Matches/AdmissionEvidence somente selecionados. Literal permanece compatível. EvidenceVersion evidence-v17-discovery-coverage.


## Passo 0057
EvidenceVersion passa a evidence-v18-function-behavior.


## Adendo 0057 — Origem documental
EvidenceVersion passa a evidence-v19-function-behavior-origin para invalidar projeções anteriores.

## Passo 0058
EvidenceVersion passa a evidence-v20-local-query-plan para identificar a integração
local das hipóteses lexicais de ação e objeto do plano existente.

## Passo 0059
EvidenceVersion passa a evidence-v21-documentation-budget para identificar a política
de reserva documental limitada, sem alterar autoridade ou validação de comportamento.

## Passo 0060
EvidenceVersion passa a evidence-v22-path-query pela tradução local explícita
caminho/caminhos para path, com provenance lexical e comportamento não verificado.


## Passo 0061
EvidenceVersion = evidence-v23-action-hints.
