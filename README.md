# 🔍 Bateia — Motor de Inteligência Arquitetural Multi-Linguagem

**Bateia existe para produzir oportunidades de reúso e insights de engenharia a partir da investigação de software.** Uma ideia ou um projeto conhecido orienta a descoberta de alternativas, mecanismos, restrições e evidências úteis para decidir como construir ou evoluir uma solução.

A decisão entre `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE`, `BUILD` ou `INCONCLUSIVE` é uma saída. O projeto também possui extração de mecanismos e síntese de hipóteses de melhoria com premissas, lacunas e planos de validação. Uma hipótese não comprova efeito no alvo nem autoriza adoção; essas capacidades ainda têm dependências e limites descritos no mapa arquitetural.

O funil **Descoberta → Triagem → Extração**, os parsers de manifestos, os coletores de código e os adaptadores AST/CFG são meios de obter evidência para essa jornada. Um repositório analisado é um candidato ou corpus de validação; não define a identidade do Bateia. Evidências de dependências sugerem aspectos a investigar, sem comprovar automaticamente capacidades ou necessidades.

A qualidade do produto deve ser avaliada pela utilidade e fundamentação das decisões em diferentes ideias, projetos e ecossistemas. Quantidade de arquivos, símbolos ou grafos extraídos e aprovação de testes de infraestrutura não substituem essa avaliação.

Intenção de origem: [D-001 — Da ideia à decisão de reúso](00_nucleo/desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [D-002 — Alternativas de software com evidência](00_nucleo/desejos/D-002-alternativas-de-software-com-evidencia.md).

## Organização arquitetural

Tekt governa o programa. O [mapa de responsabilidades](00_nucleo/arquitetura/README.md) liga capacidade, código, contrato e verificação. O [ADR-0008](00_nucleo/ADR/ADR-0008-hierarquia-de-arquiteturas-sob-tekt.md) delimita organizações locais candidatas; a primeira separação distingue análise de insights de seu protocolo de execução durável. A composição de descoberta usa um registro validado em L1 e o MCP anuncia suas fontes configuradas ([ADR-0009](00_nucleo/ADR/ADR-0009-admissao-de-conectores-e-projecao-mcp.md)). Para agentes, comece por [AGENTS.md](AGENTS.md).

## Documentação do projeto

Comece pelo [mapa de documentação](00_nucleo/documentacao/README.md), que reúne a [visão do produto](00_nucleo/documentacao/produto.md), a [arquitetura atual](00_nucleo/documentacao/arquitetura.md), o [guia de engenharia](00_nucleo/documentacao/engenharia.md) e o [plano de organização](00_nucleo/documentacao/organizacao.md).

## Investigação de insights sem decisão de reúso

A jornada offline recebe pergunta, contexto do alvo e evidências congeladas, sem exigir política, decisão ou run de recomendação:

```bash
go build -o /tmp/bateia-insights ./04_wiring
/tmp/bateia-insights idea insights analyze \
  --input lab/bancada/insights/0105/queue-documentation-input.json \
  --output /tmp/insight.json
```

O [exemplo de bancada](lab/bancada/insights/0105/README.md) é ilustrativo e não representa fonte remota autenticada. Produz hipótese, proveniência, lacunas e plano de baseline/aplicabilidade/experimento. A jornada usa avaliações fornecidas e verifica seu vínculo com os bytes; não reavalia automaticamente sua semântica nem comprova benefício no alvo. As classes suportadas são log append-only, fila limitada/backpressure, armazenamento por conteúdo e isolamento de portas/adaptadores. Adoção de código continua exigindo investigação e decisão próprias ([contrato 0105](00_nucleo/passos/0105-investigar-insights-sem-decisao-de-reuso.md)).

O [ensaio com fonte real no GitHub](lab/bancada/insights/0106/README.md) coleta `sourcegraph/conc` em commit fixado, registra avaliações manuais e verifica replay/integridade. A coleta online passou; o mecanismo permanece PARTIAL e a hipótese INCONCLUSIVE. Consulte o [laudo 0106](00_nucleo/laudos/0106-correcoes-e-ensaio-github.md).

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
| **$L_0$** | `00_nucleo/` | Intenção, prompts, ADRs, laudos e sínteses declarativas | Regime declarativo: sem código ou executores. |
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

### **Busca em workspaces públicos do Bitbucket**

Defina `BATEIA_BITBUCKET_WORKSPACES` com os IDs dos workspaces separados por vírgula antes de iniciar `./bateia mcp` (ou no arquivo `.env`):

```bash
BATEIA_BITBUCKET_WORKSPACES=atlassian,example ./bateia mcp
```

O conector `bitbucket-repository-catalog` entra na descoberta de mercado quando essa variável está definida. Ele consulta nome e descrição dos repositórios públicos nesses workspaces pela API Bitbucket Cloud. A busca não cobre todos os workspaces do Bitbucket nem usa a API de busca de código.

### **Pesquisa de software acadêmico**

O Bateia consulta dois catálogos públicos além dos repositórios convencionais:

- **zbMATH Open / swMATH** indexa software matemático e informa referências bibliográficas. O sinal `INDEXED` registra apenas a presença no catálogo; não representa revisão de código.
- **JOSS** publica software de pesquisa após revisão aberta. O sinal `PEER_REVIEWED` inclui DOI e link para a revisão. O endereço do código é mantido como vínculo não verificado até uma comprovação posterior.

```bash
./bateia idea research discover --query "computer algebra" --output pesquisa.json
./bateia idea research discover --input ideia.json --provider zbmath --output pesquisa.json
```

O snapshot guarda consulta, fonte, hash da resposta, URL do registro e candidatos. `--provider` aceita `all`, `zbmath` ou `joss`; `--max-per-query` limita os resultados por fonte. Os dois conectores também estão disponíveis na descoberta de mercado do MCP, sem credenciais. Consultas sem resultado e falhas do provedor são registradas separadamente.

Registros do zbMATH também podem trazer `related_software`, uma lista limitada de pistas com nome e URL do catálogo. Essa relação é declarada pelo catálogo; não confirma dependência, equivalência nem revisão do código. Por exemplo, a consulta `symbolic calculator` retorna Symja, cujo registro relaciona Giac.

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

## Piloto de observação de concorrência no Lab

O [piloto 0109](lab/bancada/concorrencia/0109/README.md) observa sintaxe Go por
propriedade, condição e escopo, com hashes/localizadores e lacunas explícitas.
Não promove mecanismos nem comprova limite global ou cancelamento completo.

```sh
go build -o /tmp/bateia-observer ./lab/bancada/concorrencia/0109/cmd/observer
/tmp/bateia-observer observe --input lab/bancada/concorrencia/0109/testdata/conc-input.json
/tmp/bateia-observer evaluate --corpus lab/bancada/concorrencia/0109/testdata/corpus.json
```

A CLI é experimental e permanece isolada do produto. Consulte o
[laudo 0109](00_nucleo/laudos/0109-observador-de-concorrencia-no-lab.md) para
verificação e limites do corpus.

O piloto de concorrência evoluiu para [resolução limitada de canais nomeados](lab/bancada/concorrencia/0110/README.md), com tipos em memória, casos de contraste e lacunas explícitas para pacotes incompletos.

O [passo 0111](lab/bancada/concorrencia/0111/README.md) fornece dependências congeladas ao observador: o tipo limiter de conc é resolvido offline, com hashes e controles de ausência/integridade. Análise de fluxo permanece pendente.

O [passo 0112](lab/bancada/concorrencia/0112/README.md) acrescenta relações tipadas entre envio no case, lançador e defer/receive, com ramos e lacunas explícitos. O grafo ainda não certifica limite global de concorrência.

O [passo 0113](lab/bancada/concorrencia/0113/README.md) rastreia parâmetros de callback até sua invocação em Catcher.Try, com identidade de tipos, escopos e controles de mutação. A cadeia estrutural mantém lacunas de execução/alcance.

O [passo 0114](lab/bancada/concorrencia/0114/README.md) distingue sites após terminadores e ordem de registro de defer, com alcance limitado e lacunas globais preservadas.

O [passo 0115](lab/bancada/concorrencia/0115/README.md) audita escritas, exposições e capturas do campo de canal, mantendo identidade de instância inconclusiva e distinguindo inventário de defeito.

O [passo 0116](lab/bancada/concorrencia/0116/README.md) localiza a guarda anterior à escrita de configuração de conc, separando checagem sequencial de atomicidade concorrente.

O [passo 0117](lab/bancada/concorrencia/0117/README.md) executa um modelo autoral que separa checagem sequencial, ausência de data race e atomicidade, sem executar conc ou promover garantias do candidato.

O [passo 0118](lab/bancada/concorrencia/0118/README.md) confronta oito alegações com documentação e testes congelados: reutilizar o pool após Wait não implica reconfigurá-lo. A revisão é autoral e mantém inconclusivos o suporte à sobreposição e a atomicidade.

O [passo 0119](lab/bancada/concorrencia/0119/README.md) executa os testes originais de conc com detector de race e confirma quatro comportamentos da revisão no recorte selecionado. Essa referência externa de execução não mede a precisão semântica geral do classificador.

O [passo 0120](lab/bancada/concorrencia/0120/README.md) avalia configuração/inicialização em tunny e pond com observador congelado. Localiza uma referência estrutural, mantém atomicidade sem promoção e revela que vereditos de contrato ainda não são implementados; referência anterior às respostas permanece autoral.

O [passo 0121](lab/bancada/concorrencia/0121/README.md) acrescenta uma vista documental: Go Doc é localizado automaticamente por dono AST e README pode receber associações fornecidas e validadas. Evidências ficam ligadas por contenção lexical, sem classificar contratos ou alterar o observador v8.

A [pesquisa 0122](00_nucleo/relatorios/pesquisa-modelos-de-decisao-0122.md) usa o coletor do Bateia para investigar Laya, Arbiter e avaliação de decisões tipadas. Recomenda piloto local por HTTP e prepara quatro pedidos diagnósticos, sem inferência ou calibração executadas.

O [passo 0123](lab/bancada/decisoes/0123/README.md) executa Laya real em CPU isolada com pesos fixados: três concordâncias em quatro diagnósticos autorais, estabilidade de opções/repetição e plano de treinamento por famílias de projetos. Ainda não há treino, calibração ou avaliação independente.

Pesquisa do Laya com o produto atual: [passo 0124](00_nucleo/passos/0124-investigar-melhorias-do-laya-com-bateia.md) e [laudo](00_nucleo/laudos/0124-pesquisa-do-laya-com-bateia.md). Coleta rastreável, limites do catálogo e controles de contexto/calibração, sem treino ou promoção produtiva.
