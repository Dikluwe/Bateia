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
