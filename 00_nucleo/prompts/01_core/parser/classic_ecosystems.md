# Prompt: Parser de Manifestos Clássicos e Legados (parser/classic_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/classic_ecosystems.go`, `01_core/parser/classic_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas clássicos (Pascal/Delphi, Fortran e BASIC), os manifestos nativos (`.dproj`, `.lpi`, `fpm.toml`, `.vbp`) precisam ser analisados para extrair o nome do módulo e suas dependências/referências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParsePascalDelphi(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `DCCReference` ou `PackageImport` em arquivos `.dproj` (Delphi) / `.lpi` (Lazarus).
   - Preencher `Language: "pascal"`.

2. `func ParseFortranFpm(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name = "pkg"` e dependências declaradas sob `[dependencies]` em `fpm.toml`.
   - Preencher `Language: "fortran"`.

3. `func ParseBasicVbp(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas em `Reference=...` ou `Object=...` em arquivos `.vbp` (Visual Basic).
   - Preencher `Language: "basic"`.

---

## Critérios de Verificação

```
Dado um fpm.toml com name = "lapack_wrapper" e [dependencies] stdlib = "*"
Quando ParseFortranFpm for executada
Então retorna um Manifest com Language "fortran", ModuleName "lapack_wrapper" e dependência "stdlib".
```

---

## Resultado Esperado

- `01_core/parser/classic_ecosystems.go`: Implementação dos parsers Pascal/Delphi, Fortran e BASIC.
- `01_core/parser/classic_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Pascal/Delphi, Fortran e BASIC | classic_ecosystems.go, classic_ecosystems_test.go |
