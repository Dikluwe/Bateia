/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-07-25
 */

package usecase

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// SearchFunnel UseCase que implementa a orquestração do funil de redução de 3 estágios em L1.
type SearchFunnel struct {
	searcher domain.SearchProvider
	fetcher  domain.ContentProvider
	store    domain.CacheStore
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
	// Normalizar query antes de iniciar a busca
	if len(filter.DirectRepos) == 0 {
		normalized := parser.NormalizeQuery(filter.Query)
		if normalized != filter.Query {
			log.Printf("[QUERY NORMALIZER] Query '%s' sanitizada para '%s'", filter.Query, normalized)
			filter.Query = normalized
		}
	}

	visited := make(map[string]bool)
	return f.executeRecursive(ctx, filter, visited)
}

func (f *SearchFunnel) executeRecursive(ctx context.Context, filter domain.QueryFilter, visited map[string]bool) (*domain.SearchReport, error) {
	report := &domain.SearchReport{
		Query:       filter,
		GeneratedAt: time.Now().UTC(),
		Results:     make([]domain.ExtractedRepository, 0),
	}

	queryKey := fmt.Sprintf("%s_%s_%d", filter.Query, filter.Language, filter.MaxResults)

	// ==========================================
	// Estágio 1: Descoberta (Topo do Funil)
	// ==========================================
	var candidates []domain.RepositoryCandidate
	var hit bool
	var err error

	if len(filter.DirectRepos) > 0 {
		log.Printf("[ESTÁGIO 1] Modo Benchmark: Injetando %d repositórios diretamente...", len(filter.DirectRepos))
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
		log.Printf("[ESTÁGIO 1] Buscando candidatos no topo do funil...")
		candidates, hit, err = f.store.GetDiscovery(queryKey, 2*time.Hour)
		if err != nil {
			log.Printf("[WARNING] Erro ao ler cache de descoberta: %v", err)
		}

		if !hit || len(candidates) == 0 {
			log.Printf("[ESTÁGIO 1] Cache miss/vazio. Consultando provedor de busca...")
			candidates, err = f.searcher.SearchCandidates(ctx, filter)
			if err != nil {
				return nil, fmt.Errorf("falha no estágio 1 (Descoberta): %w", err)
			}
			if len(candidates) > 0 {
				_ = f.store.SaveDiscovery(queryKey, candidates)
			}
		} else {
			log.Printf("[ESTÁGIO 1] Cache HIT para a busca '%s' (%d candidatos).", filter.Query, len(candidates))
		}
	}

	report.TopCandidates = len(candidates)
	if len(candidates) == 0 {
		return report, nil
	}

	// ==========================================
	// Estágio 2: Triagem (Meio do Funil)
	// ==========================================
	log.Printf("[ESTÁGIO 2] Triando manifestos e estrutura (%d candidatos)...", len(candidates))
	triaged := make([]domain.TriagedRepository, 0)
	discardedCount := 0

	for _, cand := range candidates {
		ref := domain.CommitRef{
			Owner:  cand.Owner,
			Repo:   cand.Repo,
			Commit: cand.CommitID,
		}
		
		refKey := cand.Owner + "/" + cand.Repo
		if visited[refKey] {
			discardedCount++
			continue
		}
		visited[refKey] = true

		// 1. Tentar manifesto via Cache
		manifest, hitM, _ := f.store.GetManifest(ref)
		if !hitM {
			manifest, err = f.fetcher.FetchManifest(ctx, cand)
			if err != nil {
				if err.Error() == "MANIFEST_NOT_FOUND (404/410)" {
					cand.Stale = true
					discardedCount++
					continue
				}
				log.Printf("[WARNING] Falha ao buscar manifesto para %s: %v", ref.Key(), err)
			} else {
				_ = f.store.SaveManifest(ref, manifest)
			}
		}

		// 2. Tentar estrutura de diretórios se filtros foram especificados
		tree, err := f.fetcher.FetchTree(ctx, cand)
		score := parser.StructuralScore{Approved: true, Reason: "Sem filtro estrutural"}
		if err == nil && len(filter.MustHave) > 0 {
			score = parser.EvaluateStructure(tree, filter.MustHave)
		}

		// 3. Avaliar Score de Relevância e Import Matcher
		rel := parser.EvaluateRelevance(filter.Query, manifest, cand, tree)

		if score.Approved {
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
		} else {
			discardedCount++
		}
	}

	report.TriagedCount = len(triaged)
	report.DiscardedCount = discardedCount
	report.DiscardReason = "Filtro estrutural ou repositório indisponível (404/stale)"

	// Ordenar repositórios aprovados por RelevanceScore decrescente
	sort.SliceStable(triaged, func(i, j int) bool {
		return triaged[i].RelevanceScore > triaged[j].RelevanceScore
	})

	log.Printf("[ESTÁGIO 2] Triagem concluída: %d aprovados de %d candidatos (ordenados por relevância).", len(triaged), len(candidates))

	// ==========================================
	// Estágio 3: Extração (Fundo do Funil)
	// Limita o fundo do funil a no máximo 10 repositórios validados
	// ==========================================
	maxExtract := 10
	if len(triaged) < maxExtract {
		maxExtract = len(triaged)
	}

	log.Printf("[ESTÁGIO 3] Extraindo até %d repositórios mais relevantes...", maxExtract)

	for i := 0; i < maxExtract; i++ {
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
		extracted, hitE, _ := f.store.GetExtracted(ref)
		if !hitE {
			log.Printf("[ESTÁGIO 3] Baixando conteúdo para %s (Score: %d)...", ref.Key(), tRepo.RelevanceScore)
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

			extracted, err = f.fetcher.ExtractRepoFiles(ctx, tRepo.Candidate, targetPaths)
			if err != nil {
				log.Printf("[WARNING] Falha na extração de %s: %v", ref.Key(), err)
				continue
			}
			extracted.Candidate = tRepo.Candidate
			extracted.Manifest = tRepo.Manifest
			extracted.RelevanceScore = tRepo.RelevanceScore
			extracted.RelevanceLabel = tRepo.RelevanceLabel
			extracted.RelevanceReason = tRepo.RelevanceReason
			_ = f.store.SaveExtracted(ref, extracted)
		} else {
			extracted.Candidate = tRepo.Candidate
			extracted.RelevanceScore = tRepo.RelevanceScore
			extracted.RelevanceLabel = tRepo.RelevanceLabel
			extracted.RelevanceReason = tRepo.RelevanceReason
			log.Printf("[ESTÁGIO 3] Cache HIT na extração para %s (Score: %d).", ref.Key(), tRepo.RelevanceScore)
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

		report.Results = append(report.Results, *extracted)
	}

	report.ExtractedCount = len(report.Results)

	// ==========================================
	// Estágio 4: Recursividade (Malha de Dependências)
	// ==========================================
	if filter.Depth > 1 {
		log.Printf("[ESTÁGIO 4] Aprofundando busca recursiva (Restam %d rodadas) para query '%s'...", filter.Depth-1, filter.Query)
		
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
			if err == nil && childReport != nil {
				report.Results = append(report.Results, childReport.Results...)
				report.TriagedCount += childReport.TriagedCount
				report.DiscardedCount += childReport.DiscardedCount
				report.TopCandidates += childReport.TopCandidates
			}
		}
		report.ExtractedCount = len(report.Results)
	}

	// Extrair Necessidades Ocultas & Insights Arquiteturais dos repositórios validados
	insights := parser.ExtractArchitecturalInsights(report.Results)
	report.Insights = &insights

	log.Printf("[CONCLUÍDO] Funil executado com sucesso! %d repositórios entregues no relatório com insights computados.", report.ExtractedCount)

	return report, nil
}
