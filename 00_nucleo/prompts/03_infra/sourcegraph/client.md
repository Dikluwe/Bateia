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
