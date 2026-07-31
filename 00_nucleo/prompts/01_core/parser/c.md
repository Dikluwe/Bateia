# Prompt L0: Parser de Manifestos C (parser/c)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/c.go`, `01_core/parser/c_test.go`  

---

## Contexto

O ecossistema C puro é distinto do C++. Os seus manifestos canônicos são:

- `CMakeLists.txt` — também usado por C++, mas a linguagem declarada é `c` quando não há extensões C++;
- `Makefile` — sistema de build clássico sem gerenciador de pacotes;
- `meson.build` — sistema de build moderno, usado por projetos como GLib e GStreamer;
- `conan2.txt` / `conanfile.txt` — gerenciador de dependências (também compartilhado com C++, mas a label de linguagem difere).

O Bateia precisava de um parser **dedicado para C**, com `Language: "c"`, separado do `cpp.go` que marca projetos como `Language: "cpp"`. Sem esse arquivo, repositórios C puro eram ou não reconhecidos ou incorretamente classificados como C++.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Funções estritamente puras. Zero I/O, zero leitura de disco, zero chamadas HTTP.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar L₂, L₃ ou L₄.
- Não duplicar regex já declaradas em `cpp.go` — os parsers de Makefile e Meson são exclusivos deste arquivo.

---

## Instrução

Criar as seguintes funções em `01_core/parser/c.go`:

### 1. `func ParseMakefile(ref domain.CommitRef, content string) (*domain.Manifest, error)`

- Extrair a variável `TARGET` ou `PROG` para `ModuleName`; fallback para `ref.Repo`.
- Extrair dependências declaradas em `LIBS` ou flags `-l<lib>` no `LDFLAGS` / `LDLIBS`.
- Marcar `Language: "c"`.

### 2. `func ParseMesonBuild(ref domain.CommitRef, content string) (*domain.Manifest, error)`

- Extrair `project('name', ...)` para `ModuleName`.
- Extrair chamadas `dependency('lib')` para a lista de dependências.
- Marcar `Language: "c"`.

---

## Critérios de Verificação

```
Dado um Makefile com TARGET = myapp e LDLIBS = -lpthread -lm
Quando ParseMakefile for executada
Então retorna um Manifest com ModuleName "myapp", Language "c" e dependências ["pthread", "m"].

Dado um meson.build com project('libuv', 'c') e dependency('threads')
Quando ParseMesonBuild for executada
Então retorna um Manifest com ModuleName "libuv", Language "c" e 1 dependência "threads".
```

---

## Resultado Esperado

- `01_core/parser/c.go`: Implementação dos parsers Makefile e Meson para C puro.
- `01_core/parser/c_test.go`: Suíte de testes cobrindo os dois formatos de manifesto.

---

## Integrações Necessárias

- `manifest.go → DetectManifestType`: adicionar casos para `"makefile"` → `"c"` e `"meson.build"` → `"meson"`.
- `manifest.go → ParseManifest`: adicionar `case "c_makefile"`, `case "meson"`.
- `03_infra/github/client.go → manifestFiles`: adicionar `"Makefile"` e `"meson.build"` à lista de alvos de busca.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial — separação de C puro de C++ (ADR-0003) | c.go, c_test.go, manifest.go |
