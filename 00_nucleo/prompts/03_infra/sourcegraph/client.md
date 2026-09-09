# Prompt: Adaptador de Descoberta Sourcegraph (sourcegraph/client)

**Camada**: L₃ — Infraestrutura  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `03_infra/sourcegraph/client.go`, `03_infra/sourcegraph/client_test.go`  

---

## Contexto

Conforme ADR-0001 (rev 1), o topo do funil de descoberta do Bateia utiliza o Sourcegraph como fonte primária, evitando consumir a quota de 30 req/min da Search API do GitHub.

---

## Restrições Estruturais

- **L₃ (Infraestrutura)**: Implementa a interface `domain.SearchProvider` de $L_1$.
- Realiza requisições HTTP POST para a API GraphQL do Sourcegraph (`https://sourcegraph.com/.api/graphql`).
- Trata respostas `Repository` e `FileMatch` no payload GraphQL unificado.
- Proibido importar $L_2$ ou $L_4$.

---

## Instrução

1. Definir `Client` que implementa `domain.SearchProvider`:
   ```go
   type Client struct {
       endpoint   string
       apiToken   string
       httpClient *http.Client
   }
   ```
2. Implementar `NewClient(apiToken string) *Client` com tempo limite HTTP de 45 segundos.
3. Implementar `SearchCandidates(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, error)`:
   - Formatar a query GraphQL com suporte a `repo:^github\.com/` e filtro por linguagem.
   - Tratar `__typename == "Repository"` e `__typename == "FileMatch"`.
   - Preencher `CommitID` com `"HEAD"` quando omitido e atribuir `ConsultadoEm` UTC.
   - Limitar os resultados retornado ao `MaxResults` informado.

---

## Critérios de Verificação

```
Dado uma resposta da API GraphQL com fragmentos FileMatch e Repository
Quando SearchCandidates for executado
Então converte os itens em []domain.RepositoryCandidate sem duplicatas.
```

---

## Resultado Esperado

- `03_infra/sourcegraph/client.go`: Adaptador Sourcegraph `Client`.
- `03_infra/sourcegraph/client_test.go`: Testes unitários com simulador de servidor HTTP GraphQL.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | client.go, client_test.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0054
Aplicar o contrato de `00_nucleo/passos/0054-preservar-intencao-na-descoberta.md`: consulta integral única, sem fallback destrutivo, erros explícitos e versão de evidência nova para invalidar cache anterior.

## Passo 0055 — hipótese rastreável

Aplicar o contrato congelado 00_nucleo/passos/0055-planejar-consultas-de-codigo.md, que substitui a tentativa única por original seguido de no máximo uma hipótese apenas após vazio válido. QueryPlan, RawCandidates e RejectedCandidates registram somente a derivada; contagens finais são admitidos. Plano preserva pergunta/restrições não verificadas. Ambas portas usam catálogo puro; linguagem e escopo regex literal persistem. Admitir somente FileMatch com caminho e ação+objeto em snippets efetivamente retidos do mesmo arquivo; reter todos os previews limitados individualmente a 240 runes, sem descrição/nome. Não combinar arquivos. Preservar limites, erros, timeouts, cancelamento e rastros nos candidatos. Validar envelope GraphQL antes de considerá-lo vazio. EvidenceVersion muda para evidence-v15-planned-discovery.

Envelope mínimo válido: data.search.results.results deve existir como array JSON (inclusive []); data/search/results nulos ou ausentes não são resposta vazia. Errors não vazio continua sendo erro, antes da validação de envelope.

Compatibilidade de MaxResults: valores <=0 preservam o limite efetivo anterior de um candidato; valores positivos limitam candidatos brutos retidos, mantendo previews posteriores desses mesmos repositórios. RawCandidates contabiliza os candidatos retidos antes da admissão lexical.

Adendo 0055: aplicar 00_nucleo/passos/0055-adendo-evidencia.md. SearchAttempt.AdmissionEvidence ([]ProviderCandidateMatch, JSON admission_evidence,omitempty) somente na derivada, preserva matches admitidos e os snippets efetivamente examinados da porta Discover, também nos rastros anexados aos candidatos. Não publicar rejeitados como evidência de admissão.

## Revisão de sintaxe 0055-v2
Aplicar 00_nucleo/passos/0055-adendo-sintaxe.md: emitir `patterntype:regexp (a|b) AND (x|y)`, sem aspas/content. O catálogo, filtros, bloqueios e admissão permanecem. Versão do plano 0055-v2 e EvidenceVersion evidence-v16-planned-discovery-syntax. Esta revisão substitui a serialização anterior do contrato.

## Passo 0056 — cobertura observada

Aplicar o contrato congelado `00_nucleo/passos/0056-explicar-cobertura-e-perdas.md`. Somente na derivada: count limitado20..200 sem overflow, coleta de até N repositórios distintos na ordem da resposta e evidências posteriores desses mesmos repositórios, admissão lexical antes do limite final. Cobertura unknown, matchCount nullable, contagens observadas e diagnósticos com motivos, no máximo32 previews retidos de240 runes normalizadas por repositório; omissões e truncamentos explícitos. Diagnósticos usam somente evidência de arquivo retida, exigindo ação e objeto no mesmo arquivo. Selected/insufficient_evidence/candidate_limit distinguem admissão e seleção. Propagar Coverage e CandidateDiagnostics nas tentativas e rastros; Matches/AdmissionEvidence somente selecionados. Literal permanece compatível. EvidenceVersion evidence-v17-discovery-coverage.
