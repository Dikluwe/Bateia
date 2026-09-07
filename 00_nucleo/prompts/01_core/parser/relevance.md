# Prompt: Evaluator de Relevância e Import Matcher (parser/relevance)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/relevance.go`, `01_core/parser/relevance_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia, os repositórios candidatos precisam ser avaliados quanto à sua relevância real em relação à query pesquisada pelo usuário. O objetivo é diferenciar repositórios que importam diretamente uma biblioteca de repositórios que apenas possuem menções textuais genéricas.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP.
- Depende apenas de `01_core/domain` e bibliotecas padrão (`strings`, `path`).
- Proibido importar L₂, L₃ ou L₄.

---

## Instrução

Criar o arquivo `01_core/parser/relevance.go` com a struct `RelevanceResult` e a função pura:
`func EvaluateRelevance(query string, manifest *domain.Manifest, cand domain.RepositoryCandidate, tree []domain.DirectoryNode) RelevanceResult`

Regras do cálculo de relevância (Score 0 a 100):
1. **Normalização**: Limpar a query e termos comparados (lowercase, trim whitespace, remover sufixos de versão como `/v2`, `/v3`).
2. **DIRECT_IMPORT_MATCH (100 pts)**: O termo bate com o nome de uma dependência direta no `go.mod`.
3. **MODULE_NAME_MATCH (85 pts)**: O termo bate com o nome do módulo declarado no `go.mod` ou com o último segmento do nome do módulo.
4. **INDIRECT_IMPORT_MATCH (70 pts)**: O termo bate com uma dependência indireta (`Indirect: true`) no `go.mod`.
5. **STRUCTURAL_MATCH (50 pts)**: O termo bate com o nome de um diretório na árvore do repositório.
6. **GENERIC_TEXT_MATCH (20 pts)**: Nenhuma correspondência específica de import ou módulo encontrada.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 | relevance.go, relevance_test.go |

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.
