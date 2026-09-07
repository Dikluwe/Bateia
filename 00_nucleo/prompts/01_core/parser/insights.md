# Prompt: Extrator de Insights Arquiteturais e Necessidades Ocultas (parser/insights)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/insights.go`, `01_core/parser/insights_test.go`  

---

## Contexto

Os relatórios do Bateia contêm a seção 3 destinada ao "Descobrimento de Necessidades Ocultas". Em vez de textos estáticos genéricos, esta funcionalidade analisa o conteúdo dos manifestos (`go.mod`, `package.json`, `pyproject.toml`, `requirements.txt`, `Cargo.toml`, `pom.xml`, `build.gradle`, `composer.json`, `CMakeLists.txt`, `vcpkg.json`, `.csproj`, `Project.toml`, `DESCRIPTION`, `.rockspec`, `dune-project`, `mix.exs`, `Gemfile`, `Package.swift`, `*.cabal`, `pubspec.yaml`, `build.sbt`, `project.clj`, `build.zig.zon`, `.dproj`, `.lpi`, `fpm.toml`, `.vbp`, `.fsproj`, `cpanfile`, `.cbl`, `*.nimble`, `shard.yml`, `v.mod`, `rebar.config`, `dub.json`, `gleam.toml`, `box.json`, `flex-config.xml`, `alire.toml`, `mojoproject.toml`, `Podfile`, `haxelib.json`, `META6.json`, `info.rkt`, `mason.toml`, `corral.json`, `ols.json`) e as árvores de arquivos dos repositórios extraídos para extrair métricas e ocorrências de 5 categorias essenciais:
1. Estratégia de Testes (`testify`, `jest`, `pytest`, `criterion`, `junit`, `silly-unittest`, `gleam/test`, `mxunit`, `flexunit`, `ahunit`, `mojotest`, `xctest`, `utest`, `zef`, `rackunit`, `ponytest`, etc.)
2. Configuração e Ambiente (`godotenv`, `viper`, `dotenv`, `spring-boot`, `dub`, `gleam_env`, `commandbox`, `mxml`, `alire`, `pixi`, `cocoapods`, `haxelib`, `mason`, `corral`, etc.)
3. Banco de Dados e Migrações (`pgx`, `gorm`, `prisma`, `vibe-d`, `pgo`, `cborm`, `as3commons`, `gnatcoll_db`, `max`, `sqlite`, `dbd-pg`, `db`, etc.)
4. Observabilidade e Comunicação (`zerolog`, `zap`, `fastapi`, `vibe-d`, `lustre`, `coldbox`, `robotlegs`, `aws_sdk`, `light_async`, `AFNetworking`, `heaps`, `cro`, `spiffy`, etc.)
5. Automação e DevOps (`.github/workflows/`, `Dockerfile`, `Makefile`, `dub`, `gleam format`, `box`, `mxmlc`, `alr`, `mojo`, `pod`, `haxe`, `zef`, `raco`, `mason`, `corral`, `ols`, etc.)

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

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.
