# 🔍 Bateia — Motor de Inteligência Arquitetural Multi-Linguagem

**Bateia** é uma ferramenta CLI e plataforma de inteligência arquitetural projetada para buscar, triar e analisar repositórios de código open-source. Através de um funil de três estágios (**Descoberta → Triagem → Extração**), o Bateia extrai dependências, identifica padrões e gera relatórios com o **Descobrimento de Necessidades Ocultas** (insights de testes, banco de dados, observabilidade, devops e configuração).

---

## 🏛️ **Arquitetura Tekt (Five-Layer Lattice v1.4)**

O Bateia é construído sobre os princípios da **Arquitetura Cristalina / Tekt**, garantindo separação rígida de responsabilidades e pureza de domínio:

| Camada | Diretório | Responsabilidade | Restrições |
|--------|-----------|------------------|------------|
| **$L_0$** | `00_nucleo/prompts/` | Especificações e Prompts perenes | Fonte imutável da verdade. Não é código. |
| **$L_1$** | `01_core/` | Domínio Puro, Parsers e Insights | **Zero I/O**, zero chamadas HTTP ou leituras de disco. |
| **$L_2$** | `02_shell/` | CLI e Gerador de Relatórios | Conhece apenas $L_1$. Formatadores de saída. |
| **$L_3$** | `03_infra/` | Adaptadores de Infraestrutura | I/O, GitHub API, Sourcegraph e Cache em disco. |
| **$L_4$** | `04_wiring/` | Ponto de entrada (`main.go`) | Composição sem lógica de negócio. |

---

## 🌍 **Ecossistemas e Manifestos Suportados**

O Bateia utiliza o **Dispatcher Pattern** genérico em $L_1$ para parsear manifestos nativos sem acoplamento de infraestrutura:

| Ecossistema | Flag `-lang` | Manifesto Alvo | Ferramentas & Insights Rastreáveis |
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

---

## 📊 **Seção 3: Descobrimento de Necessidades Ocultas**

Em todos os relatórios gerados, o Bateia analisa o ecossistema dos repositórios extraídos e compila métricas de 5 categorias essenciais:

1. 🧪 **Estratégia de Testes**: Frameworks de testes unitários, integração e mocks.
2. ⚙️ **Configuração e Ambiente**: Gerenciadores de variáveis de ambiente e arquivos de configuração.
3. 🗄️ **Banco de Dados e Migrações**: Drivers de banco, ORMs e ferramentas de migration.
4. 📊 **Observabilidade e Comunicação**: Frameworks web/gRPC, loggers e bibliotecas de telemetria.
5. 🔄 **Automação e DevOps**: Presença de `Dockerfile`, GitHub Actions, linters e formatadores.

---

## 🚀 **Como Compilar e Usar**

### **Pré-requisitos**
- Go 1.22+

### **Compilação**
```bash
go build -o bateia ./04_wiring/main.go
```

### **Exemplos de Uso**

#### 1. Buscar Projetos Go usando Fiber
```bash
./bateia -query "github.com/gofiber/fiber" -lang go -max 10 -output relatorios/fiber.md
```

#### 2. Buscar Projetos TypeScript usando NestJS
```bash
./bateia -query "@nestjs/core" -lang typescript -max 10 -output relatorios/nestjs.md
```

#### 3. Buscar Projetos Python usando FastAPI
```bash
./bateia -query "fastapi" -lang python -max 10 -output relatorios/fastapi.md
```

#### 4. Buscar Projetos Rust usando Axum
```bash
./bateia -query "axum" -lang rust -max 10 -output relatorios/axum.md
```

#### 5. Buscar Projetos Java usando Spring Boot
```bash
./bateia -query "spring-boot" -lang java -max 10 -output relatorios/spring.md
```

#### 6. Buscar Projetos PHP usando Laravel
```bash
./bateia -query "laravel" -lang php -max 10 -output relatorios/laravel.md
```

#### 7. Buscar Projetos C++ usando CMake / Boost
```bash
./bateia -query "boost" -lang cpp -max 10 -output relatorios/cpp.md
```

#### 8. Buscar Projetos C# / .NET usando Entity Framework
```bash
./bateia -query "EntityFrameworkCore" -lang csharp -max 10 -output relatorios/dotnet.md
```

---

## 🧪 **Garantia de Qualidade Tekt**

O projeto mantém **100% de aprovação nos testes** e conformidade rigorosa com a trava de arquitetura:

```bash
# Executar suíte completa de testes unitários
go test ./...

# Validar conformidade arquitetural Tekt
crystalline-lint .
```
