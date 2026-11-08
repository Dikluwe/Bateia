# Prompt L0: Parser de Manifestos C++ (parser/cpp)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Revisado em**: 2026-07-31 — separação C puro (ADR-0003, Opção A)  
**Arquivos gerados**: `01_core/parser/cpp.go`, `01_core/parser/cpp_test.go`  

---

## Contexto

`cpp.go` trata os manifestos **compartilhados ou exclusivos de C++**: `CMakeLists.txt`, `vcpkg.json` e `conanfile.txt`. Dada a origem histórica, estes manifestos podiam ser usados tanto em projetos C puro quanto C++. Esta spec formaliza a heurística de detecção automática (ADR-0003, Opção A) que determina `Language: "c"` ou `Language: "cpp"` em tempo de parse.

Manifestos **exclusivos de C puro** (`Makefile`, `meson.build`) pertencem a `c.go`.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Funções estritamente puras. Zero I/O, zero leitura de disco, zero chamadas HTTP.
- Depende apenas de `01_core/domain`, `strings`, `regexp`, `encoding/json` e `fmt`.
- Proibido importar L₂, L₃ ou L₄.

---

## Heurísticas de Detecção de Linguagem (ADR-0003 Opção A)

### CMakeLists.txt — `detectCMakeLanguage(content)`

Inspeciona o corpo da declaração `project()`:

| Conteúdo do `project()` | `Language` produzido |
|---|---|
| `project(name CXX)` | `"cpp"` |
| `project(name C CXX)` | `"cpp"` (CXX presente) |
| `project(name LANGUAGES C)` | `"c"` (apenas C) |
| `project(name C)` | `"c"` |
| `project(name)` — sem linguagem | `"cpp"` (padrão histórico) |

Prioridade: se `\bCXX\b` presente → `"cpp"`; senão se `\bC\b` presente → `"c"`; senão → `"cpp"`.

### conanfile.txt — `detectConanLanguage(content)`

- Presença de `compiler.libcxx` no bloco `[settings]` → `"cpp"` (runtime C++ declarado).
- Ausência → `"c"` (projeto C puro sem runtime C++).

### vcpkg.json

Exclusivo de C++. Sempre produz `Language: "cpp"`. Sem heurística necessária.

---

## Instrução

Implementar em `cpp.go`:

1. `func detectCMakeLanguage(content string) string` — helper privado, regex sobre o corpo do `project()`.
2. `func detectConanLanguage(content string) string` — helper privado, regex sobre `compiler.libcxx`.
3. `func ParseCMakeListsTxt(ref domain.CommitRef, content string) (*domain.Manifest, error)` — usa `detectCMakeLanguage`.
4. `func ParseVcpkgJson(ref domain.CommitRef, content string) (*domain.Manifest, error)` — sempre `Language: "cpp"`.
5. `func ParseConanfileTxt(ref domain.CommitRef, content string) (*domain.Manifest, error)` — usa `detectConanLanguage`.

---

## Critérios de Verificação

```
project(libuv C)         → Language "c"
project(my_engine CXX)   → Language "cpp"
project(mixed C CXX)     → Language "cpp"
project(name)            → Language "cpp" (padrão)
LANGUAGES C              → Language "c"

conanfile com compiler.libcxx=libstdc++11  → Language "cpp"
conanfile sem compiler.libcxx              → Language "c"

vcpkg.json qualquer      → Language "cpp"
```

---

## Resultado Esperado

- `01_core/parser/cpp.go`: Parsers CMake, Vcpkg e Conan com detecção automática C/C++.
- `01_core/parser/cpp_test.go`: 8 testes cobrindo todas as combinações de heurística.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para C/C++ | cpp.go, cpp_test.go |
| 2026-07-31 | Opção A: heurística de detecção C vs C++ (ADR-0003) | cpp.go, cpp_test.go |

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.
