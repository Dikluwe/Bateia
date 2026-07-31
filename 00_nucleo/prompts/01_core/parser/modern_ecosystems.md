# Prompt: Parser de Manifestos Modernos e JVM Alt (parser/modern_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/modern_ecosystems.go`, `01_core/parser/modern_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas Dart/Flutter, Scala, Clojure e Zig, os manifestos nativos (`pubspec.yaml`, `build.sbt`, `project.clj`, `build.zig.zon`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParsePubspecYaml(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name: pkg_name` para `ModuleName`.
   - Extrair dependências sob `dependencies:` (tipo `"direct"`) e `dev_dependencies:` (tipo `"dev"`).
   - Preencher `Language: "dart"`.

2. `func ParseBuildSbt(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name := "module"` para `ModuleName`.
   - Extrair dependências de `libraryDependencies += "group" %% "artifact" % "version"`.
   - Preencher `Language: "scala"`.

3. `func ParseProjectClj(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `defproject group/name` para `ModuleName`.
   - Extrair dependências de `:dependencies [[group/artifact "version"]]`.
   - Preencher `Language: "clojure"`.

4. `func ParseBuildZigZon(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `.name = "pkg"` para `ModuleName`.
   - Extrair dependências declaradas sob `.dependencies = .{ ... }`.
   - Preencher `Language: "zig"`.

---

## Critérios de Verificação

```
Dado um pubspec.yaml com name: flutter_app e provider: ^6.0.0 sob dependencies
Quando ParsePubspecYaml for executada
Então retorna um Manifest com Language "dart", ModuleName "flutter_app" e 1 dependência "provider".
```

---

## Resultado Esperado

- `01_core/parser/modern_ecosystems.go`: Implementação dos parsers Dart, Scala, Clojure e Zig.
- `01_core/parser/modern_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Dart, Scala, Clojure e Zig | modern_ecosystems.go, modern_ecosystems_test.go |
