# Prompt: Parser de Manifestos Dinâmicos e Funcionais (parser/dynamic_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/dynamic_ecosystems.go`, `01_core/parser/dynamic_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas F#, Perl, Groovy e COBOL, os manifestos/código nativos (`*.fsproj`, `cpanfile`, `Makefile.PL`, `build.gradle`, `*.cbl`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseFsProj(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `<PackageReference Include="PackageName" Version="1.0" />` em `.fsproj`.
   - Preencher `Language: "fsharp"`.

2. `func ParseCpanfile(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `requires 'Package::Name'` ou `PREREQ_PM` em `cpanfile` / `Makefile.PL`.
   - Preencher `Language: "perl"`.

3. `func ParseGroovyBuild(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas em `implementation`, `api`, `testImplementation` em `build.gradle`.
   - Preencher `Language: "groovy"`.

4. `func ParseCobolSource(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências de `COPY BOOKNAME` ou `CALL "PROGRAM"` em arquivos `.cbl` / `.cob`.
   - Preencher `Language: "cobol"`.

---

## Critérios de Verificação

```
Dado um cpanfile com requires 'Dancer2', '0.400000';
Quando ParseCpanfile for executada
Então retorna um Manifest com Language "perl" e dependência "Dancer2".
```

---

## Resultado Esperado

- `01_core/parser/dynamic_ecosystems.go`: Implementação dos parsers F#, Perl, Groovy e COBOL.
- `01_core/parser/dynamic_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para F#, Perl, Groovy e COBOL | dynamic_ecosystems.go, dynamic_ecosystems_test.go |
