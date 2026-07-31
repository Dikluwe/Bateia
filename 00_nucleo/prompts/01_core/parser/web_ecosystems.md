# Prompt: Parser de Manifestos Web/Ecossistemas (parser/web_ecosystems)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/web_ecosystems.go`, `01_core/parser/web_ecosystems_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para os ecossistemas Elixir, Ruby, Swift e Haskell, os manifestos nativos (`mix.exs`, `Gemfile`, `Package.swift`, `*.cabal`) precisam ser analisados para extrair o nome do módulo e suas dependências.

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `strings`, `regexp` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar as funções:
1. `func ParseMixExs(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `app: :name` para `ModuleName`.
   - Extrair dependências declaradas em `{:dep_name, "version"}`.
   - Preencher `Language: "elixir"`.

2. `func ParseGemfile(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair dependências declaradas via `gem 'name'`, `gem "name", "~> 1.0"`.
   - Distinguir dependências sob `group :development, :test` como tipo `"dev"`.
   - Preencher `Language: "ruby"`.

3. `func ParsePackageSwift(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name: "MyPackage"` para `ModuleName`.
   - Extrair dependências de `.package(url: ..., from: ...)` ou `.package(url: ..., "1.0.0")`.
   - Preencher `Language: "swift"`.

4. `func ParseHaskellCabal(ref domain.CommitRef, content string) (*domain.Manifest, error)`
   - Extrair `name: pkgname` para `ModuleName`.
   - Extrair dependências declaradas sob `build-depends:`.
   - Preencher `Language: "haskell"`.

---

## Critérios de Verificação

```
Dado um Gemfile com gem 'rails', '~> 7.0' e gem 'rspec-rails', group: :test
Quando ParseGemfile for executada
Então retorna um Manifest com Language "ruby", 1 dependência "direct" (rails) e 1 dependência "dev" (rspec-rails).
```

---

## Resultado Esperado

- `01_core/parser/web_ecosystems.go`: Implementação dos parsers Elixir, Ruby, Swift e Haskell.
- `01_core/parser/web_ecosystems_test.go`: Suite de testes unitários.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para Elixir, Ruby, Swift e Haskell | web_ecosystems.go, web_ecosystems_test.go |
