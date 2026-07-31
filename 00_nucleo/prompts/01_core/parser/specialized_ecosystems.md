# Prompt: Parser de Manifestos Especializados (parser/specialized_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/specialized_ecosystems.go`, `01_core/parser/specialized_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas especializados (Ada, Mojo, Objective-C, Haxe, Raku, Racket, Chapel, Pony e Odin), os manifestos nativos (`alire.toml`, `mojoproject.toml`, `Podfile`, `haxelib.json`, `META6.json`, `info.rkt`, `mason.toml`, `corral.json`, `ols.json`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp`, `encoding/json` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseAlireToml(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Ada - `alire.toml`)
2. `func ParseMojoProjectToml(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Mojo - `mojoproject.toml` / `pixi.toml`)
3. `func ParsePodfile(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Objective-C - `Podfile`)
4. `func ParseHaxelibJson(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Haxe - `haxelib.json`)
5. `func ParseMeta6Json(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Raku - `META6.json`)
6. `func ParseInfoRkt(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Racket - `info.rkt`)
7. `func ParseMasonToml(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Chapel - `mason.toml`)
8. `func ParseCorralJson(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Pony - `corral.json`)
9. `func ParseOlsJson(ref domain.CommitRef, content string) (*domain.Manifest, error)` (Odin - `ols.json`)

---

## Critérios de Verificação

```
Dado um alire.toml com name = "my_ada_app" e [dependencies] gnatcoll = "*"
Quando ParseAlireToml for executada
Então retorna um Manifest com Language "ada", ModuleName "my_ada_app" e dependência "gnatcoll".
```

---

## Resultado Esperado

- `01_core/parser/specialized_ecosystems.go`: Implementação dos parsers Ada, Mojo, Objective-C, Haxe, Raku, Racket, Chapel, Pony e Odin.
- `01_core/parser/specialized_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Ada, Mojo, Objective-C, Haxe, Raku, Racket, Chapel, Pony e Odin | specialized_ecosystems.go, specialized_ecosystems_test.go |
