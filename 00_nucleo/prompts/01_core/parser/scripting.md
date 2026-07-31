# Prompt: Parser de Manifestos Científicos/Scripting (parser/scripting)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/scripting.go`, `01_core/parser/scripting_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para as linguagens Julia, R, Lua e OCaml, os manifestos nativos (`Project.toml`, `DESCRIPTION`, `*.rockspec`, `dune-project`) precisam ser analisados para extrair o nome do pacote e suas dependências declaradas.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseJuliaProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name` sob `[deps]` ou topo do `Project.toml`.
   - Extrair dependências declaradas no bloco `[deps]`.
   - Preencher `Language: "julia"`.

2. `func ParseRDescription(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `Package:` para `ModuleName`.
   - Extrair dependências listadas nos campos `Imports:`, `Depends:` (tipo `"direct"`) e `Suggests:` (tipo `"dev"`).
   - Preencher `Language: "r"`.

3. `func ParseLuaRockspec(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `package = "..."` para `ModuleName`.
   - Extrair dependências de `dependencies = { ... }`.
   - Preencher `Language: "lua"`.

4. `func ParseOCamlDune(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `(name ...)` para `ModuleName`.
   - Extrair dependências declaradas sob `(depends ...)`.
   - Preencher `Language: "ocaml"`.

---

## Critérios de Verificação

```
Dado um Project.toml com name = "DataAnalysis.jl" e DataFrames = "a93c63e6-5b32-5aa1-634d-178039b1d722" sob [deps]
Quando ParseJuliaProjectToml for executada
Então retorna um Manifest com ModuleName "DataAnalysis.jl", Language "julia" e dependência "DataFrames".
```

---

## Resultado Esperado

- `01_core/parser/scripting.go`: Implementação dos parsers Julia, R, Lua e OCaml.
- `01_core/parser/scripting_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Julia, R, Lua e OCaml | scripting.go, scripting_test.go |
