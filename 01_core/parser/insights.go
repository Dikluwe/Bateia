/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/insights.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type toolSpec struct {
	Key      string
	Patterns []string
}

var testingTools = []toolSpec{
	{Key: "testify", Patterns: []string{"testify"}},
	{Key: "testcontainers", Patterns: []string{"testcontainers-go"}},
	{Key: "pgtestdb", Patterns: []string{"pgtestdb"}},
	{Key: "miniredis", Patterns: []string{"miniredis"}},
	{Key: "go-sqlmock", Patterns: []string{"sqlmock"}},
	{Key: "mockery", Patterns: []string{"mockery"}},
}

var configTools = []toolSpec{
	{Key: "godotenv", Patterns: []string{"godotenv"}},
	{Key: "viper", Patterns: []string{"viper"}},
	{Key: "envconfig", Patterns: []string{"envconfig"}},
	{Key: "fsnotify", Patterns: []string{"fsnotify"}},
	{Key: "koanf", Patterns: []string{"koanf"}},
}

var databaseTools = []toolSpec{
	{Key: "pgx", Patterns: []string{"pgx"}},
	{Key: "goose", Patterns: []string{"goose"}},
	{Key: "bun", Patterns: []string{"bun"}},
	{Key: "ent", Patterns: []string{"ent"}},
	{Key: "gorm", Patterns: []string{"gorm"}},
	{Key: "sqlx", Patterns: []string{"sqlx"}},
	{Key: "mysql-driver", Patterns: []string{"go-sql-driver/mysql"}},
	{Key: "redis", Patterns: []string{"go-redis"}},
	{Key: "mongo", Patterns: []string{"mongo-driver"}},
}

var observabilityTools = []toolSpec{
	{Key: "zerolog", Patterns: []string{"zerolog"}},
	{Key: "zap", Patterns: []string{"zap"}},
	{Key: "opentelemetry", Patterns: []string{"opentelemetry", "otel"}},
	{Key: "grpc", Patterns: []string{"grpc"}},
	{Key: "prometheus", Patterns: []string{"prometheus"}},
	{Key: "asynq", Patterns: []string{"asynq"}},
}

var devopsFiles = []toolSpec{
	{Key: "GitHub Actions", Patterns: []string{".github/workflows", "workflows"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "Makefile", Patterns: []string{"Makefile"}},
	{Key: "golangci-lint", Patterns: []string{".golangci", "golangci"}},
	{Key: "docker-compose", Patterns: []string{"docker-compose", "compose.yml"}},
}

var typeScriptTestingTools = []toolSpec{
	{Key: "jest", Patterns: []string{"jest"}},
	{Key: "vitest", Patterns: []string{"vitest"}},
	{Key: "mocha", Patterns: []string{"mocha"}},
	{Key: "testing-library", Patterns: []string{"@testing-library"}},
	{Key: "cypress", Patterns: []string{"cypress"}},
	{Key: "playwright", Patterns: []string{"@playwright/test"}},
}

var typeScriptConfigTools = []toolSpec{
	{Key: "dotenv", Patterns: []string{"dotenv"}},
	{Key: "nestjs-config", Patterns: []string{"@nestjs/config"}},
	{Key: "config", Patterns: []string{"config"}},
	{Key: "convict", Patterns: []string{"convict"}},
	{Key: "envalid", Patterns: []string{"envalid"}},
}

var typeScriptDatabaseTools = []toolSpec{
	{Key: "prisma", Patterns: []string{"prisma", "@prisma/client"}},
	{Key: "drizzle-orm", Patterns: []string{"drizzle-orm"}},
	{Key: "typeorm", Patterns: []string{"typeorm"}},
	{Key: "mongoose", Patterns: []string{"mongoose"}},
	{Key: "sequelize", Patterns: []string{"sequelize"}},
	{Key: "knex", Patterns: []string{"knex"}},
}

var typeScriptObservabilityTools = []toolSpec{
	{Key: "winston", Patterns: []string{"winston"}},
	{Key: "pino", Patterns: []string{"pino"}},
	{Key: "morgan", Patterns: []string{"morgan"}},
	{Key: "axios", Patterns: []string{"axios"}},
	{Key: "opentelemetry", Patterns: []string{"@opentelemetry/api"}},
	{Key: "nestjs-core", Patterns: []string{"@nestjs/core"}},
	{Key: "express", Patterns: []string{"express"}},
	{Key: "fastify", Patterns: []string{"fastify"}},
}

var typeScriptDevopsFiles = []toolSpec{
	{Key: "eslint", Patterns: []string{"eslint", ".eslintrc"}},
	{Key: "prettier", Patterns: []string{"prettier", ".prettierrc"}},
	{Key: "husky", Patterns: []string{"husky"}},
	{Key: "lint-staged", Patterns: []string{"lint-staged"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var pythonTestingTools = []toolSpec{
	{Key: "pytest", Patterns: []string{"pytest"}},
	{Key: "unittest", Patterns: []string{"unittest"}},
	{Key: "coverage", Patterns: []string{"coverage"}},
	{Key: "httpx", Patterns: []string{"httpx"}},
	{Key: "factory-boy", Patterns: []string{"factory-boy", "factory_boy"}},
}

var pythonConfigTools = []toolSpec{
	{Key: "pydantic-settings", Patterns: []string{"pydantic-settings", "pydantic"}},
	{Key: "python-dotenv", Patterns: []string{"python-dotenv", "dotenv"}},
	{Key: "dynaconf", Patterns: []string{"dynaconf"}},
}

var pythonDatabaseTools = []toolSpec{
	{Key: "sqlalchemy", Patterns: []string{"sqlalchemy"}},
	{Key: "alembic", Patterns: []string{"alembic"}},
	{Key: "django-orm", Patterns: []string{"django"}},
	{Key: "tortoise-orm", Patterns: []string{"tortoise-orm"}},
	{Key: "peewee", Patterns: []string{"peewee"}},
	{Key: "psycopg2", Patterns: []string{"psycopg2", "psycopg"}},
	{Key: "asyncpg", Patterns: []string{"asyncpg"}},
	{Key: "redis", Patterns: []string{"redis"}},
	{Key: "motor", Patterns: []string{"motor", "pymongo"}},
}

var pythonObservabilityTools = []toolSpec{
	{Key: "fastapi", Patterns: []string{"fastapi"}},
	{Key: "flask", Patterns: []string{"flask"}},
	{Key: "celery", Patterns: []string{"celery"}},
	{Key: "structlog", Patterns: []string{"structlog"}},
	{Key: "sentry-sdk", Patterns: []string{"sentry-sdk"}},
	{Key: "opentelemetry", Patterns: []string{"opentelemetry"}},
	{Key: "prometheus-client", Patterns: []string{"prometheus-client"}},
}

var pythonDevopsFiles = []toolSpec{
	{Key: "ruff", Patterns: []string{"ruff", ".ruff.toml"}},
	{Key: "black", Patterns: []string{"black"}},
	{Key: "flake8", Patterns: []string{"flake8", ".flake8"}},
	{Key: "mypy", Patterns: []string{"mypy", "mypy.ini"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var rustTestingTools = []toolSpec{
	{Key: "criterion", Patterns: []string{"criterion"}},
	{Key: "tokio-test", Patterns: []string{"tokio-test"}},
	{Key: "proptest", Patterns: []string{"proptest"}},
	{Key: "rstest", Patterns: []string{"rstest"}},
	{Key: "wiremock", Patterns: []string{"wiremock"}},
}

var rustConfigTools = []toolSpec{
	{Key: "config-rs", Patterns: []string{"config"}},
	{Key: "dotenvy", Patterns: []string{"dotenvy", "dotenv"}},
	{Key: "clap", Patterns: []string{"clap"}},
	{Key: "figment", Patterns: []string{"figment"}},
}

var rustDatabaseTools = []toolSpec{
	{Key: "sqlx", Patterns: []string{"sqlx"}},
	{Key: "diesel", Patterns: []string{"diesel"}},
	{Key: "sea-orm", Patterns: []string{"sea-orm"}},
	{Key: "tokio-postgres", Patterns: []string{"tokio-postgres"}},
	{Key: "redis", Patterns: []string{"redis"}},
	{Key: "mongodb", Patterns: []string{"mongodb"}},
}

var rustObservabilityTools = []toolSpec{
	{Key: "tokio", Patterns: []string{"tokio"}},
	{Key: "axum", Patterns: []string{"axum"}},
	{Key: "actix-web", Patterns: []string{"actix-web"}},
	{Key: "tracing", Patterns: []string{"tracing", "tracing-subscriber"}},
	{Key: "serde", Patterns: []string{"serde"}},
	{Key: "reqwest", Patterns: []string{"reqwest"}},
	{Key: "opentelemetry", Patterns: []string{"opentelemetry"}},
}

var rustDevopsFiles = []toolSpec{
	{Key: "clippy", Patterns: []string{"clippy", "clippy.toml"}},
	{Key: "cargo-deny", Patterns: []string{"deny.toml"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var javaTestingTools = []toolSpec{
	{Key: "junit", Patterns: []string{"junit"}},
	{Key: "mockito", Patterns: []string{"mockito"}},
	{Key: "kotest", Patterns: []string{"kotest"}},
	{Key: "assertj", Patterns: []string{"assertj"}},
	{Key: "testcontainers", Patterns: []string{"testcontainers"}},
}

var javaConfigTools = []toolSpec{
	{Key: "spring-boot-starter", Patterns: []string{"spring-boot-starter"}},
	{Key: "archaius", Patterns: []string{"archaius"}},
	{Key: "jackson", Patterns: []string{"jackson"}},
}

var javaDatabaseTools = []toolSpec{
	{Key: "hibernate", Patterns: []string{"hibernate"}},
	{Key: "spring-data", Patterns: []string{"spring-data"}},
	{Key: "flyway", Patterns: []string{"flyway"}},
	{Key: "liquibase", Patterns: []string{"liquibase"}},
	{Key: "r2dbc", Patterns: []string{"r2dbc"}},
	{Key: "mybatis", Patterns: []string{"mybatis"}},
	{Key: "h2", Patterns: []string{"h2"}},
}

var javaObservabilityTools = []toolSpec{
	{Key: "spring-web", Patterns: []string{"spring-boot-starter-web", "spring-web"}},
	{Key: "quarkus", Patterns: []string{"quarkus"}},
	{Key: "micronaut", Patterns: []string{"micronaut"}},
	{Key: "slf4j", Patterns: []string{"slf4j"}},
	{Key: "logback", Patterns: []string{"logback"}},
	{Key: "micrometer", Patterns: []string{"micrometer"}},
}

var javaDevopsFiles = []toolSpec{
	{Key: "spotless", Patterns: []string{"spotless"}},
	{Key: "checkstyle", Patterns: []string{"checkstyle"}},
	{Key: "detekt", Patterns: []string{"detekt"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var phpTestingTools = []toolSpec{
	{Key: "phpunit", Patterns: []string{"phpunit"}},
	{Key: "pest", Patterns: []string{"pest"}},
	{Key: "mockery", Patterns: []string{"mockery"}},
	{Key: "codeception", Patterns: []string{"codeception"}},
	{Key: "behat", Patterns: []string{"behat"}},
}

var phpConfigTools = []toolSpec{
	{Key: "vlucas/phpdotenv", Patterns: []string{"vlucas/phpdotenv"}},
	{Key: "symfony/dotenv", Patterns: []string{"symfony/dotenv"}},
}

var phpDatabaseTools = []toolSpec{
	{Key: "doctrine", Patterns: []string{"doctrine/orm", "doctrine/dbal"}},
	{Key: "eloquent", Patterns: []string{"illuminate/database"}},
	{Key: "propel", Patterns: []string{"propel"}},
}

var phpObservabilityTools = []toolSpec{
	{Key: "laravel", Patterns: []string{"laravel/framework"}},
	{Key: "symfony", Patterns: []string{"symfony/http-kernel", "symfony/framework-bundle"}},
	{Key: "monolog", Patterns: []string{"monolog/monolog"}},
	{Key: "guzzle", Patterns: []string{"guzzlehttp/guzzle"}},
	{Key: "sentry", Patterns: []string{"sentry/sentry"}},
}

var phpDevopsFiles = []toolSpec{
	{Key: "phpstan", Patterns: []string{"phpstan", "phpstan.neon"}},
	{Key: "psalm", Patterns: []string{"psalm", "psalm.xml"}},
	{Key: "php-cs-fixer", Patterns: []string{"php-cs-fixer", ".php-cs-fixer.php"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var cppTestingTools = []toolSpec{
	{Key: "gtest", Patterns: []string{"gtest", "googletest"}},
	{Key: "catch2", Patterns: []string{"catch2"}},
	{Key: "doctest", Patterns: []string{"doctest"}},
	{Key: "boost.test", Patterns: []string{"boost_test"}},
}

var cppConfigTools = []toolSpec{
	{Key: "fmt", Patterns: []string{"fmt"}},
	{Key: "spdlog", Patterns: []string{"spdlog"}},
	{Key: "boost", Patterns: []string{"boost"}},
}

var cppDatabaseTools = []toolSpec{
	{Key: "sqlite3", Patterns: []string{"sqlite", "sqlite3"}},
	{Key: "pqxx", Patterns: []string{"pqxx"}},
	{Key: "mysqlclient", Patterns: []string{"mysql"}},
	{Key: "hiredis", Patterns: []string{"hiredis"}},
}

var cppObservabilityTools = []toolSpec{
	{Key: "qt", Patterns: []string{"qt", "qt5", "qt6"}},
	{Key: "opencv", Patterns: []string{"opencv"}},
	{Key: "grpc", Patterns: []string{"grpc"}},
	{Key: "protobuf", Patterns: []string{"protobuf"}},
	{Key: "opentelemetry", Patterns: []string{"opentelemetry"}},
}

var cppDevopsFiles = []toolSpec{
	{Key: "clang-format", Patterns: []string{".clang-format"}},
	{Key: "clang-tidy", Patterns: []string{".clang-tidy"}},
	{Key: "CMakeLists", Patterns: []string{"CMakeLists.txt"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var csharpTestingTools = []toolSpec{
	{Key: "xunit", Patterns: []string{"xunit"}},
	{Key: "nunit", Patterns: []string{"nunit"}},
	{Key: "mstest", Patterns: []string{"mstest"}},
	{Key: "moq", Patterns: []string{"moq"}},
	{Key: "fluentassertions", Patterns: []string{"fluentassertions"}},
}

var csharpConfigTools = []toolSpec{
	{Key: "configuration", Patterns: []string{"microsoft.extensions.configuration"}},
	{Key: "options", Patterns: []string{"microsoft.extensions.options"}},
}

var csharpDatabaseTools = []toolSpec{
	{Key: "entity-framework", Patterns: []string{"microsoft.entityframeworkcore"}},
	{Key: "dapper", Patterns: []string{"dapper"}},
	{Key: "npgsql", Patterns: []string{"npgsql"}},
	{Key: "sqlclient", Patterns: []string{"microsoft.data.sqlclient"}},
	{Key: "mongodb", Patterns: []string{"mongodb.driver"}},
}

var csharpObservabilityTools = []toolSpec{
	{Key: "serilog", Patterns: []string{"serilog"}},
	{Key: "nlog", Patterns: []string{"nlog"}},
	{Key: "mediatr", Patterns: []string{"mediatr"}},
	{Key: "fluentvalidation", Patterns: []string{"fluentvalidation"}},
	{Key: "automapper", Patterns: []string{"automapper"}},
	{Key: "polly", Patterns: []string{"polly"}},
	{Key: "opentelemetry", Patterns: []string{"opentelemetry"}},
}

var csharpDevopsFiles = []toolSpec{
	{Key: "stylecop", Patterns: []string{"stylecop"}},
	{Key: "editorconfig", Patterns: []string{".editorconfig"}},
	{Key: "Dockerfile", Patterns: []string{"Dockerfile"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var juliaTestingTools = []toolSpec{
	{Key: "Test", Patterns: []string{"test"}},
}

var juliaConfigTools = []toolSpec{
	{Key: "ArgParse", Patterns: []string{"argparse"}},
	{Key: "Configurations", Patterns: []string{"configurations"}},
}

var juliaDatabaseTools = []toolSpec{
	{Key: "LibPQ", Patterns: []string{"libpq"}},
	{Key: "SQLite", Patterns: []string{"sqlite"}},
	{Key: "MySQL", Patterns: []string{"mysql"}},
}

var juliaObservabilityTools = []toolSpec{
	{Key: "Flux", Patterns: []string{"flux"}},
	{Key: "Plots", Patterns: []string{"plots"}},
	{Key: "HTTP", Patterns: []string{"http"}},
}

var juliaDevopsFiles = []toolSpec{
	{Key: "Julia", Patterns: []string{"Project.toml"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var rTestingTools = []toolSpec{
	{Key: "testthat", Patterns: []string{"testthat"}},
	{Key: "covr", Patterns: []string{"covr"}},
}

var rConfigTools = []toolSpec{
	{Key: "config", Patterns: []string{"config"}},
}

var rDatabaseTools = []toolSpec{
	{Key: "DBI", Patterns: []string{"dbi"}},
	{Key: "RSQLite", Patterns: []string{"rsqlite"}},
	{Key: "RPostgres", Patterns: []string{"rpostgres"}},
}

var rObservabilityTools = []toolSpec{
	{Key: "shiny", Patterns: []string{"shiny"}},
	{Key: "ggplot2", Patterns: []string{"ggplot2"}},
}

var rDevopsFiles = []toolSpec{
	{Key: "styler", Patterns: []string{"styler"}},
	{Key: "GitHub Actions", Patterns: []string{".github/workflows"}},
}

var luaTestingTools = []toolSpec{
	{Key: "busted", Patterns: []string{"busted"}},
}

var luaConfigTools = []toolSpec{
	{Key: "penlight", Patterns: []string{"penlight"}},
}

var luaDatabaseTools = []toolSpec{
	{Key: "luasql", Patterns: []string{"luasql"}},
}

var luaObservabilityTools = []toolSpec{
	{Key: "lapis", Patterns: []string{"lapis"}},
}

var luaDevopsFiles = []toolSpec{
	{Key: "luacheck", Patterns: []string{".luacheckrc"}},
}

var ocamlTestingTools = []toolSpec{
	{Key: "alcotest", Patterns: []string{"alcotest"}},
}

var ocamlConfigTools = []toolSpec{
	{Key: "cmdliner", Patterns: []string{"cmdliner"}},
}

var ocamlDatabaseTools = []toolSpec{
	{Key: "caqti", Patterns: []string{"caqti"}},
}

var ocamlObservabilityTools = []toolSpec{
	{Key: "dream", Patterns: []string{"dream"}},
	{Key: "cohttp", Patterns: []string{"cohttp"}},
}

var ocamlDevopsFiles = []toolSpec{
	{Key: "dune", Patterns: []string{"dune-project"}},
}

// ExtractArchitecturalInsights analisa a lista de repositórios extraídos para calcular ocorrências nas 5 categorias de necessidades ocultas.
func ExtractArchitecturalInsights(results []domain.ExtractedRepository) domain.ArchitecturalInsights {
	totalRepos := len(results)
	if totalRepos == 0 {
		return domain.ArchitecturalInsights{
			Testing:       []domain.PatternMatch{},
			Configuration: []domain.PatternMatch{},
			Database:      []domain.PatternMatch{},
			Observability: []domain.PatternMatch{},
			DevOps:        []domain.FilePatternMatch{},
		}
	}

	// Detectar linguagem pré-dominante dos resultados
	lang := "go"
	for _, repo := range results {
		if repo.Manifest != nil && repo.Manifest.Language != "" {
			lang = repo.Manifest.Language
			break
		}
	}

	tTools := testingTools
	cTools := configTools
	dTools := databaseTools
	oTools := observabilityTools
	devOps := devopsFiles

	if lang == "typescript" || lang == "javascript" || lang == "ts" || lang == "node" {
		tTools = typeScriptTestingTools
		cTools = typeScriptConfigTools
		dTools = typeScriptDatabaseTools
		oTools = typeScriptObservabilityTools
		devOps = typeScriptDevopsFiles
	} else if lang == "python" || lang == "py" {
		tTools = pythonTestingTools
		cTools = pythonConfigTools
		dTools = pythonDatabaseTools
		oTools = pythonObservabilityTools
		devOps = pythonDevopsFiles
	} else if lang == "rust" || lang == "rs" {
		tTools = rustTestingTools
		cTools = rustConfigTools
		dTools = rustDatabaseTools
		oTools = rustObservabilityTools
		devOps = rustDevopsFiles
	} else if lang == "java" || lang == "kotlin" || lang == "kt" {
		tTools = javaTestingTools
		cTools = javaConfigTools
		dTools = javaDatabaseTools
		oTools = javaObservabilityTools
		devOps = javaDevopsFiles
	} else if lang == "php" {
		tTools = phpTestingTools
		cTools = phpConfigTools
		dTools = phpDatabaseTools
		oTools = phpObservabilityTools
		devOps = phpDevopsFiles
	} else if lang == "cpp" || lang == "c" || lang == "c++" {
		tTools = cppTestingTools
		cTools = cppConfigTools
		dTools = cppDatabaseTools
		oTools = cppObservabilityTools
		devOps = cppDevopsFiles
	} else if lang == "csharp" || lang == "cs" || lang == "dotnet" {
		tTools = csharpTestingTools
		cTools = csharpConfigTools
		dTools = csharpDatabaseTools
		oTools = csharpObservabilityTools
		devOps = csharpDevopsFiles
	} else if lang == "julia" {
		tTools = juliaTestingTools
		cTools = juliaConfigTools
		dTools = juliaDatabaseTools
		oTools = juliaObservabilityTools
		devOps = juliaDevopsFiles
	} else if lang == "r" {
		tTools = rTestingTools
		cTools = rConfigTools
		dTools = rDatabaseTools
		oTools = rObservabilityTools
		devOps = rDevopsFiles
	} else if lang == "lua" {
		tTools = luaTestingTools
		cTools = luaConfigTools
		dTools = luaDatabaseTools
		oTools = luaObservabilityTools
		devOps = luaDevopsFiles
	} else if lang == "ocaml" {
		tTools = ocamlTestingTools
		cTools = ocamlConfigTools
		dTools = ocamlDatabaseTools
		oTools = ocamlObservabilityTools
		devOps = ocamlDevopsFiles
	}

	insights := domain.ArchitecturalInsights{
		Testing:       analyzeCategory(results, tTools, totalRepos),
		Configuration: analyzeCategory(results, cTools, totalRepos),
		Database:      analyzeCategory(results, dTools, totalRepos),
		Observability: analyzeCategory(results, oTools, totalRepos),
		DevOps:        analyzeDevOps(results, devOps),
	}

	return insights
}

func analyzeCategory(results []domain.ExtractedRepository, specs []toolSpec, totalRepos int) []domain.PatternMatch {
	matches := make([]domain.PatternMatch, 0)

	for _, spec := range specs {
		foundRepos := make([]string, 0)
		versionSet := make(map[string]bool)

		for _, repo := range results {
			if repo.Manifest == nil {
				continue
			}

			repoName := repo.Candidate.Owner + "/" + repo.Candidate.Repo
			repoMatched := false

			for _, dep := range repo.Manifest.Dependencies {
				for _, pattern := range spec.Patterns {
					if strings.Contains(strings.ToLower(dep.Name), strings.ToLower(pattern)) {
						if !repoMatched {
							foundRepos = append(foundRepos, repoName)
							repoMatched = true
						}
						if dep.Version != "" {
							versionSet[dep.Version] = true
						}
					}
				}
			}
		}

		if len(foundRepos) > 0 {
			versions := make([]string, 0, len(versionSet))
			for v := range versionSet {
				versions = append(versions, v)
			}
			sort.Strings(versions)

			pct := (float64(len(foundRepos)) / float64(totalRepos)) * 100.0

			matches = append(matches, domain.PatternMatch{
				Name:       spec.Key,
				Count:      len(foundRepos),
				Percentage: pct,
				FoundIn:    foundRepos,
				Versions:   versions,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Count > matches[j].Count
	})

	return matches
}

func analyzeDevOps(results []domain.ExtractedRepository, specs []toolSpec) []domain.FilePatternMatch {
	matches := make([]domain.FilePatternMatch, 0)

	for _, spec := range specs {
		foundRepos := make([]string, 0)

		for _, repo := range results {
			repoName := repo.Candidate.Owner + "/" + repo.Candidate.Repo
			repoMatched := false

			// Verificar arquivos extraídos ou nos da árvore
			for _, file := range repo.Files {
				for _, pattern := range spec.Patterns {
					if strings.Contains(strings.ToLower(file.Path), strings.ToLower(pattern)) {
						foundRepos = append(foundRepos, repoName)
						repoMatched = true
						break
					}
				}
				if repoMatched {
					break
				}
			}
		}

		if len(foundRepos) > 0 {
			matches = append(matches, domain.FilePatternMatch{
				Name:    spec.Key,
				Count:   len(foundRepos),
				FoundIn: foundRepos,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Count > matches[j].Count
	})

	return matches
}
