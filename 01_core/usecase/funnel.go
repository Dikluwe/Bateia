/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-08-04
 */

package usecase

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// SearchFunnel UseCase que implementa a orquestração do funil de redução de 3 estágios em L1.
type SearchFunnel struct {
	searcher domain.SearchProvider
	fetcher  domain.ContentProvider
	store    domain.CacheStore
	reporter domain.ProgressReporter
	mu       sync.Mutex
}

// WithProgressReporter configura um observador de progresso e mantém o núcleo sem I/O direto.
func (f *SearchFunnel) WithProgressReporter(reporter domain.ProgressReporter) *SearchFunnel {
	f.reporter = reporter
	return f
}

func (f *SearchFunnel) progress(stage, format string, args ...any) {
	if f.reporter != nil {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.reporter.ReportProgress(domain.ProgressEvent{Stage: stage, Message: fmt.Sprintf(format, args...)})
	}
}

func (f *SearchFunnel) cacheWarning(report *domain.SearchReport, stage, operation, target string, err error) {
	if err == nil {
		return
	}
	f.mu.Lock()
	report.Warnings = append(report.Warnings, domain.OperationalWarning{
		Stage: stage, Operation: operation, Target: target, Message: err.Error(),
	})
	f.mu.Unlock()
	f.progress(stage, "cache degradado em %s (%s): %v", operation, target, err)
}

func (f *SearchFunnel) recordCache(telemetry *domain.FunnelTelemetry, hit bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if hit {
		telemetry.CacheHits++
	} else {
		telemetry.CacheMisses++
	}
}

func parallelFor(ctx context.Context, count, limit int, fn func(int)) {
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(index int) {
			defer wg.Done()
			defer func() { <-sem }()
			fn(index)
		}(i)
	}
	wg.Wait()
}

func NewSearchFunnel(searcher domain.SearchProvider, fetcher domain.ContentProvider, store domain.CacheStore) *SearchFunnel {
	return &SearchFunnel{
		searcher: searcher,
		fetcher:  fetcher,
		store:    store,
	}
}

// Execute executa a busca através dos 3 estágios do funil com cache e invariante de redução.
func (f *SearchFunnel) Execute(ctx context.Context, filter domain.QueryFilter) (*domain.SearchReport, error) {
	if err := ValidateQueryFilter(filter); err != nil {
		return nil, fmt.Errorf("filtro inválido: %w", err)
	}
	startTotal := time.Now()

	langs := ParseLanguageList(filter.Language)
	normalizedLanguage := strings.ToLower(strings.TrimSpace(filter.Language))
	if len(langs) > 1 && (strings.Contains(normalizedLanguage, ",") || normalizedLanguage == "sweep" || normalizedLanguage == "all" || normalizedLanguage == "polyglot") {
		f.progress("polyglot", "iniciando varredura em %d ecossistemas (%s)", len(langs), strings.Join(langs, ", "))
		var mergedResults []domain.ExtractedRepository
		var topCandidates, triagedCount, discardedCount int
		var failures []domain.ExecutionFailure
		var warnings []domain.OperationalWarning
		var discoveryMs, triageMs, extractionMs int64
		var cacheHits, cacheMisses int
		visited := make(map[string]bool)

		for _, l := range langs {
			subFilter := filter
			subFilter.Language = l
			subReport, err := f.executeRecursive(ctx, subFilter, visited)
			if err != nil {
				failures = append(failures, domain.ExecutionFailure{Scope: l, Stage: "sweep", Message: err.Error()})
				continue
			}
			if subReport != nil {
				mergedResults = append(mergedResults, subReport.Results...)
				topCandidates += subReport.TopCandidates
				triagedCount += subReport.TriagedCount
				discardedCount += subReport.DiscardedCount
				failures = append(failures, subReport.Failures...)
				warnings = append(warnings, subReport.Warnings...)
				if subReport.Telemetry != nil {
					discoveryMs += subReport.Telemetry.DiscoveryDurationMs
					triageMs += subReport.Telemetry.TriageDurationMs
					extractionMs += subReport.Telemetry.ExtractionDurationMs
					cacheHits += subReport.Telemetry.CacheHits
					cacheMisses += subReport.Telemetry.CacheMisses
				}
			}
		}

		telemetry := &domain.FunnelTelemetry{
			DiscoveryDurationMs:  discoveryMs,
			TriageDurationMs:     triageMs,
			ExtractionDurationMs: extractionMs,
			TotalDurationMs:      time.Since(startTotal).Milliseconds(),
			CacheHits:            cacheHits,
			CacheMisses:          cacheMisses,
		}

		finalReport := &domain.SearchReport{
			Query:          filter,
			GeneratedAt:    time.Now().UTC(),
			TopCandidates:  topCandidates,
			TriagedCount:   triagedCount,
			ExtractedCount: len(mergedResults),
			DiscardedCount: discardedCount,
			Results:        mergedResults,
			Telemetry:      telemetry,
			Failures:       failures,
			Warnings:       warnings,
		}

		insights := parser.ExtractArchitecturalInsights(finalReport.Results)
		finalReport.Insights = &insights
		finalReport.PolyglotMatrix = parser.EvaluatePolyglotMatrix(finalReport.Results)

		if len(finalReport.Results) > 0 {
			targetOnt := ExtractOntologyFromIdea(filter.Query)
			finalReport.TargetOntology = &targetOnt

			var simResults []domain.IdeaSimilarityResult
			for _, res := range finalReport.Results {
				candOnt := ExtractOntologyFromRepo(res)
				sim := CalculateIdeaSimilarity(targetOnt, candOnt)
				simResults = append(simResults, sim)
			}
			finalReport.SimilarityResults = simResults
			finalReport.Recommendations = MineOntologicalImprovements(targetOnt, simResults, finalReport.Results)
			finalReport.LanguageDistribution = CalculateLanguageDistribution(finalReport.Results, simResults)
		}

		return finalReport, nil
	}

	visited := make(map[string]bool)
	rep, err := f.executeRecursive(ctx, filter, visited)
	if err != nil {
		return nil, err
	}

	if rep.Telemetry == nil {
		rep.Telemetry = &domain.FunnelTelemetry{}
	}
	rep.Telemetry.TotalDurationMs = time.Since(startTotal).Milliseconds()
	return rep, nil
}

func (f *SearchFunnel) executeRecursive(ctx context.Context, filter domain.QueryFilter, visited map[string]bool) (*domain.SearchReport, error) {
	startStage1 := time.Now()
	telemetry := &domain.FunnelTelemetry{}

	report := &domain.SearchReport{
		Query:       filter,
		GeneratedAt: time.Now().UTC(),
		Results:     make([]domain.ExtractedRepository, 0),
		Telemetry:   telemetry,
	}

	queryKey := fmt.Sprintf("%s_%s_%d", filter.Query, filter.Language, filter.MaxResults)

	// ==========================================
	// Estágio 1: Descoberta (Topo do Funil)
	// ==========================================
	var candidates []domain.RepositoryCandidate
	var hit bool
	var err error

	if filter.LocalPath != "" {
		f.progress("discovery", "analisando worktree local %s", filter.LocalPath)
		candidates = append(candidates, domain.RepositoryCandidate{
			Owner: "local", Repo: filepath.Base(filter.LocalPath), URL: "file://" + filter.LocalPath,
			Description: "Worktree local", CommitID: "WORKTREE", ConsultadoEm: time.Now().UTC(),
		})
	} else if len(filter.DirectRepos) > 0 {
		f.progress("discovery", "modo benchmark: injetando %d repositórios", len(filter.DirectRepos))
		for _, repoStr := range filter.DirectRepos {
			parts := strings.Split(repoStr, "/")
			if len(parts) == 2 {
				candidates = append(candidates, domain.RepositoryCandidate{
					Owner:        parts[0],
					Repo:         parts[1],
					URL:          "https://github.com/" + repoStr,
					Description:  "Repositório Injetado (Benchmark)",
					CommitID:     "HEAD",
					ConsultadoEm: time.Now().UTC(),
				})
			}
		}
	} else {
		f.progress("discovery", "buscando candidatos")
		candidates, hit, err = f.store.GetDiscovery(queryKey, 2*time.Hour)
		if err != nil {
			f.cacheWarning(report, "discovery", "get", queryKey, err)
			hit = false
		}

		if !hit || len(candidates) == 0 {
			telemetry.CacheMisses++
			f.progress("discovery", "cache vazio; consultando provedor")
			candidates, err = f.searcher.SearchCandidates(ctx, filter)
			if err != nil {
				return nil, fmt.Errorf("falha no estágio 1 (Descoberta): %w", err)
			}
			if len(candidates) > 0 {
				f.cacheWarning(report, "discovery", "save", queryKey, f.store.SaveDiscovery(queryKey, candidates))
			}
		} else {
			telemetry.CacheHits++
			f.progress("discovery", "cache hit para %q (%d candidatos)", filter.Query, len(candidates))
		}
	}

	telemetry.DiscoveryDurationMs = time.Since(startStage1).Milliseconds()

	report.TopCandidates = len(candidates)
	if len(candidates) == 0 {
		return report, nil
	}

	// ==========================================
	// Estágio 2: Triagem (Meio do Funil)
	// ==========================================
	startStage2 := time.Now()
	f.progress("triage", "triando manifestos e estrutura de %d candidatos", len(candidates))
	triaged := make([]domain.TriagedRepository, 0)
	discardedCount := 0

	parallelFor(ctx, len(candidates), 4, func(candidateIndex int) {
		cand := candidates[candidateIndex]
		ref := domain.CommitRef{
			Owner:  cand.Owner,
			Repo:   cand.Repo,
			Commit: cand.CommitID,
		}

		refKey := cand.Owner + "/" + cand.Repo
		f.mu.Lock()
		alreadyVisited := visited[refKey]
		if !alreadyVisited {
			visited[refKey] = true
		}
		f.mu.Unlock()
		if alreadyVisited {
			f.mu.Lock()
			discardedCount++
			f.mu.Unlock()
			return
		}

		// 1. Tentar manifesto via Cache
		var manifest *domain.Manifest
		var hitM bool
		if filter.LocalPath == "" {
			var cacheErr error
			manifest, hitM, cacheErr = f.store.GetManifest(ref)
			if cacheErr != nil {
				f.cacheWarning(report, "triage", "get_manifest", ref.Key(), cacheErr)
				hitM = false
			}
		}
		if hitM {
			f.recordCache(telemetry, true)
		} else {
			if filter.LocalPath == "" {
				f.recordCache(telemetry, false)
			}
			fetchedManifest, fetchErr := f.fetcher.FetchManifest(ctx, cand)
			if fetchErr != nil {
				if errors.Is(fetchErr, domain.ErrManifestNotFound) {
					cand.Stale = true
					f.mu.Lock()
					discardedCount++
					f.mu.Unlock()
					return
				}
				f.progress("triage", "falha ao buscar manifesto para %s: %v", ref.Key(), fetchErr)
			} else {
				manifest = fetchedManifest
				if filter.LocalPath == "" {
					f.cacheWarning(report, "triage", "save_manifest", ref.Key(), f.store.SaveManifest(ref, manifest))
				}
			}
		}

		// 2. Tentar estrutura de diretórios se filtros foram especificados
		tree, treeErr := f.fetcher.FetchTree(ctx, cand)
		score := parser.StructuralScore{Approved: true, Reason: "Sem filtro estrutural"}
		if treeErr == nil && len(filter.MustHave) > 0 {
			score = parser.EvaluateStructure(tree, filter.MustHave)
		}

		// 3. Avaliar Score de Relevância e Import Matcher
		rel := parser.EvaluateRelevance(filter.Query, manifest, cand, tree)

		if score.Approved {
			f.mu.Lock()
			triaged = append(triaged, domain.TriagedRepository{
				Candidate:       cand,
				Manifest:        manifest,
				Tree:            tree,
				Approved:        true,
				ScoreReason:     score.Reason,
				RelevanceScore:  rel.Score,
				RelevanceLabel:  rel.Label,
				RelevanceReason: rel.Reason,
			})
			f.mu.Unlock()
		} else {
			f.mu.Lock()
			discardedCount++
			f.mu.Unlock()
		}
	})

	telemetry.TriageDurationMs = time.Since(startStage2).Milliseconds()

	report.TriagedCount = len(triaged)
	report.DiscardedCount = discardedCount
	report.DiscardReason = "Filtro estrutural ou repositório indisponível (404/stale)"

	// Ordenar repositórios aprovados por RelevanceScore decrescente
	sort.SliceStable(triaged, func(i, j int) bool {
		return triaged[i].RelevanceScore > triaged[j].RelevanceScore
	})

	f.progress("triage", "%d aprovados de %d candidatos", len(triaged), len(candidates))

	// ==========================================
	// Estágio 3: Extração (Fundo do Funil)
	// Limita o fundo do funil a no máximo 10 repositórios validados
	// ==========================================
	startStage3 := time.Now()
	maxExtract := 10
	if len(triaged) < maxExtract {
		maxExtract = len(triaged)
	}

	f.progress("extraction", "extraindo até %d repositórios", maxExtract)

	parallelFor(ctx, maxExtract, 4, func(i int) {
		tRepo := triaged[i]

		// Enriquecer metadados (estrelas e descrição) via GitHub API
		enrichedCand, errMeta := f.fetcher.FetchRepoMetadata(ctx, tRepo.Candidate)
		if errMeta == nil {
			tRepo.Candidate = enrichedCand
		}

		ref := domain.CommitRef{
			Owner:  tRepo.Candidate.Owner,
			Repo:   tRepo.Candidate.Repo,
			Commit: tRepo.Candidate.CommitID,
		}

		// Tentar extração via cache imutável (repo@commit)
		var extracted *domain.ExtractedRepository
		var hitE bool
		if filter.LocalPath == "" {
			var cacheErr error
			extracted, hitE, cacheErr = f.store.GetExtracted(ref)
			if cacheErr != nil {
				f.cacheWarning(report, "extraction", "get_extracted", ref.Key(), cacheErr)
				hitE = false
			}
		}
		if hitE {
			f.recordCache(telemetry, true)
			extracted.Candidate = tRepo.Candidate
			if extracted.Manifest != nil && extracted.Manifest.Content != nil {
				rel := parser.EvaluateRelevance(filter.Query, extracted.Manifest, tRepo.Candidate, tRepo.Tree)
				extracted.RelevanceScore = rel.Score
				extracted.RelevanceLabel = rel.Label
				extracted.RelevanceReason = rel.Reason
			} else {
				extracted.RelevanceScore = tRepo.RelevanceScore
				extracted.RelevanceLabel = tRepo.RelevanceLabel
				extracted.RelevanceReason = tRepo.RelevanceReason
			}
			f.progress("extraction", "cache hit para %s (score %d)", ref.Key(), extracted.RelevanceScore)
		} else {
			if filter.LocalPath == "" {
				f.recordCache(telemetry, false)
			}
			f.progress("extraction", "baixando conteúdo de %s (score %d)", ref.Key(), tRepo.RelevanceScore)
			targetPaths := []string{"go.mod", "README.md", "package.json", "Cargo.toml"}

			// Selecionar dinamicamente os 3 primeiros arquivos de código fonte da árvore
			codeCount := 0
			for _, node := range tRepo.Tree {
				if node.Type == "file" || node.Type == "blob" {
					ext := strings.ToLower(filepath.Ext(node.Path))
					if (ext == ".go" || ext == ".rs" || ext == ".ts" || ext == ".py" || ext == ".java") && !strings.Contains(node.Path, "_test.go") {
						targetPaths = append(targetPaths, node.Path)
						codeCount++
						if codeCount >= 3 {
							break
						}
					}
				}
			}

			if codeCount == 0 {
				targetPaths = append(targetPaths, "gin.go", "main.go", "app.go", "lib.rs", "index.ts", "cmd/main.go", "types.go", "domain.go", "server.go", "client.go", "01_core/domain/types.go", "01_core/parser/ast.go")
			}

			downloaded, extractErr := f.fetcher.ExtractRepoFiles(ctx, tRepo.Candidate, targetPaths)
			if extractErr != nil {
				f.progress("extraction", "falha na extração de %s: %v", ref.Key(), extractErr)
				f.mu.Lock()
				report.Failures = append(report.Failures, domain.ExecutionFailure{Scope: ref.Key(), Stage: "extraction", Message: extractErr.Error()})
				f.mu.Unlock()
				return
			}
			extracted = downloaded
			extracted.Candidate = tRepo.Candidate
			if tRepo.Manifest != nil {
				extracted.Manifest = tRepo.Manifest
			} else {
				// Tentar recuperar manifesto a partir dos arquivos extraídos (ex: go.mod)
				for _, f := range extracted.Files {
					if f.Path == "go.mod" || f.Path == "package.json" || f.Path == "Cargo.toml" {
						m, errM := parser.ParseManifest(f.Path, ref, f.Content)
						if errM == nil {
							extracted.Manifest = m
							break
						}
					}
				}
			}

			// Se o manifesto tiver sido recuperado ou existir, anexar README.md e reavaliar Relevância
			if extracted.Manifest != nil {
				for _, f := range extracted.Files {
					if strings.EqualFold(f.Path, "README.md") {
						extracted.Manifest.Content = []byte(f.Content)
						break
					}
				}
				rel := parser.EvaluateRelevance(filter.Query, extracted.Manifest, tRepo.Candidate, tRepo.Tree)
				extracted.RelevanceScore = rel.Score
				extracted.RelevanceLabel = rel.Label
				extracted.RelevanceReason = rel.Reason
			} else {
				extracted.RelevanceScore = tRepo.RelevanceScore
				extracted.RelevanceLabel = tRepo.RelevanceLabel
				extracted.RelevanceReason = tRepo.RelevanceReason
			}

			if filter.LocalPath == "" {
				f.cacheWarning(report, "extraction", "save_extracted", ref.Key(), f.store.SaveExtracted(ref, extracted))
			}
		}

		// Executar Análise AST nos arquivos extraídos (seja HIT ou MISS)
		var allRepoSymbols []domain.ASTSymbol
		for idxFile, file := range extracted.Files {
			syms := parser.ExtractASTSymbols(file.Path, []byte(file.Content))
			if len(syms) > 0 {
				extracted.Files[idxFile].Symbols = syms
				allRepoSymbols = append(allRepoSymbols, syms...)
			}
		}
		extracted.Symbols = allRepoSymbols
		extracted.ArchViolations = parser.EvaluateArchViolations(extracted.Manifest, tRepo.Tree)
		extracted.HealthScore = parser.CalculateHealthScore(extracted.ArchViolations)
		extracted.SecurityAdvisories = parser.EvaluateSecurityAdvisories(extracted.Manifest)

		f.mu.Lock()
		report.Results = append(report.Results, *extracted)
		f.mu.Unlock()
	})
	sort.SliceStable(report.Results, func(i, j int) bool {
		return report.Results[i].RelevanceScore > report.Results[j].RelevanceScore
	})

	telemetry.ExtractionDurationMs = time.Since(startStage3).Milliseconds()

	report.PolyglotMatrix = parser.EvaluatePolyglotMatrix(report.Results)

	report.ExtractedCount = len(report.Results)

	// ==========================================
	// Estágio 4: Recursividade (Malha de Dependências)
	// ==========================================
	if filter.Depth > 1 {
		f.progress("recursion", "aprofundando busca; restam %d rodadas para %q", filter.Depth-1, filter.Query)

		depsToSearch := make([]string, 0)
		for _, res := range report.Results {
			if res.Manifest != nil {
				count := 0
				for _, dep := range res.Manifest.Dependencies {
					if !dep.Indirect && count < 3 { // Top 3 dependências diretas
						depKey := "dep:" + dep.Name
						if !visited[depKey] {
							visited[depKey] = true
							depsToSearch = append(depsToSearch, dep.Name)
							count++
						}
					}
				}
			}
		}

		for _, depName := range depsToSearch {
			childFilter := filter
			childFilter.Depth = filter.Depth - 1
			childFilter.Query = depName

			childReport, err := f.executeRecursive(ctx, childFilter, visited)
			if err != nil {
				report.Failures = append(report.Failures, domain.ExecutionFailure{Scope: depName, Stage: "recursion", Message: err.Error()})
				continue
			}
			if childReport != nil {
				report.Results = append(report.Results, childReport.Results...)
				report.TriagedCount += childReport.TriagedCount
				report.DiscardedCount += childReport.DiscardedCount
				report.TopCandidates += childReport.TopCandidates
				report.Failures = append(report.Failures, childReport.Failures...)
				report.Warnings = append(report.Warnings, childReport.Warnings...)
			}
		}
		report.ExtractedCount = len(report.Results)
	}

	// Extrair Necessidades Ocultas & Insights Arquiteturais dos repositórios validados
	insights := parser.ExtractArchitecturalInsights(report.Results)
	report.Insights = &insights

	// Estágio 5: Ontologia e Score de Similaridade de Ideias
	if len(report.Results) > 0 {
		var targetOnt domain.ProjectOntology
		if len(filter.DirectRepos) > 0 {
			for _, res := range report.Results {
				if res.Ref.Key() == filter.DirectRepos[0] || (res.Ref.Owner+"/"+res.Ref.Repo) == filter.DirectRepos[0] {
					targetOnt = ExtractOntologyFromRepo(res)
					break
				}
			}
		}
		if targetOnt.ProblemDomain == "" {
			targetOnt = ExtractOntologyFromIdea(filter.Query)
		}
		report.TargetOntology = &targetOnt

		var simResults []domain.IdeaSimilarityResult
		for _, res := range report.Results {
			candOnt := ExtractOntologyFromRepo(res)
			sim := CalculateIdeaSimilarity(targetOnt, candOnt)
			simResults = append(simResults, sim)
		}
		report.SimilarityResults = simResults
		report.Recommendations = MineOntologicalImprovements(targetOnt, simResults, report.Results)
		report.LanguageDistribution = CalculateLanguageDistribution(report.Results, simResults)
	}

	f.progress("completed", "%d repositórios entregues com insights computados", report.ExtractedCount)

	return report, nil
}
