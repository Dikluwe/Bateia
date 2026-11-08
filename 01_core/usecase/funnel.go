/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"errors"
	"fmt"
	pathpkg "path"
	"sort"
	"strings"
	"sync"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// SearchFunnel UseCase que implementa a orquestração do funil de redução de 3 estágios em L1.
type SearchFunnel struct {
	runtime  domain.ExecutionRuntime
	semantic domain.SemanticProvider
	searcher domain.SearchProvider
	fetcher  domain.ContentProvider
	store    domain.CacheStore
	reporter domain.ProgressReporter
	mu       sync.Mutex
}

// WithSemanticProvider injects bounded static analysis without adding I/O to the core.
func (f *SearchFunnel) WithSemanticProvider(provider domain.SemanticProvider) *SearchFunnel {
	f.semantic = provider
	return f
}

// WithProgressReporter configura um observador de progresso e mantém o núcleo sem I/O direto.
func (f *SearchFunnel) WithProgressReporter(reporter domain.ProgressReporter) *SearchFunnel {
	f.reporter = reporter
	return f
}

func (f *SearchFunnel) WithExecutionRuntime(runtime domain.ExecutionRuntime) *SearchFunnel {
	f.runtime = runtime
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
		sem <- struct{}{}
		index := i
		wg.Go(func() {
			defer func() { <-sem }()
			fn(index)
		})
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
	ctx = domain.WithExecutionRuntime(ctx, f.runtime)
	if err := ValidateQueryFilter(filter); err != nil {
		return nil, fmt.Errorf("filtro inválido: %w", err)
	}
	startTotal := domain.RuntimeNow(ctx)

	langs := ParseLanguageList(filter.Language)
	normalizedLanguage := strings.ToLower(strings.TrimSpace(filter.Language))
	if len(langs) > 1 && (strings.Contains(normalizedLanguage, ",") || normalizedLanguage == "sweep" || normalizedLanguage == "all" || normalizedLanguage == "polyglot") {
		f.progress("polyglot", "iniciando varredura em %d ecossistemas (%s)", len(langs), strings.Join(langs, ", "))
		var mergedResults []domain.ExtractedRepository
		var topCandidates, triagedCount, discardedCount int
		var failures []domain.ExecutionFailure
		var warnings []domain.OperationalWarning
		var attempts []domain.SearchAttempt
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
				attempts = append(attempts, subReport.SearchAttempts...)
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
			TotalDurationMs:      domain.ElapsedMillis(domain.RuntimeNow(ctx), startTotal),
			CacheHits:            cacheHits,
			CacheMisses:          cacheMisses,
		}

		finalReport := &domain.SearchReport{
			Query:           filter,
			GeneratedAt:     domain.RuntimeNow(ctx),
			AnalyzerVersion: parser.EvidenceVersion,
			TopCandidates:   topCandidates,
			TriagedCount:    triagedCount,
			ExtractedCount:  len(mergedResults),
			DiscardedCount:  discardedCount,
			Results:         mergedResults,
			Telemetry:       telemetry,
			Failures:        failures,
			Warnings:        warnings,
			SearchAttempts:  attempts,
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
	rep.Telemetry.TotalDurationMs = domain.ElapsedMillis(domain.RuntimeNow(ctx), startTotal)
	return rep, nil
}

func (f *SearchFunnel) executeRecursive(ctx context.Context, filter domain.QueryFilter, visited map[string]bool) (*domain.SearchReport, error) {
	startStage1 := domain.RuntimeNow(ctx)
	telemetry := &domain.FunnelTelemetry{}

	report := &domain.SearchReport{
		Query:           filter,
		GeneratedAt:     domain.RuntimeNow(ctx),
		AnalyzerVersion: parser.EvidenceVersion,
		Results:         make([]domain.ExtractedRepository, 0),
		Telemetry:       telemetry,
	}

	queryKey := fmt.Sprintf("%s_%s_%s_%d", parser.EvidenceVersion, filter.Query, filter.Language, filter.MaxResults)

	// ==========================================
	// Estágio 1: Descoberta (Topo do Funil)
	// ==========================================
	var candidates []domain.RepositoryCandidate
	var hit bool
	var err error

	if filter.LocalPath != "" {
		f.progress("discovery", "analisando worktree local %s", filter.LocalPath)
		candidates = append(candidates, domain.RepositoryCandidate{
			Owner: "local", Repo: pathpkg.Base(strings.ReplaceAll(filter.LocalPath, "\\", "/")), URL: "file://" + filter.LocalPath,
			Description: "Worktree local", CommitID: "WORKTREE", ConsultadoEm: domain.RuntimeNow(ctx),
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
					ConsultadoEm: domain.RuntimeNow(ctx),
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
			if traced, ok := f.searcher.(interface {
				SearchCandidatesWithTrace(context.Context, domain.QueryFilter) ([]domain.RepositoryCandidate, []domain.SearchAttempt, error)
			}); ok {
				candidates, report.SearchAttempts, err = traced.SearchCandidatesWithTrace(ctx, filter)
			} else {
				candidates, err = f.searcher.SearchCandidates(ctx, filter)
				report.SearchAttempts = append(report.SearchAttempts, domain.SearchAttempt{OriginalQuery: filter.Query, Provider: "não informado", Relaxation: "provedor sem rastreio; consulta enviada não verificada", Candidates: len(candidates)})
			}
			if err != nil {
				return nil, fmt.Errorf("falha no estágio 1 (Descoberta): %w", err)
			}
			if len(candidates) > 0 {
				f.cacheWarning(report, "discovery", "save", queryKey, f.store.SaveDiscovery(queryKey, candidates))
			}
		} else {
			telemetry.CacheHits++
			if len(candidates) > 0 && len(candidates[0].SearchAttempts) > 0 {
				report.SearchAttempts = append(report.SearchAttempts, candidates[0].SearchAttempts...)
			} else {
				report.SearchAttempts = append(report.SearchAttempts, domain.SearchAttempt{OriginalQuery: filter.Query, Provider: "cache", Relaxation: "consulta histórica ao provedor não registrada", Candidates: len(candidates)})
			}
			f.progress("discovery", "cache hit para %q (%d candidatos)", filter.Query, len(candidates))
		}
	}

	telemetry.DiscoveryDurationMs = domain.ElapsedMillis(domain.RuntimeNow(ctx), startStage1)

	report.TopCandidates = len(candidates)
	if len(candidates) == 0 {
		return report, nil
	}

	// ==========================================
	// Estágio 2: Triagem (Meio do Funil)
	// ==========================================
	startStage2 := domain.RuntimeNow(ctx)
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

	telemetry.TriageDurationMs = domain.ElapsedMillis(domain.RuntimeNow(ctx), startStage2)

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
	startStage3 := domain.RuntimeNow(ctx)
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

		// Selection depends on the query; legacy commit-only extraction cache is unsafe.
		var extracted *domain.ExtractedRepository
		implementationSearch := RetrieveImplementationEvidence(ctx, f.fetcher, tRepo.Candidate, filter.Query)
		targetPaths := SelectImplementationPaths(filter.Query, tRepo.Tree, implementationSearch, 12)
		{
			if filter.LocalPath == "" {
				f.recordCache(telemetry, false)
			}
			f.progress("extraction", "baixando conteúdo de %s (score %d)", ref.Key(), tRepo.RelevanceScore)

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
			// Reparse the bytes actually extracted so hashes identify analysis inputs.
			extracted.Manifest = nil
			for _, file := range extracted.Files {
				if parser.DetectManifestType(file.Path) != "" {
					if m, err := parser.ParseManifestAt(file.Path, extracted.Ref, file.Content, extracted.ConsultadoEm); err == nil {
						extracted.Manifest = m
						break
					}
				}
			}

			// Se o manifesto tiver sido recuperado ou existir, anexar README.md e reavaliar Relevância
			if extracted.Manifest != nil {
				for _, f := range extracted.Files {
					if strings.EqualFold(f.Path, "README.md") {
						extracted.Manifest.Content = []byte(f.Content)
						extracted.Manifest.ContentPath = f.Path
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

		}

		parser.AnnotateEvidence(extracted, filter.Query, targetPaths)
		implementationSearch = ValidateImplementationRequirements(filter.Query, implementationSearch, extracted.Files)
		extracted.ImplementationSearch = &implementationSearch
		if len(implementationSearch.Hits) > 0 {
			extracted.Coverage.Reason = "Definições Go localizadas por AST e releitura de bytes; pertinência heurística e comportamento não verificado."
			extracted.Coverage.ChecksRun = append(extracted.Coverage.ChecksRun, "recuperação Go por símbolos/conteúdo com identidade conferida")
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
		semantic := domain.SemanticReport{Status: "not_evaluated", Reason: "Análise Clang requer provedor e checkout local com compile_commands.json."}
		if f.semantic != nil && filter.LocalPath != "" {
			f.progress("semantic", "analisando unidades C/C++ de %s", ref.Key())
			var semanticErr error
			semantic, semanticErr = f.semantic.Analyze(ctx, domain.SemanticRequest{Root: filter.LocalPath, Query: filter.Query, Paths: nil})
			if semanticErr != nil {
				semantic.Status = "partial"
				semantic.Reason = semanticErr.Error()
			}
		}

		// Feed compiler-selected units back into textual extraction through the content port.
		var semanticPaths []string
		known := map[string]bool{}
		for _, file := range extracted.Files {
			known[file.Path] = true
		}
		root := strings.TrimSuffix(strings.TrimPrefix(extracted.Candidate.URL, "file://"), "/")
		for _, unit := range semantic.Units {
			p := strings.TrimPrefix(unit.Path, root+"/")
			if unit.Status == "complete" && p != unit.Path && !known[p] {
				semanticPaths = append(semanticPaths, p)
				known[p] = true
			}
		}
		if len(semanticPaths) > 0 {
			extra, err := f.fetcher.ExtractRepoFiles(ctx, tRepo.Candidate, semanticPaths)
			if err == nil && extra != nil {
				extracted.Files = append(extracted.Files, extra.Files...)
			} else {
				semantic.Reason += "; conteúdo das unidades não recuperado"
			}
			requested := append(append([]string(nil), extracted.Coverage.RequestedFiles...), semanticPaths...)
			parser.AnnotateEvidence(extracted, filter.Query, requested)
		}
		for i := range extracted.Files {
			file := &extracted.Files[i]
			for u := range semantic.Units {
				unit := &semantic.Units[u]
				for _, input := range unit.Inputs {
					normalizedRoot := strings.TrimSuffix(strings.TrimPrefix(extracted.Candidate.URL, "file://"), "/")
					inputPath := strings.TrimPrefix(strings.ReplaceAll(input.Path, "\\", "/"), normalizedRoot+"/")
					if inputPath == file.Path && file.Evidence.SHA256 != input.SHA256 {
						unit.Status = "partial"
						semantic.Status = "partial"
						unit.Diagnostics += "; bytes divergem da extração inicial"
					} else if inputPath == file.Path && unit.Status == "complete" {
						file.Evidence.State = "static_analyzed"
						file.Evidence.Method += "; Clang com contexto de compilação"
					}
				}
			}
		}
		extracted.Semantic = &semantic
		extracted.QueryEvidence = parser.BuildQueryEvidence(filter.Query, &semantic, extracted.Files, root)
		if len(semantic.Units) > 0 {
			extracted.Coverage.SupportedLanguages = append(extracted.Coverage.SupportedLanguages, "C/C++ (Clang; configuração explícita)")
			extracted.Coverage.ChecksRun = append(extracted.Coverage.ChecksRun, "tentativa Clang: "+semantic.Status)
		}

		extracted.ArchViolations = parser.EvaluateArchViolations(extracted.Manifest, tRepo.Tree)
		extracted.HealthScore = -1 // No validated architectural coverage metric.
		securityScan, advisories := parser.EvaluateSecurityScan(extracted.Manifest)
		extracted.SecurityScan = &securityScan
		extracted.SecurityAdvisories = advisories

		f.mu.Lock()
		report.Results = append(report.Results, *extracted)
		f.mu.Unlock()
	})
	sort.SliceStable(report.Results, func(i, j int) bool {
		return report.Results[i].RelevanceScore > report.Results[j].RelevanceScore
	})

	telemetry.ExtractionDurationMs = domain.ElapsedMillis(domain.RuntimeNow(ctx), startStage3)

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
				for i := range childReport.Results {
					res := &childReport.Results[i]
					rel := parser.EvaluateRelevance(filter.Query, res.Manifest, res.Candidate, nil)
					res.RelevanceScore, res.RelevanceLabel, res.RelevanceReason = rel.Score, rel.Label, rel.Reason
					parser.AnnotateEvidence(res, filter.Query, res.Coverage.RequestedFiles)
					res.QueryEvidence = parser.BuildQueryEvidence(filter.Query, res.Semantic, res.Files, strings.TrimSuffix(strings.TrimPrefix(res.Candidate.URL, "file://"), "/"))
				}
				report.SearchAttempts = append(report.SearchAttempts, childReport.SearchAttempts...)
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
