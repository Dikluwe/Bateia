# Prompt: Parser de Manifestos Emergentes e Sistemas (parser/emergent_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/emergent_ecosystems.go`, `01_core/parser/emergent_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas emergentes e de sistemas (Nim, Crystal, V e Erlang), os manifestos nativos (`*.nimble`, `shard.yml`, `v.mod`, `rebar.config`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseNimble(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `requires "package >= version"` em `*.nimble`.
   - Preencher `Language: "nim"`.

2. `func ParseShardYaml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name: ...` e dependências sob `dependencies:` em `shard.yml`.
   - Preencher `Language: "crystal"`.

3. `func ParseVMod(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name: '...'` e `deps: [...]` em `v.mod`.
   - Preencher `Language: "v"`.

4. `func ParseRebarConfig(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas sob `{deps, [...]}` em `rebar.config`.
   - Preencher `Language: "erlang"`.

---

## Critérios de Verificação

```
Dado um shard.yml com name: kemal e dependência github: kemalcr/kemal
Quando ParseShardYaml for executada
Então retorna um Manifest com Language "crystal" e ModuleName "kemal".
```

---

## Resultado Esperado

- `01_core/parser/emergent_ecosystems.go`: Implementação dos parsers Nim, Crystal, V e Erlang.
- `01_core/parser/emergent_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Nim, Crystal, V e Erlang | emergent_ecosystems.go, emergent_ecosystems_test.go |
