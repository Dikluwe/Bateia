# Prompt: Parser de Truncamento Inteligente de Conteúdo (parser/truncate)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/truncate.go`, `01_core/parser/truncate_test.go`  

---

## Contexto

Nos relatórios Markdown do Bateia, arquivos extraídos do fundo do funil (como `go.mod` e `README.md`) podem ser extensos. O truncamento ingênuo por fatiamento de caracteres/bytes pode cortar palavras ou sequências UTF-8 no meio. O truncamento inteligente garante que o corte ocorra estritamente ao final de quebras de linha (`\n`), preservando a validade UTF-8 e a legibilidade.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero chamadas de rede ou disco.
- Depende apenas da biblioteca padrão Go (`strings`, `unicode/utf8`, `fmt`, `net/url`).
- Proibido importar L₂, L₃ ou L₄.

---

## Instrução

Criar o arquivo `01_core/parser/truncate.go` com as seguintes funções puras:

1. `ShouldTruncate(content string, maxChars, maxLines int) bool`:
   Retorna `true` se o número de runes de `content` exceder `maxChars` OU se o número de linhas (`\n`) exceder `maxLines`.

2. `TruncateSmart(content string, maxChars, maxLines int) (string, bool)`:
   - Se `ShouldTruncate` for falso, retorna o conteúdo original e `false`.
   - Acumula caracteres até `maxChars` ou `maxLines` respeitando runes UTF-8.
   - Recua o ponto de corte para a última quebra de linha `\n` encontrada, evitando truncar palavras no meio.
   - Retorna o texto truncado limpo e o booleano `true`.

3. `SanitizeCommitRef(commit string) string`:
   Sanitiza a referência de commit para geração de URLs públicas no GitHub (converte `HEAD` ou vazio para `main`).

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 | truncate.go, truncate_test.go |
