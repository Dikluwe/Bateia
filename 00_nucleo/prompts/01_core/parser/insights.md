# Prompt: Extrator de Insights Arquiteturais e Necessidades Ocultas (parser/insights)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/insights.go`, `01_core/parser/insights_test.go`  

---

## Contexto

Os relatórios do Bateia contêm a seção 3 destinada ao "Descobrimento de Necessidades Ocultas". Em vez de textos estáticos genéricos, esta funcionalidade analisa o conteúdo dos manifestos (`go.mod`, `package.json`, `pyproject.toml`, `requirements.txt`, `Cargo.toml`, `pom.xml`, `build.gradle`, `composer.json`, `CMakeLists.txt`, `vcpkg.json`, `.csproj`, `Project.toml`, `DESCRIPTION`, `.rockspec`, `dune-project`, `mix.exs`, `Gemfile`, `Package.swift`, `*.cabal`) e as árvores de arquivos dos repositórios extraídos para extrair métricas e ocorrências de 5 categorias essenciais:
1. Estratégia de Testes (`testify`, `testcontainers`, `jest`, `vitest`, `pytest`, `criterion`, `junit`, `phpunit`, `pest`, `mockery`, `gtest`, `catch2`, `xunit`, `nunit`, `test`, `testthat`, `busted`, `alcotest`, `ex_unit`, `rspec`, `quickcheck`, `hspec`, etc.)
2. Configuração e Ambiente (`godotenv`, `viper`, `dotenv`, `vlucas/phpdotenv`, `symfony/dotenv`, `spring-boot`, `fmt`, `spdlog`, `microsoft.extensions.configuration`, `pydantic`, `config.jl`, `argparse`, `dune`, `vapor`, `swiftyjson`, etc.)
3. Banco de Dados e Migrações (`pgx`, `gorm`, `prisma`, `sqlalchemy`, `sqlx`, `hibernate`, `doctrine`, `illuminate/database`, `eloquent`, `sqlite3`, `microsoft.entityframeworkcore`, `dapper`, `libpq.jl`, `dbi`, `luasql`, `caqti`, `ecto`, `active_record`, `realm`, `persistent`, etc.)
4. Observabilidade e Comunicação (`zerolog`, `zap`, `winston`, `fastapi`, `tokio`, `laravel`, `symfony`, `monolog`, `guzzlehttp`, `sentry`, `spring-web`, `boost`, `qt`, `serilog`, `nlog`, `mediatr`, `flux.jl`, `shiny`, `lapis`, `dream`, `phoenix`, `rails`, `sinatra`, `alamofire`, `scotty`, `yesod`, etc.)
5. Automação e DevOps (`.github/workflows/`, `Dockerfile`, `Makefile`, `eslint`, `ruff`, `clippy`, `phpstan`, `psalm`, `php-cs-fixer`, `clang-format`, `stylecop`, `juliac`, `styler`, `luacheck`, `dune`, `credo`, `rubocop`, `swiftlint`, `hlint`, etc.)

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero requisições HTTP ou acessos de disco.
- Depende apenas da biblioteca padrão de Go e de `01_core/domain`.
- Proibido importar L₂, L₃ ou L₄.

---

## Instrução

Criar o arquivo `01_core/parser/insights.go` com a função pura:
`func ExtractArchitecturalInsights(results []domain.ExtractedRepository) domain.ArchitecturalInsights`

1. Processar cada repositório extraído.
2. Mapear ocorrências de pacotes no manifesto correspondente (`go.mod` ou `package.json`), selecionando o dicionário apropriado com base em `manifest.Language`.
3. Agregar contagem, porcentagem de repositórios e versões comuns encontradas.
4. Mapear presença de arquivos de devops e automação na estrutura de diretórios (`Tree`).
5. Retornar a struct `ArchitecturalInsights` preenchida com as 5 categorias.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 | insights.go, insights_test.go |
| 2026-07-31 | Suporte Multi-Linguagem (TypeScript/Node.js) | insights.go |
