# Prompt: Parser de Manifestos Estendidos (parser/extended_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/extended_ecosystems.go`, `01_core/parser/extended_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas adicionais (D, Gleam, ColdFusion e ActionScript/Flex), os manifestos nativos (`dub.json`, `gleam.toml`, `box.json`, `flex-config.xml`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp`, `encoding/json` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseDubJson(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `"name"` e bloco `"dependencies"` em `dub.json` (D language).
   - Preencher `Language: "d"`.

2. `func ParseGleamToml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name = "..."` e bloco `[dependencies]` em `gleam.toml` (Gleam).
   - Preencher `Language: "gleam"`.

3. `func ParseBoxJson(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `"name"` e bloco `"dependencies"` em `box.json` (ColdFusion CommandBox).
   - Preencher `Language: "coldfusion"`.

4. `func ParseFlexConfig(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências de `<path-element>...</path-element>` em `flex-config.xml` (ActionScript / Flex).
   - Preencher `Language: "actionscript"`.

---

## Critérios de Verificação

```
Dado um gleam.toml com name = "my_gleam_app" e [dependencies] lustre = ">= 1.0.0"
Quando ParseGleamToml for executada
Então retorna um Manifest com Language "gleam", ModuleName "my_gleam_app" e dependência "lustre".
```

---

## Resultado Esperado

- `01_core/parser/extended_ecosystems.go`: Implementação dos parsers D, Gleam, ColdFusion e ActionScript.
- `01_core/parser/extended_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para D, Gleam, ColdFusion e ActionScript | extended_ecosystems.go, extended_ecosystems_test.go |
