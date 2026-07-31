# Prompt: Parser de Manifestos C / C++ (parser/cpp)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/cpp.go`, `01_core/parser/cpp_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas C e C++, o conteúdo de manifestos como `CMakeLists.txt` (CMake), `vcpkg.json` (Vcpkg) ou `conanfile.txt` (Conan) precisa ser analisado para extrair o nome do projeto e suas bibliotecas/dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp`, `encoding/json` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseCMakeListsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `project(...)` para `ModuleName`.
   - Extrair dependências de `find_package(...)` e `target_link_libraries(...)`.
   - Marcar `Language: "cpp"`.

2. `func ParseVcpkgJson(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Unmarshal JSON de `vcpkg.json`. Extrair `name` e a lista `dependencies`.
   - Marcar `Language: "cpp"`.

3. `func ParseConanfileTxt(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas sob o bloco `[requires]`.
   - Marcar `Language: "cpp"`.

---

## Critérios de Verificação

```
Dado um CMakeLists.txt com project(my_cpp_app) e find_package(Boost REQUIRED)
Quando ParseCMakeListsTxt for executada
Então retorna um Manifest com ModuleName "my_cpp_app", Language "cpp" e 1 dependência "boost".
```

---

## Resultado Esperado

- `01_core/parser/cpp.go`: Implementação dos parsers CMake, Vcpkg e Conan.
- `01_core/parser/cpp_test.go`: Suite de testes cobrindo manifestos C/C++.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para C/C++ | cpp.go, cpp_test.go |
