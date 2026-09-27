# 🔍 Bateia — Motor de Inteligência Arquitetural Multi-Linguagem

**Bateia existe para ajudar a decidir o que reutilizar antes de começar um software do zero.** A partir de uma ideia ou de um software conhecido, descobre alternativas, compara capacidades e restrições e reúne evidências para uma decisão de engenharia.

A jornada principal recebe uma `IdeaSpec` e produz uma decisão entre `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE`, `BUILD` ou `INCONCLUSIVE`. O dossiê deve explicar o que cada candidato resolve, o que pode ser aproveitado, quais lacunas e bloqueios permanecem e por que uma alternativa é preferível às outras.

O funil **Descoberta → Triagem → Extração**, os parsers de manifestos, os coletores de código e os adaptadores AST/CFG são meios de obter evidência para essa jornada. Um repositório analisado é um candidato ou corpus de validação; não define a identidade do Bateia. Evidências de dependências sugerem aspectos a investigar, sem comprovar automaticamente capacidades ou necessidades.

A qualidade do produto deve ser avaliada pela utilidade e fundamentação das decisões em diferentes ideias, projetos e ecossistemas. Quantidade de arquivos, símbolos ou grafos extraídos e aprovação de testes de infraestrutura não substituem essa avaliação.

Intenção de origem: [D-001 — Da ideia à decisão de reúso](00_nucleo/desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [D-002 — Alternativas de software com evidência](00_nucleo/desejos/D-002-alternativas-de-software-com-evidencia.md).

## Da ideia à decisão de reúso

Para avaliar a própria worktree, incluindo alterações não commitadas:

```bash
./bateia idea recommend \
  --idea lab/bancada/history/benchmarks/0012/experimento/bateia-idea.json \
  --policy lab/bancada/history/benchmarks/0012/experimento/reuse-policy.json \
  --target-name "Bateia" \
  --target-repo Dikluwe/Bateia \
  --mode local \
  --local . \
  --output /dev/shm/bateia-runs/dossie-local
```

O diretório de saída contém a expansão da ideia, snapshot de descoberta, triagem, planos e bundles de evidência, assessment, decisão canônica, `alternatives.json`, `alternatives.md` e uma projeção legível em `decision.md`. `--target-name` é opcional e usa `IdeaSpec.problem` como fallback; `--target-repo` aceita `owner/repo` ou URL GitHub e exclui o próprio alvo do ranking independentemente do commit. Uma execução interrompida ou uma mudança apenas de política ou alvo pode usar `--resume`; o Bateia só reaproveita evidência quando os inputs causais e os hashes dos checkpoints continuam válidos.

Para pesquisa remota, use `--mode live`; para reproduzir uma decisão sem rede, use `--mode replay --replay <dossiê-fonte>`. Quando identidade, coleta ou evidência são opacas, o resultado permanece `INCONCLUSIVE` em vez de virar `BUILD` por ausência presumida.

No modo avançado, `--trace` adiciona `trace-graph.json` e `trace-index.json` ao dossiê; `idea decide --trace-output <arquivo>` emite somente o grafo. `idea trace build`, `idea trace impact` e `idea trace explain` constroem o índice de dependências, planejam revisões e localizam declarações no relatório original, offline. `idea trace compare` compara dois grafos congelados e planeja a revisão das diferenças registradas. Isso não troca o veredito nem autentica as evidências. Consulte os contratos [0098](00_nucleo/passos/0098-rastrear-dependencias-e-planejar-revisao.md), [0099](00_nucleo/passos/0099-localizar-declaracoes-no-relatorio-original.md) e [0100](00_nucleo/passos/0100-comparar-grafos-e-planejar-revisao.md) para formatos, limites e pertencimento ao manifest.

`idea trace review` registra progresso declarado das tarefas (pendente, em andamento, bloqueada ou concluída), com recibos vinculados ao plano por hash. Concluir tarefas não aprova nem troca a decisão histórica. Consulte o contrato [0101](00_nucleo/passos/0101-registrar-progresso-de-revisao.md).

`idea trace review-context` liga essas tarefas às declarações originais e aos ancestrais alterados, com localizadores compartilhados e progresso preservado. Consulte o contrato [0102](00_nucleo/passos/0102-ligar-revisao-as-declaracoes-originais.md).

---

## 🏛️ **Arquitetura Tekt (Five-Layer Lattice v1.4)**

O Bateia é construído sobre os princípios da **Arquitetura Cristalina / Tekt**, com separação de responsabilidades como regra de projeto. A tabela descreve a intenção das camadas; lint e testes verificam aspectos específicos, sem garantir pureza funcional de todo o núcleo:

| Camada | Diretório | Responsabilidade | Restrições |
|--------|-----------|------------------|------------|
| **$L_0$** | `00_nucleo/` | Intenção, prompts, ADRs, laudos e evidência passiva | Regime declarativo: sem código ou executores. |
| **$L_1$** | `01_core/` | Domínio, parsers e casos de uso | Sem adaptadores concretos; orquestração por interfaces. Há consultas ao relógio em parsers, portanto nem todas as funções são puras. |
| **$L_2$** | `02_shell/` | CLI e Gerador de Relatórios | Conhece apenas $L_1$. Formatadores de saída. |
| **$L_3$** | `03_infra/` | Adaptadores de Infraestrutura | I/O, GitHub API, Sourcegraph e Cache em disco. |
| **$L_4$** | `04_wiring/` | Ponto de entrada (`main.go`) | Composição sem lógica de negócio. |
| **Lab** | `lab/` | Arena, bancada, mutantes e experimentos | Não pode ser importado pelo produto; promoção exige contrato em $L_0$. |

No brownfield, a convergência para essas fronteiras é feita por fatias pequenas e
reversíveis. Relatórios e ADRs continuam em $L_0$; executores, gates e corpora
ativos pertencem ao Lab. A decisão está registrada no
[ADR-0004](00_nucleo/ADR/ADR-0004-lab-e-contencao-brownfield.md).

---

## 🌍 **Cobertura de leitura de manifestos**

O dispatcher em $L_1$ reconhece os formatos abaixo. Esta tabela descreve leitura de manifestos e identificação heurística de dependências, com limitações por parser; não comprova análise semântica completa desses ecossistemas. A coleta local de manifesto usa o primeiro manifesto reconhecido na raiz, sem inventariar todo monorepo.

Há análise AST de Go e um adaptador opcional Clang para C/C++, condicionado à configuração de compilação e a limites de coleta. Eles não dão cobertura equivalente às demais linguagens nem comprovam comportamento em execução:

| Ecossistema | Flag `-lang` | Manifesto Alvo | Exemplos de dependências reconhecidas |
|-------------|--------------|----------------|------------------------------------|
| 🐹 **Go** | `go` | `go.mod` | Testify, Pgx, Goose, Viper, Zap, Testcontainers |
| 🟦 **TypeScript / Node** | `typescript`, `ts`, `node` | `package.json` | Jest, Vitest, Prisma, NestJS, Express, Fastify, ESLint |
| 🐍 **Python** | `python`, `py` | `pyproject.toml`, `requirements.txt` | Pytest, FastAPI, Django, SQLAlchemy, Alembic, Ruff, Black |
| 🦀 **Rust** | `rust`, `rs` | `Cargo.toml` | Tokio, Serde, Axum, Actix-web, SQLx, Criterion |
| ☕ **Java / Kotlin** | `java`, `kotlin`, `kt` | `pom.xml`, `build.gradle`, `build.gradle.kts` | Spring Boot, JUnit, Mockito, Hibernate, Flyway, Quarkus |
| 🐘 **PHP** | `php` | `composer.json` | Laravel, Symfony, PHPUnit, Pest, Doctrine, Eloquent |
| ⚙️ **C / C++** | `c`, `cpp`, `c++` | `CMakeLists.txt`, `vcpkg.json`, `conanfile.txt` | Boost, GoogleTest, Catch2, fmt, spdlog, Qt, OpenCV |
| 💜 **C# / .NET** | `csharp`, `cs`, `dotnet` | `.csproj` | EntityFramework, Dapper, Serilog, xUnit, NUnit, MediatR |
| 🟣 **Julia** | `julia` | `Project.toml` | DataFrames, Flux, Plots, Test, LibPQ |
| 📊 **R** | `r` | `DESCRIPTION` | ggplot2, shiny, testthat, dplyr, DBI |
| 🌙 **Lua** | `lua` | `*.rockspec` | Busted, Lapis, Penlight, LuaSQL |
| 🐫 **OCaml** | `ocaml` | `dune-project`, `*.opam` | Dune, Alcotest, Dream, Cohttp, Caqti |
| 💧 **Elixir** | `elixir` | `mix.exs` | Phoenix, Ecto, ExUnit, Credo, Telemetry |
| 💎 **Ruby** | `ruby` | `Gemfile` | Rails, RSpec, ActiveRecord, Sinatra, RuboCop |
| 🐦 **Swift** | `swift` | `Package.swift` | Vapor, Alamofire, Realm, XCTest, SwiftLint |
| 📜 **Haskell** | `haskell` | `*.cabal`, `stack.yaml` | Yesod, Servant, Persistent, HSpec, HLint |
| 🎯 **Dart / Flutter** | `dart`, `flutter` | `pubspec.yaml` | Provider, Dio, Sqflite, Flutter_test, Flutter_lints |
| 🔴 **Scala** | `scala` | `build.sbt` | Http4s, Scalatest, Doobie, Akka-http, Pureconfig |
| 🍃 **Clojure** | `clojure`, `clj` | `project.clj`, `deps.edn` | Ring, Next.jdbc, Compojure, Clojure.test, Clj-kondo |
| ⚡ **Zig** | `zig` | `build.zig.zon`, `build.zig` | Httpz, Zap, Zig-sqlite, Std.testing, Zig fmt |
| 🏛️ **Pascal / Delphi** | `pascal`, `delphi` | `*.dproj`, `*.lpi` | FireDAC, Zeos, Indy, VCL, DUnit, IniFiles |
| 🧮 **Fortran** | `fortran` | `fpm.toml` | LAPACK, FPM, FUnit, OpenMP, MPI |
| 💻 **BASIC / VB** | `basic`, `vb` | `*.vbp` | ADODB, DAO, MSComCtl, VBWatch |
| 🔷 **F#** | `fsharp`, `fs` | `*.fsproj` | Giraffe, Expecto, FSharp.Configuration, Fantomas |
| 🐪 **Perl** | `perl` | `cpanfile`, `Makefile.PL` | Dancer2, DBI, Mojolicious, Test::More, PerlCritic |
| 🍃 **Groovy** | `groovy` | `build.gradle` | Spock, Ratpack, GORM, CodeNarc |
| 📜 **COBOL** | `cobol` | `*.cbl`, `*.cob` | GnuCOBOL, CobUnit, ESQL, CICS, COBC |
| 👑 **Nim** | `nim` | `*.nimble` | Jester, Norm, UnitTest, Parsimonious, Nimble |
| 💎 **Crystal** | `crystal` | `shard.yml` | Kemal, Crystal-DB, Spec, Shards |
| ⚡ **V (Vlang)** | `v`, `vlang` | `v.mod` | VWeb, V-SQLite, VTest, V fmt |
| 🔴 **Erlang** | `erlang` | `rebar.config` | Cowboy, Epgsql, Mnesia, EUnit, Rebar3 |
| 🎯 **D** | `d`, `dlang` | `dub.json` | Vibe-d, Silly-UnitTest, DUB |
| ⭐️ **Gleam** | `gleam` | `gleam.toml` | Lustre, Wisp, PGo, Gleam Test |
| ❄️ **ColdFusion** | `coldfusion`, `cf` | `box.json` | ColdBox, CBORM, CommandBox, MXUnit |
| ⚡ **ActionScript / Flex** | `actionscript`, `flex` | `flex-config.xml` | FlexUnit, Robotlegs, PureMVC, MXMLC |
| 🛡️ **Ada** | `ada` | `alire.toml` | GNATColl, AHunit, Alire, Alr |
| 🚀 **Mojo** | `mojo` | `mojoproject.toml`, `pixi.toml` | MojoTest, Pixi, Max, Light_Async |
| 🍎 **Objective-C** | `objective-c`, `objc` | `Podfile` | AFNetworking, Realm, SDWebImage, XCTest, CocoaPods |
| 🎯 **Haxe** | `haxe` | `haxelib.json` | Heaps, UTest, Haxelib |
| 🦋 **Raku** | `raku` | `META6.json` | Cro, DBD-Pg, Zef |
| 🌐 **Racket** | `racket` | `info.rkt` | RackUnit, Spiffy, Raco |
| 🧱 **Chapel** | `chapel` | `mason.toml` | Mason, Mason Test |
| 🐴 **Pony** | `pony` | `corral.json` | PonyTest, Corral |
| ⚔️ **Odin** | `odin` | `ols.json`, `odin.json` | OLS, Odin Testing |

---

## 📊 **Dependências observadas e hipóteses de investigação**

O relatório do funil pode agrupar dependências e arquivos encontrados nas categorias abaixo. A presença de uma biblioteca não demonstra sua utilização, qualidade arquitetural ou necessidade para a ideia do usuário. Sem candidatos ou evidência, essas categorias não são observações do ecossistema:

1. 🧪 **Estratégia de Testes**: Frameworks de testes unitários, integração e mocks.
2. ⚙️ **Configuração e Ambiente**: Gerenciadores de variáveis de ambiente e arquivos de configuração.
3. 🗄️ **Banco de Dados e Migrações**: Drivers de banco, ORMs e ferramentas de migration.
4. 📊 **Observabilidade e Comunicação**: Frameworks web/gRPC, loggers e bibliotecas de telemetria.
5. 🔄 **Automação e DevOps**: Presença de `Dockerfile`, GitHub Actions, linters e formatadores.

---

## 🚀 **Como Compilar e Usar**

### **Pré-requisitos**
- Go 1.25+

### **Compilação**
```bash
go build -o bateia ./04_wiring
```

### **Exemplos de Uso**

#### Analisar uma worktree local (incluindo alterações não commitadas)
```bash
./bateia analyze -path . -query "arquitetura do projeto" -format markdown -output /dev/shm/bateia-runs/analise_local.md
```

O modo local não depende do GitHub e ignora diretórios gerados ou volumosos como `.git`, `node_modules`, `vendor`, `target`, `dist` e `build`.

#### 1. Buscar Projetos Go usando Fiber
```bash
./bateia -query "github.com/gofiber/fiber" -lang go -max 10 -output /dev/shm/bateia-runs/fiber.md
```

#### 2. Buscar Projetos TypeScript usando NestJS
```bash
./bateia -query "@nestjs/core" -lang typescript -max 10 -output /dev/shm/bateia-runs/nestjs.md
```

#### 3. Buscar Projetos Python usando FastAPI
```bash
./bateia -query "fastapi" -lang python -max 10 -output /dev/shm/bateia-runs/fastapi.md
```

#### 4. Buscar Projetos Rust usando Axum
```bash
./bateia -query "axum" -lang rust -max 10 -output /dev/shm/bateia-runs/axum.md
```

#### 5. Buscar Projetos Java usando Spring Boot
```bash
./bateia -query "spring-boot" -lang java -max 10 -output /dev/shm/bateia-runs/spring.md
```

#### 6. Buscar Projetos PHP usando Laravel
```bash
./bateia -query "laravel" -lang php -max 10 -output /dev/shm/bateia-runs/laravel.md
```

#### 7. Buscar Projetos C++ usando CMake / Boost
```bash
./bateia -query "boost" -lang cpp -max 10 -output /dev/shm/bateia-runs/cpp.md
```

#### 8. Buscar Projetos C# / .NET usando Entity Framework
```bash
./bateia -query "EntityFrameworkCore" -lang csharp -max 10 -output /dev/shm/bateia-runs/dotnet.md
```

---

## 🧪 **Verificações de desenvolvimento**

Execute as verificações para a revisão em uso. Resultados são específicos do checkout e dos casos exercitados; aprovação não demonstra precisão de descoberta, utilidade das decisões ou ausência de defeitos. O [passo 0041](00_nucleo/relatorios/autofagia-katamari-0041.md) registra a validação daquela alteração.

```bash
# Executar suíte completa de testes unitários
go test ./...

# Validar conformidade arquitetural Tekt
crystalline-lint --fail-on warning .
```

## Limites conhecidos da decisão

A [execução 0043](00_nucleo/relatorios/execucao-promessas-buscas-0043.md) corrigiu os contraexemplos de autoridade, negação, vínculo de mecanismos e leitura local registrados no diagnóstico 0042. Essas correções têm escopo declarado: admissão de fontes não é prova de comportamento, e o catálogo pode não possuir detector para uma pergunta. Estimativas de esforço e complexidade são identificadas como heurísticas; não medem integração ou desempenho dos projetos.

Autofagia e Katamari produzem inventários e propostas para investigação. A implementação de planejamento não demonstra integração executável. Toda conclusão precisa distinguir declaração documental, indício estático e comportamento verificado.

O [plano de promessas e buscas verificáveis — 0043](00_nucleo/passos/0043-delimitar-promessas-e-planejar-buscas-verificaveis.md) organiza as próximas correções: conter leituras locais, uniformizar autoridade de evidência, vincular buscas a requisitos e produzir planos de reúso com condições de validação. A implementação e os limites da validação estão no relatório de execução 0043; a avaliação de utilidade com usuários permanece pendente.

A [execução 0044](00_nucleo/relatorios/recuperacao-implementacoes-0044.md) acrescenta recuperação local de definições Go antes da seleção de arquivos, com limites, linhas e hashes conferidos. O controle por símbolo exato foi recuperado; pertinência em linguagem natural continua heurística e encontrar uma declaração não comprova uma capacidade.

A [execução 0045](00_nucleo/relatorios/recuperacao-funcional-0045.md) prioriza funções por identificadores AST e inclui um salto de chamadas diretas com referências conferidas. Houve ganho no corpus sintético e na seleção de dois arquivos antes ausentes na consulta funcional local; compreensão semântica geral e capacidades permanecem não demonstradas.

A [execução 0046](00_nucleo/relatorios/recuperacao-real-0046.md) acrescenta evidências e lacunas por requisito em JSON/Markdown e um corpus real reproduzível. Na amostra de três projetos, o ranking recuperou apenas 1/6 arquivos-alvo nos primeiros 12 caminhos; uma tentativa sem ganho foi retirada. A projeção comunica incerteza e não comprova satisfação de capacidades.

A [execução 0047](00_nucleo/relatorios/ranking-corpus-0047.md), inspirada em funções que o Bateia localizou no Bleve, pondera raridade e comprimento e registra contribuições por termo. No corpus real congelado, a recuperação passou de 1/6 para 4/6 arquivos-alvo nos primeiros 12 caminhos, sem regressão no projeto adicional. O ganho é limitado à amostra; ambiguidades lexicais e comprovação de capacidades permanecem em aberto.

A [execução 0048](00_nucleo/relatorios/consulta-simbolos-0048.md) distingue consultas naturais e símbolos explícitos na recuperação de implementações Go. Marcadores como “sem” deixam de pontuar como identificadores em frases e aparecem como restrições não verificadas; símbolos explícitos continuam pesquisáveis. As métricas anteriores foram preservadas. Reconhecimento de negação não comprova seu cumprimento, e a pontuação de descoberta por README permanece fora deste recorte.

A [execução 0050](00_nucleo/relatorios/recuperacao-remota-0050.md) habilita recuperação de implementações Go em candidatos GitHub por commit e blobs verificados, com limites explícitos. No Snappy, a recuperação remota reproduziu os hits locais do mesmo snapshot e a CLI confirmou definições após releitura. Isso remove uma lacuna do caminho remoto; interpretação funcional, cobertura do catálogo e consumo agregado de cota continuam limitados.

A [execução 0051](00_nucleo/relatorios/documentacao-declaracoes-0051.md) recupera comentários Doc vinculados a funções Go, com fontes documentais separadas, intervalos e hashes revalidados. Ensaios locais/remotos localizaram descrições antes invisíveis. Documentação não comprova comportamento; a pergunta portuguesa de descompressão continua sem recuperação.

A [execução 0052](00_nucleo/relatorios/expansao-com-admissao-0052.md) registra expansões lexicais portuguesas como hipóteses e exige ação/contexto para admitir alternativas ambíguas. A pergunta original de descompressão recuperou Decode no Snappy local/remoto; uma versão irrestrita trouxe falsos indícios no UUID e foi substituída após controle negativo. O filtro continua lexical e o catálogo finito não demonstra compreensão geral.

A [avaliação 0053](00_nucleo/relatorios/avaliacao-jornada-0053.md) congelou nove perguntas de três novos domínios. O Bateia leu 6/6 arquivos-alvo e extraiu 5/6, mas recuperou somente1/6 funções esperadas. Descoberta pública:0/2 projetos esperados nas buscas avaliáveis e uma desconhecida. Nenhum algoritmo foi alterado; os dados priorizam correção de consulta/fallback e seleção de símbolos pertinentes.
