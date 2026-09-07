# Prompt: Orquestrador do Funil de Busca (usecase/funnel)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-25  
**Arquivos gerados**: `01_core/usecase/funnel.go`, `01_core/usecase/funnel_test.go`  

---

## Contexto

O `SearchFunnel` é a regra de negócio central do Bateia. Ele coordena a busca em três estágios (Descoberta -> Triagem -> Extração), garantindo a invariante de redução de volume (1000 -> 50 -> 5) e a verificação prévia do cache por `repo@commit` para idempotência. Agora, ele também suporta busca recursiva por grafo de dependências (Estágio 4) para exploração em malha fechada.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Não realiza chamadas diretas de I/O.
- **Injeção de Abstrações**: O `SearchFunnel` depende **exclusivamente** das interfaces de $L_1$ (`domain.SearchProvider`, `domain.ContentProvider`, `domain.CacheStore`).
- **Proibido Importar $L_3$**: Não importa `03_infra/sourcegraph`, `03_infra/github` ou `03_infra/cache`. As implementações concretas serão injetadas por $L_4$ (Fiação).

---

## Instrução

1. Definir o struct:
   ```go
   type SearchFunnel struct {
       searcher domain.SearchProvider
       fetcher  domain.ContentProvider
       store    domain.CacheStore
   }
   ```
2. Implementar `NewSearchFunnel(searcher domain.SearchProvider, fetcher domain.ContentProvider, store domain.CacheStore) *SearchFunnel`.
3. Implementar `Execute(ctx context.Context, filter domain.QueryFilter) (*domain.SearchReport, error)`:
   - **Estágio 1 (Descoberta)**: Consultar `store.GetDiscovery`. Se der cache miss ou retornar lista vazia, consultar `searcher.SearchCandidates`. Salvar no cache se `len > 0`.
   - **Estágio 2 (Triagem)**: Para cada candidato, tentar `store.GetManifest`. Se miss, buscar via `fetcher.FetchManifest`. Se retornar 404/stale, marcar candidato e descartar. Avaliar estrutura via `parser.EvaluateStructure`.
   - **Estágio 3 (Extração)**: Limitar a no máximo 10 repositórios validados. Tentar `store.GetExtracted`. Se miss, buscar via `fetcher.ExtractRepoFiles` e salvar no cache.
   - **Estágio 4 (Recursividade)**: Se `filter.Depth > 1`, iterar pelas principais dependências (ex: as top 3 não-indiretas) encontradas nos manifestos extraídos que não estejam num set de visitados (`visited map[string]bool`). Para cada dependência, gerar uma nova query, decrescer `filter.Depth - 1`, e executar `Execute` recursivamente. Anexar os resultados no relatório atual.
   - Retornar o `SearchReport` totalmente montado.

---

## Critérios de Verificação

```
Dado mocks das interfaces SearchProvider, ContentProvider e CacheStore
Quando Execute for chamado com um filtro válido
Então orquestra os 3 estágios e retorna o SearchReport sem erros.

Dado um cache pré-existente para a busca (Cache HIT no topo)
Quando Execute for chamado
Então utiliza os candidatos cacheados sem invocar o searcher de rede.
```

---

## Resultado Esperado

- `01_core/usecase/funnel.go`: Implementação pura do caso de uso `SearchFunnel`.
- `01_core/usecase/funnel_test.go`: Suite de testes unitários utilizando mocks em memória para comprovação de fluxo e retenção de cache.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-25 | Criação inicial Tekt L0 | funnel.go, funnel_test.go |
| 2026-08-03 | Adição da especificação de busca recursiva (Estágio 4) | funnel.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.

## Passo 0035 — Análise semântica C/C++

Cumprir `00_nucleo/passos/0035-corrigir-lint-e-analisar-c-cpp-com-clang.md`. Porta L1 SemanticProvider, execução Clang limitada em L3, injeção L4. Contexto explícito em compile_commands.arguments, sem executar command, plugins ou response files; capturar hashes de configuração, fonte e headers, argumentos, versão, símbolos e chamadas com intervalos. Falta de evidência nunca vira complete. Exposição L2 separa resolução estática e adequação à consulta.

## Passo 0036
Implementar o contrato de `00_nucleo/passos/0036-projetar-evidencias-da-consulta.md`: evidência estática vinculada a corpos/hashes, seleção limitada, ausência de dimensões indeterminada, recibos completos separados da projeção. Não inferir execução ou solução exata.
