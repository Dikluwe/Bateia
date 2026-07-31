# Prompt: Extrator de Insights Arquiteturais e Necessidades Ocultas (parser/insights)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/insights.go`, `01_core/parser/insights_test.go`  

---

## Contexto

Os relatórios do Bateia contêm a seção 3 destinada ao "Descobrimento de Necessidades Ocultas". Em vez de textos estáticos genéricos, esta funcionalidade analisa o conteúdo dos manifestos (`go.mod`, `package.json`, `pyproject.toml`, `requirements.txt`, `Cargo.toml`, `pom.xml`, `build.gradle`, `composer.json`, `CMakeLists.txt`, `vcpkg.json`, `.csproj`, `Project.toml`, `DESCRIPTION`, `.rockspec`, `dune-project`, `mix.exs`, `Gemfile`, `Package.swift`, `*.cabal`, `pubspec.yaml`, `build.sbt`, `project.clj`, `build.zig.zon`) e as árvores de arquivos dos repositórios extraídos para extrair métricas e ocorrências de 5 categorias essenciais:
1. Estratégia de Testes (`testify`, `jest`, `pytest`, `criterion`, `junit`, `flutter_test`, `scalatest`, `clojure.test`, `zig test`, etc.)
2. Configuração e Ambiente (`godotenv`, `viper`, `dotenv`, `spring-boot`, `pureconfig`, `environ`, etc.)
3. Banco de Dados e Migrações (`pgx`, `gorm`, `prisma`, `sqflite`, `doobie`, `slick`, `next.jdbc`, `toucan`, etc.)
4. Observabilidade e Comunicação (`zerolog`, `zap`, `fastapi`, `dio`, `http4s`, `akka-http`, `ring`, `comport`, `httpz`, etc.)
5. Automação e DevOps (`.github/workflows/`, `Dockerfile`, `Makefile`, `eslint`, `flutter_lints`, `scalafmt`, `clj-kondo`, `zig fmt`, etc.)

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
